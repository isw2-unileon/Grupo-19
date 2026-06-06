package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// setupTestDB build a temp database for our tests
func setupTestDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Error al conectar con la base de datos de prueba")
	}

	err2 := db.AutoMigrate(&models.User{})
	if err2 != nil {
		panic("Error al migrar la base de datos de prueba")
	}

	database.DB = db
}

// mockAuthMiddleware simulate real JWT tokens
func mockAuthMiddleware(userID uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	}
}

func TestLoginHandler(t *testing.T) {
	setupTestDB()
	gin.SetMode(gin.TestMode)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
	testUser := models.User{
		Username: "usuario_test",
		Email:    "test@test.com",
		Password: string(hashedPassword),
	}
	database.DB.Create(&testUser)

	tests := []struct {
		name           string
		payload        LoginRequest
		expectedStatus int
	}{
		{"Caso 1: Login exitoso", LoginRequest{Email: "test@test.com", Password: "1234"}, http.StatusOK},
		{"Caso 2: Contraseña incorrecta", LoginRequest{Email: "test@test.com", Password: "mala_contraseña"}, http.StatusUnauthorized},
		{"Caso 3: Usuario no existe", LoginRequest{Email: "noexiste@test.com", Password: "1234"}, http.StatusUnauthorized},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.Default()
			r.POST("/api/login", LoginHandler)

			jsonValue, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonValue))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("Esperábamos %d, pero obtuvimos %d", tc.expectedStatus, w.Code)
			}
		})
	}
}

func TestRegisterHandler(t *testing.T) {
	setupTestDB()
	gin.SetMode(gin.TestMode)

	database.DB.Create(&models.User{Username: "usuario_existente", Email: "duplicado@test.com"})

	tests := []struct {
		name           string
		payload        RegisterRequest
		expectedStatus int
	}{
		{"Caso 1: Registro exitoso", RegisterRequest{Username: "nuevo_user", Email: "nuevo@test.com", Password: "1234"}, http.StatusCreated},
		{"Caso 2: Correo ya registrado", RegisterRequest{Username: "otro_user", Email: "duplicado@test.com", Password: "1234"}, http.StatusConflict},
		{"Caso 3: Nombre de usuario ya registrado", RegisterRequest{Username: "usuario_existente", Email: "otro_correo@test.com", Password: "1234"}, http.StatusConflict},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.Default()
			r.POST("/api/register", RegisterHandler)

			jsonValue, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonValue))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("Esperábamos %d, pero obtuvimos %d", tc.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetProfileHandler(t *testing.T) {
	setupTestDB()
	gin.SetMode(gin.TestMode)

	testUser := models.User{Username: "perfil_user", Email: "perfil@test.com"}
	database.DB.Create(&testUser)

	t.Run("Caso 1: Obtener perfil con éxito", func(t *testing.T) {
		r := gin.Default()
		// Inyectamos el middleware falso para simular que estamos logueados
		r.Use(mockAuthMiddleware(testUser.UserID))
		r.GET("/api/user/profile", GetProfileHandler)

		req, _ := http.NewRequest("GET", "/api/user/profile", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Esperábamos 200, obtuvimos %d", w.Code)
		}
	})

	t.Run("Caso 2: Usuario no autenticado", func(t *testing.T) {
		r := gin.Default()
		r.GET("/api/user/profile", GetProfileHandler) // Sin middleware

		req, _ := http.NewRequest("GET", "/api/user/profile", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Esperábamos 401, obtuvimos %d", w.Code)
		}
	})
}

func TestUpdateProfileHandler(t *testing.T) {
	setupTestDB()
	gin.SetMode(gin.TestMode)

	user1 := models.User{Username: "usuario_1", Email: "user1@test.com"}
	user2 := models.User{Username: "usuario_2", Email: "user2@test.com"}
	database.DB.Create(&user1)
	database.DB.Create(&user2)

	tests := []struct {
		name           string
		payload        UpdateProfileRequest
		expectedStatus int
	}{
		{"Caso 1: Actualización exitosa", UpdateProfileRequest{Username: "usuario_1_nuevo", Email: "nuevo@test.com"}, http.StatusOK},
		{"Caso 2: Conflicto con email de otro", UpdateProfileRequest{Username: "usuario_1", Email: "user2@test.com"}, http.StatusConflict},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.Default()
			r.Use(mockAuthMiddleware(user1.UserID))
			r.PUT("/api/user/profile", UpdateProfileHandler)

			jsonValue, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("PUT", "/api/user/profile", bytes.NewBuffer(jsonValue))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("Esperábamos %d, pero obtuvimos %d", tc.expectedStatus, w.Code)
			}
		})
	}
}

func TestUpdatePasswordHandler(t *testing.T) {
	setupTestDB()
	gin.SetMode(gin.TestMode)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password_vieja"), bcrypt.DefaultCost)
	testUser := models.User{Username: "pass_user", Password: string(hashedPassword)}
	database.DB.Create(&testUser)

	tests := []struct {
		name           string
		payload        UpdatePasswordRequest
		expectedStatus int
	}{
		{"Caso 1: Cambio exitoso", UpdatePasswordRequest{CurrentPassword: "password_vieja", NewPassword: "nueva_password"}, http.StatusOK},
		{"Caso 2: Contraseña actual incorrecta", UpdatePasswordRequest{CurrentPassword: "mala", NewPassword: "nueva_password"}, http.StatusUnauthorized},
		{"Caso 3: Nueva contraseña muy corta", UpdatePasswordRequest{CurrentPassword: "password_vieja", NewPassword: "123"}, http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.Default()
			r.Use(mockAuthMiddleware(testUser.UserID))
			r.PUT("/api/user/profile/password", UpdatePasswordHandler)

			jsonValue, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("PUT", "/api/user/profile/password", bytes.NewBuffer(jsonValue))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("Esperábamos %d, pero obtuvimos %d en %s", tc.expectedStatus, w.Code, tc.name)
			}
		})
	}
}

func TestLogoutHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/api/logout", LogoutHandler)

	req, _ := http.NewRequest("POST", "/api/logout", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Esperábamos 200, obtuvimos %d", w.Code)
	}

	// Comprobamos que el servidor manda destruir la cookie
	cookie := w.Header().Get("Set-Cookie")
	if cookie == "" {
		t.Errorf("Esperábamos que el servidor enviara un Set-Cookie para borrar el token")
	}
}
