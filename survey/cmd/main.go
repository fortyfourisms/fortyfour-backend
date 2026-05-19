package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"fortyfour-backend/pkg/logger"
	_ "survey/docs"
	"survey/internal/cache"
	"survey/internal/config"
	"survey/internal/handlers"
	"survey/internal/middleware"
	"survey/internal/repository"
	"survey/internal/routes"
	"survey/internal/services"
	"survey/pkg/database"

	"github.com/joho/godotenv"
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
	// Load env
	if err := godotenv.Load(); err != nil {
		os.Stdout.WriteString("Warning: .env file not found\n")
	}

	cfg := config.Load()

	// Initialize structured logger
	logger.Init(cfg.LogLevel, cfg.Environment)
	logger.Info("Logger initialized successfully")

	// Init DB
	db, err := database.NewMySQLConnection(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	})
	if err != nil {
		logger.FatalErr(err, "Failed to connect to database")
	}
	defer db.Close()

	// INIT REDIS
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// test connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.FatalErr(err, "Redis connection failed")
	}
	logger.Info("Redis initialized successfully")

	// REPOSITORY
	respondenRepo := repository.NewRespondenRepository(db)
	risikoRepo := repository.NewRisikoRepository(db)

	// SERVICE
	redisCache := cache.NewRedisCache(rdb)
	respondenService := services.NewRespondenService(respondenRepo, services.DefaultValidator{}, redisCache)
	risikoService := services.NewRisikoService(risikoRepo, redisCache)

	// HANDLER
	respondenHandler := handlers.NewRespondenHandler(respondenService)
	risikoHandler := handlers.NewRisikoHandler(risikoService)

	// MIDDLEWARE
	authMiddleware := middleware.NewAuthMiddleware(cfg.InternalGatewayKey)

	// ROUTER
	mux := routes.InitRouter(
		respondenHandler,
		risikoHandler,
		authMiddleware,
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
		logger.Infof("Survey API running on port %s", cfg.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.FatalErr(err, "Server failed")
		}
	}()

	// GRACEFUL SHUTDOWN
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctxShutdown); err != nil {
		logger.FatalErr(err, "Server forced to shutdown")
	}

	logger.Info("Server exited properly")
}
