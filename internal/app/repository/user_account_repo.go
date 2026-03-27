package repository

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"gorm.io/gorm"
)

type UserAccountRepository interface {
	ForceLogoutAllDevices(userID string) (*dto.ForceLogoutUserDevicesResponse, error)
	DeleteAccount(userID string, req dto.DeleteAccountRequest) (*dto.DeleteAccountResponse, error)
}

type userAccountRepository struct {
	db *gorm.DB
}

func NewUserAccountRepository(db *gorm.DB) UserAccountRepository {
	return &userAccountRepository{db: db}
}

func (r *userAccountRepository) ForceLogoutAllDevices(userID string) (*dto.ForceLogoutUserDevicesResponse, error) {
	var res dto.ForceLogoutUserDevicesResponse
	// Real GORM logic
	resTx := r.db.Where("user_id = ?", userID).Delete(&database.Session{})
	if resTx.Error != nil {
		return nil, resTx.Error
	}
	res.Message = "All user devices logged out successfully"
	res.Data.UserID = userID
	res.Data.DevicesTerminated = int(resTx.RowsAffected)
	return &res, nil
}

func (r *userAccountRepository) DeleteAccount(userID string, req dto.DeleteAccountRequest) (*dto.DeleteAccountResponse, error) {
	var res dto.DeleteAccountResponse
	// Real GORM logic
	err := r.db.Where("id = ?", userID).Delete(&database.User{}).Error
	if err != nil {
		return nil, err
	}
	res.Message = "Account deleted successfully"
	return &res, nil
}
