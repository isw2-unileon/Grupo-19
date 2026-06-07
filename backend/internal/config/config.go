// Package config handles application configuration from environment variables.
package config

import (
	"os"

	"github.com/joho/godotenv" // 1. Importamos godotenv
)

// Config holds the application configuration loaded from environment variables.
type Config struct {
	Port            string
	GinMode         string
	CORSAllowOrigin string

	DatabaseURL string
	JWTSecret   string // Added the field for the token's secret key
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:            getEnv("PORT", "8080"),
		GinMode:         getEnv("GIN_MODE", "debug"),
		CORSAllowOrigin: getEnv("CORS_ALLOW_ORIGIN", "*"),

		DatabaseURL: getEnv("DATABASE_URL", "postgres://admin:secretpassword@db:5432/rastreador_precios?sslmode=disable"),

		// Add secret key reading with a safe default value for local
		JWTSecret: getEnv("JWT_SECRET", "clave_temporal_local_grupo19_2026"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
