package service

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
)

type UserAccountService interface {
	ForceLogoutAllDevices(userID string) (*dto.ForceLogoutUserDevicesResponse, error)
	DeleteAccount(userID string, req dto.DeleteAccountRequest) (*dto.DeleteAccountResponse, error)
}

type userAccountService struct {
	repo repository.UserAccountRepository
}

func NewUserAccountService(repo repository.UserAccountRepository) UserAccountService {
	return &userAccountService{repo: repo}
}

func (s *userAccountService) ForceLogoutAllDevices(userID string) (*dto.ForceLogoutUserDevicesResponse, error) {
	return s.repo.ForceLogoutAllDevices(userID)
}

func (s *userAccountService) DeleteAccount(userID string, req dto.DeleteAccountRequest) (*dto.DeleteAccountResponse, error) {
	return s.repo.DeleteAccount(userID, req)
}
