package main

import (
	"log"
	"net/http"

	"fortyfour-backend/internal/config"
	"fortyfour-backend/internal/handlers"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/routes"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/pkg/database"
)

func main() {
	cfg := config.Load()

	// Initialize MySQL database
	db, err := database.NewMySQLConnection(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize Redis
	// redisClient, err := cache.NewRedisClient(cache.RedisConfig{
	// 	Host:     cfg.Redis.Host,
	// 	Port:     cfg.Redis.Port,
	// 	Password: cfg.Redis.Password,
	// 	DB:       cfg.Redis.DB,
	// })
	// if err != nil {
	// 	log.Fatal("Failed to connect to Redis:", err)
	// }
	// defer redisClient.Close()

	// Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	perusahaanRepo := repository.NewPerusahaanRepository(db)
	picRepo := repository.NewPICPerusahaanRepository(db)
	identifikasiRepo := repository.NewIdentifikasiRepository(db)
	proteksiRepo := repository.NewProteksiRepository(db)
	deteksiRepo := repository.NewDeteksiRepository(db)
	gulihRepo := repository.NewGulihRepository(db)
	ikasRepo := repository.NewIkasRepository(db)
	
	// Initialize services
	tokenService := services.NewTokenService(nil, cfg.JWTSecret)
	authService := services.NewAuthService(userRepo, tokenService)
	postService := services.NewPostService(postRepo)
	perusahaanService := services.NewPerusahaanService(perusahaanRepo)
	picService := services.NewPICPerusahaanService(picRepo)
	identifikasiService := services.NewIdentifikasiService(identifikasiRepo)
	proteksiService := services.NewProteksiService(proteksiRepo)
	deteksiService := services.NewDeteksiService(deteksiRepo)
	gulihService := services.NewGulihService(gulihRepo)
	ikasService := services.NewIkasService(ikasRepo)

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService, tokenService)
	postHandler := handlers.NewPostHandler(postService)
	perusahaanHandler := handlers.NewPerusahaanHandler(perusahaanService)
	picHandler := handlers.NewPICPerusahaanHandler(picService)
	identifikasiHandler := handlers.NewIdentifikasiHandler(identifikasiService)
	proteksiHandler := handlers.NewProteksiHandler(proteksiService)
	deteksiHandler := handlers.NewDeteksiHandler(deteksiService)
	gulihHandler := handlers.NewGulihHandler(gulihService)
	ikasHandler := handlers.NewIkasHandler(ikasService)

	// Initialize Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// Setup routes
	mux := routes.InitRouter(authHandler, postHandler, perusahaanHandler, picHandler, identifikasiHandler, proteksiHandler, deteksiHandler, gulihHandler, ikasHandler, authMiddleware)

	// Start server
	log.Printf("Server starting on %s", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}