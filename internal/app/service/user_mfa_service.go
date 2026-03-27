package service

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
)

type UserMfaService interface {
	Send(req dto.MfaSendRequest) (*dto.MfaSendResponse, error)
	Respond(req dto.MfaRespondRequest) (*dto.MfaRespondResponse, error)
}

type userMfaService struct {
	repo repository.UserMfaRepository
}

func NewUserMfaService(repo repository.UserMfaRepository) UserMfaService {
	return &userMfaService{repo: repo}
}

func (s *userMfaService) Send(req dto.MfaSendRequest) (*dto.MfaSendResponse, error) {
	return s.repo.Send(req)
}

func (s *userMfaService) Respond(req dto.MfaRespondRequest) (*dto.MfaRespondResponse, error) {
	return s.repo.Respond(req)
}
