package handlers

import (
	"net/http"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardSvc service.DashboardService
}

// NewDashboardHandler creates a new dashboard handler instance
// Requires a DashboardService for analytics operations
func NewDashboardHandler(dashboardSvc service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardSvc: dashboardSvc,
	}
}

// GetAuthenticationTrends handles GET /api/v1/admin/dashboard/auth-trends
// Returns authentication activity trends over the last 24 hours, structured for time-series visualization
// Query Parameters:
//   - interval (optional, default: "hour"): Aggregation interval ("minute" or "hour")
//
// Response: Time-series data with success and failure counts for each interval
// Requires: Admin role and valid JWT + session cache validation
func (h *DashboardHandler) GetAuthenticationTrends(c *gin.Context) {
	// Extract and validate interval parameter
	interval := c.DefaultQuery("interval", "hour")

	// Get auth trends from service
	trends, err := h.dashboardSvc.GetAuthenticationTrends(c.Request.Context(), interval)
	if err != nil {
		// Check if it's a validation error (invalid interval)
		if err.Error() == "Invalid interval. Supported values are 'minute' or 'hour'" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "BadRequest",
				"message": "Invalid interval. Supported values are 'minute' or 'hour'",
			})
			return
		}

		// Internal server error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "InternalServerError",
			"message": "Failed to fetch authentication trends",
		})
		return
	}

	// Return successful response
	c.JSON(http.StatusOK, trends)
}

// GetDashboardMetrics handles dashboard metrics summary endpoint
// Returns aggregated authentication metrics for the last 24 hours
// Requires: Admin role
func (h *DashboardHandler) GetDashboardMetrics(c *gin.Context) {
	metrics, err := h.dashboardSvc.GetDashboardMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "InternalServerError",
			"message": "Failed to fetch dashboard metrics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   metrics,
	})
}
