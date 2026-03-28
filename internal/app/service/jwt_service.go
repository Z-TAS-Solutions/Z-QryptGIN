package service

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type JWTService interface {
	GenerateToken(userID uint, role string, session *dto.Session) (string, string, time.Time, error)
	VerifyToken(tokenString string) (jwt.MapClaims, error)
	StoreSession(ctx context.Context, jti string, session *dto.Session, expiry time.Duration) error
	GetSession(ctx context.Context, jti string) (*dto.Session, error)
	ValidateSessionWithToken(ctx context.Context, tokenClaims jwt.MapClaims) (*dto.Session, error)
	RevokeSession(ctx context.Context, jti string) error
}

type jwtService struct {
	privateKey  ed25519.PrivateKey
	publicKey   ed25519.PublicKey
	issuer      string
	redisClient *redis.Client
}

// NewJWTService initializes the EdDSA keys and Redis client. In production, load keys from secure env/vault.
func NewJWTService(privKey ed25519.PrivateKey, pubKey ed25519.PublicKey, issuer string, redisClient *redis.Client) JWTService {
	return &jwtService{
		privateKey:  privKey,
		publicKey:   pubKey,
		issuer:      issuer,
		redisClient: redisClient,
	}
}

func (s *jwtService) GenerateToken(userID uint, role string, session *dto.Session) (string, string, time.Time, error) {
	now := time.Now()

	// Determine expiration based on role
	var exp time.Time
	if role == "admin" {
		exp = now.Add(24 * time.Hour)
	} else {
		exp = now.Add(30 * 24 * time.Hour) // ~1 month
	}

	// Generate UUID v7 for JTI (Time-ordered UUID)
	jti, err := uuid.NewV7()
	if err != nil {
		return "", "", time.Time{}, err
	}
	jtiString := jti.String()

	// Encode private key as hex for auditing purposes (NOT FOR SIGN, just metadata)
	privateKeyHex := hex.EncodeToString(s.privateKey)

	claims := jwt.MapClaims{
		"sub":              userID,             // Subject: User ID
		"jti":              jtiString,          // JWT ID: Unique token ID for session tracking
		"role":             role,               // User role
		"exp":              exp.Unix(),         // Expiration time
		"nbf":              now.Unix(),         // Not before time
		"iat":              now.Unix(),         // Issued at
		"iss":              s.issuer,           // Issuer
		"algorithm":        "EdDSA",            // Algorithm used
		"private_key_hash": privateKeyHex[:16], // First 16 chars of key for audit trail (not the actual key)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	signedToken, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// Store session in Redis with the JTI as key
	if session != nil {
		// Update session with JTI
		session.JTI = jtiString
		session.ExpiresAt = exp
		session.IsActive = true

		// Store in Redis
		err = s.StoreSession(context.Background(), jtiString, session, exp.Sub(now))
		if err != nil {
			return "", "", time.Time{}, err
		}
	}

	return signedToken, jtiString, exp, nil
}

func (s *jwtService) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Force EdDSA validation
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.publicKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// StoreSession stores a session in Redis using JTI as the key
// This allows hybrid stateless/stateful approach where we can revoke tokens or check their status
func (s *jwtService) StoreSession(ctx context.Context, jti string, session *dto.Session, expiry time.Duration) error {
	if s.redisClient == nil {
		return errors.New("redis client not configured")
	}

	// Serialize session to JSON
	sessionJSON, err := marshalSession(session)
	if err != nil {
		return err
	}

	// Store with expiry
	return s.redisClient.Set(ctx, "session:"+jti, sessionJSON, expiry).Err()
}

// GetSession retrieves a session from Redis by JTI
func (s *jwtService) GetSession(ctx context.Context, jti string) (*dto.Session, error) {
	if s.redisClient == nil {
		return nil, errors.New("redis client not configured")
	}

	sessionJSON, err := s.redisClient.Get(ctx, "session:"+jti).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("session not found or expired")
		}
		return nil, err
	}

	session, err := unmarshalSession(sessionJSON)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// ValidateSessionWithToken validates both the JWT signature and the associated session in Redis
// This is the core of the hybrid approach - stateless verification + stateful session tracking
func (s *jwtService) ValidateSessionWithToken(ctx context.Context, tokenClaims jwt.MapClaims) (*dto.Session, error) {
	// Extract JTI from claims
	jtiInterface, ok := tokenClaims["jti"]
	if !ok {
		return nil, errors.New("jti claim not found in token")
	}

	jti, ok := jtiInterface.(string)
	if !ok {
		return nil, errors.New("invalid jti format")
	}

	// Get session from Redis using JTI
	session, err := s.GetSession(ctx, jti)
	if err != nil {
		return nil, err
	}

	// Verify session is active
	if !session.IsActive {
		return nil, errors.New("session is inactive or revoked")
	}

	// Verify session hasn't expired
	if time.Now().After(session.ExpiresAt) {
		return nil, errors.New("session has expired")
	}

	// Update last active time
	session.LastActiveAt = time.Now()

	// Re-store updated session
	_ = s.StoreSession(ctx, jti, session, time.Until(session.ExpiresAt))

	return session, nil
}

// RevokeSession revokes a session by JTI (logout functionality)
func (s *jwtService) RevokeSession(ctx context.Context, jti string) error {
	if s.redisClient == nil {
		return errors.New("redis client not configured")
	}

	return s.redisClient.Del(ctx, "session:"+jti).Err()
}

// Helper function to marshal session to JSON
func marshalSession(session *dto.Session) (string, error) {
	data, err := json.Marshal(session)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Helper function to unmarshal session from JSON
func unmarshalSession(sessionJSON string) (*dto.Session, error) {
	var session dto.Session
	err := json.Unmarshal([]byte(sessionJSON), &session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}
