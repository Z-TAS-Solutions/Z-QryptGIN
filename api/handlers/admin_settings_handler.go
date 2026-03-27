package handlers

import (
	"net/http"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
	"github.com/gin-gonic/gin"
)

type AdminSettingsHandler struct {
	svc service.AdminSettingsService
}

func NewAdminSettingsHandler(svc service.AdminSettingsService) *AdminSettingsHandler {
	return &AdminSettingsHandler{svc: svc}
}

func (h *AdminSettingsHandler) GetSettings(c *gin.Context) {
	resp, err := h.svc.GetSettings()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
