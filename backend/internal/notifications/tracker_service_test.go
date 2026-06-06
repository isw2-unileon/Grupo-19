package handlers

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/database"
	"github.com/isw2-unileon/proyect-scaffolding/backend/internal/models"
	"gorm.io/gorm"
)

// setupTrackerTestDB creates a fresh in-memory database and migrates needed tables
func setupTrackerTestDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the test database")
	}

	err2 := db.AutoMigrate(&models.Product{}, &models.Tracking{}, &models.Notification{})
	if err2 != nil {
		panic("Failed to migrate the test database")
	}

	database.DB = db
}

func TestEvaluatePriceDropAndNotify_NoPriceDrop(t *testing.T) {
	setupTrackerTestDB()

	// If current price is >= old price, it should return early without errors
	// and without querying the database (meaning it won't fail even if the DB is empty).
	err := EvaluatePriceDropAndNotify(1, 100.0, 90.0)
	if err != nil {
		t.Errorf("Expected nil error when price goes up, got: %v", err)
	}
}

func TestEvaluatePriceDropAndNotify_ProductNotFound(t *testing.T) {
	setupTrackerTestDB()

	// Price dropped, but the product ID does not exist in the DB
	err := EvaluatePriceDropAndNotify(999, 40.0, 50.0)
	if err == nil {
		t.Errorf("Expected an error when the product is not found, got nil")
	}
}

func TestEvaluatePriceDropAndNotify_Logic(t *testing.T) {
	setupTrackerTestDB()

	// 1. Create a mock product
	product := models.Product{Name: "Gaming Laptop"}
	database.DB.Create(&product)

	// 2. Create mock trackings for different scenarios
	// User 1: Wants to be notified on ANY price drop
	database.DB.Create(&models.Tracking{
		UserID:             1,
		ProductID:          product.ProductID,
		NotifyPriceChanges: true,
	})

	// User 2: Wants to be notified only if it drops to 500 or below
	database.DB.Create(&models.Tracking{
		UserID:            2,
		ProductID:         product.ProductID,
		NotifyTargetPrice: 500.0,
	})

	// User 3: Wants to be notified only if it drops to 300 or below
	database.DB.Create(&models.Tracking{
		UserID:            3,
		ProductID:         product.ProductID,
		NotifyTargetPrice: 300.0,
	})

	// 3. Trigger a price drop: Old = 600.0, New = 450.0
	err := EvaluatePriceDropAndNotify(product.ProductID, 450.0, 600.0)
	if err != nil {
		t.Fatalf("Unexpected error executing the service: %v", err)
	}

	// 4. Verify the generated notifications
	var notifications []models.Notification
	database.DB.Find(&notifications)

	// We expect exactly 2 notifications:
	// - User 1 gets it because of general price drop
	// - User 2 gets it because 450.0 <= 500.0
	// - User 3 DOES NOT get it because 450.0 > 300.0
	if len(notifications) != 2 {
		t.Fatalf("Expected exactly 2 notifications to be generated, got %d", len(notifications))
	}

	// Helper to check what type of notification each user got
	userNotifs := make(map[uint]models.Notification)
	for _, n := range notifications {
		userNotifs[n.UserID] = n
	}

	// Check User 1
	n1, exists := userNotifs[1]
	if !exists {
		t.Errorf("User 1 should have received a notification")
	} else if n1.Type != "price_drop" {
		t.Errorf("Expected User 1 notification type to be 'price_drop', got '%s'", n1.Type)
	}

	// Check User 2
	n2, exists := userNotifs[2]
	if !exists {
		t.Errorf("User 2 should have received a notification")
	} else if n2.Type != "target_price" {
		t.Errorf("Expected User 2 notification type to be 'target_price', got '%s'", n2.Type)
	}

	// Check User 3
	_, exists = userNotifs[3]
	if exists {
		t.Errorf("User 3 SHOULD NOT have received a notification")
	}
}
