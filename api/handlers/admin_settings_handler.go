package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
)

type AdminSettingsHandler struct {
	svc service.AdminSettingsService
}

func NewAdminSettingsHandler(svc service.AdminSettingsService) *AdminSettingsHandler {
	return &AdminSettingsHandler{svc: svc}
}

func (h *AdminSettingsHandler) GetSettings(c *gin.Context) {
	settings, err := h.svc.GetSettings(c.Request.Context())
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to fetch settings",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   settings,
	})
}

func (h *AdminSettingsHandler) UpdateSettings(c *gin.Context) {
	var req dto.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid JSON payload: " + err.Error(),
		})
		return
	}

	resp, err := h.svc.UpdateSettings(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": resp.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   resp,
	})
}

func (h *AdminSettingsHandler) EnforceMFA(c *gin.Context) {
	var req dto.AdminMFAEnforcementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid JSON payload: " + err.Error(),
		})
		return
	}

	resp, err := h.svc.EnforceMFA(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to enforce MFA",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   resp,
	})
}
