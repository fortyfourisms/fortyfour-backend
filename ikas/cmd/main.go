package main

import (
	"ikas/internal/config"
	"ikas/internal/handlers"
	"ikas/internal/middleware"
	"ikas/internal/repository"
	"ikas/internal/routes"
	"ikas/internal/services"
	"ikas/pkg/cache"
	"ikas/pkg/database"
	"log"
	"net/http"

	"github.com/rollbar/rollbar-go"
)

func main() {

	// Load env
	// if err := godotenv.Load(); err != nil {
	// 	log.Println("No .env file found")
	// }

	cfg := config.Load()

	rollbar.SetToken(cfg.Rollbar.Token)
	rollbar.SetEnvironment(cfg.Rollbar.Env)
	// Send a test message
	rollbar.Info("Rollbar Go SDK initialized successfully!")

	// call rollbar.Close() before the application exits to flush error message queue
	// rollbar.Close()

	// Initialize MySQL database
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

	// Initialize Redis
	redisClient, err := cache.NewRedisClient(cache.RedisConfig{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		log.Println("Failed to connect to Redis:", err)
		rollbar.Error(err)
		log.Fatal(err)
	}
	log.Println("Redis initialized successfully")
	rollbar.Info("Redis initialized successfully")
	defer redisClient.Close()

	// repository
	ikasRepo := repository.NewIkasRepository(db)
	ruangLingkupRepo := repository.NewRuangLingkupRepository(db)
	domainRepo := repository.NewDomainRepository(db)
	kategoriRepo := repository.NewKategoriRepository(db)
	subKategoriRepo := repository.NewSubKategoriRepository(db)

	// services
	ikasService := services.NewIkasService(ikasRepo)
	ruangLingkupService := services.NewRuangLingkupService(ruangLingkupRepo)
	domainService := services.NewDomainService(domainRepo)
	kategoriService := services.NewKategoriService(kategoriRepo)
	subKategoriService := services.NewSubKategoriService(subKategoriRepo)

	// handlers
	ikasHandler := handlers.NewIkasHandler(ikasService)
	ruangLingkupHandler := handlers.NewRuangLingkupHandler(ruangLingkupService)
	domainHandler := handlers.NewDomainHandler(domainService)
	kategoriHandler := handlers.NewKategoriHandler(kategoriService)
	subKategoriHandler := handlers.NewSubKategoriHandler(subKategoriService)

	// Router
	mux := routes.InitRouter(
		ikasHandler,
		ruangLingkupHandler,
		domainHandler,
		kategoriHandler,
		subKategoriHandler,
	)

	auth := middleware.NewAuthMiddleware(cfg.JWTSecret)
	handler := auth.Authenticate(mux)

	log.Println("IKAS service running on", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, handler))
}
