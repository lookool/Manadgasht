package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"manadgasht/internal/middleware"
	"manadgasht/internal/service"
)

func (h *Handlers) AgentRegisterForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		render(w, r, "agent_register.html", map[string]interface{}{"Title": "ثبت‌نام همکاران"})
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "درخواست نامعتبر است", http.StatusBadRequest)
		return
	}

	input := service.RegisterAgentInput{
		FullName:   r.FormValue("full_name"),
		NationalID: r.FormValue("national_id"),
		Email:      r.FormValue("email"),
		Phone:      r.FormValue("phone"),
		CardNumber: r.FormValue("card_number"),
		Password:   r.FormValue("password"),
	}

	agent, err := h.Auth.RegisterAgent(input)
	if err != nil {
		var verr *service.ValidationError
		if errors.As(err, &verr) {
			render(w, r, "agent_register.html", map[string]interface{}{
				"Title":  "ثبت‌نام همکاران",
				"Errors": verr.Fields,
				"Form":   input,
			})
			return
		}
		logError("ثبت‌نام همکار", err)
		http.Error(w, "خطای داخلی سرور", http.StatusInternalServerError)
		return
	}

	token, expiresAt, err := h.Auth.LoginAgent(agent.Email, input.Password)
	if err != nil {
		http.Redirect(w, r, "/agent/login", http.StatusSeeOther)
		return
	}
	setSessionCookie(w, middleware.AgentCookieName, token, expiresAt)
	http.Redirect(w, r, "/agent/dashboard", http.StatusSeeOther)
}

func (h *Handlers) AgentLoginForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		render(w, r, "agent_login.html", map[string]interface{}{"Title": "ورود همکاران"})
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "درخواست نامعتبر است", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	token, expiresAt, err := h.Auth.LoginAgent(email, password)
	if err != nil {
		render(w, r, "agent_login.html", map[string]interface{}{
			"Title": "ورود همکاران",
			"Error": "ایمیل یا رمز عبور اشتباه است",
			"Email": email,
		})
		return
	}

	setSessionCookie(w, middleware.AgentCookieName, token, expiresAt)
	http.Redirect(w, r, "/agent/dashboard", http.StatusSeeOther)
}

func (h *Handlers) AgentLogout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(middleware.AgentCookieName); err == nil {
		_ = h.Auth.Logout(cookie.Value)
	}
	clearSessionCookie(w, middleware.AgentCookieName)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// AgentDashboard پنل همکار: کد معرف، موجودی، تاریخچه‌ی کمیسیون‌ها و
// درخواست‌های برداشت.
func (h *Handlers) AgentDashboard(w http.ResponseWriter, r *http.Request) {
	agent := middleware.CurrentAgent(r)

	orders, err := h.Orders.ListByReferringAgent(agent.ID)
	if err != nil {
		logError("بارگذاری کمیسیون‌های همکار", err)
	}

	withdrawals, err := h.Wallet.ListByAgent(agent.ID)
	if err != nil {
		logError("بارگذاری درخواست‌های برداشت", err)
	}

	render(w, r, "agent_dashboard.html", map[string]interface{}{
		"Title":       "پنل همکاران",
		"Agent":       agent,
		"Orders":      orders,
		"Withdrawals": withdrawals,
	})
}

// AgentRequestWithdrawal درخواست برداشت جدید ثبت می‌کند.
func (h *Handlers) AgentRequestWithdrawal(w http.ResponseWriter, r *http.Request) {
	agent := middleware.CurrentAgent(r)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "درخواست نامعتبر است", http.StatusBadRequest)
		return
	}

	amount, scanErr := strconv.ParseFloat(r.FormValue("amount"), 64)

	input := service.RequestWithdrawalInput{
		AgentID:    agent.ID,
		Amount:     amount,
		CardNumber: r.FormValue("card_number"),
		IBAN:       r.FormValue("iban"),
	}

	_, err := h.Wallet.RequestWithdrawal(input)
	if err != nil || scanErr != nil {
		errFields := map[string]string{}
		var verr *service.ValidationError
		if errors.As(err, &verr) {
			errFields = verr.Fields
		} else if scanErr != nil {
			errFields["Amount"] = "مبلغ معتبر نیست"
		} else if err != nil {
			logError("ثبت درخواست برداشت", err)
			http.Error(w, "خطای داخلی سرور", http.StatusInternalServerError)
			return
		}

		orders, _ := h.Orders.ListByReferringAgent(agent.ID)
		withdrawals, _ := h.Wallet.ListByAgent(agent.ID)
		render(w, r, "agent_dashboard.html", map[string]interface{}{
			"Title":            "پنل همکاران",
			"Agent":            agent,
			"Orders":           orders,
			"Withdrawals":      withdrawals,
			"WithdrawalErrors": errFields,
		})
		return
	}

	http.Redirect(w, r, "/agent/dashboard", http.StatusSeeOther)
}
