package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/cache"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
)

type WebAuthnHandler struct {
	logger       zerolog.Logger
	webauthnSvc  *service.WebAuthnService
	userRepo     repository.UserRepository
	sessionCache *cache.WebAuthnSessionCache
	redisClient  *redis.Client
}

// NewWebAuthnHandler creates a new WebAuthnHandler instance
func NewWebAuthnHandler(
	logger zerolog.Logger,
	webauthnSvc *service.WebAuthnService,
	userRepo repository.UserRepository,
	sessionCache *cache.WebAuthnSessionCache,
	redisClient *redis.Client,
) *WebAuthnHandler {
	return &WebAuthnHandler{
		logger:       logger,
		webauthnSvc:  webauthnSvc,
		userRepo:     userRepo,
		sessionCache: sessionCache,
		redisClient:  redisClient,
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

	// Query the user registration cache by email
	// The registration cache links webauthn data with the newly registered, MFA validated user
	regCacheEntry, err := h.getRegistrationCacheByEmail(c.Request.Context(), req.Email)
	if err != nil {
		h.logger.Warn().Str("email", string(req.Email)).Err(err).Msg("Failed to retrieve registration cache entry")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to retrieve registration data",
		})
		return
	}
	if regCacheEntry == nil {
		h.logger.Warn().Str("email", string(req.Email)).Msg("Registration cache entry not found")
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"code":    http.StatusNotFound,
			"message": "User registration not found or expired",
		})
		return
	}

	// Verify that OTP has been verified (MFA status is approved)
	if regCacheEntry.MfaStatus != database.MfaApproved {
		h.logger.Warn().Str("email", string(req.Email)).Str("mfa_status", string(regCacheEntry.MfaStatus)).Msg("User MFA not verified for passkey registration")
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "User must complete OTP verification before registering passkey",
		})
		return
	}

	// Create a PendingUser from the registration cache data
	// This implements the webauthn.User interface for registration ceremony
	pendingUser := &cache.PendingUser{
		ID:          regCacheEntry.UserID,
		Email:       regCacheEntry.Email,
		Name:        regCacheEntry.Name,
		Credentials: make([]webauthn.Credential, 0), // Empty credentials during registration
	}

	// Begin the registration ceremony
	sessionData, creationData, err := h.webauthnSvc.BeginRegistration(c.Request.Context(), pendingUser)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", string(regCacheEntry.UserID)).Msg("Failed to begin WebAuthn registration")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to initiate registration: " + err.Error(),
		})
		return
	}

	// Generate a unique session token
	sessionToken := generateSessionToken()

	// Store session data in cache with TTL, linking it to the registration cache entry via email and userID
	err = h.sessionCache.StoreRegistrationSession(c.Request.Context(), sessionToken, sessionData, regCacheEntry.Email, regCacheEntry.UserID)
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
		Str("user_id", string(regCacheEntry.UserID)).
		Str("email", string(req.Email)).
		Msg("Registration ceremony initiated for MFA-verified user from cache")

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": dto.RegisterBeginResponse{
			SessionToken: sessionToken,
			CreationData: creationData,
		},
	})
}

// getRegistrationCacheByEmail retrieves user registration data from Redis cache by email
func (h *WebAuthnHandler) getRegistrationCacheByEmail(ctx context.Context, email database.Email) (*cache.RegistrationCache, error) {
	key := "registration:email:" + string(email)
	val, err := h.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var entry cache.RegistrationCache
	if err := json.Unmarshal([]byte(val), &entry); err != nil {
		h.logger.Warn().Err(err).Str("key", key).Msg("Malformed user registration cache entry")
		// Delete stale key
		_ = h.redisClient.Del(ctx, key).Err()
		return nil, nil
	}

	return &entry, nil
}

// generateSessionToken creates a unique session token for WebAuthn ceremonies
func generateSessionToken() string {
	return uuid.New().String() + ":" + time.Now().Format("20060102150405")
}
