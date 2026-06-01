package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Clave de respaldo por si acaso no se ha cargado el entorno en algún test
		return []byte("clave_local_por_defecto_grupo19")
	}
	return []byte(secret)
}

// CustomClaims define la estructura de la información que meteremos dentro del token
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}

// GenerateToken crea un token JWT firmado para un usuario que acaba de loguearse con éxito
func GenerateToken(userID uint, userType string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // El token expira en 24 horas
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(getJWTSecret())
}

// ValidateToken recibe el string del token, comprueba si la firma es válida y extrae sus datos
func ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}

	// Verificación
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}
