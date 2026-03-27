package handlers

import (
	"net/http"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
	"github.com/gin-gonic/gin"
)

type UserAuthHandler struct {
	svc service.UserAuthService
}

func NewUserAuthHandler(svc service.UserAuthService) *UserAuthHandler {
	return &UserAuthHandler{svc: svc}
}

func (h *UserAuthHandler) RegisterOptions(c *gin.Context) {
	var req dto.PasskeyRegisterOptionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.svc.RegisterOptions(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserAuthHandler) RegisterVerify(c *gin.Context) {
	var req dto.PasskeyRegisterVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.svc.RegisterVerify(req)
	if err != nil {
		c.Error(err)
		return
	}
	// Note: prompt wants 201 for register verification success
	c.JSON(http.StatusCreated, resp)
}

func (h *UserAuthHandler) LoginOptions(c *gin.Context) {
	var req dto.PasskeyLoginOptionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.svc.LoginOptions(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *UserAuthHandler) LoginVerify(c *gin.Context) {
	var req dto.PasskeyLoginVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.svc.LoginVerify(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
