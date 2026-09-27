package handlers

import (
	"net/http"

	"manadgasht/internal/middleware"
)

func (h *Handlers) AdminLoginForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderAdmin(w, r, "admin_login.html", map[string]interface{}{"Title": "ورود ادمین"})
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "درخواست نامعتبر است", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	token, expiresAt, err := h.Auth.LoginAdmin(username, password)
	if err != nil {
		renderAdmin(w, r, "admin_login.html", map[string]interface{}{
			"Title":    "ورود ادمین",
			"Error":    "نام‌کاربری یا رمز عبور اشتباه است",
			"Username": username,
		})
		return
	}

	setSessionCookie(w, middleware.AdminCookieName, token, expiresAt)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *Handlers) AdminLogout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(middleware.AdminCookieName); err == nil {
		_ = h.Auth.Logout(cookie.Value)
	}
	clearSessionCookie(w, middleware.AdminCookieName)
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

func (h *Handlers) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	packages, _ := h.Packages.ListAll()
	posts, _ := h.Blog.ListAll()

	renderAdmin(w, r, "admin_dashboard.html", map[string]interface{}{
		"Title":        "پنل ادمین",
		"PackageCount": len(packages),
		"PostCount":    len(posts),
	})
}