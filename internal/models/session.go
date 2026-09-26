package models

import "time"

// SubjectType مشخص می‌کند یک نشست (Session) متعلق به چه نوع حسابی است.
type SubjectType string

const (
	SubjectUser  SubjectType = "user"
	SubjectAgent SubjectType = "agent"
	SubjectAdmin SubjectType = "admin"
)

// Session یک نشست ورود فعال است. توکن آن داخل کوکی مرورگر کاربر ذخیره
// می‌شود؛ خودِ داده‌های نشست همیشه سمت سرور (دیتابیس) می‌ماند.
type Session struct {
	Token       string      `gorm:"primaryKey;size:64"`
	SubjectType SubjectType `gorm:"size:20;not null"`
	SubjectID   uint        `gorm:"not null"`
	ExpiresAt   time.Time   `gorm:"not null"`
	CreatedAt   time.Time
}
