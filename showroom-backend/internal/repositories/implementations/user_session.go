package implementations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/user"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

type userSessionRepository struct {
	db *sql.DB
}

// NewUserSessionRepository creates a new user session repository
func NewUserSessionRepository(db *sql.DB) interfaces.UserSessionRepository {
	return &userSessionRepository{db: db}
}

// Create creates a new user session
func (r *userSessionRepository) Create(ctx context.Context, session *user.UserSession) (*user.UserSession, error) {
	query := `
		INSERT INTO user_sessions (user_id, session_token, ip_address, user_agent)
		VALUES ($1, $2, $3, $4)
		RETURNING session_id, login_at`

	err := r.db.QueryRowContext(ctx, query,
		session.UserID, session.SessionToken, session.IPAddress, session.UserAgent,
	).Scan(&session.SessionID, &session.LoginAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	session.IsActive = true
	return session, nil
}

// GetByToken retrieves a session by token
func (r *userSessionRepository) GetByToken(ctx context.Context, token string) (*user.UserSession, error) {
	query := `
		SELECT session_id, user_id, session_token, login_at, logout_at, ip_address, user_agent, is_active
		FROM user_sessions
		WHERE session_token = $1 AND is_active = TRUE`

	session := &user.UserSession{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&session.SessionID, &session.UserID, &session.SessionToken,
		&session.LoginAt, &session.LogoutAt, &session.IPAddress,
		&session.UserAgent, &session.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	session.CalculateDuration()
	return session, nil
}

// GetByID retrieves a session by ID
func (r *userSessionRepository) GetByID(ctx context.Context, id int) (*user.UserSession, error) {
	query := `
		SELECT session_id, user_id, session_token, login_at, logout_at, ip_address, user_agent, is_active
		FROM user_sessions
		WHERE session_id = $1`

	session := &user.UserSession{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&session.SessionID, &session.UserID, &session.SessionToken,
		&session.LoginAt, &session.LogoutAt, &session.IPAddress,
		&session.UserAgent, &session.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	session.CalculateDuration()
	return session, nil
}

// GetActiveByUserID retrieves active sessions for a user
func (r *userSessionRepository) GetActiveByUserID(ctx context.Context, userID int) ([]user.UserSession, error) {
	query := `
		SELECT session_id, user_id, session_token, login_at, logout_at, ip_address, user_agent, is_active
		FROM user_sessions
		WHERE user_id = $1 AND is_active = TRUE
		ORDER BY login_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sessions: %w", err)
	}
	defer rows.Close()

	var sessions []user.UserSession
	for rows.Next() {
		var session user.UserSession
		err := rows.Scan(
			&session.SessionID, &session.UserID, &session.SessionToken,
			&session.LoginAt, &session.LogoutAt, &session.IPAddress,
			&session.UserAgent, &session.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		session.CalculateDuration()
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// GetRecentByUserID retrieves recent sessions for a user
func (r *userSessionRepository) GetRecentByUserID(ctx context.Context, userID int, limit int) ([]user.UserSession, error) {
	query := `
		SELECT session_id, user_id, session_token, login_at, logout_at, ip_address, user_agent, is_active
		FROM user_sessions
		WHERE user_id = $1
		ORDER BY login_at DESC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent sessions: %w", err)
	}
	defer rows.Close()

	var sessions []user.UserSession
	for rows.Next() {
		var session user.UserSession
		err := rows.Scan(
			&session.SessionID, &session.UserID, &session.SessionToken,
			&session.LoginAt, &session.LogoutAt, &session.IPAddress,
			&session.UserAgent, &session.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		session.CalculateDuration()
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// UpdateLogout updates session logout time
func (r *userSessionRepository) UpdateLogout(ctx context.Context, sessionID int) error {
	query := `
		UPDATE user_sessions 
		SET logout_at = NOW(), is_active = FALSE 
		WHERE session_id = $1`

	result, err := r.db.ExecContext(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to update session logout: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

// RevokeAllUserSessions revokes all active sessions for a user
func (r *userSessionRepository) RevokeAllUserSessions(ctx context.Context, userID int) error {
	query := `
		UPDATE user_sessions 
		SET logout_at = NOW(), is_active = FALSE 
		WHERE user_id = $1 AND is_active = TRUE`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke user sessions: %w", err)
	}

	return nil
}

// ListByUserID retrieves sessions for a user with pagination
func (r *userSessionRepository) ListByUserID(ctx context.Context, userID int, page, limit int) ([]user.UserSession, int, error) {
	// Count total sessions
	countQuery := `SELECT COUNT(*) FROM user_sessions WHERE user_id = $1`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count sessions: %w", err)
	}

	// Calculate offset
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	// Get sessions with pagination
	query := `
		SELECT session_id, user_id, session_token, login_at, logout_at, ip_address, user_agent, is_active
		FROM user_sessions
		WHERE user_id = $1
		ORDER BY login_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list sessions: %w", err)
	}
	defer rows.Close()

	var sessions []user.UserSession
	for rows.Next() {
		var session user.UserSession
		err := rows.Scan(
			&session.SessionID, &session.UserID, &session.SessionToken,
			&session.LoginAt, &session.LogoutAt, &session.IPAddress,
			&session.UserAgent, &session.IsActive,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan session: %w", err)
		}
		session.CalculateDuration()
		sessions = append(sessions, session)
	}

	return sessions, total, nil
}

// DeleteExpiredSessions deletes expired sessions
func (r *userSessionRepository) DeleteExpiredSessions(ctx context.Context) error {
	// Delete sessions that have been inactive for more than 30 days
	query := `
		DELETE FROM user_sessions 
		WHERE is_active = FALSE 
		AND logout_at < NOW() - INTERVAL '30 days'`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	return nil
}

// DeleteInactiveSessions deletes inactive sessions older than specified days
func (r *userSessionRepository) DeleteInactiveSessions(ctx context.Context, days int) error {
	query := `
		DELETE FROM user_sessions 
		WHERE is_active = FALSE 
		AND logout_at < NOW() - INTERVAL '%d days'`

	_, err := r.db.ExecContext(ctx, fmt.Sprintf(query, days))
	if err != nil {
		return fmt.Errorf("failed to delete inactive sessions: %w", err)
	}

	return nil
}