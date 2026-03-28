package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
	"github.com/rs/zerolog/log"
)

type AdminUserService interface {
	ListUsers(ctx context.Context, page, limit int) (*dto.AdminUserListAllResponse, error)
	GetUserDetails(ctx context.Context, userID string) (*dto.AdminUserDetailsResponse, error)
	UpdateLockStatus(ctx context.Context, userID string, req dto.UpdateLockStatusRequest) (*dto.UpdateLockStatusResponse, error)
	DeleteUser(ctx context.Context, userID string, req dto.AdminDeleteUserRequest) (*dto.AdminDeleteUserResponse, error)
}

type adminUserService struct {
	userRepo repository.UserRepository
}

func NewAdminUserService(userRepo repository.UserRepository) AdminUserService {
	return &adminUserService{
		userRepo: userRepo,
	}
}

func (s *adminUserService) ListUsers(ctx context.Context, page, limit int) (*dto.AdminUserListAllResponse, error) {
	log.Info().Int("page", page).Int("limit", limit).Msg("Listing users")

	// Mock users list
	users := []dto.AdminUserListResponse{
		{
			ID:        "U001",
			Email:     "john@example.com",
			Name:      "John Doe",
			Role:      "USER",
			Status:    "active",
			CreatedAt: time.Now().AddDate(0, -6, 0),
			LastLogin: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:        "U002",
			Email:     "jane@example.com",
			Name:      "Jane Smith",
			Role:      "USER",
			Status:    "active",
			CreatedAt: time.Now().AddDate(0, -3, 0),
			LastLogin: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:        "U003",
			Email:     "bob@example.com",
			Name:      "Bob Wilson",
			Role:      "ADMIN",
			Status:    "locked",
			CreatedAt: time.Now().AddDate(0, -1, 0),
			LastLogin: time.Now().Add(-7 * 24 * time.Hour),
		},
	}

	return &dto.AdminUserListAllResponse{
		Success: true,
		Users:   users,
		Total:   245, // Mock total
		Page:    page,
		Limit:   limit,
	}, nil
}

func (s *adminUserService) GetUserDetails(ctx context.Context, userID string) (*dto.AdminUserDetailsResponse, error) {
	log.Info().Str("user_id", userID).Msg("Fetching user details")

	// Mock user details
	user := &dto.AdminUserDetailsResponse{
		ID:              userID,
		Email:           "john@example.com",
		Name:            "John Doe",
		PhoneNo:         "+1234567890",
		NIC:             "12345-1234567-1",
		Role:            "USER",
		Status:          "active",
		MFAEnabled:      true,
		WebAuthnEnabled: true,
		CreatedAt:       time.Now().AddDate(0, -6, 0),
		UpdatedAt:       time.Now().Add(-2 * 24 * time.Hour),
		LastLogin:       time.Now().Add(-1 * time.Hour),
		LoginAttempts:   2,
		SecurityLevel:   "HIGH",
	}

	return user, nil
}

func (s *adminUserService) UpdateLockStatus(ctx context.Context, userID string, req dto.UpdateLockStatusRequest) (*dto.UpdateLockStatusResponse, error) {
	log.Info().Str("user_id", userID).Bool("locked", req.Locked).Msg("Updating user lock status")

	status := "unlocked"
	if req.Locked {
		status = "locked"
	}

	return &dto.UpdateLockStatusResponse{
		Success: true,
		Message: fmt.Sprintf("user successfully %s", status),
		Status:  status,
	}, nil
}

func (s *adminUserService) DeleteUser(ctx context.Context, userID string, req dto.AdminDeleteUserRequest) (*dto.AdminDeleteUserResponse, error) {
	log.Info().Str("user_id", userID).Str("reason", req.Reason).Msg("Deleting user")

	return &dto.AdminDeleteUserResponse{
		Success: true,
		Message: fmt.Sprintf("user %s has been deleted", userID),
	}, nil
}
