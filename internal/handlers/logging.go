package handlers

import "log/slog"

// logError یک خطای غیربحرانی (مثلاً شکست یک کوئری برای یک بخش غیرضروری
// از صفحه) را با فرمت یکدست لاگ می‌کند، بدون اینکه درخواست کاربر رو قطع کنه.
func logError(action string, err error) {
	slog.Error("خطا در "+action, "error", err)
}
