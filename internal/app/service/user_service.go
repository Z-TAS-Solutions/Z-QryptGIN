package service

import (
	"crypto/rand"
	"io"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
	"github.com/rs/zerolog/log"
)

type UserService interface {
	RegisterUser(req dto.CreateUserRequest) (*dto.UserResponse, error)
}

type userService struct {
	repo    repository.UserRepository
	Session repository.SessionRepository
	email   EmailService
}

func NewUserService(
	repo repository.UserRepository,
	session repository.SessionRepository,
	email EmailService,
) UserService {
	return &userService{
		repo:    repo,
		Session: session,
		email:   email,
	}
}

func (s *userService) RegisterUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// 1. TODO: Hash the password using argon2id (Gotta do this before deploying anything)
	hashedPassword := "hashed_password_placeholder"

	// Map DTO to GORM Model
	user := &database.User{
		Name:         req.Name,
		Email:        database.Email(req.Email),
		PasswordHash: hashedPassword,
		Status:       database.StatusActive,
		Role:         database.RoleClient,
	}

	// Save via Repository
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	otpCode := generateOTP(6)

	err := s.email.SendOTPEmail(string(user.Email), user.Name, otpCode)
	if err != nil {
		log.Error().Err(err).Str("user_email", string(user.Email)).Msg("Failed to send registration OTP")
	} else {
		log.Info().Str("user_email", string(user.Email)).Msg("Registration OTP sent Successfully")
	}

	// Map GORM Model back to Response DTO
	return &dto.UserResponse{
		Success: true,
		Name:    user.Name,
		Email:   string(user.Email),
	}, nil
}

func generateOTP(max int) string {
	var table = []byte{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

	b := make([]byte, max)

	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max || err != nil {
		return "0Aj3kL"
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
