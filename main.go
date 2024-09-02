package main

import (
	"log"
	"users-api/src/config/builder"
	"users-api/src/config/envs"
)

func main() {
	env := envs.LoadEnvs(".env")

	app := builder.NewAppBuilder().
		BuildDBConnection().
		BuildUserRepo().
		BuildUserService().
		BuildUserController().
		BuildRouter()

	defer app.DisconnectDB()

	port := env.Get("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Iniciando servidor en el puerto %s...", port)
	if err := app.GetRouter().Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
