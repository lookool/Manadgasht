package middleware

const (
	UserCookieName  = "manadgasht_user_session"
	AgentCookieName = "manadgasht_agent_session"
	AdminCookieName = "manadgasht_admin_session"
)

type contextKey string

const (
	userContextKey  contextKey = "current_user"
	agentContextKey contextKey = "current_agent"
	adminContextKey contextKey = "current_admin"
)