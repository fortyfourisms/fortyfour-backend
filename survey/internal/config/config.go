package config

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Port            string
	JWTSecret       string
	Database        DatabaseConfig
	Redis           RedisConfig
	RabbitMQ        RabbitMQConfig
	CasbinModelPath string
	LogLevel        string
	Environment     string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Vhost    string
}

func Load() *Config {
	absPath, _ := filepath.Abs("casbin/casbin_model.conf")

	return &Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "secret"),

		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "fortyfour-backend_db"),
		},

		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},

		RabbitMQ: RabbitMQConfig{
			Host:     getEnv("RABBITMQ_HOST", "localhost"),
			Port:     getEnv("RABBITMQ_PORT", "5672"),
			User:     getEnv("RABBITMQ_USER", "guest"),
			Password: getEnv("RABBITMQ_PASSWORD", "guest"),
			Vhost:    getEnv("RABBITMQ_VHOST", "/"),
		},

		CasbinModelPath: getEnv("CASBIN_MODEL_PATH", absPath),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		Environment:     getEnv("ENVIRONMENT", "development"),
	}
}

func InitDB(cfg *Config) (*sql.DB, error) {

	dsn := cfg.Database.GetDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open error: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping error: %w", err)
	}

	return db, nil
}

func (r *RabbitMQConfig) GetURL() string {
	vhost := r.Vhost
	if vhost == "" {
		vhost = "/"
	}

	if vhost[0] != '/' {
		vhost = "/" + vhost
	}

	return "amqp://" + r.User + ":" + r.Password + "@" + r.Host + ":" + r.Port + vhost
}

// MySQL DSN
func (d *DatabaseConfig) GetDSN() string {
	return d.User + ":" + d.Password + "@tcp(" + d.Host + ":" + d.Port + ")/" + d.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}