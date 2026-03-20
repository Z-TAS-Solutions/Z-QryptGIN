package service

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
)

type UserService interface {
	RegisterUser(req dto.CreateUserRequest) (*dto.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// 1. Hash the password (using bcrypt/argon2)
	// 2. Map DTO to GORM Model
	user := &database.User{
		Name:         req.Name,
		Email:        database.Email(req.Email),
		PasswordHash: "hashed_password_here", 
		// Set other fields...
	}

	// 3. Save via Repository
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// 4. Map GORM Model back to Response DTO
	return &dto.UserResponse{
		Name:  user.Name,
		Email: string(user.Email),
	}, nil
}