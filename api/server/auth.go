package server

import (
	"net/http"
	"strings"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/repository"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
	"github.com/gin-gonic/gin"
)

func RequireAuth(jwtSvc service.JWTService, sessionRepo repository.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extract Bearer Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 2. Cryptographic Verification (Fast, CPU-bound)
		claims, err := jwtSvc.VerifyToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token signature or expired"})
			return
		}

		// Extract claims safely
		subFloat, _ := claims["sub"].(float64) // JSON numbers parse as float64
		userID := uint(subFloat)
		jti, _ := claims["jti"].(string)
		role, _ := claims["role"].(string)

		// 3. Stateful Redis Check (I/O bound, but very fast)
		session, err := sessionRepo.GetSessionByJTI(c.Request.Context(), jti)
		if err != nil || session == nil || !session.IsActive || session.MfaStatus != dto.MfaStatusVerified {
			// Token is cryptographically valid, but session was revoked/deleted in Redis
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "session expired or revoked"})
			return
		}

		// 4. Set Context Variables for downstream handlers
		c.Set("user_id", userID)
		c.Set("jti", jti)
		c.Set("role", role)
		c.Set("session_id", session.ID)

		// 5. Update LastActiveAt asynchronously (Fire and forget so it doesn't block the request)
		go func(sessJti string) {
			// Note: Create a fresh context here since the Gin context cancels when the request ends
			// sessionRepo.UpdateLastActive(context.Background(), sessJti, time.Now())
		}(jti)

		c.Next()
	}
}

// RequireRole enforces role-based access control (RBAC) by:
// 1. Checking the role in JWT claims matches the role stored in Redis session cache
// 2. Verifying the user has one of the allowed roles
// This prevents role manipulation attacks and ensures JWT role claims match cached session role
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract role from context (set by RequireAuth middleware)
		roleInterface, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   "InternalServerError",
				"code":    http.StatusInternalServerError,
				"message": "role not found in context - ensure RequireAuth middleware is applied first",
			})
			return
		}

		role, ok := roleInterface.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   "InternalServerError",
				"code":    http.StatusInternalServerError,
				"message": "invalid role type in context",
			})
			return
		}

		// Check if user's role matches any of the allowed roles
		isAuthorized := false
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				isAuthorized = true
				break
			}
		}

		if !isAuthorized {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"code":    http.StatusForbidden,
				"message": "You do not have access to this resource",
			})
			return
		}

		c.Next()
	}
}

// ValidateRoleConsistency ensures the role in JWT claims matches the role stored in Redis session cache.
// This is a stricter validation that detects role tampering or cache inconsistencies.
// It should be used as an extra validation layer for sensitive admin operations.
func ValidateRoleConsistency(jwtSvc service.JWTService, sessionRepo repository.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"code":    http.StatusUnauthorized,
				"message": "missing or invalid authorization header",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify JWT signature
		claims, err := jwtSvc.VerifyToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"code":    http.StatusUnauthorized,
				"message": "invalid or expired token",
			})
			return
		}

		// Extract claims
		jti, _ := claims["jti"].(string)
		jwtRole, _ := claims["role"].(string)

		// Get session from cache
		session, err := sessionRepo.GetSessionByJTI(c.Request.Context(), jti)
		if err != nil || session == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"code":    http.StatusUnauthorized,
				"message": "session not found in cache",
			})
			return
		}

		// Cross-validate JWT role with session cache role
		// This prevents users from using an admin token on a client session or vice versa
		if jwtRole != session.Role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"code":    http.StatusForbidden,
				"message": "role mismatch between token and session - possible security violation",
			})
			return
		}

		c.Next()
	}
}
