package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware intercepts the requests, reads the cookie and saves the UserID in the context
//
//nolint:revive
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// We try to read the cookie that we saved in the login
		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Inicie sesión para acceder a este recurso"})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesión inválida o expirada"})
			c.Abort()
			return
		}

		// Injection of the UserID in the context of the request
		c.Set("userID", claims.UserID)
		c.Set("userType", claims.UserType)

		// Allow the request to continue to the final controller
		c.Next()
	}
}
