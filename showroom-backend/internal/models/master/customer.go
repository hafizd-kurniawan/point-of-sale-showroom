package master

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// CustomerType represents the type of customer
type CustomerType string

const (
	CustomerTypeIndividual CustomerType = "individual"
	CustomerTypeCorporate  CustomerType = "corporate"
)

// IsValid checks if the customer type is valid
func (ct CustomerType) IsValid() bool {
	switch ct {
	case CustomerTypeIndividual, CustomerTypeCorporate:
		return true
	default:
		return false
	}
}

// String returns the string representation of the customer type
func (ct CustomerType) String() string {
	return string(ct)
}

// Customer represents a customer in the system
type Customer struct {
	CustomerID   int          `json:"customer_id" db:"customer_id"`
	CustomerCode string       `json:"customer_code" db:"customer_code"`
	CustomerName string       `json:"customer_name" db:"customer_name"`
	CustomerType CustomerType `json:"customer_type" db:"customer_type"`
	Phone        string       `json:"phone" db:"phone"`
	Email        *string      `json:"email,omitempty" db:"email"`
	Address      string       `json:"address" db:"address"`
	City         string       `json:"city" db:"city"`
	PostalCode   *string      `json:"postal_code,omitempty" db:"postal_code"`
	TaxNumber    *string      `json:"tax_number,omitempty" db:"tax_number"`
	ContactPerson *string     `json:"contact_person,omitempty" db:"contact_person"`
	Notes        *string      `json:"notes,omitempty" db:"notes"`
	IsActive     bool         `json:"is_active" db:"is_active"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
	CreatedBy    int          `json:"created_by" db:"created_by"`
}

// CustomerListItem represents a simplified customer for list views
type CustomerListItem struct {
	CustomerID   int          `json:"customer_id" db:"customer_id"`
	CustomerCode string       `json:"customer_code" db:"customer_code"`
	CustomerName string       `json:"customer_name" db:"customer_name"`
	CustomerType CustomerType `json:"customer_type" db:"customer_type"`
	Phone        string       `json:"phone" db:"phone"`
	Email        *string      `json:"email,omitempty" db:"email"`
	City         string       `json:"city" db:"city"`
	IsActive     bool         `json:"is_active" db:"is_active"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
}

// CustomerCreateRequest represents a request to create a customer
type CustomerCreateRequest struct {
	CustomerName  string       `json:"customer_name" binding:"required,max=255"`
	CustomerType  CustomerType `json:"customer_type" binding:"required"`
	Phone         string       `json:"phone" binding:"required,max=20"`
	Email         *string      `json:"email,omitempty" binding:"omitempty,email,max=100"`
	Address       string       `json:"address" binding:"required,max=500"`
	City          string       `json:"city" binding:"required,max=100"`
	PostalCode    *string      `json:"postal_code,omitempty" binding:"omitempty,max=10"`
	TaxNumber     *string      `json:"tax_number,omitempty" binding:"omitempty,max=30"`
	ContactPerson *string      `json:"contact_person,omitempty" binding:"omitempty,max=255"`
	Notes         *string      `json:"notes,omitempty"`
}

// CustomerUpdateRequest represents a request to update a customer
type CustomerUpdateRequest struct {
	CustomerName  *string       `json:"customer_name,omitempty" binding:"omitempty,max=255"`
	CustomerType  *CustomerType `json:"customer_type,omitempty"`
	Phone         *string       `json:"phone,omitempty" binding:"omitempty,max=20"`
	Email         *string       `json:"email,omitempty" binding:"omitempty,email,max=100"`
	Address       *string       `json:"address,omitempty" binding:"omitempty,max=500"`
	City          *string       `json:"city,omitempty" binding:"omitempty,max=100"`
	PostalCode    *string       `json:"postal_code,omitempty" binding:"omitempty,max=10"`
	TaxNumber     *string       `json:"tax_number,omitempty" binding:"omitempty,max=30"`
	ContactPerson *string       `json:"contact_person,omitempty" binding:"omitempty,max=255"`
	Notes         *string       `json:"notes,omitempty"`
	IsActive      *bool         `json:"is_active,omitempty"`
}

// CustomerFilterParams represents filtering parameters for customer queries
type CustomerFilterParams struct {
	CustomerType *CustomerType `json:"customer_type,omitempty" form:"customer_type"`
	IsActive     *bool         `json:"is_active,omitempty" form:"is_active"`
	Search       string        `json:"search,omitempty" form:"search"`
	City         string        `json:"city,omitempty" form:"city"`
	common.PaginationParams
}