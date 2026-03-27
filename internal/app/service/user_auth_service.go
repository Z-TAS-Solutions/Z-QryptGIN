package service

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
)

type UserAuthService interface {
	RegisterOptions(req dto.PasskeyRegisterOptionsRequest) (*dto.PasskeyRegisterOptionsResponse, error)
	RegisterVerify(req dto.PasskeyRegisterVerifyRequest) (*dto.PasskeyRegisterVerifyResponse, error)
	LoginOptions(req dto.PasskeyLoginOptionsRequest) (*dto.PasskeyLoginOptionsResponse, error)
	LoginVerify(req dto.PasskeyLoginVerifyRequest) (*dto.PasskeyLoginVerifyResponse, error)
}

type userAuthService struct {
	repo repository.UserAuthRepository
}

func NewUserAuthService(repo repository.UserAuthRepository) UserAuthService {
	return &userAuthService{repo: repo}
}

func (s *userAuthService) RegisterOptions(req dto.PasskeyRegisterOptionsRequest) (*dto.PasskeyRegisterOptionsResponse, error) {
	return s.repo.RegisterOptions(req)
}

func (s *userAuthService) RegisterVerify(req dto.PasskeyRegisterVerifyRequest) (*dto.PasskeyRegisterVerifyResponse, error) {
	return s.repo.RegisterVerify(req)
}

func (s *userAuthService) LoginOptions(req dto.PasskeyLoginOptionsRequest) (*dto.PasskeyLoginOptionsResponse, error) {
	return s.repo.LoginOptions(req)
}

func (s *userAuthService) LoginVerify(req dto.PasskeyLoginVerifyRequest) (*dto.PasskeyLoginVerifyResponse, error) {
	return s.repo.LoginVerify(req)
}
