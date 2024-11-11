package client

import (
	"context"
	"users-api/src/errors"
	"users-api/src/models"

	"net/http"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	ReadAll(ctx context.Context) ([]models.User, error)
	GetUsersList(ctx context.Context, ids []string) ([]models.User, error)
	ReadByEmail(ctx context.Context, email string) (*models.User, error)
	ReadOne(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, id string, user *models.User) error
	Delete(ctx context.Context, id string) error
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{
		db: db,
	}
}

func (r *gormUserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *gormUserRepository) ReadAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}

func (r *gormUserRepository) GetUsersList(ctx context.Context, ids []string) ([]models.User, error) {
	var users []models.User
	result := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&users)

	if result.Error != nil {
		return nil, errors.NewError("DB_ERROR", "Error al recuperar los usuarios de la base de datos", http.StatusInternalServerError)
	}

	if len(users) != len(ids) {
		return nil, errors.NewError("NOT_FOUND", "No se encontraron todos los usuarios solicitados", http.StatusNotFound)
	}

	return users, nil
}

func (r *gormUserRepository) ReadByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewError("NOT_FOUND", "Usuario no encontrado", http.StatusNotFound)
		}
		return nil, errors.NewError("DB_ERROR", "Error al recuperar el usuario de la base de datos", http.StatusInternalServerError)
	}
	return &user, nil
}

func (r *gormUserRepository) ReadOne(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewError("NOT_FOUND", "Usuario no encontrado", http.StatusNotFound)
		}
		return nil, errors.NewError("DB_ERROR", "Error al recuperar el usuario de la base de datos", http.StatusInternalServerError)
	}
	return &user, nil
}

func (r *gormUserRepository) Update(ctx context.Context, id string, user *models.User) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(user).Error
}

func (r *gormUserRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}
