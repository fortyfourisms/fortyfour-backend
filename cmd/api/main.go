package main

import (
	"log"
	"net/http"

	"fortyfour-backend/internal/config"
	"fortyfour-backend/internal/handlers"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/pkg/database"
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
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	postService := services.NewPostService(postRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	postHandler := handlers.NewPostHandler(postService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// Setup routes
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/api/register", authHandler.Register)
	mux.HandleFunc("/api/login", authHandler.Login)
	mux.HandleFunc("/api/posts", postHandler.GetPosts)
	mux.HandleFunc("/api/posts/single", postHandler.GetPost)

	// Protected routes
	mux.HandleFunc("/api/posts/create", authMiddleware.Authenticate(postHandler.CreatePost))
	mux.HandleFunc("/api/posts/update", authMiddleware.Authenticate(postHandler.UpdatePost))
	mux.HandleFunc("/api/posts/delete", authMiddleware.Authenticate(postHandler.DeletePost))

	// Start server
	log.Printf("Server starting on %s", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}
