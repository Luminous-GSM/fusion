package docker

import (
	"context"
	"sync"

	"github.com/docker/docker/client"
	"github.com/luminous-gsm/fusion/event"
	"github.com/luminous-gsm/fusion/model"
	"github.com/luminous-gsm/fusion/model/request"
	eventModel "github.com/vmware/transport-go/model"
	"go.uber.org/zap"
)

var (
	_conce        sync.Once
	dockerService *DockerService
)

type DockerService struct {
	ctx    context.Context
	client *client.Client
}

// Returns a docker client.
func InitDockerService(context context.Context) *DockerService {
	_conce.Do(func() {
		dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			zap.S().Named("docker").DPanicw("could not create client", "error", err)
		}
		dockerService = &DockerService{
			client: dockerClient,
			ctx:    context,
		}
		zap.S().Info("docker: configured docker client")
	})

	return dockerService
}

func (ds DockerService) InitEventListeners() {

	handler, err := event.Instance().Bus().ListenRequestStream(event.EVENT_REQUEST_POD_CREATE)
	if err != nil {
		ds.log().Panicw("could not create handler for channel", "error", err, "channel", event.EVENT_REQUEST_POD_CREATE)
	}
	handler.Handle(
		func(m *eventModel.Message) {
			ds.log().Debugw("new event received", "message", m, "event", event.EVENT_REQUEST_POD_CREATE)
			fusionEvent := m.Payload.(event.FusionEvent[map[string]interface{}])
			var podDescription model.PodDescription
			err = event.UnmarshalUnknown(fusionEvent.Data, &podDescription)
			ds.log().Debugw("marshalled pod description", "podDescription", podDescription)
			ds.CreateContainer(request.PodCreateRequest{
				PodDescription: podDescription,
			})
		},
		event.Instance().DefaultErrorHandler,
	)

}

func Instance() *DockerService {
	return dockerService
}

func (ds DockerService) log() *zap.SugaredLogger {
	return zap.S().Named("docker")
}
