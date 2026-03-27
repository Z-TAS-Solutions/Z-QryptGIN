package dto

// -- MFA Send --
type MfaSendRequest struct {
	Method   string `json:"method" binding:"required"`
	DeviceID string `json:"deviceId" binding:"required"`
}

type MfaSendResponse struct {
	Message string `json:"message"`
}

// -- MFA Respond --
type MfaRespondRequest struct {
	NotificationID string `json:"notificationId" binding:"required"`
	Action         string `json:"action" binding:"required"`
	DeviceID       string `json:"deviceId" binding:"required"`
	Timestamp      int64  `json:"timestamp"`
}

type MfaRespondResponse struct {
	Message string `json:"message"`
}
