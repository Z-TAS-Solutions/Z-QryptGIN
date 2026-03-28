package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
	"github.com/rs/zerolog/log"
)

type ProfileService interface {
	GetProfile(ctx context.Context, userID string) (*dto.UserProfileResponse, error)
	UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error)
}

type profileService struct {
	userRepo repository.UserRepository
}

func NewProfileService(userRepo repository.UserRepository) ProfileService {
	return &profileService{
		userRepo: userRepo,
	}
}

func (s *profileService) GetProfile(ctx context.Context, userID string) (*dto.UserProfileResponse, error) {
	// Just mock some data for now
	log.Info().Str("user_id", userID).Msg("Fetching user profile")
	
	profile := &dto.UserProfileResponse{
		UserID:    "U123",
		Name:      "John Doe",
		Email:     "john@example.com",
		PhoneNo:   "+1234567890",
		NIC:       "12345-1234567-1",
		Role:      "USER",
		CreatedAt: time.Now().AddDate(0, -6, 0).Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	return profile, nil
}

func (s *profileService) UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error) {
	log.Info().Str("user_id", userID).Msg("Updating user profile")
	
	// Simulate some validation
	if req.Name == nil && req.PhoneNo == nil {
		return &dto.UpdateProfileResponse{
			Success: false,
			Message: "no fields to update",
		}, fmt.Errorf("no fields provided for update")
	}

	// Mock updated profile
	profile := &dto.UserProfileResponse{
		UserID:    "U123",
		Name:      "John Doe",
		Email:     "john@example.com",
		PhoneNo:   "+1234567890",
		NIC:       "12345-1234567-1",
		Role:      "USER",
		CreatedAt: time.Now().AddDate(0, -6, 0).Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	return &dto.UpdateProfileResponse{
		Success: true,
		Message: "profile updated successfully",
		Data:    profile,
	}, nil
}
