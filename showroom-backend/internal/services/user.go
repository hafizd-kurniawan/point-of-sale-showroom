package services

import (
	"context"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/user"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/utils"
)

// UserService handles user management business logic
type UserService struct {
	userRepo    interfaces.UserRepository
	sessionRepo interfaces.UserSessionRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo interfaces.UserRepository, sessionRepo interfaces.UserSessionRepository) *UserService {
	return &UserService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *user.UserCreateRequest, createdBy int) (*user.User, error) {
	// Validate role
	if !req.Role.IsValid() {
		return nil, fmt.Errorf("invalid role: %s", req.Role)
	}

	// Check if username already exists
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("username already exists")
	}

	// Check if email already exists
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email already exists")
	}

	// Validate password
	if err := utils.IsValidPassword(req.Password); err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Parse hire date if provided
	var hireDate *time.Time
	if req.HireDate != nil && *req.HireDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.HireDate)
		if err != nil {
			return nil, fmt.Errorf("invalid hire date format, use YYYY-MM-DD")
		}
		hireDate = &parsed
	}

	// Create user object
	newUser := &user.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Phone:        req.Phone,
		Address:      req.Address,
		Role:         req.Role,
		Salary:       req.Salary,
		HireDate:     hireDate,
		CreatedBy:    createdBy,
		IsActive:     true,
		ProfileImage: req.ProfileImage,
		Notes:        req.Notes,
	}

	// Save user
	return s.userRepo.Create(ctx, newUser)
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id int) (*user.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, id int, req *user.UserUpdateRequest) (*user.User, error) {
	// Get existing user
	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Validate role if provided
	if req.Role != nil && !req.Role.IsValid() {
		return nil, fmt.Errorf("invalid role: %s", *req.Role)
	}

	// Check username uniqueness if changed
	if req.Username != nil && *req.Username != existingUser.Username {
		exists, err := s.userRepo.ExistsByUsernameExcludingID(ctx, *req.Username, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check username: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("username already exists")
		}
		existingUser.Username = *req.Username
	}

	// Check email uniqueness if changed
	if req.Email != nil && *req.Email != existingUser.Email {
		exists, err := s.userRepo.ExistsByEmailExcludingID(ctx, *req.Email, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check email: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("email already exists")
		}
		existingUser.Email = *req.Email
	}

	// Update fields if provided
	if req.FullName != nil {
		existingUser.FullName = *req.FullName
	}
	if req.Phone != nil {
		existingUser.Phone = *req.Phone
	}
	if req.Address != nil {
		existingUser.Address = req.Address
	}
	if req.Role != nil {
		existingUser.Role = *req.Role
	}
	if req.Salary != nil {
		existingUser.Salary = req.Salary
	}
	if req.IsActive != nil {
		existingUser.IsActive = *req.IsActive
	}
	if req.ProfileImage != nil {
		existingUser.ProfileImage = req.ProfileImage
	}
	if req.Notes != nil {
		existingUser.Notes = req.Notes
	}

	// Update user
	return s.userRepo.Update(ctx, id, existingUser)
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	return s.userRepo.Delete(ctx, id)
}

// ListUsers retrieves users with filtering and pagination
func (s *UserService) ListUsers(ctx context.Context, params *user.UserFilterParams) (*common.PaginatedResponse, error) {
	// Validate pagination parameters
	params.Validate()

	// Get users
	users, total, err := s.userRepo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Calculate pagination metadata
	totalPages := params.GetTotalPages(total)
	hasMore := params.GetHasMore(total)

	return &common.PaginatedResponse{
		Data:       users,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
		HasMore:    hasMore,
	}, nil
}

// GetUsersByRole retrieves users by role with pagination
func (s *UserService) GetUsersByRole(ctx context.Context, role string, page, limit int) (*common.PaginatedResponse, error) {
	users, total, err := s.userRepo.GetByRole(ctx, role, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by role: %w", err)
	}

	// Calculate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	totalPages := total / limit
	if total%limit > 0 {
		totalPages++
	}
	hasMore := page < totalPages

	return &common.PaginatedResponse{
		Data:       users,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasMore:    hasMore,
	}, nil
}

// GetUserSessions retrieves user sessions with pagination
func (s *UserService) GetUserSessions(ctx context.Context, userID int, page, limit int) (*common.PaginatedResponse, error) {
	sessions, total, err := s.sessionRepo.ListByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}

	// Clear sensitive session tokens
	for i := range sessions {
		sessions[i].SessionToken = ""
	}

	// Calculate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	totalPages := total / limit
	if total%limit > 0 {
		totalPages++
	}
	hasMore := page < totalPages

	return &common.PaginatedResponse{
		Data:       sessions,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasMore:    hasMore,
	}, nil
}

// RevokeUserSessions revokes all active sessions for a user
func (s *UserService) RevokeUserSessions(ctx context.Context, userID int) error {
	return s.sessionRepo.RevokeAllUserSessions(ctx, userID)
}