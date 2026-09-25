package models

import "time"

// CommissionSetting درصد کمیسیونی است که ادمین برای هر نوع محصول
// (تور / هتل / ترکیبی) تعیین می‌کند. این درصد روی قیمت نمایش‌داده‌شده
// به مشتری هیچ تاثیری ندارد؛ فقط تعیین می‌کند چه سهمی از سود به
// تورلیدر معرف تعلق می‌گیرد.
type CommissionSetting struct {
	ID         uint        `gorm:"primaryKey"`
	Type       PackageType `gorm:"size:20;uniqueIndex;not null"`
	Percentage float64     `gorm:"not null"` // مثلا 20 یعنی ۲۰٪
	UpdatedAt  time.Time
}
