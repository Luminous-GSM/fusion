package router

import (
	"github.com/luminous-gsm/fusion/controllers"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(AttachRequestID(), CaptureErrors(), SetAccessControlHeaders())
	router.Use(AdvancedLogging())

	health := new(controllers.HealthController)

	// The following routes require no authorization.
	router.GET("/health", health.Status)

	// The following routes require authorization
	router.Use(RequireAuthorization())
	{
		configurationGroup := router.Group("configuration")
		{
			configuration := new(controllers.ConfigurationController)
			configurationGroup.GET("/", configuration.Get)
		}
	}

	return router
}
