package repository

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	GetNotificationsByUserID(userID uint, limit int, offset int, unreadOnly bool, sortOrder string) ([]database.Notification, int64, error)
	GetNotificationByIDAndUserID(notificationID string, userID uint) (*database.Notification, error)
	UpdateNotificationStatus(notificationID string, userID uint, isRead bool) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

// GetNotificationsByUserID retrieves notifications for a user with filtering and pagination
// Returns notifications, total count, and error
func (r *notificationRepository) GetNotificationsByUserID(userID uint, limit int, offset int, unreadOnly bool, sortOrder string) ([]database.Notification, int64, error) {
	var notifications []database.Notification
	var totalCount int64

	query := r.db.Where("user_id = ?", userID)

	// Apply unread filter if requested
	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}

	// Get total count before pagination
	if err := query.Model(&database.Notification{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	orderDirection := "DESC" // default to newest first
	if sortOrder == "asc" {
		orderDirection = "ASC"
	}

	// Apply pagination and sorting
	if err := query.
		Order("created_at " + orderDirection).
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, totalCount, nil
}

// GetNotificationByIDAndUserID retrieves a notification by ID and verifies it belongs to the user
func (r *notificationRepository) GetNotificationByIDAndUserID(notificationID string, userID uint) (*database.Notification, error) {
	var notification database.Notification
	err := r.db.Where("notif_id = ? AND user_id = ?", notificationID, userID).First(&notification).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// UpdateNotificationStatus updates the read status of a notification (verifies ownership)
func (r *notificationRepository) UpdateNotificationStatus(notificationID string, userID uint, isRead bool) error {
	result := r.db.
		Where("notif_id = ? AND user_id = ?", notificationID, userID).
		Model(&database.Notification{}).
		Update("is_read", isRead)

	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were affected (notification not found or unauthorized)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
