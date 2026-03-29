package dto

// NotificationStatus represents the read status of a notification
type NotificationStatus string

const (
	NotificationStatusRead   NotificationStatus = "read"
	NotificationStatusUnread NotificationStatus = "unread"
)

// NotificationResponse represents a single notification in the response
type NotificationResponse struct {
	Title     string             `json:"title"`
	Details   string             `json:"details"`
	Timestamp int64              `json:"timestamp"` // Unix milliseconds
	Status    NotificationStatus `json:"status"`
}

// PaginationInfo contains pagination metadata
type PaginationInfo struct {
	Limit    int  `json:"limit"`
	Offset   int  `json:"offset"`
	Returned int  `json:"returned"`
	HasMore  bool `json:"has_more"`
}

// NotificationsResponseData contains the response payload
type NotificationsResponseData struct {
	Notifications []NotificationResponse `json:"notifications"`
	Pagination    PaginationInfo         `json:"pagination"`
}

// GetNotificationsResponse is the complete response for the notifications endpoint
type GetNotificationsResponse struct {
	Message string                    `json:"message"`
	Data    NotificationsResponseData `json:"data"`
}

// FetchNotificationsRequest holds query parameters for fetching notifications
type FetchNotificationsRequest struct {
	Limit      int    `form:"limit,default=20" binding:"min=1,max=100"`
	Offset     int    `form:"offset,default=0" binding:"min=0"`
	UnreadOnly bool   `form:"unread_only,default=false"`
	SortOrder  string `form:"sort_order,default=desc" binding:"oneof=asc desc"`
}

// UpdateNotificationStatusRequest holds the update request for notification status
type UpdateNotificationStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=read unread"`
}

// UpdateNotificationStatusResponse is the response for updating notification status
type UpdateNotificationStatusResponse struct {
	Message string `json:"message"`
	Data    struct {
		NotificationID string `json:"notification_id"`
		Status         string `json:"status"`
	} `json:"data"`
}

// MarkAllAsReadResponse is the response for marking all notifications as read
type MarkAllAsReadResponse struct {
	Message string `json:"message"`
	Data    struct {
		UpdatedCount int64 `json:"updated_count"`
	} `json:"data"`
}
