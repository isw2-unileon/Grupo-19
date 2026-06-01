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

// setupTestDB crea una base de datos temporal en blanco y la conecta a tu variable global
func setupTestDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Error al conectar con la base de datos de prueba")
	}

	// Creamos las tablas necesarias
	err2 := db.AutoMigrate(&models.User{})
	if err2 != nil {
		panic("Error al migrar la base de datos de prueba")
	}

	// Sustituimos la base de datos real por la de prueba
	database.DB = db
}

func TestLoginHandler(t *testing.T) {
	setupTestDB()

	// Preparamos un usuario de prueba en la base de datos temporal
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
	testUser := models.User{
		Username: "usuario_test",
		Email:    "test@test.com",
		Password: string(hashedPassword),
	}
	database.DB.Create(&testUser)

	// Cambiamos Gin a modo "test" para que no sature la terminal de mensajes
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		payload        LoginRequest
		expectedStatus int
	}{
		{
			name:           "Caso 1: Login exitoso",
			payload:        LoginRequest{Email: "test@test.com", Password: "1234"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Caso 2: Contraseña incorrecta",
			payload:        LoginRequest{Email: "test@test.com", Password: "mala_contraseña"},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Caso 3: Usuario no existe",
			payload:        LoginRequest{Email: "noexiste@test.com", Password: "1234"},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Preparamos el servidor de pruebas con nuestra ruta
			r := gin.Default()
			r.POST("/api/login", LoginHandler)

			// Convertimos los datos de prueba a JSON
			jsonValue, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonValue))

			// Grabadora que captura la respuesta del servidor
			w := httptest.NewRecorder()

			// Ejecutamos la petición
			r.ServeHTTP(w, req)

			// Comprobamos si el código de estado es el que esperábamos (Ej: 200 vs 401)
			if w.Code != tc.expectedStatus {
				t.Errorf("Esperábamos el código %d, pero obtuvimos %d", tc.expectedStatus, w.Code)
			}
		})
	}
}

func TestRegisterHandler(t *testing.T) {
	// Limpiamos y recreamos la base de datos antes de este test
	setupTestDB()
	gin.SetMode(gin.TestMode)

	// Metemos un usuario previo para probar que bloquea duplicados
	database.DB.Create(&models.User{
		Username: "usuario_existente",
		Email:    "duplicado@test.com",
	})

	tests := []struct {
		name           string
		payload        RegisterRequest
		expectedStatus int
	}{
		{
			name:           "Caso 1: Registro exitoso",
			payload:        RegisterRequest{Username: "nuevo_user", Email: "nuevo@test.com", Password: "1234"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Caso 2: Correo ya registrado",
			payload:        RegisterRequest{Username: "otro_user", Email: "duplicado@test.com", Password: "1234"},
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Caso 3: Nombre de usuario ya registrado",
			payload:        RegisterRequest{Username: "usuario_existente", Email: "otro_correo@test.com", Password: "1234"},
			expectedStatus: http.StatusConflict,
		},
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
				t.Errorf("Esperábamos el código %d, pero obtuvimos %d", tc.expectedStatus, w.Code)
			}
		})
	}
}
