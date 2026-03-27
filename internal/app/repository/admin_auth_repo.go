package repository

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"gorm.io/gorm"
)

type AdminAuthRepository interface {
	Login(req dto.AdminLoginRequest) (*dto.AdminLoginResponse, error)
	Refresh(req dto.AdminRefreshRequest) (*dto.AdminRefreshResponse, error)
}

type adminAuthRepository struct {
	db *gorm.DB
}

func NewAdminAuthRepository(db *gorm.DB) AdminAuthRepository {
	return &adminAuthRepository{db: db}
}

func (r *adminAuthRepository) Login(req dto.AdminLoginRequest) (*dto.AdminLoginResponse, error) {
	var res dto.AdminLoginResponse
	res.Message = "Login successful"
	res.Data.AccessToken = "eyJ..."
	res.Data.RefreshToken = "def502..."
	res.Data.AdminID = "admin_001"
	res.Data.Role = "super_admin"
	return &res, nil
}

func (r *adminAuthRepository) Refresh(req dto.AdminRefreshRequest) (*dto.AdminRefreshResponse, error) {
	var res dto.AdminRefreshResponse
	res.Message = "Token refreshed successfully"
	res.Data.AccessToken = "eyJ..."
	return &res, nil
}
