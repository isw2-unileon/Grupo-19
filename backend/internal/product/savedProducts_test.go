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
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// setupSavedProductsTestDB creates a fresh in-memory database for testing
func setupSavedProductsTestDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the test database")
	}

	err2 := db.AutoMigrate(&models.User{}, &models.Product{})
	if err2 != nil {
		panic("Failed to migrate the test database")
	}

	database.DB = db
}

func TestGetSavedProducts(t *testing.T) {
	setupSavedProductsTestDB()
	gin.SetMode(gin.TestMode)

	// 1. Create mock products in the database
	prod1 := models.Product{Name: "Monitor UltraWide"}
	prod2 := models.Product{Name: "Teclado Mecánico"}
	database.DB.Create(&prod1)
	database.DB.Create(&prod2)

	// 2. Create mock users
	// User with 2 saved products
	userWithItems := models.User{
		Username:      "user_with_items",
		SavedProducts: pq.Int64Array{int64(prod1.ProductID), int64(prod2.ProductID)},
	}
	database.DB.Create(&userWithItems)

	// User with an empty saved products array
	userWithoutItems := models.User{
		Username:      "user_empty",
		SavedProducts: pq.Int64Array{},
	}
	database.DB.Create(&userWithoutItems)

	tests := []struct {
		name           string
		authUserID     uint
		useMiddleware  bool
		expectedStatus int
		expectedLen    int
	}{
		{"Successfully retrieve saved products", userWithItems.UserID, true, http.StatusOK, 2},
		{"User has no saved products (returns empty list)", userWithoutItems.UserID, true, http.StatusOK, 0},
		{"User not found in database", 999, true, http.StatusNotFound, 0},
		{"Unauthorized request (no context injected)", userWithItems.UserID, false, http.StatusUnauthorized, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.Default()

			// Apply the mock middleware only if the test case requires it
			if tc.useMiddleware {
				r.Use(mockAuthMiddleware(tc.authUserID))
			}
			r.GET("/api/user/saved-products", GetSavedProducts)

			req, _ := http.NewRequest("GET", "/api/user/saved-products", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			// If we expect a 200 OK, verify the payload length
			if tc.expectedStatus == http.StatusOK {
				// The response is wrapped in a "data" object: {"data": [...]}
				var response struct {
					Data []models.Product `json:"data"`
				}
				_ = json.Unmarshal(w.Body.Bytes(), &response)

				if len(response.Data) != tc.expectedLen {
					t.Errorf("Expected %d products, got %d", tc.expectedLen, len(response.Data))
				}
			}
		})
	}
}
