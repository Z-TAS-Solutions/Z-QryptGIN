package service

import (
	"context"
	"errors"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
)

// DashboardService handles business logic for admin dashboard analytics
type DashboardService interface {
	// GetAuthenticationTrends retrieves authentication trends for the specified interval
	// Validates interval parameter and returns time-series data
	GetAuthenticationTrends(ctx context.Context, interval string) (*dto.GetAuthTrendsResponse, error)

	// GetDashboardMetrics retrieves aggregated dashboard metrics for summary display
	GetDashboardMetrics(ctx context.Context) (*dto.DashboardMetrics, error)
}

type dashboardService struct {
	dashboardRepo repository.DashboardRepository
}

// NewDashboardService creates a new dashboard service instance
func NewDashboardService(dashboardRepo repository.DashboardRepository) DashboardService {
	return &dashboardService{
		dashboardRepo: dashboardRepo,
	}
}

// GetAuthenticationTrends retrieves auth trends with validation of interval parameter
func (s *dashboardService) GetAuthenticationTrends(ctx context.Context, interval string) (*dto.GetAuthTrendsResponse, error) {
	// Validate interval parameter
	if interval != "minute" && interval != "hour" {
		return nil, errors.New("Invalid interval. Supported values are 'minute' or 'hour'")
	}

	// Determine how many hours to fetch based on interval
	// For granularity:
	// - minute interval: fetch last 24 hours (1440 data points)
	// - hour interval: fetch last 24 hours (24 data points)
	lastHours := 24

	// Fetch trends from repository
	trends, err := s.dashboardRepo.GetAuthTrendsByInterval(ctx, interval, lastHours)
	if err != nil {
		return nil, err
	}

	return &dto.GetAuthTrendsResponse{
		Interval: interval,
		Data:     trends,
	}, nil
}

// GetDashboardMetrics retrieves aggregated metrics for dashboard display
func (s *dashboardService) GetDashboardMetrics(ctx context.Context) (*dto.DashboardMetrics, error) {
	return s.dashboardRepo.GetAuthTrendsMetrics(ctx, 24)
}
