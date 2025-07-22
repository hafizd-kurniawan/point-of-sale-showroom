package master

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/user"
)

// ProductCategory represents a product category in the system
type ProductCategory struct {
	CategoryID   int                   `json:"category_id" db:"category_id"`
	CategoryCode string                `json:"category_code" db:"category_code"`
	CategoryName string                `json:"category_name" db:"category_name"`
	Description  *string               `json:"description,omitempty" db:"description"`
	IsActive     bool                  `json:"is_active" db:"is_active"`
	CreatedAt    time.Time             `json:"created_at" db:"created_at"`
	CreatedBy    int                   `json:"created_by" db:"created_by"`
	Creator      *user.UserCreatorInfo `json:"creator,omitempty" db:"-"`
}

// ProductCategoryCreateRequest represents a request to create a product category
type ProductCategoryCreateRequest struct {
	CategoryName string  `json:"category_name" binding:"required,max=255"`
	Description  *string `json:"description,omitempty"`
}

// ProductCategoryUpdateRequest represents a request to update a product category
type ProductCategoryUpdateRequest struct {
	CategoryName *string `json:"category_name,omitempty" binding:"omitempty,max=255"`
	Description  *string `json:"description,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}