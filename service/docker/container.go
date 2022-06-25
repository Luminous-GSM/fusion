package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"time"

	"emperror.dev/errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/daemon/logger/local"
	"github.com/docker/go-connections/nat"
	"github.com/luminous-gsm/fusion/config"
	"github.com/luminous-gsm/fusion/model"
	"github.com/luminous-gsm/fusion/model/request"
	"go.uber.org/zap"
)

const (
	containerPullTimeout    = 15
	containerCreateTimeout  = 5
	containerStartTimeout   = 5
	containerStopTimeout    = 5
	containerRemoveTimmeout = 5
	isFusionManaged         = "is-fusion-managed"
	manifestFileUsed        = "manifest-file-used"
	friendlyName            = "friendly-name"
)

func (ds DockerService) ListContainers() ([]types.Container, error) {
	ctx, cancel := context.WithCancel(ds.ctx)
	defer cancel()
	// TODO List only fusion pods. Probably filter on labels.
	options := types.ContainerListOptions{
		All: true,
	}
	containers, err := ds.client.ContainerList(ctx, options)

	return containers, errors.Wrap(err, "docker: could not list containers")
}

// Create the container for the specific pod
func (ds DockerService) CreateContainer(podCreateRequest request.PodCreateRequest) (string, error) {
	imageRef := podCreateRequest.PodDescription.Image + ":" + podCreateRequest.PodDescription.Tag

	zap.S().Infow("creating container", "image", imageRef)

	if err := ds.ensureImageExists(imageRef); err != nil {
		return "", err
	}

	podUniqueId := generateUniqueFolderName(podCreateRequest.PodDescription)

	// Cancel after containerPullTimeout of time
	ctx, cancel := context.WithTimeout(ds.ctx, time.Minute*containerCreateTimeout)
	defer cancel()

	exposed, bindings, err := getBindsFromPortMaps(podCreateRequest.PodDescription.PortMaps)
	if err != nil {
		return "", err
	}

	containerConfig := &container.Config{
		Hostname:     config.Get().Node.Hostname,
		Domainname:   "",
		User:         strconv.Itoa(config.Get().System.User.Uid) + ":" + strconv.Itoa(config.Get().System.User.Gid),
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		Tty:          true,
		ExposedPorts: exposed,
		Image:        imageRef,
		Env:          getEnvironmentVariablesFromMaps(podCreateRequest.PodDescription.EnvironmentMaps),
		Labels: map[string]string{
			manifestFileUsed: podCreateRequest.PodDescription.ManifestFileUsed,
			isFusionManaged:  "true",
			friendlyName:     podCreateRequest.PodDescription.Name,
		},
	}

	zap.S().Debugw("container configuration", "containerConfiguration", containerConfig)

	tmpfsSize := config.Get().Pod.TmpfsSize

	hostConfig := &container.HostConfig{
		PortBindings: bindings,

		Mounts: getMountsFromMountMaps(podCreateRequest.PodDescription, podUniqueId),

		DNS: config.Get().Pod.Dns,

		// Temp storage for server downloadable assets
		Tmpfs: map[string]string{
			"/tmp": "rw,exec,nosuid,size=" + tmpfsSize + "M",
		},

		LogConfig: container.LogConfig{
			Type: local.Name,
			Config: map[string]string{
				"max-size": "5m",
				"max-file": "1",
				"compress": "false",
				"mode":     "non-blocking",
			},
		},
		Resources: container.Resources{
			Memory: int64(podCreateRequest.PodDescription.Limit.Memory * 1_000_000),
		},

		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
		SecurityOpt:    []string{"no-new-privileges"},
		ReadonlyRootfs: true,
	}

	zap.S().Debugw("host configuration", "hostConfiguration", hostConfig)

	result, err := ds.client.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, podUniqueId)
	if err != nil {
		zap.S().Errorw("error creating container", "error", err)
		return "", err
	}

	for _, warning := range result.Warnings {
		zap.S().Warnw("creating container completed, but there were some warnings",
			"warning", warning,
			"image", imageRef,
		)
	}

	return result.ID, nil
}

// Starts the specific container
// Request container ID
func (ds DockerService) StartContainer(podStartRequest request.PodStartRequest) (string, error) {

	// Cancel after containerPullTimeout of time
	ctx, cancel := context.WithTimeout(ds.ctx, time.Minute*containerStartTimeout)
	defer cancel()

	if err := ds.client.ContainerStart(ctx, podStartRequest.ContainerId.ContainerId, types.ContainerStartOptions{}); err != nil {
		zap.S().Errorw("could not start container", "podStartRequest", podStartRequest)
		return "", err
	}

	return podStartRequest.ContainerId.ContainerId, nil

}

// Stops the specific container
// Request container ID
func (ds DockerService) StopContainer(podStopRequest request.PodStopRequest) (string, error) {

	// Cancel after containerPullTimeout of time
	ctx, cancel := context.WithTimeout(ds.ctx, time.Minute*containerStopTimeout)
	defer cancel()

	duration, err := time.ParseDuration("30s")
	if err != nil {
		zap.S().Errorw("could not parse duration", "duration", "30s")
		return "", err
	}

	if err := ds.client.ContainerStop(ctx, podStopRequest.ContainerId.ContainerId, &duration); err != nil {
		zap.S().Errorw("could not stop container", "PodStopRequest", podStopRequest)
		return "", err
	}

	return podStopRequest.ContainerId.ContainerId, nil

}

// Remove the specific container
// Request container ID
func (ds DockerService) RemoveContainer(podRemoveRequest request.PodRemoveRequest) (string, error) {

	// Cancel after containerPullTimeout of time
	ctx, cancel := context.WithTimeout(ds.ctx, time.Minute*containerRemoveTimmeout)
	defer cancel()

	if err := ds.client.ContainerRemove(ctx, podRemoveRequest.ContainerId.ContainerId, types.ContainerRemoveOptions{Force: true}); err != nil {
		zap.S().Errorw("could not remove container", "PodRemoveRequest", podRemoveRequest)
		return "", err
	}

	return podRemoveRequest.ContainerId.ContainerId, nil

}

func getBindsFromPortMaps(ports []model.PortMap) (map[nat.Port]struct{}, map[nat.Port][]nat.PortBinding, error) {
	var dockerStandardPorts []string
	for _, port := range ports {
		dockerStandardPorts = append(dockerStandardPorts, fmt.Sprintf("127.0.0.1:%v:%v/%v", port.Exposed, port.Binding, port.Protocol))
		zap.S().Debugw("preparing bings from port maps", "port", port)
	}
	exposed, bindings, err := nat.ParsePortSpecs(dockerStandardPorts)
	if err != nil {
		zap.S().Errorw("error parsing port specifications", "error", err)
		return nil, nil, err
	}
	return exposed, bindings, nil
}

func getEnvironmentVariablesFromMaps(environmentVariables []model.EnvironmentMap) []string {
	var dockerStandardEnvironmentVariables []string
	for _, env := range environmentVariables {
		dockerStandardEnvironmentVariables = append(dockerStandardEnvironmentVariables, fmt.Sprintf("%s=%s", env.Name, env.Value))
		zap.S().Debugw("preparing environment from environment maps", "environment", env)
	}
	return dockerStandardEnvironmentVariables
}

func generateUniqueFolderName(description model.PodDescription) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return description.Name + "-" + string(b)
}

func getMountsFromMountMaps(description model.PodDescription, podUniqueId string) []mount.Mount {
	var dockerStandardMountVariables []mount.Mount
	for _, mountItem := range description.MountMaps {
		dockerStandardMountVariables = append(dockerStandardMountVariables, mount.Mount{
			Type:     mount.TypeBind,
			Source:   config.Get().System.DataDirectory + podUniqueId + mountItem.Destination,
			Target:   mountItem.Destination,
			ReadOnly: false,
		})
		zap.S().Debugw("preparing mount from mount maps", "mount", mountItem)
	}
	return dockerStandardMountVariables
}

// See if the image already exist.
// If the image does not exist, pull it.
// If the image does exist locally, do nothing.
func (ds DockerService) ensureImageExists(imageRef string) error {

	// Cancel after containerPullTimeout of time
	ctx, cancel := context.WithTimeout(ds.ctx, time.Minute*containerPullTimeout)
	defer cancel()

	zap.S().Info("ensuring that the image exist")

	// Try and pull the image from the registry
	pullOptions := types.ImagePullOptions{All: false}
	out, err := ds.client.ImagePull(ds.ctx, imageRef, pullOptions)
	if err != nil {
		// Image pull did not succeed
		images, ierr := ds.client.ImageList(ctx, types.ImageListOptions{})
		if ierr != nil {
			return errors.Wrap(ierr, "docker: failed to list images")
		}

		for _, img := range images {
			for _, t := range img.RepoTags {
				if t != imageRef {
					continue
				}

				zap.S().Warnw("unable to pull requested image from remote source, however the image exists locally",
					"image", imageRef,
					"error", err.Error(),
				)

				// Matching container image found locally.
				// Return from this fuction as the image is still available for use
				return nil
			}
		}

		return errors.Wrapf(err, "docker: failed to pull \"%s\" image for server", imageRef)
	}
	defer out.Close()

	zap.S().Debugw("pulling docker images", "image", imageRef)

	d := json.NewDecoder(out)

	type Event struct {
		Status         string `json:"status"`
		Error          string `json:"error"`
		Progress       string `json:"progress"`
		ProgressDetail struct {
			Current int `json:"current"`
			Total   int `json:"total"`
		} `json:"progressDetail"`
	}

	var event *Event
	for {
		if err := d.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}

			panic(err)
		}

		zap.S().Debugf("docker event: %+v\n", event)
	}

	zap.S().Debugw("complete docker image pull", "image", imageRef)

	return nil

}
