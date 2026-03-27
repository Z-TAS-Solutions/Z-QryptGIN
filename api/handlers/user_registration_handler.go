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
		c.Error(err)
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
