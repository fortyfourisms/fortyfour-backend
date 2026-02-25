package main

import (
	"context"
	"net/http"
	"os"

	_ "fortyfour-backend/docs"
	"fortyfour-backend/internal/config"
	"fortyfour-backend/internal/handlers"
	"fortyfour-backend/internal/middleware"
	usersRmq "fortyfour-backend/internal/rabbitmq"
	"fortyfour-backend/internal/repository"
	"fortyfour-backend/internal/routes"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"fortyfour-backend/pkg/cache"
	"fortyfour-backend/pkg/database"
	"fortyfour-backend/pkg/logger"
	pkgRmq "fortyfour-backend/pkg/rabbitmq"

	"github.com/joho/godotenv"
)

// @title Fortyfour Backend API
// @version 1.0
// @description API documentation for Fortyfour Backend - main auth and management service.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	// Load env
	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found, using system env")
	}

	cfg := config.Load()

	// Initialize structured logger
	logger.Init(cfg.LogLevel, cfg.Environment)
	logger.Info("Logger initialized successfully")

	// Initialize MySQL database
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

	// Run database migrations
	if err := database.RunMigrations(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	}, "./migrations"); err != nil {
		logger.FatalErr(err, "Failed to run database migrations")
	}

	// Initialize Redis
	redisClient, err := cache.NewRedisClient(cache.RedisConfig{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		logger.FatalErr(err, "Failed to connect to Redis")
	}
	defer redisClient.Close()

	// Initialize RabbitMQ
	rmq, err := pkgRmq.NewRabbitMQ(cfg.RabbitMQ.GetURL())
	if err != nil {
		logger.FatalErr(err, "Failed to connect to RabbitMQ")
	}
	defer rmq.Close()
	logger.Info("RabbitMQ initialized successfully")

	// Setup RabbitMQ infrastructure (Users)
	if err := usersRmq.SetupInfrastructure(rmq); err != nil {
		logger.FatalErr(err, "Failed to setup Users RabbitMQ infrastructure")
	}

	logger.Info("RabbitMQ infrastructure initialized successfully")

	// Create Shared Producer and Consumer
	sharedProducer := pkgRmq.NewProducer(rmq.GetChannel())
	sharedConsumer := pkgRmq.NewConsumer(rmq.GetChannel())

	// Wrap with specific Producer
	usersProducer := usersRmq.NewProducer(sharedProducer)

	// Wrap with specific Consumer
	usersConsumer := usersRmq.NewConsumer(sharedConsumer)

	// Start consumers in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := usersConsumer.StartAllConsumers(ctx); err != nil {
		logger.FatalErr(err, "Failed to start Users consumers")
	}

	// Initialize Casbin Service with GORM Adapter
	casbinService, err := services.NewCasbinService(cfg.Database.GetDSN(), cfg.CasbinModelPath)
	if err != nil {
		logger.FatalErr(err, "Failed to initialize Casbin")
	}
	logger.Info("Casbin RBAC initialized successfully with GORM adapter")

	// Initialize SSE Service
	sseService := services.NewSSEService()
	logger.Info("SSE Service initialized successfully")

	// Initialize Gemini Client
	geminiClient := utils.NewGeminiClient(cfg.GeminiAPIKey)
	logger.Info("Gemini client initialized successfully")

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
	roleRepo := repository.NewRoleRepository(db)
	chatRepo := repository.NewInMemoryChatRepo()
	sektorRepo := repository.NewSektorRepository(db)
	subSektorRepo := repository.NewSubSektorRepository(db)
	seRepo := repository.NewSERepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)

	// Initialize services
	tokenService := services.NewTokenService(redisClient, cfg.JWTSecret, true, cfg.Domain)
	authService := services.NewAuthService(userRepo, tokenService)
	perusahaanService := services.NewPerusahaanService(perusahaanRepo, subSektorRepo)
	picService := services.NewPICService(picRepo, redisClient)
	identifikasiService := services.NewIdentifikasiService(identifikasiRepo)
	jabatanService := services.NewJabatanService(jabatanRepo)
	proteksiService := services.NewProteksiService(proteksiRepo)
	deteksiService := services.NewDeteksiService(deteksiRepo)
	gulihService := services.NewGulihService(gulihRepo)
	csirtService := services.NewCsirtService(csirtRepo, redisClient)
	sdmCsirtService := services.NewSdmCsirtService(sdmCsirtRepo, redisClient)
	userService := services.NewUserService(userRepo, "./uploads", usersProducer)
	roleService := services.NewRoleService(roleRepo)
	chatService := services.NewChatService(chatRepo, geminiClient, db)
	sektorService := services.NewSektorService(sektorRepo)
	subSektorService := services.NewSubSektorService(subSektorRepo)
	seService := services.NewSEService(seRepo, redisClient)
	dashboardService := services.NewDashboardService(dashboardRepo)

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService, tokenService, perusahaanService)
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
	roleHandler := handlers.NewRoleHandler(roleService, sseService)
	casbinHandler := handlers.NewCasbinHandler(casbinService, sseService)
	sseHandler := handlers.NewSSEHandler(sseService)
	chatHandler := handlers.NewChatHandler(chatService)
	sektorHandler := handlers.NewSektorHandler(sektorService)
	subSektorHandler := handlers.NewSubSektorHandler(subSektorService)
	seHandler := handlers.NewSEHandler(seService, sseService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	// Initialize Middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService)
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
		chatHandler,
		sektorHandler,
		subSektorHandler,
		seHandler,
		dashboardHandler,
	)

	// Start server
	logger.Infof("Server starting on %s", cfg.Port)
	logger.Info("Rate limiting enabled:")
	logger.Info("  - Auth endpoints: 5 requests/minute per IP")
	logger.Info("  - Public posts: 1000 requests/minute per IP")
	logger.Info("  - Protected posts: 100 requests/minute per user")
	logger.Info("SSE endpoint available at /api/events")

	logger.FatalErr(http.ListenAndServe(cfg.Port, mux), "Server stopped")
}
