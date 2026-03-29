package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/rs/zerolog/log"
)

type AdminAuthService interface {
	AdminLogin(ctx context.Context, req dto.AdminLoginRequest) (*dto.AdminLoginResponse, error)
	AdminRefresh(ctx context.Context, req dto.AdminRefreshRequest) (*dto.AdminRefreshResponse, error)
	GetAuthLogs(ctx context.Context, page, limit int) (*dto.AdminGetAuthLogsResponse, error)
	GetAuthAnalytics(ctx context.Context, period string) (*dto.AdminAuthAnalyticsResponse, error)
}

type adminAuthService struct {
}

func NewAdminAuthService() AdminAuthService {
	return &adminAuthService{}
}

func (s *adminAuthService) AdminLogin(ctx context.Context, req dto.AdminLoginRequest) (*dto.AdminLoginResponse, error) {
	log.Info().Str("email", req.Email).Msg("Admin login attempt")

	// Mock validation
	if req.Email == "" || req.Password == "" {
		return &dto.AdminLoginResponse{
			Success: false,
			Message: "invalid email or password",
		}, fmt.Errorf("invalid credentials")
	}

	return &dto.AdminLoginResponse{
		Success:      true,
		Message:      "login successful",
		AccessToken:  "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJBMDAxIiwibmFtZSI6IkFkbWluIn0.mock_token",
		RefreshToken: "refresh_token_mock_value_12345",
		ExpiresIn:    3600,
	}, nil
}

func (s *adminAuthService) AdminRefresh(ctx context.Context, req dto.AdminRefreshRequest) (*dto.AdminRefreshResponse, error) {
	log.Info().Msg("Admin refresh token request")

	if req.RefreshToken == "" {
		return &dto.AdminRefreshResponse{
			Success: false,
		}, fmt.Errorf("invalid refresh token")
	}

	return &dto.AdminRefreshResponse{
		Success:     true,
		AccessToken: "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJBMDAxIiwibmFtZSI6IkFkbWluIn0.new_mock_token",
		ExpiresIn:   3600,
	}, nil
}

func (s *adminAuthService) GetAuthLogs(ctx context.Context, page, limit int) (*dto.AdminGetAuthLogsResponse, error) {
	log.Info().Int("page", page).Int("limit", limit).Msg("Fetching auth logs")

	logs := []dto.AuthLogEntry{
		{
			ID:        "log-001",
			UserID:    "U001",
			Email:     "john@example.com",
			Type:      "login",
			Status:    "success",
			IPAddress: "192.168.1.100",
			UserAgent: "Chrome on Windows",
			Timestamp: time.Now().Add(-5 * time.Minute),
		},
		{
			ID:        "log-002",
			UserID:    "U002",
			Email:     "jane@example.com",
			Type:      "mfa_attempt",
			Status:    "success",
			IPAddress: "203.0.113.45",
			UserAgent: "Safari on iPhone",
			Timestamp: time.Now().Add(-15 * time.Minute),
		},
		{
			ID:        "log-003",
			UserID:    "U003",
			Email:     "bob@example.com",
			Type:      "login",
			Status:    "failed",
			Reason:    "invalid password",
			IPAddress: "198.51.100.22",
			UserAgent: "Firefox on Mac",
			Timestamp: time.Now().Add(-25 * time.Minute),
		},
	}

	return &dto.AdminGetAuthLogsResponse{
		Success: true,
		Logs:    logs,
		Total:   48392,
		Page:    page,
		Limit:   limit,
	}, nil
}

func (s *adminAuthService) GetAuthAnalytics(ctx context.Context, period string) (*dto.AdminAuthAnalyticsResponse, error) {
	log.Info().Str("period", period).Msg("Fetching auth analytics")

	analytics := make([]dto.AuthAnalyticsPoint, 7)
	now := time.Now()

	for i := 0; i < 7; i++ {
		analytics[i] = dto.AuthAnalyticsPoint{
			Timestamp: now.AddDate(0, 0, -(6-i)),
			Logins:    int64(5000 + i*300),
			Failures:  int64(250 + i*20),
			MFASent:   int64(4800 + i*280),
			MFAFailed: int64(150 + i*10),
		}
	}

	return &dto.AdminAuthAnalyticsResponse{
		Success:       true,
		Period:        period,
		Analytics:     analytics,
		TotalLogins:   48392,
		TotalFailures: 2341,
	}, nil
}
