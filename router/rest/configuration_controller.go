package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luminous-gsm/fusion/config"
)

type ConfigurationController struct{}

func (cc ConfigurationController) GetConfiguration(c *gin.Context) {
	config := config.Get()
	c.JSON(http.StatusOK, config)
}
