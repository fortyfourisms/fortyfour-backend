package main

import (
	"log"
	"net/http"
	"os"

	_ "fortyfour-backend/docs"
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

	// Initialize Casbin Service with GORM Adapter
	casbinService, err := services.NewCasbinService(cfg.Database.GetDSN(), cfg.CasbinModelPath)
	if err != nil {
		log.Fatal("Failed to initialize Casbin:", err)
	}
	log.Println("Casbin RBAC initialized successfully with GORM adapter")

	// Initialize SSE Service
	sseService := services.NewSSEService()
	log.Println("SSE Service initialized successfully")

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
	sdmCsirtRepo := repository.NewSdmCsirtRepository(db)
	seCsirtRepo := repository.NewSeCsirtRepository(db)
	roleRepo := repository.NewRoleRepository(db)

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
	sdmCsirtService := services.NewSdmCsirtService(sdmCsirtRepo)
	seCsirtService := services.NewSeCsirtService(seCsirtRepo)
	userService := services.NewUserService(userRepo, "./uploads")
	roleService := services.NewRoleService(roleRepo)

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService, tokenService)
	userHandler := handlers.NewUserHandler(userService, "./uploads", sseService)
	postHandler := handlers.NewPostHandler(postService, sseService)
	uploadPath := "./uploads"
	os.MkdirAll(uploadPath, os.ModePerm)
	perusahaanHandler := handlers.NewPerusahaanHandler(perusahaanService, uploadPath, sseService)
	picHandler := handlers.NewPICHandler(picService, sseService)
	identifikasiHandler := handlers.NewIdentifikasiHandler(identifikasiService, sseService)
	jabatanHandler := handlers.NewJabatanHandler(jabatanService, sseService)
	proteksiHandler := handlers.NewProteksiHandler(proteksiService, sseService)
	deteksiHandler := handlers.NewDeteksiHandler(deteksiService, sseService)
	gulihHandler := handlers.NewGulihHandler(gulihService, sseService)
	ikasHandler := handlers.NewIkasHandler(ikasService, sseService)
	csirtHandler := handlers.NewCsirtHandler(csirtService)
	sdmCsirtHandler := handlers.NewSdmCsirtHandler(sdmCsirtService)
	seCsirtHandler := handlers.NewSeCsirtHandler(seCsirtService)
	roleHandler := handlers.NewRoleHandler(roleService, sseService)
	casbinHandler := handlers.NewCasbinHandler(casbinService, sseService)
	sseHandler := handlers.NewSSEHandler(sseService)

	// Initialize Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)
	casbinMiddleware := middleware.NewCasbinMiddleware(casbinService.GetEnforcer())

	// Initialize rate limiters with different configurations
	rateLimitConfigs := middleware.GetRateLimitConfigs()

	strictLimiter := middleware.NewRateLimiter(redisClient, rateLimitConfigs.Strict)
	moderateLimiter := middleware.NewRateLimiter(redisClient, rateLimitConfigs.Moderate)
	lenientLimiter := middleware.NewRateLimiter(redisClient, rateLimitConfigs.Lenient)

	// Setup routes
	mux := routes.InitRouter(
		authHandler,
		userHandler,
		postHandler,
		perusahaanHandler,
		picHandler,
		jabatanHandler,
		identifikasiHandler,
		deteksiHandler,
		gulihHandler,
		ikasHandler,
		proteksiHandler,
		roleHandler,
		casbinHandler,
		sseHandler,
		authMiddleware,
		casbinMiddleware,
		strictLimiter,
		moderateLimiter,
		lenientLimiter,
		csirtHandler,
		sdmCsirtHandler,
		seCsirtHandler,
	)

	// Start server
	log.Printf("Server starting on %s", cfg.Port)
	log.Println("Rate limiting enabled:")
	log.Println("  - Auth endpoints: 5 requests/minute per IP")
	log.Println("  - Public posts: 1000 requests/minute per IP")
	log.Println("  - Protected posts: 100 requests/minute per user")
	log.Println("SSE endpoint available at /api/events")

	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}
