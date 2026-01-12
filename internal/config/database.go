package config

import (
	"fmt"
	"log"
	"os"
	"slices"
)

type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// SupportedDrivers lists all supported database drivers
var SupportedDrivers = []string{"mysql", "postgres"}

// Validate checks if the database configuration is valid
func (c *DatabaseConfig) Validate() error {
	// Check driver
	validDriver := slices.Contains(SupportedDrivers, c.Driver)
	if !validDriver {
		return fmt.Errorf("unsupported driver: %s, supported: %v", c.Driver, SupportedDrivers)
	}

	// Check required fields
	if c.User == "" || c.Host == "" || c.Port == 0 || c.DBName == "" {
		return fmt.Errorf("invalid database config: host, port, user, and dbname are required")
	}

	return nil
}

func LoadDatabaseOnly() (*DatabaseConfig, error) {
	v := newViper()

	bindDatabaseEnv(v)

	v.SetDefault("database.driver", "mysql")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.user", "root")
	v.SetDefault("database.password", "")
	v.SetDefault("database.dbname", "test")
	v.SetDefault("database.sslmode", "")

	// Don't read YAML anymore
	// readConfig(v)

	// Map directly to struct
	db := &DatabaseConfig{
		Driver:   v.GetString("database.driver"),
		Host:     v.GetString("database.host"),
		Port:     v.GetInt("database.port"),
		User:     v.GetString("database.user"),
		Password: v.GetString("database.password"),
		DBName:   v.GetString("database.dbname"),
		SSLMode:  v.GetString("database.sslmode"),
	}

	// Debug log
	log.Printf("[DEBUG] DB CONFIG LOADED: %+v\n", db)
	log.Printf("[DEBUG] ENV DB_HOST=%s DB_PORT=%s DB_USER=%s DB_PASSWORD=%s DB_NAME=%s DB_DRIVER=%s\n",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_DRIVER"),
	)

	// Validate configuration
	if err := db.Validate(); err != nil {
		return nil, err
	}

	return db, nil
}
