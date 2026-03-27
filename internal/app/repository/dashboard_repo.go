package repository

import (
	"fmt"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetDashboardAnalytics() (*dto.DashboardAnalyticsResponse, error)
	GetDashboardAuthTrends(interval string) (*dto.DashboardAuthTrendsResponse, error)
	GetRecentAuthActivity(page int, limit int, status string) (*dto.RecentAuthActivityResponse, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetDashboardAnalytics() (*dto.DashboardAnalyticsResponse, error) {
	var res dto.DashboardAnalyticsResponse

	var totalUsers int64
	var activeSessions int64
	var successCount int64
	var failedCount int64
	var suspiciousCount int64

	// Get Total Users
	r.db.Model(&database.User{}).Count(&totalUsers)

	// Get Active Sessions
	r.db.Model(&database.Session{}).Count(&activeSessions)

	// Get Auth Rates (Success/Failed)
	r.db.Model(&database.ActivityLog{}).Where("type = ?", "login_Success").Count(&successCount)
	r.db.Model(&database.ActivityLog{}).Where("type = ?", "Failed_Login").Count(&failedCount)

	// Get Suspicious Activities
	r.db.Model(&database.ActivityLog{}).Where("is_critical = ?", true).Count(&suspiciousCount)

	totalAttempts := successCount + failedCount
	successRate := 0.0
	failedRate := 0.0

	if totalAttempts > 0 {
		successRate = (float64(successCount) / float64(totalAttempts)) * 100
		failedRate = (float64(failedCount) / float64(totalAttempts)) * 100
	}

	res.TotalUsers = totalUsers
	res.ActiveSessions = activeSessions
	res.SuccessRate = successRate
	res.FailedRate = failedRate
	res.SuspiciousActivity = suspiciousCount

	return &res, nil
}

func (r *dashboardRepository) GetDashboardAuthTrends(interval string) (*dto.DashboardAuthTrendsResponse, error) {
	var res dto.DashboardAuthTrendsResponse
	res.Interval = interval
	res.Data = []dto.AuthTrendDataPoint{}

	// Calculate trends based on last 7 days as default interval
	// To make this fully DB agnostic in raw GORM, we'll fetch logs from last 7 days
	// and aggregate them by day in-memory rather than relying on DATE_TRUNC or DATE()
	// that vary between Postgres/MySQL/SQLite.

	startDate := time.Now().AddDate(0, 0, -7)
	var logs []database.ActivityLog

	r.db.Where("time_label >= ?", startDate).
		Where("type IN ?", []string{"login_Success", "Failed_Login"}).
		Order("time_label asc").
		Find(&logs)

	// Group by Year/Month/Day
	trendsMap := make(map[string]*dto.AuthTrendDataPoint)
	var order []string

	for i := 0; i <= 7; i++ {
		day := startDate.AddDate(0, 0, i).Format("2006-01-02")
		trendsMap[day] = &dto.AuthTrendDataPoint{
			Timestamp:    startDate.AddDate(0, 0, i),
			SuccessCount: 0,
			FailureCount: 0,
		}
		order = append(order, day)
	}

	for _, log := range logs {
		day := log.TimeLabel.Format("2006-01-02")
		if point, exists := trendsMap[day]; exists {
			if string(log.Type) == "login_Success" {
				point.SuccessCount++
			} else if string(log.Type) == "Failed_Login" {
				point.FailureCount++
			}
		}
	}

	for _, day := range order {
		res.Data = append(res.Data, *trendsMap[day])
	}

	return &res, nil
}

func (r *dashboardRepository) GetRecentAuthActivity(page int, limit int, status string) (*dto.RecentAuthActivityResponse, error) {
	var logs []database.ActivityLog
	var total int64
	var res dto.RecentAuthActivityResponse

	query := r.db.Model(&database.ActivityLog{})

	if status != "" {
		if status == "success" {
			query = query.Where("type = ?", "login_Success")
		} else if status == "failed" {
			query = query.Where("type = ?", "Failed_Login")
		}
	} else {
		query = query.Where("type IN ?", []string{"login_Success", "Failed_Login"})
	}

	query.Count(&total)

	if limit <= 0 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	query.Order("time_label desc").Limit(limit).Offset(offset).Find(&logs)

	res.Page = page
	res.Limit = limit
	res.Total = total
	res.Data = []dto.AuthActivityItem{}

	for _, log := range logs {
		logStatus := "unknown"
		if string(log.Type) == "login_Success" {
			logStatus = "success"
		} else if string(log.Type) == "Failed_Login" {
			logStatus = "failed"
		}

		res.Data = append(res.Data, dto.AuthActivityItem{
			UserID:    fmt.Sprint(log.UserID),
			Device:    log.Device,
			Method:    "Biometric/Passkey", // Assuming primary auth mode for this ecosystem
			Status:    logStatus,
			Timestamp: log.TimeLabel,
		})
	}

	return &res, nil
}
