package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware intercepta las peticiones, lee la cookie y guarda el UserID en el contexto //nolint:revive
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Intentamos leer la cookie que guardamos en el login
		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Inicie sesión para acceder a este recurso"})
			c.Abort()
			return
		}

		// Validamos el token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesión inválida o expirada"})
			c.Abort()
			return
		}

		// Inyectamos el UserID en el contexto de la petición
		c.Set("userID", claims.UserID)
		c.Set("userType", claims.UserType)

		// Permitimos que la petición continúe hacia el controlador final
		c.Next()
	}
}
