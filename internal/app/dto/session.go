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
	DeviceName   string    `json:"device_name"`
	DeviceID     string    `json:"device_id"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	IsActive     bool      `json:"is_active"`
	MfaStatus    MfaStatus `json:"mfa_status"`
	LastActiveAt time.Time `json:"last_active_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}