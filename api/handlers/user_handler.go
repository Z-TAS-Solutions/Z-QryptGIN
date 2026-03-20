package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err) // Passes to your ErrorHandler middleware
		return
	}

	resp, err := h.svc.RegisterUser(req)
	if err != nil {
		c.Error(err) 
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": resp})
}