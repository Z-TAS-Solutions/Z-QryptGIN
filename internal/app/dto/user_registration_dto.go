package dto

import "github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"

type UserRegistrationDetailsRequest struct {
	Name    database.Name        `json:"name"`
	Email   database.Email       `json:"email"`
	PhoneNo database.PhoneNumber `json:"phone_no"`
	Nic     database.NIC         `json:"nic"`
	Role    database.UserRole    `json:"role"`
}

type UserRegistrationDetailsResponse struct {
	Success  bool                  `json:"success"`
	CustomID database.UserCustomID `json:"custom_id"`
}

type UserRegistrationOTPRequest struct {
	UserID database.UserCustomID `json:"user_id"`
	OTP    string                `json:"otp"`
}

type UserRegistrationOTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ResendOTPRequest struct {
	UserID database.UserCustomID `json:"user_id"`
}

type ResendOTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
