package master

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/user"
)

// Supplier represents a supplier in the system
type Supplier struct {
	SupplierID      int                    `json:"supplier_id" db:"supplier_id"`
	SupplierCode    string                 `json:"supplier_code" db:"supplier_code"`
	SupplierName    string                 `json:"supplier_name" db:"supplier_name"`
	ContactPerson   string                 `json:"contact_person" db:"contact_person"`
	Phone           string                 `json:"phone" db:"phone"`
	Email           *string                `json:"email,omitempty" db:"email"`
	Address         *string                `json:"address,omitempty" db:"address"`
	City            *string                `json:"city,omitempty" db:"city"`
	Province        *string                `json:"province,omitempty" db:"province"`
	PostalCode      *string                `json:"postal_code,omitempty" db:"postal_code"`
	TaxNumber       *string                `json:"tax_number,omitempty" db:"tax_number"`
	SupplierType    common.SupplierType    `json:"supplier_type" db:"supplier_type"`
	CreditLimit     *float64               `json:"credit_limit,omitempty" db:"credit_limit"`
	CreditTermDays  *int                   `json:"credit_term_days,omitempty" db:"credit_term_days"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at" db:"updated_at"`
	CreatedBy       int                    `json:"created_by" db:"created_by"`
	IsActive        bool                   `json:"is_active" db:"is_active"`
	Notes           *string                `json:"notes,omitempty" db:"notes"`
	Creator         *user.UserCreatorInfo  `json:"creator,omitempty" db:"-"`
}

// SupplierCreateRequest represents a request to create a supplier
type SupplierCreateRequest struct {
	SupplierName   string              `json:"supplier_name" binding:"required,max=255"`
	ContactPerson  string              `json:"contact_person" binding:"required,max=255"`
	Phone          string              `json:"phone" binding:"required,max=20"`
	Email          *string             `json:"email,omitempty" binding:"omitempty,email,max=100"`
	Address        *string             `json:"address,omitempty"`
	City           *string             `json:"city,omitempty" binding:"omitempty,max=100"`
	Province       *string             `json:"province,omitempty" binding:"omitempty,max=100"`
	PostalCode     *string             `json:"postal_code,omitempty" binding:"omitempty,max=10"`
	TaxNumber      *string             `json:"tax_number,omitempty" binding:"omitempty,max=50"`
	SupplierType   common.SupplierType `json:"supplier_type" binding:"required"`
	CreditLimit    *float64            `json:"credit_limit,omitempty"`
	CreditTermDays *int                `json:"credit_term_days,omitempty"`
	Notes          *string             `json:"notes,omitempty"`
}

// SupplierUpdateRequest represents a request to update a supplier
type SupplierUpdateRequest struct {
	SupplierName   *string              `json:"supplier_name,omitempty" binding:"omitempty,max=255"`
	ContactPerson  *string              `json:"contact_person,omitempty" binding:"omitempty,max=255"`
	Phone          *string              `json:"phone,omitempty" binding:"omitempty,max=20"`
	Email          *string              `json:"email,omitempty" binding:"omitempty,email,max=100"`
	Address        *string              `json:"address,omitempty"`
	City           *string              `json:"city,omitempty" binding:"omitempty,max=100"`
	Province       *string              `json:"province,omitempty" binding:"omitempty,max=100"`
	PostalCode     *string              `json:"postal_code,omitempty" binding:"omitempty,max=10"`
	SupplierType   *common.SupplierType `json:"supplier_type,omitempty"`
	CreditLimit    *float64             `json:"credit_limit,omitempty"`
	CreditTermDays *int                 `json:"credit_term_days,omitempty"`
	IsActive       *bool                `json:"is_active,omitempty"`
	Notes          *string              `json:"notes,omitempty"`
}

// SupplierFilterParams represents filtering parameters for supplier queries
type SupplierFilterParams struct {
	SupplierType *common.SupplierType `json:"supplier_type,omitempty" form:"supplier_type"`
	IsActive     *bool                `json:"is_active,omitempty" form:"is_active"`
	Search       string               `json:"search,omitempty" form:"search"`
	City         string               `json:"city,omitempty" form:"city"`
	Province     string               `json:"province,omitempty" form:"province"`
	common.PaginationParams
}

// SupplierListItem represents a simplified supplier for list views
type SupplierListItem struct {
	SupplierID    int                 `json:"supplier_id" db:"supplier_id"`
	SupplierCode  string              `json:"supplier_code" db:"supplier_code"`
	SupplierName  string              `json:"supplier_name" db:"supplier_name"`
	ContactPerson string              `json:"contact_person" db:"contact_person"`
	Phone         string              `json:"phone" db:"phone"`
	Email         *string             `json:"email,omitempty" db:"email"`
	SupplierType  common.SupplierType `json:"supplier_type" db:"supplier_type"`
	City          *string             `json:"city,omitempty" db:"city"`
	IsActive      bool                `json:"is_active" db:"is_active"`
	CreatedAt     time.Time           `json:"created_at" db:"created_at"`
}