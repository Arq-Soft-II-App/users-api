package controllers

import (
	"net/http"
	"users-api/src/dto"
	"users-api/src/errors"
	"users-api/src/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthController struct {
	service services.AuthService
	logger  *zap.Logger
}

func NewAuthController(service services.AuthService, logger *zap.Logger) *AuthController {
	return &AuthController{
		service: service,
		logger:  logger,
	}
}

func (ac *AuthController) Login(c *gin.Context) {
	loginDTO := &dto.LoginDTO{}
	if err := c.BindJSON(loginDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error al procesar la solicitud",
			"code":  "INVALID_REQUEST",
		})
		return
	}

	user, err := ac.service.Login(c.Request.Context(), loginDTO)
	if err != nil {
		if customErr, ok := err.(*errors.Error); ok {
			c.JSON(customErr.HTTPStatusCode, gin.H{
				"error": customErr.Message,
				"code":  customErr.Code,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
				"code":  "INTERNAL_SERVER_ERROR",
			})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}
