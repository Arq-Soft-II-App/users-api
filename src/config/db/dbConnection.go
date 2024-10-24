package db

import (
	"log"
	"sync"
	"users-api/src/config/envs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var once sync.Once
var dbInstance *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	var err error
	POSTGRES_URI := envs.LoadEnvs(".env").Get("POSTGRES_URI")

	once.Do(func() {
		dbInstance, err = gorm.Open(postgres.Open(POSTGRES_URI), &gorm.Config{})
		if err != nil {
			log.Fatalf("Error al conectar con PostgreSQL: %v", err)
		}

		/* 		err = dbInstance.AutoMigrate(&models.User{})
		   		if err != nil {
		   			log.Fatalf("Error al realizar la migración automática: %v", err)
		   		} */

		log.Println("Conexión a PostgreSQL establecida")
	})

	return dbInstance, err
}
