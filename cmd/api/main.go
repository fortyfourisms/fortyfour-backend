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
	"fortyfour-backend/internal/seeder"
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
	// if err := database.RunMigrations(database.Config{
	// 	Host:     cfg.Database.Host,
	// 	Port:     cfg.Database.Port,
	// 	User:     cfg.Database.User,
	// 	Password: cfg.Database.Password,
	// 	DBName:   cfg.Database.DBName,
	// }, "./migrations"); err != nil {
	// 	logger.FatalErr(err, "Failed to run database migrations")
	// }

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

	// Initialize SSE Service
	sseService := services.NewSSEService()
	logger.Info("SSE Service initialized successfully")

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
	usersConsumer := usersRmq.NewConsumer(sharedConsumer, sseService)

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

	// Seed default Casbin policies (aman dijalankan berulang kali)
	seeder.SeedCasbinPolicies(casbinService)
	logger.Info("Casbin policies seeded")

	// Initialize Gemini Client
	geminiClient := utils.NewGeminiClient(cfg.GeminiAPIKey)
	logger.Info("Gemini client initialized successfully")

	uploadPath := os.Getenv("UPLOAD_DIR")
	if uploadPath == "" {
		uploadPath = "/app/uploads" // default fallback
	}
	os.MkdirAll(uploadPath, os.ModePerm)

	// Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	perusahaanRepo := repository.NewPerusahaanRepository(db)
	picRepo := repository.NewPICRepository(db)
	jabatanRepo := repository.NewJabatanRepository(db)
	csirtRepo := repository.NewCsirtRepository(db)
	sdmCsirtRepo := repository.NewSdmCsirtRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	chatRepo := repository.NewInMemoryChatRepo()
	sektorRepo := repository.NewSektorRepository(db)
	subSektorRepo := repository.NewSubSektorRepository(db)
	seRepo := repository.NewSERepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)
	kelasRepo := repository.NewKelasRepository(db)
	materiRepo := repository.NewMateriRepository(db)
	soalRepo := repository.NewSoalRepository(db)
	kuisRepo := repository.NewKuisAttemptRepository(db)
	progressRepo := repository.NewProgressRepository(db)

	// Initialize services
	tokenService := services.NewTokenService(redisClient, cfg.JWTSecret, true, cfg.Domain)
	notificationService := services.NewNotificationService(redisClient)
	strExpiryService := services.NewSTRExpiryService(csirtRepo, notificationService)
	authService := services.NewAuthService(userRepo, roleRepo, tokenService, notificationService, strExpiryService)
	perusahaanService := services.NewPerusahaanService(perusahaanRepo, subSektorRepo, redisClient)
	picService := services.NewPICService(picRepo, redisClient)
	jabatanService := services.NewJabatanService(jabatanRepo, redisClient)
	csirtService := services.NewCsirtService(csirtRepo, redisClient)
	csirtExportService := services.NewCsirtExportService(csirtService)
	sdmCsirtService := services.NewSdmCsirtService(sdmCsirtRepo, redisClient)
	userService := services.NewUserService(userRepo, uploadPath, usersProducer)
	roleService := services.NewRoleService(roleRepo, redisClient)
	chatService := services.NewChatService(chatRepo, geminiClient, db)
	sektorService := services.NewSektorService(sektorRepo, redisClient)
	subSektorService := services.NewSubSektorService(subSektorRepo, redisClient)
	seService := services.NewSEService(seRepo, redisClient)
	seExportService := services.NewSEExportService(seService)
	dashboardService := services.NewDashboardService(dashboardRepo, redisClient)
	kelasSvc := services.NewKelasService(kelasRepo, materiRepo, progressRepo, redisClient)
	materiSvc := services.NewMateriService(materiRepo, kelasRepo, progressRepo, redisClient)
	soalSvc := services.NewSoalService(soalRepo, materiRepo, redisClient)
	kuisSvc := services.NewKuisService(kuisRepo, soalRepo, materiRepo, progressRepo, redisClient)

	// Initialize Handler
	authHandler := handlers.NewAuthHandler(authService, tokenService, perusahaanService, userService, uploadPath)
	userHandler := handlers.NewUserHandler(userService, uploadPath, sseService)
	perusahaanHandler := handlers.NewPerusahaanHandler(perusahaanService, uploadPath, sseService)
	picHandler := handlers.NewPICHandler(picService, sseService)
	jabatanHandler := handlers.NewJabatanHandler(jabatanService, sseService)
	csirtHandler := handlers.NewCsirtHandler(csirtService, sseService)
	csirtExportHandler := handlers.NewCsirtExportHandler(csirtExportService)
	csirtHandler.SetExportHandler(csirtExportHandler)
	sdmCsirtHandler := handlers.NewSdmCsirtHandler(sdmCsirtService, csirtService, sseService)
	roleHandler := handlers.NewRoleHandler(roleService, sseService)
	casbinHandler := handlers.NewCasbinHandler(casbinService, sseService)
	sseHandler := handlers.NewSSEHandler(sseService)
	chatHandler := handlers.NewChatHandler(chatService)
	sektorHandler := handlers.NewSektorHandler(sektorService)
	subSektorHandler := handlers.NewSubSektorHandler(subSektorService)
	seHandler := handlers.NewSEHandler(seService, sseService)
	seExportHandler := handlers.NewSEExportHandler(seExportService)
	seHandler.SetExportHandler(seExportHandler)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	lmsHandler := handlers.NewLMSHandler(kelasSvc, materiSvc, soalSvc, kuisSvc, sseService)

	// Proxy Handler for IKAS
	ikasProxyHandler := handlers.NewProxyHandler("http://ikas:8081", cfg.InternalGatewayKey)

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
		roleHandler,
		casbinHandler,
		sseHandler,
		authMiddleware,
		casbinMiddleware,
		strictLimiter,
		moderateLimiter,
		lenientLimiter,
		csirtHandler,
		csirtExportHandler,
		sdmCsirtHandler,
		chatHandler,
		sektorHandler,
		subSektorHandler,
		seHandler,
		seExportHandler,
		dashboardHandler,
		notificationHandler,
		ikasProxyHandler,
		lmsHandler,
	)

	// Start server
	logger.Infof("Server starting on %s", cfg.Port)
	logger.Info("Rate limiting enabled:")
	logger.Info("  - Auth endpoints: 10 requests/minute per IP")
	logger.Info("  - Public posts: 1000 requests/minute per IP")
	logger.Info("  - Protected posts: 100 requests/minute per user")
	logger.Info("SSE endpoint available at /api/events")

	logger.FatalErr(http.ListenAndServe(cfg.Port, mux), "Server stopped")
}
