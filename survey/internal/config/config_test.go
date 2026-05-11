package config

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// Load Config
func TestLoad_DefaultValues(t *testing.T) {
	os.Clearenv()

	cfg := Load()

	if cfg.Port != "8082" {
		t.Errorf("expected default port 8082, got %s", cfg.Port)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("expected default db host localhost")
	}

	if cfg.Redis.DB != 0 {
		t.Errorf("expected default redis db 0")
	}
}

func TestLoad_WithEnv(t *testing.T) {
	os.Setenv("PORT", "9000")
	os.Setenv("DB_HOST", "dbhost")
	os.Setenv("REDIS_DB", "2")

	defer os.Clearenv()

	cfg := Load()

	if cfg.Port != "9000" {
		t.Errorf("expected 9000, got %s", cfg.Port)
	}

	if cfg.Database.Host != "dbhost" {
		t.Errorf("expected dbhost")
	}

	if cfg.Redis.DB != 2 {
		t.Errorf("expected redis db 2")
	}
}

// getEnv
func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_KEY", "value")
	defer os.Unsetenv("TEST_KEY")

	val := getEnv("TEST_KEY", "fallback")
	if val != "value" {
		t.Errorf("expected value, got %s", val)
	}

	val = getEnv("UNKNOWN", "fallback")
	if val != "fallback" {
		t.Errorf("expected fallback")
	}
}

// getEnvAsInt
func TestGetEnvAsInt(t *testing.T) {
	os.Setenv("INT_KEY", "5")
	defer os.Unsetenv("INT_KEY")

	val := getEnvAsInt("INT_KEY", 1)
	if val != 5 {
		t.Errorf("expected 5")
	}

	os.Setenv("INT_KEY", "invalid")
	val = getEnvAsInt("INT_KEY", 1)
	if val != 1 {
		t.Errorf("expected fallback 1")
	}
}

// GetDSN
func TestGetDSN(t *testing.T) {
	db := DatabaseConfig{
		User:     "root",
		Password: "pass",
		Host:     "localhost",
		Port:     "3306",
		DBName:   "testdb",
	}

	dsn := db.GetDSN()

	expected := "root:pass@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"

	if dsn != expected {
		t.Errorf("expected %s, got %s", expected, dsn)
	}
}

// RabbitMQ URL
func TestRabbitMQ_GetURL(t *testing.T) {
	r := RabbitMQConfig{
		User:     "guest",
		Password: "guest",
		Host:     "localhost",
		Port:     "5672",
		Vhost:    "test",
	}

	url := r.GetURL()
	expected := "amqp://guest:guest@localhost:5672/test"

	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}
}

func TestRabbitMQ_GetURL_DefaultVhost(t *testing.T) {
	r := RabbitMQConfig{
		User:     "guest",
		Password: "guest",
		Host:     "localhost",
		Port:     "5672",
		Vhost:    "",
	}

	url := r.GetURL()
	expected := "amqp://guest:guest@localhost:5672/"

	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}
}

// InitDB Success
func TestInitDB_Success(t *testing.T) {
	dbMock, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("error mock: %v", err)
	}
	defer dbMock.Close()

	// override sqlOpen
	originalOpen := sqlOpen
	defer func() { sqlOpen = originalOpen }()

	sqlOpen = func(driverName, dsn string) (*sql.DB, error) {
		return dbMock, nil
	}

	mock.ExpectPing()

	cfg := &Config{
		Database: DatabaseConfig{
			User:   "root",
			Host:   "localhost",
			Port:   "3306",
			DBName: "test",
		},
	}

	db, err := InitDB(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if db == nil {
		t.Errorf("expected db")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// InitDB Ping Error
func TestInitDB_PingError(t *testing.T) {
	dbMock, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("error mock: %v", err)
	}
	defer dbMock.Close()

	originalOpen := sqlOpen
	defer func() { sqlOpen = originalOpen }()

	sqlOpen = func(driverName, dsn string) (*sql.DB, error) {
		return dbMock, nil
	}

	mock.ExpectPing().WillReturnError(errors.New("ping error"))

	cfg := &Config{
		Database: DatabaseConfig{
			User:   "root",
			Host:   "localhost",
			Port:   "3306",
			DBName: "test",
		},
	}

	db, err := InitDB(cfg)
	if err == nil {
		t.Errorf("expected error")
	}
	if db != nil {
		t.Errorf("expected nil db")
	}
}
