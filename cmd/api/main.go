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
	proteksiRepo := repository.NewProteksiRepository(db)
	deteksiRepo := repository.NewDeteksiRepository(db)
	gulihRepo := repository.NewGulihRepository(db)
	ikasRepo := repository.NewIkasRepository(db)
	csirtRepo := repository.NewCsirtRepository(db)

	// Initialize services
	tokenService := services.NewTokenService(redisClient, cfg.JWTSecret)
	authService := services.NewAuthService(userRepo, tokenService)
	postService := services.NewPostService(postRepo)
	perusahaanService := services.NewPerusahaanService(perusahaanRepo)
	picService := services.NewPICService(picRepo)
	identifikasiService := services.NewIdentifikasiService(identifikasiRepo)
	jabatanService := services.NewJabatanService(jabatanRepo)
	proteksiService := services.NewProteksiService(proteksiRepo)
	deteksiService := services.NewDeteksiService(deteksiRepo)
	gulihService := services.NewGulihService(gulihRepo)
	ikasService := services.NewIkasService(ikasRepo)
	csirtService := services.NewCsirtService(csirtRepo)

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService, tokenService)
	postHandler := handlers.NewPostHandler(postService)
	uploadPath := "./uploads"
	os.MkdirAll(uploadPath, os.ModePerm)
	perusahaanHandler := handlers.NewPerusahaanHandler(perusahaanService, uploadPath)
	picHandler := handlers.NewPICHandler(picService)
	identifikasiHandler := handlers.NewIdentifikasiHandler(identifikasiService)
	jabatanHandler := handlers.NewJabatanHandler(jabatanService)
	proteksiHandler := handlers.NewProteksiHandler(proteksiService)
	deteksiHandler := handlers.NewDeteksiHandler(deteksiService)
	gulihHandler := handlers.NewGulihHandler(gulihService)
	ikasHandler := handlers.NewIkasHandler(ikasService)
	csirtHandler := handlers.NewCsirtHandler(csirtService)

	// Initialize Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// Initialize rate limiters with different configurations
	rateLimitConfigs := middleware.GetRateLimitConfigs()

	strictLimiter := middleware.NewRateLimiter(redisClient, rateLimitConfigs.Strict)
	moderateLimiter := middleware.NewRateLimiter(redisClient, rateLimitConfigs.Moderate)
	lenientLimiter := middleware.NewRateLimiter(redisClient, rateLimitConfigs.Lenient)

	// Setup routes
	mux := routes.InitRouter(
		authHandler,
		postHandler,
		perusahaanHandler,
		picHandler,
		jabatanHandler,
		identifikasiHandler,
		deteksiHandler,
		gulihHandler,
		ikasHandler,
		proteksiHandler,
		authMiddleware,
		strictLimiter,
		moderateLimiter,
		lenientLimiter,
		csirtHandler,
	)

	// Start server
	log.Printf("Server starting on %s", cfg.Port)
	log.Println("Rate limiting enabled:")
	log.Println("  - Auth endpoints: 5 requests/minute per IP")
	log.Println("  - Public posts: 1000 requests/minute per IP")
	log.Println("  - Protected posts: 100 requests/minute per user")

	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}
