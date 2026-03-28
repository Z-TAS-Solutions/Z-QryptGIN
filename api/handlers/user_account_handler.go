package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type UserAccountHandler struct {
	logger zerolog.Logger
	// Add services as needed (e.g., sessionService, userService, etc.)
}

// NewUserAccountHandler initializes a new UserAccountHandler
func NewUserAccountHandler(logger zerolog.Logger) *UserAccountHandler {
	return &UserAccountHandler{
		logger: logger,
	}
}

// ForceLogoutAllDevices force logs out a user from all their devices
// Requires: admin role in JWT token
func (h *UserAccountHandler) ForceLogoutAllDevices(c *gin.Context) {
	// Extract role from context (set by RequireAuth middleware)
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "role not found in token"})
		return
	}

	// Verify admin role
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin role required"})
		return
	}

	// Parse request body
	var req struct {
		Force bool `json:"force" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: force field is required"})
		return
	}

	// TODO: Implement force logout logic
	// 1. Get user_id from context
	// 2. Query all active sessions for the user
	// 3. Revoke/delete them
	// 4. Return success response

	c.JSON(http.StatusOK, gin.H{"message": "force logout initiated"})
}

// DeleteAccount deletes a user account
// TODO: Implement this endpoint
func (h *UserAccountHandler) DeleteAccount(c *gin.Context) {
	// Extract user_id from context (set by RequireAuth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}

	// TODO: Implement account deletion logic
	// 1. Delete user from database
	// 2. Revoke all sessions
	// 3. Delete credentials/devices
	// 4. Return success response

	c.JSON(http.StatusOK, gin.H{"message": "account deletion initiated", "user_id": userID})
}
