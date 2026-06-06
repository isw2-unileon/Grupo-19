package handlers

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
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// setupTrackingTestDB creates a fresh in-memory database and migrates needed tables
func setupTrackingTestDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the test database")
	}

	err2 := db.AutoMigrate(&models.User{}, &models.Product{}, &models.Tracking{})
	if err2 != nil {
		panic("Failed to migrate the test database")
	}

	database.DB = db
}

// mockAuthMiddleware simulates a valid JWT session by injecting a userID
func mockAuthMiddleware(userID uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	}
}

func TestUpdateTracking(t *testing.T) {
	setupTrackingTestDB()
	gin.SetMode(gin.TestMode)

	// Create a test user and product
	user := models.User{Username: "tracker_user", SavedProducts: pq.Int64Array{}}
	database.DB.Create(&user)

	product := models.Product{Name: "PlayStation 5"}
	database.DB.Create(&product)

	t.Run("Create a new tracking record", func(t *testing.T) {
		r := gin.Default()
		r.Use(mockAuthMiddleware(user.UserID))
		r.POST("/api/tracking", UpdateTracking)

		payload := TrackingRequest{ProductID: product.ProductID, TargetPrice: 450.0}
		jsonValue, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/tracking", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Verify the tracking was created in the DB
		var tracking models.Tracking
		if err := database.DB.First(&tracking).Error; err != nil {
			t.Errorf("Tracking record was not created in the database")
		} else if tracking.NotifyTargetPrice != 450.0 {
			t.Errorf("Expected target price 450.0, got %f", tracking.NotifyTargetPrice)
		}

		// Verify the product was added to the user's SavedProducts
		var updatedUser models.User
		database.DB.First(&updatedUser, user.UserID)
		if len(updatedUser.SavedProducts) != 1 || updatedUser.SavedProducts[0] != int64(product.ProductID) {
			t.Errorf("Product was not properly added to User's SavedProducts array")
		}
	})

	t.Run("Update an existing tracking record", func(t *testing.T) {
		r := gin.Default()
		r.Use(mockAuthMiddleware(user.UserID))
		r.POST("/api/tracking", UpdateTracking)

		// Change the target price to 400.0 for the same product
		payload := TrackingRequest{ProductID: product.ProductID, TargetPrice: 400.0}
		jsonValue, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/tracking", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Verify it was updated, not duplicated
		var trackings []models.Tracking
		database.DB.Find(&trackings)
		if len(trackings) != 1 {
			t.Errorf("Expected 1 tracking record, got %d (might have duplicated)", len(trackings))
		} else if trackings[0].NotifyTargetPrice != 400.0 {
			t.Errorf("Expected target price to be updated to 400.0, got %f", trackings[0].NotifyTargetPrice)
		}
	})
}

func TestCheckTracking(t *testing.T) {
	setupTrackingTestDB()
	gin.SetMode(gin.TestMode)

	user := models.User{Username: "check_user"}
	database.DB.Create(&user)

	product1 := models.Product{Name: "Tracked Item"}
	product2 := models.Product{Name: "Untracked Item"}
	database.DB.Create(&product1)
	database.DB.Create(&product2)

	// User follows product 1
	database.DB.Create(&models.Tracking{UserID: user.UserID, ProductID: product1.ProductID, NotifyTargetPrice: 19.99})

	t.Run("Product is being tracked", func(t *testing.T) {
		r := gin.Default()
		r.Use(mockAuthMiddleware(user.UserID))
		r.GET("/api/tracking/:id", CheckTracking)

		req, _ := http.NewRequest("GET", "/api/tracking/1", nil) // Product 1
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &response)

		if isFollowing, ok := response["is_following"].(bool); !ok || !isFollowing {
			t.Errorf("Expected is_following to be true")
		}
		if targetPrice, ok := response["target_price"].(float64); !ok || targetPrice != 19.99 {
			t.Errorf("Expected target_price to be 19.99")
		}
	})

	t.Run("Product is NOT being tracked", func(t *testing.T) {
		r := gin.Default()
		r.Use(mockAuthMiddleware(user.UserID))
		r.GET("/api/tracking/:id", CheckTracking)

		req, _ := http.NewRequest("GET", "/api/tracking/2", nil) // Product 2
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &response)

		if isFollowing, ok := response["is_following"].(bool); !ok || isFollowing {
			t.Errorf("Expected is_following to be false")
		}
	})
}

func TestUnfollowProduct(t *testing.T) {
	setupTrackingTestDB()
	gin.SetMode(gin.TestMode)

	// Create user with pre-populated SavedProducts
	user := models.User{Username: "unfollow_user", SavedProducts: pq.Int64Array{1, 2}}
	database.DB.Create(&user)

	// Create tracking records
	database.DB.Create(&models.Tracking{UserID: user.UserID, ProductID: 1})
	database.DB.Create(&models.Tracking{UserID: user.UserID, ProductID: 2})

	r := gin.Default()
	r.Use(mockAuthMiddleware(user.UserID))
	r.DELETE("/api/tracking/:id", UnfollowProduct)

	// Action: Unfollow product 1
	req, _ := http.NewRequest("DELETE", "/api/tracking/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify tracking record is deleted
	var trackings []models.Tracking
	database.DB.Find(&trackings)
	if len(trackings) != 1 || trackings[0].ProductID != 2 {
		t.Errorf("Expected only tracking for product 2 to remain")
	}

	// Verify user's SavedProducts array is updated
	var updatedUser models.User
	database.DB.First(&updatedUser, user.UserID)
	if len(updatedUser.SavedProducts) != 1 || updatedUser.SavedProducts[0] != 2 {
		t.Errorf("Expected SavedProducts to only contain ID 2, got %v", updatedUser.SavedProducts)
	}
}

func TestUpdateUserSavedProducts(t *testing.T) {
	setupTrackingTestDB()

	user := models.User{Username: "helper_user", SavedProducts: pq.Int64Array{}}
	database.DB.Create(&user)

	// 1. Add product 10
	_ = updateUserSavedProducts(user.UserID, 10, true)
	database.DB.First(&user, user.UserID)
	if len(user.SavedProducts) != 1 || user.SavedProducts[0] != 10 {
		t.Errorf("Expected product 10 to be added")
	}

	// 2. Add product 10 again (should not duplicate)
	_ = updateUserSavedProducts(user.UserID, 10, true)
	database.DB.First(&user, user.UserID)
	if len(user.SavedProducts) != 1 {
		t.Errorf("Expected no duplicates, length should still be 1")
	}

	// 3. Remove product 10
	_ = updateUserSavedProducts(user.UserID, 10, false)
	database.DB.First(&user, user.UserID)
	if len(user.SavedProducts) != 0 {
		t.Errorf("Expected array to be empty after removal")
	}
}
