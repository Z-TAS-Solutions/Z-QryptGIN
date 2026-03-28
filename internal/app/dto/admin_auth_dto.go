package dto

import "time"

// Admin Auth DTOs
type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminLoginResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type AdminRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AdminRefreshResponse struct {
	Success     bool   `json:"success"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// Auth Logs DTOs
type AuthLogEntry struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Type      string    `json:"type"` // "login", "logout", "mfa_attempt", "registration"
	Status    string    `json:"status"` // "success", "failed"
	Reason    string    `json:"reason,omitempty"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Timestamp time.Time `json:"timestamp"`
}

type AdminGetAuthLogsResponse struct {
	Success bool           `json:"success"`
	Logs    []AuthLogEntry `json:"logs"`
	Total   int            `json:"total"`
	Page    int            `json:"page"`
	Limit   int            `json:"limit"`
}

// Auth Analytics DTOs
type AuthAnalyticsPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Logins    int64     `json:"logins"`
	Failures  int64     `json:"failures"`
	MFASent   int64     `json:"mfa_sent"`
	MFAFailed int64     `json:"mfa_failed"`
}

type AdminAuthAnalyticsResponse struct {
	Success       bool                    `json:"success"`
	Period        string                  `json:"period"` // "24h", "7d", "30d"
	Analytics     []AuthAnalyticsPoint    `json:"analytics"`
	TotalLogins   int64                   `json:"total_logins"`
	TotalFailures int64                   `json:"total_failures"`
}
