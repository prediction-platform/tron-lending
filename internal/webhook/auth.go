package webhook

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var authToken = "mysecrettoken"

func init() {
	if v := os.Getenv("WEBHOOK_AUTH_TOKEN"); v != "" {
		authToken = v
	}
}

// AuthMiddleware gin 鉴权中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Auth-Token")
		if token != authToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
