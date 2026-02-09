package database_test

import (
	"fortyfour-backend/pkg/database"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

const (
	DriverName = "mysql"
)

func TestNewMySQLConnection_Validation(t *testing.T) {
	tests := []struct {
		name    string
		config  database.Config
		wantErr string
	}{
		{
			name: "Missing Host",
			config: database.Config{
				Port:   "3306",
				User:   "root",
				DBName: "test",
			},
			wantErr: "host is required",
		},
		{
			name: "Missing Port",
			config: database.Config{
				Host:   "localhost",
				User:   "root",
				DBName: "test",
			},
			wantErr: "port is required",
		},
		{
			name: "Missing User",
			config: database.Config{
				Host:   "localhost",
				Port:   "3306",
				DBName: "test",
			},
			wantErr: "user is required",
		},
		{
			name: "Missing DBName",
			config: database.Config{
				Host: "localhost",
				Port: "3306",
				User: "root",
			},
			wantErr: "database name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := database.NewMySQLConnection(tt.config)
			assert.Error(t, err)
			assert.Equal(t, tt.wantErr, err.Error())
		})
	}
}

func TestNewMySQLConnection_MockOnly(t *testing.T) {
	// ... existing code ...
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Successful ping", func(t *testing.T) {
		mock.ExpectPing()
		err = db.Ping()
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed ping", func(t *testing.T) {
		mock.ExpectPing().WillReturnError(assert.AnError)
		err = db.Ping()
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
