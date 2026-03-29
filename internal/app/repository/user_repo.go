package repository

import (
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *database.User) error
	FindByID(id uint) (*database.User, error)
	FindByEmail(email string) (*database.User, error)
	FindByPhoneNo(phone string) (*database.User, error)
	FindByNic(nic string) (*database.User, error)
	FindByCustomID(customID string) (*database.User, error)
	GetPaginatedUsers(limit, offset int, search, status string, mfaEnabled *bool, lastLoginFrom, lastLoginTo *time.Time, sortBy, order string) ([]database.User, int64, error)
	UpdateLastLogin(userID uint, lastLogin time.Time) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *database.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*database.User, error) {
	var user database.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByEmail(email string) (*database.User, error) {
	var user database.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByPhoneNo(phone string) (*database.User, error) {
	var user database.User
	err := r.db.Where("phone_no = ?", phone).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByNic(nic string) (*database.User, error) {
	var user database.User
	err := r.db.Where("nic = ?", nic).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByCustomID(customID string) (*database.User, error) {
	var user database.User
	err := r.db.Where("custom_id = ?", customID).First(&user).Error
	return &user, err
}

func (r *userRepository) GetPaginatedUsers(limit, offset int, search, status string, mfaEnabled *bool, lastLoginFrom, lastLoginTo *time.Time, sortBy, order string) ([]database.User, int64, error) {
	var users []database.User
	var total int64

	query := r.db.Model(&database.User{})

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("LOWER(name) LIKE LOWER(?) OR LOWER(email) LIKE LOWER(?) OR phone_no LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if mfaEnabled != nil {
		query = query.Where("mfa_status = ?", *mfaEnabled)
	}

	if lastLoginFrom != nil {
		query = query.Where("last_login >= ?", *lastLoginFrom)
	}

	if lastLoginTo != nil {
		query = query.Where("last_login <= ?", *lastLoginTo)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	sortField := "created_at"
	switch sortBy {
	case "lastLogin":
		sortField = "last_login"
	case "name":
		sortField = "name"
	case "email":
		sortField = "email"
	}

	if order == "asc" {
		query = query.Order(sortField + " asc")
	} else {
		query = query.Order(sortField + " desc")
	}

	err = query.Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) UpdateLastLogin(userID uint, lastLogin time.Time) error {
	return r.db.Model(&database.User{}).Where("id = ?", userID).Update("last_login", lastLogin).Error
}
