package master

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/user"
)

// Customer represents a customer in the system
type Customer struct {
	CustomerID     int                    `json:"customer_id" db:"customer_id"`
	CustomerCode   string                 `json:"customer_code" db:"customer_code"`
	CustomerName   string                 `json:"customer_name" db:"customer_name"`
	Phone          string                 `json:"phone" db:"phone"`
	Email          *string                `json:"email,omitempty" db:"email"`
	Address        *string                `json:"address,omitempty" db:"address"`
	City           *string                `json:"city,omitempty" db:"city"`
	Province       *string                `json:"province,omitempty" db:"province"`
	PostalCode     *string                `json:"postal_code,omitempty" db:"postal_code"`
	IDCardNumber   *string                `json:"id_card_number,omitempty" db:"id_card_number"`
	TaxNumber      *string                `json:"tax_number,omitempty" db:"tax_number"`
	CustomerType   common.CustomerType    `json:"customer_type" db:"customer_type"`
	BirthDate      *time.Time             `json:"birth_date,omitempty" db:"birth_date"`
	Occupation     *string                `json:"occupation,omitempty" db:"occupation"`
	IncomeRange    *float64               `json:"income_range,omitempty" db:"income_range"`
	CreatedAt      time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at" db:"updated_at"`
	CreatedBy      int                    `json:"created_by" db:"created_by"`
	IsActive       bool                   `json:"is_active" db:"is_active"`
	Notes          *string                `json:"notes,omitempty" db:"notes"`
	Creator        *user.UserCreatorInfo  `json:"creator,omitempty" db:"-"`
}

// CustomerCreateRequest represents a request to create a customer
type CustomerCreateRequest struct {
	CustomerName string                 `json:"customer_name" binding:"required,max=255"`
	Phone        string                 `json:"phone" binding:"required,max=20"`
	Email        *string                `json:"email,omitempty" binding:"omitempty,email,max=100"`
	Address      *string                `json:"address,omitempty"`
	City         *string                `json:"city,omitempty" binding:"omitempty,max=100"`
	Province     *string                `json:"province,omitempty" binding:"omitempty,max=100"`
	PostalCode   *string                `json:"postal_code,omitempty" binding:"omitempty,max=10"`
	IDCardNumber *string                `json:"id_card_number,omitempty" binding:"omitempty,max=20"`
	TaxNumber    *string                `json:"tax_number,omitempty" binding:"omitempty,max=50"`
	CustomerType common.CustomerType    `json:"customer_type" binding:"required"`
	BirthDate    *string                `json:"birth_date,omitempty"`
	Occupation   *string                `json:"occupation,omitempty" binding:"omitempty,max=100"`
	IncomeRange  *float64               `json:"income_range,omitempty"`
	Notes        *string                `json:"notes,omitempty"`
}

// CustomerUpdateRequest represents a request to update a customer
type CustomerUpdateRequest struct {
	CustomerName *string                `json:"customer_name,omitempty" binding:"omitempty,max=255"`
	Phone        *string                `json:"phone,omitempty" binding:"omitempty,max=20"`
	Email        *string                `json:"email,omitempty" binding:"omitempty,email,max=100"`
	Address      *string                `json:"address,omitempty"`
	City         *string                `json:"city,omitempty" binding:"omitempty,max=100"`
	Province     *string                `json:"province,omitempty" binding:"omitempty,max=100"`
	PostalCode   *string                `json:"postal_code,omitempty" binding:"omitempty,max=10"`
	CustomerType *common.CustomerType   `json:"customer_type,omitempty"`
	Occupation   *string                `json:"occupation,omitempty" binding:"omitempty,max=100"`
	IncomeRange  *float64               `json:"income_range,omitempty"`
	IsActive     *bool                  `json:"is_active,omitempty"`
	Notes        *string                `json:"notes,omitempty"`
}

// CustomerFilterParams represents filtering parameters for customer queries
type CustomerFilterParams struct {
	CustomerType *common.CustomerType `json:"customer_type,omitempty" form:"customer_type"`
	IsActive     *bool                `json:"is_active,omitempty" form:"is_active"`
	Search       string               `json:"search,omitempty" form:"search"`
	City         string               `json:"city,omitempty" form:"city"`
	Province     string               `json:"province,omitempty" form:"province"`
	common.PaginationParams
}

// CustomerListItem represents a simplified customer for list views
type CustomerListItem struct {
	CustomerID   int                 `json:"customer_id" db:"customer_id"`
	CustomerCode string              `json:"customer_code" db:"customer_code"`
	CustomerName string              `json:"customer_name" db:"customer_name"`
	Phone        string              `json:"phone" db:"phone"`
	Email        *string             `json:"email,omitempty" db:"email"`
	CustomerType common.CustomerType `json:"customer_type" db:"customer_type"`
	City         *string             `json:"city,omitempty" db:"city"`
	IsActive     bool                `json:"is_active" db:"is_active"`
	CreatedAt    time.Time           `json:"created_at" db:"created_at"`
}