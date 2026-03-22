package server

import (
	"net/http"
	"strings"

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