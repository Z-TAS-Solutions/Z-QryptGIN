package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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

// WebAuthnSessionCache handles caching of WebAuthn session data during registration/authentication
type WebAuthnSessionCache struct {
	client *redis.Client
}

// NewWebAuthnSessionCache creates a new WebAuthnSessionCache
func NewWebAuthnSessionCache(client *redis.Client) *WebAuthnSessionCache {
	return &WebAuthnSessionCache{client: client}
}

// StoreRegistrationSession stores the registration session data in Redis with a TTL
func (c *WebAuthnSessionCache) StoreRegistrationSession(ctx context.Context, sessionToken string, sessionData *webauthn.SessionData) error {
	// Serialize session data to JSON
	data, err := json.Marshal(sessionData)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}

	// Store in Redis with 15-minute expiration (enough time for user to complete registration)
	key := fmt.Sprintf("webauthn:registration:session:%s", sessionToken)
	err = c.client.Set(ctx, key, data, 15*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to store registration session in cache: %w", err)
	}

	return nil
}

// GetRegistrationSession retrieves the registration session data from Redis
func (c *WebAuthnSessionCache) GetRegistrationSession(ctx context.Context, sessionToken string) (*webauthn.SessionData, error) {
	key := fmt.Sprintf("webauthn:registration:session:%s", sessionToken)
	
	data, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("registration session not found or expired")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve registration session: %w", err)
	}

	var sessionData webauthn.SessionData
	err = json.Unmarshal([]byte(data), &sessionData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session data: %w", err)
	}

	return &sessionData, nil
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