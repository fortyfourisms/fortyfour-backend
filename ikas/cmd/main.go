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

	"fortyfour-backend/pkg/logger"
	pkgRmq "fortyfour-backend/pkg/rabbitmq"

	"github.com/joho/godotenv"
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
	// load env
	if err := godotenv.Load(); err != nil {
		os.Stdout.WriteString("Warning: .env file not found\n")
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
		repository.NewIdentifikasiRepository(db),
		repository.NewPertanyaanIdentifikasiRepository(db),
		repository.NewJawabanIdentifikasiRepository(db),
		repository.NewProteksiRepository(db),
		repository.NewPertanyaanProteksiRepository(db),
		repository.NewJawabanProteksiRepository(db),
		repository.NewDeteksiRepository(db),
		repository.NewPertanyaanDeteksiRepository(db),
		repository.NewJawabanDeteksiRepository(db),
		repository.NewGulihRepository(db),
		repository.NewPertanyaanGulihRepository(db),
		repository.NewJawabanGulihRepository(db),
		repository.NewDomainRepository(db),
		repository.NewRuangLingkupRepository(db),
		repository.NewKategoriRepository(db),
		repository.NewSubKategoriRepository(db),
		repository.NewAuditLogRepository(db),
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

	identifikasiRepo := repository.NewIdentifikasiRepository(db)
	proteksiRepo := repository.NewProteksiRepository(db)
	deteksiRepo := repository.NewDeteksiRepository(db)
	gulihRepo := repository.NewGulihRepository(db)

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
	ruangLingkupService := services.NewRuangLingkupService(ruangLingkupRepo, msgProducer)
	domainService := services.NewDomainService(domainRepo, msgProducer)
	kategoriService := services.NewKategoriService(kategoriRepo, msgProducer)
	subKategoriService := services.NewSubKategoriService(subKategoriRepo, msgProducer)
	identifikasiService := services.NewIdentifikasiService(identifikasiRepo)
	proteksiService := services.NewProteksiService(proteksiRepo)
	deteksiService := services.NewDeteksiService(deteksiRepo)
	gulihService := services.NewGulihService(gulihRepo)
	pertanyaanIdentifikasiService := services.NewPertanyaanIdentifikasiService(pertanyaanIdentifikasiRepo, msgProducer)
	pertanyaanProteksiService := services.NewPertanyaanProteksiService(pertanyaanProteksiRepo, msgProducer)
	pertanyaanDeteksiService := services.NewPertanyaanDeteksiService(pertanyaanDeteksiRepo, msgProducer)
	pertanyaanGulihService := services.NewPertanyaanGulihService(pertanyaanGulihRepo, msgProducer)
	jawabanIdentifikasiService := services.NewJawabanIdentifikasiService(jawabanIdentifikasiRepo, ikasRepo, msgProducer)
	jawabanProteksiService := services.NewJawabanProteksiService(jawabanProteksiRepo, ikasRepo, msgProducer)
	jawabanDeteksiService := services.NewJawabanDeteksiService(jawabanDeteksiRepo, ikasRepo, msgProducer)
	jawabanGulihService := services.NewJawabanGulihService(jawabanGulihRepo, ikasRepo, msgProducer)

	// handlers
	ikasHandler := handlers.NewIkasHandler(ikasService)
	ruangLingkupHandler := handlers.NewRuangLingkupHandler(ruangLingkupService)
	domainHandler := handlers.NewDomainHandler(domainService)
	kategoriHandler := handlers.NewKategoriHandler(kategoriService)
	subKategoriHandler := handlers.NewSubKategoriHandler(subKategoriService)
	identifikasiHandler := handlers.NewIdentifikasiHandler(identifikasiService)
	proteksiHandler := handlers.NewProteksiHandler(proteksiService)
	deteksiHandler := handlers.NewDeteksiHandler(deteksiService)
	gulihHandler := handlers.NewGulihHandler(gulihService)
	pertanyaanIdentifikasiHandler := handlers.NewPertanyaanIdentifikasiHandler(pertanyaanIdentifikasiService)
	pertanyaanProteksiHandler := handlers.NewPertanyaanProteksiHandler(pertanyaanProteksiService)
	pertanyaanDeteksiHandler := handlers.NewPertanyaanDeteksiHandler(pertanyaanDeteksiService)
	pertanyaanGulihHandler := handlers.NewPertanyaanGulihHandler(pertanyaanGulihService)
	jawabanIdentifikasiHandler := handlers.NewJawabanIdentifikasiHandler(jawabanIdentifikasiService)
	jawabanProteksiHandler := handlers.NewJawabanProteksiHandler(jawabanProteksiService)
	jawabanDeteksiHandler := handlers.NewJawabanDeteksiHandler(jawabanDeteksiService)
	jawabanGulihHandler := handlers.NewJawabanGulihHandler(jawabanGulihService)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.InternalGatewayKey)

	// Router
	mux := routes.InitRouter(
		ikasHandler,
		ruangLingkupHandler,
		domainHandler,
		kategoriHandler,
		subKategoriHandler,
		identifikasiHandler,
		proteksiHandler,
		deteksiHandler,
		gulihHandler,
		pertanyaanIdentifikasiHandler,
		pertanyaanProteksiHandler,
		pertanyaanDeteksiHandler,
		pertanyaanGulihHandler,
		jawabanIdentifikasiHandler,
		jawabanProteksiHandler,
		jawabanDeteksiHandler,
		jawabanGulihHandler,
		authMiddleware,
	)

	go func() {
		logger.Infof("IKAS service running on %s", cfg.Port)
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
		logger.Info("  - jawaban.deteksi.created")
		logger.Info("  - jawaban.deteksi.updated")
		logger.Info("  - jawaban.deteksi.deleted")
		logger.Info("  - jawaban.gulih.created")
		logger.Info("  - jawaban.gulih.updated")
		logger.Info("  - jawaban.gulih.deleted")
		logger.Info("  - domain.created")
		logger.Info("  - domain.updated")
		logger.Info("  - domain.deleted")
		logger.Info("  - ruang_lingkup.created")
		logger.Info("  - ruang_lingkup.updated")
		logger.Info("  - ruang_lingkup.deleted")
		logger.Info("  - kategori.deleted")
		logger.Info("  - sub_kategori.created")
		logger.Info("  - sub_kategori.updated")
		logger.Info("  - sub_kategori.deleted")
		logger.Info("  - pertanyaan_identifikasi.created")
		logger.Info("  - pertanyaan_identifikasi.updated")
		logger.Info("  - pertanyaan_identifikasi.deleted")
		logger.Info("  - pertanyaan_proteksi.created")
		logger.Info("  - pertanyaan_proteksi.updated")
		logger.Info("  - pertanyaan_proteksi.deleted")
		logger.Info("  - pertanyaan_deteksi.created")
		logger.Info("  - pertanyaan_deteksi.updated")
		logger.Info("  - pertanyaan_deteksi.deleted")
		logger.Info("  - pertanyaan_gulih.created")
		logger.Info("  - pertanyaan_gulih.updated")
		logger.Info("  - pertanyaan_gulih.deleted")
		logger.Info("  - ikas.audit_logs")

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
