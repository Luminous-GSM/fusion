package router

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"emperror.dev/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RequestError struct {
	err     error
	uuid    string
	message string
}

func WithError(c *gin.Context, err error) error {
	return c.Error(errors.WithStackDepthIf(err, 1))
}

func NewError(err error) *RequestError {
	return &RequestError{
		err:  err,
		uuid: uuid.Must(uuid.NewRandom()).String(),
	}
}

func (e *RequestError) SetMessage(msg string) *RequestError {
	e.message = msg
	return e
}

func (e *RequestError) AbortWithStatus(c *gin.Context, status int) {
	if c.Writer.Status() != 200 {
		status = c.Writer.Status()
	}

	if errors.Is(e.err, os.ErrNotExist) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "The requested resource was not found on the system.",
		})
		return
	}

	if strings.HasPrefix(e.err.Error(), "invalid URL escape") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Some of the data provided in the request appears to be escaped improperly.",
		})
		return
	}

	if e.message == "" {
		e.message = "An unexpected error was encountered while processing this request."
	}

	c.AbortWithStatusJSON(status, gin.H{"error": e.message, "error_id": e.uuid})
}

func (e *RequestError) Abort(c *gin.Context) {
	e.AbortWithStatus(c, http.StatusInternalServerError)
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("%v (uuid: %s)", e.err, e.uuid)
}
