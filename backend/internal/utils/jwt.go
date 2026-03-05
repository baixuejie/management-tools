package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/baixuejie/key-management-tool/backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims structure
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for the given username
// If rememberMe is true, the token expiry is multiplied by 4 (7 days → 30 days)
func GenerateToken(username string, rememberMe bool) (string, error) {
	cfg := config.GetConfig()
	if cfg == nil {
		return "", errors.New("config not initialized")
	}

	if username == "" {
		return "", errors.New("username cannot be empty")
	}

	// Calculate expiry duration
	expiryHours := cfg.JWT.ExpiryHours
	if rememberMe {
		expiryHours *= 4
	}
	expiryTime := time.Now().Add(time.Duration(expiryHours) * time.Hour)

	// Create claims
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*Claims, error) {
	cfg := config.GetConfig()
	if cfg == nil {
		return nil, errors.New("config not initialized")
	}

	if tokenString == "" {
		return nil, errors.New("token string cannot be empty")
	}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		// Check for specific error types
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New("token is malformed")
		}
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("token signature is invalid")
		}
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
