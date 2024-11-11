package router

import (
	"net/http"
	"users-api/src/controllers"
	"users-api/src/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userController *controllers.UserController, authController *controllers.AuthController) {
	// Middleware para autenticación con API Key
	router.Use(middlewares.APIKeyAuthMiddleware())

	// Configurar rutas para el servicio de usuarios
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", userController.GetUsers)
		userRoutes.GET("/email/:email", userController.GetUserByEmail)
		userRoutes.GET("/list", userController.GetUsersList)
		userRoutes.GET("/:id", userController.GetUserByID)
		userRoutes.POST("/", userController.CreateUser)
		userRoutes.POST("/login", authController.Login)
		userRoutes.PUT("/:id", userController.UpdateUser)
		userRoutes.DELETE("/:id", userController.DeleteUser)
	}

	// Handler para rutas no encontradas
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ruta no encontrada"})

	})
}
