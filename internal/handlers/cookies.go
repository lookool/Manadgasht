package handlers

import (
	"net/http"
	"time"
)

func setSessionCookie(w http.ResponseWriter, name, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure عمداً تنظیم نشده چون توسعه‌ی محلی روی http:// انجام
		// میشه؛ موقع دیپلوی روی سرور واقعی (با HTTPS) حتماً اضافه بشه.
	})
}

func clearSessionCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
