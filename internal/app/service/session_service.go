package service

import (
	"context"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
)

type SessionService interface {
	GetActiveSessionsByUserID(ctx context.Context, userID uint) (*dto.GetActiveSessionsResponse, error)
}

type sessionService struct {
	sessionRepo repository.SessionRepository
}

func NewSessionService(sessionRepo repository.SessionRepository) SessionService {
	return &sessionService{
		sessionRepo: sessionRepo,
	}
}

// GetActiveSessionsByUserID retrieves all active sessions for a user and returns formatted response
func (s *sessionService) GetActiveSessionsByUserID(ctx context.Context, userID uint) (*dto.GetActiveSessionsResponse, error) {
	sessions, err := s.sessionRepo.GetActiveSessionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	activeSessionResponses := make([]dto.ActiveSessionResponse, 0, len(sessions))
	now := time.Now()

	for _, session := range sessions {
		// Check if the session is expired
		isExpired := now.After(session.ExpiresAt)

		activeSessionResponse := dto.ActiveSessionResponse{
			JTI:          session.JTI,
			DeviceName:   session.DeviceName,
			DeviceID:     session.DeviceID,
			IPAddress:    session.IPAddress,
			UserAgent:    session.UserAgent,
			MfaStatus:    session.MfaStatus,
			LastActiveAt: session.LastActiveAt,
			ExpiresAt:    session.ExpiresAt,
			IsExpired:    isExpired,
		}

		activeSessionResponses = append(activeSessionResponses, activeSessionResponse)
	}

	return &dto.GetActiveSessionsResponse{
		TotalActiveCount: len(activeSessionResponses),
		Sessions:         activeSessionResponses,
	}, nil
}
