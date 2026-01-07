package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gostructure/app/internal/config"
)

// RunMigrations executes the database migrations
func RunMigrations(db *sql.DB, cfg *config.DatabaseConfig) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("could not create mysql driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		cfg.DBName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migration instance: %w", err)
	}

	// Get current version
	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("could not get migration version: %w", err)
	}
	log.Printf("Current migration version: %d, dirty: %v", version, dirty)

	// Verify if there are migrations to run
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Database is up to date regarding migrations")
			return nil
		}
		return fmt.Errorf("could not run up migrations: %w", err)
	}

	log.Println("Successfully ran database migrations")
	return nil
}
