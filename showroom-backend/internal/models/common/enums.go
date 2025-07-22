package common

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleSales    UserRole = "sales"
	RoleCashier  UserRole = "cashier"
	RoleMechanic UserRole = "mechanic"
	RoleManager  UserRole = "manager"
)

// CustomerType represents the type of customer
type CustomerType string

const (
	CustomerTypeIndividual CustomerType = "individual"
	CustomerTypeCorporate  CustomerType = "corporate"
)

// SupplierType represents the type of supplier
type SupplierType string

const (
	SupplierTypeParts   SupplierType = "parts"
	SupplierTypeVehicle SupplierType = "vehicle"
	SupplierTypeBoth    SupplierType = "both"
)

// FuelType represents the fuel type of a vehicle
type FuelType string

const (
	FuelTypeGasoline FuelType = "gasoline"
	FuelTypeDiesel   FuelType = "diesel"
	FuelTypeElectric FuelType = "electric"
	FuelTypeHybrid   FuelType = "hybrid"
)

// TransmissionType represents the transmission type of a vehicle
type TransmissionType string

const (
	TransmissionManual    TransmissionType = "manual"
	TransmissionAutomatic TransmissionType = "automatic"
	TransmissionCVT       TransmissionType = "cvt"
)

// IsValid checks if the user role is valid
func (r UserRole) IsValid() bool {
	switch r {
	case RoleAdmin, RoleSales, RoleCashier, RoleMechanic, RoleManager:
		return true
	default:
		return false
	}
}

// IsValid checks if the customer type is valid
func (c CustomerType) IsValid() bool {
	switch c {
	case CustomerTypeIndividual, CustomerTypeCorporate:
		return true
	default:
		return false
	}
}

// IsValid checks if the supplier type is valid
func (s SupplierType) IsValid() bool {
	switch s {
	case SupplierTypeParts, SupplierTypeVehicle, SupplierTypeBoth:
		return true
	default:
		return false
	}
}

// IsValid checks if the fuel type is valid
func (f FuelType) IsValid() bool {
	switch f {
	case FuelTypeGasoline, FuelTypeDiesel, FuelTypeElectric, FuelTypeHybrid:
		return true
	default:
		return false
	}
}

// IsValid checks if the transmission type is valid
func (t TransmissionType) IsValid() bool {
	switch t {
	case TransmissionManual, TransmissionAutomatic, TransmissionCVT:
		return true
	default:
		return false
	}
}

// String returns the string representation of the role
func (r UserRole) String() string {
	return string(r)
}

// String returns the string representation of the customer type
func (c CustomerType) String() string {
	return string(c)
}

// String returns the string representation of the supplier type
func (s SupplierType) String() string {
	return string(s)
}

// String returns the string representation of the fuel type
func (f FuelType) String() string {
	return string(f)
}

// String returns the string representation of the transmission type
func (t TransmissionType) String() string {
	return string(t)
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

// Value implements the driver.Valuer interface for CustomerType
func (c CustomerType) Value() (driver.Value, error) {
	return string(c), nil
}

// Scan implements the sql.Scanner interface for CustomerType
func (c *CustomerType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch s := value.(type) {
	case string:
		*c = CustomerType(s)
	case []byte:
		*c = CustomerType(s)
	default:
		return fmt.Errorf("cannot scan %T into CustomerType", value)
	}
	return nil
}

// Value implements the driver.Valuer interface for SupplierType
func (s SupplierType) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan implements the sql.Scanner interface for SupplierType
func (s *SupplierType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case string:
		*s = SupplierType(v)
	case []byte:
		*s = SupplierType(v)
	default:
		return fmt.Errorf("cannot scan %T into SupplierType", value)
	}
	return nil
}

// Value implements the driver.Valuer interface for FuelType
func (f FuelType) Value() (driver.Value, error) {
	return string(f), nil
}

// Scan implements the sql.Scanner interface for FuelType
func (f *FuelType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch s := value.(type) {
	case string:
		*f = FuelType(s)
	case []byte:
		*f = FuelType(s)
	default:
		return fmt.Errorf("cannot scan %T into FuelType", value)
	}
	return nil
}

// Value implements the driver.Valuer interface for TransmissionType
func (t TransmissionType) Value() (driver.Value, error) {
	return string(t), nil
}

// Scan implements the sql.Scanner interface for TransmissionType
func (t *TransmissionType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch s := value.(type) {
	case string:
		*t = TransmissionType(s)
	case []byte:
		*t = TransmissionType(s)
	default:
		return fmt.Errorf("cannot scan %T into TransmissionType", value)
	}
	return nil
}