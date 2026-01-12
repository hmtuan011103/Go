package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gostructure/app/internal/config"
	_ "github.com/lib/pq"
)

// PostgresDatabase implements the Database interface for PostgreSQL
type PostgresDatabase struct {
	db  *sql.DB
	cfg *config.DatabaseConfig
}

// NewPostgresDatabase creates a new PostgreSQL database connection
func NewPostgresDatabase(cfg *config.DatabaseConfig, timezone string) (*PostgresDatabase, error) {
	// Build connection string for PostgreSQL
	sslMode := cfg.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		sslMode,
		timezone,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Connected to PostgreSQL successfully")
	return &PostgresDatabase{db: db, cfg: cfg}, nil
}

// GetDB returns the underlying *sql.DB connection
func (p *PostgresDatabase) GetDB() *sql.DB {
	return p.db
}

// Close closes the database connection
func (p *PostgresDatabase) Close() error {
	return p.db.Close()
}

// DriverName returns "postgres"
func (p *PostgresDatabase) DriverName() string {
	return "postgres"
}

// RunMigrations runs database migrations for PostgreSQL
func (p *PostgresDatabase) RunMigrations() error {
	driver, err := postgres.WithInstance(p.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres migration driver: %w", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		"file://db/postgres/migrations",
		p.cfg.DBName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("PostgreSQL migrations completed successfully")
	return nil
}
