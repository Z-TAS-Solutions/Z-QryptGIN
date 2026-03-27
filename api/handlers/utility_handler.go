package handlers

import (
	"net/http"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/gin-gonic/gin"
)

type UtilityHandler struct {
}

func NewUtilityHandler() *UtilityHandler {
	return &UtilityHandler{}
}

func (h *UtilityHandler) Ping(c *gin.Context) {
	resp := dto.PingResponse{
		Message: "pong",
	}
	resp.Data.ServerTimestamp = time.Now().UnixMilli()

	c.JSON(http.StatusOK, resp)
}
