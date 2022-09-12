package event

import "github.com/luminous-gsm/fusion/model"

const (
	EVENT_REQUEST_POD_CREATE = "event:request:pod:create"
	EVENT_DOCKER_POD_CREATE  = "event:docker:pod:create"

	OPERATION_CONTAINER_CREATE_START  = "container-create-start"
	OPERATION_CONTAINER_CREATE_FINISH = "container-create-finish"
	OPERATION_IMAGE_DOWNLOAD_START    = "image-download-start"
	OPERATION_IMAGE_DOWNLOAD_PROGRESS = "image-download-progress"
	OPERATION_IMAGE_DOWNLOAD_FINISH   = "image-download-finish"
)

type FusionEventData interface {
	model.PodDescription | FusionDockerEventData | map[string]interface{}
}

type FusionEvent[T FusionEventData] struct {
	Entity []*string `json:"entity"`
	Event  string    `json:"event"`
	Data   T         `json:"data"`
}

type FusionDockerEventData struct {
	Operation string `json:"operation"`
	Message   string `json:"message"`
}

type PublishManualEventRequest struct {
	Topic string `json:"topic"`
}
