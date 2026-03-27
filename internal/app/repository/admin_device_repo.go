package repository

import (
	"fmt"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"gorm.io/gorm"
)

type AdminDeviceRepository interface {
	ListDevices(req dto.ListDevicesRequest) (*dto.ListDevicesResponse, error)
	ForceLogout(deviceID string) (*dto.ForceLogoutDeviceResponse, error)
}

type adminDeviceRepository struct {
	db *gorm.DB
}

func NewAdminDeviceRepository(db *gorm.DB) AdminDeviceRepository {
	return &adminDeviceRepository{db: db}
}

func (r *adminDeviceRepository) ListDevices(req dto.ListDevicesRequest) (*dto.ListDevicesResponse, error) {
	var sessions []database.Session
	var total int64

	query := r.db.Model(&database.Session{})

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	limit := req.Limit
	if limit == 0 {
		limit = 10
	}
	offset := req.Offset

	if err := query.Limit(limit).Offset(offset).Find(&sessions).Error; err != nil {
		return nil, err
	}

	var res dto.ListDevicesResponse
	res.Message = "Devices retrieved successfully"
	res.Data.Devices = []dto.DeviceItem{}

	for _, s := range sessions {
		res.Data.Devices = append(res.Data.Devices, dto.DeviceItem{
			DeviceID:   fmt.Sprint(s.ID),
			DeviceName: s.DeviceName,
			Location:   s.Location,
			LastActive: s.LastActive.UnixMilli(),
		})
	}

	res.Data.Pagination.Limit = limit
	res.Data.Pagination.Offset = offset
	res.Data.Pagination.Returned = len(res.Data.Devices)
	res.Data.Pagination.HasMore = int64(offset+limit) < total

	return &res, nil
}

func (r *adminDeviceRepository) ForceLogout(deviceID string) (*dto.ForceLogoutDeviceResponse, error) {
	err := r.db.Where("id = ?", deviceID).Delete(&database.Session{}).Error
	if err != nil {
		return nil, err
	}

	var res dto.ForceLogoutDeviceResponse
	res.Message = "Device logged out successfully"
	res.Data.DeviceID = deviceID
	res.Data.LoggedOut = true
	return &res, nil
}
