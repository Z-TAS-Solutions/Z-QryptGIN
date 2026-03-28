package dto

import "github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"

// Profile DTOs
type UserProfileResponse struct {
	UserID    database.UserCustomID `json:"user_id"`
	Name      database.Name         `json:"name"`
	Email     database.Email        `json:"email"`
	PhoneNo   database.PhoneNumber  `json:"phone_no"`
	NIC       database.NIC          `json:"nic"`
	Role      database.UserRole     `json:"role"`
	CreatedAt string                `json:"created_at"`
	UpdatedAt string                `json:"updated_at"`
}

type UpdateProfileRequest struct {
	Name    *database.Name       `json:"name,omitempty"`
	PhoneNo *database.PhoneNumber `json:"phone_no,omitempty"`
}

type UpdateProfileResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *UserProfileResponse `json:"data,omitempty"`
}
