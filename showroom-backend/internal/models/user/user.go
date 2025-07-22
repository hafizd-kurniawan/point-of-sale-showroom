package user

import (
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// User represents a user in the system
type User struct {
	UserID       int                `json:"user_id" db:"user_id"`
	Username     string             `json:"username" db:"username"`
	Email        string             `json:"email" db:"email"`
	PasswordHash string             `json:"-" db:"password_hash"`
	FullName     string             `json:"full_name" db:"full_name"`
	Phone        string             `json:"phone" db:"phone"`
	Address      *string            `json:"address,omitempty" db:"address"`
	Role         common.UserRole    `json:"role" db:"role"`
	Salary       *float64           `json:"salary,omitempty" db:"salary"`
	HireDate     *time.Time         `json:"hire_date,omitempty" db:"hire_date"`
	CreatedAt    time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" db:"updated_at"`
	CreatedBy    int                `json:"created_by" db:"created_by"`
	IsActive     bool               `json:"is_active" db:"is_active"`
	ProfileImage *string            `json:"profile_image,omitempty" db:"profile_image"`
	Notes        *string            `json:"notes,omitempty" db:"notes"`
	Creator      *UserCreatorInfo   `json:"creator,omitempty" db:"-"`
}

// UserCreatorInfo represents minimal creator information
type UserCreatorInfo struct {
	UserID   int    `json:"user_id" db:"creator_user_id"`
	Username string `json:"username" db:"creator_username"`
	FullName string `json:"full_name" db:"creator_full_name"`
}

// UserSession represents a user session
type UserSession struct {
	SessionID    int        `json:"session_id" db:"session_id"`
	UserID       int        `json:"user_id" db:"user_id"`
	SessionToken string     `json:"-" db:"session_token"`
	LoginAt      time.Time  `json:"login_at" db:"login_at"`
	LogoutAt     *time.Time `json:"logout_at,omitempty" db:"logout_at"`
	IPAddress    *string    `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent    *string    `json:"user_agent,omitempty" db:"user_agent"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	Duration     string     `json:"duration,omitempty" db:"-"`
}

// UserProfile represents user profile with sessions
type UserProfile struct {
	User     User          `json:"user"`
	Sessions []UserSession `json:"sessions"`
}

// UserListItem represents a simplified user for list views
type UserListItem struct {
	UserID    int             `json:"user_id" db:"user_id"`
	Username  string          `json:"username" db:"username"`
	Email     string          `json:"email" db:"email"`
	FullName  string          `json:"full_name" db:"full_name"`
	Phone     string          `json:"phone" db:"phone"`
	Address   *string         `json:"address,omitempty" db:"address"`
	Role      common.UserRole `json:"role" db:"role"`
	Salary    *float64        `json:"salary,omitempty" db:"salary"`
	HireDate  *time.Time      `json:"hire_date,omitempty" db:"hire_date"`
	IsActive  bool            `json:"is_active" db:"is_active"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
	CreatedBy int             `json:"created_by" db:"created_by"`
}

// UserCreateRequest represents a request to create a user
type UserCreateRequest struct {
	Username     string          `json:"username" binding:"required,min=3,max=50"`
	Email        string          `json:"email" binding:"required,email,max=100"`
	Password     string          `json:"password" binding:"required,min=6,max=100"`
	FullName     string          `json:"full_name" binding:"required,max=255"`
	Phone        string          `json:"phone" binding:"required,max=20"`
	Address      *string         `json:"address,omitempty" binding:"omitempty,max=500"`
	Role         common.UserRole `json:"role" binding:"required"`
	Salary       *float64        `json:"salary,omitempty"`
	HireDate     *string         `json:"hire_date,omitempty"`
	ProfileImage *string         `json:"profile_image,omitempty" binding:"omitempty,max=500"`
	Notes        *string         `json:"notes,omitempty"`
}

// UserUpdateRequest represents a request to update a user
type UserUpdateRequest struct {
	Username     *string         `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	Email        *string         `json:"email,omitempty" binding:"omitempty,email,max=100"`
	FullName     *string         `json:"full_name,omitempty" binding:"omitempty,max=255"`
	Phone        *string         `json:"phone,omitempty" binding:"omitempty,max=20"`
	Address      *string         `json:"address,omitempty" binding:"omitempty,max=500"`
	Role         *common.UserRole `json:"role,omitempty"`
	Salary       *float64        `json:"salary,omitempty"`
	IsActive     *bool           `json:"is_active,omitempty"`
	ProfileImage *string         `json:"profile_image,omitempty" binding:"omitempty,max=500"`
	Notes        *string         `json:"notes,omitempty"`
}

// UserFilterParams represents filtering parameters for user queries
type UserFilterParams struct {
	Role     *common.UserRole `json:"role,omitempty" form:"role"`
	IsActive *bool            `json:"is_active,omitempty" form:"is_active"`
	Search   string           `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// CalculateDuration calculates session duration
func (s *UserSession) CalculateDuration() {
	if s.LogoutAt != nil {
		duration := s.LogoutAt.Sub(s.LoginAt)
		s.Duration = formatDuration(duration)
	} else {
		duration := time.Since(s.LoginAt)
		s.Duration = formatDuration(duration)
	}
}

// formatDuration formats duration to human readable format
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh%dm%ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	} else {
		return fmt.Sprintf("%ds", seconds)
	}
}