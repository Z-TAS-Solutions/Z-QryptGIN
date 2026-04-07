package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
)

type ProfileHandler struct {
	svc service.ProfileService
}

func NewProfileHandler(svc service.ProfileService) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// Extract user ID from JWT context (would be set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "user not authenticated",
		})
		return
	}

	profile, err := h.svc.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to fetch profile",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   profile,
	})
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "user not authenticated",
		})
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    http.StatusBadRequest,
			"message": "invalid JSON payload: " + err.Error(),
		})
		return
	}

	resp, err := h.svc.UpdateProfile(c.Request.Context(), userID.(string), req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": resp.Message,
		})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": resp.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   resp,
	})
}
