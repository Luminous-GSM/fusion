package rest

import (
	"github.com/luminous-gsm/fusion/config"
	"github.com/luminous-gsm/fusion/router/middleware"
	"github.com/luminous-gsm/fusion/server"

	"github.com/gin-gonic/gin"
)

func NewRouter(mgr *server.ServerManager) *gin.Engine {
	if config.Get().Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.AttachRequestID(), middleware.CaptureErrors(), middleware.SetAccessControlHeaders())
	router.Use(middleware.AttachServerManager(mgr))
	router.Use(middleware.AdvancedLogging())
	router.Use(middleware.LocalDelay())

	health := new(HealthController)

	// The following routes require no authorization.
	router.GET("/health", health.Status)

	tempAuthGroup := router.Group("ws", middleware.RequireTemporaryAuthorization())
	{
		tempAuthGroup.GET("/node", RunWebSocket)
	}

	// The following routes require authorization
	router.Use(middleware.RequireAuthorization())
	{
		agentGroup := router.Group("agent")
		{
			agent := new(AgentController)
			agentGroup.GET("/ping", agent.PingAgent)
			agentGroup.GET("/dashboard", agent.Dashboard)
			agentGroup.GET("/system-load", agent.GetSystemLoad)
			agentGroup.GET("/temp-auth", agent.TemporaryAuthentication)
			agentGroup.POST("/manual-event", agent.PublishManualEvent)
			agentGroup.GET("/allocated-ports", agent.GetAllocatedPorts)
		}
		configurationGroup := router.Group("configuration")
		{
			configuration := new(ConfigurationController)
			configurationGroup.GET("/", configuration.GetConfiguration)
		}
		podGroup := router.Group("pods")
		{
			pod := new(PodController)
			podGroup.GET("/", pod.ListPods)
			podGroup.GET("/inspect/:containerId", pod.InfoPod)
			podGroup.POST("/create", pod.CreatePod)
			podGroup.POST("/stop", pod.StopPod)
			podGroup.POST("/remove", pod.RemovePod)
			podGroup.POST("/start", pod.StartPod)
			podGroup.GET("/logs/:containerId", pod.GetLogsPod)
		}
	}

	return router
}
