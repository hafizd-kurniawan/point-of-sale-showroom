package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/user"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) interfaces.UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, u *user.User) (*user.User, error) {
	query := `
		INSERT INTO users (username, email, password_hash, full_name, phone, address, role, salary, hire_date, created_by, profile_image, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING user_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		u.Username, u.Email, u.PasswordHash, u.FullName, u.Phone, u.Address,
		u.Role, u.Salary, u.HireDate, u.CreatedBy, u.ProfileImage, u.Notes,
	).Scan(&u.UserID, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return u, nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id int) (*user.User, error) {
	query := `
		SELECT u.user_id, u.username, u.email, u.password_hash, u.full_name, u.phone, u.address,
		       u.role, u.salary, u.hire_date, u.created_at, u.updated_at, u.created_by,
		       u.is_active, u.profile_image, u.notes,
		       c.user_id as creator_user_id, c.username as creator_username, c.full_name as creator_full_name
		FROM users u
		LEFT JOIN users c ON u.created_by = c.user_id
		WHERE u.user_id = $1`

	u := &user.User{}
	creator := &user.UserCreatorInfo{}
	var creatorUserID sql.NullInt64
	var creatorUsername, creatorFullName sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.UserID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName, &u.Phone, &u.Address,
		&u.Role, &u.Salary, &u.HireDate, &u.CreatedAt, &u.UpdatedAt, &u.CreatedBy,
		&u.IsActive, &u.ProfileImage, &u.Notes,
		&creatorUserID, &creatorUsername, &creatorFullName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Set creator info if available
	if creatorUserID.Valid {
		creator.UserID = int(creatorUserID.Int64)
		creator.Username = creatorUsername.String
		creator.FullName = creatorFullName.String
		u.Creator = creator
	}

	return u, nil
}

// GetByUsername retrieves a user by username
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	query := `
		SELECT user_id, username, email, password_hash, full_name, phone, address,
		       role, salary, hire_date, created_at, updated_at, created_by,
		       is_active, profile_image, notes
		FROM users
		WHERE username = $1`

	u := &user.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&u.UserID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName, &u.Phone, &u.Address,
		&u.Role, &u.Salary, &u.HireDate, &u.CreatedAt, &u.UpdatedAt, &u.CreatedBy,
		&u.IsActive, &u.ProfileImage, &u.Notes,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return u, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT user_id, username, email, password_hash, full_name, phone, address,
		       role, salary, hire_date, created_at, updated_at, created_by,
		       is_active, profile_image, notes
		FROM users
		WHERE email = $1`

	u := &user.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.UserID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName, &u.Phone, &u.Address,
		&u.Role, &u.Salary, &u.HireDate, &u.CreatedAt, &u.UpdatedAt, &u.CreatedBy,
		&u.IsActive, &u.ProfileImage, &u.Notes,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return u, nil
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, id int, u *user.User) (*user.User, error) {
	query := `
		UPDATE users 
		SET username = $1, email = $2, full_name = $3, phone = $4, address = $5,
		    role = $6, salary = $7, is_active = $8, profile_image = $9, notes = $10,
		    updated_at = NOW()
		WHERE user_id = $11
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		u.Username, u.Email, u.FullName, u.Phone, u.Address,
		u.Role, u.Salary, u.IsActive, u.ProfileImage, u.Notes, id,
	).Scan(&u.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	u.UserID = id
	return u, nil
}

// Delete soft deletes a user
func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE users SET is_active = FALSE, updated_at = NOW() WHERE user_id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// List retrieves users with filtering and pagination
func (r *userRepository) List(ctx context.Context, params *user.UserFilterParams) ([]user.UserListItem, int, error) {
	// Build WHERE clause
	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if params.Role != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("role = $%d", argIndex))
		args = append(args, *params.Role)
		argIndex++
	}

	if params.IsActive != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *params.IsActive)
		argIndex++
	}

	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		whereConditions = append(whereConditions, fmt.Sprintf("(username ILIKE $%d OR email ILIKE $%d OR full_name ILIKE $%d OR phone ILIKE $%d)", argIndex, argIndex, argIndex, argIndex))
		args = append(args, searchPattern)
		argIndex++
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Build main query with pagination
	params.Validate()
	offset := params.GetOffset()
	
	query := fmt.Sprintf(`
		SELECT user_id, username, email, full_name, phone, address, role, salary, hire_date,
		       is_active, created_at, updated_at, created_by
		FROM users %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)
	
	args = append(args, params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []user.UserListItem
	for rows.Next() {
		var u user.UserListItem
		err := rows.Scan(
			&u.UserID, &u.Username, &u.Email, &u.FullName, &u.Phone, &u.Address,
			&u.Role, &u.Salary, &u.HireDate, &u.IsActive, &u.CreatedAt, &u.UpdatedAt, &u.CreatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}

	return users, total, nil
}

// GetByRole retrieves users by role with pagination
func (r *userRepository) GetByRole(ctx context.Context, role string, page, limit int) ([]user.UserListItem, int, error) {
	params := &user.UserFilterParams{
		Role: (*common.UserRole)(&role),
		PaginationParams: common.PaginationParams{
			Page:  page,
			Limit: limit,
		},
	}
	return r.List(ctx, params)
}

// Search searches users by query
func (r *userRepository) Search(ctx context.Context, query string, page, limit int) ([]user.UserListItem, int, error) {
	params := &user.UserFilterParams{
		Search: query,
		PaginationParams: common.PaginationParams{
			Page:  page,
			Limit: limit,
		},
	}
	return r.List(ctx, params)
}

// ExistsByUsername checks if a user exists by username
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	return exists, err
}

// ExistsByEmail checks if a user exists by email
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	return exists, err
}

// ExistsByUsernameExcludingID checks if a username exists excluding a specific user ID
func (r *userRepository) ExistsByUsernameExcludingID(ctx context.Context, username string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND user_id != $2)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, username, excludeID).Scan(&exists)
	return exists, err
}

// ExistsByEmailExcludingID checks if an email exists excluding a specific user ID
func (r *userRepository) ExistsByEmailExcludingID(ctx context.Context, email string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND user_id != $2)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, email, excludeID).Scan(&exists)
	return exists, err
}