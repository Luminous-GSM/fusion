package controllers

import (
	"net/http"

	"github.com/luminous-gsm/fusion/models"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

var userModel = new(models.User)

func (u UserController) Retrieve(c *gin.Context) {
	if c.Param("id") != "" {
		c.JSON(http.StatusOK, gin.H{"message": "User founded!", "user": "Fake User"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
	c.Abort()
	return
}
