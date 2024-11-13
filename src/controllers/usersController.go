package controllers

import (
	"net/http"
	"users-api/src/dto"
	"users-api/src/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	service services.UserService
	logger  *zap.Logger
}

func NewUserController(service services.UserService, logger *zap.Logger) *UserController {
	return &UserController{
		service: service,
		logger:  logger,
	}
}

// GetUsers maneja la solicitud GET /users/ para obtener todos los usuarios o aplicar un filtro
func (uc *UserController) GetUsers(c *gin.Context) {
	filter := make(map[string]interface{})
	if err := c.BindJSON(&filter); err != nil && err.Error() != "EOF" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar el filtro"})
		return
	}

	users, err := uc.service.GetAllUsers(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (uc *UserController) GetUsersList(c *gin.Context) {
	var requestBody struct {
		IDs []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		uc.logger.Error("Error al procesar los IDs", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Los IDs son requeridos y deben ser un array de strings"})
		return
	}

	if len(requestBody.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La lista de IDs no puede estar vacía"})
		return
	}

	users, err := uc.service.GetUsersList(c.Request.Context(), requestBody.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (uc *UserController) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	user, err := uc.service.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserByID maneja la solicitud GET /users/:id para obtener un usuario por su ID
func (uc *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := uc.service.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser maneja la solicitud POST /users/ para crear un nuevo usuario
func (uc *UserController) CreateUser(c *gin.Context) {
	var createUserDTO dto.CreateUserDTO

	if err := c.ShouldBindJSON(&createUserDTO); err != nil {
		uc.logger.Error("Error al procesar datos de usuario", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	userResponse, err := uc.service.CreateUser(c.Request.Context(), &createUserDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el usuario"})
		return
	}

	c.JSON(http.StatusCreated, userResponse)
}

// UpdateUser maneja la solicitud PUT /users/:id para actualizar un usuario existente
func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var updateUserDTO dto.UpdateUserDTO
	if err := c.ShouldBindJSON(&updateUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	userResponse, err := uc.service.UpdateUser(c.Request.Context(), id, &updateUserDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el usuario"})
		return
	}

	c.JSON(http.StatusOK, userResponse)
}

// DeleteUser maneja la solicitud DELETE /users/:id para eliminar un usuario existente
func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := uc.service.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el usuario"})
		return
	}

	c.Status(http.StatusNoContent)
}
