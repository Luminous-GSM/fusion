package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luminous-gsm/fusion/model/response"
	"github.com/luminous-gsm/fusion/router/middleware"
)

type AgentController struct{}

func (agent AgentController) PingAgent(c *gin.Context) {
	s := middleware.GetServerManager(c)

	containers, err := s.ServiceManager().DockerService().ListContainers([]string{})
	if err != nil {
		NewError(err).SetMessage("Could not get container count. See server logs").Abort(c)
		return
	}

	nodeDescription := s.ServiceManager().NodeService().GetNodeDescription()
	nodeDescription.ActivePods = len(containers)

	nodeResponse := &response.NodeDescriptionResponse{
		NodeDescriptionModel: nodeDescription,
	}

	c.JSON(http.StatusOK, nodeResponse)
}

func (agent AgentController) Dashboard(c *gin.Context) {
	s := middleware.GetServerManager(c)

	images, err := s.ServiceManager().DockerService().GetImages()
	if err != nil {
		NewError(err).SetMessage("Could not get images. See server logs").Abort(c)
		return
	}

	containers, err := s.ServiceManager().DockerService().ListContainers([]string{})
	if err != nil {
		NewError(err).SetMessage("Could not get containers. See server logs").Abort(c)
		return
	}

	description := s.ServiceManager().NodeService().GetNodeDescription()

	dashboardResponse := &response.DashboardResponse{
		NodeDescription: description,
		Images:          images,
		Pods:            containers,
	}

	c.JSON(http.StatusOK, dashboardResponse)

}

func (agent AgentController) GetSystemLoad(c *gin.Context) {
	s := middleware.GetServerManager(c)

	systemLoad, err := s.ServiceManager().NodeService().GetSystemLoad()
	if err != nil {
		NewError(err).SetMessage("Could not get containers. See server logs").Abort(c)
		return
	}

	c.JSON(http.StatusOK, systemLoad)
}

func (agent AgentController) TemporaryAuthentication(c *gin.Context) {
	s := middleware.GetServerManager(c)

	token, err := s.ServiceManager().NodeService().TemporaryAuthentication()
	if err != nil {
		NewError(err).SetMessage("Could not generate temporary authentication token").Abort(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
