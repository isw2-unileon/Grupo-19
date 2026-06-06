package handlers

import (
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

// setupTestDB creates a blank in-memory database for testing
func setupTestDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the test database")
	}

	// Migrate the required tables
	err2 := db.AutoMigrate(&models.Notification{})
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

func TestGetUserNotifications(t *testing.T) {
	setupTestDB()
	gin.SetMode(gin.TestMode)

	userID := uint(1)
	otherUserID := uint(2)

	// Create test notifications
	// Note: Adjust the field names (like NotificationID or IsRead) if your model defines them slightly differently.
	database.DB.Create(&models.Notification{UserID: userID, Title: "Read Notif", IsRead: true, CreatedAt: time.Now().Add(-1 * time.Hour)})
	database.DB.Create(&models.Notification{UserID: userID, Title: "Unread Notif", IsRead: false, CreatedAt: time.Now()})
	database.DB.Create(&models.Notification{UserID: otherUserID, Title: "Someone Else's Notif", IsRead: false})

	t.Run("Successfully fetch user notifications", func(t *testing.T) {
		r := gin.Default()
		r.Use(mockAuthMiddleware(userID))
		r.GET("/api/notifications", GetUserNotifications)

		req, _ := http.NewRequest("GET", "/api/notifications", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response []models.Notification
		_ = json.Unmarshal(w.Body.Bytes(), &response)

		// Check if it only retrieved the 2 notifications belonging to userID = 1
		if len(response) != 2 {
			t.Errorf("Expected 2 notifications, got %d", len(response))
		}

		// Check the sorting logic (IsRead == false should come first)
		if response[0].Title != "Unread Notif" {
			t.Errorf("Expected unread notification to be first, but got: %s", response[0].Title)
		}
	})

	t.Run("Unauthorized request", func(t *testing.T) {
		r := gin.Default()
		r.GET("/api/notifications", GetUserNotifications) // No middleware applied

		req, _ := http.NewRequest("GET", "/api/notifications", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})
}

func TestMarkNotificationAsRead(t *testing.T) {
	setupTestDB()
	gin.SetMode(gin.TestMode)

	userID := uint(1)
	notif := models.Notification{UserID: userID, Title: "Update Me", IsRead: false}
	database.DB.Create(&notif)

	// In GORM, after Create(), the struct is populated with the generated primary key.
	// We will format it as a string to pass it into the URL.
	// Example: "1"
	notifIDStr := "1"

	tests := []struct {
		name           string
		paramID        string
		authUserID     uint
		useMiddleware  bool
		expectedStatus int
	}{
		{"Success marking as read", notifIDStr, userID, true, http.StatusOK},
		{"Invalid notification ID type", "abc", userID, true, http.StatusBadRequest},
		{"Notification not found or belongs to another user", "999", userID, true, http.StatusNotFound},
		{"Unauthorized request", notifIDStr, userID, false, http.StatusUnauthorized},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.Default()
			if tc.useMiddleware {
				r.Use(mockAuthMiddleware(tc.authUserID))
			}
			r.PATCH("/api/notifications/:id", MarkNotificationAsRead)

			req, _ := http.NewRequest("PATCH", "/api/notifications/"+tc.paramID, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}

func TestDeleteNotification(t *testing.T) {
	setupTestDB()
	gin.SetMode(gin.TestMode)

	userID := uint(1)
	notif := models.Notification{UserID: userID, Title: "Delete Me"}
	database.DB.Create(&notif)

	notifIDStr := "1"

	tests := []struct {
		name           string
		paramID        string
		authUserID     uint
		useMiddleware  bool
		expectedStatus int
	}{
		{"Success deleting notification", notifIDStr, userID, true, http.StatusOK},
		{"Invalid notification ID type", "abc", userID, true, http.StatusBadRequest},
		{"Notification not found or belongs to another user", "999", userID, true, http.StatusNotFound},
		{"Unauthorized request", notifIDStr, userID, false, http.StatusUnauthorized},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.Default()
			if tc.useMiddleware {
				r.Use(mockAuthMiddleware(tc.authUserID))
			}
			r.DELETE("/api/notifications/:id", DeleteNotification)

			req, _ := http.NewRequest("DELETE", "/api/notifications/"+tc.paramID, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}
