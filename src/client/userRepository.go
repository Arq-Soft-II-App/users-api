package client

import (
	"context"
	"users-api/src/models"

	"go.uber.org/zap"
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

type userRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserRepository(db *gorm.DB, logger *zap.Logger) UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	r.logger.Info("[USERS-API][Repository]: Iniciando creación de usuario en BD",
		zap.String("email", user.Email))

	if err := r.db.Create(user).Error; err != nil {
		r.logger.Error("[USERS-API][Repository]: Error al crear usuario en BD",
			zap.String("email", user.Email),
			zap.Error(err))
		return err
	}

	r.logger.Info("[USERS-API][Repository]: Usuario creado exitosamente en BD",
		zap.String("id", user.ID))
	return nil
}

func (r *userRepository) ReadAll(ctx context.Context) ([]models.User, error) {
	r.logger.Info("[USERS-API][Repository]: Iniciando búsqueda de todos los usuarios en BD")

	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		r.logger.Error("[USERS-API][Repository]: Error al obtener todos los usuarios de BD",
			zap.Error(err))
		return nil, err
	}

	r.logger.Info("[USERS-API][Repository]: Usuarios obtenidos exitosamente de BD",
		zap.Int("count", len(users)))
	return users, nil
}

func (r *userRepository) GetUsersList(ctx context.Context, ids []string) ([]models.User, error) {
	r.logger.Info("[USERS-API][Repository]: Buscando lista de usuarios por IDs en BD",
		zap.Strings("ids", ids))

	var users []models.User
	if err := r.db.Where("id IN ?", ids).Find(&users).Error; err != nil {
		r.logger.Error("[USERS-API][Repository]: Error al obtener lista de usuarios de BD",
			zap.Strings("ids", ids),
			zap.Error(err))
		return nil, err
	}

	r.logger.Info("[USERS-API][Repository]: Lista de usuarios obtenida exitosamente de BD",
		zap.Int("count", len(users)))
	return users, nil
}

func (r *userRepository) ReadByEmail(ctx context.Context, email string) (*models.User, error) {
	r.logger.Info("[USERS-API][Repository]: Buscando usuario por email en BD",
		zap.String("email", email))

	var user models.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		r.logger.Error("[USERS-API][Repository]: Error al buscar usuario por email en BD",
			zap.String("email", email),
			zap.Error(err))
		return nil, err
	}

	r.logger.Info("[USERS-API][Repository]: Usuario encontrado exitosamente en BD",
		zap.String("email", email))
	return &user, nil
}

func (r *userRepository) ReadOne(ctx context.Context, id string) (*models.User, error) {
	r.logger.Info("[USERS-API][Repository]: Buscando usuario por ID en BD",
		zap.String("id", id))

	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		r.logger.Error("[USERS-API][Repository]: Error al buscar usuario por ID en BD",
			zap.String("id", id),
			zap.Error(err))
		return nil, err
	}

	r.logger.Info("[USERS-API][Repository]: Usuario encontrado exitosamente en BD",
		zap.String("id", id))
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, id string, user *models.User) error {
	r.logger.Info("[USERS-API][Repository]: Iniciando actualización de usuario en BD",
		zap.String("id", id))

	if err := r.db.Where("id = ?", id).Updates(user).Error; err != nil {
		r.logger.Error("[USERS-API][Repository]: Error al actualizar usuario en BD",
			zap.String("id", id),
			zap.Error(err))
		return err
	}

	r.logger.Info("[USERS-API][Repository]: Usuario actualizado exitosamente en BD",
		zap.String("id", id))
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	r.logger.Info("[USERS-API][Repository]: Iniciando eliminación de usuario en BD",
		zap.String("id", id))

	if err := r.db.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		r.logger.Error("[USERS-API][Repository]: Error al eliminar usuario en BD",
			zap.String("id", id),
			zap.Error(err))
		return err
	}

	r.logger.Info("[USERS-API][Repository]: Usuario eliminado exitosamente de BD",
		zap.String("id", id))
	return nil
}
