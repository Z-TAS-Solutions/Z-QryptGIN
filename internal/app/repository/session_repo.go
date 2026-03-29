package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/redis/go-redis/v9"
)

// SessionRepository handles operations for session tracking and analytics
type SessionRepository interface {
	StoreSession(ctx context.Context, sessionID string, userID uint, expiration time.Duration) error
	GetUserIDBySession(ctx context.Context, sessionID string) (uint, error)
	IncrementAnalytics(ctx context.Context, metricKey string) error
	CreateSession(ctx context.Context, session *dto.Session) error
	GetSessionByJTI(ctx context.Context, jti string) (*dto.Session, error)
	GetActiveSessionsByUserID(ctx context.Context, userID uint) ([]dto.Session, error)
}

type sessionRepository struct {
	redis *redis.Client
}

func NewSessionRepository(redisClient *redis.Client) SessionRepository {
	return &sessionRepository{redis: redisClient}
}

// StoreSession saves a token/session in Redis with an automatic TTL
func (r *sessionRepository) StoreSession(ctx context.Context, sessionID string, userID uint, expiration time.Duration) error {
	// e.g., SET session:xyz123 45 EX 3600
	return r.redis.Set(ctx, "session:"+sessionID, userID, expiration).Err()
}

// GetUserIDBySession retrieves the user ID to validate a request
func (r *sessionRepository) GetUserIDBySession(ctx context.Context, sessionID string) (uint, error) {
	// e.g., GET session:xyz123
	userIDStr, err := r.redis.Get(ctx, "session:"+sessionID).Uint64()
	if err != nil {
		return 0, err // Returns redis.Nil if the session doesn't exist or expired
	}
	return uint(userIDStr), nil
}

// IncrementAnalytics is a lightning-fast counter for API metrics
func (r *sessionRepository) IncrementAnalytics(ctx context.Context, metricKey string) error {
	// e.g., INCR api:logins:today
	return r.redis.Incr(ctx, "analytics:"+metricKey).Err()
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *dto.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	ttl := time.Until(session.ExpiresAt)
	return r.redis.Set(ctx, "session:"+session.JTI, data, ttl).Err()
}

func (r *sessionRepository) GetSessionByJTI(ctx context.Context, jti string) (*dto.Session, error) {
	data, err := r.redis.Get(ctx, "session:"+jti).Bytes()
	if err != nil {
		return nil, err // Returns redis.Nil if not found (logged out / revoked)
	}

	var session dto.Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// GetActiveSessionsByUserID retrieves all active sessions for a given user from Redis
func (r *sessionRepository) GetActiveSessionsByUserID(ctx context.Context, userID uint) ([]dto.Session, error) {
	var sessions []dto.Session

	// Use SCAN to iterate through all keys matching pattern "session:*"
	// This is more efficient than KEYS for large datasets
	iter := r.redis.Scan(ctx, 0, "session:*", 0).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()

		// Get session data from Redis
		data, err := r.redis.Get(ctx, key).Bytes()
		if err != nil {
			// Skip if we can't retrieve this session
			continue
		}

		var session dto.Session
		if err := json.Unmarshal(data, &session); err != nil {
			// Skip malformed sessions
			continue
		}

		// Only include sessions that match the userID and are active
		if session.UserID == userID && session.IsActive {
			sessions = append(sessions, session)
		}
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}
