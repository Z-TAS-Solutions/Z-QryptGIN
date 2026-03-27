package service

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
)

type AdminDeviceService interface {
	ListDevices(req dto.ListDevicesRequest) (*dto.ListDevicesResponse, error)
	ForceLogout(deviceID string) (*dto.ForceLogoutDeviceResponse, error)
}

type adminDeviceService struct {
	repo repository.AdminDeviceRepository
}

func NewAdminDeviceService(repo repository.AdminDeviceRepository) AdminDeviceService {
	return &adminDeviceService{repo: repo}
}

func (s *adminDeviceService) ListDevices(req dto.ListDevicesRequest) (*dto.ListDevicesResponse, error) {
	if req.Limit < 1 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.SortBy == "" {
		req.SortBy = "last_active"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}
	return s.repo.ListDevices(req)
}

func (s *adminDeviceService) ForceLogout(deviceID string) (*dto.ForceLogoutDeviceResponse, error) {
	// Add business logic to invalidate session inside repo
	return s.repo.ForceLogout(deviceID)
}
