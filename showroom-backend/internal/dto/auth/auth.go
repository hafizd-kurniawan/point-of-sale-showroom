package auth

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/user"
)

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response payload
type LoginResponse struct {
	Token     string              `json:"token"`
	TokenType string              `json:"token_type"`
	ExpiresIn int                 `json:"expires_in"`
	ExpiresAt time.Time           `json:"expires_at"`
	User      UserInfo            `json:"user"`
	Message   string              `json:"message"`
	SessionID int                 `json:"session_id"`
}

// UserInfo represents basic user information in responses
type UserInfo struct {
	UserID       int             `json:"user_id"`
	Username     string          `json:"username"`
	Email        string          `json:"email"`
	FullName     string          `json:"full_name"`
	Role         common.UserRole `json:"role"`
	IsActive     bool            `json:"is_active"`
	ProfileImage *string         `json:"profile_image"`
}

// LogoutResponse represents the logout response payload
type LogoutResponse struct {
	Message string `json:"message"`
}

// RefreshTokenResponse represents the refresh token response payload
type RefreshTokenResponse struct {
	Token     string    `json:"token"`
	TokenType string    `json:"token_type"`
	ExpiresIn int       `json:"expires_in"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserInfo  `json:"user"`
	Message   string    `json:"message"`
	SessionID int       `json:"session_id"`
}

// ChangePasswordRequest represents the change password request payload
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6,max=100"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

// ChangePasswordResponse represents the change password response payload
type ChangePasswordResponse struct {
	Message string `json:"message"`
}

// ProfileResponse represents the profile response payload
type ProfileResponse struct {
	User     user.User          `json:"user"`
	Sessions []user.UserSession `json:"sessions"`
}

// TokenClaims represents the JWT token claims
type TokenClaims struct {
	UserID    int             `json:"user_id"`
	Username  string          `json:"username"`
	Email     string          `json:"email"`
	Role      common.UserRole `json:"role"`
	SessionID int             `json:"session_id"`
	IssuedAt  time.Time       `json:"iat"`
	ExpiresAt time.Time       `json:"exp"`
}

// UserFromModel converts user model to UserInfo DTO
func UserFromModel(u *user.User) UserInfo {
	return UserInfo{
		UserID:       u.UserID,
		Username:     u.Username,
		Email:        u.Email,
		FullName:     u.FullName,
		Role:         u.Role,
		IsActive:     u.IsActive,
		ProfileImage: u.ProfileImage,
	}
}