package repository

import "errors"

// ErrInsufficientBalance وقتی برگردانده می‌شود که تلاش برای کسر مبلغی
// بیشتر از موجودی فعلی همکار انجام شود (چه برای خرید با کیف‌پول، چه
// برای درخواست برداشت).
var ErrInsufficientBalance = errors.New("موجودی کافی نیست")
