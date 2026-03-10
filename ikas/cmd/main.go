package main

import (
	"context"
	"ikas/internal/config"
	"ikas/internal/handlers"
	"ikas/internal/middleware"
	internalRmq "ikas/internal/rabbitmq"
	"ikas/internal/repository"
	"ikas/internal/routes"
	"ikas/internal/services"
	"ikas/pkg/cache"
	"ikas/pkg/database"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pkgRmq "fortyfour-backend/pkg/rabbitmq"

	"fortyfour-backend/pkg/logger"
)

// @title IKAS API
// @version 1.0
// @description API documentation for IKAS (Indeks Kematangan Keamanan Siber) service.
// @host localhost:8081
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

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
	logger.Info("Redis initialized successfully")
	defer redisClient.Close()

	// Initialize RabbitMQ (Shared Package)
	rmq, err := pkgRmq.NewRabbitMQ(cfg.RabbitMQ.GetURL())
	if err != nil {
		logger.FatalErr(err, "Failed to connect to RabbitMQ")
	}
	defer rmq.Close()

	// Setup RabbitMQ infrastructure (Specific to IKAS)
	if err := internalRmq.SetupInfrastructure(rmq); err != nil {
		logger.FatalErr(err, "Failed to setup RabbitMQ infrastructure")
	}
	logger.Info("RabbitMQ initialized successfully")

	// Create Shared Producer and Consumer
	sharedProducer := pkgRmq.NewProducer(rmq.GetChannel())
	sharedConsumer := pkgRmq.NewConsumer(rmq.GetChannel())

	// Wrap with IKAS specific Producer and Consumer
	msgProducer := internalRmq.NewProducer(sharedProducer)
	msgConsumer := internalRmq.NewConsumer(
		sharedConsumer,
		repository.NewIkasRepository(db),
		repository.NewJawabanIdentifikasiRepository(db),
		repository.NewPertanyaanIdentifikasiRepository(db),
		repository.NewJawabanProteksiRepository(db),
		repository.NewPertanyaanProteksiRepository(db),
	)

	// Start consumers in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := msgConsumer.StartAllConsumers(ctx); err != nil {
		logger.FatalErr(err, "Failed to start consumers")
	}

	// repository
	ikasRepo := repository.NewIkasRepository(db)
	ruangLingkupRepo := repository.NewRuangLingkupRepository(db)
	domainRepo := repository.NewDomainRepository(db)
	kategoriRepo := repository.NewKategoriRepository(db)
	subKategoriRepo := repository.NewSubKategoriRepository(db)

	pertanyaanIdentifikasiRepo := repository.NewPertanyaanIdentifikasiRepository(db)
	pertanyaanProteksiRepo := repository.NewPertanyaanProteksiRepository(db)
	pertanyaanDeteksiRepo := repository.NewPertanyaanDeteksiRepository(db)
	pertanyaanGulihRepo := repository.NewPertanyaanGulihRepository(db)

	jawabanIdentifikasiRepo := repository.NewJawabanIdentifikasiRepository(db)
	jawabanProteksiRepo := repository.NewJawabanProteksiRepository(db)
	jawabanDeteksiRepo := repository.NewJawabanDeteksiRepository(db)
	jawabanGulihRepo := repository.NewJawabanGulihRepository(db)

	// services
	ikasService := services.NewIkasService(ikasRepo, msgProducer)
	ruangLingkupService := services.NewRuangLingkupService(ruangLingkupRepo)
	domainService := services.NewDomainService(domainRepo)
	kategoriService := services.NewKategoriService(kategoriRepo)
	subKategoriService := services.NewSubKategoriService(subKategoriRepo)
	pertanyaanIdentifikasiService := services.NewPertanyaanIdentifikasiService(pertanyaanIdentifikasiRepo)
	pertanyaanProteksiService := services.NewPertanyaanProteksiService(pertanyaanProteksiRepo)
	pertanyaanDeteksiService := services.NewPertanyaanDeteksiService(pertanyaanDeteksiRepo)
	pertanyaanGulihService := services.NewPertanyaanGulihService(pertanyaanGulihRepo)
	jawabanIdentifikasiService := services.NewJawabanIdentifikasiService(jawabanIdentifikasiRepo, msgProducer)
	jawabanProteksiService := services.NewJawabanProteksiService(jawabanProteksiRepo, msgProducer)
	jawabanDeteksiService := services.NewJawabanDeteksiService(jawabanDeteksiRepo)
	jawabanGulihService := services.NewJawabanGulihService(jawabanGulihRepo)

	// handlers
	ikasHandler := handlers.NewIkasHandler(ikasService)
	ruangLingkupHandler := handlers.NewRuangLingkupHandler(ruangLingkupService)
	domainHandler := handlers.NewDomainHandler(domainService)
	kategoriHandler := handlers.NewKategoriHandler(kategoriService)
	subKategoriHandler := handlers.NewSubKategoriHandler(subKategoriService)
	pertanyaanIdentifikasiHandler := handlers.NewPertanyaanIdentifikasiHandler(pertanyaanIdentifikasiService)
	pertanyaanProteksiHandler := handlers.NewPertanyaanProteksiHandler(pertanyaanProteksiService)
	pertanyaanDeteksiHandler := handlers.NewPertanyaanDeteksiHandler(pertanyaanDeteksiService)
	pertanyaanGulihHandler := handlers.NewPertanyaanGulihHandler(pertanyaanGulihService)
	jawabanIdentifikasiHandler := handlers.NewJawabanIdentifikasiHandler(jawabanIdentifikasiService)
	jawabanProteksiHandler := handlers.NewJawabanProteksiHandler(jawabanProteksiService)
	jawabanDeteksiHandler := handlers.NewJawabanDeteksiHandler(jawabanDeteksiService)
	jawabanGulihHandler := handlers.NewJawabanGulihHandler(jawabanGulihService)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// Initialize rate limiters with different configurations
	rateLimitConfigs := middleware.GetRateLimitConfigs()

	strictLimiter := middleware.NewRateLimiter(redisClient, rateLimitConfigs.Strict)
	moderateLimiter := middleware.NewRateLimiter(redisClient, rateLimitConfigs.Moderate)
	lenientLimiter := middleware.NewRateLimiter(redisClient, rateLimitConfigs.Lenient)

	// Router
	mux := routes.InitRouter(
		ikasHandler,
		ruangLingkupHandler,
		domainHandler,
		kategoriHandler,
		subKategoriHandler,
		pertanyaanIdentifikasiHandler,
		pertanyaanProteksiHandler,
		pertanyaanDeteksiHandler,
		pertanyaanGulihHandler,
		jawabanIdentifikasiHandler,
		jawabanProteksiHandler,
		jawabanDeteksiHandler,
		jawabanGulihHandler,
		authMiddleware,
		strictLimiter,
		moderateLimiter,
		lenientLimiter,
	)

	go func() {
		logger.Infof("IKAS service running on %s", cfg.Port)
		logger.Info("Rate limiting enabled:")
		logger.Info("  - Auth endpoints: 5 requests/minute per IP")
		logger.Info("  - Public posts: 60 requests/minute per IP")
		logger.Info("  - Protected posts: 20 requests/minute per user")
		logger.Info("RabbitMQ consumers running:")
		logger.Info("  - ikas.created")
		logger.Info("  - ikas.updated")
		logger.Info("  - ikas.deleted")
		logger.Info("  - ikas.imported")
		logger.Info("  - notifications.email")
		logger.Info("  - jawaban.identifikasi.created")
		logger.Info("  - jawaban.identifikasi.updated")
		logger.Info("  - jawaban.identifikasi.deleted")
		logger.Info("  - jawaban.proteksi.created")
		logger.Info("  - jawaban.proteksi.updated")
		logger.Info("  - jawaban.proteksi.deleted")

		if err := http.ListenAndServe(cfg.Port, mux); err != nil && err != http.ErrServerClosed {
			logger.FatalErr(err, "Server failed")
		}
	}()

	// Wait for interrupt signal untuk graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	cancel() // Stop consumers
	logger.Info("Server stopped gracefully")

}
