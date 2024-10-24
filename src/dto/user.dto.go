package dto

import "time"

type UserDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Lastname  string    `json:"lastname"`
	Birthdate time.Time `json:"birthdate"`
	Role      string    `json:"role"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Avatar    string    `json:"avatar"`
}
