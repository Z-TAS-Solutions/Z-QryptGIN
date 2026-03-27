package dto

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
)

// CreateUserRequest is what the client sends
type CreateUserRequest struct {
	Name     database.Name        `json:"name" binding:"required"`
	Email    database.Email       `json:"email" binding:"required,email"`
	PhoneNo  database.PhoneNumber `json:"phone_no" binding:"required"`
	Nic      database.NIC         `json:"nic" binding:"required"`
	Password string               `json:"password" binding:"required,min=8"`
}

type CreateUserRequestOTP struct {
	UserID database.UserCustomID `json:"user_id" binding:"required"`
	OTP    string                `json:"otp" binding:"required,len=6"`
}

type CreateUserResponseOTP struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CreateUserResponse is what the client receives (no passwords!)
type CreateUserResponse struct {
	Success  bool                  `json:"success"`
	CustomID database.UserCustomID `json:"custom_id"`
	Name     database.Name         `json:"name"`
	Email    database.Email        `json:"email"`
}
