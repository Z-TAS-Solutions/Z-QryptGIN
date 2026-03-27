package handlers

import (
	"net/http"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
	"github.com/gin-gonic/gin"
)

type AdminUserHandler struct {
	svc service.AdminUserService
}

func NewAdminUserHandler(svc service.AdminUserService) *AdminUserHandler {
	return &AdminUserHandler{svc: svc}
}

func (h *AdminUserHandler) ListUsers(c *gin.Context) {
	var req dto.ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.svc.ListUsers(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AdminUserHandler) GetUser(c *gin.Context) {
	userID := c.Param("userId")

	resp, err := h.svc.GetUser(userID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AdminUserHandler) UpdateLockStatus(c *gin.Context) {
	userID := c.Param("userId")
	var req dto.LockUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.svc.UpdateLockStatus(userID, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AdminUserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("userId")

	resp, err := h.svc.DeleteUser(userID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
