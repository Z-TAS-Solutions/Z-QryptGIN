package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/service"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if isValidationError(err) || isJSONDecodeError(err) {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"code":    http.StatusBadRequest,
					"message": err.Error(),
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"code":    500,
				"message": "an internal server error occurred",
			})
		}
	}
}

func isValidationError(err error) bool {
	validationErrors := []error{
		database.ErrInvalidEmail,
		database.ErrInvalidPhone,
		database.ErrInvalidIPv4,
		database.ErrInvalidNIC,
		database.ErrInvalidUserID,
		database.ErrInvalidNotificationID,
		database.ErrInvalidMfaID,
		database.ErrInvalidActivityID,
		database.ErrInvalidPasskeyID,
		database.ErrInvalidSessionID,
		database.ErrInvalidStatus,
		database.ErrInvalidSecurityLevel,
		database.ErrInvalidNotifyType,
		database.ErrInvalidActivityType,
		database.ErrInvalidMfaStatus,
		database.ErrInvalidMfaDecision,
		database.ErrInvalidRole,
	}

	for _, vErr := range validationErrors {
		if errors.Is(err, vErr) {
			return true
		}
	}
	return false
}

func isJSONDecodeError(err error) bool {
	var syntaxErr *json.SyntaxError
	var unmarshalTypeErr *json.UnmarshalTypeError
	if errors.As(err, &syntaxErr) || errors.As(err, &unmarshalTypeErr) {
		return true
	}
	return false
}

// JWTAuthMiddleware verifies JWT token, validates session in Redis, and extracts user identity
// This is the core of the hybrid JWT + Redis session tracking approach
// It handles:
// 1. Token signature verification (EdDSA)
// 2. Token expiry checks
// 3. Cross-validation with Redis session cache
// 4. Extraction of 'sub' (user ID) and storing in Gin context
func JWTAuthMiddleware(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"code":    http.StatusUnauthorized,
				"message": "missing authorization header",
			})
			c.Abort()
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"code":    http.StatusUnauthorized,
				"message": "invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Step 1: Verify JWT signature and structure
		tokenClaims, err := jwtService.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"code":    http.StatusUnauthorized,
				"message": "invalid or expired token",
			})
			c.Abort()
			return
		}

		// Step 2: Extract and validate claims
		userID, ok := tokenClaims["sub"]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"code":    http.StatusUnauthorized,
				"message": "invalid token claims: missing sub",
			})
			c.Abort()
			return
		}

		// Step 3: Validate session in Redis and get session data
		// This is the hybrid approach - we verify the token is still valid in our session store
		session, err := jwtService.ValidateSessionWithToken(c.Request.Context(), tokenClaims)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"code":    http.StatusUnauthorized,
				"message": "session invalid or revoked: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Step 4: Store user identity and session in Gin context for handlers to use
		// This allows handlers to identify the user without passing tokens around
		c.Set("userID", userID)
		c.Set("session", session)
		c.Set("tokenClaims", tokenClaims)

		// Continue to next handler
		c.Next()
	}
}
