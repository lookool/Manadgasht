package models

import "time"

// PackageType نوع محصولی است که فروخته می‌شود.
type PackageType string

const (
	PackageTypeTour      PackageType = "tour"       // فقط تور
	PackageTypeHotel     PackageType = "hotel"      // فقط هتل
	PackageTypeTourHotel PackageType = "tour_hotel" // پکیج ترکیبی تور + هتل
)

// Package یک آیتم قابل‌فروش در سایت است: یک تور، یک هتل یا ترکیب هر دو.
type Package struct {
	ID          uint        `gorm:"primaryKey"`
	Title       string      `gorm:"size:200;not null"`
	Slug        string      `gorm:"size:220;uniqueIndex;not null"`
	Type        PackageType `gorm:"size:20;not null"`
	Description string      `gorm:"type:text"`
	Location    string      `gorm:"size:150"`

	// IsInternational مشخص می‌کند این تور/هتل در دسته‌ی «تورهای خارجی»
	// نمایش داده شود یا «تورهای داخلی» (دقیقاً مطابق دو منوی سایت فعلی).
	IsInternational bool `gorm:"default:false"`

	DurationDays int        `gorm:"default:0"`
	Price       float64     `gorm:"not null"`
	CoverImage  string      `gorm:"size:255"`

	// Gallery آدرس عکس‌های بیشتر را به صورت رشته جدا شده با کاما نگه می‌دارد.
	// در آینده می‌شود این را به یک جدول جدا (PackageImage) تبدیل کرد.
	Gallery string `gorm:"type:text"`

	IsActive  bool `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
