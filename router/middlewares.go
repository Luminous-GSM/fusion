package router

import (
	"crypto/subtle"
	"io"
	"net/http"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/luminous-gsm/fusion/config"
)

// Attaches a unique ID to the incoming HTTP request.
// The request id is also added to the Gin Context for future reference of the request id.
func AttachRequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := uuid.New().String()
		ctx.Set("request_id", id)
		ctx.Set("logger", log.WithField("request_id", id))
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
			log.WithFields(log.Fields{
				"client_ip":  params.ClientIP,
				"status":     params.StatusCode,
				"latency":    params.Latency,
				"request_id": params.Keys["request_id"],
			}).Debugf("%s %s", params.MethodColor()+params.Method+params.ResetColor(), params.Path)

			return ""
		},
	)
}
