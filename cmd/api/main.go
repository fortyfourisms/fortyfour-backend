package main

import (
	"log"
	"net/http"
	"os"

	"fortyfour-backend/internal/config"
	"fortyfour-backend/internal/handlers"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/routes"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/pkg/cache"
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
	redisClient, err := cache.NewRedisClient(cache.RedisConfig{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisClient.Close()

	// Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	perusahaanRepo := repository.NewPerusahaanRepository(db)
	picRepo := repository.NewPICRepository(db)
	identifikasiRepo := repository.NewIdentifikasiRepository(db)
	jabatanRepo := repository.NewJabatanRepository(db)
	ikasRepo := repository.NewIkasRepository(db)
	proteksiRepo := repository.NewProteksiRepository(db)
	gulihRepo := repository.NewGulihRepository(db)

	// Initialize services
	tokenService := services.NewTokenService(redisClient, cfg.JWTSecret)
	authService := services.NewAuthService(userRepo, tokenService)
	postService := services.NewPostService(postRepo)
	perusahaanService := services.NewPerusahaanService(perusahaanRepo)
	picService := services.NewPICService(picRepo)
	identifikasiService := services.NewIdentifikasiService(identifikasiRepo)
	jabatanService := services.NewJabatanService(jabatanRepo)
	ikasService := services.NewIkasService(ikasRepo)
	proteksiService := services.NewProteksiService(proteksiRepo)
	gulihService := services.NewGulihService(gulihRepo)

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService, tokenService)
	postHandler := handlers.NewPostHandler(postService)
	uploadPath := "./uploads"
	os.MkdirAll(uploadPath, os.ModePerm)
	perusahaanHandler := handlers.NewPerusahaanHandler(perusahaanService, uploadPath)
	picHandler := handlers.NewPICHandler(picService)
	identifikasiHandler := handlers.NewIdentifikasiHandler(identifikasiService)
	jabatanHandler := handlers.NewJabatanHandler(jabatanService)
	ikasHandler := handlers.NewIkasHandler(ikasService)
	proteksiHandler := handlers.NewProteksiHandler(proteksiService)
	gulihHandler := handlers.NewGulihHandler(gulihService)

	// Initialize Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// Setup routes
	mux := routes.InitRouter(authHandler, postHandler, perusahaanHandler, picHandler, identifikasiHandler, jabatanHandler, gulihHandler, ikasHandler, proteksiHandler, authMiddleware)

	// Start server
	log.Printf("Server starting on %s", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}
