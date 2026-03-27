package dto

import "github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"

type UserRegistrationCache struct {
	ID        string                `json:"id"`
	UserID    database.UserCustomID `json:"user_id"`
	Email     database.Email        `json:"email"`
	PhoneNo   database.PhoneNumber  `json:"phone_no"`
	Nic       database.NIC          `json:"nic"`
	Role      database.UserRole     `json:"role"`
	OTP       string                `json:"otp"`
	MfaStatus database.MfaStatus    `json:"mfa_status",default:"Pending"`
	SecurityLevel database.UserSecurityLevel `json:"security_level",default:"Low"`

}
