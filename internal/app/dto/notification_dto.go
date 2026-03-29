package dto

import "time"

// Notification DTOs
type NotificationResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"` // "auth", "account", "security"
	Status    string    `json:"status"` // "unread", "read"
	CreatedAt time.Time `json:"created_at"`
}

type FetchNotificationsResponse struct {
	Success       bool                     `json:"success"`
	Notifications []NotificationResponse   `json:"notifications"`
	Total         int                      `json:"total"`
}

type MarkAllReadRequest struct {
}

type MarkAllReadResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Updated  int    `json:"updated"`
}

type UpdateNotificationStatusRequest struct {
	Status string `json:"status"` // "read", "archived"
}

type UpdateNotificationStatusResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    NotificationResponse `json:"data"`
}
