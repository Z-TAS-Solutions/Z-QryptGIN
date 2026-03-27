package repository

import (
	"fmt"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"gorm.io/gorm"
)

type UserSessionsRepository interface {
	FetchActiveSessions(userID string, req dto.FetchSessionsRequest) (*dto.FetchSessionsResponse, error)
	SignOutOtherDevices(userID string, currentSessionID string) (*dto.LogoutOtherSessionsResponse, error)
}

type userSessionsRepository struct {
	db *gorm.DB
}

func NewUserSessionsRepository(db *gorm.DB) UserSessionsRepository {
	return &userSessionsRepository{db: db}
}

func (r *userSessionsRepository) FetchActiveSessions(userID string, req dto.FetchSessionsRequest) (*dto.FetchSessionsResponse, error) {
	var sessions []database.Session
	var total int64

	query := r.db.Model(&database.Session{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset

	if err := query.Limit(limit).Offset(offset).Order("last_active desc").Find(&sessions).Error; err != nil {
		return nil, err
	}

	var res dto.FetchSessionsResponse
	res.Message = "Active sessions retrieved successfully"
	res.Data.Sessions = []dto.SessionItem{}

	for _, s := range sessions {
		res.Data.Sessions = append(res.Data.Sessions, dto.SessionItem{
			SessionID:  fmt.Sprint(s.SessionNo),
			DeviceID:   fmt.Sprint(s.ID),
			DeviceName: s.DeviceName,
			Location:   s.Location,
			IpAddress:  string(s.IpAddress),
			LastActive: s.LastActive.UnixMilli(),
			Current:    false, // Needs the request context to truly determine the current session being used
		})
	}

	res.Data.Pagination.Limit = limit
	res.Data.Pagination.Offset = offset
	res.Data.Pagination.Returned = len(res.Data.Sessions)
	res.Data.Pagination.HasMore = int64(offset+limit) < total

	return &res, nil
}

func (r *userSessionsRepository) SignOutOtherDevices(userID string, currentSessionID string) (*dto.LogoutOtherSessionsResponse, error) {
	tx := r.db.Where("user_id = ? AND session_no != ?", userID, currentSessionID).Delete(&database.Session{})
	if tx.Error != nil {
		return nil, tx.Error
	}

	var res dto.LogoutOtherSessionsResponse
	res.Message = "All other sessions logged out successfully"
	res.Data.SessionsTerminated = int(tx.RowsAffected)
	return &res, nil
}
