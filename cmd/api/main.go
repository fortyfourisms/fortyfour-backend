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

	"github.com/rollbar/rollbar-go"
)

func main() {
	cfg := config.Load()

	// Debug: check what Rollbar token was loaded
	rawToken := cfg.Rollbar.Token
	log.Printf("ROLLBAR_TOKEN loaded: len=%d, first4=%q, raw_env=%q", len(rawToken), rawToken[:min(4, len(rawToken))], os.Getenv("ROLLBAR_TOKEN"))

	rollbar.SetToken(cfg.Rollbar.Token)
	rollbar.SetEnvironment(cfg.Rollbar.Env)

	// Ensure all items are sent before the app exits
	defer rollbar.Wait()
	// Close Rollbar client when the application exits
	defer rollbar.Close()

	// Send a test message
	rollbar.Info("Rollbar Go SDK initialized successfully!")

	// Initialize MySQL database
	db, err := database.NewMySQLConnection(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	})
	if err != nil {
		rollbar.Error(err)
		rollbar.Wait()
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
		rollbar.Error(err)
		rollbar.Wait()
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisClient.Close()

	// Initialize Casbin Service with GORM Adapter
	casbinService, err := services.NewCasbinService(cfg.Database.GetDSN(), cfg.CasbinModelPath)
	if err != nil {
		rollbar.Error(err)
		rollbar.Wait()
		log.Fatal("Failed to initialize Casbin:", err)
	}
	log.Println("Casbin RBAC initialized successfully with GORM adapter")
	rollbar.Info("Casbin RBAC initialized successfully with GORM adapter")

	// Initialize SSE Service
	sseService := services.NewSSEService()
	log.Println("SSE Service initialized successfully")
	rollbar.Info("SSE Service initialized successfully")

	// Initialize Repositories
	userRepo := repository.NewUserRepository(db)
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
	roleRepo := repository.NewRoleRepository(db)
	sektorRepo := repository.NewSektorRepository(db)
	subSektorRepo := repository.NewSubSektorRepository(db)
	seRepo := repository.NewSERepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)

	// Initialize services
	tokenService := services.NewTokenService(redisClient, cfg.JWTSecret, true, cfg.Domain)
	authService := services.NewAuthService(userRepo, tokenService)
	perusahaanService := services.NewPerusahaanService(perusahaanRepo, subSektorRepo)
	picService := services.NewPICService(picRepo)
	identifikasiService := services.NewIdentifikasiService(identifikasiRepo)
	jabatanService := services.NewJabatanService(jabatanRepo)
	proteksiService := services.NewProteksiService(proteksiRepo)
	deteksiService := services.NewDeteksiService(deteksiRepo)
	gulihService := services.NewGulihService(gulihRepo)
	ikasService := services.NewIkasService(ikasRepo)
	csirtService := services.NewCsirtService(csirtRepo)
	sdmCsirtService := services.NewSdmCsirtService(sdmCsirtRepo)
	userService := services.NewUserService(userRepo, "./uploads")
	roleService := services.NewRoleService(roleRepo)
	sektorService := services.NewSektorService(sektorRepo)
	subSektorService := services.NewSubSektorService(subSektorRepo)
	seService := services.NewSEService(seRepo)
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
	ikasHandler := handlers.NewIkasHandler(ikasService, sseService)
	csirtHandler := handlers.NewCsirtHandler(csirtService)
	sdmCsirtHandler := handlers.NewSdmCsirtHandler(sdmCsirtService)
	roleHandler := handlers.NewRoleHandler(roleService, sseService)
	casbinHandler := handlers.NewCasbinHandler(casbinService, sseService)
	sseHandler := handlers.NewSSEHandler(sseService)
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
		sektorHandler,
		subSektorHandler,
		seHandler,
		dashboardHandler,
	)

	// Start server
	log.Printf("Server starting on %s", cfg.Port)
	rollbar.Info("Server starting on %s", cfg.Port)
	log.Println("Rate limiting enabled:")
	log.Println("  - Auth endpoints: 5 requests/minute per IP")
	log.Println("  - Public posts: 1000 requests/minute per IP")
	log.Println("  - Protected posts: 100 requests/minute per user")
	log.Println("SSE endpoint available at /api/events")

	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}