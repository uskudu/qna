package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	HTTPPort    string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using environment variables")
	}

	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		HTTPPort:    os.Getenv("HTTP_PORT"),
	}
}
