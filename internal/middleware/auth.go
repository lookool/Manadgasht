package middleware

import (
	"context"
	"net/http"

	"manadgasht/internal/models"
	"manadgasht/internal/service"
)

type Auth struct {
	auth *service.AuthService
}

func NewAuth(auth *service.AuthService) *Auth {
	return &Auth{auth: auth}
}

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

func (a *Auth) WithAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie(AdminCookieName); err == nil {
			if admin, err := a.auth.CurrentAdmin(cookie.Value); err == nil {
				ctx := context.WithValue(r.Context(), adminContextKey, admin)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (a *Auth) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if CurrentUser(r) == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func (a *Auth) RequireAgent(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if CurrentAgent(r) == nil {
			http.Redirect(w, r, "/agent/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func (a *Auth) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if CurrentAdmin(r) == nil {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func CurrentUser(r *http.Request) *models.User {
	user, _ := r.Context().Value(userContextKey).(*models.User)
	return user
}

func CurrentAgent(r *http.Request) *models.Agent {
	agent, _ := r.Context().Value(agentContextKey).(*models.Agent)
	return agent
}

func CurrentAdmin(r *http.Request) *models.Admin {
	admin, _ := r.Context().Value(adminContextKey).(*models.Admin)
	return admin
}