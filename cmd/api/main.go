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
	"fortyfour-backend/internal/utils"
	"fortyfour-backend/pkg/cache"
	"fortyfour-backend/pkg/database"

	"github.com/joho/godotenv"
	"github.com/rollbar/rollbar-go"
	"github.com/rs/cors"
)

func main() {

	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

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
		log.Fatal("Failed to connect to database:", err)
		rollbar.Error(err)
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
		rollbar.Error(err)
	}
	defer redisClient.Close()

	// Initialize Casbin Service with GORM Adapter
	casbinService, err := services.NewCasbinService(cfg.Database.GetDSN(), cfg.CasbinModelPath)
	if err != nil {
		log.Fatal("Failed to initialize Casbin:", err)
		rollbar.Error(err)
	}
	log.Println("Casbin RBAC initialized successfully with GORM adapter")
	rollbar.Info("Casbin RBAC initialized successfully with GORM adapter")

	// Initialize SSE Service
	sseService := services.NewSSEService()
	log.Println("SSE Service initialized successfully")
	rollbar.Info("SSE Service initialized successfully")

	// Initialize Gemini Client
	geminiClient := utils.NewGeminiClient()
	if err != nil {
		log.Fatal("Failed to initialize Gemini client:", err)
	}
	log.Println("Gemini client initialized successfully")

	// Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	perusahaanRepo := repository.NewPerusahaanRepository(db)
	picRepo := repository.NewPICRepository(db)
	identifikasiRepo := repository.NewIdentifikasiRepository(db)
	jabatanRepo := repository.NewJabatanRepository(db)
	proteksiRepo := repository.NewProteksiRepository(db)
	deteksiRepo := repository.NewDeteksiRepository(db)
	gulihRepo := repository.NewGulihRepository(db)
	csirtRepo := repository.NewCsirtRepository(db)
	sdmCsirtRepo := repository.NewSdmCsirtRepository(db)
	seCsirtRepo := repository.NewSeCsirtRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	chatRepo := repository.NewInMemoryChatRepo()

	// Initialize services
	tokenService := services.NewTokenService(redisClient, cfg.JWTSecret)
	authService := services.NewAuthService(userRepo, tokenService)
	perusahaanService := services.NewPerusahaanService(perusahaanRepo)
	picService := services.NewPICService(picRepo)
	identifikasiService := services.NewIdentifikasiService(identifikasiRepo)
	jabatanService := services.NewJabatanService(jabatanRepo)
	proteksiService := services.NewProteksiService(proteksiRepo)
	deteksiService := services.NewDeteksiService(deteksiRepo)
	gulihService := services.NewGulihService(gulihRepo)
	csirtService := services.NewCsirtService(csirtRepo)
	sdmCsirtService := services.NewSdmCsirtService(sdmCsirtRepo)
	seCsirtService := services.NewSeCsirtService(seCsirtRepo)
	userService := services.NewUserService(userRepo, "./uploads")
	roleService := services.NewRoleService(roleRepo)
	chatService := services.NewChatService(chatRepo, geminiClient, db)

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService, tokenService)
	userHandler := handlers.NewUserHandler(userService, "./uploads", sseService)
	uploadPath := "./uploads"
	os.MkdirAll(uploadPath, os.ModePerm)
	perusahaanHandler := handlers.NewPerusahaanHandler(perusahaanService, uploadPath, sseService)
	picHandler := handlers.NewPICHandler(picService, sseService)
	identifikasiHandler := handlers.NewIdentifikasiHandler(identifikasiService, sseService)
	jabatanHandler := handlers.NewJabatanHandler(jabatanService, sseService)
	proteksiHandler := handlers.NewProteksiHandler(proteksiService, sseService)
	deteksiHandler := handlers.NewDeteksiHandler(deteksiService, sseService)
	gulihHandler := handlers.NewGulihHandler(gulihService, sseService)
	csirtHandler := handlers.NewCsirtHandler(csirtService)
	sdmCsirtHandler := handlers.NewSdmCsirtHandler(sdmCsirtService)
	seCsirtHandler := handlers.NewSeCsirtHandler(seCsirtService)
	roleHandler := handlers.NewRoleHandler(roleService, sseService)
	casbinHandler := handlers.NewCasbinHandler(casbinService, sseService)
	sseHandler := handlers.NewSSEHandler(sseService)
	chatHandler := handlers.NewChatHandler(chatService)

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
		perusahaanHandler,
		picHandler,
		jabatanHandler,
		identifikasiHandler,
		deteksiHandler,
		gulihHandler,
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
		chatHandler,
	)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://admin.kssindustri.site", "https://fortyfouris.netlify.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	mCors := c.Handler(mux)

	// Start server
	log.Printf("Server starting on %s", cfg.Port)
	rollbar.Info("Server starting on %s", cfg.Port)
	log.Println("Rate limiting enabled:")
	log.Println("  - Auth endpoints: 5 requests/minute per IP")
	log.Println("  - Public posts: 1000 requests/minute per IP")
	log.Println("  - Protected posts: 100 requests/minute per user")
	log.Println("SSE endpoint available at /api/events")

	log.Fatal(http.ListenAndServe(cfg.Port, mCors))

	// Ensure all items are sent before the app exits
	rollbar.Wait()
}
