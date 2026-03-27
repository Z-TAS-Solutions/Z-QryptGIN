package repository

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"gorm.io/gorm"
)

type AdminSecurityRepository interface {
	EnforceMfa(enabled bool) (*dto.EnforceMfaResponse, error)
}

type adminSecurityRepository struct {
	db *gorm.DB
}

func NewAdminSecurityRepository(db *gorm.DB) AdminSecurityRepository {
	return &adminSecurityRepository{db: db}
}

func (r *adminSecurityRepository) EnforceMfa(enabled bool) (*dto.EnforceMfaResponse, error) {
	var config database.SystemConfiguration

	// Get or Create Configuration Let's do first
	if err := r.db.First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			config = database.SystemConfiguration{
				SessionTimeoutMinutes:  30,
				MaxFailedLoginAttempts: 5,
				AllowedAdminIpRanges:   "192.168.1.0/24",
				EnforceMfa:             enabled,
			}
			r.db.Create(&config)
		} else {
			return nil, err
		}
	} else {
		// Update existing configuration
		err = r.db.Model(&config).Update("enforce_mfa", enabled).Error
		if err != nil {
			return nil, err
		}
	}

	var res dto.EnforceMfaResponse
	res.Message = "Two-factor authentication enforcement updated successfully"
	res.Data.Enabled = enabled
	return &res, nil
}
