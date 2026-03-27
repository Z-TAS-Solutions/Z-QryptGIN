package handlers

import (
	"net/http"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
	"github.com/gin-gonic/gin"
)

type UserAccountHandler struct {
	svc service.UserAccountService
}

func NewUserAccountHandler(svc service.UserAccountService) *UserAccountHandler {
	return &UserAccountHandler{svc: svc}
}

func (h *UserAccountHandler) ForceLogoutAllDevices(c *gin.Context) {
	// From admin or User
	var req dto.ForceLogoutUserDevicesRequest
	// This could come from query or body depending on spec.
	userID := c.Query("userId")
	if userID == "" {
		if err := c.ShouldBindJSON(&req); err == nil {
			userID = req.UserID
		}
	}

	resp, err := h.svc.ForceLogoutAllDevices(userID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserAccountHandler) DeleteAccount(c *gin.Context) {
	userID := "user_123"

	var req dto.DeleteAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.svc.DeleteAccount(userID, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
