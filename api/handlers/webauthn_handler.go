package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/cache"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
)

type WebAuthnHandler struct {
	logger       zerolog.Logger
	webauthnSvc  *service.WebAuthnService
	userRepo     repository.UserRepository
	sessionCache *cache.WebAuthnSessionCache
}

// NewWebAuthnHandler creates a new WebAuthnHandler instance
func NewWebAuthnHandler(
	logger zerolog.Logger,
	webauthnSvc *service.WebAuthnService,
	userRepo repository.UserRepository,
	sessionCache *cache.WebAuthnSessionCache,
) *WebAuthnHandler {
	return &WebAuthnHandler{
		logger:       logger,
		webauthnSvc:  webauthnSvc,
		userRepo:     userRepo,
		sessionCache: sessionCache,
	}
}

// RegisterBegin initiates the passkey registration process
// POST /api/v1/webauthn/register/begin
func (h *WebAuthnHandler) RegisterBegin(c *gin.Context) {
	var req dto.RegisterBeginRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid request payload for RegisterBegin")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    http.StatusBadRequest,
			"message": "Invalid request payload: " + err.Error(),
		})
		return
	}

	// Find user by email
	user, err := h.userRepo.FindByEmail(string(req.Email))
	if err != nil {
		h.logger.Warn().Str("email", string(req.Email)).Err(err).Msg("User not found for registration")
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"code":    http.StatusNotFound,
			"message": "User not found",
		})
		return
	}

	// Begin the registration ceremony
	sessionData, creationData, err := h.webauthnSvc.BeginRegistration(c.Request.Context(), user)
	if err != nil {
		h.logger.Error().Err(err).Uint("user_id", user.ID).Msg("Failed to begin WebAuthn registration")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to initiate registration: " + err.Error(),
		})
		return
	}

	// Generate a unique session token
	sessionToken := generateSessionToken()

	// Store session data in cache with TTL
	err = h.sessionCache.StoreRegistrationSession(c.Request.Context(), sessionToken, sessionData)
	if err != nil {
		h.logger.Error().Err(err).Str("session_token", sessionToken).Msg("Failed to store registration session in cache")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to store session data",
		})
		return
	}

	h.logger.Info().
		Str("session_token", sessionToken).
		Uint("user_id", user.ID).
		Str("email", string(req.Email)).
		Msg("Registration ceremony initiated")

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": dto.RegisterBeginResponse{
			SessionToken: sessionToken,
			CreationData: creationData,
		},
	})
}

// generateSessionToken creates a unique session token for WebAuthn ceremonies
func generateSessionToken() string {
	return uuid.New().String() + ":" + time.Now().Format("20060102150405")
}
