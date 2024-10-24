package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"users-api/src/client"
	"users-api/src/dto"
	"users-api/src/errors"
	"users-api/src/utils"

	"github.com/go-redis/redis/v8"
)

type AuthService interface {
	Login(ctx context.Context, loginDTO *dto.LoginDTO) (*dto.UserResponseDTO, error)
}

type authService struct {
	repo        client.UserRepository
	redisClient *redis.Client
}

func NewAuthService(repo client.UserRepository, redisClient *redis.Client) AuthService {
	return &authService{repo: repo, redisClient: redisClient}
}

func (s *authService) Login(ctx context.Context, loginDTO *dto.LoginDTO) (*dto.UserResponseDTO, error) {
	cacheKey := fmt.Sprintf("user_email:%s", loginDTO.Email)
	var user *dto.UserDTO

	if s.redisClient != nil {
		// Intentar obtener de caché
		cachedUser, err := s.redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			err = json.Unmarshal([]byte(cachedUser), &user)
			if err == nil {
				// Verificar la contraseña
				if !utils.CheckPasswordHash(loginDTO.Password, user.Password) {
					return nil, errors.NewError("INVALID CREDENTIALS", "Invalid credentials", 401)
				}
				return &dto.UserResponseDTO{
					ID:        user.ID,
					Name:      user.Name,
					Lastname:  user.Lastname,
					Birthdate: user.Birthdate,
					Role:      user.Role,
					Email:     user.Email,
					Avatar:    user.Avatar,
				}, nil
			}
		}
	}

	// Si no está en caché o Redis no está disponible, obtener de la base de datos
	dbUser, err := s.repo.ReadByEmail(ctx, loginDTO.Email)
	if err != nil {
		return nil, err
	}

	// Verificar la contraseña
	if !utils.CheckPasswordHash(loginDTO.Password, dbUser.Password) {
		return nil, errors.NewError("INVALID CREDENTIALS", "Invalid credentials", 401)
	}

	userResponse := &dto.UserResponseDTO{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Lastname:  dbUser.Lastname,
		Birthdate: dbUser.Birthdate,
		Role:      dbUser.Role,
		Email:     dbUser.Email,
		Avatar:    dbUser.Avatar,
	}

	if s.redisClient != nil {
		// Guardar en caché
		userJSON, _ := json.Marshal(dbUser)
		s.redisClient.Set(ctx, cacheKey, userJSON, 5*time.Minute)
	}

	return userResponse, nil
}
