package docker

import (
	"context"

	"emperror.dev/errors"
	"github.com/docker/docker/api/types"
)

func (ds DockerService) ListContainers(context context.Context) ([]types.Container, error) {
	// TODO List only fusion pods. Probably filter on labels.
	options := types.ContainerListOptions{
		All: true,
	}
	containers, err := ds.client.ContainerList(context, options)

	return containers, errors.Wrap(err, "docker: could not list containers")
}
