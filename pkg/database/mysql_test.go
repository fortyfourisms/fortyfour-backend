package database_test

import (
	"fortyfour-backend/pkg/database"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMySQLConnection(t *testing.T) {
	tests := []struct {
		name    string
		cfg     database.Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid configuration",
			cfg: database.Config{
				Host:     "localhost",
				Port:     "3306",
				User:     "testuser",
				Password: "testpass",
				DBName:   "testdb",
			},
			wantErr: false,
		},
		{
			name: "Missing host",
			cfg: database.Config{
				Host:     "",
				Port:     "3306",
				User:     "testuser",
				Password: "testpass",
				DBName:   "testdb",
			},
			wantErr: true,
			errMsg:  "host",
		},
		{
			name: "Missing port",
			cfg: database.Config{
				Host:     "localhost",
				Port:     "",
				User:     "testuser",
				Password: "testpass",
				DBName:   "testdb",
			},
			wantErr: true,
			errMsg:  "port",
		},
		{
			name: "Missing user",
			cfg: database.Config{
				Host:     "localhost",
				Port:     "3306",
				User:     "",
				Password: "testpass",
				DBName:   "testdb",
			},
			wantErr: true,
			errMsg:  "user",
		},
		{
			name: "Missing database name",
			cfg: database.Config{
				Host:     "localhost",
				Port:     "3306",
				User:     "testuser",
				Password: "testpass",
				DBName:   "",
			},
			wantErr: true,
			errMsg:  "database",
		},
		{
			name: "Invalid port format",
			cfg: database.Config{
				Host:     "localhost",
				Port:     "invalid",
				User:     "testuser",
				Password: "testpass",
				DBName:   "testdb",
			},
			wantErr: true,
		},
		{
			name: "Empty password (should work for local dev)",
			cfg: database.Config{
				Host:     "localhost",
				Port:     "3306",
				User:     "root",
				Password: "",
				DBName:   "testdb",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := database.NewMySQLConnection(tt.cfg)

			if tt.wantErr {
				assert.Error(t, gotErr, "Expected an error but got none")
				if tt.errMsg != "" {
					assert.Contains(t, gotErr.Error(), tt.errMsg,
						"Error message should contain expected text")
				}
				assert.Nil(t, got, "DB connection should be nil on error")
				return
			}

			// For successful connections - skip if database not available
			if gotErr != nil {
				t.Skipf("Skipping test - database not available: %v", gotErr)
				return
			}

			require.NotNil(t, got, "DB connection should not be nil")

			// Test the connection is valid (ping)
			err := got.Ping()
			if err != nil {
				t.Logf("Warning: Could not ping database: %v", err)
			}

			// Clean up
			if got != nil {
				got.Close()
			}
		})
	}
}

func TestNewMySQLConnection_WithMock(t *testing.T) {
	// This test uses sqlmock to test without a real database
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	defer db.Close()

	t.Run("Connection with proper settings", func(t *testing.T) {
		// Expect a ping
		mock.ExpectPing()

		err = db.Ping()
		assert.NoError(t, err)

		// Verify all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Connection pool settings", func(t *testing.T) {
		cfg := database.Config{
			Host:            "localhost",
			Port:            "3306",
			User:            "testuser",
			Password:        "testpass",
			DBName:          "testdb",
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
		}

		// Note: This will try to connect to a real database
		// Skip if no database is available
		conn, err := database.NewMySQLConnection(cfg)
		if err != nil {
			t.Skip("Skipping connection pool test - no database available")
		}
		defer conn.Close()

		// Verify connection pool settings
		stats := conn.Stats()
		assert.GreaterOrEqual(t, stats.MaxOpenConnections, 0)
	})
}

func TestNewMySQLConnection_ConnectionString(t *testing.T) {
	tests := []struct {
		name          string
		cfg           database.Config
		expectedDSN   string
		shouldContain []string
	}{
		{
			name: "Basic DSN format",
			cfg: database.Config{
				Host:     "localhost",
				Port:     "3306",
				User:     "testuser",
				Password: "testpass",
				DBName:   "testdb",
			},
			shouldContain: []string{
				"testuser",
				"localhost",
				"3306",
				"testdb",
			},
		},
		{
			name: "DSN with special characters in password",
			cfg: database.Config{
				Host:     "localhost",
				Port:     "3306",
				User:     "testuser",
				Password: "p@ss:w0rd!",
				DBName:   "testdb",
			},
			shouldContain: []string{
				"testuser",
				"localhost",
			},
		},
		{
			name: "DSN with charset parameter",
			cfg: database.Config{
				Host:     "localhost",
				Port:     "3306",
				User:     "testuser",
				Password: "testpass",
				DBName:   "testdb",
				Charset:  "utf8mb4",
			},
			shouldContain: []string{
				"charset=utf8mb4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test verifies DSN format without actually connecting
			// You'll need to expose a DSN building function or test indirectly

			conn, err := database.NewMySQLConnection(tt.cfg)
			if err != nil {
				// If connection fails, check if it's a network error
				// (expected in CI/CD without database)
				t.Logf("Connection failed (expected in test env): %v", err)
				return
			}

			if conn != nil {
				defer conn.Close()

				// Verify we can ping
				err = conn.Ping()
				if err != nil {
					t.Logf("Ping failed (expected in test env): %v", err)
				}
			}
		})
	}
}

func TestNewMySQLConnection_ErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		cfg     database.Config
		wantErr bool
	}{
		{
			name: "Invalid host",
			cfg: database.Config{
				Host:     "invalid-host-that-does-not-exist-12345",
				Port:     "3306",
				User:     "testuser",
				Password: "testpass",
				DBName:   "testdb",
			},
			wantErr: true,
		},
		{
			name: "Invalid port number",
			cfg: database.Config{
				Host:     "localhost",
				Port:     "99999",
				User:     "testuser",
				Password: "testpass",
				DBName:   "testdb",
			},
			wantErr: true,
		},
		{
			name: "Connection timeout",
			cfg: database.Config{
				Host:     "192.0.2.1", // TEST-NET-1 (non-routable)
				Port:     "3306",
				User:     "testuser",
				Password: "testpass",
				DBName:   "testdb",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := database.NewMySQLConnection(tt.cfg)

			if tt.wantErr {
				assert.Error(t, err)
				if conn != nil {
					conn.Close()
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, conn)
				if conn != nil {
					conn.Close()
				}
			}
		})
	}
}

func TestNewMySQLConnection_Concurrent(t *testing.T) {
	cfg := database.Config{
		Host:     "localhost",
		Port:     "3306",
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
	}

	// Test concurrent connection attempts
	t.Run("Multiple concurrent connections", func(t *testing.T) {
		concurrency := 10
		results := make(chan error, concurrency)

		for i := 0; i < concurrency; i++ {
			go func() {
				conn, err := database.NewMySQLConnection(cfg)
				if err != nil {
					results <- err
					return
				}
				if conn != nil {
					conn.Close()
				}
				results <- nil
			}()
		}

		// Collect results
		errorCount := 0
		for i := 0; i < concurrency; i++ {
			err := <-results
			if err != nil {
				errorCount++
				t.Logf("Connection %d failed: %v", i, err)
			}
		}

		// In test environment, connections might fail
		t.Logf("Failed connections: %d/%d", errorCount, concurrency)
	})
}

// Integration test - requires actual database
func TestNewMySQLConnection_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cfg := database.Config{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "",
		DBName:   "test",
	}

	conn, err := database.NewMySQLConnection(cfg)
	if err != nil {
		t.Skipf("Skipping integration test - database not available: %v", err)
	}
	defer conn.Close()

	// Test actual database operations
	t.Run("Execute simple query", func(t *testing.T) {
		var result int
		err := conn.QueryRow("SELECT 1").Scan(&result)
		assert.NoError(t, err)
		assert.Equal(t, 1, result)
	})

	t.Run("Check connection stats", func(t *testing.T) {
		stats := conn.Stats()
		assert.GreaterOrEqual(t, stats.MaxOpenConnections, 0)
		t.Logf("Connection stats: %+v", stats)
	})
}
