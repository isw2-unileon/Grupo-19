package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
	"gorm.io/gorm"
)

// setupAddProductTestDB creates a fresh in-memory database and migrates needed tables
func setupAddProductTestDB() {
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

func TestAddProduct_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.POST("/api/track", AddProduct)

	// Case 1: Send invalid JSON (missing URL field)
	// We test this because the success path is tightly coupled to the external scraper
	reqBody := []byte(`{"wrong_field": "http://example.com"}`)
	req, _ := http.NewRequest("POST", "/api/track", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request when JSON is invalid, got %d", w.Code)
	}
}

func TestGetPriceHistory(t *testing.T) {
	setupAddProductTestDB()
	gin.SetMode(gin.TestMode)

	// 1. Create a mock product
	product := models.Product{Name: "Test Laptop"}
	database.DB.Create(&product)

	// 2. Create price history records at different times
	now := time.Now()

	// Record from today
	database.DB.Create(&models.PriceHistory{ProductID: product.ProductID, Price: 1000.0, RegisterDate: now})
	// Record from 5 days ago
	database.DB.Create(&models.PriceHistory{ProductID: product.ProductID, Price: 1100.0, RegisterDate: now.AddDate(0, 0, -5)})
	// Record from 10 days ago
	database.DB.Create(&models.PriceHistory{ProductID: product.ProductID, Price: 1200.0, RegisterDate: now.AddDate(0, 0, -10)})

	// In an empty DB, the first inserted product usually gets ID 1
	productIDStr := "1"

	t.Run("Default 7 days history", func(t *testing.T) {
		r := gin.Default()
		r.GET("/api/history/:id", GetPriceHistory)

		req, _ := http.NewRequest("GET", "/api/history/"+productIDStr, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200 OK, got %d", w.Code)
		}

		var history []models.PriceHistory
		_ = json.Unmarshal(w.Body.Bytes(), &history)

		// We expect 2 records (Today and 5 days ago).
		// The one from 10 days ago should be excluded.
		if len(history) != 2 {
			t.Errorf("Expected 2 history records, got %d", len(history))
		}
	})

	t.Run("Custom 15 days history parameter", func(t *testing.T) {
		r := gin.Default()
		r.GET("/api/history/:id", GetPriceHistory)

		// Passing 'days=15' via query param
		req, _ := http.NewRequest("GET", "/api/history/"+productIDStr+"?days=15", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200 OK, got %d", w.Code)
		}

		var history []models.PriceHistory
		_ = json.Unmarshal(w.Body.Bytes(), &history)

		// We expect all 3 records to be returned
		if len(history) != 3 {
			t.Errorf("Expected 3 history records, got %d", len(history))
		}
	})
}
