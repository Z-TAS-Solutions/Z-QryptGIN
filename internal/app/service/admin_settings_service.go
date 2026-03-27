package service

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
)

type AdminSettingsService interface {
	GetSettings() (*dto.SystemSettingsResponse, error)
}

type adminSettingsService struct {
	repo repository.AdminSettingsRepository
}

func NewAdminSettingsService(repo repository.AdminSettingsRepository) AdminSettingsService {
	return &adminSettingsService{repo: repo}
}

func (s *adminSettingsService) GetSettings() (*dto.SystemSettingsResponse, error) {
	return s.repo.GetSettings()
}
