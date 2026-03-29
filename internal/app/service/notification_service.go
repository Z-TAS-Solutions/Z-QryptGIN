package service

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
)

type NotificationService interface {
	GetNotificationsByUserID(userID uint, limit int, offset int, unreadOnly bool, sortOrder string) (*dto.GetNotificationsResponse, error)
	UpdateNotificationStatus(userID uint, notificationID string, status string) (*dto.UpdateNotificationStatusResponse, error)
	MarkAllAsRead(userID uint) (*dto.MarkAllAsReadResponse, error)
}

type notificationService struct {
	notificationRepo repository.NotificationRepository
}

func NewNotificationService(notificationRepo repository.NotificationRepository) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
	}
}

// GetNotificationsByUserID retrieves notifications for a user with filtering, pagination, and sorting
func (s *notificationService) GetNotificationsByUserID(userID uint, limit int, offset int, unreadOnly bool, sortOrder string) (*dto.GetNotificationsResponse, error) {
	// Fetch notifications from repository
	notifications, totalCount, err := s.notificationRepo.GetNotificationsByUserID(userID, limit, offset, unreadOnly, sortOrder)
	if err != nil {
		return nil, err
	}

	// Transform database models to DTOs
	notificationResponses := make([]dto.NotificationResponse, 0, len(notifications))
	for _, notif := range notifications {
		status := dto.NotificationStatusRead
		if !notif.IsRead {
			status = dto.NotificationStatusUnread
		}

		notificationResponses = append(notificationResponses, dto.NotificationResponse{
			ID:        string(notif.NotifiID),
			Title:     notif.Title,
			Details:   notif.Message,
			Timestamp: notif.CreatedAt.UnixMilli(),
			Status:    status,
		})
	}

	// Calculate pagination info
	returned := len(notificationResponses)
	hasMore := int64(offset+returned) < totalCount

	return &dto.GetNotificationsResponse{
		Message: "Notifications retrieved successfully",
		Data: dto.NotificationsResponseData{
			Notifications: notificationResponses,
			Pagination: dto.PaginationInfo{
				Limit:    limit,
				Offset:   offset,
				Returned: returned,
				HasMore:  hasMore,
			},
		},
	}, nil
}

// UpdateNotificationStatus updates the read status of a notification
// Verifies the notification belongs to the user before updating
func (s *notificationService) UpdateNotificationStatus(userID uint, notificationID string, status string) (*dto.UpdateNotificationStatusResponse, error) {
	// Verify notification exists and belongs to user
	notification, err := s.notificationRepo.GetNotificationByIDAndUserID(notificationID, userID)
	if err != nil {
		return nil, err
	}

	// Convert status string to boolean
	isRead := status == "read"

	// Update the notification status
	if err := s.notificationRepo.UpdateNotificationStatus(notificationID, userID, isRead); err != nil {
		return nil, err
	}

	// Build and return response
	response := &dto.UpdateNotificationStatusResponse{
		Message: "Notification status updated successfully",
	}
	response.Data.NotificationID = string(notification.NotifiID)
	response.Data.Status = status

	return response, nil
}

// MarkAllAsRead marks all unread notifications for a user as read
// Fully idempotent - returns 0 if all are already read
func (s *notificationService) MarkAllAsRead(userID uint) (*dto.MarkAllAsReadResponse, error) {
	// Update all unread notifications for the user
	updatedCount, err := s.notificationRepo.MarkAllNotificationsAsRead(userID)
	if err != nil {
		return nil, err
	}

	// Build and return response
	response := &dto.MarkAllAsReadResponse{
		Message: "All notifications marked as read",
	}
	response.Data.UpdatedCount = updatedCount

	return response, nil
}
