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
	JWTSecret   string // 2. Añadimos el campo para la clave secreta del token
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:            getEnv("PORT", "8080"),
		GinMode:         getEnv("GIN_MODE", "debug"),
		CORSAllowOrigin: getEnv("CORS_ALLOW_ORIGIN", "*"),

		DatabaseURL: getEnv("DATABASE_URL", "postgres://admin:secretpassword@db:5432/rastreador_precios?sslmode=disable"),

		// Añadimos la lectura de la clave secreta con un valor por defecto seguro para local
		JWTSecret: getEnv("JWT_SECRET", "clave_temporal_local_grupo19_2026"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
