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
	UserID             uint `gorm:"primaryKey;autoIncrement:false"` // PK, FK
	ProductID          uint `gorm:"primaryKey;autoIncrement:false"` // PK, FK
	NotifyPriceChanges bool
	NotifyTargetPrice  float64 `gorm:"type:decimal(10,2)"`
	TrackingStartDate  time.Time
}

// PriceHistory represent a change of price (for graphics)
type PriceHistory struct {
	PriceHistoryID uint    `gorm:"primaryKey;autoIncrement"` // PK
	ProductID      uint    // FK ProductID
	Price          float64 `gorm:"type:decimal(10,2)"`
	RegisterDate   time.Time
}
