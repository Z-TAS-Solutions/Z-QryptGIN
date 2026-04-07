package dto

import "time"

// Admin User Management DTOs
type AdminUserListResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Status    string    `json:"status"` // "active", "locked"
	CreatedAt time.Time `json:"created_at"`
	LastLogin time.Time `json:"last_login"`
}

type AdminUserListAllResponse struct {
	Success bool                    `json:"success"`
	Users   []AdminUserListResponse `json:"users"`
	Total   int                     `json:"total"`
	Page    int                     `json:"page"`
	Limit   int                     `json:"limit"`
}

type AdminUserDetailsResponse struct {
	ID                string    `json:"id"`
	Email             string    `json:"email"`
	Name              string    `json:"name"`
	PhoneNo           string    `json:"phone_no"`
	NIC               string    `json:"nic"`
	Role              string    `json:"role"`
	Status            string    `json:"status"` // "active", "locked"
	MFAEnabled        bool      `json:"mfa_enabled"`
	WebAuthnEnabled   bool      `json:"webauthn_enabled"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	LastLogin         time.Time `json:"last_login"`
	LoginAttempts     int       `json:"login_attempts"`
	SecurityLevel     string    `json:"security_level"`
}

type UpdateLockStatusRequest struct {
	Locked bool   `json:"locked"`
	Reason string `json:"reason,omitempty"`
}

type UpdateLockStatusResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type AdminDeleteUserRequest struct {
	Reason string `json:"reason,omitempty"`
}

type AdminDeleteUserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
