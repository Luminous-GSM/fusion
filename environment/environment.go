package environment

import (
	"emperror.dev/errors"
	"github.com/luminous-gsm/fusion/docker"
)

type Environment struct {
	dockerService *docker.DockerService
}

func NewEnvironment() (*Environment, error) {
	dkr, err := docker.NewDocker()
	if err != nil {
		return nil, errors.Wrap(err, "environment: could not create environment")
	}

	return &Environment{
		dockerService: dkr,
	}, nil
}

func (e Environment) DockerService() *docker.DockerService {
	return e.dockerService
}
