package models

import "time"

// User respresent each user of the app
type User struct {
	UserID     uint `gorm:"primaryKey;autoIncrement"` // PK
	Username   string
	Email      string
	Password   string
	UserType   string
	RegisterAt time.Time
}

// Product represent the products the scapper have benn asked to scrap
type Product struct {
	ProductID   uint `gorm:"primaryKey;autoIncrement"` // PK
	Name        string
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
	SourceURL   string
	LastPrice   float64 `gorm:"type:decimal(10,2)"`
	LowestPrice float64 `gorm:"type:decimal(10,2)"`
	CreatedBy   uint    // FK UserID
	CreateAt    time.Time
	UpdatedAt   time.Time
}

// Tracking represent the User intention of being notificated when price drop
type Tracking struct {
	UserID             uint    `gorm:"primaryKey;autoIncrement:false" json:"user_id"`
	ProductID          uint    `gorm:"primaryKey;autoIncrement:false" json:"product_id"`
	NotifyPriceChanges bool    `json:"notify_price_changes"`
	NotifyTargetPrice  float64 `gorm:"type:decimal(10,2)" json:"target_price"`
	TrackingStartDate  time.Time
}

// PriceHistory represent a change of price (for graphics)
type PriceHistory struct {
	PriceHistoryID uint    `gorm:"primaryKey;autoIncrement"` // PK
	ProductID      uint    // FK ProductID
	Price          float64 `gorm:"type:decimal(10,2)"`
	RegisterDate   time.Time
}

type Notification struct {
	NotificationID uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint      `gorm:"not null;index" json:"userId"`
	ProductID      uint      `gorm:"not null" json:"productId"`
	Type           string    `gorm:"type:varchar(20);default:'price_drop'" json:"type"`
	Title          string    `gorm:"type:varchar(100);not null" json:"title"`
	Description    string    `gorm:"type:text;not null" json:"description"`
	IsRead         bool      `gorm:"default:false" json:"isRead"`
	CreatedAt      time.Time `json:"time"`
}
