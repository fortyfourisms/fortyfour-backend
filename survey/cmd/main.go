package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "survey/docs"
	"survey/internal/cache"
	"survey/internal/config"
	"survey/internal/handlers"
	"survey/internal/middleware"
	"survey/internal/repository"
	"survey/internal/routes"
	"survey/internal/services"
	"survey/pkg/database"

	"github.com/redis/go-redis/v9"
)

// @title						Survey API
// @version					1.0
// @description				API documentation for Survey service.
// @host						localhost:8082
// @BasePath					/
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
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

	// INIT REDIS
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// test connection (optional tapi bagus)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("Redis connection failed: " + err.Error())
	}

	// REPOSITORY
	respondenRepo := repository.NewRespondenRepository(db)
	risikoRepo := repository.NewRisikoRepository(db)

	// SERVICE
	cache := cache.NewRedisCache(rdb)
	respondenService := services.NewRespondenService(respondenRepo, services.DefaultValidator{}, cache)
	risikoService := services.NewRisikoService(risikoRepo, cache)

	// HANDLER
	respondenHandler := handlers.NewRespondenHandler(respondenService)
	risikoHandler := handlers.NewRisikoHandler(risikoService)

	// ROUTER
	mux := routes.InitRouter(
		respondenHandler,
		risikoHandler,
		middleware.Auth,
	)

	// SERVER
	addr := cfg.Port
	if !strings.Contains(addr, ":") {
		addr = ":" + addr
	}

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		println("=======================================")
		println("Survey API running on port", cfg.Port)
		println("=======================================")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// GRACEFUL SHUTDOWN
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	println("Shutting down server...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctxShutdown); err != nil {
		panic(err)
	}

	println("Server exited properly")
}
