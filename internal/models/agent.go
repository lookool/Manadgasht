package models

import "time"

// Agent همان «تورلیدر» است: یک همکار فروش که یک کد معرف مخصوص به خودش
// دارد و از هر خریدی که با آن کد ثبت شود، درصدی کمیسیون می‌گیرد.
type Agent struct {
	ID           uint   `gorm:"primaryKey"`
	FullName     string `gorm:"size:150;not null"`
	NationalID   string `gorm:"size:15;uniqueIndex;not null"`
	Email        string `gorm:"size:150;uniqueIndex;not null"`
	Phone        string `gorm:"size:20;uniqueIndex;not null"`
	PasswordHash string `gorm:"size:255;not null"`

	// CardNumber شماره کارتی است که همکار موقع ثبت‌نام می‌دهد (برای
	// مرجع اولیه)؛ در هر درخواست برداشت هم می‌تواند شماره‌ی جدید بدهد.
	CardNumber string `gorm:"size:30"`

	// ReferralCode کدی است که مشتری هنگام خرید وارد می‌کند تا خرید
	// به این تورلیدر نسبت داده شود.
	ReferralCode string `gorm:"size:30;uniqueIndex;not null"`

	// Balance موجودی کیف‌پول همکار است: مجموع کمیسیون‌های واریزشده
	// منهای برداشت‌های انجام‌شده. با هر خرید موفقی که از کد معرف این
	// همکار استفاده شود، به‌صورت خودکار افزایش پیدا می‌کند.
	Balance float64 `gorm:"not null;default:0"`

	IsActive  bool `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
