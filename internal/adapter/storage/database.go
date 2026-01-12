package storage

import (
	"database/sql"
	"fmt"

	"github.com/gostructure/app/internal/config"
)

// Database is the common interface for all SQL database connections
type Database interface {
	// GetDB returns the underlying *sql.DB connection
	GetDB() *sql.DB
	// Close closes the database connection
	Close() error
	// RunMigrations runs database migrations
	RunMigrations() error
	// DriverName returns the driver name (mysql, postgres, etc.)
	DriverName() string
}

// NewDatabase creates a database connection based on the driver type specified in config
func NewDatabase(cfg *config.DatabaseConfig, timezone string) (Database, error) {
	switch cfg.Driver {
	case "mysql":
		return NewMySQLDatabase(cfg, timezone)
	case "postgres":
		return NewPostgresDatabase(cfg, timezone)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s (supported: mysql, postgres)", cfg.Driver)
	}
}
