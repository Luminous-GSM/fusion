package middlewares

import (
	"strings"

	"github.com/luminous-gsm/fusion/config"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestToken := c.Request.Header.Get("X-Auth-Key")

		var token string
		if token = config.Get().Api.Security.Token; len(strings.TrimSpace(token)) == 0 {
			c.AbortWithStatus(500)
		}
		if token != requestToken {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}
