package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	ServerPort  string
	CORSOrigins string
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func LoadConfig() *Config {
	// Load .env file (ignore error if not found, e.g. in Docker)
	_ = godotenv.Load()

	return &Config{
		DBHost:	  getEnv("DB_HOST", "localhost"),
		DBPort:	  getEnv("DB_PORT", "5432"),
		DBUser:	  getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName: getEnv("DB_NAME", "menus_db"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		CORSOrigins: getEnv("CORS_ORIGINS", "http://localhost:3000"),
	}
}
