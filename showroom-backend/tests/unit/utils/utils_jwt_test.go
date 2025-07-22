package utils_test

import (
	"testing"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/auth"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestJWTManager_GenerateToken(t *testing.T) {
	jwtManager := utils.NewJWTManager("test-secret", time.Hour)
	
	claims := &auth.TokenClaims{
		UserID:    1,
		Username:  "testuser",
		Email:     "test@example.com",
		Role:      common.RoleAdmin,
		SessionID: 123,
	}

	token, err := jwtManager.GenerateToken(claims)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJWTManager_ValidateToken(t *testing.T) {
	jwtManager := utils.NewJWTManager("test-secret", time.Hour)
	
	originalClaims := &auth.TokenClaims{
		UserID:    1,
		Username:  "testuser",
		Email:     "test@example.com",
		Role:      common.RoleAdmin,
		SessionID: 123,
	}

	token, err := jwtManager.GenerateToken(originalClaims)
	assert.NoError(t, err)

	validatedClaims, err := jwtManager.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, originalClaims.UserID, validatedClaims.UserID)
	assert.Equal(t, originalClaims.Username, validatedClaims.Username)
	assert.Equal(t, originalClaims.Email, validatedClaims.Email)
	assert.Equal(t, originalClaims.Role, validatedClaims.Role)
	assert.Equal(t, originalClaims.SessionID, validatedClaims.SessionID)
}

func TestJWTManager_ValidateToken_InvalidToken(t *testing.T) {
	jwtManager := utils.NewJWTManager("test-secret", time.Hour)
	
	_, err := jwtManager.ValidateToken("invalid-token")
	assert.Error(t, err)
}

func TestJWTManager_ValidateToken_ExpiredToken(t *testing.T) {
	jwtManager := utils.NewJWTManager("test-secret", -time.Hour) // Expired duration
	
	claims := &auth.TokenClaims{
		UserID:    1,
		Username:  "testuser",
		Email:     "test@example.com",
		Role:      common.RoleAdmin,
		SessionID: 123,
	}

	token, err := jwtManager.GenerateToken(claims)
	assert.NoError(t, err)

	_, err = jwtManager.ValidateToken(token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
}