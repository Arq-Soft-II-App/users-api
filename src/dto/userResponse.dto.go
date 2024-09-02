package dto

import "time"

// UserResponseDTO es el DTO para enviar la respuesta al cliente
type UserResponseDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Lastname  string    `json:"lastname"`
	Birthdate time.Time `json:"birthdate"`
	Role      string    `json:"role"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
}

type UsersResponseDto []UserResponseDTO
