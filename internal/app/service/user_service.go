package service

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
)

type UserService interface {
	ListUsers(req *dto.ListUsersRequest) (*dto.ListUsersResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// ListUsers queries the user repository and maps the DB models to ListUsersResponse DTO.
func (s *userService) ListUsers(req *dto.ListUsersRequest) (*dto.ListUsersResponse, error) {
	users, total, err := s.userRepo.GetPaginatedUsers(
		req.Limit,
		req.Offset,
		req.Search,
		req.Status,
		req.MfaEnabled,
		req.LastLoginFrom,
		req.LastLoginTo,
		req.SortBy,
		req.Order,
	)

	if err != nil {
		return nil, err
	}

	dataList := make([]dto.UserDashboardResponse, len(users))
	for i, u := range users {
		dataList[i] = dto.UserDashboardResponse{
			UserID:     string(u.CustomID),
			Name:       string(u.Name),
			Email:      string(u.Email),
			Phone:      string(u.PhoneNo),
			MfaEnabled: u.MFAStatus,
			LastLogin:  u.LastLogin,
			Status:     string(u.Status),
		}
	}

	response := &dto.ListUsersResponse{
		Limit:  req.Limit,
		Offset: req.Offset,
		Total:  total,
		Data:   dataList,
	}

	return response, nil
}
