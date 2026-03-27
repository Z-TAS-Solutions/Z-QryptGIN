package service

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
)

type AdminAuthService interface {
	Login(req dto.AdminLoginRequest) (*dto.AdminLoginResponse, error)
	Refresh(req dto.AdminRefreshRequest) (*dto.AdminRefreshResponse, error)
}

type adminAuthService struct {
	repo repository.AdminAuthRepository
}

func NewAdminAuthService(repo repository.AdminAuthRepository) AdminAuthService {
	return &adminAuthService{repo: repo}
}

func (s *adminAuthService) Login(req dto.AdminLoginRequest) (*dto.AdminLoginResponse, error) {
	return s.repo.Login(req)
}

func (s *adminAuthService) Refresh(req dto.AdminRefreshRequest) (*dto.AdminRefreshResponse, error) {
	return s.repo.Refresh(req)
}
