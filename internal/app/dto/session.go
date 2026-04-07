package dto

import "time"

type MfaStatus string

const (
	MfaStatusVerified MfaStatus = "verified"
	MfaStatusPending  MfaStatus = "pending"
	MfaStatusRequired MfaStatus = "required"
)

type Session struct {
	ID           string    `json:"id"`
	UserID       uint      `json:"user_id"`
	JTI          string    `json:"jti"`
	Role         string    `json:"role"`
	DeviceName   string    `json:"device_name"`
	DeviceID     string    `json:"device_id"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	IsActive     bool      `json:"is_active"`
	MfaStatus    MfaStatus `json:"mfa_status"`
	LastActiveAt time.Time `json:"last_active_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	Location     string    `json:"location"`
}

// ActiveSessionResponse represents a single active session in the response
type ActiveSessionResponse struct {
	SessionID  string `json:"session_id"`
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	Location   string `json:"location"`
	IPAddress  string `json:"ip_address"`
	LastActive int64  `json:"last_active"` // Unix timestamp in milliseconds
	Current    bool   `json:"current"`
}

// GetActiveSessionsResponseData is the data wrapper for the get active sessions endpoint
type GetActiveSessionsResponseData struct {
	Sessions []ActiveSessionResponse `json:"sessions"`
	// Pagination PaginationInfo          `json:"pagination"`
	Pagination PaginationInfo `json:"pagination"`
}

// LogoutOthersResponseData is the data wrapper for the logout others endpoint
type LogoutOthersResponseData struct {
	SessionsTerminated int `json:"sessions_terminated"`
}
