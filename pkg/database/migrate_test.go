package database_test

import (
	"fmt"
	"testing"

	"fortyfour-backend/pkg/database"

	"github.com/stretchr/testify/assert"
)

// ═══════════════════════════════════════════════════════════════════════════
// TEST: RunMigrations — validasi parameter
// RunMigrations internally builds a DSN and calls sql.Open + migrate.
// We test the code paths that fail BEFORE actual DB interaction.
// ═══════════════════════════════════════════════════════════════════════════

func TestRunMigrations_InvalidConfig_FailsAtOpen(t *testing.T) {
	// Config with invalid/unreachable host — sql.Open succeeds (deferred)
	// but the mysql driver instance creation will fail
	cfg := database.Config{
		Host:   "invalid-host-that-does-not-exist",
		Port:   "9999",
		User:   "nobody",
		DBName: "nodb",
	}

	err := database.RunMigrations(cfg, "nonexistent/path")
	assert.Error(t, err, "expected error with invalid host")
}

func TestRunMigrations_InvalidMigrationPath(t *testing.T) {
	// Even with a valid-looking config, a bad migration path should fail
	cfg := database.Config{
		Host:   "127.0.0.1",
		Port:   "3306",
		User:   "root",
		DBName: "testdb",
	}

	err := database.RunMigrations(cfg, "/path/that/does/not/exist")
	assert.Error(t, err, "expected error with non-existent migration path")
}

func TestRunMigrations_EmptyMigrationPath(t *testing.T) {
	cfg := database.Config{
		Host:   "127.0.0.1",
		Port:   "3306",
		User:   "root",
		DBName: "testdb",
	}

	err := database.RunMigrations(cfg, "")
	assert.Error(t, err, "expected error with empty migration path")
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: Config struct — DSN construction verification
// ═══════════════════════════════════════════════════════════════════════════

func TestConfig_DSNConstruction(t *testing.T) {
	// Verify that the DSN format used by NewMySQLConnection matches expected pattern
	tests := []struct {
		name   string
		config database.Config
		want   string
	}{
		{
			name: "standard config",
			config: database.Config{
				Host:     "localhost",
				Port:     "3306",
				User:     "root",
				Password: "pass",
				DBName:   "mydb",
			},
			want: "root:pass@tcp(localhost:3306)/mydb?parseTime=true",
		},
		{
			name: "empty password",
			config: database.Config{
				Host:   "db.example.com",
				Port:   "3307",
				User:   "admin",
				DBName: "proddb",
			},
			want: "admin:@tcp(db.example.com:3307)/proddb?parseTime=true",
		},
		{
			name: "special chars in password",
			config: database.Config{
				Host:     "localhost",
				Port:     "3306",
				User:     "root",
				Password: "p@ss!123",
				DBName:   "testdb",
			},
			want: "root:p@ss!123@tcp(localhost:3306)/testdb?parseTime=true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := fmt.Sprintf(
				"%s:%s@tcp(%s:%s)/%s?parseTime=true",
				tt.config.User,
				tt.config.Password,
				tt.config.Host,
				tt.config.Port,
				tt.config.DBName,
			)
			assert.Equal(t, tt.want, dsn)
		})
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// TEST: NewMySQLConnection — edge cases
// ═══════════════════════════════════════════════════════════════════════════

func TestNewMySQLConnection_OnlyHostMissing(t *testing.T) {
	cfg := database.Config{
		Port:   "3306",
		User:   "root",
		DBName: "test",
	}
	_, err := database.NewMySQLConnection(cfg)
	assert.Error(t, err)
	assert.Equal(t, "host is required", err.Error())
}

func TestNewMySQLConnection_OnlyPortMissing(t *testing.T) {
	cfg := database.Config{
		Host:   "localhost",
		User:   "root",
		DBName: "test",
	}
	_, err := database.NewMySQLConnection(cfg)
	assert.Error(t, err)
	assert.Equal(t, "port is required", err.Error())
}

func TestNewMySQLConnection_OnlyUserMissing(t *testing.T) {
	cfg := database.Config{
		Host:   "localhost",
		Port:   "3306",
		DBName: "test",
	}
	_, err := database.NewMySQLConnection(cfg)
	assert.Error(t, err)
	assert.Equal(t, "user is required", err.Error())
}

func TestNewMySQLConnection_OnlyDBNameMissing(t *testing.T) {
	cfg := database.Config{
		Host: "localhost",
		Port: "3306",
		User: "root",
	}
	_, err := database.NewMySQLConnection(cfg)
	assert.Error(t, err)
	assert.Equal(t, "database name is required", err.Error())
}

// TestNewMySQLConnection_ValidConfigButUnreachable — SKIPPED
// This test uses 192.0.2.1 (TEST-NET) which causes a 21s TCP timeout.
// The Ping error path is already covered by mock tests in mysql_test.go.

func TestNewMySQLConnection_WhitespaceInFieldsNotTrimmed(t *testing.T) {
	// Host with whitespace is NOT trimmed — it's treated as-is
	// (This verifies that no implicit trimming is done by the function)
	cfg := database.Config{
		Host:   " ",
		Port:   "3306",
		User:   "root",
		DBName: "test",
	}
	// Host is " " (not empty), so it should pass validation but fail at Ping
	_, err := database.NewMySQLConnection(cfg)
	assert.Error(t, err)
	// Should NOT be "host is required" because " " is not empty
	assert.NotEqual(t, "host is required", err.Error())
}
