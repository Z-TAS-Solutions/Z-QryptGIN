package service

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
)

type UserNotificationsService interface {
	FetchNotifications(userID string, req dto.FetchNotificationsRequest) (*dto.FetchNotificationsResponse, error)
	UpdateStatus(userID string, notificationID string, req dto.UpdateNotificationStatusRequest) (*dto.UpdateNotificationStatusResponse, error)
	MarkAllRead(userID string) (*dto.MarkAllNotificationsReadResponse, error)
}

type userNotificationsService struct {
	repo repository.UserNotificationsRepository
}

func NewUserNotificationsService(repo repository.UserNotificationsRepository) UserNotificationsService {
	return &userNotificationsService{repo: repo}
}

func (s *userNotificationsService) FetchNotifications(userID string, req dto.FetchNotificationsRequest) (*dto.FetchNotificationsResponse, error) {
	if req.Limit < 1 {
		req.Limit = 20
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}
	return s.repo.FetchNotifications(userID, req)
}

func (s *userNotificationsService) UpdateStatus(userID string, notificationID string, req dto.UpdateNotificationStatusRequest) (*dto.UpdateNotificationStatusResponse, error) {
	return s.repo.UpdateStatus(userID, notificationID, req.Status)
}

func (s *userNotificationsService) MarkAllRead(userID string) (*dto.MarkAllNotificationsReadResponse, error) {
	return s.repo.MarkAllRead(userID)
}
