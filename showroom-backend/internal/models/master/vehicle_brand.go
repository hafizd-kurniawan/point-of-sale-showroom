package master

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/user"
)

// VehicleBrand represents a vehicle brand in the system
type VehicleBrand struct {
	BrandID       int                   `json:"brand_id" db:"brand_id"`
	BrandCode     string                `json:"brand_code" db:"brand_code"`
	BrandName     string                `json:"brand_name" db:"brand_name"`
	CountryOrigin *string               `json:"country_origin,omitempty" db:"country_origin"`
	LogoImage     *string               `json:"logo_image,omitempty" db:"logo_image"`
	IsActive      bool                  `json:"is_active" db:"is_active"`
	CreatedAt     time.Time             `json:"created_at" db:"created_at"`
	CreatedBy     int                   `json:"created_by" db:"created_by"`
	Creator       *user.UserCreatorInfo `json:"creator,omitempty" db:"-"`
}

// VehicleBrandCreateRequest represents a request to create a vehicle brand
type VehicleBrandCreateRequest struct {
	BrandName     string  `json:"brand_name" binding:"required,max=255"`
	CountryOrigin *string `json:"country_origin,omitempty" binding:"omitempty,max=100"`
	LogoImage     *string `json:"logo_image,omitempty" binding:"omitempty,max=500"`
}

// VehicleBrandUpdateRequest represents a request to update a vehicle brand
type VehicleBrandUpdateRequest struct {
	BrandName     *string `json:"brand_name,omitempty" binding:"omitempty,max=255"`
	CountryOrigin *string `json:"country_origin,omitempty" binding:"omitempty,max=100"`
	LogoImage     *string `json:"logo_image,omitempty" binding:"omitempty,max=500"`
	IsActive      *bool   `json:"is_active,omitempty"`
}