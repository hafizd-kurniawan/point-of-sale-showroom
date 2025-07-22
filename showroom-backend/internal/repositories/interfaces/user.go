package interfaces

import (
	"context"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/user"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, user *user.User) (*user.User, error)
	GetByID(ctx context.Context, id int) (*user.User, error)
	GetByUsername(ctx context.Context, username string) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	Update(ctx context.Context, id int, user *user.User) (*user.User, error)
	Delete(ctx context.Context, id int) error
	
	// List and filtering operations
	List(ctx context.Context, params *user.UserFilterParams) ([]user.UserListItem, int, error)
	GetByRole(ctx context.Context, role string, page, limit int) ([]user.UserListItem, int, error)
	
	// Search operations
	Search(ctx context.Context, query string, page, limit int) ([]user.UserListItem, int, error)
	
	// Existence checks
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsernameExcludingID(ctx context.Context, username string, excludeID int) (bool, error)
	ExistsByEmailExcludingID(ctx context.Context, email string, excludeID int) (bool, error)
}

// UserSessionRepository defines the interface for user session operations
type UserSessionRepository interface {
	// Session management
	Create(ctx context.Context, session *user.UserSession) (*user.UserSession, error)
	GetByToken(ctx context.Context, token string) (*user.UserSession, error)
	GetByID(ctx context.Context, id int) (*user.UserSession, error)
	GetActiveByUserID(ctx context.Context, userID int) ([]user.UserSession, error)
	GetRecentByUserID(ctx context.Context, userID int, limit int) ([]user.UserSession, error)
	UpdateLogout(ctx context.Context, sessionID int) error
	RevokeAllUserSessions(ctx context.Context, userID int) error
	
	// List operations with pagination
	ListByUserID(ctx context.Context, userID int, page, limit int) ([]user.UserSession, int, error)
	
	// Cleanup operations
	DeleteExpiredSessions(ctx context.Context) error
	DeleteInactiveSessions(ctx context.Context, days int) error
}