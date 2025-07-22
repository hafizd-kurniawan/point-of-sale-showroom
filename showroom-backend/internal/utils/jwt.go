package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/auth"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// JWTManager handles JWT token operations
type JWTManager struct {
	secretKey string
	duration  time.Duration
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey string, duration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey: secretKey,
		duration:  duration,
	}
}

// GenerateToken generates a new JWT token
func (manager *JWTManager) GenerateToken(claims *auth.TokenClaims) (string, error) {
	// Set standard claims
	now := time.Now()
	claims.IssuedAt = now
	claims.ExpiresAt = now.Add(manager.duration)

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    claims.UserID,
		"username":   claims.Username,
		"email":      claims.Email,
		"role":       claims.Role,
		"session_id": claims.SessionID,
		"iat":        claims.IssuedAt.Unix(),
		"exp":        claims.ExpiresAt.Unix(),
	})

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(manager.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates and parses a JWT token
func (manager *JWTManager) ValidateToken(tokenString string) (*auth.TokenClaims, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(manager.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if token is valid
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Check expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, ErrExpiredToken
		}
	}

	// Parse claims into struct
	tokenClaims := &auth.TokenClaims{}

	if userID, ok := claims["user_id"].(float64); ok {
		tokenClaims.UserID = int(userID)
	}

	if username, ok := claims["username"].(string); ok {
		tokenClaims.Username = username
	}

	if email, ok := claims["email"].(string); ok {
		tokenClaims.Email = email
	}

	if role, ok := claims["role"].(string); ok {
		tokenClaims.Role = common.UserRole(role)
	}

	if sessionID, ok := claims["session_id"].(float64); ok {
		tokenClaims.SessionID = int(sessionID)
	}

	if iat, ok := claims["iat"].(float64); ok {
		tokenClaims.IssuedAt = time.Unix(int64(iat), 0)
	}

	if exp, ok := claims["exp"].(float64); ok {
		tokenClaims.ExpiresAt = time.Unix(int64(exp), 0)
	}

	return tokenClaims, nil
}

// RefreshToken generates a new token with updated expiration
func (manager *JWTManager) RefreshToken(oldTokenString string) (string, *auth.TokenClaims, error) {
	// Validate old token (but allow expired tokens for refresh)
	token, err := jwt.Parse(oldTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(manager.secretKey), nil
	})

	if err != nil {
		// For refresh, we allow expired tokens
		return "", nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil, ErrInvalidToken
	}

	// Create new claims
	newClaims := &auth.TokenClaims{}
	if userID, ok := claims["user_id"].(float64); ok {
		newClaims.UserID = int(userID)
	}
	if username, ok := claims["username"].(string); ok {
		newClaims.Username = username
	}
	if email, ok := claims["email"].(string); ok {
		newClaims.Email = email
	}
	if role, ok := claims["role"].(string); ok {
		newClaims.Role = common.UserRole(role)
	}
	if sessionID, ok := claims["session_id"].(float64); ok {
		newClaims.SessionID = int(sessionID)
	}

	// Generate new token
	newToken, err := manager.GenerateToken(newClaims)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate new token: %w", err)
	}

	return newToken, newClaims, nil
}

// GetExpirationDuration returns the token expiration duration
func (manager *JWTManager) GetExpirationDuration() time.Duration {
	return manager.duration
}