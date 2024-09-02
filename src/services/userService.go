package services

import (
	"context"
	"users-api/src/client"
	"users-api/src/dto"
	"users-api/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserService define la interfaz que implementará el servicio de usuarios
type UserService interface {
	GetAllUsers(ctx context.Context, filter map[string]interface{}) ([]dto.UserResponseDTO, error)
	GetUserByID(ctx context.Context, id string) (*dto.UserResponseDTO, error)
	CreateUser(ctx context.Context, createUserDTO *dto.CreateUserDTO) (*dto.UserResponseDTO, error)
	UpdateUser(ctx context.Context, id string, updateUserDTO *dto.UpdateUserDTO) (*dto.UserResponseDTO, error)
	DeleteUser(ctx context.Context, id string) error
}

// userService es la implementación del UserService
type userService struct {
	repo client.UserRepository
}

// NewUserService crea una nueva instancia de UserService
func NewUserService(repo client.UserRepository) UserService {
	return &userService{repo: repo}
}

// GetAllUsers devuelve una lista de usuarios según un filtro
func (s *userService) GetAllUsers(ctx context.Context, filter map[string]interface{}) ([]dto.UserResponseDTO, error) {
	// Convierte el filtro map[string]interface{} a bson.M
	bsonFilter := make(map[string]interface{})
	for k, v := range filter {
		bsonFilter[k] = v
	}

	users, err := s.repo.ReadAll(ctx, bsonFilter)
	if err != nil {
		return nil, err
	}

	// Convertir la lista de usuarios en una lista de DTOs de respuesta
	var userResponses []dto.UserResponseDTO
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponseDTO{
			ID:        user.ID.Hex(),
			Name:      user.Name,
			Lastname:  user.Lastname,
			Birthdate: user.Birthdate,
			Role:      user.Role,
			Email:     user.Email,
			Avatar:    user.Avatar,
		})
	}

	return userResponses, nil
}

// GetUserByID devuelve un usuario por su ID
func (s *userService) GetUserByID(ctx context.Context, id string) (*dto.UserResponseDTO, error) {
	user, err := s.repo.ReadOne(ctx, id)
	if err != nil {
		return nil, err
	}

	userResponse := &dto.UserResponseDTO{
		ID:        user.ID.Hex(),
		Name:      user.Name,
		Lastname:  user.Lastname,
		Birthdate: user.Birthdate,
		Role:      user.Role,
		Email:     user.Email,
		Avatar:    user.Avatar,
	}

	return userResponse, nil
}

// CreateUser crea un nuevo usuario
func (s *userService) CreateUser(ctx context.Context, createUserDTO *dto.CreateUserDTO) (*dto.UserResponseDTO, error) {
	hashedPassword, err := createUserDTO.ValidateAndHash()
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:        primitive.NewObjectID(),
		Name:      createUserDTO.Name,
		Lastname:  createUserDTO.Lastname,
		Birthdate: createUserDTO.Birthdate, // Convierte a time.Time si es necesario
		Role:      createUserDTO.Role,
		Email:     createUserDTO.Email,
		Password:  hashedPassword,
		Avatar:    createUserDTO.Avatar,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	userResponse := &dto.UserResponseDTO{
		ID:        user.ID.Hex(),
		Name:      user.Name,
		Lastname:  user.Lastname,
		Birthdate: user.Birthdate, // Convierte a time.Time si es necesario
		Role:      user.Role,
		Email:     user.Email,
		Avatar:    user.Avatar,
	}

	return userResponse, nil
}

// UpdateUser actualiza un usuario existente
func (s *userService) UpdateUser(ctx context.Context, id string, updateUserDTO *dto.UpdateUserDTO) (*dto.UserResponseDTO, error) {
	user, err := s.repo.ReadOne(ctx, id)
	if err != nil {
		return nil, err
	}

	// Actualizar solo los campos que están presentes en el DTO
	if updateUserDTO.Name != nil {
		user.Name = *updateUserDTO.Name
	}
	if updateUserDTO.Lastname != nil {
		user.Lastname = *updateUserDTO.Lastname
	}
	if updateUserDTO.Birthdate != nil {
		user.Birthdate = *updateUserDTO.Birthdate // Convierte a time.Time si es necesario
	}
	if updateUserDTO.Role != nil {
		user.Role = *updateUserDTO.Role
	}
	if updateUserDTO.Email != nil {
		user.Email = *updateUserDTO.Email
	}
	if updateUserDTO.Password != nil {
		hashedPassword, err := updateUserDTO.ValidateAndHashPassword()
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}
	if updateUserDTO.Avatar != nil {
		user.Avatar = *updateUserDTO.Avatar
	}

	if err := s.repo.Update(ctx, id, user); err != nil {
		return nil, err
	}

	userResponse := &dto.UserResponseDTO{
		ID:        user.ID.Hex(),
		Name:      user.Name,
		Lastname:  user.Lastname,
		Birthdate: user.Birthdate, // Convierte a time.Time si es necesario
		Role:      user.Role,
		Email:     user.Email,
		Avatar:    user.Avatar,
	}

	return userResponse, nil
}

// DeleteUser elimina un usuario existente por su ID
func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
