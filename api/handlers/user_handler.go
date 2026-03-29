package handlers

import (
	"net/http"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	sessionSvc service.SessionService
}

func NewUserHandler(sessionSvc service.SessionService) *UserHandler {
	return &UserHandler{
		sessionSvc: sessionSvc,
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
