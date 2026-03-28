package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
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
	logger                  zerolog.Logger
	webauthnSvc             *service.WebAuthnService
	userRepo                repository.UserRepository
	credentialRepo          repository.WebAuthnCredentialRepository
	sessionCache            *cache.WebAuthnSessionCache
	userRegistrationService service.UserRegistrationService
	redisClient             *redis.Client
	jwtService              service.JWTService
}

// NewWebAuthnHandler creates a new WebAuthnHandler instance
func NewWebAuthnHandler(
	logger zerolog.Logger,
	webauthnSvc *service.WebAuthnService,
	userRepo repository.UserRepository,
	sessionCache *cache.WebAuthnSessionCache,
	redisClient *redis.Client,
	jwtService service.JWTService,
) *WebAuthnHandler {
	return &WebAuthnHandler{
		logger:       logger,
		webauthnSvc:  webauthnSvc,
		userRepo:     userRepo,
		sessionCache: sessionCache,
		redisClient:  redisClient,
		jwtService:   jwtService,
	}
}

// NewWebAuthnHandlerWithRegistration creates a new WebAuthnHandler with user registration service support
func NewWebAuthnHandlerWithRegistration(
	logger zerolog.Logger,
	webauthnSvc *service.WebAuthnService,
	userRepo repository.UserRepository,
	credentialRepo repository.WebAuthnCredentialRepository,
	sessionCache *cache.WebAuthnSessionCache,
	userRegistrationService service.UserRegistrationService,
	redisClient *redis.Client,
	jwtService service.JWTService,
) *WebAuthnHandler {
	return &WebAuthnHandler{
		logger:                  logger,
		webauthnSvc:             webauthnSvc,
		userRepo:                userRepo,
		credentialRepo:          credentialRepo,
		sessionCache:            sessionCache,
		userRegistrationService: userRegistrationService,
		redisClient:             redisClient,
		jwtService:              jwtService,
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

// RegisterFinish completes the WebAuthn registration ceremony
// POST /api/v1/webauthn/register/finish
// Expects: X-Session-Token header + raw credential creation response in body
func (h *WebAuthnHandler) RegisterFinish(c *gin.Context) {
	// 1. Extract session token from header
	sessionToken := c.GetHeader("X-Session-Token")
	if sessionToken == "" {
		h.logger.Warn().Msg("Missing X-Session-Token header")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    http.StatusBadRequest,
			"message": "Missing X-Session-Token header",
		})
		return
	}

	// 2. Retrieve the session from cache
	wrappedSession, err := h.sessionCache.GetRegistrationSession(c.Request.Context(), sessionToken)
	if err != nil {
		h.logger.Warn().Str("session_token", sessionToken).Err(err).Msg("Failed to retrieve registration session")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Session expired or invalid",
		})
		return
	}

	// 3. Parse the credential creation response from the request body
	credentialResponse, err := protocol.ParseCredentialCreationResponseBody(c.Request.Body)
	if err != nil {
		h.logger.Error().Err(err).Str("session_token", sessionToken).Msg("Failed to parse credential creation response")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    http.StatusBadRequest,
			"message": "Invalid credential response: " + err.Error(),
		})
		return
	}

	// 4. Call the registration service to finish registration
	// The service will:
	// - Verify the credential with the stored session data
	// - Create the user in the database
	// - Save the WebAuthn credential
	// - Clean up the caches
	err = h.getUserRegistrationService().FinishRegistration(
		c.Request.Context(),
		wrappedSession.UserID,
		credentialResponse,
		wrappedSession.SessionData,
	)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("session_token", sessionToken).
			Str("user_id", string(wrappedSession.UserID)).
			Msg("Registration finish failed")

		// Return appropriate error based on error type
		if err.Error() == "registration entry not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"code":    http.StatusNotFound,
				"message": "Registration not found - please start over",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Registration failed: " + err.Error(),
		})
		return
	}

	// 5. Delete the session from cache after successful registration
	_ = h.sessionCache.DeleteRegistrationSession(c.Request.Context(), sessionToken)

	// 6. Load the newly created user from database
	registeredUser, err := h.userRepo.FindByCustomID(string(wrappedSession.UserID))
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", string(wrappedSession.UserID)).Msg("Failed to load registered user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Registration completed but failed to generate authentication token",
		})
		return
	}

	// 7. Create session object for Redis session tracking
	// This is the key component of the hybrid approach for the newly registered user
	registrationSession := &dto.Session{
		ID:           uuid.New().String(),
		UserID:       registeredUser.ID,
		DeviceName:   c.GetHeader("User-Agent"),
		DeviceID:     fmt.Sprintf("%s-registration", uuid.New().String()),
		IPAddress:    c.ClientIP(),
		UserAgent:    c.GetHeader("User-Agent"),
		IsActive:     true,
		MfaStatus:    dto.MfaStatusVerified,
		LastActiveAt: time.Now(),
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour), // 30 days expiry
	}

	// 8. Generate JWT token with session
	// This will also store the session in Redis via the JWT service
	userRole := string(registeredUser.Role)
	jwtToken, jti, expiry, err := h.jwtService.GenerateToken(registeredUser.ID, userRole, registrationSession)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", string(registeredUser.ID)).Msg("Failed to generate JWT token after registration")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Registration completed but failed to generate authentication token",
		})
		return
	}

	h.logger.Info().
		Str("session_token", sessionToken).
		Str("user_id", string(wrappedSession.UserID)).
		Str("jti", jti).
		Str("ip_address", c.ClientIP()).
		Msg("User successfully registered with WebAuthn credential and authenticated")

	// 9. Return success response with JWT token
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"code":   http.StatusCreated,
		"data": gin.H{
			"token":      jwtToken,
			"token_type": "Bearer",
			"expires_in": expiry.Unix(),
			"user_id":    registeredUser.ID,
			"email":      registeredUser.Email,
			"message":    "Registration completed successfully",
		},
	})
}

// LoginBegin initiates the WebAuthn authentication ceremony
// POST /api/v1/webauthn/login/begin
// Expects: Optional username in JSON body (for username-based flow)
//
//	Can be empty for usernameless/discoverable credential flow
func (h *WebAuthnHandler) LoginBegin(c *gin.Context) {
	var req dto.LoginBeginRequest

	// Bind and validate request (username is optional)
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid request payload for LoginBegin")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    http.StatusBadRequest,
			"message": "Invalid request payload: " + err.Error(),
		})
		return
	}

	var sessionData *webauthn.SessionData
	var assertionData *protocol.CredentialAssertion
	var err error

	// Determine if this is a username-based or usernameless flow
	if req.Username != "" {
		// Username-based authentication: Look up the user
		user, err := h.userRepo.FindByEmail(req.Username)
		if err != nil {
			h.logger.Warn().Str("username", req.Username).Err(err).Msg("User not found for login")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"code":    http.StatusUnauthorized,
				"message": "User not found",
			})
			return
		}

		// Begin login for this specific user
		sessionData, assertionData, err = h.webauthnSvc.BeginLogin(c.Request.Context(), user)
		if err != nil {
			h.logger.Error().Err(err).Str("username", req.Username).Msg("Failed to begin WebAuthn login")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"code":    http.StatusInternalServerError,
				"message": "Failed to initiate login: " + err.Error(),
			})
			return
		}
	} else {
		// Usernameless/Discoverable Credential flow
		// Let the authenticator device choose which account to sign in with
		sessionData, assertionData, err = h.webauthnSvc.BeginLogin(c.Request.Context(), nil)
		if err != nil {
			h.logger.Error().Err(err).Msg("Failed to begin WebAuthn usernameless login")
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"code":    http.StatusInternalServerError,
				"message": "Failed to initiate login: " + err.Error(),
			})
			return
		}
	}

	// Generate a unique session token
	sessionToken := generateSessionToken()

	// Store session data in cache with TTL
	err = h.sessionCache.StoreAuthenticationSession(c.Request.Context(), sessionToken, sessionData)
	if err != nil {
		h.logger.Error().Err(err).Str("session_token", sessionToken).Msg("Failed to store authentication session in cache")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to store session data",
		})
		return
	}

	h.logger.Info().
		Str("session_token", sessionToken).
		Str("username", req.Username).
		Msg("Authentication ceremony initiated")

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": dto.LoginBeginResponse{
			SessionToken:  sessionToken,
			AssertionData: assertionData,
		},
	})
}

// LoginFinish completes the WebAuthn authentication ceremony and generates JWT token
// POST /api/v1/webauthn/login/finish
// Expects: X-Session-Token header + raw credential assertion response in body
// Returns: JWT token with session data for hybrid-stateless/stateful authentication
func (h *WebAuthnHandler) LoginFinish(c *gin.Context) {
	// 1. Extract session token from header
	sessionToken := c.GetHeader("X-Session-Token")
	if sessionToken == "" {
		h.logger.Warn().Msg("Missing X-Session-Token header")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    http.StatusBadRequest,
			"message": "Missing X-Session-Token header",
		})
		return
	}

	// 2. Retrieve the session from cache
	sessionData, err := h.sessionCache.GetAuthenticationSession(c.Request.Context(), sessionToken)
	if err != nil {
		h.logger.Warn().Str("session_token", sessionToken).Err(err).Msg("Failed to retrieve authentication session")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Session expired or invalid",
		})
		return
	}

	// 3. Parse the credential assertion response from the request body
	assertionResponse, err := protocol.ParseCredentialRequestResponseBody(c.Request.Body)
	if err != nil {
		h.logger.Error().Err(err).Str("session_token", sessionToken).Msg("Failed to parse credential assertion response")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    http.StatusBadRequest,
			"message": "Invalid assertion response: " + err.Error(),
		})
		return
	}

	// 4. Find the credential by credential ID
	credentialID := assertionResponse.RawID
	credential, err := h.credentialRepo.FindCredentialByID(credentialID)
	if err != nil {
		h.logger.Warn().Bytes("credential_id", credentialID).Err(err).Msg("Credential not found for login")
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Authentication failed: credential not recognized",
		})
		return
	}

	// 5. Load the user from the database by ID
	authenticatedUser, err := h.userRepo.FindByID(credential.UserID)
	if err != nil {
		h.logger.Error().Err(err).Uint("user_id", credential.UserID).Msg("Failed to load user for authentication")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Authentication failed",
		})
		return
	}

	// 6. Load user credentials for WebAuthn verification
	credentials, err := h.credentialRepo.FindCredentialsByUserID(authenticatedUser.ID)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", string(credential.UserID)).Msg("Failed to load user credentials")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Authentication failed",
		})
		return
	}

	// 7. Create webauthn.User interface for verification
	webauthnUser := &cache.PendingUser{
		ID:          authenticatedUser.CustomID,
		Email:       authenticatedUser.Email,
		Name:        authenticatedUser.Name,
		Credentials: credentialsToWebAuthnCredentials(credentials),
	}

	// 8. Complete the WebAuthn login ceremony
	userData, err := h.webauthnSvc.FinishLogin(c.Request.Context(), webauthnUser, assertionResponse, sessionData)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", string(authenticatedUser.ID)).Msg("WebAuthn login verification failed")
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Authentication failed: " + err.Error(),
		})
		return
	}

	// 9. Check for cloned credentials (sign count anomaly)
	if userData.Sign < credential.SignCount {
		h.logger.Warn().
			Uint("user_id", authenticatedUser.ID).
			Uint32("stored_sign_count", credential.SignCount).
			Uint32("new_sign_count", userData.Sign).
			Msg("Potential cloned credential detected")
		_ = h.credentialRepo.UpdateCloneWarning(credentialID, true)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Authentication rejected: potential credential cloning detected",
		})
		return
	}

	// 10. Update the sign count for this credential
	err = h.credentialRepo.UpdateSignCount(credentialID, userData.Sign)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to update credential sign count")
		// Non-fatal, continue with login
	}

	// 11. Create session object for Redis session tracking
	// This is the key component of the hybrid approach
	deviceID := string(credential.AuthenticatorName)
	if deviceID == "" {
		// Use credential ID hex representation as device identifier
		deviceID = fmt.Sprintf("device-%d", credential.ID)
	}
	sessionInfo := &dto.Session{
		ID:           uuid.New().String(),
		UserID:       authenticatedUser.ID,
		DeviceName:   c.GetHeader("User-Agent"), // User-Agent as device identifier
		DeviceID:     deviceID,
		IPAddress:    c.ClientIP(),
		UserAgent:    c.GetHeader("User-Agent"),
		IsActive:     true,
		MfaStatus:    dto.MfaStatusVerified,
		LastActiveAt: time.Now(),
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour), // 30 days expiry
	}

	// 12. Generate JWT token with session
	// This will also store the session in Redis via the JWT service
	userRole := string(authenticatedUser.Role)
	jwtToken, jti, expiry, err := h.jwtService.GenerateToken(authenticatedUser.ID, userRole, sessionInfo)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", string(authenticatedUser.ID)).Msg("Failed to generate JWT token")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to generate authentication token",
		})
		return
	}

	// 13. Clean up WebAuthn session from cache
	_ = h.sessionCache.DeleteAuthenticationSession(c.Request.Context(), sessionToken)

	h.logger.Info().
		Str("user_id", string(authenticatedUser.ID)).
		Str("jti", jti).
		Str("ip_address", c.ClientIP()).
		Msg("User successfully authenticated via WebAuthn")

	// 14. Return JWT token to client
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"code":   http.StatusOK,
		"data": gin.H{
			"token":      jwtToken,
			"token_type": "Bearer",
			"expires_in": expiry.Unix(),
			"user_id":    authenticatedUser.ID,
			"email":      authenticatedUser.Email,
		},
	})
}

// getUserRegistrationService is a helper to get the registration service
// This would normally be injected, but for now we use a pattern to access it
func (h *WebAuthnHandler) getUserRegistrationService() service.UserRegistrationService {
	// This is a placeholder - in the actual implementation, it should be injected in NewWebAuthnHandler
	// For now, we'll need to make sure this is available
	return h.userRegistrationService
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

// credentialsToWebAuthnCredentials converts database WebAuthn credentials to webauthn.Credential format
func credentialsToWebAuthnCredentials(dbCredentials []database.WebAuthnCredential) []webauthn.Credential {
	var waCredentials []webauthn.Credential
	for _, dbCred := range dbCredentials {
		transports := []string{}
		if len(dbCred.Transport) > 0 {
			transports = dbCred.Transport
		}
		waCredentials = append(waCredentials, webauthn.Credential{
			ID:        dbCred.CredentialID,
			PublicKey: dbCred.PublicKey,
			Sign:      dbCred.SignCount,
		})
	}
	return waCredentials
}
