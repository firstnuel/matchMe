package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping returns a middleware that responds with "pong" to GET requests on "/ping".
func Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet && c.Request.URL.Path == "/ping" {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
			c.Abort()
			return
		}
		c.Next()
	}
}
