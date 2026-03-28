package service

import (
	"context"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/rs/zerolog/log"
)

type AdminDashboardService interface {
	GetAnalytics(ctx context.Context) (*dto.AnalyticsResponse, error)
	GetAuthTrends(ctx context.Context, period string) (*dto.AuthTrendsResponse, error)
	GetRecentAuthActivity(ctx context.Context) (*dto.RecentAuthActivityResponse, error)
}

type adminDashboardService struct {
}

func NewAdminDashboardService() AdminDashboardService {
	return &adminDashboardService{}
}

func (s *adminDashboardService) GetAnalytics(ctx context.Context) (*dto.AnalyticsResponse, error) {
	log.Info().Msg("Fetching dashboard analytics")

	analytics := &dto.AnalyticsResponse{
		TotalUsers:        1245,
		ActiveUsers:       892,
		NewUsersThisMonth: 123,
		AuthAttempts:      48392,
		FailedAttempts:    2341,
		SuccessRate:       95.16,
		AverageAuthTime:   245.5,
	}

	return analytics, nil
}

func (s *adminDashboardService) GetAuthTrends(ctx context.Context, period string) (*dto.AuthTrendsResponse, error) {
	log.Info().Str("period", period).Msg("Fetching auth trends")

	// Mock 7 data points for the period
	trends := make([]dto.AuthTrendPoint, 7)
	now := time.Now()

	for i := 0; i < 7; i++ {
		trends[i] = dto.AuthTrendPoint{
			Timestamp:  now.AddDate(0, 0, -(6-i)),
			Count:      int64(6000 + i*500),
			SuccessSum: int64(5700 + i*450),
			FailureSum: int64(300 + i*50),
		}
	}

	return &dto.AuthTrendsResponse{
		Period: period,
		Trends: trends,
	}, nil
}

func (s *adminDashboardService) GetRecentAuthActivity(ctx context.Context) (*dto.RecentAuthActivityResponse, error) {
	log.Info().Msg("Fetching recent auth activity")

	activities := []dto.RecentAuthActivity{
		{
			ID:        "log-001",
			UserID:    "U001",
			Email:     "user1@example.com",
			Type:      "login",
			Status:    "success",
			IPAddress: "192.168.1.100",
			UserAgent: "Chrome on Windows",
			Timestamp: time.Now().Add(-5 * time.Minute),
		},
		{
			ID:        "log-002",
			UserID:    "U002",
			Email:     "user2@example.com",
			Type:      "mfa_attempt",
			Status:    "success",
			IPAddress: "203.0.113.45",
			UserAgent: "Safari on iPhone",
			Timestamp: time.Now().Add(-15 * time.Minute),
		},
		{
			ID:        "log-003",
			UserID:    "U003",
			Email:     "user3@example.com",
			Type:      "login",
			Status:    "failed",
			IPAddress: "198.51.100.22",
			UserAgent: "Firefox on Mac",
			Timestamp: time.Now().Add(-25 * time.Minute),
		},
		{
			ID:        "log-004",
			UserID:    "U004",
			Email:     "user4@example.com",
			Type:      "registration",
			Status:    "success",
			IPAddress: "192.168.0.50",
			UserAgent: "Chrome on Android",
			Timestamp: time.Now().Add(-1 * time.Hour),
		},
	}

	return &dto.RecentAuthActivityResponse{
		Activities: activities,
		Total:      len(activities),
	}, nil
}
