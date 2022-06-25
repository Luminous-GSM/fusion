package docker

import (
	"context"
	"sync"

	"emperror.dev/errors"
	"github.com/docker/docker/client"
	"go.uber.org/zap"
)

var (
	_conce sync.Once
)

type DockerService struct {
	ctx    context.Context
	client *client.Client
}

// Returns a docker client.
func NewDockerService(context context.Context) (*DockerService, error) {
	var err error
	var dockerClient *client.Client
	_conce.Do(func() {
		dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	})
	if err != nil {
		return nil, errors.Wrap(err, "docker: could not create client")
	}
	zap.S().Info("docker: configured docker client")
	return &DockerService{
		client: dockerClient,
		ctx:    context,
	}, nil
}
