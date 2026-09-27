package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all the environment-based settings for the app.
type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	AdminUsername string
	AdminPassword string
}

// Load reads the .env file (if present) and returns a populated Config.
// If a variable is missing, a sensible default is used so the app can
// still boot during local development.
func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "manadgasht"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		AdminUsername: getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "change-me-now"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
