package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"users-api/src/client"
	"users-api/src/dto"
	"users-api/src/models"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService interface {
	GetAllUsers(ctx context.Context, filter map[string]interface{}) ([]dto.UserResponseDTO, error)
	GetUserByEmail(ctx context.Context, email string) (*dto.UserResponseDTO, error)
	GetUserByID(ctx context.Context, id string) (*dto.UserResponseDTO, error)
	GetUsersList(ctx context.Context, ids []string) ([]dto.UserResponseDTO, error)
	CreateUser(ctx context.Context, createUserDTO *dto.CreateUserDTO) (*dto.UserResponseDTO, error)
	UpdateUser(ctx context.Context, id string, updateUserDTO *dto.UpdateUserDTO) (*dto.UserResponseDTO, error)
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	repo        client.UserRepository
	redisClient *redis.Client
	logger      *zap.Logger
}

func NewUserService(repo client.UserRepository, redisClient *redis.Client, logger *zap.Logger) UserService {
	return &userService{
		repo:        repo,
		redisClient: redisClient,
		logger:      logger,
	}
}

func (s *userService) GetAllUsers(ctx context.Context, filter map[string]interface{}) ([]dto.UserResponseDTO, error) {
	cacheKey := "all_users"
	var userResponses []dto.UserResponseDTO

	if s.redisClient != nil {
		cachedUsers, err := s.redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			s.logger.Info("Usuarios obtenidos desde caché")
			err = json.Unmarshal([]byte(cachedUsers), &userResponses)
			if err == nil {
				return userResponses, nil
			}
		}
	}

	users, err := s.repo.ReadAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponseDTO{
			ID:        user.ID,
			Name:      user.Name,
			Lastname:  user.Lastname,
			Birthdate: user.Birthdate,
			Role:      user.Role,
			Email:     user.Email,
			Avatar:    user.Avatar,
		})
	}

	if s.redisClient != nil {
		usersJSON, _ := json.Marshal(userResponses)
		s.redisClient.Set(ctx, cacheKey, usersJSON, 5*time.Minute)
	}

	return userResponses, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponseDTO, error) {
	cacheKey := fmt.Sprintf("user_email:%s", email)

	cachedUser, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var userResponse dto.UserResponseDTO
		err = json.Unmarshal([]byte(cachedUser), &userResponse)
		if err == nil {
			return &userResponse, nil
		}
	}

	user, err := s.repo.ReadByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	userResponse := &dto.UserResponseDTO{
		ID:        user.ID,
		Name:      user.Name,
		Lastname:  user.Lastname,
		Birthdate: user.Birthdate,
		Role:      user.Role,
		Email:     user.Email,
		Avatar:    user.Avatar,
	}

	userJSON, _ := json.Marshal(userResponse)
	s.redisClient.Set(ctx, cacheKey, userJSON, 5*time.Minute)

	return userResponse, nil
}

func (s *userService) GetUsersList(ctx context.Context, ids []string) ([]dto.UserResponseDTO, error) {
	users, err := s.repo.GetUsersList(ctx, ids)
	if err != nil {
		return nil, err
	}

	var userResponses []dto.UserResponseDTO
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponseDTO{
			ID:        user.ID,
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

func (s *userService) GetUserByID(ctx context.Context, id string) (*dto.UserResponseDTO, error) {
	cacheKey := fmt.Sprintf("user_id:%s", id)

	cachedUser, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var userResponse dto.UserResponseDTO
		err = json.Unmarshal([]byte(cachedUser), &userResponse)
		if err == nil {
			return &userResponse, nil
		}
	}

	user, err := s.repo.ReadOne(ctx, id)
	if err != nil {
		return nil, err
	}

	userResponse := &dto.UserResponseDTO{
		ID:        user.ID,
		Name:      user.Name,
		Lastname:  user.Lastname,
		Birthdate: user.Birthdate,
		Role:      user.Role,
		Email:     user.Email,
		Avatar:    user.Avatar,
	}

	userJSON, _ := json.Marshal(userResponse)
	s.redisClient.Set(ctx, cacheKey, userJSON, 5*time.Minute)

	return userResponse, nil
}

func (s *userService) CreateUser(ctx context.Context, createUserDTO *dto.CreateUserDTO) (*dto.UserResponseDTO, error) {
	hashedPassword, err := createUserDTO.ValidateAndHash()
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:        uuid.New().String(),
		Name:      createUserDTO.Name,
		Lastname:  createUserDTO.Lastname,
		Birthdate: createUserDTO.Birthdate,
		Role:      createUserDTO.Role,
		Email:     createUserDTO.Email,
		Password:  hashedPassword,
		Avatar:    createUserDTO.Avatar,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	userResponse := &dto.UserResponseDTO{
		ID:        user.ID,
		Name:      user.Name,
		Lastname:  user.Lastname,
		Birthdate: user.Birthdate,
		Role:      user.Role,
		Email:     user.Email,
		Avatar:    user.Avatar,
	}

	if s.redisClient != nil {
		s.redisClient.Del(ctx, "all_users")
		s.redisClient.Del(ctx, fmt.Sprintf("user_email:%s", user.Email))
		s.redisClient.Del(ctx, fmt.Sprintf("user_id:%s", user.ID))
	}

	return userResponse, nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, updateUserDTO *dto.UpdateUserDTO) (*dto.UserResponseDTO, error) {
	user, err := s.repo.ReadOne(ctx, id)
	if err != nil {
		return nil, err
	}

	if updateUserDTO.Name != nil {
		user.Name = *updateUserDTO.Name
	}
	if updateUserDTO.Lastname != nil {
		user.Lastname = *updateUserDTO.Lastname
	}
	if updateUserDTO.Birthdate != nil {
		user.Birthdate = *updateUserDTO.Birthdate
	}
	if updateUserDTO.Role != nil {
		user.Role = *updateUserDTO.Role
	}
	if updateUserDTO.Email != nil {
		user.Email = *updateUserDTO.Email
	}
	if updateUserDTO.Password != nil {
		user.Password = *updateUserDTO.Password
	}
	if updateUserDTO.Avatar != nil {
		user.Avatar = *updateUserDTO.Avatar
	}

	if err := s.repo.Update(ctx, id, user); err != nil {
		return nil, err
	}

	userResponse := &dto.UserResponseDTO{
		ID:        user.ID,
		Name:      user.Name,
		Lastname:  user.Lastname,
		Birthdate: user.Birthdate,
		Role:      user.Role,
		Email:     user.Email,
		Avatar:    user.Avatar,
	}

	if s.redisClient != nil {
		s.redisClient.Del(ctx, "all_users")
		s.redisClient.Del(ctx, fmt.Sprintf("user_email:%s", user.Email))
		s.redisClient.Del(ctx, fmt.Sprintf("user_id:%s", user.ID))
	}

	return userResponse, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	user, err := s.repo.ReadOne(ctx, id)
	if err != nil {
		return err
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	if s.redisClient != nil {
		s.redisClient.Del(ctx, "all_users")
		s.redisClient.Del(ctx, fmt.Sprintf("user_email:%s", user.Email))
		s.redisClient.Del(ctx, fmt.Sprintf("user_id:%s", user.ID))
	}

	return nil
}
