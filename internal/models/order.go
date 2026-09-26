package models

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order یک خرید ثبت‌شده در سایت است. دقیقاً یکی از دو فیلد UserID یا
// PurchasedByAgentID پر می‌شود:
//   - UserID: یک مشتری عادی خریده (احتمالاً با کد معرف یک همکار)
//   - PurchasedByAgentID: خودِ همکار با موجودی کیف‌پولش خریده
//
// ReferringAgentID فقط وقتی پر می‌شود که خریدار یک مشتری عادی باشد و
// موقع ثبت‌نام یک کد معرف معتبر وارد کرده باشد.
type Order struct {
	ID uint `gorm:"primaryKey"`

	UserID             *uint
	PurchasedByAgentID *uint

	PackageID uint `gorm:"not null"`

	ReferringAgentID *uint
	ReferralCodeUsed string `gorm:"size:30"`

	TotalPrice           float64
	CommissionPercentage float64 // در لحظه خرید از CommissionSetting کپی می‌شود
	CommissionAmount     float64 // TotalPrice * CommissionPercentage / 100

	Status    OrderStatus `gorm:"size:20;default:confirmed"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
