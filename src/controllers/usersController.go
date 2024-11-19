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
	uc.logger.Info("[USERS-API]: Iniciando obtención de usuarios")

	filter := make(map[string]interface{})
	if err := c.BindJSON(&filter); err != nil && err.Error() != "EOF" {
		uc.logger.Error("[USERS-API]: Error al procesar filtro", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar el filtro"})
		return
	}

	users, err := uc.service.GetAllUsers(c.Request.Context(), filter)
	if err != nil {
		uc.logger.Error("[USERS-API]: Error al obtener usuarios", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios"})
		return
	}

	uc.logger.Info("[USERS-API]: Usuarios obtenidos exitosamente", zap.Int("count", len(users)))
	c.JSON(http.StatusOK, users)
}

func (uc *UserController) GetUsersList(c *gin.Context) {
	uc.logger.Info("[USERS-API]: Iniciando obtención de lista de usuarios por IDs")

	var requestBody struct {
		IDs []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		uc.logger.Error("[USERS-API]: Error al procesar los IDs", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Los IDs son requeridos y deben ser un array de strings"})
		return
	}

	if len(requestBody.IDs) == 0 {
		uc.logger.Warn("[USERS-API]: Se recibió una lista vacía de IDs")
		c.JSON(http.StatusBadRequest, gin.H{"error": "La lista de IDs no puede estar vacía"})
		return
	}

	users, err := uc.service.GetUsersList(c.Request.Context(), requestBody.IDs)
	if err != nil {
		uc.logger.Error("[USERS-API]: Error al obtener lista de usuarios", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios"})
		return
	}

	uc.logger.Info("[USERS-API]: Lista de usuarios obtenida exitosamente", zap.Int("count", len(users)))
	c.JSON(http.StatusOK, users)
}

func (uc *UserController) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	uc.logger.Info("[USERS-API]: Buscando usuario por email", zap.String("email", email))

	user, err := uc.service.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		uc.logger.Error("[USERS-API]: Usuario no encontrado por email", zap.String("email", email), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	uc.logger.Info("[USERS-API]: Usuario encontrado exitosamente por email", zap.String("email", email))
	c.JSON(http.StatusOK, user)
}

// GetUserByID maneja la solicitud GET /users/:id para obtener un usuario por su ID
func (uc *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	uc.logger.Info("[USERS-API]: Buscando usuario por ID", zap.String("id", id))

	user, err := uc.service.GetUserByID(c.Request.Context(), id)
	if err != nil {
		uc.logger.Error("[USERS-API]: Usuario no encontrado por ID", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	uc.logger.Info("[USERS-API]: Usuario encontrado exitosamente por ID", zap.String("id", id))
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
	uc.logger.Info("[USERS-API]: Iniciando actualización de usuario", zap.String("id", id))

	var updateUserDTO dto.UpdateUserDTO
	if err := c.ShouldBindJSON(&updateUserDTO); err != nil {
		uc.logger.Error("[USERS-API]: Error al procesar datos de actualización", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	userResponse, err := uc.service.UpdateUser(c.Request.Context(), id, &updateUserDTO)
	if err != nil {
		uc.logger.Error("[USERS-API]: Error al actualizar usuario", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el usuario"})
		return
	}

	uc.logger.Info("[USERS-API]: Usuario actualizado exitosamente", zap.String("id", id))
	c.JSON(http.StatusOK, userResponse)
}

// DeleteUser maneja la solicitud DELETE /users/:id para eliminar un usuario existente
func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	uc.logger.Info("[USERS-API]: Iniciando eliminación de usuario", zap.String("id", id))

	if err := uc.service.DeleteUser(c.Request.Context(), id); err != nil {
		uc.logger.Error("[USERS-API]: Error al eliminar usuario", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el usuario"})
		return
	}

	uc.logger.Info("[USERS-API]: Usuario eliminado exitosamente", zap.String("id", id))
	c.Status(http.StatusNoContent)
}
