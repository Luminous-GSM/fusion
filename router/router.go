package router

import (
	"github.com/luminous-gsm/fusion/router/middleware"
	"github.com/luminous-gsm/fusion/server"

	"github.com/gin-gonic/gin"
)

func NewRouter(mgr *server.ServerManager) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.AttachRequestID(), middleware.CaptureErrors(), middleware.SetAccessControlHeaders())
	router.Use(middleware.AttachServerManager(mgr))
	router.Use(middleware.AdvancedLogging())

	health := new(HealthController)

	// The following routes require no authorization.
	router.GET("/health", health.Status)

	// The following routes require authorization
	router.Use(middleware.RequireAuthorization())
	{
		agentGroup := router.Group("agent")
		{
			agent := new(AgentController)
			agentGroup.GET("/ping", agent.PingAgent)
		}
		configurationGroup := router.Group("configuration")
		{
			configuration := new(ConfigurationController)
			configurationGroup.GET("/", configuration.Get)
		}
		podGroup := router.Group("pods")
		{
			pod := new(PodController)
			podGroup.GET("/", pod.ListPods)
			podGroup.POST("/create", pod.CreatePod)
			podGroup.POST("/stop", pod.StopPod)
			podGroup.POST("/remove", pod.RemovePod)
			podGroup.POST("/start", pod.StartPod)
			podGroup.GET("/logs/:containerId", pod.GetLogsPod)
		}
	}

	return router
}
