package dto

import "time"

// Session Management DTOs
type SessionDTO struct {
	SessionID    string    `json:"session_id"`
	DeviceInfo   string    `json:"device_info"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	LastActivity time.Time `json:"last_activity"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	Current      bool      `json:"current"`
}

type FetchActiveSessionsResponse struct {
	Success  bool         `json:"success"`
	Sessions []SessionDTO `json:"sessions"`
	Total    int          `json:"total"`
}

type LogoutOtherDevicesRequest struct {
	ExcludeSessionID *string `json:"exclude_session_id,omitempty"`
}

type LogoutOtherDevicesResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Revoked  int    `json:"revoked"`
}
