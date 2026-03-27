package service

import (
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
)

type AdminSecurityService interface {
	EnforceMfa(req dto.EnforceMfaRequest) (*dto.EnforceMfaResponse, error)
}

type adminSecurityService struct {
	repo repository.AdminSecurityRepository
}

func NewAdminSecurityService(repo repository.AdminSecurityRepository) AdminSecurityService {
	return &adminSecurityService{repo: repo}
}

func (s *adminSecurityService) EnforceMfa(req dto.EnforceMfaRequest) (*dto.EnforceMfaResponse, error) {
	return s.repo.EnforceMfa(req.Enabled)
}
