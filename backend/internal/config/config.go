package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         int
	Env          string
	DatabasePath string

	PayPalClientID string
	PayPalSecret   string
	PayPalBaseURL  string

	AnthropicAPIKey string

	FrontendURL string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	port, _ := strconv.Atoi(getEnv("PORT", "8080"))

	return &Config{
		Port:            port,
		Env:             getEnv("ENV", "development"),
		DatabasePath:    getEnv("DATABASE_PATH", "./paypilot.db"),
		PayPalClientID:  getEnv("PAYPAL_CLIENT_ID", ""),
		PayPalSecret:    getEnv("PAYPAL_SECRET", ""),
		PayPalBaseURL:   getEnv("PAYPAL_BASE_URL", "https://api-m.sandbox.paypal.com"),
		AnthropicAPIKey: getEnv("ANTHROPIC_API_KEY", ""),
		FrontendURL:     getEnv("FRONTEND_URL", "http://localhost:5173"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
