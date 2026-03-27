package dto

import "time"

// -- Analytics --

type DashboardAnalyticsResponse struct {
	TotalUsers         int64   `json:"totalUsers"`
	ActiveSessions     int64   `json:"activeSessions"`
	SuccessRate        float64 `json:"successRate"`
	FailedRate         float64 `json:"failedRate"`
	SuspiciousActivity int64   `json:"suspiciousActivity"`
}

// -- Auth Trends --

type AuthTrendDataPoint struct {
	Timestamp    time.Time `json:"timestamp"`
	SuccessCount int64     `json:"successCount"`
	FailureCount int64     `json:"failureCount"`
}

type DashboardAuthTrendsResponse struct {
	Interval string               `json:"interval"`
	Data     []AuthTrendDataPoint `json:"data"`
}

// -- Recent Auth Activity --

type RecentAuthActivityRequest struct {
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1"`
	Status string `form:"status" binding:"omitempty"`
}

type AuthActivityItem struct {
	UserID    string    `json:"userId"`
	Device    string    `json:"device"`
	Method    string    `json:"method"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type RecentAuthActivityResponse struct {
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
	Total int64              `json:"total"`
	Data  []AuthActivityItem `json:"data"`
}
