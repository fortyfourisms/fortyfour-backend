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
	// Load config
	cfg := config.Load()

	// Init DB
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

	// REPOSITORY
	respondenRepo := repository.NewRespondenRepository(db)
	risikoRepo := repository.NewRisikoRepository(db)

	// SERVICE
	respondenService := services.NewRespondenService(respondenRepo)
	risikoService := services.NewRisikoService(risikoRepo)

	// HANDLER
	respondenHandler := handlers.NewRespondenHandler(respondenService)
	risikoHandler := handlers.NewRisikoHandler(risikoService)

	// ROUTER
	mux := routes.InitRouter(
		respondenHandler,
		risikoHandler,
	)

	// SERVER
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	// Run server
	go func() {
		println("=======================================")
		println("Survey API running on port :", cfg.Port)
		println("=======================================")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}

	println("Server exited properly")
} 