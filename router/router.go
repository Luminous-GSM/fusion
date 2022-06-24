package router

import (
	"github.com/luminous-gsm/fusion/router/middleware"
	"github.com/luminous-gsm/fusion/server"

	"github.com/gin-gonic/gin"
)

func NewRouter(m *server.ServerManager) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.AttachRequestID(), middleware.CaptureErrors(), middleware.SetAccessControlHeaders())
	router.Use(middleware.AttachServerManager(m))
	router.Use(middleware.AdvancedLogging())

	health := new(HealthController)

	// The following routes require no authorization.
	router.GET("/health", health.Status)

	// The following routes require authorization
	router.Use(middleware.RequireAuthorization())
	{
		configurationGroup := router.Group("configuration")
		{
			configuration := new(ConfigurationController)
			configurationGroup.GET("/", configuration.Get)
		}
		podGroup := router.Group("pods")
		{
			pod := new(PodController)
			podGroup.GET("/", pod.ListPods)
		}
	}

	return router
}
