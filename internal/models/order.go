package models

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order یک خرید ثبت‌شده در سایت است. اگر مشتری هنگام خرید کد معرف یک
// تورلیدر را وارد کرده باشد، AgentID و درصد/مبلغ کمیسیون محاسبه و
// ذخیره می‌شود.
type Order struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	PackageID uint

	// AgentID اگر خرید از طریق کد معرف یک تورلیدر انجام شده باشد پر می‌شود.
	AgentID          *uint
	ReferralCodeUsed string `gorm:"size:30"`

	TotalPrice           float64
	CommissionPercentage float64 // در لحظه خرید از CommissionSetting کپی می‌شود
	CommissionAmount     float64 // TotalPrice * CommissionPercentage / 100

	Status    OrderStatus `gorm:"size:20;default:pending"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
