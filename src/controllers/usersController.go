package controllers

import (
	"log"
	"net/http"
	"users-api/src/dto"
	"users-api/src/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
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
		log.Fatal(err)
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
