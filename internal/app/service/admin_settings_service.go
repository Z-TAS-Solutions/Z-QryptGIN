package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/rs/zerolog/log"
)

type AdminSettingsService interface {
	GetSettings(ctx context.Context) (*dto.AdminSettingsResponse, error)
	UpdateSettings(ctx context.Context, req dto.UpdateSettingsRequest) (*dto.UpdateSettingsResponse, error)
	EnforceMFA(ctx context.Context, req dto.AdminMFAEnforcementRequest) (*dto.AdminMFAEnforcementResponse, error)
}

type adminSettingsService struct {
}

func NewAdminSettingsService() AdminSettingsService {
	return &adminSettingsService{}
}

func (s *adminSettingsService) GetSettings(ctx context.Context) (*dto.AdminSettingsResponse, error) {
	log.Info().Msg("Fetching admin settings")

	settings := &dto.AdminSettingsResponse{
		Success: true,
		SecurityPolicy: dto.SecurityPolicy{
			MFARequired:            true,
			WebAuthnRequired:       true,
			PasswordExpiryDays:     90,
			MaxLoginAttempts:       5,
			LockoutDurationMinutes: 15,
			SessionTimeoutMinutes:  30,
			AllowWeakPasswords:     false,
			require2FAForAdmin:     true,
		},
		NotificationSettings: dto.NotificationSettings{
			EmailOnNewLogin:        true,
			EmailOnSecurityChange:  true,
			EmailOnAccountActivity: true,
			SMSNotifications:       false,
		},
		AppVersion:  "v1.0.0",
		LastUpdated: time.Now().Add(-7 * 24 * time.Hour).Format(time.RFC3339),
	}

	return settings, nil
}

func (s *adminSettingsService) UpdateSettings(ctx context.Context, req dto.UpdateSettingsRequest) (*dto.UpdateSettingsResponse, error) {
	log.Info().Msg("Updating admin settings")

	if req.SecurityPolicy == nil && req.NotificationSettings == nil {
		return &dto.UpdateSettingsResponse{
			Success: false,
			Message: "no settings to update",
		}, fmt.Errorf("no settings provided")
	}

	return &dto.UpdateSettingsResponse{
		Success: true,
		Message: "settings updated successfully",
	}, nil
}

func (s *adminSettingsService) EnforceMFA(ctx context.Context, req dto.AdminMFAEnforcementRequest) (*dto.AdminMFAEnforcementResponse, error) {
	log.Info().Bool("enabled", req.Enabled).Str("reason", req.Reason).Msg("Updating MFA enforcement")

	status := "disabled"
	if req.Enabled {
		status = "enabled"
	}

	return &dto.AdminMFAEnforcementResponse{
		Success: true,
		Message: fmt.Sprintf("MFA enforcement %s", status),
		Status:  req.Enabled,
	}, nil
}
