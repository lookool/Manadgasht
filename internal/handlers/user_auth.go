package handlers

import (
	"errors"
	"net/http"

	"manadgasht/internal/middleware"
	"manadgasht/internal/service"
)

// UserRegisterForm فرم ثبت‌نام کاربر عادی را نمایش می‌دهد و ثبت می‌کند.
func (h *Handlers) UserRegisterForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		render(w, r, "user_register.html", map[string]interface{}{"Title": "ثبت‌نام"})
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "درخواست نامعتبر است", http.StatusBadRequest)
		return
	}

	input := service.RegisterUserInput{
		FullName:     r.FormValue("full_name"),
		Email:        r.FormValue("email"),
		Phone:        r.FormValue("phone"),
		Password:     r.FormValue("password"),
		ReferralCode: r.FormValue("referral_code"),
	}

	user, err := h.Auth.RegisterUser(input)
	if err != nil {
		var verr *service.ValidationError
		if errors.As(err, &verr) {
			render(w, r, "user_register.html", map[string]interface{}{
				"Title":  "ثبت‌نام",
				"Errors": verr.Fields,
				"Form":   input,
			})
			return
		}
		logError("ثبت‌نام کاربر", err)
		http.Error(w, "خطای داخلی سرور", http.StatusInternalServerError)
		return
	}

	token, expiresAt, err := h.Auth.LoginUser(user.Email, input.Password)
	if err != nil {
		// ثبت‌نام موفق بود ولی ورود خودکار ناموفق شد؛ کاربر رو به صفحه‌ی ورود می‌فرستیم.
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	setSessionCookie(w, middleware.UserCookieName, token, expiresAt)
	http.Redirect(w, r, "/account", http.StatusSeeOther)
}

// UserLoginForm فرم ورود کاربر عادی را نمایش می‌دهد و پردازش می‌کند.
func (h *Handlers) UserLoginForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		render(w, r, "user_login.html", map[string]interface{}{"Title": "ورود"})
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "درخواست نامعتبر است", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	token, expiresAt, err := h.Auth.LoginUser(email, password)
	if err != nil {
		render(w, r, "user_login.html", map[string]interface{}{
			"Title": "ورود",
			"Error": "ایمیل یا رمز عبور اشتباه است",
			"Email": email,
		})
		return
	}

	setSessionCookie(w, middleware.UserCookieName, token, expiresAt)
	http.Redirect(w, r, "/account", http.StatusSeeOther)
}

// UserLogout نشست کاربر را باطل می‌کند.
func (h *Handlers) UserLogout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(middleware.UserCookieName); err == nil {
		_ = h.Auth.Logout(cookie.Value)
	}
	clearSessionCookie(w, middleware.UserCookieName)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Account پنل کاربری: پروفایل کوتاه + تاریخچه‌ی خریدها.
func (h *Handlers) Account(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)

	orders, err := h.Orders.ListByUser(user.ID)
	if err != nil {
		logError("بارگذاری سفارش‌های کاربر", err)
	}

	render(w, r, "account.html", map[string]interface{}{
		"Title":  "پنل کاربری",
		"User":   user,
		"Orders": orders,
	})
}
