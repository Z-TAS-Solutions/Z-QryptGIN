package handlers

import (
	"net/http"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	sessionSvc      service.SessionService
	notificationSvc service.NotificationService
}

func NewUserHandler(sessionSvc service.SessionService, notificationSvc service.NotificationService) *UserHandler {
	return &UserHandler{
		sessionSvc:      sessionSvc,
		notificationSvc: notificationSvc,
	}
}

// GetActiveSessions retrieves all active sessions for the authenticated user
// It extracts the user_id from the JWT claims (via context) and fetches active sessions from Redis
func (h *UserHandler) GetActiveSessions(c *gin.Context) {
	// Extract user_id from context (set by RequireAuth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "user_id not found in context",
		})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "invalid user_id type in context",
		})
		return
	}

	// Fetch active sessions for the user
	response, err := h.sessionSvc.GetActiveSessionsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}

// GetNotifications retrieves notifications for the authenticated user with filtering and pagination
// It extracts the user_id from JWT context and supports limit, offset, unread filtering, and sorting
func (h *UserHandler) GetNotifications(c *gin.Context) {
	// Extract user_id from context (set by RequireAuth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Authentication required",
		})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": "Failed to retrieve notifications",
		})
		return
	}

	// Parse and validate query parameters
	var req dto.FetchNotificationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid query parameters",
			"details": []map[string]string{
				{
					"field":   "limit",
					"message": "Must be between 1 and 100",
				},
			},
		})
		return
	}

	// Fetch notifications from service
	response, err := h.notificationSvc.GetNotificationsByUserID(userID, req.Limit, req.Offset, req.UnreadOnly, req.SortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": "Failed to retrieve notifications",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateNotificationStatus updates the read status of a specific notification
// It extracts the user_id from JWT context and verifies ownership before updating
func (h *UserHandler) UpdateNotificationStatus(c *gin.Context) {
	// Extract user_id from context (set by RequireAuth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Authentication required",
		})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": "Failed to update notification status",
		})
		return
	}

	// Extract notification ID from URL path
	notificationID := c.Param("notificationId")
	if notificationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": []map[string]string{
				{
					"field":   "notificationId",
					"message": "Notification ID is required",
				},
			},
		})
		return
	}

	// Parse request body
	var req dto.UpdateNotificationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": []map[string]string{
				{
					"field":   "status",
					"message": "Must be either 'read' or 'unread'",
				},
			},
		})
		return
	}

	// Update notification status
	response, err := h.notificationSvc.UpdateNotificationStatus(userID, notificationID, req.Status)
	if err != nil {
		// Check if it's a "not found" error
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Notification not found",
				"message": "No notification exists with the given ID",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": "Failed to update notification status",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// MarkAllAsRead marks all notifications for the authenticated user as read
// It extracts the user_id from JWT context and performs a bulk update operation
func (h *UserHandler) MarkAllAsRead(c *gin.Context) {
	// Extract user_id from context (set by RequireAuth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Authentication required",
		})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": "Failed to update notifications",
		})
		return
	}

	// Mark all notifications as read
	response, err := h.notificationSvc.MarkAllAsRead(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": "Failed to update notifications",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
