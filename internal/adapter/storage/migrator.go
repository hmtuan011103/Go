package storage

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gostructure/app/internal/config"
)

// Migrator wraps the migrate.Migrate instance
type Migrator struct {
	m *migrate.Migrate
}

// NewMigrator creates a migrator based on the database driver
func NewMigrator(db *sql.DB, cfg *config.DatabaseConfig) (*Migrator, error) {
	var m *migrate.Migrate
	var err error

	migrationPath := fmt.Sprintf("file://db/%s/migrations", cfg.Driver)

	switch cfg.Driver {
	case "mysql":
		driver, dErr := mysql.WithInstance(db, &mysql.Config{})
		if dErr != nil {
			return nil, fmt.Errorf("create mysql migration driver failed: %w", dErr)
		}
		m, err = migrate.NewWithDatabaseInstance(migrationPath, cfg.DBName, driver)
	case "postgres":
		driver, dErr := postgres.WithInstance(db, &postgres.Config{})
		if dErr != nil {
			return nil, fmt.Errorf("create postgres migration driver failed: %w", dErr)
		}
		m, err = migrate.NewWithDatabaseInstance(migrationPath, cfg.DBName, driver)
	default:
		return nil, fmt.Errorf("unsupported driver for migration: %s", cfg.Driver)
	}
	if err != nil {
		return nil, fmt.Errorf("create migrate instance failed: %w", err)
	}

	return &Migrator{m: m}, nil
}

func (mg *Migrator) Up() error {
	if err := mg.m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func (mg *Migrator) Down() error {
	if err := mg.m.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func (mg *Migrator) Version() (uint, bool, error) {
	return mg.m.Version()
}

func (mg *Migrator) Close() error {
	srcErr, dbErr := mg.m.Close()
	if srcErr != nil {
		return srcErr
	}
	return dbErr
}
