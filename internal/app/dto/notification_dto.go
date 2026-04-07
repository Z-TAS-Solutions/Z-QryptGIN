package dto

import "time"

const (
	NotificationStatusRead   NotificationStatus = "read"
	NotificationStatusUnread NotificationStatus = "unread"
)

// PaginationInfo contains pagination metadata
type PaginationInfo struct {
	Limit    int  `json:"limit"`
	Offset   int  `json:"offset"`
	Returned int  `json:"returned"`
	HasMore  bool `json:"has_more"`
}

// NotificationResponse represents a single notification in the response
type NotificationResponse struct {
	ID        string             `json:"id"`
	Title     string             `json:"title"`
	Details   string             `json:"details"`
	Timestamp int64              `json:"timestamp"` // Unix milliseconds
	Status    NotificationStatus `json:"status"`
}

// NotificationsResponseData contains the response payload
type NotificationsResponseData struct {
	Notifications []NotificationResponse `json:"notifications"`
	Pagination    PaginationInfo         `json:"pagination"`
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
