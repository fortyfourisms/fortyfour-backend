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
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pkgRmq "fortyfour-backend/pkg/rabbitmq"

	"github.com/rollbar/rollbar-go"
)

func main() {

	// Load env
	// if err := godotenv.Load(); err != nil {
	// 	log.Println("No .env file found")
	// }

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
		log.Fatal(err)
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
		log.Println("Failed to connect to Redis:", err)
		rollbar.Error(err)
		log.Fatal(err)
	}
	log.Println("Redis initialized successfully")
	rollbar.Info("Redis initialized successfully")
	defer redisClient.Close()

	// Initialize RabbitMQ (Shared Package)
	rmq, err := pkgRmq.NewRabbitMQ(cfg.RabbitMQ.GetURL())
	if err != nil {
		log.Println("Failed to connect to RabbitMQ:", err)
		rollbar.Error(err)
		log.Fatal(err)
	}
	defer rmq.Close()

	// Setup RabbitMQ infrastructure (Specific to IKAS)
	if err := internalRmq.SetupInfrastructure(rmq); err != nil {
		log.Println("Failed to setup RabbitMQ infrastructure:", err)
		rollbar.Error(err)
		log.Fatal(err)
	}
	log.Println("RabbitMQ initialized successfully")
	rollbar.Info("RabbitMQ initialized successfully")

	// Create Shared Producer and Consumer
	sharedProducer := pkgRmq.NewProducer(rmq.GetChannel())
	sharedConsumer := pkgRmq.NewConsumer(rmq.GetChannel())

	// Wrap with IKAS specific Producer and Consumer
	msgProducer := internalRmq.NewProducer(sharedProducer)
	msgConsumer := internalRmq.NewConsumer(sharedConsumer)

	// Start consumers in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := msgConsumer.StartAllConsumers(ctx); err != nil {
		log.Println("Failed to start consumers:", err)
		rollbar.Error(err)
		log.Fatal(err)
	}

	// repository
	ikasRepo := repository.NewIkasRepository(db)
	ruangLingkupRepo := repository.NewRuangLingkupRepository(db)
	domainRepo := repository.NewDomainRepository(db)
	kategoriRepo := repository.NewKategoriRepository(db)
	subKategoriRepo := repository.NewSubKategoriRepository(db)
	pertanyaanIdentifikasiRepo := repository.NewPertanyaanIdentifikasiRepository(db)
	jawabanIdentifikasiRepo := repository.NewJawabanIdentifikasiRepository(db)

	// services
	ikasService := services.NewIkasService(ikasRepo, msgProducer)
	ruangLingkupService := services.NewRuangLingkupService(ruangLingkupRepo)
	domainService := services.NewDomainService(domainRepo)
	kategoriService := services.NewKategoriService(kategoriRepo)
	subKategoriService := services.NewSubKategoriService(subKategoriRepo)
	pertanyaanIdentifikasiService := services.NewPertanyaanIdentifikasiService(pertanyaanIdentifikasiRepo)
	jawabanIdentifikasiService := services.NewJawabanIdentifikasiService(jawabanIdentifikasiRepo)

	// handlers
	ikasHandler := handlers.NewIkasHandler(ikasService)
	ruangLingkupHandler := handlers.NewRuangLingkupHandler(ruangLingkupService)
	domainHandler := handlers.NewDomainHandler(domainService)
	kategoriHandler := handlers.NewKategoriHandler(kategoriService)
	subKategoriHandler := handlers.NewSubKategoriHandler(subKategoriService)
	pertanyaanIdentifikasiHandler := handlers.NewPertanyaanIdentifikasiHandler(pertanyaanIdentifikasiService)
	jawabanIdentifikasiHandler := handlers.NewJawabanIdentifikasiHandler(jawabanIdentifikasiService)

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
		jawabanIdentifikasiHandler,
		authMiddleware,
		strictLimiter,
		moderateLimiter,
		lenientLimiter,
	)

	go func() {
		log.Println("IKAS service running on", cfg.Port)
		log.Println("Rate limiting enabled:")
		log.Println("  - Auth endpoints: 5 requests/minute per IP")
		log.Println("  - Public posts: 60 requests/minute per IP")
		log.Println("  - Protected posts: 20 requests/minute per user")
		log.Println("RabbitMQ consumers running:")
		log.Println("  - ikas.created")
		log.Println("  - ikas.updated")
		log.Println("  - ikas.deleted")
		log.Println("  - ikas.imported")
		log.Println("  - notifications.email")

		if err := http.ListenAndServe(cfg.Port, mux); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal untuk graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	cancel() // Stop consumers
	log.Println("Server stopped gracefully")

}
