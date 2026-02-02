package main

import (
	"ikas/internal/config"
	"ikas/internal/handlers"
	"ikas/internal/middleware"
	"ikas/internal/repository"
	"ikas/internal/routes"
	"ikas/internal/services"
	"ikas/pkg/database"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {

	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()

	// DB
	db, err := database.NewMySQLConnection(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Layering IKAS
	ikasRepo := repository.NewIkasRepository(db)
	ikasService := services.NewIkasService(ikasRepo)
	ikasHandler := handlers.NewIkasHandler(ikasService)

	// Router
	mux := routes.InitRouter(ikasHandler)

	auth := middleware.NewAuthMiddleware(cfg.JWTSecret)
	handler := auth.Authenticate(mux)

	log.Println("IKAS service running on", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, handler))
}
