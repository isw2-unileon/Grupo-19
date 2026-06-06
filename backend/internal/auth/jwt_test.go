package auth

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGetJWTSecret(t *testing.T) {
	// Test what happend if there is not enviroment variable
	os.Unsetenv("JWT_SECRET")
	secret := getJWTSecret()
	if string(secret) != "clave_local_por_defecto_grupo19" {
		t.Errorf("Esperaba la clave por defecto, pero obtuve: %s", string(secret))
	}

	os.Setenv("JWT_SECRET", "mi_clave_super_secreta_test")
	defer os.Unsetenv("JWT_SECRET") // Clean after test

	secretConEnv := getJWTSecret()
	if string(secretConEnv) != "mi_clave_super_secreta_test" {
		t.Errorf("Esperaba 'mi_clave_super_secreta_test', pero obtuve: %s", string(secretConEnv))
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	userID := uint(42)
	userType := "admin"

	tokenString, err := GenerateToken(userID, userType)
	if err != nil {
		t.Fatalf("Error inesperado al generar el token: %v", err)
	}
	if tokenString == "" {
		t.Fatalf("El token generado está completamente vacío")
	}

	claims, err := ValidateToken(tokenString)
	if err != nil {
		t.Fatalf("Error inesperado al validar un token que debería ser correcto: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("Esperaba el UserID %d, pero obtuve %d", userID, claims.UserID)
	}
	if claims.UserType != userType {
		t.Errorf("Esperaba el UserType '%s', pero obtuve '%s'", userType, claims.UserType)
	}
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	invalidToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VyX3R5cGUiOiJ1c2VyIn0.firma_falsa_e_inventada"

	_, err := ValidateToken(invalidToken)
	if err == nil {
		t.Error("Se esperaba un error al intentar validar un token con firma falsa, pero pasó como válido")
	}
}

func TestValidateToken_Expired(t *testing.T) {
	claims := CustomClaims{
		UserID:   1,
		UserType: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredTokenString, _ := token.SignedString(getJWTSecret())

	_, err := ValidateToken(expiredTokenString)
	if err == nil {
		t.Error("Se esperaba un error al validar un token caducado, pero el sistema lo aceptó")
	}
}
