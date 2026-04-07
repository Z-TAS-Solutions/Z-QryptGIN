package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/rs/zerolog/log"
)

type AdminDeviceService interface {
	ListDevices(ctx context.Context, page, limit int) (*dto.AdminListDevicesResponse, error)
	ForceLogout(ctx context.Context, deviceID string, req dto.AdminForceLogoutRequest) (*dto.AdminForceLogoutResponse, error)
}

type adminDeviceService struct {
}

func NewAdminDeviceService() AdminDeviceService {
	return &adminDeviceService{}
}

func (s *adminDeviceService) ListDevices(ctx context.Context, page, limit int) (*dto.AdminListDevicesResponse, error) {
	log.Info().Int("page", page).Int("limit", limit).Msg("Listing devices")

	devices := []dto.AdminDeviceResponse{
		{
			DeviceID:     "dev-001",
			UserID:       "U001",
			UserEmail:    "john@example.com",
			DeviceName:   "Chrome on Windows",
			UserAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			IPAddress:    "192.168.1.100",
			LastActivity: time.Now().Add(-5 * time.Minute),
			CreatedAt:    time.Now().AddDate(0, 0, -2),
			Status:       "active",
		},
		{
			DeviceID:     "dev-002",
			UserID:       "U001",
			UserEmail:    "john@example.com",
			DeviceName:   "Safari on iPhone",
			UserAgent:    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0)",
			IPAddress:    "203.0.113.45",
			LastActivity: time.Now().Add(-2 * time.Hour),
			CreatedAt:    time.Now().AddDate(0, 0, -10),
			Status:       "active",
		},
		{
			DeviceID:     "dev-003",
			UserID:       "U002",
			UserEmail:    "jane@example.com",
			DeviceName:   "Firefox on Mac",
			UserAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
			IPAddress:    "198.51.100.22",
			LastActivity: time.Now().Add(-1 * time.Day),
			CreatedAt:    time.Now().AddDate(0, 0, -5),
			Status:       "inactive",
		},
	}

	return &dto.AdminListDevicesResponse{
		Success: true,
		Devices: devices,
		Total:   892, // Mock total
		Page:    page,
		Limit:   limit,
	}, nil
}

func (s *adminDeviceService) ForceLogout(ctx context.Context, deviceID string, req dto.AdminForceLogoutRequest) (*dto.AdminForceLogoutResponse, error) {
	log.Info().Str("device_id", deviceID).Str("reason", req.Reason).Msg("Force logging out device")

	return &dto.AdminForceLogoutResponse{
		Success:  true,
		Message:  fmt.Sprintf("device %s has been forcefully logged out", deviceID),
		DeviceID: deviceID,
	}, nil
}
