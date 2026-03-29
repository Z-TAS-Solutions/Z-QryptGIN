package service

import (
	"context"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
)

type SessionService interface {
	GetActiveSessionsByUserID(ctx context.Context, userID uint, currentJTI string, limit int, offset int) (*dto.GetActiveSessionsResponseData, error)
	LogoutOtherSessions(ctx context.Context, userID uint, currentJTI string) (int, error)
}

type sessionService struct {
	sessionRepo repository.SessionRepository
}

func NewSessionService(sessionRepo repository.SessionRepository) SessionService {
	return &sessionService{
		sessionRepo: sessionRepo,
	}
}

// GetActiveSessionsByUserID retrieves all active sessions for a user and returns formatted response with pagination
func (s *sessionService) GetActiveSessionsByUserID(ctx context.Context, userID uint, currentJTI string, limit int, offset int) (*dto.GetActiveSessionsResponseData, error) {
	// Validate and set default pagination values
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	sessions, err := s.sessionRepo.GetActiveSessionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Apply pagination
	total := len(sessions)
	start := offset
	end := offset + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedSessions := sessions[start:end]

	activeSessionResponses := make([]dto.ActiveSessionResponse, 0, len(paginatedSessions))

	for _, session := range paginatedSessions {
		// Check if this is the current session
		isCurrent := session.JTI == currentJTI

		activeSessionResponse := dto.ActiveSessionResponse{
			SessionID:  session.JTI,
			DeviceID:   session.DeviceID,
			DeviceName: session.DeviceName,
			Location:   session.Location,
			IPAddress:  session.IPAddress,
			LastActive: session.LastActiveAt.UnixMilli(), // Convert to milliseconds
			Current:    isCurrent,
		}

		activeSessionResponses = append(activeSessionResponses, activeSessionResponse)
	}

	hasMore := end < total

	return &dto.GetActiveSessionsResponseData{
		Sessions: activeSessionResponses,
		Pagination: dto.PaginationInfo{
			Limit:    limit,
			Offset:   offset,
			Returned: len(activeSessionResponses),
			HasMore:  hasMore,
		},
	}, nil
}

// LogoutOtherSessions logs out all other sessions for a user except the current one
func (s *sessionService) LogoutOtherSessions(ctx context.Context, userID uint, currentJTI string) (int, error) {
	sessions, err := s.sessionRepo.GetActiveSessionsByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	terminatedCount := 0

	for _, session := range sessions {
		// Skip the current session
		if session.JTI == currentJTI {
			continue
		}

		// Revoke the session
		err := s.sessionRepo.RevokeSession(ctx, session.JTI)
		if err != nil {
			// Log the error but continue to revoke other sessions
			continue
		}

		terminatedCount++
	}

	return terminatedCount, nil
}
