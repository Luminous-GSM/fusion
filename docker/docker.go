package docker

import (
	"sync"

	"emperror.dev/errors"
	"github.com/apex/log"
	"github.com/docker/docker/client"
)

var (
	_conce sync.Once
)

type DockerService struct {
	client *client.Client
}

// Returns a docker client.
func NewDocker() (*DockerService, error) {
	var err error
	var dockerClient *client.Client
	_conce.Do(func() {
		dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	})
	log.Info("docker: configured docker client")
	return &DockerService{
		client: dockerClient,
	}, errors.Wrap(err, "docker: could not create client")
}
