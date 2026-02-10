package config

import (
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	Port            string
	JWTSecret       string
	Domain          string
	Database        DatabaseConfig
	Redis           RedisConfig
	CasbinModelPath string
	Rollbar         RollbarConfig
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

type RollbarConfig struct {
	Token string
	Env   string
}

func Load() *Config {
	absPath, _ := filepath.Abs("casbin/casbin_model.conf")
	return &Config{
		Port:      getEnv("PORT", ":8080"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
		Domain:    getEnv("DOMAIN", "https://admin.kssindustri.site"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "fortyfour-backend_db"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		CasbinModelPath: getEnv("CASBIN_MODEL_PATH", absPath),
		Rollbar: RollbarConfig{
			Token: getEnv("ROLLBAR_TOKEN", "0eddf8fb05e44067a12a8bb36ccc3ef9"),
			Env:   getEnv("ROLLBAR_STATUS", "production"),
		},
	}
}

// GetDSN returns MySQL DSN for GORM
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
