package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
	RedisURL    string
	LogLevel    string
	Environment string
}

func Load() (*Config, error) {
	// Try loading .env from multiple locations (order matters: closest first)
	// godotenv.Load loads all files that exist, ignores missing ones
	godotenv.Load("../.env", "../../.env", ".env")

	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost/taskdb?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "localhost:6379"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
