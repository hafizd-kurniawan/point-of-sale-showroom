package master

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/user"
)

// VehicleModel represents a vehicle model in the system
type VehicleModel struct {
	ModelID            int                      `json:"model_id" db:"model_id"`
	BrandID            int                      `json:"brand_id" db:"brand_id"`
	CategoryID         int                      `json:"category_id" db:"category_id"`
	ModelCode          string                   `json:"model_code" db:"model_code"`
	ModelName          string                   `json:"model_name" db:"model_name"`
	YearStart          int                      `json:"year_start" db:"year_start"`
	YearEnd            *int                     `json:"year_end,omitempty" db:"year_end"`
	FuelType           common.FuelType          `json:"fuel_type" db:"fuel_type"`
	Transmission       common.TransmissionType  `json:"transmission" db:"transmission"`
	EngineCapacity     int                      `json:"engine_capacity" db:"engine_capacity"`
	SpecificationsJSON *string                  `json:"specifications_json,omitempty" db:"specifications_json"`
	IsActive           bool                     `json:"is_active" db:"is_active"`
	CreatedAt          time.Time                `json:"created_at" db:"created_at"`
	CreatedBy          int                      `json:"created_by" db:"created_by"`
	Creator            *user.UserCreatorInfo    `json:"creator,omitempty" db:"-"`
	Brand              *VehicleBrandInfo        `json:"brand,omitempty" db:"-"`
	Category           *VehicleCategoryInfo     `json:"category,omitempty" db:"-"`
}

// VehicleBrandInfo represents minimal brand information for vehicle models
type VehicleBrandInfo struct {
	BrandID       int     `json:"brand_id" db:"brand_id"`
	BrandName     string  `json:"brand_name" db:"brand_name"`
	CountryOrigin *string `json:"country_origin,omitempty" db:"country_origin"`
}

// VehicleCategoryInfo represents minimal category information for vehicle models
type VehicleCategoryInfo struct {
	CategoryID   int    `json:"category_id" db:"category_id"`
	CategoryName string `json:"category_name" db:"category_name"`
}

// VehicleModelCreateRequest represents a request to create a vehicle model
type VehicleModelCreateRequest struct {
	BrandID            int                      `json:"brand_id" binding:"required"`
	CategoryID         int                      `json:"category_id" binding:"required"`
	ModelName          string                   `json:"model_name" binding:"required,max=255"`
	YearStart          int                      `json:"year_start" binding:"required,min=1900,max=2100"`
	YearEnd            *int                     `json:"year_end,omitempty" binding:"omitempty,min=1900,max=2100"`
	FuelType           common.FuelType          `json:"fuel_type" binding:"required"`
	Transmission       common.TransmissionType  `json:"transmission" binding:"required"`
	EngineCapacity     int                      `json:"engine_capacity" binding:"required,min=100"`
	SpecificationsJSON *string                  `json:"specifications_json,omitempty"`
}

// VehicleModelUpdateRequest represents a request to update a vehicle model
type VehicleModelUpdateRequest struct {
	BrandID            *int                     `json:"brand_id,omitempty"`
	CategoryID         *int                     `json:"category_id,omitempty"`
	ModelName          *string                  `json:"model_name,omitempty" binding:"omitempty,max=255"`
	YearStart          *int                     `json:"year_start,omitempty" binding:"omitempty,min=1900,max=2100"`
	YearEnd            *int                     `json:"year_end,omitempty" binding:"omitempty,min=1900,max=2100"`
	FuelType           *common.FuelType         `json:"fuel_type,omitempty"`
	Transmission       *common.TransmissionType `json:"transmission,omitempty"`
	EngineCapacity     *int                     `json:"engine_capacity,omitempty" binding:"omitempty,min=100"`
	SpecificationsJSON *string                  `json:"specifications_json,omitempty"`
	IsActive           *bool                    `json:"is_active,omitempty"`
}

// VehicleModelFilterParams represents filtering parameters for vehicle model queries
type VehicleModelFilterParams struct {
	BrandID      *int                     `json:"brand_id,omitempty" form:"brand_id"`
	CategoryID   *int                     `json:"category_id,omitempty" form:"category_id"`
	FuelType     *common.FuelType         `json:"fuel_type,omitempty" form:"fuel_type"`
	Transmission *common.TransmissionType `json:"transmission,omitempty" form:"transmission"`
	IsActive     *bool                    `json:"is_active,omitempty" form:"is_active"`
	YearStart    *int                     `json:"year_start,omitempty" form:"year_start"`
	Search       string                   `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// VehicleModelListItem represents a simplified vehicle model for list views
type VehicleModelListItem struct {
	ModelID        int                      `json:"model_id" db:"model_id"`
	ModelCode      string                   `json:"model_code" db:"model_code"`
	ModelName      string                   `json:"model_name" db:"model_name"`
	BrandID        int                      `json:"brand_id" db:"brand_id"`
	BrandName      string                   `json:"brand_name" db:"brand_name"`
	CategoryID     int                      `json:"category_id" db:"category_id"`
	CategoryName   string                   `json:"category_name" db:"category_name"`
	YearStart      int                      `json:"year_start" db:"year_start"`
	YearEnd        *int                     `json:"year_end,omitempty" db:"year_end"`
	FuelType       common.FuelType          `json:"fuel_type" db:"fuel_type"`
	Transmission   common.TransmissionType  `json:"transmission" db:"transmission"`
	EngineCapacity int                      `json:"engine_capacity" db:"engine_capacity"`
	IsActive       bool                     `json:"is_active" db:"is_active"`
	CreatedAt      time.Time                `json:"created_at" db:"created_at"`
}