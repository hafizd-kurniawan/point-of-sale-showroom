package utils_test

import (
	"testing"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	
	hash, err := utils.HashPassword(password)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"
	wrongPassword := "wrongpassword"
	
	hash, err := utils.HashPassword(password)
	assert.NoError(t, err)
	
	// Test correct password
	assert.True(t, utils.CheckPassword(password, hash))
	
	// Test wrong password
	assert.False(t, utils.CheckPassword(wrongPassword, hash))
}

func TestGenerateSecureToken(t *testing.T) {
	token, err := utils.GenerateSecureToken(32)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, 64, len(token)) // 32 bytes = 64 hex characters
}

func TestGenerateSessionToken(t *testing.T) {
	token, err := utils.GenerateSessionToken()
	
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, 64, len(token)) // 32 bytes = 64 hex characters
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		password string
		valid    bool
	}{
		{"123", false},        // Too short
		{"password123", true}, // Valid
		{"", false},           // Empty
		{string(make([]byte, 101)), false}, // Too long
		{"validpass", true},   // Valid
	}
	
	for _, test := range tests {
		err := utils.IsValidPassword(test.password)
		if test.valid {
			assert.NoError(t, err, "Password %s should be valid", test.password)
		} else {
			assert.Error(t, err, "Password %s should be invalid", test.password)
		}
	}
}