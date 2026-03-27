package handlers

import (
	"net/http"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
	"github.com/gin-gonic/gin"
)

type AdminDeviceHandler struct {
	svc service.AdminDeviceService
}

func NewAdminDeviceHandler(svc service.AdminDeviceService) *AdminDeviceHandler {
	return &AdminDeviceHandler{svc: svc}
}

func (h *AdminDeviceHandler) ListDevices(c *gin.Context) {
	var req dto.ListDevicesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.svc.ListDevices(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AdminDeviceHandler) ForceLogout(c *gin.Context) {
	deviceID := c.Param("deviceId")

	resp, err := h.svc.ForceLogout(deviceID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
