package middlewares

import (
	"net/http"
	"users-api/src/config/envs"

	"github.com/gin-gonic/gin"
)

func APIKeyAuthMiddleware() gin.HandlerFunc {
	KEY := envs.LoadEnvs(".env").Get("USERS_API_KEY")
	return func(c *gin.Context) {
		apiKey := c.GetHeader("Authorization")

		if apiKey != KEY {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
			c.Abort()
			return
		}

		c.Next()
	}
}
