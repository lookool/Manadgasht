package models

import "time"

// User یک مشتری عادی سایت است (کسی که تور/هتل خرید می‌کند).
type User struct {
	ID           uint   `gorm:"primaryKey"`
	FullName     string `gorm:"size:150;not null"`
	Email        string `gorm:"size:150;uniqueIndex;not null"`
	Phone        string `gorm:"size:20;uniqueIndex"`
	PasswordHash string `gorm:"size:255;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
