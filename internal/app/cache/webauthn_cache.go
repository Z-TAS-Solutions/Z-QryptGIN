package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/redis/go-redis/v9"
)

type CredentialCache struct {
	UserID          uint     `json:"u_id"`
	CredentialID    []byte   `json:"c_id"`
	PublicKey       []byte   `json:"pk"`
	AttestationType string   `json:"at"`
	Transport       []string `json:"tr"`
	
	// FIDO2 State & Validation
	UserPresent    bool   `json:"up"`
	UserVerified   bool   `json:"uv"`
	BackupEligible bool   `json:"be"`
	BackupState    bool   `json:"bs"`
	
	// Security & Metadata (Crucial for Z-TAS)
	AAGUID       []byte `json:"aa"` // Needed to identify the device type (Yubikey vs Apple)
	SignCount    uint32 `json:"sc"` // MUST be cached to detect cloned keys
	CloneWarning bool   `json:"cw"` // Persistent warning state
	
	// Optional: Metadata for UI
	AuthenticatorName string `json:"n"`
}

// WebAuthnRegistrationSession wraps WebAuthn session data with registration cache linkage
type WebAuthnRegistrationSession struct {
	SessionData *webauthn.SessionData `json:"session_data"`
	Email       database.Email        `json:"email"`
	UserID      database.UserCustomID `json:"user_id"`
}

// PendingUser is an adapter that implements the webauthn.User interface
// for users during the registration cache phase (before database insert)
type PendingUser struct {
	ID          database.UserCustomID `json:"id"`
	Email       database.Email        `json:"email"`
	Name        database.Name         `json:"name"`
	Credentials []webauthn.Credential `json:"credentials"` // Empty during registration
}

// WebAuthnID returns the unique identifier for WebAuthn
func (p *PendingUser) WebAuthnID() []byte {
	return []byte(string(p.ID))
}

// WebAuthnName returns the user's email (name) for WebAuthn
func (p *PendingUser) WebAuthnName() string {
	return string(p.Email)
}

// WebAuthnDisplayName returns the user's full name for WebAuthn
func (p *PendingUser) WebAuthnDisplayName() string {
	return string(p.Name)
}

// WebAuthnIcon returns the user's icon URL (not used in registration)
func (p *PendingUser) WebAuthnIcon() string {
	return ""
}

// WebAuthnCredentials returns existing credentials (empty during registration cache phase)
func (p *PendingUser) WebAuthnCredentials() []webauthn.Credential {
	return p.Credentials
}

// WebAuthnSessionCache handles caching of WebAuthn session data during registration/authentication
type WebAuthnSessionCache struct {
	client *redis.Client
}

// NewWebAuthnSessionCache creates a new WebAuthnSessionCache
func NewWebAuthnSessionCache(client *redis.Client) *WebAuthnSessionCache {
	return &WebAuthnSessionCache{client: client}
}

// StoreRegistrationSession stores the registration session data with email linkage in Redis with a TTL
// This links WebAuthn session data with the user registration cache entry
func (c *WebAuthnSessionCache) StoreRegistrationSession(ctx context.Context, sessionToken string, sessionData *webauthn.SessionData, email database.Email, userID database.UserCustomID) error {
	// Create a wrapped session with email/userID linkage to registration cache
	wrappedSession := WebAuthnRegistrationSession{
		SessionData: sessionData,
		Email:       email,
		UserID:      userID,
	}

	// Serialize wrapped session data to JSON
	data, err := json.Marshal(wrappedSession)
	if err != nil {
		return fmt.Errorf("failed to marshal webauthn registration session: %w", err)
	}

	// Store in Redis with 15-minute expiration (enough time for user to complete registration)
	key := fmt.Sprintf("webauthn:registration:session:%s", sessionToken)
	err = c.client.Set(ctx, key, data, 15*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to store registration session in cache: %w", err)
	}

	return nil
}

// GetRegistrationSession retrieves the registration session data with email linkage from Redis
func (c *WebAuthnSessionCache) GetRegistrationSession(ctx context.Context, sessionToken string) (*WebAuthnRegistrationSession, error) {
	key := fmt.Sprintf("webauthn:registration:session:%s", sessionToken)
	
	data, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("registration session not found or expired")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve registration session: %w", err)
	}

	var wrappedSession WebAuthnRegistrationSession
	err = json.Unmarshal([]byte(data), &wrappedSession)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal webauthn registration session: %w", err)
	}

	return &wrappedSession, nil
}

// GetRegistrationSessionData retrieves only the WebAuthn session data from the wrapped session
func (c *WebAuthnSessionCache) GetRegistrationSessionData(ctx context.Context, sessionToken string) (*webauthn.SessionData, error) {
	wrappedSession, err := c.GetRegistrationSession(ctx, sessionToken)
	if err != nil {
		return nil, err
	}
	return wrappedSession.SessionData, nil
}

// GetRegistrationSessionEmail retrieves the email linked to the WebAuthn registration session
// This queries the user_registration_cache using the email to get the full registration data
func (c *WebAuthnSessionCache) GetRegistrationSessionEmail(ctx context.Context, sessionToken string) (database.Email, error) {
	wrappedSession, err := c.GetRegistrationSession(ctx, sessionToken)
	if err != nil {
		return "", err
	}
	return wrappedSession.Email, nil
}

// GetRegistrationSessionUserID retrieves the user ID linked to the WebAuthn registration session
func (c *WebAuthnSessionCache) GetRegistrationSessionUserID(ctx context.Context, sessionToken string) (database.UserCustomID, error) {
	wrappedSession, err := c.GetRegistrationSession(ctx, sessionToken)
	if err != nil {
		return "", err
	}
	return wrappedSession.UserID, nil
}

// DeleteRegistrationSession removes the registration session from Redis (called after successful registration)
func (c *WebAuthnSessionCache) DeleteRegistrationSession(ctx context.Context, sessionToken string) error {
	key := fmt.Sprintf("webauthn:registration:session:%s", sessionToken)
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete registration session: %w", err)
	}
	return nil
}

// StoreAuthenticationSession stores the authentication session data in Redis with a TTL
func (c *WebAuthnSessionCache) StoreAuthenticationSession(ctx context.Context, sessionToken string, sessionData *webauthn.SessionData) error {
	// Serialize session data to JSON
	data, err := json.Marshal(sessionData)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}

	// Store in Redis with 10-minute expiration (enough time for user to complete authentication)
	key := fmt.Sprintf("webauthn:authentication:session:%s", sessionToken)
	err = c.client.Set(ctx, key, data, 10*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to store authentication session in cache: %w", err)
	}

	return nil
}

// GetAuthenticationSession retrieves the authentication session data from Redis
func (c *WebAuthnSessionCache) GetAuthenticationSession(ctx context.Context, sessionToken string) (*webauthn.SessionData, error) {
	key := fmt.Sprintf("webauthn:authentication:session:%s", sessionToken)
	
	data, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("authentication session not found or expired")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve authentication session: %w", err)
	}

	var sessionData webauthn.SessionData
	err = json.Unmarshal([]byte(data), &sessionData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session data: %w", err)
	}

	return &sessionData, nil
}

// DeleteAuthenticationSession removes the authentication session from Redis
func (c *WebAuthnSessionCache) DeleteAuthenticationSession(ctx context.Context, sessionToken string) error {
	key := fmt.Sprintf("webauthn:authentication:session:%s", sessionToken)
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete authentication session: %w", err)
	}
	return nil
}