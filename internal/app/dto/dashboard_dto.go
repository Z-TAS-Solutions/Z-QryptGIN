package dto

import "time"

// AuthTrendDataPoint represents a single time-series data point for authentication trends
type AuthTrendDataPoint struct {
	Timestamp    time.Time `json:"timestamp"`
	SuccessCount int64     `json:"successCount"`
	FailureCount int64     `json:"failureCount"`
}

// GetAuthTrendsResponse represents the complete response for the authentication trends endpoint
type GetAuthTrendsResponse struct {
	Interval string               `json:"interval"`
	Data     []AuthTrendDataPoint `json:"data"`
}

// DashboardMetrics represents an aggregated view of authentication metrics
type DashboardMetrics struct {
	TotalSuccessfulAuthentications int64
	TotalFailedAuthentications     int64
	AverageSuccessPerInterval      float64
	AverageFailurePerInterval      float64
	PeakSuccessCount               int64
	PeakFailureCount               int64
	PeakSuccessTime                time.Time
	PeakFailureTime                time.Time
}
