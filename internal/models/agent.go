package models

import "time"

// Agent همان «تورلیدر» است: یک همکار فروش که یک کد معرف مخصوص به خودش
// دارد و از هر خریدی که با آن کد ثبت شود، درصدی کمیسیون می‌گیرد.
type Agent struct {
	ID           uint   `gorm:"primaryKey"`
	FullName     string `gorm:"size:150;not null"`
	Email        string `gorm:"size:150;uniqueIndex;not null"`
	Phone        string `gorm:"size:20;uniqueIndex"`
	PasswordHash string `gorm:"size:255;not null"`

	// ReferralCode کدی است که مشتری هنگام خرید وارد می‌کند تا خرید
	// به این تورلیدر نسبت داده شود.
	ReferralCode string `gorm:"size:30;uniqueIndex;not null"`

	IsActive  bool `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
