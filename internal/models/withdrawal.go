package models

import "time"

type WithdrawalStatus string

const (
	WithdrawalPending  WithdrawalStatus = "pending"
	WithdrawalPaid     WithdrawalStatus = "paid"
	WithdrawalRejected WithdrawalStatus = "rejected"
)

// WithdrawalRequest درخواست برداشت موجودی یک همکار از کیف‌پولش است.
// به‌محض ثبت درخواست، مبلغ از Balance همکار کسر می‌شود تا امکان خرج
// دوباره‌ی همان مبلغ وجود نداشته باشد؛ اگر ادمین درخواست را رد کند
// (در فاز پنل ادمین)، مبلغ باید به موجودی برگردد.
type WithdrawalRequest struct {
	ID      uint `gorm:"primaryKey"`
	AgentID uint `gorm:"not null"`
	Amount  float64 `gorm:"not null"`

	CardNumber string `gorm:"size:30"`
	IBAN       string `gorm:"size:34"` // شماره شبا

	Status WithdrawalStatus `gorm:"size:20;default:pending"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
