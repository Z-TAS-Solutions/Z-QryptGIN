package dto

import "time"

// Admin Dashboard DTOs
type AnalyticsResponse struct {
	TotalUsers          int64     `json:"total_users"`
	ActiveUsers         int64     `json:"active_users"`
	NewUsersThisMonth   int64     `json:"new_users_this_month"`
	AuthAttempts        int64     `json:"auth_attempts"`
	FailedAttempts      int64     `json:"failed_attempts"`
	SuccessRate         float64   `json:"success_rate"`
	AverageAuthTime     float64   `json:"average_auth_time_ms"`
}

type AuthTrendPoint struct {
	Timestamp  time.Time `json:"timestamp"`
	Count      int64     `json:"count"`
	SuccessSum int64     `json:"success_sum"`
	FailureSum int64     `json:"failure_sum"`
}

type AuthTrendsResponse struct {
	Period string            `json:"period"` // "24h", "7d", "30d"
	Trends []AuthTrendPoint  `json:"trends"`
}

type RecentAuthActivity struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Type      string    `json:"type"` // "login", "registration", "logout"
	Status    string    `json:"status"` // "success", "failed"
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Timestamp time.Time `json:"timestamp"`
}

type RecentAuthActivityResponse struct {
	Activities []RecentAuthActivity `json:"activities"`
	Total      int                  `json:"total"`
}
