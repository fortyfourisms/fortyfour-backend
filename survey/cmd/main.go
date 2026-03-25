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

	cfg := config.Load()

	db, err := database.NewMySQLConnection(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Initialize Repositories
	respondenRepo := repository.NewRespondenRepository(db)
	risikoRepo := repository.NewRisikoRepository(db)

	// Initialize Services
	respondenService := services.NewRespondenService(respondenRepo)
	risikoService := services.NewRisikoService(risikoRepo)

	// Initialize Handler
	respondenHandler := handlers.NewRespondenHandler(respondenService)
	risikoHandler := handlers.NewRisikoHandler(risikoService)

	// ROUTER
	mux := routes.InitRouter(
		respondenHandler,
		risikoHandler,
	)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	go func() {

		println("Survey API running on", cfg.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)
}