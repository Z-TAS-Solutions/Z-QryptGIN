package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
)

type AdminSecurityHandler struct {
	svc service.AdminSettingsService
}

func NewAdminSecurityHandler(svc service.AdminSettingsService) *AdminSecurityHandler {
	return &AdminSecurityHandler{svc: svc}
}

func (h *AdminSecurityHandler) EnforceMFA(c *gin.Context) {
	// Delegating to admin settings handler - enforce MFA method
	handler := NewAdminSettingsHandler(h.svc)
	handler.EnforceMFA(c)
}
