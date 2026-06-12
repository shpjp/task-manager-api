package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort             string
	JWTSecret           string
	TokenTTL            time.Duration
	FrontendOrigin      string
	CookieSecure        bool
	AdminEmails         []string
	MaxUploadMB         int64
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
	CloudinaryFolder    string
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

	maxUploadMB := int64(5)
	if raw := os.Getenv("MAX_UPLOAD_MB"); raw != "" {
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || parsed < 1 {
			return nil, fmt.Errorf("MAX_UPLOAD_MB must be a positive integer")
		}
		maxUploadMB = parsed
	}

	var adminEmails []string
	if raw := os.Getenv("ADMIN_EMAILS"); raw != "" {
		adminEmails = strings.Split(raw, ",")
	}

	return &Config{
		AppPort:             getEnv("APP_PORT", "8080"),
		JWTSecret:           secret,
		TokenTTL:            time.Duration(ttlHours) * time.Hour,
		FrontendOrigin:      getEnv("FRONTEND_ORIGIN", "http://localhost:3000"),
		CookieSecure:        os.Getenv("COOKIE_SECURE") == "true",
		AdminEmails:         adminEmails,
		MaxUploadMB:         maxUploadMB,
		CloudinaryCloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryAPIKey:    os.Getenv("CLOUDINARY_API_KEY"),
		CloudinaryAPISecret: os.Getenv("CLOUDINARY_API_SECRET"),
		CloudinaryFolder:    getEnv("CLOUDINARY_FOLDER", "tasktheteam"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
