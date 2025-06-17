package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// AppConfig holds the application configuration
var AppConfig struct {
	DBSource         string
	Port             string
	JWTSecret        string
	JWTRefreshSecret string
}

// LoadConfig loads configuration from .env file
func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	AppConfig.DBSource = os.Getenv("DB_DSN")
	if AppConfig.DBSource == "" {
		log.Fatal("DB_DSN not set in .env file")
	}

	AppConfig.Port = os.Getenv("PORT")
	if AppConfig.Port == "" {
		log.Fatal("PORT not set in .env file")
	}

	AppConfig.JWTSecret = os.Getenv("JWT_SECRET")
	if AppConfig.JWTSecret == "" {
		log.Fatal("JWT_SECRET not set in .env file")
	}

	AppConfig.JWTRefreshSecret = os.Getenv("JWT_REFRESH_SECRET")
	if AppConfig.JWTRefreshSecret == "" {
		log.Fatal("JWT_REFRESH_SECRET not set in .env file")
	}
}

// GetConfig returns a pointer to the AppConfig struct
func GetConfig() *struct {
	DBSource         string
	Port             string
	JWTSecret        string
	JWTRefreshSecret string
} {
	return &AppConfig
}
