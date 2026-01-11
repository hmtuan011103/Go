package config

import (
	"fmt"
	"log"
	"os"
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
	log.Printf("[DEBUG] ENV DB_HOST=%s DB_PORT=%s DB_USER=%s DB_PASSWORD=%s DB_NAME=%s\n",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// Validate information
	if db.User == "" || db.Host == "" || db.Port == 0 || db.DBName == "" {
		return nil, fmt.Errorf("invalid database config: %+v", db)
	}

	return db, nil
}
