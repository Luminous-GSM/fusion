package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luminous-gsm/fusion/model/request"
	"github.com/luminous-gsm/fusion/router/middleware"
)

type PodController struct{}

func (p PodController) ListPods(c *gin.Context) {
	s := middleware.GetServerManager(c)

	containers, err := s.ServiceManager().DockerService().ListContainers()
	if err != nil {
		NewError(err).Abort(c)
		return
	}

	c.JSON(http.StatusOK, containers)

}

func (p PodController) CreatePod(c *gin.Context) {
	s := middleware.GetServerManager(c)

	var podCreateRequest request.PodCreateRequest
	if err := s.BindAndValidate(c, &podCreateRequest); err != nil {
		NewError(err).SetMessage("Bind or data struct validation error. See server logs").AbortWithStatus(c, http.StatusBadRequest)
		return
	}

	id, err := s.ServiceManager().DockerService().CreateContainer(podCreateRequest)
	if err != nil {
		NewError(err).SetMessage("Could not create container. See server logs").Abort(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"containerId": id})

}

func (p PodController) StartPod(c *gin.Context) {
	s := middleware.GetServerManager(c)

	var podStartRequest request.PodStartRequest
	if err := s.BindAndValidate(c, &podStartRequest); err != nil {
		NewError(err).SetMessage("Bind or data struct validation error. See server logs").AbortWithStatus(c, http.StatusBadRequest)
		return
	}

	id, err := s.ServiceManager().DockerService().StartContainer(podStartRequest)
	if err != nil {
		NewError(err).SetMessage("Could not start container. See server logs").Abort(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"containerId": id})

}

func (p PodController) StopPod(c *gin.Context) {
	s := middleware.GetServerManager(c)

	var podStopRequest request.PodStopRequest
	if err := s.BindAndValidate(c, &podStopRequest); err != nil {
		NewError(err).SetMessage("Bind or data struct validation error. See server logs").AbortWithStatus(c, http.StatusBadRequest)
		return
	}

	id, err := s.ServiceManager().DockerService().StopContainer(podStopRequest)
	if err != nil {
		NewError(err).SetMessage("Could not stop container. See server logs").Abort(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"containerId": id})

}

func (p PodController) RemovePod(c *gin.Context) {
	s := middleware.GetServerManager(c)

	var podRemoveRequest request.PodRemoveRequest
	if err := s.BindAndValidate(c, &podRemoveRequest); err != nil {
		NewError(err).SetMessage("Bind or data struct validation error. See server logs").AbortWithStatus(c, http.StatusBadRequest)
		return
	}

	id, err := s.ServiceManager().DockerService().RemoveContainer(podRemoveRequest)
	if err != nil {
		NewError(err).SetMessage("Could not remove container. See server logs").Abort(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"containerId": id})

}
