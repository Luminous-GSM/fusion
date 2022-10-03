package docker

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"emperror.dev/errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/daemon/logger/local"
	"github.com/docker/go-connections/nat"
	"github.com/kennygrant/sanitize"
	"github.com/luminous-gsm/fusion/config"
	"github.com/luminous-gsm/fusion/event"
	"github.com/luminous-gsm/fusion/model"
	"github.com/luminous-gsm/fusion/model/domain"
	"github.com/luminous-gsm/fusion/model/request"
	"github.com/luminous-gsm/fusion/utils"
	"github.com/luminous-gsm/fusion/variables"
	"go.uber.org/zap"
)

const (
	containerPullTimeout    = 15
	containerCreateTimeout  = 5
	containerStartTimeout   = 5
	containerStopTimeout    = 5
	containerRemoveTimmeout = 5
)

func (ds DockerService) publishEvent(operation, message string) {
	ds.log().Debugw("publish docker container event", "message", message)
	event.FireEvent(
		event.EVENT_DOCKER_POD_CREATE,
		event.FusionEvent[event.FusionDockerEventData]{
			Entity: []*string{},
			Event:  event.EVENT_DOCKER_POD_CREATE,
			Data: event.FusionDockerEventData{
				Operation: operation,
				Message:   message,
			},
		},
	)
}

func (ds DockerService) Info() (types.Info, error) {
	info, err := ds.client.Info(ds.ctx)
	if err != nil {
		ds.log().Errorw("could not list containers", "error", err)
		return types.Info{}, err
	}

	return info, nil

}

func (ds DockerService) InspectContainer(containerId string) (domain.FusionContainerModel, error) {
	ctx, cancel := context.WithCancel(ds.ctx)
	defer cancel()
	inspect, err := ds.client.ContainerInspect(ctx, containerId)
	if err != nil {
		ds.log().Errorw("Could not inspect container", "error", err, "containerId", containerId)
		return domain.FusionContainerModel{}, err
	}

	ports := make([]domain.ContainerPort, 0)

	for key := range inspect.Config.ExposedPorts {
		ports = append(ports, domain.ContainerPort{
			Ip:          "",
			PrivatePort: key.Port(),
			PublicPort:  "0",
			Type:        key.Proto(),
		})
	}

	labels := make([]string, 0)
	for key, value := range inspect.Config.Labels {
		labels = append(labels, fmt.Sprintf("%v=%v", key, value))
	}

	return domain.FusionContainerModel{
		Id:     inspect.ID,
		Ports:  ports,
		Labels: labels,
	}, nil
}

func (ds DockerService) ListContainers(containerIds, containerTypes []string) ([]domain.FusionContainerModel, error) {
	ctx, cancel := context.WithCancel(ds.ctx)
	defer cancel()

	// Filter the containers that is managed by fusion.
	// This is important as the user might manually create and use
	// other containers, and we don't want to manage those.
	// TODO -low : We can add a query parameter to list all (including fusion non-managed) containers.
	filters := filters.NewArgs()
	filters.Add("label", "fusion_is-managed")
	for _, containerType := range containerTypes {
		filters.Add("label", fmt.Sprintf("fusion_type=%v", containerType))
	}

	for _, id := range containerIds {
		filters.Add("id", id)
	}

	options := types.ContainerListOptions{
		All:     true,
		Filters: filters,
	}
	containers, err := ds.client.ContainerList(ctx, options)
	if err != nil {
		ds.log().Errorw("could not list containers", "error", err)
		return nil, err
	}

	consoleContainers := []domain.FusionContainerModel{}

	for _, container := range containers {

		inspectedContainer, err := ds.InspectContainer(container.ID)
		if err != nil {
			return nil, err
		}

		ports := []domain.ContainerPort{}

		if len(container.Ports) == 0 {
			ports = append(ports, inspectedContainer.Ports...)
		} else {
			for _, port := range container.Ports {
				ports = append(ports, domain.ContainerPort{
					Ip:          port.IP,
					PrivatePort: strconv.FormatUint(uint64(port.PrivatePort), 10),
					PublicPort:  strconv.FormatUint(uint64(port.PublicPort), 10),
					Type:        port.Type,
				})
			}
		}

		consoleContainers = append(consoleContainers, domain.FusionContainerModel{
			Id:      container.ID,
			Command: container.Command,
			Created: int(container.Created),
			Image:   container.Image,
			ImageId: container.ImageID,
			Names:   container.Names,
			Status:  container.Status,
			State:   domain.FusionContainerState(container.State),
			Ports:   ports,
			Labels:  inspectedContainer.Labels,
		})
	}

	// ds.client.ContainerInspectWithRaw()

	zap.S().Debugw("docker: got containers")

	return consoleContainers, nil
}

// Create the container for the specific pod
func (ds DockerService) CreateContainer(podCreateRequest request.PodCreateRequest) (string, error) {
	ds.log().Infow("Creating container", "podCreateRequest", podCreateRequest)

	ds.publishEvent(event.OPERATION_CONTAINER_CREATE_START, "Starting Fusion Pod Creation")
	defer ds.publishEvent(event.OPERATION_CONTAINER_CREATE_FINISH, "Finished Fusion Pod Creation")

	// Create a secure random string for pod suffix, as well as the scope for fusion variables.
	podIdSuffix, err := utils.GenerateSecureRandomString(5, false)
	if err != nil {
		return "", err
	}
	variableScope := podIdSuffix

	// Make sure to remove the current scope from variables.
	// Also refersh the generated variables for fresh random values
	defer variables.Instance().RemoveScopedVariables(variableScope)
	variables.Instance().RefreshGeneratedVariables()

	// Cancel after containerPullTimeout of time
	ctx, cancel := context.WithTimeout(ds.ctx, time.Minute*containerCreateTimeout)
	defer cancel()

	imageRef := podCreateRequest.PodDescription.Image + ":" + podCreateRequest.PodDescription.Tag
	ds.log().Debugw("creating container", "image", imageRef)

	containerName := podCreateRequest.PodDescription.ManifestInfo.Id + "_" + sanitize.BaseName(podCreateRequest.PodDescription.Name)
	podId := containerName + "_" + podIdSuffix
	ds.log().Debugw("creating pod id", "podId", podId, "containerName", containerName)

	// Add pod ID to variables that can be used via fusion variables
	variables.Instance().AddVariableToScopedReplacer(variableScope, "fusion.pod.id", podId)

	if err := ds.ensureImageExists(imageRef); err != nil {
		return "", err
	}

	bindings, exposed := getPortBindsFromPortMaps(podCreateRequest.PodDescription.PortMaps)
	if err != nil {
		return "", err
	}

	environmentVariables, err := getEnvironmentVariablesFromMaps(variableScope, podCreateRequest.PodDescription.EnvironmentMaps)
	if err != nil {
		return "", err
	}

	containerConfig := &container.Config{
		Hostname:   config.Get().NodeHostname,
		Domainname: "",
		// User:         strconv.Itoa(config.Get().SystemUserUid) + ":" + strconv.Itoa(config.Get().SystemUserGid),
		ExposedPorts: exposed,
		Image:        imageRef,
		Tty:          true,
		Env:          environmentVariables,
		Labels: map[string]string{
			"fusion_is-managed":    "true",
			"fusion_friendly-name": podCreateRequest.PodDescription.Name,
			"fusion_pod-id":        podId,
			"fusion_type":          podCreateRequest.PodDescription.ManifestInfo.Type,
			"fusion_manifest-id":   podCreateRequest.PodDescription.ManifestInfo.Id,
			"fusion_photon-id":     podCreateRequest.PodDescription.ManifestInfo.Id,
		},
	}

	ds.log().Debugw("container configuration", "containerConfiguration", containerConfig)

	tmpfsSize := config.Get().Pod.TmpfsSize

	mounts, err := getMountsFromMountMaps(variableScope, podCreateRequest.PodDescription)
	if err != nil {
		return "", err
	}

	hostConfig := &container.HostConfig{
		PortBindings: bindings,

		Mounts: mounts,

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
		ReadonlyRootfs: false,

		CapDrop: []string{
			"setpcap", "mknod", "audit_write", "net_raw", "dac_override",
			"fowner", "fsetid", "net_bind_service", "sys_chroot", "setfcap",
		},
	}

	ds.log().Debugw("host configuration", "hostConfiguration", hostConfig)

	err = ds.createDirectories(hostConfig.Mounts)
	if err != nil {
		return "", err
	}

	result, err := ds.client.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		ds.log().Errorw("error creating container", "error", err)
		return "", err
	}

	for _, warning := range result.Warnings {
		ds.log().Debugw("completed creating container, but there were some warnings",
			"warning", warning,
			"image", imageRef,
		)
	}

	ds.log().Infow("completed creating container",
		"imageRef", imageRef,
		"containerId", result.ID,
	)

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

func getPortBindsFromPortMaps(ports []model.PortMap) (nat.PortMap, nat.PortSet) {
	portMap := nat.PortMap{}

	for _, port := range ports {
		binding := nat.PortBinding{
			HostIP:   "0.0.0.0",
			HostPort: strconv.Itoa(port.Binding),
		}

		natPort := nat.Port(fmt.Sprintf("%d/%v", port.Binding, port.Protocol))

		portMap[natPort] = append(portMap[natPort], binding)
	}

	portSet := nat.PortSet{}

	for port := range portMap {
		portSet[port] = struct{}{}
	}

	return portMap, portSet
}

// func getBindsFromPortMaps(ports []model.PortMap) (map[nat.Port]struct{}, map[nat.Port][]nat.PortBinding, error) {
// 	var dockerStandardPorts []string
// 	for _, port := range ports {
// 		dockerStandardPorts = append(dockerStandardPorts, fmt.Sprintf("0.0.0.0:%v:%v/%v", port.Exposed, port.Binding, port.Protocol))
// 		zap.S().Debugw("preparing bings from port maps", "port", port)
// 	}
// 	exposed, bindings, err := nat.ParsePortSpecs(dockerStandardPorts)
// 	if err != nil {
// 		zap.S().Errorw("error parsing port specifications", "error", err)
// 		return nil, nil, err
// 	}
// 	return exposed, bindings, nil
// }

func getEnvironmentVariablesFromMaps(variableScope string, environmentVariables []model.EnvironmentMap) ([]string, error) {
	var dockerStandardEnvironmentVariables []string
	for _, env := range environmentVariables {
		environmentVariableValue, err := variables.Instance().ReplaceGlobalAndScopedVariablesInString(variableScope, env.Value, true)
		if err != nil {
			return nil, err
		}

		dockerStandardEnvironmentVariables = append(dockerStandardEnvironmentVariables, fmt.Sprintf("%s=%s", env.Name, environmentVariableValue))
		zap.S().Debugw("preparing environment from environment maps", "environment", env)
	}
	return dockerStandardEnvironmentVariables, nil
}

func getMountsFromMountMaps(variableScope string, description model.PodDescription) ([]mount.Mount, error) {
	var dockerStandardMountVariables []mount.Mount
	for _, mountItem := range description.MountMaps {

		sourceDirectory, err := variables.Instance().ReplaceGlobalAndScopedVariablesInString(variableScope, mountItem.Source, true)
		if err != nil {
			return nil, err
		}

		dockerStandardMountVariables = append(dockerStandardMountVariables, mount.Mount{
			Type:     mount.TypeBind,
			Source:   sourceDirectory,
			Target:   mountItem.Destination,
			ReadOnly: false,
		})
		zap.S().Debugw("preparing mount from mount maps", "mount", mountItem)
	}
	return dockerStandardMountVariables, nil
}

// See if the image already exist.
// If the image does not exist, pull it.
// If the image does exist locally, do nothing.
func (ds DockerService) ensureImageExists(imageRef string) error {
	ds.publishEvent(event.OPERATION_IMAGE_DOWNLOAD_START, "Downloading Pod Image - Starting")
	defer ds.publishEvent(event.OPERATION_IMAGE_DOWNLOAD_FINISH, "Downloading Pod Image - Completed")
	// Cancel after containerPullTimeout of time
	ctx, cancel := context.WithTimeout(ds.ctx, time.Minute*containerPullTimeout)
	defer cancel()

	ds.log().Debugw("ensuring that the image exist", "imageRef", imageRef)

	// Try and pull the image from the registry
	pullOptions := types.ImagePullOptions{All: false}
	out, err := ds.client.ImagePull(ds.ctx, imageRef, pullOptions)
	if err != nil {
		ds.log().Debugw("image pull failed, checking if image exists locally", "imageRef", imageRef)
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

				ds.log().Debugw("unable to pull requested image from remote source, however the image exists locally",
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

	ds.log().Debugw("pulling docker images", "image", imageRef)

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

	var imageEvent *Event
	for {
		if err := d.Decode(&imageEvent); err != nil {
			if err == io.EOF {
				break
			}

			panic(err)
		}

		// ds.log().Debugf("docker event: %+v\n", event)

		if imageEvent.Status == "Downloading" {
			// percentage := ((event.ProgressDetail.Current / event.ProgressDetail.Total) * 100)
			// publishEvent(eventContainerOpertaionDownoad, "Downloading Pod Image - "+strconv.Itoa(percentage)+"%")

			ds.publishEvent(event.OPERATION_IMAGE_DOWNLOAD_PROGRESS, imageEvent.Progress)
		}

	}

	ds.log().Debugw("completed docker image pull", "image", imageRef)

	return nil

}

func (ds DockerService) GetLogs(containerId, limit string) ([]string, error) {
	zap.S().Infow("started get logs", "containerId", containerId)

	// Cancel after containerPullTimeout of time
	ctx, cancel := context.WithCancel(ds.ctx)
	defer cancel()

	logOptions := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Tail:       limit,
	}

	r, err := ds.client.ContainerLogs(ctx, containerId, logOptions)
	if err != nil {
		zap.S().Errorw("could not retrieve container logs",
			"error", err,
			"containerId", containerId,
		)
		return []string{}, err
	}
	defer r.Close()

	var logs []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		logs = append(logs, scanner.Text())
	}

	zap.S().Infow("returning container logs", "containerId", containerId)
	return logs, nil

}

func (ds DockerService) GetImages() ([]domain.FusionImageModel, error) {
	ds.log().Infow("listing images")

	// Cancel after containerPullTimeout of time
	ctx, cancel := context.WithCancel(ds.ctx)
	defer cancel()

	options := types.ImageListOptions{
		All: true,
	}
	images, err := ds.client.ImageList(ctx, options)
	if err != nil {
		zap.S().Errorw("could not retrieve images",
			"error", err,
		)
		return nil, err
	}

	consoleImages := []domain.FusionImageModel{}
	for _, image := range images {
		consoleImages = append(consoleImages, domain.FusionImageModel{
			Id:         image.ID,
			Created:    int(image.Created),
			Size:       int(image.Size),
			Containers: int(image.Containers),
		})
	}

	ds.log().Infow("listing images completed")

	return consoleImages, nil

}

func (ds DockerService) createDirectories(mounts []mount.Mount) error {
	zap.S().Debugw("docker: creating directories", "directories", mounts)

	for _, mount := range mounts {
		err := os.MkdirAll(mount.Source, os.ModePerm)
		if err != nil {
			zap.S().Errorw("docker: cannot create logging directory", "error", err)
			return err
		}
	}

	zap.S().Debug("docker: created directories")
	return nil
}
