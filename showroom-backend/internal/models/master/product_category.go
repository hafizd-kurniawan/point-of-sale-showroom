package master

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// ProductCategory represents a product category in the system
type ProductCategory struct {
	CategoryID   int       `json:"category_id" db:"category_id"`
	CategoryCode string    `json:"category_code" db:"category_code"`
	CategoryName string    `json:"category_name" db:"category_name"`
	Description  *string   `json:"description,omitempty" db:"description"`
	ParentID     *int      `json:"parent_id,omitempty" db:"parent_id"`
	Level        int       `json:"level" db:"level"`
	Path         string    `json:"path" db:"path"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy    int       `json:"created_by" db:"created_by"`
	
	// Related data
	ParentName *string `json:"parent_name,omitempty" db:"parent_name"`
	Children   []ProductCategory `json:"children,omitempty" db:"-"`
}

// ProductCategoryListItem represents a simplified product category for list views
type ProductCategoryListItem struct {
	CategoryID   int       `json:"category_id" db:"category_id"`
	CategoryCode string    `json:"category_code" db:"category_code"`
	CategoryName string    `json:"category_name" db:"category_name"`
	ParentName   *string   `json:"parent_name,omitempty" db:"parent_name"`
	Level        int       `json:"level" db:"level"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// ProductCategoryCreateRequest represents a request to create a product category
type ProductCategoryCreateRequest struct {
	CategoryName string  `json:"category_name" binding:"required,max=100"`
	Description  *string `json:"description,omitempty"`
	ParentID     *int    `json:"parent_id,omitempty"`
}

// ProductCategoryUpdateRequest represents a request to update a product category
type ProductCategoryUpdateRequest struct {
	CategoryName *string `json:"category_name,omitempty" binding:"omitempty,max=100"`
	Description  *string `json:"description,omitempty"`
	ParentID     *int    `json:"parent_id,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

// ProductCategoryFilterParams represents filtering parameters for product category queries
type ProductCategoryFilterParams struct {
	ParentID *int   `json:"parent_id,omitempty" form:"parent_id"`
	Level    *int   `json:"level,omitempty" form:"level"`
	IsActive *bool  `json:"is_active,omitempty" form:"is_active"`
	Search   string `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// ProductCategoryTree represents a hierarchical view of product categories
type ProductCategoryTree struct {
	CategoryID   int                    `json:"category_id"`
	CategoryCode string                 `json:"category_code"`
	CategoryName string                 `json:"category_name"`
	Description  *string                `json:"description,omitempty"`
	Level        int                    `json:"level"`
	IsActive     bool                   `json:"is_active"`
	Children     []ProductCategoryTree  `json:"children,omitempty"`
}