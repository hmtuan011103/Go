package config

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func newViper() *viper.Viper {
	loadEnv()

	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return v
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using env vars")
	}
}

func readConfig(v *viper.Viper) error {
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	return nil
}

func bindCommonEnv(v *viper.Viper) {
	// App
	v.BindEnv("app.name", "APP_NAME")
	v.BindEnv("app.environment", "APP_ENV")
	v.BindEnv("app.timezone", "APP_TIMEZONE")

	// Server
	v.BindEnv("server.address", "SERVER_ADDRESS")
	v.BindEnv("server.read_timeout", "SERVER_READ_TIMEOUT")
	v.BindEnv("server.write_timeout", "SERVER_WRITE_TIMEOUT")
	v.BindEnv("server.idle_timeout", "SERVER_IDLE_TIMEOUT")

	// JWT
	v.BindEnv("jwt.secret", "JWT_SECRET")
	v.BindEnv("jwt.expiration", "JWT_EXPIRATION")
}

func bindDatabaseEnv(v *viper.Viper) {
	v.BindEnv("database.driver", "DB_DRIVER")
	v.BindEnv("database.host", "DB_HOST")
	v.BindEnv("database.port", "DB_PORT")
	v.BindEnv("database.user", "DB_USER")
	v.BindEnv("database.password", "DB_PASSWORD")
	v.BindEnv("database.dbname", "DB_NAME")
	v.BindEnv("database.sslmode", "DB_SSLMODE")
}
