package service

import (
	"context"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/luminous-gsm/fusion/service/docker"
	"go.uber.org/zap"
)

type ServiceManager struct {
	dockerService *docker.DockerService
	validator     *validator.Validate
}

func NewServiceManager(ctx context.Context) (*ServiceManager, error) {
	dkr, err := docker.NewDockerService(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "environment: could not create environment")
	}

	return &ServiceManager{
		dockerService: dkr,
		validator:     validator.New(),
	}, nil
}

func (e ServiceManager) DockerService() *docker.DockerService {
	return e.dockerService
}

func (e ServiceManager) ValidateData(data interface{}) error {
	// Validate the configuration according to validation tags in the structs.
	if err := e.validator.Struct(data); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			zap.S().Errorw("request field error",
				"field", err.Field(),
				"value", err.Value(),
				"validation_type", err.Tag(),
				"field_type", err.Type(),
			)
		}
		zap.S().Debugw("request error", "error", err)
		return err
	}
	return nil
}
