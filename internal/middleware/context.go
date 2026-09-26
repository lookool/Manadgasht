package middleware

// نام کوکی‌های نشست؛ کاربر عادی و همکار جدا نگه داشته می‌شن تا یک
// مرورگر بتونه هم‌زمان هم به‌عنوان کاربر و هم همکار وارد باشه.
const (
	UserCookieName  = "manadgasht_user_session"
	AgentCookieName = "manadgasht_agent_session"
)

type contextKey string

const (
	userContextKey  contextKey = "current_user"
	agentContextKey contextKey = "current_agent"
)
