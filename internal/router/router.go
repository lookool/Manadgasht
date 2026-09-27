package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"manadgasht/internal/handlers"
	appmw "manadgasht/internal/middleware"
)

func New(h *handlers.Handlers, auth *appmw.Auth) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Use(auth.WithUser)
	r.Use(auth.WithAgent)
	r.Use(auth.WithAdmin)

	fileServer := http.FileServer(http.Dir("web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	r.Get("/", h.Home)
	r.Get("/packages", h.PackagesList)
	r.Get("/packages/{slug}", h.PackageDetail)
	r.Get("/blog", h.BlogList)
	r.Get("/blog/{slug}", h.BlogDetail)
	r.Get("/about", h.About)
	r.Get("/contact", h.Contact)

	r.Get("/register", h.UserRegisterForm)
	r.Post("/register", h.UserRegisterForm)
	r.Get("/login", h.UserLoginForm)
	r.Post("/login", h.UserLoginForm)
	r.Post("/logout", h.UserLogout)
	r.Get("/account", auth.RequireUser(h.Account))
	r.Post("/packages/{slug}/buy", auth.RequireUser(h.BuyAsUser))

	r.Get("/agent/register", h.AgentRegisterForm)
	r.Post("/agent/register", h.AgentRegisterForm)
	r.Get("/agent/login", h.AgentLoginForm)
	r.Post("/agent/login", h.AgentLoginForm)
	r.Post("/agent/logout", h.AgentLogout)
	r.Get("/agent/dashboard", auth.RequireAgent(h.AgentDashboard))
	r.Post("/agent/withdrawals", auth.RequireAgent(h.AgentRequestWithdrawal))
	r.Post("/packages/{slug}/buy-wallet", auth.RequireAgent(h.BuyWithAgentWallet))

	r.Get("/admin/login", h.AdminLoginForm)
	r.Post("/admin/login", h.AdminLoginForm)
	r.Post("/admin/logout", h.AdminLogout)
	r.Get("/admin", auth.RequireAdmin(h.AdminDashboard))

	r.Get("/admin/packages", auth.RequireAdmin(h.AdminPackagesList))
	r.Get("/admin/packages/new", auth.RequireAdmin(h.AdminPackageNewForm))
	r.Post("/admin/packages", auth.RequireAdmin(h.AdminPackageCreate))
	r.Get("/admin/packages/{id}/edit", auth.RequireAdmin(h.AdminPackageEditForm))
	r.Post("/admin/packages/{id}/update", auth.RequireAdmin(h.AdminPackageUpdate))
	r.Post("/admin/packages/{id}/delete", auth.RequireAdmin(h.AdminPackageDelete))

	r.Get("/admin/blog", auth.RequireAdmin(h.AdminBlogList))
	r.Get("/admin/blog/new", auth.RequireAdmin(h.AdminBlogNewForm))
	r.Post("/admin/blog", auth.RequireAdmin(h.AdminBlogCreate))
	r.Get("/admin/blog/{id}/edit", auth.RequireAdmin(h.AdminBlogEditForm))
	r.Post("/admin/blog/{id}/update", auth.RequireAdmin(h.AdminBlogUpdate))
	r.Post("/admin/blog/{id}/delete", auth.RequireAdmin(h.AdminBlogDelete))

	return r
}