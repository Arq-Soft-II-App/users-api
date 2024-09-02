package router

import (
	"net/http"
	"users-api/src/config/envs"
	"users-api/src/controllers"
	"users-api/src/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userController *controllers.UserController) {
	// Aplicar middleware de verificación de API key a todas las rutas
	env := envs.LoadEnvs(".env")
	API_KEY := env.Get("USERS_API_KEY")
	router.Use(middlewares.APIKeyAuthMiddleware(API_KEY))

	// Configurar rutas para el servicio de usuarios
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", userController.GetUsers)
		userRoutes.GET("/:id", userController.GetUserByID)
		userRoutes.POST("/", userController.CreateUser)
		userRoutes.PUT("/:id", userController.UpdateUser)
		userRoutes.DELETE("/:id", userController.DeleteUser)
	}
	// Handler para rutas no encontradas
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ruta no encontrada"})
	})
}
