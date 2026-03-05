package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/baixuejie/key-management-tool/backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims payload.
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken signs a JWT for the specified user.
func GenerateToken(userID uint, username string, rememberMe bool) (string, error) {
	cfg := config.GetConfig()
	if cfg == nil {
		return "", errors.New("config not initialized")
	}
	if userID == 0 {
		return "", errors.New("user_id cannot be empty")
	}
	if username == "" {
		return "", errors.New("username cannot be empty")
	}

	expiryHours := cfg.JWT.ExpiryHours
	if rememberMe {
		expiryHours *= 4
	}
	now := time.Now()
	expiryTime := now.Add(time.Duration(expiryHours) * time.Hour)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

// ValidateToken validates and parses a JWT token.
func ValidateToken(tokenString string) (*Claims, error) {
	cfg := config.GetConfig()
	if cfg == nil {
		return nil, errors.New("config not initialized")
	}
	if tokenString == "" {
		return nil, errors.New("token string cannot be empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWT.Secret), nil
	})
	if err != nil {
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

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
