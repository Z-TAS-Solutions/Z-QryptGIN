package handlers

import (
	"net/http"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
	"github.com/gin-gonic/gin"
)

type AdminUserHandler struct {
	userService service.UserService
}

func NewAdminUserHandler(userService service.UserService) *AdminUserHandler {
	return &AdminUserHandler{
		userService: userService,
	}
}

// ListUsers handles the GET /api/v1/admin/users endpoint.
func (h *AdminUserHandler) ListUsers(c *gin.Context) {
	var req dto.ListUsersRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "BadRequest",
			"message": "Invalid query parameters: " + err.Error(),
		})
		return
	}

	response, err := h.userService.ListUsers(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "InternalServerError",
			"message": "Failed to fetch users: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
