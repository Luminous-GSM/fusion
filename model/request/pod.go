package request

import "github.com/luminous-gsm/fusion/model"

type ContainerId struct {
	ContainerId string `validate:"required" json:"containerId"`
}

type PodCreateRequest struct {
	PodDescription model.PodDescription `validate:"required" json:"podDescription"`
}

type PodStartRequest struct {
	ContainerId `validate:"required" json:"container"`
}

type PodStopRequest struct {
	ContainerId `validate:"required" json:"container"`
}

type PodRemoveRequest struct {
	ContainerId `validate:"required" json:"container"`
}
