package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const envPath = "./.env"

// DB represents the configuration for the database.
type DB struct {
	DBName   string `env:"POSTGRES_DB" env-required:"true"`
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     int    `env:"DB_PORT" env-required:"true"`
	DSN      string `env:"-"`
}

// HTTP represents the configuration for the http server.
type HTTP struct {
	Port    int    `env:"HTTP_PORT" env-required:"true"`
	Host    string `env:"HTTP_HOST" env-required:"true"`
	Address string `env:"-"`
}

// Config represents the overall application configuration.
type Config struct {
	DB   DB
	HTTP HTTP
}

// Load reads configuration from .env file.
func Load() (*Config, error) {
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return nil, fmt.Errorf(".env file does not exist in project's root")
	}

	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("cannot read config from environment variables: %w", err)
	}

	cfg.DB.DSN = fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.DBName,
		cfg.DB.User,
		cfg.DB.Password,
	)

	cfg.HTTP.Address = fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)

	return &cfg, nil
}
