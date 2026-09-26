package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"manadgasht/internal/middleware"
	"manadgasht/internal/repository"
)

// BuyAsUser خرید یک کاربر عادی را ثبت می‌کند. فعلاً به هیچ درگاه
// پرداختی وصل نیست؛ سفارش مستقیم «تأییدشده» ثبت می‌شود تا جریان
// کمیسیون از الان قابل تست باشد.
func (h *Handlers) BuyAsUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	slug := chi.URLParam(r, "slug")

	_, err := h.Orders.PlaceOrderForUser(user.ID, slug)
	if err != nil {
		logError("ثبت سفارش کاربر", err)
		http.Redirect(w, r, "/packages/"+slug+"?error=order", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/account", http.StatusSeeOther)
}

// BuyWithAgentWallet خرید همکار از موجودی کیف‌پول خودش را ثبت می‌کند.
func (h *Handlers) BuyWithAgentWallet(w http.ResponseWriter, r *http.Request) {
	agent := middleware.CurrentAgent(r)
	slug := chi.URLParam(r, "slug")

	_, err := h.Orders.PlaceOrderForAgentWallet(agent.ID, slug)
	if err != nil {
		if errors.Is(err, repository.ErrInsufficientBalance) {
			http.Redirect(w, r, "/packages/"+slug+"?error=balance", http.StatusSeeOther)
			return
		}
		logError("ثبت سفارش همکار از کیف‌پول", err)
		http.Redirect(w, r, "/packages/"+slug+"?error=order", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/agent/dashboard", http.StatusSeeOther)
}
