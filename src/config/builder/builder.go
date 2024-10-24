package builder

import (
	"log"
	"users-api/src/client"
	"users-api/src/config/db"
	"users-api/src/config/redis"
	"users-api/src/controllers"
	"users-api/src/router"
	"users-api/src/services"

	"github.com/gin-gonic/gin"
	redisClient "github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type AppBuilder struct {
	db             *gorm.DB
	redisClient    *redisClient.Client
	userRepo       client.UserRepository
	userService    services.UserService
	authService    services.AuthService
	userController *controllers.UserController
	authController *controllers.AuthController
	router         *gin.Engine
}

func NewAppBuilder() *AppBuilder {
	return &AppBuilder{}
}

func BuildApp() *AppBuilder {
	return NewAppBuilder().
		BuildDBConnection().
		BuildUserRepo().
		BuildUserService().
		BuildUserController().
		BuildRouter()
}

func (b *AppBuilder) BuildDBConnection() *AppBuilder {
	var err error
	b.db, err = db.ConnectDB()
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	b.redisClient = redis.ConnectRedis()
	if b.redisClient == nil {
		log.Println("Redis no está disponible. La aplicación funcionará sin caché.")
	}
	return b
}

func (b *AppBuilder) DisconnectDB() {
	if b.db != nil {
		sqlDB, err := b.db.DB()
		if err != nil {
			log.Printf("Error al obtener la conexión SQL: %v", err)
		} else {
			if err := sqlDB.Close(); err != nil {
				log.Printf("Error al desconectar de la base de datos: %v", err)
			}
		}
	}
	if b.redisClient != nil {
		if err := b.redisClient.Close(); err != nil {
			log.Printf("Error al cerrar la conexión de Redis: %v", err)
		}
	}
}

func (b *AppBuilder) BuildUserRepo() *AppBuilder {
	b.userRepo = client.NewGormUserRepository(b.db)
	return b
}

func (b *AppBuilder) BuildUserService() *AppBuilder {
	b.userService = services.NewUserService(b.userRepo, b.redisClient)
	b.authService = services.NewAuthService(b.userRepo, b.redisClient)
	return b
}

func (b *AppBuilder) BuildUserController() *AppBuilder {
	b.userController = controllers.NewUserController(b.userService)
	b.authController = controllers.NewAuthController(b.authService)
	return b
}

func (b *AppBuilder) BuildRouter() *AppBuilder {
	b.router = gin.Default()
	router.SetupRoutes(b.router, b.userController, b.authController)
	return b
}

func (b *AppBuilder) GetRouter() *gin.Engine {
	return b.router
}
