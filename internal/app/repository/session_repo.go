package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// SessionRepository handles operations for session tracking and analytics
type SessionRepository interface {
	StoreSession(ctx context.Context, sessionID string, userID uint, expiration time.Duration) error
	GetUserIDBySession(ctx context.Context, sessionID string) (uint, error)
	IncrementAnalytics(ctx context.Context, metricKey string) error
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