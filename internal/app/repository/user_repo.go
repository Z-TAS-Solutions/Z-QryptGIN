package repository

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *database.User) error
	FindByEmail(email string) (*database.User, error)
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