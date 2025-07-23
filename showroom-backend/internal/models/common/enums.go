package common

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleSuperAdmin UserRole = "super_admin"  // Full system management
	RoleKasir      UserRole = "kasir"        // Operational: purchase, sales, approvals
	RoleMekanik    UserRole = "mekanik"      // Repair work, damage assessment, parts usage
)

// IsValid checks if the user role is valid
func (r UserRole) IsValid() bool {
	switch r {
	case RoleSuperAdmin, RoleKasir, RoleMekanik:
		return true
	default:
		return false
	}
}

// String returns the string representation of the role
func (r UserRole) String() string {
	return string(r)
}

// BaseModel represents common fields for all models
type BaseModel struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Total      int  `json:"total"`
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	TotalPages int  `json:"total_pages"`
	HasMore    bool `json:"has_more"`
}

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

// Validate validates pagination parameters
func (p *PaginationParams) Validate() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = 10
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
}

// GetOffset calculates the offset for database queries
func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// GetTotalPages calculates total pages from total records
func (p *PaginationParams) GetTotalPages(total int) int {
	if total == 0 {
		return 0
	}
	pages := total / p.Limit
	if total%p.Limit > 0 {
		pages++
	}
	return pages
}

// GetHasMore calculates if there are more pages
func (p *PaginationParams) GetHasMore(total int) bool {
	return p.Page < p.GetTotalPages(total)
}

// Value implements the driver.Valuer interface for UserRole
func (r UserRole) Value() (driver.Value, error) {
	return string(r), nil
}

// Scan implements the sql.Scanner interface for UserRole
func (r *UserRole) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch s := value.(type) {
	case string:
		*r = UserRole(s)
	case []byte:
		*r = UserRole(s)
	default:
		return fmt.Errorf("cannot scan %T into UserRole", value)
	}
	return nil
}