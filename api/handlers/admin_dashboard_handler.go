package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
)

type AdminDashboardHandler struct {
	svc service.AdminDashboardService
}

func NewAdminDashboardHandler(svc service.AdminDashboardService) *AdminDashboardHandler {
	return &AdminDashboardHandler{svc: svc}
}

func (h *AdminDashboardHandler) GetAnalytics(c *gin.Context) {
	analytics, err := h.svc.GetAnalytics(c.Request.Context())
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to fetch analytics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   analytics,
	})
}

func (h *AdminDashboardHandler) GetAuthTrends(c *gin.Context) {
	period := c.DefaultQuery("period", "7d")
	if period == "" {
		period = "7d"
	}

	trends, err := h.svc.GetAuthTrends(c.Request.Context(), period)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to fetch auth trends",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   trends,
	})
}

func (h *AdminDashboardHandler) GetRecentAuthActivity(c *gin.Context) {
	activity, err := h.svc.GetRecentAuthActivity(c.Request.Context())
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to fetch recent auth activity",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   activity,
	})
}
