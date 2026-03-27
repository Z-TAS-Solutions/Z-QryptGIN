package handlers

import (
	"net/http"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
	"github.com/gin-gonic/gin"
)

type AdminSecurityHandler struct {
	svc service.AdminSecurityService
}

func NewAdminSecurityHandler(svc service.AdminSecurityService) *AdminSecurityHandler {
	return &AdminSecurityHandler{svc: svc}
}

func (h *AdminSecurityHandler) EnforceMfa(c *gin.Context) {
	var req dto.EnforceMfaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.svc.EnforceMfa(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
