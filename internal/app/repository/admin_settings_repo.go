package repository

import (
	"strings"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"gorm.io/gorm"
)

type AdminSettingsRepository interface {
	GetSettings() (*dto.SystemSettingsResponse, error)
}

type adminSettingsRepository struct {
	db *gorm.DB
}

func NewAdminSettingsRepository(db *gorm.DB) AdminSettingsRepository {
	return &adminSettingsRepository{db: db}
}

func (r *adminSettingsRepository) GetSettings() (*dto.SystemSettingsResponse, error) {
	var config database.SystemConfiguration
	if err := r.db.First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Seed default configuration if none exists
			config = database.SystemConfiguration{
				SessionTimeoutMinutes:  30,
				MaxFailedLoginAttempts: 5,
				AllowedAdminIpRanges:   "192.168.1.0/24",
				EnforceMfa:             false,
			}
			r.db.Create(&config)
		} else {
			return nil, err
		}
	}

	var res dto.SystemSettingsResponse
	res.Message = "System settings retrieved successfully"
	res.Data.SessionTimeoutMinutes = config.SessionTimeoutMinutes
	res.Data.MaxFailedLoginAttempts = config.MaxFailedLoginAttempts

	if config.AllowedAdminIpRanges != "" {
		res.Data.AllowedAdminIpRanges = strings.Split(config.AllowedAdminIpRanges, ",")
	} else {
		res.Data.AllowedAdminIpRanges = []string{}
	}

	return &res, nil
}
