package middleware

import (
	"context"
	"net/http"

	"manadgasht/internal/models"
	"manadgasht/internal/service"
)

// Auth تمام میان‌افزارهای مربوط به احراز هویت را می‌سازد. نگه‌داشتن
// AuthService این‌جا باعث می‌شه هندلرها مستقیم درگیر خوندن کوکی نباشن.
type Auth struct {
	auth *service.AuthService
}

func NewAuth(auth *service.AuthService) *Auth {
	return &Auth{auth: auth}
}

// WithUser اگر کوکی نشستِ معتبر کاربر وجود داشته باشه، کاربر رو توی
// context درخواست می‌ذاره؛ در غیر این صورت بدون خطا رد میشه (برای
// صفحاتی که هم برای مهمون هم برای کاربر لاگین‌کرده باز هستن).
func (a *Auth) WithUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie(UserCookieName); err == nil {
			if user, err := a.auth.CurrentUser(cookie.Value); err == nil {
				ctx := context.WithValue(r.Context(), userContextKey, user)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}

// WithAgent مشابه WithUser ولی برای همکار.
func (a *Auth) WithAgent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie(AgentCookieName); err == nil {
			if agent, err := a.auth.CurrentAgent(cookie.Value); err == nil {
				ctx := context.WithValue(r.Context(), agentContextKey, agent)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}

// RequireUser فقط به کاربرِ واردشده اجازه‌ی عبور می‌ده؛ در غیر این
// صورت به صفحه‌ی ورود کاربر هدایت می‌کنه.
func (a *Auth) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if CurrentUser(r) == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

// RequireAgent فقط به همکارِ واردشده اجازه‌ی عبور می‌ده.
func (a *Auth) RequireAgent(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if CurrentAgent(r) == nil {
			http.Redirect(w, r, "/agent/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

// CurrentUser کاربر لاگین‌کرده‌ی داخل context رو برمی‌گردونه (یا nil).
func CurrentUser(r *http.Request) *models.User {
	user, _ := r.Context().Value(userContextKey).(*models.User)
	return user
}

// CurrentAgent همکارِ لاگین‌کرده‌ی داخل context رو برمی‌گردونه (یا nil).
func CurrentAgent(r *http.Request) *models.Agent {
	agent, _ := r.Context().Value(agentContextKey).(*models.Agent)
	return agent
}
