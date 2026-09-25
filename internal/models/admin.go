package models

import "time"

// Admin مدیر سایت است که از پنل مدیریت، تورها، بلاگ، تورلیدرها و
// درصد کمیسیون‌ها را کنترل می‌کند.
type Admin struct {
	ID           uint   `gorm:"primaryKey"`
	FullName     string `gorm:"size:150;not null"`
	Username     string `gorm:"size:100;uniqueIndex;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
