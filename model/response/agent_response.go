package response

import (
	"github.com/luminous-gsm/fusion/model/domain"
)

type NodeDescriptionResponse struct {
	domain.NodeDescriptionModel
}

type DashboardResponse struct {
	NodeDescription domain.NodeDescriptionModel   `json:"nodeDescription"`
	Pods            []domain.FusionContainerModel `json:"pods"`
	Images          []domain.FusionImageModel     `json:"images"`
}
