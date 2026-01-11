package mysql

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gostructure/app/internal/config"
)

type Migrator struct {
	m *migrate.Migrate
}

func NewMigrator(db *sql.DB, cfg *config.DatabaseConfig) (*Migrator, error) {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return nil, fmt.Errorf("create mysql driver failed: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		cfg.DBName,
		driver,
	)
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
