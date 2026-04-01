package database_test

import (
	"fortyfour-backend/pkg/database"
	"testing"
	"time"

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

func TestNewMySQLConnection_Validation_AllFieldsMissing(t *testing.T) {
	// Config kosong total → harus error di validasi pertama (host)
	_, err := database.NewMySQLConnection(database.Config{})
	assert.Error(t, err)
	assert.Equal(t, "host is required", err.Error())
}

func TestNewMySQLConnection_Validation_PasswordOptional(t *testing.T) {
	// Password kosong bukan merupakan error validasi — koneksi tetap dicoba
	// (akan gagal di Ping karena tidak ada DB nyata, bukan di validasi)
	cfg := database.Config{
		Host:   "localhost",
		Port:   "3306",
		User:   "root",
		DBName: "testdb",
		// Password sengaja dikosongkan
	}
	_, err := database.NewMySQLConnection(cfg)
	// Error harus terjadi di Ping (koneksi gagal), bukan di validasi field
	assert.Error(t, err)
	assert.NotEqual(t, "host is required", err.Error())
	assert.NotEqual(t, "port is required", err.Error())
	assert.NotEqual(t, "user is required", err.Error())
	assert.NotEqual(t, "database name is required", err.Error())
}

func TestNewMySQLConnection_Validation_ErrorOrdering(t *testing.T) {
	// Validasi harus mengikuti urutan: host → port → user → dbname
	// Jika host dan port keduanya kosong, error host yang muncul duluan
	_, err := database.NewMySQLConnection(database.Config{
		User:   "root",
		DBName: "test",
		// Host dan Port kosong
	})
	assert.Error(t, err)
	assert.Equal(t, "host is required", err.Error(),
		"host seharusnya divalidasi sebelum port")
}

func TestNewMySQLConnection_Ping_SuccessAndFailure(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	t.Run("ping sukses tidak return error", func(t *testing.T) {
		mock.ExpectPing()
		assert.NoError(t, db.Ping())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ping gagal return error", func(t *testing.T) {
		mock.ExpectPing().WillReturnError(assert.AnError)
		assert.Error(t, db.Ping())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("beberapa ping berurutan", func(t *testing.T) {
		mock.ExpectPing()
		mock.ExpectPing()
		assert.NoError(t, db.Ping())
		assert.NoError(t, db.Ping())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestConfig_DSNFieldsReflectConfig(t *testing.T) {
	// Memverifikasi bahwa struct Config bisa diisi dengan field lengkap tanpa error
	cfg := database.Config{
		Host:            "db.example.com",
		Port:            "3306",
		User:            "appuser",
		Password:        "s3cr3t",
		DBName:          "appdb",
		Charset:         "utf8mb4",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * 60 * 1000000000, // 5 menit dalam nanoseconds
	}

	assert.Equal(t, "db.example.com", cfg.Host)
	assert.Equal(t, "3306", cfg.Port)
	assert.Equal(t, "appuser", cfg.User)
	assert.Equal(t, "s3cr3t", cfg.Password)
	assert.Equal(t, "appdb", cfg.DBName)
	assert.Equal(t, "utf8mb4", cfg.Charset)
	assert.Equal(t, 10, cfg.MaxOpenConns)
	assert.Equal(t, 5, cfg.MaxIdleConns)
	assert.Equal(t, 5*time.Minute, cfg.ConnMaxLifetime)
}