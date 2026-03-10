package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"survey/internal/config"
	"survey/internal/handlers"
	"survey/internal/repository"
	"survey/internal/routes"
	"survey/internal/services"
	"survey/pkg/database"
)

func main() {

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewMySQLConnection(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	})
	if err != nil {
		panic("Failed to connect database: " + err.Error())
	}
	defer db.Close()

	// Repository
	respondenRepo := repository.NewRespondenRepository(db)

	// Services
	respondenService := services.NewRespondenService(respondenRepo)

	// Handlers
	respondenHandler := handlers.NewRespondenHandler(respondenService)

	// Router
	mux := routes.InitRouter(
		respondenHandler,
	)

	// HTTP Server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	// Run server
	go func() {
		println("=======================================")
		println("Survey API running on port :", cfg.Port)
		println("Endpoint:")
		println("GET  /api/health")
		println("GET  /api/responden")
		println("POST /api/responden")
		println("=======================================")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("Server error: " + err.Error())
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic("Server forced shutdown: " + err.Error())
	}

	println("Server exited properly")
}