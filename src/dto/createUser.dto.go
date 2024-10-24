package dto

import (
	"errors"
	"strings"
	"time"
	"users-api/src/utils"

	"github.com/google/uuid"
)

type CreateUserDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Lastname  string    `json:"lastname" binding:"required"`
	Birthdate time.Time `json:"birthdate" binding:"required"`
	Role      string    `json:"role"`
	Email     string    `json:"email" binding:"required,email"`
	Password  string    `json:"password" binding:"required,min=6"`
	Avatar    string    `json:"avatar"`
}

func (dto *CreateUserDTO) ValidateAndHash() (string, error) {
	if strings.TrimSpace(dto.Name) == "" {
		return "", errors.New("el nombre es obligatorio")
	}

	if strings.TrimSpace(dto.Lastname) == "" {
		return "", errors.New("el apellido es obligatorio")
	}

	if strings.TrimSpace(dto.Email) == "" {
		return "", errors.New("el email es obligatorio")
	}

	if strings.TrimSpace(dto.Password) == "" {
		return "", errors.New("la contraseña es obligatoria y debe tener al menos 6 caracteres")
	}

	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(dto.Role) == "" {
		dto.Role = "user"
	}

	if strings.TrimSpace(dto.Avatar) == "" {
		dto.Avatar = "https://i.postimg.cc/wTgNFWhR/profile.png"
	}

	return hashedPassword, nil
}
