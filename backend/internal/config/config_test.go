package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	os.Unsetenv("TEST_KEY_NO_EXISTE")
	if val := getEnv("TEST_KEY_NO_EXISTE", "fallback_default"); val != "fallback_default" {
		t.Errorf("Esperaba 'fallback_default', obtuve '%s'", val)
	}

	os.Setenv("TEST_KEY_EXISTE", "valor_real")
	defer os.Unsetenv("TEST_KEY_EXISTE")

	if val := getEnv("TEST_KEY_EXISTE", "fallback_default"); val != "valor_real" {
		t.Errorf("Esperaba 'valor_real', obtuve '%s'", val)
	}

	os.Setenv("TEST_KEY_VACIA", "")
	defer os.Unsetenv("TEST_KEY_VACIA")

	if val := getEnv("TEST_KEY_VACIA", "fallback_default"); val != "fallback_default" {
		t.Errorf("Esperaba 'fallback_default' para una cadena vacía, obtuve '%s'", val)
	}
}

func TestLoad(t *testing.T) {
	variables := []string{"PORT", "GIN_MODE", "CORS_ALLOW_ORIGIN", "DATABASE_URL", "JWT_SECRET"}
	backup := make(map[string]string)
	for _, v := range variables {
		backup[v] = os.Getenv(v)
		os.Unsetenv(v)
	}

	defer func() {
		for key, val := range backup {
			if val != "" {
				os.Setenv(key, val)
			}
		}
	}()

	t.Run("Carga de valores por defecto (Fallbacks)", func(t *testing.T) {
		cfg := Load()
		if cfg.Port != "8080" {
			t.Errorf("Esperaba puerto '8080', obtuve '%s'", cfg.Port)
		}
		if cfg.GinMode != "debug" {
			t.Errorf("Esperaba modo 'debug', obtuve '%s'", cfg.GinMode)
		}
		if cfg.JWTSecret != "clave_temporal_local_grupo19_2026" {
			t.Errorf("Esperaba secreto por defecto, obtuve '%s'", cfg.JWTSecret)
		}
	})

	t.Run("Carga con variables de entorno modificadas", func(t *testing.T) {
		os.Setenv("PORT", "9090")
		os.Setenv("GIN_MODE", "release")
		os.Setenv("JWT_SECRET", "super_secreto_test")

		cfg := Load()

		if cfg.Port != "9090" {
			t.Errorf("Esperaba puerto '9090', obtuve '%s'", cfg.Port)
		}
		if cfg.GinMode != "release" {
			t.Errorf("Esperaba modo 'release', obtuve '%s'", cfg.GinMode)
		}
		if cfg.JWTSecret != "super_secreto_test" {
			t.Errorf("Esperaba secreto 'super_secreto_test', obtuve '%s'", cfg.JWTSecret)
		}

		os.Unsetenv("PORT")
		os.Unsetenv("GIN_MODE")
		os.Unsetenv("JWT_SECRET")
	})
}
