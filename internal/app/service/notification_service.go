package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type NotificationService interface {
	FetchNotifications(ctx context.Context, userID string) (*dto.FetchNotificationsResponse, error)
	MarkAllRead(ctx context.Context, userID string) (*dto.MarkAllReadResponse, error)
	UpdateNotificationStatus(ctx context.Context, userID, notificationID string, req dto.UpdateNotificationStatusRequest) (*dto.UpdateNotificationStatusResponse, error)
}

type notificationService struct {
	redisClient *redis.Client
}

func NewNotificationService(redisClient *redis.Client) NotificationService {
	return &notificationService{
		redisClient: redisClient,
	}
}

func (s *notificationService) FetchNotifications(ctx context.Context, userID string) (*dto.FetchNotificationsResponse, error) {
	log.Info().Str("user_id", userID).Msg("Fetching notifications")

	// Mock some notifications
	notifications := []dto.NotificationResponse{
		{
			ID:        "notif-001",
			UserID:    userID,
			Title:     "New Login Detected",
			Message:   "Your account was accessed from a new device",
			Type:      "auth",
			Status:    "unread",
			CreatedAt: time.Now().AddDate(0, 0, -1),
		},
		{
			ID:        "notif-002",
			UserID:    userID,
			Title:     "Security Update",
			Message:   "Consider updating your security settings",
			Type:      "security",
			Status:    "read",
			CreatedAt: time.Now().AddDate(0, 0, -2),
		},
		{
			ID:        "notif-003",
			UserID:    userID,
			Title:     "Account Activity",
			Message:   "Your password was changed 30 days ago",
			Type:      "account",
			Status:    "read",
			CreatedAt: time.Now().AddDate(0, 0, -30),
		},
	}

	return &dto.FetchNotificationsResponse{
		Success:       true,
		Notifications: notifications,
		Total:         len(notifications),
	}, nil
}

func (s *notificationService) MarkAllRead(ctx context.Context, userID string) (*dto.MarkAllReadResponse, error) {
	log.Info().Str("user_id", userID).Msg("Marking all notifications as read")

	// Simulate marking notifications
	updated := 3 // mocking 3 notifications updated

	return &dto.MarkAllReadResponse{
		Success: true,
		Message: "all notifications marked as read",
		Updated: updated,
	}, nil
}

func (s *notificationService) UpdateNotificationStatus(ctx context.Context, userID, notificationID string, req dto.UpdateNotificationStatusRequest) (*dto.UpdateNotificationStatusResponse, error) {
	log.Info().Str("user_id", userID).Str("notification_id", notificationID).Str("status", req.Status).Msg("Updating notification status")

	if req.Status != "read" && req.Status != "archived" {
		return nil, fmt.Errorf("invalid status: %s", req.Status)
	}

	notif := dto.NotificationResponse{
		ID:        notificationID,
		UserID:    userID,
		Title:     "Sample Notification",
		Message:   "This is a sample notification",
		Type:      "account",
		Status:    req.Status,
		CreatedAt: time.Now(),
	}

	return &dto.UpdateNotificationStatusResponse{
		Success: true,
		Message: fmt.Sprintf("notification marked as %s", req.Status),
		Data:    notif,
	}, nil
}
