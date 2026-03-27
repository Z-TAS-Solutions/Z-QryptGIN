package handlers

import (
	"net/http"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
	"github.com/gin-gonic/gin"
)

type UserNotificationsHandler struct {
	svc service.UserNotificationsService
}

func NewUserNotificationsHandler(svc service.UserNotificationsService) *UserNotificationsHandler {
	return &UserNotificationsHandler{svc: svc}
}

func (h *UserNotificationsHandler) FetchNotifications(c *gin.Context) {
	userID := "user_123" // normally from token context

	var req dto.FetchNotificationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.svc.FetchNotifications(userID, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserNotificationsHandler) UpdateStatus(c *gin.Context) {
	userID := "user_123" // token context
	notificationID := c.Param("notificationId")

	var req dto.UpdateNotificationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.svc.UpdateStatus(userID, notificationID, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserNotificationsHandler) MarkAllRead(c *gin.Context) {
	userID := "user_123" // token

	resp, err := h.svc.MarkAllRead(userID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
