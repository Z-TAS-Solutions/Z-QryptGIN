package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/cache"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type UserRegistrationService interface {
	RegisterUser(ctx context.Context, req dto.UserRegistrationDetailsRequest) (*dto.UserRegistrationDetailsResponse, error)
	VerifyOTP(ctx context.Context, req dto.UserRegistrationOTPRequest) (*dto.UserRegistrationOTPResponse, error)
	ResendOTP(ctx context.Context, req dto.ResendOTPRequest) (*dto.ResendOTPResponse, error)
}

type userRegistrationService struct {
	repo        repository.UserRepository
	redisClient *redis.Client
	email       EmailService
	rateLimiter *ratelimit.OTPRateLimiter
}

var (
	ErrRegistrationNotFound = fmt.Errorf("registration entry not found")
	ErrInvalidOTP           = fmt.Errorf("invalid otp")
)

func NewUserRegistrationService(repo repository.UserRepository, redisClient *redis.Client, email EmailService) UserRegistrationService {
	return &userRegistrationService{
		repo:        repo,
		redisClient: redisClient,
		email:       email,
		rateLimiter: ratelimit.NewOTPRateLimiter(redisClient),
	}
}

func (s *userRegistrationService) RegisterUser(ctx context.Context, req dto.UserRegistrationDetailsRequest) (*dto.UserRegistrationDetailsResponse, error) {
	// 1) Check existing registration in Redis cache for email/phone/nic
	found, err := s.findExistingCacheEntry(ctx, req)
	if err != nil {
		return nil, err
	}

	if found != nil {
		if found.MfaStatus == database.MfaApproved {
			// already verified registration
			return &dto.UserRegistrationDetailsResponse{Success: true, CustomID: found.UserID}, nil
		}

		if s.entirelyMatchesRequest(found, req) {
			return &dto.UserRegistrationDetailsResponse{Success: true, CustomID: found.UserID}, nil
		}

		return &dto.UserRegistrationDetailsResponse{Success: false, CustomID: ""}, nil
	}

	// 2) if not in cache, check database for existing user by any key
	if existing, err := s.repo.FindByEmail(string(req.Email)); err == nil && existing != nil && existing.ID != 0 {
		return &dto.UserRegistrationDetailsResponse{Success: false, CustomID: ""}, nil
	}
	if existing, err := s.repo.FindByPhoneNo(string(req.PhoneNo)); err == nil && existing != nil && existing.ID != 0 {
		return &dto.UserRegistrationDetailsResponse{Success: false, CustomID: ""}, nil
	}
	if existing, err := s.repo.FindByNic(string(req.Nic)); err == nil && existing != nil && existing.ID != 0 {
		return &dto.UserRegistrationDetailsResponse{Success: false, CustomID: ""}, nil
	}

	// 3) Create new registration cache entry with generated user ID and OTP
	customID := generateCustomID()
	otpCode := generateOTP(6)
	cacheEntry := cache.RegistrationCache{
		ID:            fmt.Sprintf("reg-%s", customID),
		UserID:        customID,
		Name:          req.Name,
		Email:         req.Email,
		PhoneNo:       req.PhoneNo,
		Nic:           req.Nic,
		Role:          req.Role,
		OTP:           otpCode,
		MfaStatus:     database.MfaPending,
		SecurityLevel: database.SecurityLow,
		CreatedAt:     time.Now().UTC(),
	}

	if err = s.saveCacheEntry(ctx, cacheEntry); err != nil {
		return nil, err
	}

	// 4) Send OTP email
	err = s.email.SendOTPEmail(string(req.Email), string(req.Name), otpCode)
	if err != nil {
		log.Error().Err(err).Str("user_email", string(req.Email)).Msg("Failed to send registration OTP")
	} else {
		log.Info().Str("user_email", string(req.Email)).Msg("Registration OTP sent Successfully")
	}

	return &dto.UserRegistrationDetailsResponse{Success: true, CustomID: customID}, nil
}

func (s *userRegistrationService) findExistingCacheEntry(ctx context.Context, req dto.UserRegistrationDetailsRequest) (*cache.RegistrationCache, error) {
	keys := []string{
		registrationEmailKey(string(req.Email)),
		registrationPhoneKey(string(req.PhoneNo)),
		registrationNicKey(string(req.Nic)),
	}

	seen := map[string]struct{}{}
	for _, key := range keys {
		val, err := s.redisClient.Get(ctx, key).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			return nil, err
		}

		var entry cache.RegistrationCache
		if err := json.Unmarshal([]byte(val), &entry); err != nil {
			// If the cache value is malformed (e.g., non-JSON payload), remove stale key and continue lookup
			log.Warn().Err(err).Str("key", key).Msg("Malformed user registration cache entry, deleting stale value")
			_ = s.redisClient.Del(ctx, key).Err()
			continue
		}

		if entry.UserID == "" {
			continue
		}

		if _, exists := seen[string(entry.UserID)]; exists {
			continue
		}
		seen[string(entry.UserID)] = struct{}{}

		return &entry, nil
	}

	return nil, nil
}

func (s *userRegistrationService) findCacheEntryByUserID(ctx context.Context, userID database.UserCustomID) (*cache.RegistrationCache, error) {
	val, err := s.redisClient.Get(ctx, registrationUserIDKey(string(userID))).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var entry cache.RegistrationCache
	if err := json.Unmarshal([]byte(val), &entry); err != nil {
		log.Warn().Err(err).Str("user_id", string(userID)).Msg("Malformed user registration cache entry, deleting stale value")
		_ = s.redisClient.Del(ctx, registrationUserIDKey(string(userID))).Err()
		return nil, nil
	}

	return &entry, nil
}

func (s *userRegistrationService) VerifyOTP(ctx context.Context, req dto.UserRegistrationOTPRequest) (*dto.UserRegistrationOTPResponse, error) {
	entry, err := s.findCacheEntryByUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, ErrRegistrationNotFound
	}

	if entry.OTP != req.OTP {
		return nil, ErrInvalidOTP
	}

	if entry.MfaStatus == database.MfaApproved {
		return &dto.UserRegistrationOTPResponse{Success: true, Message: "OTP already verified"}, nil
	}

	entry.MfaStatus = database.MfaApproved
	if err := s.saveCacheEntry(ctx, *entry); err != nil {
		return nil, err
	}

	return &dto.UserRegistrationOTPResponse{Success: true, Message: "OTP verified successfully"}, nil
}

func (s *userRegistrationService) ResendOTP(ctx context.Context, req dto.ResendOTPRequest) (*dto.ResendOTPResponse, error) {
	// Check rate limit before proceeding
	limits := ratelimit.DefaultOTPLimits()
	if err := s.rateLimiter.CheckAndRecord(ctx, string(req.UserID), limits); err != nil {
		return &dto.ResendOTPResponse{Success: false, Message: err.Error()}, nil
	}

	// Find the registration entry by UserID
	entry, err := s.findCacheEntryByUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return &dto.ResendOTPResponse{Success: false, Message: "registration info not found or expired"}, nil
	}

	// If already verified, don't resend
	if entry.MfaStatus == database.MfaApproved {
		return &dto.ResendOTPResponse{Success: false, Message: "registration already verified"}, nil
	}

	// Generate new OTP
	newOTP := generateOTP(6)
	entry.OTP = newOTP

	// Update cache with new OTP
	if err := s.saveCacheEntry(ctx, *entry); err != nil {
		return nil, err
	}

	// Send OTP email
	err = s.email.SendOTPEmail(string(entry.Email), string(entry.Name), newOTP)
	if err != nil {
		log.Error().Err(err).Str("user_email", string(entry.Email)).Str("user_id", string(req.UserID)).Msg("Failed to send resend OTP email")
		return nil, fmt.Errorf("failed to send OTP email: %w", err)
	}

	log.Info().Str("user_email", string(entry.Email)).Str("user_id", string(req.UserID)).Msg("OTP resent successfully")
	return &dto.ResendOTPResponse{Success: true, Message: "OTP sent successfully to your email"}, nil
}

func (s *userRegistrationService) entirelyMatchesRequest(entry *cache.RegistrationCache, req dto.UserRegistrationDetailsRequest) bool {
	return entry.Name == req.Name && entry.Email == req.Email && entry.PhoneNo == req.PhoneNo && entry.Nic == req.Nic && entry.Role == req.Role
}

func (s *userRegistrationService) saveCacheEntry(ctx context.Context, entry cache.RegistrationCache) error {
	payload, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	ttl := 30 * time.Minute
	if err := s.redisClient.Set(ctx, registrationEmailKey(string(entry.Email)), payload, ttl).Err(); err != nil {
		return err
	}
	if err := s.redisClient.Set(ctx, registrationPhoneKey(string(entry.PhoneNo)), payload, ttl).Err(); err != nil {
		return err
	}
	if err := s.redisClient.Set(ctx, registrationNicKey(string(entry.Nic)), payload, ttl).Err(); err != nil {
		return err
	}
	if err := s.redisClient.Set(ctx, registrationUserIDKey(string(entry.UserID)), payload, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func registrationEmailKey(email string) string {
	return "registration:email:" + email
}

func registrationPhoneKey(phone string) string {
	return "registration:phone:" + phone
}

func registrationNicKey(nic string) string {
	return "registration:nic:" + nic
}

func registrationUserIDKey(userID string) string {
	return "registration:userid:" + userID
}

func generateOTP(max int) string {
	var table = []byte{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

	b := make([]byte, max)

	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max || err != nil {
		return "0Aj3kL"
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func generateCustomID() database.UserCustomID {
	var table = []byte{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	b := make([]byte, 6)
	n, err := io.ReadAtLeast(rand.Reader, b, 6)
	if n != 6 || err != nil {
		return database.UserCustomID("default_custom_id")
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return database.UserCustomID("USR-" + string(b))
}
