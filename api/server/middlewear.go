package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
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
