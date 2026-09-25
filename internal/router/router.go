package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"manadgasht/internal/handlers"
)

// New تمام مسیرهای HTTP سایت را تعریف می‌کند. main.go فقط این تابع رو
// صدا می‌زنه؛ خودِ منطق مسیرها اینجا جدا نگه داشته می‌شه.
func New(h *handlers.Handlers) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

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

	return r
}
