package cache

import (
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
)

type RegistrationCache struct {
	ID string `json:"id"`

	//User Data
	UserID  database.UserCustomID `json:"u_id"`
	Name    database.Name         `json:"nm"`
	Email   database.Email        `json:"em"`
	PhoneNo database.PhoneNumber  `json:"ph"`
	Nic     database.NIC          `json:"nic"`
	Role    database.UserRole     `json:"ro"`

	// Verification Data
	OTP           string                     `json:"otp"`
	MfaStatus     database.MfaStatus         `json:"mfa_s"`
	SecurityLevel database.UserSecurityLevel `json:"sec_l"`

	// Expiry Tracking (Optional but recommended for Redis)
	CreatedAt time.Time `json:"c_at"`
}
