package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIKeyAuthMiddleware(expectedAPIKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("Authorization")

		if apiKey != expectedAPIKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
			c.Abort()
			return
		}

		c.Next()
	}
}
