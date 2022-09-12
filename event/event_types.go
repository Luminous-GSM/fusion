package event

import "github.com/luminous-gsm/fusion/model"

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
