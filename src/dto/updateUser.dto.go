package dto

import (
	"time"
	"users-api/src/utils"
)

type UpdateUserDTO struct {
	Name      *string    `json:"name,omitempty"`
	Lastname  *string    `json:"lastname,omitempty"`
	Birthdate *time.Time `json:"birthdate,omitempty"`
	Role      *string    `json:"role,omitempty"`
	Email     *string    `json:"email,omitempty"`
	Password  *string    `json:"password,omitempty"`
	Avatar    *string    `json:"avatar,omitempty"`
}

// ValidateAndHashPassword hashea la contraseña si es proporcionada
func (dto *UpdateUserDTO) ValidateAndHashPassword() (string, error) {
	if dto.Password != nil && *dto.Password != "" {
		hashedPassword, err := utils.HashPassword(*dto.Password)
		if err != nil {
			return "", err
		}
		return hashedPassword, nil
	}
	return "", nil
}
