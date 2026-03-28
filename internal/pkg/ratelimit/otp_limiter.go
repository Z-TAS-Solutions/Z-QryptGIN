package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// OTPRateLimiter enforces rate limits on OTP requests
// Rules: 1 OTP per 2 minutes, 3 OTPs per 12 hours
type OTPRateLimiter struct {
	redisClient *redis.Client
	userPrefix  string
}

// Limits defines the rate limiting rules
type OTPLimits struct {
	MinIntervalBetweenRequests time.Duration // Minimum time between requests (2 minutes)
	MaxRequestsPer12Hours      int           // Maximum requests per 12 hours (3 OTPs)
	WindowDuration             time.Duration // Sliding window duration (12 hours)
}

// DefaultOTPLimits returns the default OTP rate limiting rules
func DefaultOTPLimits() OTPLimits {
	return OTPLimits{
		MinIntervalBetweenRequests: 2 * time.Minute,
		MaxRequestsPer12Hours:      3,
		WindowDuration:             12 * time.Hour,
	}
}

// NewOTPRateLimiter creates a new OTP rate limiter
func NewOTPRateLimiter(redisClient *redis.Client) *OTPRateLimiter {
	return &OTPRateLimiter{
		redisClient: redisClient,
		userPrefix:  "ratelimit:otp:",
	}
}

// CheckAndRecord checks if a user can request an OTP, and if allowed, records the request timestamp
// Returns error if rate limit is exceeded
func (l *OTPRateLimiter) CheckAndRecord(ctx context.Context, userID string, limits OTPLimits) error {
    key := l.userPrefix + userID
    now := time.Now().UTC()
    timestamp := now.Unix()

    // Get all request timestamps in the 12-hour window using score range
    scores, err := l.redisClient.ZRangeByScore(ctx, key, &redis.ZRangeBy{
        Min: fmt.Sprintf("%d", now.Add(-limits.WindowDuration).Unix()),
        Max: fmt.Sprintf("%d", timestamp),
    }).Result()
    if err != nil {
        return fmt.Errorf("failed to check rate limit: %w", err)
    }

    // Clean up old entries outside the window
    l.redisClient.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%d", now.Add(-limits.WindowDuration).Unix()))

    // Check if max requests per 12 hours exceeded
    if len(scores) >= limits.MaxRequestsPer12Hours {
        return fmt.Errorf("rate limit exceeded: maximum %d OTP requests per 12 hours", limits.MaxRequestsPer12Hours)
    }

    // Check if minimum interval between requests has elapsed
    if len(scores) > 0 {
        // Get the highest score (most recent request) using ZRevRange
        lastScores, err := l.redisClient.ZRevRange(ctx, key, 0, 0).Result()
        if err != nil && err != redis.Nil {
            return fmt.Errorf("failed to get last request time: %w", err)
        }

        if len(lastScores) > 0 {
            lastScore, err := l.redisClient.ZScore(ctx, key, lastScores[0]).Result()
            if err != nil && err != redis.Nil {
                return fmt.Errorf("failed to get score: %w", err)
            }

            if err == nil {
                lastRequestTime := time.Unix(int64(lastScore), 0)
                if now.Before(lastRequestTime.Add(limits.MinIntervalBetweenRequests)) {
                    timeUntilNext := lastRequestTime.Add(limits.MinIntervalBetweenRequests).Sub(now)
                    return fmt.Errorf("rate limit exceeded: please wait %v before requesting another OTP", timeUntilNext)
                }
            }
        }
    }

    // Record this request using timestamp as both score and member
    uniqueMember := fmt.Sprintf("%d-%d", timestamp, now.Nanosecond())
    if err := l.redisClient.ZAdd(ctx, key, redis.Z{
        Score:  float64(timestamp),
        Member: uniqueMember,
    }).Err(); err != nil {
        return fmt.Errorf("failed to record request: %w", err)
    }

    // Set expiry on the sorted set to clean up old data
    l.redisClient.Expire(ctx, key, limits.WindowDuration)

    return nil
}

// Reset clears all rate limit records for a user (useful for testing or manual resets)
func (l *OTPRateLimiter) Reset(ctx context.Context, userID string) error {
	key := l.userPrefix + userID
	return l.redisClient.Del(ctx, key).Err()
}
