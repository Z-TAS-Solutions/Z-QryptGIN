package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type SessionService interface {
	FetchActiveSessions(ctx context.Context, userID string) (*dto.FetchActiveSessionsResponse, error)
	LogoutOtherDevices(ctx context.Context, userID string, req dto.LogoutOtherDevicesRequest) (*dto.LogoutOtherDevicesResponse, error)
}

type sessionService struct {
	redisClient *redis.Client
}

func NewSessionService(redisClient *redis.Client) SessionService {
	return &sessionService{
		redisClient: redisClient,
	}
}

func (s *sessionService) FetchActiveSessions(ctx context.Context, userID string) (*dto.FetchActiveSessionsResponse, error) {
	log.Info().Str("user_id", userID).Msg("Fetching active sessions")

	// Mock active sessions
	sessions := []dto.SessionDTO{
		{
			SessionID:    "sess-001-current",
			DeviceInfo:   "Chrome on Windows",
			IPAddress:    "192.168.1.100",
			UserAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			LastActivity: time.Now().Add(-5 * time.Minute),
			CreatedAt:    time.Now().AddDate(0, 0, -2),
			ExpiresAt:    time.Now().AddDate(0, 0, 28),
			Current:      true,
		},
		{
			SessionID:    "sess-002",
			DeviceInfo:   "Safari on iPhone",
			IPAddress:    "203.0.113.45",
			UserAgent:    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X)",
			LastActivity: time.Now().Add(-2 * time.Hour),
			CreatedAt:    time.Now().AddDate(0, 0, -10),
			ExpiresAt:    time.Now().AddDate(0, 0, 20),
			Current:      false,
		},
		{
			SessionID:    "sess-003",
			DeviceInfo:   "Firefox on Mac",
			IPAddress:    "198.51.100.22",
			UserAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
			LastActivity: time.Now().Add(-1 * time.Day),
			CreatedAt:    time.Now().AddDate(0, 0, -5),
			ExpiresAt:    time.Now().AddDate(0, 0, 25),
			Current:      false,
		},
	}

	return &dto.FetchActiveSessionsResponse{
		Success:  true,
		Sessions: sessions,
		Total:    len(sessions),
	}, nil
}

func (s *sessionService) LogoutOtherDevices(ctx context.Context, userID string, req dto.LogoutOtherDevicesRequest) (*dto.LogoutOtherDevicesResponse, error) {
	log.Info().Str("user_id", userID).Msg("Logging out other devices")

	// Simulate revoking sessions
	var revoked int
	if req.ExcludeSessionID != nil {
		revoked = 2 // exclude 1 session, revoke others
	} else {
		revoked = 3 // revoke all
	}

	return &dto.LogoutOtherDevicesResponse{
		Success: true,
		Message: fmt.Sprintf("successfully logged out %d device(s)", revoked),
		Revoked: revoked,
	}, nil
}
