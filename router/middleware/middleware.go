package middleware

import (
	"context"
	"crypto/subtle"
	"io"
	"net/http"
	"strings"

	"emperror.dev/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/luminous-gsm/fusion/config"
	"github.com/luminous-gsm/fusion/server"
	"go.uber.org/zap"
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
		zap.S().Errorw("error while handling HTTP request",
			"requestId", reqId,
			"url", c.Request.URL,
			"status", status,
			"Error", re.err,
		)
	} else {
		zap.S().Debugw("error handling HTTP request (not a server error)",
			"requestId", reqId,
			"url", c.Request.URL,
			"status", status,
			"Error", re.err,
		)
	}
	if re.msg == "" {
		re.msg = "An unexpected error was encountered while processing this request"
	}

	c.AbortWithStatusJSON(status, gin.H{"error": re.msg, "request_id": reqId})
}

// Attaches a unique ID to the incoming HTTP request.
// The request id is also added to the Gin Context for future reference of the request id.
func AttachRequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := uuid.New().String()
		ctx.Set("request_id", id)
		ctx.Header("X-Request-Id", id)
		ctx.Next()
	}
}

// Allows for better error handling and provides a great way for log searching.
func CaptureErrors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		err := ctx.Errors.Last()
		if err == nil || err.Err == nil {
			return
		}

		status := http.StatusInternalServerError
		if ctx.Writer.Status() != 200 {
			status = ctx.Writer.Status()
		}
		if err.Error() == io.EOF.Error() {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "The data passed in the request was not in a parsable format. Please try again."})
			return
		}
		captured := NewError(err.Err)
		captured.Abort(ctx, status)
	}
}

// Sets the access request control headers on all of
// the requests.
func SetAccessControlHeaders() gin.HandlerFunc {
	cfg := config.Get()
	location := cfg.ConsoleLocation
	allowPrivateNetwork := cfg.AllowPrivateNetwork

	return func(ctx *gin.Context) {
		ctx.Header("Actxcess-Control-Allow-Origin", location)
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Accept, Accept-Encoding, Authorization, Cache-Control, Content-Type, Content-Length, Origin, X-Real-IP, X-CSRF-Token, X-Api-Key")

		// @see https://developer.chrome.com/blog/private-network-access-update/?utm_source=devtools
		if allowPrivateNetwork {
			ctx.Header("Access-Control-Request-Private-Network", "true")
		}

		// @see https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age#Directives
		ctx.Header("Access-Control-Max-Age", "7200")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}
		ctx.Next()
	}
}

// Authenticates the request token against the given
// permission string. A unique node token is required to operate this node.
func RequireAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := config.Get().Api.Security.Token
		auth := ctx.Request.Header.Get("X-Auth-Key")
		if auth == "" {
			ctx.Header("WWW-Authenticate", "X Api Key")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "The required authorization heads were not present in the request."})
			return
		}

		if subtle.ConstantTimeCompare([]byte(auth), []byte(token)) != 1 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this endpoint."})
			return
		}
		ctx.Next()
	}
}

// Adds advanced logging to Hin for per request logging
func AdvancedLogging() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			zap.S().Debugf("%s %s", params.MethodColor()+params.Method+params.ResetColor(), params.Path,
				"client_ip", params.ClientIP,
				"status", params.StatusCode,
				"latency", params.Latency,
				"request_id", params.Keys["request_id"],
			)
			return ""
		},
	)
}

// Attach the ServerManager to the gin context for easy access to the manager down the line.
func AttachServerManager(m *server.ServerManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("manager", m)
	}
}

// Will return the server from the gin.Context or panic if it is
// not present.
func GetServerManager(c *gin.Context) *server.ServerManager {
	v, ok := c.Get("manager")
	if !ok {
		panic("Cannot extract server manager")
	}
	return v.(*server.ServerManager)
}

func RequestDataValidator() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
