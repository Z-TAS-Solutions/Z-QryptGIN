package handlers

import (
	"net/http"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
	"github.com/gin-gonic/gin"
)

type UserSessionsHandler struct {
	svc service.UserSessionsService
}

func NewUserSessionsHandler(svc service.UserSessionsService) *UserSessionsHandler {
	return &UserSessionsHandler{svc: svc}
}

func (h *UserSessionsHandler) FetchActiveSessions(c *gin.Context) {
	userID := "user_123"

	var req dto.FetchSessionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.svc.FetchActiveSessions(userID, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserSessionsHandler) SignOutOtherDevices(c *gin.Context) {
	userID := "user_123"
	currentSessionID := "sess_001" // normally resolved via middleware

	resp, err := h.svc.SignOutOtherDevices(userID, currentSessionID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
