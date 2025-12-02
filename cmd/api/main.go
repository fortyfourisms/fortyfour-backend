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

	"github.com/gorilla/mux"
	
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
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	ikasRepo := repository.NewIkasRepository(db)

	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	postService := services.NewPostService(postRepo)
	ikasService := services.NewIkasService(ikasRepo)

	authHandler := handlers.NewAuthHandler(authService)
	postHandler := handlers.NewPostHandler(postService)
	ikasHandler := handlers.NewIkasHandler(ikasService)

	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	router := mux.NewRouter()
	//mux := http.NewServeMux()

	router.HandleFunc("/api/register", authHandler.Register)
	router.HandleFunc("/api/login", authHandler.Login)
	router.HandleFunc("/api/posts", postHandler.GetPosts)
	router.HandleFunc("/api/posts/single", postHandler.GetPost)
	
	router.HandleFunc("/api/posts/create", authMiddleware.Authenticate(postHandler.CreatePost))
	router.HandleFunc("/api/posts/update", authMiddleware.Authenticate(postHandler.UpdatePost))
	router.HandleFunc("/api/posts/delete", authMiddleware.Authenticate(postHandler.DeletePost))
 
	router.HandleFunc("/api/ikas", authMiddleware.Authenticate(ikasHandler.CreateIkas)).Methods("POST")
	router.HandleFunc("/api/ikas", ikasHandler.GetAllIkas).Methods("GET")
	router.HandleFunc("/api/ikas/{id}", ikasHandler.GetIkasByID).Methods("GET")
	router.HandleFunc("/api/ikas/{id}", authMiddleware.Authenticate(ikasHandler.UpdateIkas)).Methods("PUT")
	router.HandleFunc("/api/ikas/{id}", authMiddleware.Authenticate(ikasHandler.DeleteIkas)).Methods("DELETE")

	router.PathPrefix("/").Handler(router)

	log.Printf("Server starting on %s", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, router))
}
