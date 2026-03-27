package service

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
)

type DashboardService interface {
	GetAnalytics() (*dto.DashboardAnalyticsResponse, error)
	GetAuthTrends(interval string) (*dto.DashboardAuthTrendsResponse, error)
	GetRecentAuthActivity(page int, limit int, status string) (*dto.RecentAuthActivityResponse, error)
}

type dashboardService struct {
	repo repository.DashboardRepository
}

func NewDashboardService(repo repository.DashboardRepository) DashboardService {
	return &dashboardService{repo: repo}
}

func (s *dashboardService) GetAnalytics() (*dto.DashboardAnalyticsResponse, error) {
	return s.repo.GetDashboardAnalytics()
}

func (s *dashboardService) GetAuthTrends(interval string) (*dto.DashboardAuthTrendsResponse, error) {
	// Set default mapping if not provided
	if interval == "" {
		interval = "hour"
	}
	return s.repo.GetDashboardAuthTrends(interval)
}

func (s *dashboardService) GetRecentAuthActivity(page int, limit int, status string) (*dto.RecentAuthActivityResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.GetRecentAuthActivity(page, limit, status)
}
