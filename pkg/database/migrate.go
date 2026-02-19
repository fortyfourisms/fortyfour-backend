package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations creates a dedicated connection with multiStatements enabled
// (required by golang-migrate for MySQL) and runs all pending migrations.
// If the database is in a dirty state, it will attempt to force the version
// and retry automatically.
func RunMigrations(cfg Config, migrationsPath string) error {
	// Build DSN with multiStatements=true for migration support
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("could not open migration db connection: %w", err)
	}
	defer db.Close()

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	err = m.Up()
	if err == nil || err == migrate.ErrNoChange {
		version, dirty, _ := m.Version()
		log.Printf("✅ Database migrated to version %d (dirty: %v)", version, dirty)
		return nil
	}

	// Handle dirty database state: force version back and retry
	version, dirty, verr := m.Version()
	if verr == nil && dirty {
		log.Printf("⚠️  Database is dirty at version %d, attempting to recover...", version)

		// Force back to the previous version (the last clean state)
		prevVersion := int(version) - 1
		if prevVersion < 0 {
			prevVersion = -1 // -1 means "no version" in golang-migrate
		}

		if ferr := m.Force(prevVersion); ferr != nil {
			return fmt.Errorf("migration failed and could not recover dirty state: original=%w, force=%v", err, ferr)
		}

		log.Printf("✅ Forced database to version %d, retrying migrations...", prevVersion)

		// Retry migration
		if retryErr := m.Up(); retryErr != nil && retryErr != migrate.ErrNoChange {
			return fmt.Errorf("migration retry failed: %w", retryErr)
		}

		version, dirty, _ = m.Version()
		log.Printf("✅ Database migrated to version %d (dirty: %v)", version, dirty)
		return nil
	}

	return fmt.Errorf("migration failed: %w", err)
}
