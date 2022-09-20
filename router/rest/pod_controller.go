package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luminous-gsm/fusion/docker"
	"github.com/luminous-gsm/fusion/model/request"
	"github.com/luminous-gsm/fusion/utils"
)

type PodController struct{}

func (p PodController) ListPods(c *gin.Context) {
	containers, err := docker.Instance().ListContainers([]string{})
	if err != nil {
		NewError(err).Abort(c)
		return
	}

	c.JSON(http.StatusOK, containers)

}

func (p PodController) CreatePod(c *gin.Context) {

	var podCreateRequest request.PodCreateRequest
	if err := utils.BindAndValidate(c, &podCreateRequest); err != nil {
		NewError(err).SetMessage("Bind or data struct validation error. See server logs").AbortWithStatus(c, http.StatusBadRequest)
		return
	}

	id, err := docker.Instance().CreateContainer(podCreateRequest)
	if err != nil {
		NewError(err).SetMessage("Could not create container. See server logs").Abort(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"containerId": id})

}

func (p PodController) StartPod(c *gin.Context) {

	var podStartRequest request.PodStartRequest
	if err := utils.BindAndValidate(c, &podStartRequest); err != nil {
		NewError(err).SetMessage("Bind or data struct validation error. See server logs").AbortWithStatus(c, http.StatusBadRequest)
		return
	}

	id, err := docker.Instance().StartContainer(podStartRequest)
	if err != nil {
		NewError(err).SetMessage("Could not start container. See server logs").Abort(c)
		return
	}

	containers, err := docker.Instance().ListContainers([]string{id})
	if err != nil {
		NewError(err).SetMessage("Could not start container. See server logs").Abort(c)
		return
	}
	if len(containers) != 1 {
		NewError(err).SetMessage("Could not get container based on id. See server logs").Abort(c)
		return
	}

	c.JSON(http.StatusOK, containers[0])

}

func (p PodController) StopPod(c *gin.Context) {

	var podStopRequest request.PodStopRequest
	if err := utils.BindAndValidate(c, &podStopRequest); err != nil {
		NewError(err).SetMessage("Bind or data struct validation error. See server logs").AbortWithStatus(c, http.StatusBadRequest)
		return
	}

	id, err := docker.Instance().StopContainer(podStopRequest)
	if err != nil {
		NewError(err).SetMessage("Could not stop container. See server logs").Abort(c)
		return
	}

	containers, err := docker.Instance().ListContainers([]string{id})
	if err != nil {
		NewError(err).SetMessage("Could not start container. See server logs").Abort(c)
		return
	}
	if len(containers) != 1 {
		NewError(err).SetMessage("Could not get container based on id. See server logs").Abort(c)
		return
	}

	c.JSON(http.StatusOK, containers[0])

}

func (p PodController) RemovePod(c *gin.Context) {

	var podRemoveRequest request.PodRemoveRequest
	if err := utils.BindAndValidate(c, &podRemoveRequest); err != nil {
		NewError(err).SetMessage("Bind or data struct validation error. See server logs").AbortWithStatus(c, http.StatusBadRequest)
		return
	}

	id, err := docker.Instance().RemoveContainer(podRemoveRequest)
	if err != nil {
		NewError(err).SetMessage("Could not remove container. See server logs").Abort(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"containerId": id})

}

func (p PodController) GetLogsPod(c *gin.Context) {
	containerId := c.Param("containerId")

	tail, ok := c.GetQuery("lines")
	if !ok {
		tail = "100"
	}

	logs, err := docker.Instance().GetLogs(containerId, tail)
	if err != nil {
		NewError(err).SetMessage("Could not get container logs. See server logs").Abort(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}
