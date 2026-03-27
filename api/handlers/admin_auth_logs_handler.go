package handlers

import (
	"net/http"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
	"github.com/gin-gonic/gin"
)

type AdminAuthLogsHandler struct {
	svc service.AdminAuthLogsService
}

func NewAdminAuthLogsHandler(svc service.AdminAuthLogsService) *AdminAuthLogsHandler {
	return &AdminAuthLogsHandler{svc: svc}
}

func (h *AdminAuthLogsHandler) GetAuthLogs(c *gin.Context) {
	var req dto.AuthLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.svc.GetAuthLogs(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AdminAuthLogsHandler) GetAuthAnalytics(c *gin.Context) {
	var req dto.AuthAnalyticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.svc.GetAuthAnalytics(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
