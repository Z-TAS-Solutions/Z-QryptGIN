package dto

// -- Fetch Notifications --
type FetchNotificationsRequest struct {
	Limit      int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset     int    `form:"offset" binding:"omitempty,min=0"`
	UnreadOnly bool   `form:"unread_only"`
	SortOrder  string `form:"sort_order"`
}

type NotificationItem struct {
	NotificationID string `json:"notification_id"`
	Title          string `json:"title"`
	Details        string `json:"details"`
	Timestamp      int64  `json:"timestamp"`
	Status         string `json:"status"`
}

type FetchNotificationsResponse struct {
	Message string `json:"message"`
	Data    struct {
		Notifications []NotificationItem `json:"notifications"`
		Pagination    struct {
			Limit    int  `json:"limit"`
			Offset   int  `json:"offset"`
			Returned int  `json:"returned"`
			HasMore  bool `json:"has_more"`
		} `json:"pagination"`
	} `json:"data"`
}

// -- Update Notification Status --
type UpdateNotificationStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type UpdateNotificationStatusResponse struct {
	Message string `json:"message"`
	Data    struct {
		NotificationID string `json:"notification_id"`
		Status         string `json:"status"`
	} `json:"data"`
}

// -- Mark All Read --
type MarkAllNotificationsReadResponse struct {
	Message string `json:"message"`
	Data    struct {
		UpdatedCount int `json:"updated_count"`
	} `json:"data"`
}
