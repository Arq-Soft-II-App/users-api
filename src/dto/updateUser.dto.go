package dto

import (
	"time"
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
