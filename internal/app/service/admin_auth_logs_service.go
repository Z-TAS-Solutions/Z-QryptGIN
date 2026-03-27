package service

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
)

type AdminAuthLogsService interface {
	GetAuthLogs(req dto.AuthLogsRequest) (*dto.AuthLogsResponse, error)
	GetAuthAnalytics(req dto.AuthAnalyticsRequest) (*dto.AuthAnalyticsResponse, error)
}

type adminAuthLogsService struct {
	repo repository.AdminAuthLogsRepository
}

func NewAdminAuthLogsService(repo repository.AdminAuthLogsRepository) AdminAuthLogsService {
	return &adminAuthLogsService{repo: repo}
}

func (s *adminAuthLogsService) GetAuthLogs(req dto.AuthLogsRequest) (*dto.AuthLogsResponse, error) {
	if req.Limit < 1 {
		req.Limit = 20
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}
	return s.repo.GetAuthLogs(req)
}

func (s *adminAuthLogsService) GetAuthAnalytics(req dto.AuthAnalyticsRequest) (*dto.AuthAnalyticsResponse, error) {
	return s.repo.GetAuthAnalytics(req)
}
