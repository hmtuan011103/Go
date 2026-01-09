package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gostructure/app/internal/config"
)

func NewMySQLConnection(cfg *config.DatabaseConfig, timezone string) (*sql.DB, error) {
	// 1. Convert named timezone (e.g., Asia/Ho_Chi_Minh) to offset (e.g., +07:00)
	// This ensures compatibility with MySQL even if timezone tables are not populated (common on Windows).
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone %s in config: %w", timezone, err)
	}

	_, offsetSeconds := time.Now().In(loc).Zone()
	offsetHours := offsetSeconds / 3600
	offsetMinutes := (offsetSeconds % 3600) / 60
	if offsetMinutes < 0 {
		offsetMinutes = -offsetMinutes
	}
	tzOffset := fmt.Sprintf("%+03d:%02d", offsetHours, offsetMinutes)

	// 2. Build DSN using the calculated offset for the database session
	escapedTz := url.QueryEscape(timezone)
	escapedOffset := url.QueryEscape(tzOffset)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=%s&time_zone=%%27%s%%27&multiStatements=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		escapedTz,
		escapedOffset,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Connected to MySQL successfully")
	return db, nil
}
