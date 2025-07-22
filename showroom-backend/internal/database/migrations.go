package database

import (
	"database/sql"
	"log"
)

// RunMigrations executes database migrations
func RunMigrations(db *sql.DB) error {
	migrations := []string{
		createUsersTable,
		createUserSessionsTable,
		createIndexes,
	}

	for i, migration := range migrations {
		log.Printf("Running migration %d...", i+1)
		if _, err := db.Exec(migration); err != nil {
			return err
		}
		log.Printf("Migration %d completed successfully", i+1)
	}

	log.Println("All migrations completed successfully")
	return nil
}

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    address VARCHAR(500),
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin','sales','cashier','mechanic','manager')),
    salary DECIMAL(15,2),
    hire_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    profile_image VARCHAR(500),
    notes TEXT
);`

const createUserSessionsTable = `
CREATE TABLE IF NOT EXISTS user_sessions (
    session_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    session_token VARCHAR(500) UNIQUE NOT NULL,
    login_at TIMESTAMP DEFAULT NOW(),
    logout_at TIMESTAMP,
    ip_address VARCHAR(45),
    user_agent VARCHAR(500),
    is_active BOOLEAN DEFAULT TRUE
);`

const createIndexes = `
-- Users table indexes
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_users_created_by ON users(created_by);

-- User sessions table indexes
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_token ON user_sessions(session_token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_is_active ON user_sessions(is_active);
CREATE INDEX IF NOT EXISTS idx_user_sessions_login_at ON user_sessions(login_at);

-- Foreign key constraint for users.created_by (self-referencing)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_users_created_by'
    ) THEN
        ALTER TABLE users ADD CONSTRAINT fk_users_created_by 
        FOREIGN KEY (created_by) REFERENCES users(user_id);
    END IF;
END $$;`