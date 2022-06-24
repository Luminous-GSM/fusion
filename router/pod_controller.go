package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luminous-gsm/fusion/router/middleware"
)

type PodController struct{}

func (p PodController) ListPods(c *gin.Context) {
	s := middleware.GetServerManager(c)

	containers, err := s.Environment().DockerService().ListContainers(s.Context())
	if err != nil {
		NewError(err).Abort(c)
		return
	}

	c.JSON(http.StatusOK, containers)

}
