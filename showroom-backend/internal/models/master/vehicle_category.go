package master

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/user"
)

// VehicleCategory represents a vehicle category in the system
type VehicleCategory struct {
	CategoryID   int                   `json:"category_id" db:"category_id"`
	CategoryCode string                `json:"category_code" db:"category_code"`
	CategoryName string                `json:"category_name" db:"category_name"`
	Description  *string               `json:"description,omitempty" db:"description"`
	IsActive     bool                  `json:"is_active" db:"is_active"`
	CreatedAt    time.Time             `json:"created_at" db:"created_at"`
	CreatedBy    int                   `json:"created_by" db:"created_by"`
	Creator      *user.UserCreatorInfo `json:"creator,omitempty" db:"-"`
}

// VehicleCategoryCreateRequest represents a request to create a vehicle category
type VehicleCategoryCreateRequest struct {
	CategoryName string  `json:"category_name" binding:"required,max=255"`
	Description  *string `json:"description,omitempty"`
}

// VehicleCategoryUpdateRequest represents a request to update a vehicle category
type VehicleCategoryUpdateRequest struct {
	CategoryName *string `json:"category_name,omitempty" binding:"omitempty,max=255"`
	Description  *string `json:"description,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}