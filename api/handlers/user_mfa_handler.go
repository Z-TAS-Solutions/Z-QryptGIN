package handlers

import (
	"net/http"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
	"github.com/gin-gonic/gin"
)

type UserMfaHandler struct {
	svc service.UserMfaService
}

func NewUserMfaHandler(svc service.UserMfaService) *UserMfaHandler {
	return &UserMfaHandler{svc: svc}
}

func (h *UserMfaHandler) Send(c *gin.Context) {
	var req dto.MfaSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.svc.Send(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserMfaHandler) Respond(c *gin.Context) {
	var req dto.MfaRespondRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.svc.Respond(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
