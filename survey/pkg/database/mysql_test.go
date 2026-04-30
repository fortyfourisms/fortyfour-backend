package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// CONFIG VALIDATION TEST
func TestNewMySQLConnection_InvalidConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
	}{
		{"missing host", Config{Port: "3306", User: "root", DBName: "test"}},
		{"missing port", Config{Host: "localhost", User: "root", DBName: "test"}},
		{"missing user", Config{Host: "localhost", Port: "3306", DBName: "test"}},
		{"missing db name", Config{Host: "localhost", Port: "3306", User: "root"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewMySQLConnection(tt.cfg)
			if err == nil {
				t.Errorf("expected error, got nil")
			}
			if db != nil {
				t.Errorf("expected db to be nil")
			}
		})
	}
}

// SUCCESS CONNECTION TEST
func TestNewMySQLConnection_Success(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock: %v", err)
	}
	defer dbMock.Close()

	// override sqlOpen
	originalOpen := sqlOpen
	defer func() { sqlOpen = originalOpen }()

	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return dbMock, nil
	}

	// Expect ping
	mock.ExpectPing()

	cfg := Config{
		Host:   "localhost",
		Port:   "3306",
		User:   "root",
		DBName: "test",
	}

	db, err := NewMySQLConnection(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if db == nil {
		t.Errorf("expected db, got nil")
	}

	// verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// PING FAILED TEST
func TestNewMySQLConnection_PingError(t *testing.T) {
	dbMock, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("error creating mock: %v", err)
	}
	defer dbMock.Close()

	// override sqlOpen
	originalOpen := sqlOpen
	defer func() { sqlOpen = originalOpen }()

	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return dbMock, nil
	}

	mock.ExpectPing().WillReturnError(errors.New("ping failed"))

	cfg := Config{
		Host:   "localhost",
		Port:   "3306",
		User:   "root",
		DBName: "test",
	}

	db, err := NewMySQLConnection(cfg)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if db != nil {
		t.Errorf("expected nil db")
	}
}
