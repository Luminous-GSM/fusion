package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luminous-gsm/fusion/config"
	"github.com/luminous-gsm/fusion/model/response"
	"github.com/luminous-gsm/fusion/router/middleware"
)

type AgentController struct{}

func (agent AgentController) PingAgent(c *gin.Context) {
	s := middleware.GetServerManager(c)

	config := config.Get()

	containers, err := s.ServiceManager().DockerService().ListContainers()
	if err != nil {
		NewError(err).SetMessage("Could not get container count. See server logs").Abort(c)
		return
	}

	nodeResponse := &response.NodeDescriptionResponse{
		Ip:              "0.0.0.0",
		NodeUniqueId:    config.NodeUniqueId,
		Name:            config.NodeName,
		Description:     config.NodeDescription,
		NodeStatus:      "running",
		Version:         config.Version,
		HostingPlatform: response.HostingPlatformType(config.HostingPlatform),
		ActivePods:      len(containers),
		Token:           config.ApiSecurityToken,
	}

	c.JSON(http.StatusOK, nodeResponse)
}
