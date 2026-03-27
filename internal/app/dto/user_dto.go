package dto

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
)

// CreateUserRequest is what the client sends
type CreateUserRequest struct {
	Name     string               `json:"name" binding:"required"`
	Email    database.Email       `json:"email" binding:"required,email"`
	PhoneNo  database.PhoneNumber `json:"phone_no" binding:"required"`
	Nic      database.NIC         `json:"nic" binding:"required"`
	Password string               `json:"password" binding:"required,min=8"`
}

// UserResponse is what the client receives (no passwords!)
type UserResponse struct {
	Success  bool   `json:"success"`
	CustomID string `json:"custom_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}
