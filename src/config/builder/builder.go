package builder

import (
	"context"
	"log"
	"users-api/src/client"
	"users-api/src/config/db"
	"users-api/src/controllers"
	"users-api/src/router"
	"users-api/src/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppBuilder struct {
	db             *mongo.Client
	userRepo       client.UserRepository
	userService    services.UserService
	userController *controllers.UserController
	router         *gin.Engine
}

func NewAppBuilder() *AppBuilder {
	return &AppBuilder{}
}

func (b *AppBuilder) BuildDBConnection() *AppBuilder {
	var err error
	b.db, err = db.ConnectDB()
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	return b
}

func (b *AppBuilder) DisconnectDB() {
	if b.db != nil {
		if err := b.db.Disconnect(context.Background()); err != nil {
			log.Fatalf("Error al desconectar de la base de datos: %v", err)
		}
	}
}

func (b *AppBuilder) BuildUserRepo() *AppBuilder {
	b.userRepo = client.NewMongoUserRepository(b.db)
	return b
}

func (b *AppBuilder) BuildUserService() *AppBuilder {
	b.userService = services.NewUserService(b.userRepo)
	return b
}

func (b *AppBuilder) BuildUserController() *AppBuilder {
	b.userController = controllers.NewUserController(b.userService)
	return b
}

func (b *AppBuilder) BuildRouter() *AppBuilder {
	b.router = gin.Default()
	router.SetupRoutes(b.router, b.userController)
	return b
}

func (b *AppBuilder) GetRouter() *gin.Engine {
	return b.router
}
