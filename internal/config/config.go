package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	JWTSecret      string
	TokenTTL       time.Duration
	FrontendOrigin string
	CookieSecure   bool
}

func Load() (*Config, error) {
	godotenv.Load()

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	ttlHours := 24
	if raw := os.Getenv("TOKEN_TTL_HOURS"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 {
			return nil, fmt.Errorf("TOKEN_TTL_HOURS must be a positive integer")
		}
		ttlHours = parsed
	}

	return &Config{
		AppPort:        getEnv("APP_PORT", "8080"),
		JWTSecret:      secret,
		TokenTTL:       time.Duration(ttlHours) * time.Hour,
		FrontendOrigin: getEnv("FRONTEND_ORIGIN", "http://localhost:3000"),
		CookieSecure:   os.Getenv("COOKIE_SECURE") == "true",
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
