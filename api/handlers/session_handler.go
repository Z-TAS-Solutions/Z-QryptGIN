package handlers

import (
	"net/http"
	"strconv"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionSvc service.SessionService
}

func NewSessionHandler(sessionSvc service.SessionService) *SessionHandler {
	return &SessionHandler{
		sessionSvc: sessionSvc,
	}
}

// GetActiveSessions retrieves all active sessions for the authenticated user with pagination
// GET /api/v1/user/sessions?limit=10&offset=0
func (h *SessionHandler) GetActiveSessions(c *gin.Context) {
	// Extract user_id and jti from context (set by RequireAuth middleware)
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": "Failed to retrieve sessions",
		})
		return
	}

	// Get current JTI (session ID) from context
	jtiInterface, _ := c.Get("jti")
	currentJTI, _ := jtiInterface.(string)

	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Fetch active sessions for the user
	response, err := h.sessionSvc.GetActiveSessionsByUserID(c.Request.Context(), userID, currentJTI, limit, offset)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": "Failed to retrieve sessions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Active sessions retrieved successfully",
		"data":    response,
	})
}

// LogoutOthers logs out the user from all other active sessions/devices
// POST /api/v1/user/sessions/logout-others
func (h *SessionHandler) LogoutOthers(c *gin.Context) {
	// Extract user_id and jti from context (set by RequireAuth middleware)
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": "Failed to logout other sessions",
		})
		return
	}

	// Get current JTI (session ID) from context
	jtiInterface, _ := c.Get("jti")
	currentJTI, _ := jtiInterface.(string)

	// Logout all other sessions
	sessionsTerminated, err := h.sessionSvc.LogoutOtherSessions(c.Request.Context(), userID, currentJTI)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": "Failed to logout other sessions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All other sessions logged out successfully",
		"data": gin.H{
			"sessions_terminated": sessionsTerminated,
		},
	})
}
