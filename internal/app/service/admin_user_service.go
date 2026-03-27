package service

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
)

type AdminUserService interface {
	ListUsers(req dto.ListUsersRequest) (*dto.ListUsersResponse, error)
	GetUser(userID string) (*dto.UserDetailsResponse, error)
	UpdateLockStatus(userID string, req dto.LockUserRequest) (*dto.LockUserResponse, error)
	DeleteUser(userID string) (*dto.DeleteUserResponse, error)
}

type adminUserService struct {
	repo repository.AdminUserRepository
}

func NewAdminUserService(repo repository.AdminUserRepository) AdminUserService {
	return &adminUserService{repo: repo}
}

func (s *adminUserService) ListUsers(req dto.ListUsersRequest) (*dto.ListUsersResponse, error) {
	if req.Limit < 1 {
		req.Limit = 20
	}
	if req.SortBy == "" {
		req.SortBy = "lastLogin"
	}
	if req.Order == "" {
		req.Order = "desc"
	}
	return s.repo.ListUsers(req)
}

func (s *adminUserService) GetUser(userID string) (*dto.UserDetailsResponse, error) {
	return s.repo.GetUser(userID)
}

func (s *adminUserService) UpdateLockStatus(userID string, req dto.LockUserRequest) (*dto.LockUserResponse, error) {
	return s.repo.UpdateLockStatus(userID, req.Locked)
}

func (s *adminUserService) DeleteUser(userID string) (*dto.DeleteUserResponse, error) {
	return s.repo.DeleteUser(userID)
}
