package service

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
)

type UserSessionsService interface {
	FetchActiveSessions(userID string, req dto.FetchSessionsRequest) (*dto.FetchSessionsResponse, error)
	SignOutOtherDevices(userID string, currentSessionID string) (*dto.LogoutOtherSessionsResponse, error)
}

type userSessionsService struct {
	repo repository.UserSessionsRepository
}

func NewUserSessionsService(repo repository.UserSessionsRepository) UserSessionsService {
	return &userSessionsService{repo: repo}
}

func (s *userSessionsService) FetchActiveSessions(userID string, req dto.FetchSessionsRequest) (*dto.FetchSessionsResponse, error) {
	if req.Limit < 1 {
		req.Limit = 20
	}
	return s.repo.FetchActiveSessions(userID, req)
}

func (s *userSessionsService) SignOutOtherDevices(userID string, currentSessionID string) (*dto.LogoutOtherSessionsResponse, error) {
	return s.repo.SignOutOtherDevices(userID, currentSessionID)
}
