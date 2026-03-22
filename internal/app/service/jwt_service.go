package service

import (
	"crypto/ed25519"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService interface {
	GenerateToken(userID uint, role string) (string, string, time.Time, error)
	VerifyToken(tokenString string) (jwt.MapClaims, error)
}

type jwtService struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
	issuer     string
}

// NewJWTService initializes the EdDSA keys. In production, load these from secure env/vault.
func NewJWTService(privKey ed25519.PrivateKey, pubKey ed25519.PublicKey, issuer string) JWTService {
	return &jwtService{
		privateKey: privKey,
		publicKey:  pubKey,
		issuer:     issuer,
	}
}

func (s *jwtService) GenerateToken(userID uint, role string) (string, string, time.Time, error) {
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

	claims := jwt.MapClaims{
		"sub":  userID,
		"jti":  jtiString,
		"role": role,
		"exp":  jwt.NewNumericDate(exp),
		"nbf":  jwt.NewNumericDate(now),
		"iat":  jwt.NewNumericDate(now),
		"iss":  s.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	signedToken, err := token.SignedString(s.privateKey)
	
	return signedToken, jtiString, exp, err
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