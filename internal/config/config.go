package config

import (
	"os"
	"strconv"

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

type Config struct {
	Env       string
	Port      string
	APIKey    string
	LimitRate float64
	DSN       string
}

func MustLoad() Config {
	LoadEnv()

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		panic("API_KEY not set")
	}

	limitRateStr := os.Getenv("LIMIT_RATE")
	if limitRateStr == "" {
		panic("LIMIT_RATE not set")
	}

	limitRate, err := strconv.ParseFloat(limitRateStr, 64)
	if err != nil {
		panic("invalid LIMIT_RATE: " + err.Error())
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		panic("DSN not set")
	}

	return Config{
		Env:       env,
		Port:      port,
		APIKey:    apiKey,
		LimitRate: limitRate,
		DSN:       dsn,
	}
}
