package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
)

type NotificationHandler struct {
	svc service.NotificationService
}

func NewNotificationHandler(svc service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) FetchNotifications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "user not authenticated",
		})
		return
	}

	resp, err := h.svc.FetchNotifications(c.Request.Context(), userID.(string))
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to fetch notifications",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   resp,
	})
}

func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "user not authenticated",
		})
		return
	}

	resp, err := h.svc.MarkAllRead(c.Request.Context(), userID.(string))
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to mark notifications as read",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   resp,
	})
}

func (h *NotificationHandler) UpdateStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "user not authenticated",
		})
		return
	}

	notificationID := c.Param("notificationId")
	if notificationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "notification_id is required",
		})
		return
	}

	var req dto.UpdateNotificationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid JSON payload: " + err.Error(),
		})
		return
	}

	resp, err := h.svc.UpdateNotificationStatus(c.Request.Context(), userID.(string), notificationID, req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   resp,
	})
}
