package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	Charset         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func NewMySQLConnection(cfg Config) (*sql.DB, error) {
	// Validasi Config
	if cfg.Host == "" {
		return nil, errors.New("host is required")
	}
	if cfg.Port == "" {
		return nil, errors.New("port is required")
	}
	if cfg.User == "" {
		return nil, errors.New("user is required")
	}
	if cfg.DBName == "" {
		return nil, errors.New("database name is required")
	}

	// Dsn
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Tes connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Successfully connected to MySQL database")

	return db, nil
}
