package config

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env files based on the runtime environment

func LoadEnv() string {
	_ = godotenv.Load(".env")

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	envFile := ".env." + env
	godotenv.Load(envFile)

	return env
}
