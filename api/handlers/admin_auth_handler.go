package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
)

type AdminAuthHandler struct {
	svc service.AdminAuthService
}

func NewAdminAuthHandler(svc service.AdminAuthService) *AdminAuthHandler {
	return &AdminAuthHandler{svc: svc}
}

func (h *AdminAuthHandler) Login(c *gin.Context) {
	var req dto.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid JSON payload: " + err.Error(),
		})
		return
	}

	resp, err := h.svc.AdminLogin(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "invalid credentials",
		})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusUnauthorized, gin.H{
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

func (h *AdminAuthHandler) Refresh(c *gin.Context) {
	var req dto.AdminRefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid JSON payload: " + err.Error(),
		})
		return
	}

	resp, err := h.svc.AdminRefresh(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "invalid refresh token",
		})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "token refresh failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   resp,
	})
}

func (h *AdminAuthHandler) GetAuthLogs(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	resp, err := h.svc.GetAuthLogs(c.Request.Context(), page, limit)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to fetch auth logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   resp,
	})
}

func (h *AdminAuthHandler) GetAuthAnalytics(c *gin.Context) {
	period := c.DefaultQuery("period", "7d")
	if period == "" {
		period = "7d"
	}

	resp, err := h.svc.GetAuthAnalytics(c.Request.Context(), period)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to fetch auth analytics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   resp,
	})
}
