package config

import (
	"time"
)

type Config struct {
	App    AppConfig
	Server ServerConfig
	JWT    JWTConfig
}

// Load dùng cho application server
func Load() (*Config, error) {
	v := newViper()

	// Defaults
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.timezone", "Asia/Ho_Chi_Minh")

	v.SetDefault("server.address", ":8080")
	v.SetDefault("server.read_timeout", 15*time.Second)
	v.SetDefault("server.write_timeout", 15*time.Second)
	v.SetDefault("server.idle_timeout", 60*time.Second)

	v.SetDefault("jwt.secret", "secret")
	v.SetDefault("jwt.expiration", 24*time.Hour)

	bindCommonEnv(v)

	if err := readConfig(v); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
