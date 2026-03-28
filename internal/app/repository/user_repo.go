package repository

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *database.User) error
	FindByEmail(email string) (*database.User, error)
	FindByPhoneNo(phone string) (*database.User, error)
	FindByNic(nic string) (*database.User, error)
	FindByCustomID(customID string) (*database.User, error)
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
