package repository

import (
	"fmt"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"gorm.io/gorm"
)

type AdminUserRepository interface {
	ListUsers(req dto.ListUsersRequest) (*dto.ListUsersResponse, error)
	GetUser(userID string) (*dto.UserDetailsResponse, error)
	UpdateLockStatus(userID string, locked bool) (*dto.LockUserResponse, error)
	DeleteUser(userID string) (*dto.DeleteUserResponse, error)
}

type adminUserRepository struct {
	db *gorm.DB
}

func NewAdminUserRepository(db *gorm.DB) AdminUserRepository {
	return &adminUserRepository{db: db}
}

func (r *adminUserRepository) ListUsers(req dto.ListUsersRequest) (*dto.ListUsersResponse, error) {
	var users []database.User
	var total int64

	query := r.db.Model(&database.User{})

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	limit := req.Limit
	if limit == 0 {
		limit = 10
	}
	offset := req.Offset

	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}

	res := &dto.ListUsersResponse{
		Limit:  limit,
		Offset: offset,
		Total:  total,
		Data:   []dto.UserListItem{},
	}

	for _, u := range users {
		res.Data = append(res.Data, dto.UserListItem{
			UserID:     fmt.Sprint(u.ID),
			Name:       u.Name,
			Email:      string(u.Email),
			Phone:      string(u.PhoneNo),
			MfaEnabled: u.MFAStatus,
			Status:     string(u.Status),
		})
	}

	return res, nil
}

func (r *adminUserRepository) GetUser(userID string) (*dto.UserDetailsResponse, error) {
	var user database.User

	err := r.db.Preload("Passkeys").Preload("Sessions").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}

	res := &dto.UserDetailsResponse{
		Message: "User details fetched successfully",
		Data: dto.UserDetails{
			UserID:       userID,
			Name:         user.Name,
			Email:        string(user.Email),
			Phone:        string(user.PhoneNo),
			MfaEnabled:   user.MFAStatus,
			Status:       string(user.Status),
			RegisteredAt: user.CreatedAt,
		},
	}
	return res, nil
}

func (r *adminUserRepository) UpdateLockStatus(userID string, locked bool) (*dto.LockUserResponse, error) {
	status := "Active"
	if locked {
		status = "Locked"
	}
	err := r.db.Model(&database.User{}).Where("id = ?", userID).Update("status", status).Error
	if err != nil {
		return nil, err
	}
	res := &dto.LockUserResponse{
		Message: "User lock status updated successfully",
	}
	res.Data.UserID = userID
	res.Data.Locked = locked
	return res, nil
}

func (r *adminUserRepository) DeleteUser(userID string) (*dto.DeleteUserResponse, error) {
	err := r.db.Where("id = ?", userID).Delete(&database.User{}).Error
	if err != nil {
		return nil, err
	}
	res := &dto.DeleteUserResponse{
		Message: "User deleted successfully",
	}
	res.Data.UserID = userID
	res.Data.Deleted = true
	return res, nil
}
