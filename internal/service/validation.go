package service

import "strings"

// ValidationError خطای اعتبارسنجی فرم است؛ Fields کلیدش اسم فیلد و
// مقدارش پیام خطای همان فیلد است (برای نمایش زیر همان input در فرم).
type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	var parts []string
	for field, msg := range e.Fields {
		parts = append(parts, field+": "+msg)
	}
	return strings.Join(parts, "، ")
}

func newValidationError() *ValidationError {
	return &ValidationError{Fields: map[string]string{}}
}

func (e *ValidationError) add(field, msg string) {
	e.Fields[field] = msg
}

func (e *ValidationError) hasErrors() bool {
	return len(e.Fields) > 0
}
