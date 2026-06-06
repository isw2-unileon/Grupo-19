package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	// Gin test mode
	gin.SetMode(gin.TestMode)

	// Petition without cookie
	t.Run("Falta de cookie rechaza la petición", func(t *testing.T) {
		r := gin.Default()
		r.Use(AuthMiddleware())
		r.GET("/ruta-protegida", func(c *gin.Context) {
			c.String(http.StatusOK, "¡Entraste!")
		})

		req, _ := http.NewRequest("GET", "/ruta-protegida", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Esperábamos que la petición fuera rechazada (401), pero obtuvimos %d", w.Code)
		}
	})

	// Petition with cookie with invalid token
	t.Run("Token manipulado rechaza la petición", func(t *testing.T) {
		r := gin.Default()
		r.Use(AuthMiddleware())
		r.GET("/ruta-protegida", func(c *gin.Context) {
			c.String(http.StatusOK, "¡Entraste!")
		})

		req, _ := http.NewRequest("GET", "/ruta-protegida", nil)
		req.AddCookie(&http.Cookie{
			Name:  "auth_token",
			Value: "token.falso.inventado",
		})

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Esperábamos que la petición fuera rechazada (401), pero obtuvimos %d", w.Code)
		}
	})

	// Valid token
	t.Run("Token correcto permite el paso e inyecta el contexto", func(t *testing.T) {
		r := gin.Default()
		r.Use(AuthMiddleware())
		r.GET("/ruta-protegida", func(c *gin.Context) {
			userID, existsID := c.Get("userID")
			userType, existsType := c.Get("userType")

			if !existsID || userID != uint(77) {
				t.Errorf("El middleware no inyectó correctamente el userID")
			}
			if !existsType || userType != "premium" {
				t.Errorf("El middleware no inyectó correctamente el userType")
			}

			c.String(http.StatusOK, "¡Entraste!")
		})

		tokenReal, _ := GenerateToken(77, "premium")

		req, _ := http.NewRequest("GET", "/ruta-protegida", nil)
		req.AddCookie(&http.Cookie{
			Name:  "auth_token",
			Value: tokenReal,
		})

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Esperábamos que la petición pasara (200), pero fue rechazada con %d", w.Code)
		}
	})
}
