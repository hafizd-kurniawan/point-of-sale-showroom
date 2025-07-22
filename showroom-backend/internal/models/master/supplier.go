package master

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// SupplierType represents the type of supplier
type SupplierType string

const (
	SupplierTypeParts   SupplierType = "parts"
	SupplierTypeVehicle SupplierType = "vehicle"
	SupplierTypeBoth    SupplierType = "both"
)

// IsValid checks if the supplier type is valid
func (st SupplierType) IsValid() bool {
	switch st {
	case SupplierTypeParts, SupplierTypeVehicle, SupplierTypeBoth:
		return true
	default:
		return false
	}
}

// String returns the string representation of the supplier type
func (st SupplierType) String() string {
	return string(st)
}

// Supplier represents a supplier in the system
type Supplier struct {
	SupplierID   int          `json:"supplier_id" db:"supplier_id"`
	SupplierCode string       `json:"supplier_code" db:"supplier_code"`
	SupplierName string       `json:"supplier_name" db:"supplier_name"`
	SupplierType SupplierType `json:"supplier_type" db:"supplier_type"`
	Phone        string       `json:"phone" db:"phone"`
	Email        *string      `json:"email,omitempty" db:"email"`
	Address      string       `json:"address" db:"address"`
	City         string       `json:"city" db:"city"`
	PostalCode   *string      `json:"postal_code,omitempty" db:"postal_code"`
	TaxNumber    *string      `json:"tax_number,omitempty" db:"tax_number"`
	ContactPerson string      `json:"contact_person" db:"contact_person"`
	BankAccount  *string      `json:"bank_account,omitempty" db:"bank_account"`
	PaymentTerms *string      `json:"payment_terms,omitempty" db:"payment_terms"`
	Notes        *string      `json:"notes,omitempty" db:"notes"`
	IsActive     bool         `json:"is_active" db:"is_active"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
	CreatedBy    int          `json:"created_by" db:"created_by"`
}

// SupplierListItem represents a simplified supplier for list views
type SupplierListItem struct {
	SupplierID   int          `json:"supplier_id" db:"supplier_id"`
	SupplierCode string       `json:"supplier_code" db:"supplier_code"`
	SupplierName string       `json:"supplier_name" db:"supplier_name"`
	SupplierType SupplierType `json:"supplier_type" db:"supplier_type"`
	Phone        string       `json:"phone" db:"phone"`
	Email        *string      `json:"email,omitempty" db:"email"`
	City         string       `json:"city" db:"city"`
	ContactPerson string      `json:"contact_person" db:"contact_person"`
	IsActive     bool         `json:"is_active" db:"is_active"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
}

// SupplierCreateRequest represents a request to create a supplier
type SupplierCreateRequest struct {
	SupplierName  string       `json:"supplier_name" binding:"required,max=255"`
	SupplierType  SupplierType `json:"supplier_type" binding:"required"`
	Phone         string       `json:"phone" binding:"required,max=20"`
	Email         *string      `json:"email,omitempty" binding:"omitempty,email,max=100"`
	Address       string       `json:"address" binding:"required,max=500"`
	City          string       `json:"city" binding:"required,max=100"`
	PostalCode    *string      `json:"postal_code,omitempty" binding:"omitempty,max=10"`
	TaxNumber     *string      `json:"tax_number,omitempty" binding:"omitempty,max=30"`
	ContactPerson string       `json:"contact_person" binding:"required,max=255"`
	BankAccount   *string      `json:"bank_account,omitempty" binding:"omitempty,max=100"`
	PaymentTerms  *string      `json:"payment_terms,omitempty" binding:"omitempty,max=255"`
	Notes         *string      `json:"notes,omitempty"`
}

// SupplierUpdateRequest represents a request to update a supplier
type SupplierUpdateRequest struct {
	SupplierName  *string       `json:"supplier_name,omitempty" binding:"omitempty,max=255"`
	SupplierType  *SupplierType `json:"supplier_type,omitempty"`
	Phone         *string       `json:"phone,omitempty" binding:"omitempty,max=20"`
	Email         *string       `json:"email,omitempty" binding:"omitempty,email,max=100"`
	Address       *string       `json:"address,omitempty" binding:"omitempty,max=500"`
	City          *string       `json:"city,omitempty" binding:"omitempty,max=100"`
	PostalCode    *string       `json:"postal_code,omitempty" binding:"omitempty,max=10"`
	TaxNumber     *string       `json:"tax_number,omitempty" binding:"omitempty,max=30"`
	ContactPerson *string       `json:"contact_person,omitempty" binding:"omitempty,max=255"`
	BankAccount   *string       `json:"bank_account,omitempty" binding:"omitempty,max=100"`
	PaymentTerms  *string       `json:"payment_terms,omitempty" binding:"omitempty,max=255"`
	Notes         *string       `json:"notes,omitempty"`
	IsActive      *bool         `json:"is_active,omitempty"`
}

// SupplierFilterParams represents filtering parameters for supplier queries
type SupplierFilterParams struct {
	SupplierType *SupplierType `json:"supplier_type,omitempty" form:"supplier_type"`
	IsActive     *bool         `json:"is_active,omitempty" form:"is_active"`
	Search       string        `json:"search,omitempty" form:"search"`
	City         string        `json:"city,omitempty" form:"city"`
	common.PaginationParams
}