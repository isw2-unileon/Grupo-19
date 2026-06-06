package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
	"gorm.io/gorm"
)

// setupSearchTestDB creates a fresh in-memory database and migrates needed tables
func setupSearchTestDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the test database")
	}

	err2 := db.AutoMigrate(&models.Product{}, &models.PriceHistory{})
	if err2 != nil {
		panic("Failed to migrate the test database")
	}

	database.DB = db
}

func TestGetProductByID(t *testing.T) {
	setupSearchTestDB()
	gin.SetMode(gin.TestMode)

	// Create a mock product in the database
	testProduct := models.Product{Name: "Monitor 4K", SourceURL: "http://example.com/monitor", LastPrice: 299.99}
	database.DB.Create(&testProduct)

	// Convert the primary key to string for the URL route
	// Note: In GORM, the first created element usually gets ID 1
	productIDStr := "1"

	tests := []struct {
		name           string
		paramID        string
		expectedStatus int
	}{
		{"Valid product ID", productIDStr, http.StatusOK},
		{"Invalid ID format (string)", "abc", http.StatusBadRequest},
		{"Non-existent product ID", "999", http.StatusNotFound},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.Default()
			r.GET("/api/products/:id", GetProductByID)

			req, _ := http.NewRequest("GET", "/api/products/"+tc.paramID, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}

func TestSearchProducts(t *testing.T) {
	setupSearchTestDB()
	gin.SetMode(gin.TestMode)

	// Insert mock products for the search engine to find
	database.DB.Create(&models.Product{Name: "Teclado Mecanico", SourceURL: "http://fake-url.com/1", LastPrice: 50.0})
	database.DB.Create(&models.Product{Name: "Raton Gaming", SourceURL: "http://fake-url.com/2", LastPrice: 30.0})

	tests := []struct {
		name           string
		queryParam     string
		expectedStatus int
		expectedLen    int // Expected number of products returned
	}{
		{
			name:           "Empty query term",
			queryParam:     "",
			expectedStatus: http.StatusBadRequest,
			expectedLen:    0,
		},
		{
			name:           "Query with no matches",
			queryParam:     "Impresora",
			expectedStatus: http.StatusOK,
			expectedLen:    0,
		},
		{
			name: "Query with a valid match",
			// Note: "ILIKE %Teclado%" should find "Teclado Mecanico"
			queryParam:     "Teclado",
			expectedStatus: http.StatusOK,
			expectedLen:    1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.Default()
			r.GET("/api/products/search", SearchProducts)

			req, _ := http.NewRequest("GET", "/api/products/search?q="+tc.queryParam, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			// If we expect a 200 OK, we parse the JSON array to check the results
			if tc.expectedStatus == http.StatusOK {
				var response []models.Product
				_ = json.Unmarshal(w.Body.Bytes(), &response)

				if len(response) != tc.expectedLen {
					t.Errorf("Expected %d products in the response, got %d", tc.expectedLen, len(response))
				}
			}
		})
	}
}
