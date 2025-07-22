package services

import (
	"context"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/auth"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/user"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/utils"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo    interfaces.UserRepository
	sessionRepo interfaces.UserSessionRepository
	jwtManager  *utils.JWTManager
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo interfaces.UserRepository,
	sessionRepo interfaces.UserSessionRepository,
	jwtManager *utils.JWTManager,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtManager:  jwtManager,
	}
}

// Login authenticates a user and creates a session
func (s *AuthService) Login(ctx context.Context, req *auth.LoginRequest, ipAddress, userAgent string) (*auth.LoginResponse, error) {
	// Get user by username
	foundUser, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !foundUser.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Verify password
	if !utils.CheckPassword(req.Password, foundUser.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate session token
	sessionToken, err := utils.GenerateSessionToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	// Create session
	session := &user.UserSession{
		UserID:       foundUser.UserID,
		SessionToken: sessionToken,
		IPAddress:    &ipAddress,
		UserAgent:    &userAgent,
		IsActive:     true,
	}

	session, err = s.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Generate JWT token
	claims := &auth.TokenClaims{
		UserID:    foundUser.UserID,
		Username:  foundUser.Username,
		Email:     foundUser.Email,
		Role:      foundUser.Role,
		SessionID: session.SessionID,
	}

	token, err := s.jwtManager.GenerateToken(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Calculate expiration
	duration := s.jwtManager.GetExpirationDuration()
	expiresAt := time.Now().Add(duration)

	return &auth.LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: int(duration.Seconds()),
		ExpiresAt: expiresAt,
		User:      auth.UserFromModel(foundUser),
		Message:   "Login successful",
		SessionID: session.SessionID,
	}, nil
}

// Logout terminates a user session
func (s *AuthService) Logout(ctx context.Context, sessionID int) error {
	return s.sessionRepo.UpdateLogout(ctx, sessionID)
}

// GetUserInfo retrieves current user information
func (s *AuthService) GetUserInfo(ctx context.Context, userID int) (*user.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

// GetProfile retrieves user profile with recent sessions
func (s *AuthService) GetProfile(ctx context.Context, userID int) (*auth.ProfileResponse, error) {
	// Get user info
	foundUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get recent sessions (last 5)
	sessions, err := s.sessionRepo.GetRecentByUserID(ctx, userID, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}

	// Clear sensitive session tokens
	for i := range sessions {
		sessions[i].SessionToken = ""
	}

	return &auth.ProfileResponse{
		User:     *foundUser,
		Sessions: sessions,
	}, nil
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(ctx context.Context, userID int, req *auth.ChangePasswordRequest) error {
	// Validate passwords match
	if req.NewPassword != req.ConfirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	// Validate password strength
	if err := utils.IsValidPassword(req.NewPassword); err != nil {
		return err
	}

	// Get current user
	foundUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Verify current password
	if !utils.CheckPassword(req.CurrentPassword, foundUser.PasswordHash) {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	newHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password in user object
	foundUser.PasswordHash = newHash

	// Update user in database
	_, err = s.userRepo.Update(ctx, userID, foundUser)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// RefreshToken generates a new token from an existing one
func (s *AuthService) RefreshToken(ctx context.Context, oldToken string) (*auth.RefreshTokenResponse, error) {
	// Refresh the token
	newToken, claims, err := s.jwtManager.RefreshToken(oldToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// Verify session is still active
	session, err := s.sessionRepo.GetByID(ctx, claims.SessionID)
	if err != nil || !session.IsActive {
		return nil, fmt.Errorf("session is no longer active")
	}

	// Get user info
	foundUser, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is still active
	if !foundUser.IsActive {
		return nil, fmt.Errorf("user account is deactivated")
	}

	// Calculate expiration
	duration := s.jwtManager.GetExpirationDuration()
	expiresAt := time.Now().Add(duration)

	return &auth.RefreshTokenResponse{
		Token:     newToken,
		TokenType: "Bearer",
		ExpiresIn: int(duration.Seconds()),
		ExpiresAt: expiresAt,
		User:      auth.UserFromModel(foundUser),
		Message:   "Token refreshed successfully",
		SessionID: session.SessionID,
	}, nil
}