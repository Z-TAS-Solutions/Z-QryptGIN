package handlers

import (
	"net/http"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	svc service.DashboardService
}

func NewDashboardHandler(svc service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) GetAnalytics(c *gin.Context) {
	resp, err := h.svc.GetAnalytics()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *DashboardHandler) GetAuthTrends(c *gin.Context) {
	interval := c.Query("interval") // defaults to hour normally mapped in service limits

	resp, err := h.svc.GetAuthTrends(interval)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *DashboardHandler) GetRecentAuthActivity(c *gin.Context) {
	var req dto.RecentAuthActivityRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.svc.GetRecentAuthActivity(req.Page, req.Limit, req.Status)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
