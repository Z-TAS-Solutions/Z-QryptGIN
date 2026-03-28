package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
)

type UserRegistrationHandler struct {
	svc service.UserRegistrationService
}

func NewUserRegistrationHandler(svc service.UserRegistrationService) *UserRegistrationHandler {
	return &UserRegistrationHandler{svc: svc}
}

func (h *UserRegistrationHandler) Register(c *gin.Context) {
	var req dto.UserRegistrationDetailsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    http.StatusBadRequest,
			"message": "invalid JSON payload: " + err.Error(),
		})
		return
	}

	resp, err := h.svc.RegisterUser(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	if !resp.Success {
		c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "a user with one or many of these fields already exists", "data": resp})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": resp})
}

func (h *UserRegistrationHandler) VerifyOTP(c *gin.Context) {
	var req dto.UserRegistrationOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    http.StatusBadRequest,
			"message": "invalid JSON payload: " + err.Error(),
		})
		return
	}

	resp, err := h.svc.VerifyOTP(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrRegistrationNotFound {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "registration info not found or expired"})
			return
		}
		if err == service.ErrInvalidOTP {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "invalid otp"})
			return
		}

		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": resp})
}

func (h *UserRegistrationHandler) ResendOTP(c *gin.Context) {
	var req dto.ResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    http.StatusBadRequest,
			"message": "invalid JSON payload: " + err.Error(),
		})
		return
	}

	resp, err := h.svc.ResendOTP(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	if !resp.Success {
		c.JSON(http.StatusTooManyRequests, gin.H{"status": "error", "message": resp.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": resp.Message})
}
