package repository

import (
	"fmt"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"gorm.io/gorm"
)

type AdminAuthLogsRepository interface {
	GetAuthLogs(req dto.AuthLogsRequest) (*dto.AuthLogsResponse, error)
	GetAuthAnalytics(req dto.AuthAnalyticsRequest) (*dto.AuthAnalyticsResponse, error)
}

type adminAuthLogsRepository struct {
	db *gorm.DB
}

func NewAdminAuthLogsRepository(db *gorm.DB) AdminAuthLogsRepository {
	return &adminAuthLogsRepository{db: db}
}

func (r *adminAuthLogsRepository) GetAuthLogs(req dto.AuthLogsRequest) (*dto.AuthLogsResponse, error) {
	var logs []database.ActivityLog
	var total int64

	query := r.db.Model(&database.ActivityLog{})

	if req.Status != "" {
		if req.Status == "success" {
			query = query.Where("type = ?", "login_Success")
		} else if req.Status == "failed" {
			query = query.Where("type", "Failed_Login")
		} else if req.Status == "suspicious" {
			query = query.Where("is_critical = ?", true)
		}
	}

	if req.Search != "" {
		query = query.Where("user_id = ?", req.Search) // Assuming we use Search as UserID for now, or match device etc.
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	limit := req.Limit
	if limit == 0 {
		limit = 10
	}
	offset := req.Offset

	if err := query.Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		return nil, err
	}

	var res dto.AuthLogsResponse
	res.Message = "Authentication logs retrieved successfully"
	res.Data.Logs = []dto.AuthLogItem{}

	for _, log := range logs {
		status := "unknown"
		if string(log.Type) == "login_Success" {
			status = "success"
		} else if string(log.Type) == "Failed_Login" {
			status = "failed"
		} else if log.IsCritical {
			status = "suspicious"
		}

		res.Data.Logs = append(res.Data.Logs, dto.AuthLogItem{
			LogID:     fmt.Sprint(log.ActivityNo),
			Timestamp: log.TimeLabel.UnixMilli(),
			UserName:  "", // Usually fetch user name or rely on user_id
			UserID:    fmt.Sprint(log.UserID),
			Status:    status,
			Location:  "", // Not in activity_log model yet
			Device:    log.Device,
		})
	}

	res.Data.Pagination.Limit = limit
	res.Data.Pagination.Offset = offset
	res.Data.Pagination.Returned = len(res.Data.Logs)
	res.Data.Pagination.HasMore = int64(offset+limit) < total

	return &res, nil
}

func (r *adminAuthLogsRepository) GetAuthAnalytics(req dto.AuthAnalyticsRequest) (*dto.AuthAnalyticsResponse, error) {
	var successfulLogins, failedLogins, suspiciousActivities int64

	baseQuery := r.db.Model(&database.ActivityLog{})

	if req.From > 0 {
		baseQuery = baseQuery.Where("time_label >= ?", time.UnixMilli(req.From))
	}
	if req.To > 0 {
		baseQuery = baseQuery.Where("time_label <= ?", time.UnixMilli(req.To))
	}

	baseQuery.Where("type = ?", "login_Success").Count(&successfulLogins)
	baseQuery.Where("type = ?", "Failed_Login").Count(&failedLogins)
	baseQuery.Where("is_critical = ?", true).Count(&suspiciousActivities)

	var res dto.AuthAnalyticsResponse
	res.Message = "Authentication analytics retrieved successfully"
	res.Data.TimeRange.From = req.From
	res.Data.TimeRange.To = req.To
	res.Data.Metrics.SuccessfulLogins = successfulLogins
	res.Data.Metrics.FailedLogins = failedLogins
	res.Data.Metrics.SuspiciousActivities = suspiciousActivities

	return &res, nil
}
