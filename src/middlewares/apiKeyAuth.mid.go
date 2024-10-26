package middlewares

import (
	"fmt"
	"net/http"
	"users-api/src/config/envs"

	"github.com/gin-gonic/gin"
)

func APIKeyAuthMiddleware() gin.HandlerFunc {
	KEY := envs.LoadEnvs(".env").Get("USERS_API_KEY")
	return func(c *gin.Context) {
		apiKey := c.GetHeader("Authorization")
		fmt.Println(apiKey)

		if apiKey != KEY {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
			c.Abort()
			return
		}

		c.Next()
	}
}
