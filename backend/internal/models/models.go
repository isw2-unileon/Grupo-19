package models

import "time"

type User struct {
	UserID     uint `gorm:"primaryKey;autoIncrement"` // PK
	Username   string
	Email      string
	Password   string
	UserType   string
	RegisterAt time.Time
}

type Product struct {
	ProductID   uint `gorm:"primaryKey;autoIncrement"` // PK
	Name        string
	SourceURL   string
	LastPrice   float64 `gorm:"type:decimal(10,2)"`
	LowestPrice float64 `gorm:"type:decimal(10,2)"`
	CreatedBy   uint    // FK UserID
	CreateAt    time.Time
	UpdatedAt   time.Time
}

type Tracking struct {
	UserID             uint `gorm:"primaryKey;autoIncrement:false"` // PK, FK
	ProductID          uint `gorm:"primaryKey;autoIncrement:false"` // PK, FK
	NotifyPriceChanges bool
	NotifyTargetPrice  float64 `gorm:"type:decimal(10,2)"`
	TrackingStartDate  time.Time
}

type PriceHistory struct {
	PriceHistoryID uint    `gorm:"primaryKey;autoIncrement"` // PK
	ProductID      uint    // FK ProductID
	Price          float64 `gorm:"type:decimal(10,2)"`
	RegisterDate   time.Time
}
