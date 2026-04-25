package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all configuration values loaded once at startup.
// Every field that is required will cause Load() to fail if absent.
type Config struct {
	// HTTP
	Port string

	// Database (Supabase / Postgres)
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSLMode  string // "require" in prod, "disable" for local Postgres without TLS

	// Auth
	JWTSecret      string
	InternalAPIKey string
}

// Load reads environment variables, validates required fields, and returns
// a populated Config. Fail fast: if anything required is missing the caller
// should exit immediately.
func Load() (*Config, error) {
	cfg := &Config{
		Port:           getOrDefault("PORT", "8080"),
		DBHost:         os.Getenv("DB_HOST"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         getOrDefault("DB_NAME", "postgres"),
		DBPort:         getOrDefault("DB_PORT", "5432"),
		DBSSLMode:      getOrDefault("DB_SSLMODE", "require"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		InternalAPIKey: os.Getenv("INTERNAL_API_KEY"),
	}

	var missing []string
	required := map[string]string{
		"DB_HOST":          cfg.DBHost,
		"DB_USER":          cfg.DBUser,
		"DB_PASSWORD":      cfg.DBPassword,
		"JWT_SECRET":       cfg.JWTSecret,
		"INTERNAL_API_KEY": cfg.InternalAPIKey,
	}
	for k, v := range required {
		if v == "" {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}

func getOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
