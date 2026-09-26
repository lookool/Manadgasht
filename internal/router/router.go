package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"manadgasht/internal/handlers"
	appmw "manadgasht/internal/middleware"
)

// New تمام مسیرهای HTTP سایت را تعریف می‌کند. main.go فقط این تابع رو
// صدا می‌زنه؛ خودِ منطق مسیرها اینجا جدا نگه داشته می‌شه.
func New(h *handlers.Handlers, auth *appmw.Auth) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// روی همه‌ی مسیرها اجرا می‌شه: اگه کوکی نشست معتبر باشه، کاربر/همکار
	// رو تو context می‌ذاره؛ اگه نباشه، بدون خطا رد می‌شه (صفحات عمومی
	// هم به این نیاز دارن تا هدر بدونه دکمه‌ی «ورود» نشون بده یا «پنل»).
	r.Use(auth.WithUser)
	r.Use(auth.WithAgent)

	// فایل‌های استاتیک (عکس، CSS)
	fileServer := http.FileServer(http.Dir("web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// صفحات عمومی سایت
	r.Get("/", h.Home)
	r.Get("/packages", h.PackagesList)
	r.Get("/packages/{slug}", h.PackageDetail)
	r.Get("/blog", h.BlogList)
	r.Get("/blog/{slug}", h.BlogDetail)
	r.Get("/about", h.About)
	r.Get("/contact", h.Contact)

	// احراز هویت کاربر عادی
	r.Get("/register", h.UserRegisterForm)
	r.Post("/register", h.UserRegisterForm)
	r.Get("/login", h.UserLoginForm)
	r.Post("/login", h.UserLoginForm)
	r.Post("/logout", h.UserLogout)
	r.Get("/account", auth.RequireUser(h.Account))
	r.Post("/packages/{slug}/buy", auth.RequireUser(h.BuyAsUser))

	// احراز هویت همکار (تورلیدر)
	r.Get("/agent/register", h.AgentRegisterForm)
	r.Post("/agent/register", h.AgentRegisterForm)
	r.Get("/agent/login", h.AgentLoginForm)
	r.Post("/agent/login", h.AgentLoginForm)
	r.Post("/agent/logout", h.AgentLogout)
	r.Get("/agent/dashboard", auth.RequireAgent(h.AgentDashboard))
	r.Post("/agent/withdrawals", auth.RequireAgent(h.AgentRequestWithdrawal))
	r.Post("/packages/{slug}/buy-wallet", auth.RequireAgent(h.BuyWithAgentWallet))

	return r
}
