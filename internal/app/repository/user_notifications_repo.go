package repository

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"gorm.io/gorm"
)

type UserNotificationsRepository interface {
	FetchNotifications(userID string, req dto.FetchNotificationsRequest) (*dto.FetchNotificationsResponse, error)
	UpdateStatus(userID string, notificationID string, status string) (*dto.UpdateNotificationStatusResponse, error)
	MarkAllRead(userID string) (*dto.MarkAllNotificationsReadResponse, error)
}

type userNotificationsRepository struct {
	db *gorm.DB
}

func NewUserNotificationsRepository(db *gorm.DB) UserNotificationsRepository {
	return &userNotificationsRepository{db: db}
}

func (r *userNotificationsRepository) FetchNotifications(userID string, req dto.FetchNotificationsRequest) (*dto.FetchNotificationsResponse, error) {
	var notifications []database.Notification
	var total int64

	query := r.db.Model(&database.Notification{}).Where("user_id = ?", userID)

	if req.UnreadOnly {
		query = query.Where("is_read = ?", false)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset

	order := "created_at desc"
	if req.SortOrder == "asc" {
		order = "created_at asc"
	}

	if err := query.Limit(limit).Offset(offset).Order(order).Find(&notifications).Error; err != nil {
		return nil, err
	}

	var res dto.FetchNotificationsResponse
	res.Message = "Notifications retrieved successfully"
	res.Data.Notifications = []dto.NotificationItem{}

	for _, n := range notifications {
		status := "unread"
		if n.IsRead {
			status = "read"
		}

		res.Data.Notifications = append(res.Data.Notifications, dto.NotificationItem{
			NotificationID: string(n.NotifiID),
			Title:          n.Title,
			Details:        n.Message,
			Timestamp:      n.CreatedAt.UnixMilli(),
			Status:         status,
		})
	}

	res.Data.Pagination.Limit = limit
	res.Data.Pagination.Offset = offset
	res.Data.Pagination.Returned = len(res.Data.Notifications)
	res.Data.Pagination.HasMore = int64(offset+limit) < total

	return &res, nil
}

func (r *userNotificationsRepository) UpdateStatus(userID string, notificationID string, status string) (*dto.UpdateNotificationStatusResponse, error) {
	isRead := false
	if status == "read" {
		isRead = true
	}

	err := r.db.Model(&database.Notification{}).
		Where("user_id = ? AND notifi_id = ?", userID, notificationID).
		Update("is_read", isRead).Error

	if err != nil {
		return nil, err
	}

	var res dto.UpdateNotificationStatusResponse
	res.Message = "Notification status updated successfully"
	res.Data.NotificationID = notificationID
	res.Data.Status = status
	return &res, nil
}

func (r *userNotificationsRepository) MarkAllRead(userID string) (*dto.MarkAllNotificationsReadResponse, error) {
	tx := r.db.Model(&database.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true)

	if tx.Error != nil {
		return nil, tx.Error
	}

	var res dto.MarkAllNotificationsReadResponse
	res.Message = "All notifications marked as read"
	res.Data.UpdatedCount = int(tx.RowsAffected)
	return &res, nil
}
