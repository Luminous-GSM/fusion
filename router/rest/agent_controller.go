package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luminous-gsm/fusion/docker"
	"github.com/luminous-gsm/fusion/event"
	"github.com/luminous-gsm/fusion/model/response"
	"github.com/luminous-gsm/fusion/node"
	"github.com/luminous-gsm/fusion/utils"
	"go.uber.org/zap"
)

type AgentController struct{}

func (agent AgentController) PingAgent(c *gin.Context) {

	nodeWarnings := node.Instance().GetNodeWarnings()
	nodeDescription := node.Instance().GetNodeDescription()

	containers, err := docker.Instance().ListContainers([]string{}, []string{})
	if err != nil {
		nodeDescription.ActivePods = 0
		zap.S().With("controller", "AgentController").Errorw("Could not list docker containers", "error", err)
	} else {
		nodeDescription.ActivePods = len(containers)
	}

	nodeDescription.Warnings = nodeWarnings

	nodeResponse := &response.NodeDescriptionResponse{
		NodeDescriptionModel: nodeDescription,
	}

	c.JSON(http.StatusOK, nodeResponse)
}

func (agent AgentController) Dashboard(c *gin.Context) {

	images, err := docker.Instance().GetImages()
	if err != nil {
		NewError(err).SetMessage("Could not get images. See server logs").Abort(c)
		return
	}

	containers, err := docker.Instance().ListContainers([]string{}, []string{})
	if err != nil {
		NewError(err).SetMessage("Could not get containers. See server logs").Abort(c)
		return
	}

	description := node.Instance().GetNodeDescription()

	dashboardResponse := &response.DashboardResponse{
		NodeDescription: description,
		Images:          images,
		Pods:            containers,
	}

	c.JSON(http.StatusOK, dashboardResponse)

}

func (agent AgentController) GetSystemLoad(c *gin.Context) {

	systemLoad, err := node.Instance().GetSystemLoad()
	if err != nil {
		NewError(err).SetMessage("Could not get containers. See server logs").Abort(c)
		return
	}

	c.JSON(http.StatusOK, systemLoad)
}

func (agent AgentController) TemporaryAuthentication(c *gin.Context) {

	token, err := node.Instance().TemporaryAuthentication()
	if err != nil {
		NewError(err).SetMessage("Could not generate temporary authentication token").Abort(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (agent AgentController) PublishManualEvent(c *gin.Context) {

	var publishManualEventRequest event.PublishManualEventRequest
	if err := utils.BindAndValidate(c, &publishManualEventRequest); err != nil {
		NewError(err).SetMessage("Bind or data struct validation error. See server logs").AbortWithStatus(c, http.StatusBadRequest)
		return
	}

	switch publishManualEventRequest.Topic {
	case event.EVENT_DOCKER_POD_CREATE:
		event.Instance().Bus().RequestStream(event.EVENT_DOCKER_POD_CREATE, event.FusionEvent[event.FusionDockerEventData]{})
	case event.EVENT_REQUEST_POD_CREATE:
		{
			zap.S().Debugw("sending manual request to channel", "channel", event.EVENT_REQUEST_POD_CREATE)
			req := make(map[string]interface{})
			req["Id"] = "nod"
			req["Name"] = "pod_name"
			event.FireEvent(
				event.EVENT_REQUEST_POD_CREATE, event.FusionEvent[map[string]interface{}]{
					Entity: []*string{},
					Event:  event.EVENT_REQUEST_POD_CREATE,
					Data:   req,
				},
			)
		}
	default:
		zap.S().Debugw("unknown topic", "topicRequest", publishManualEventRequest.Topic)
	}

	c.JSON(http.StatusOK, gin.H{"status": "done"})
}

func (agent AgentController) GetAllocatedPorts(c *gin.Context) {
	containers, err := docker.Instance().ListContainers([]string{}, []string{})
	if err != nil {
		NewError(err).SetMessage("listing containers error. See server logs").AbortWithStatus(c, http.StatusBadRequest)
		return
	}

	var ports []string = make([]string, 0)

	for _, container := range containers {
		for _, port := range container.Ports {
			ports = append(ports, port.PrivatePort)
		}

	}

	c.JSON(http.StatusOK, gin.H{"ports": ports})

}
