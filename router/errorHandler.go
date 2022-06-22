package router

import (
	"context"
	"net/http"
	"strings"

	"emperror.dev/errors"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
)

// RequestError is a custom error
type RequestError struct {
	err    error
	status int
	msg    string
}

func NewError(err error) *RequestError {
	return &RequestError{
		// Attach a stacktrace to the error
		err: errors.WithStackDepthIf(err, 1),
	}
}

// Aborts the request and attaches the provided error to the gin
// context so it can be reported properly.
func CaptureAndAbort(c *gin.Context, err error) {
	c.Abort()
	c.Error(errors.WithStackDepthIf(err, 1))
}

func (re *RequestError) SetMessage(m string) {
	re.msg = m
}

// Default a HTTP-500 error.
func (re *RequestError) SetStatus(s int) {
	re.status = s
}

func (re *RequestError) Cause() error {
	return re.err
}

func (re *RequestError) Error() string {
	return re.err.Error()
}

func (re *RequestError) Abort(c *gin.Context, status int) {
	reqId := c.Writer.Header().Get("X-Request-Id")

	event := log.WithField("request_id", reqId).WithField("url", c.Request.URL.String())

	if c.Writer.Status() == 200 {
		// Handle context deadlines
		if errors.Is(re.err, context.DeadlineExceeded) {
			re.SetStatus(http.StatusGatewayTimeout)
			re.SetMessage("The server could not process this request in time, please try again.")
		} else if strings.Contains(re.Cause().Error(), "context canceled") {
			re.SetStatus(http.StatusBadRequest)
			re.SetMessage("Request aborted by client.")
		}
	}

	// c.Writer.Status() will be a non-200 value. Normally marshelling issues
	if status >= 500 || c.Writer.Status() != 200 {
		event.WithField("status", status).WithField("error", re.err).Error("error while handling HTTP request")
	} else {
		event.WithField("status", status).WithField("error", re.err).Debug("error handling HTTP request (not a server error)")
	}
	if re.msg == "" {
		re.msg = "An unexpected error was encountered while processing this request"
	}

	c.AbortWithStatusJSON(status, gin.H{"error": re.msg, "request_id": reqId})
}
