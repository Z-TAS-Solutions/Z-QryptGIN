package dto

import "time"

// Admin Device Management DTOs
type AdminDeviceResponse struct {
	DeviceID     string    `json:"device_id"`
	UserID       string    `json:"user_id"`
	UserEmail    string    `json:"user_email"`
	DeviceName   string    `json:"device_name"`
	UserAgent    string    `json:"user_agent"`
	IPAddress    string    `json:"ip_address"`
	LastActivity time.Time `json:"last_activity"`
	CreatedAt    time.Time `json:"created_at"`
	Status       string    `json:"status"` // "active", "inactive"
}

type AdminListDevicesResponse struct {
	Success bool                    `json:"success"`
	Devices []AdminDeviceResponse   `json:"devices"`
	Total   int                     `json:"total"`
	Page    int                     `json:"page"`
	Limit   int                     `json:"limit"`
}

type AdminForceLogoutRequest struct {
	Reason string `json:"reason,omitempty"`
}

type AdminForceLogoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	DeviceID string `json:"device_id"`
}
