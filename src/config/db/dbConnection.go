package db

import (
	"context" // para manejar el contexto de la conexion, incluyendo la cancelacion y los tiempos de espera
	"log"     // se usa para asegurarse de que el codigo de conexion a la base de datos se ejecute solo una vez
	"sync"
	"time"

	"users-api/src/config/envs"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var once sync.Once           // variable para asegurarse de que el codigo de conexion a la base de datos se ejecute solo una vez
var dbInstance *mongo.Client // variable para almacenar la instancia de la base de datos

// conexion con MongoDB utilizando singleton
func ConnectDB() (*mongo.Client, error) {
	var err error
	env := envs.LoadEnvs(".env")
	MONGO_URI := env.Get("MONGO_URI")

	once.Do(func() { // se asegura de que el codigo de conexion a la base de datos se ejecute solo una vez
		clientOptions := options.Client().ApplyURI(MONGO_URI)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		dbInstance, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatalf("Error al conectar con MongoDB: %v", err)
		}

		// verifica la conexion
		err = dbInstance.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("No se pudo conectar a MongoDB: %v", err)
		}

		log.Println("Conexión a MongoDB establecida")
	})

	return dbInstance, err
}
