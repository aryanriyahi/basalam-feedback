package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	DatabaseURL       string
	BasicAuthUser     string
	BasicAuthPassword string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to system environment variables")
	}

	return Config{
		Port:              getEnv("PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://postgres:postgres@db:5432/feedbacks?sslmode=disable"),
		BasicAuthUser:     getEnv("BASIC_AUTH_USER", "admin"),
		BasicAuthPassword: getEnv("BASIC_AUTH_PASSWORD", "admin"),
	}
}

func (c Config) Addr() string {
	return ":" + c.Port
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
