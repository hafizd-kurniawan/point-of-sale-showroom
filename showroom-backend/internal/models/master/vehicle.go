package master

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// VehicleBrand represents a vehicle brand in the system
type VehicleBrand struct {
	BrandID     int       `json:"brand_id" db:"brand_id"`
	BrandCode   string    `json:"brand_code" db:"brand_code"`
	BrandName   string    `json:"brand_name" db:"brand_name"`
	CountryOrigin string  `json:"country_origin" db:"country_origin"`
	Description *string   `json:"description,omitempty" db:"description"`
	LogoURL     *string   `json:"logo_url,omitempty" db:"logo_url"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy   int       `json:"created_by" db:"created_by"`
}

// VehicleBrandListItem represents a simplified vehicle brand for list views
type VehicleBrandListItem struct {
	BrandID       int       `json:"brand_id" db:"brand_id"`
	BrandCode     string    `json:"brand_code" db:"brand_code"`
	BrandName     string    `json:"brand_name" db:"brand_name"`
	CountryOrigin string    `json:"country_origin" db:"country_origin"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// VehicleBrandCreateRequest represents a request to create a vehicle brand
type VehicleBrandCreateRequest struct {
	BrandName     string  `json:"brand_name" binding:"required,max=100"`
	CountryOrigin string  `json:"country_origin" binding:"required,max=100"`
	Description   *string `json:"description,omitempty"`
	LogoURL       *string `json:"logo_url,omitempty" binding:"omitempty,max=500"`
}

// VehicleBrandUpdateRequest represents a request to update a vehicle brand
type VehicleBrandUpdateRequest struct {
	BrandName     *string `json:"brand_name,omitempty" binding:"omitempty,max=100"`
	CountryOrigin *string `json:"country_origin,omitempty" binding:"omitempty,max=100"`
	Description   *string `json:"description,omitempty"`
	LogoURL       *string `json:"logo_url,omitempty" binding:"omitempty,max=500"`
	IsActive      *bool   `json:"is_active,omitempty"`
}

// VehicleBrandFilterParams represents filtering parameters for vehicle brand queries
type VehicleBrandFilterParams struct {
	IsActive      *bool  `json:"is_active,omitempty" form:"is_active"`
	Search        string `json:"search,omitempty" form:"search"`
	CountryOrigin string `json:"country_origin,omitempty" form:"country_origin"`
	common.PaginationParams
}

// VehicleCategory represents a vehicle category in the system
type VehicleCategory struct {
	CategoryID   int       `json:"category_id" db:"category_id"`
	CategoryCode string    `json:"category_code" db:"category_code"`
	CategoryName string    `json:"category_name" db:"category_name"`
	Description  *string   `json:"description,omitempty" db:"description"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy    int       `json:"created_by" db:"created_by"`
}

// VehicleCategoryListItem represents a simplified vehicle category for list views
type VehicleCategoryListItem struct {
	CategoryID   int       `json:"category_id" db:"category_id"`
	CategoryCode string    `json:"category_code" db:"category_code"`
	CategoryName string    `json:"category_name" db:"category_name"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// VehicleCategoryCreateRequest represents a request to create a vehicle category
type VehicleCategoryCreateRequest struct {
	CategoryName string  `json:"category_name" binding:"required,max=100"`
	Description  *string `json:"description,omitempty"`
}

// VehicleCategoryUpdateRequest represents a request to update a vehicle category
type VehicleCategoryUpdateRequest struct {
	CategoryName *string `json:"category_name,omitempty" binding:"omitempty,max=100"`
	Description  *string `json:"description,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

// VehicleCategoryFilterParams represents filtering parameters for vehicle category queries
type VehicleCategoryFilterParams struct {
	IsActive *bool  `json:"is_active,omitempty" form:"is_active"`
	Search   string `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// VehicleModel represents a vehicle model in the system
type VehicleModel struct {
	ModelID        int       `json:"model_id" db:"model_id"`
	ModelCode      string    `json:"model_code" db:"model_code"`
	ModelName      string    `json:"model_name" db:"model_name"`
	BrandID        int       `json:"brand_id" db:"brand_id"`
	CategoryID     int       `json:"category_id" db:"category_id"`
	ModelYear      int       `json:"model_year" db:"model_year"`
	EngineCapacity *float64  `json:"engine_capacity,omitempty" db:"engine_capacity"`
	FuelType       string    `json:"fuel_type" db:"fuel_type"`
	Transmission   string    `json:"transmission" db:"transmission"`
	SeatCapacity   int       `json:"seat_capacity" db:"seat_capacity"`
	Color          string    `json:"color" db:"color"`
	Price          float64   `json:"price" db:"price"`
	Description    *string   `json:"description,omitempty" db:"description"`
	ImageURL       *string   `json:"image_url,omitempty" db:"image_url"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy      int       `json:"created_by" db:"created_by"`
	
	// Related data
	BrandName    string `json:"brand_name,omitempty" db:"brand_name"`
	CategoryName string `json:"category_name,omitempty" db:"category_name"`
}

// VehicleModelListItem represents a simplified vehicle model for list views
type VehicleModelListItem struct {
	ModelID        int      `json:"model_id" db:"model_id"`
	ModelCode      string   `json:"model_code" db:"model_code"`
	ModelName      string   `json:"model_name" db:"model_name"`
	BrandName      string   `json:"brand_name" db:"brand_name"`
	CategoryName   string   `json:"category_name" db:"category_name"`
	ModelYear      int      `json:"model_year" db:"model_year"`
	FuelType       string   `json:"fuel_type" db:"fuel_type"`
	Transmission   string   `json:"transmission" db:"transmission"`
	SeatCapacity   int      `json:"seat_capacity" db:"seat_capacity"`
	Price          float64  `json:"price" db:"price"`
	IsActive       bool     `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// VehicleModelCreateRequest represents a request to create a vehicle model
type VehicleModelCreateRequest struct {
	ModelName      string   `json:"model_name" binding:"required,max=100"`
	BrandID        int      `json:"brand_id" binding:"required"`
	CategoryID     int      `json:"category_id" binding:"required"`
	ModelYear      int      `json:"model_year" binding:"required,min=1900,max=2100"`
	EngineCapacity *float64 `json:"engine_capacity,omitempty" binding:"omitempty,min=0"`
	FuelType       string   `json:"fuel_type" binding:"required,max=50"`
	Transmission   string   `json:"transmission" binding:"required,max=50"`
	SeatCapacity   int      `json:"seat_capacity" binding:"required,min=1,max=50"`
	Color          string   `json:"color" binding:"required,max=50"`
	Price          float64  `json:"price" binding:"required,min=0"`
	Description    *string  `json:"description,omitempty"`
	ImageURL       *string  `json:"image_url,omitempty" binding:"omitempty,max=500"`
}

// VehicleModelUpdateRequest represents a request to update a vehicle model
type VehicleModelUpdateRequest struct {
	ModelName      *string  `json:"model_name,omitempty" binding:"omitempty,max=100"`
	BrandID        *int     `json:"brand_id,omitempty"`
	CategoryID     *int     `json:"category_id,omitempty"`
	ModelYear      *int     `json:"model_year,omitempty" binding:"omitempty,min=1900,max=2100"`
	EngineCapacity *float64 `json:"engine_capacity,omitempty" binding:"omitempty,min=0"`
	FuelType       *string  `json:"fuel_type,omitempty" binding:"omitempty,max=50"`
	Transmission   *string  `json:"transmission,omitempty" binding:"omitempty,max=50"`
	SeatCapacity   *int     `json:"seat_capacity,omitempty" binding:"omitempty,min=1,max=50"`
	Color          *string  `json:"color,omitempty" binding:"omitempty,max=50"`
	Price          *float64 `json:"price,omitempty" binding:"omitempty,min=0"`
	Description    *string  `json:"description,omitempty"`
	ImageURL       *string  `json:"image_url,omitempty" binding:"omitempty,max=500"`
	IsActive       *bool    `json:"is_active,omitempty"`
}

// VehicleModelFilterParams represents filtering parameters for vehicle model queries
type VehicleModelFilterParams struct {
	BrandID      *int    `json:"brand_id,omitempty" form:"brand_id"`
	CategoryID   *int    `json:"category_id,omitempty" form:"category_id"`
	ModelYear    *int    `json:"model_year,omitempty" form:"model_year"`
	FuelType     string  `json:"fuel_type,omitempty" form:"fuel_type"`
	Transmission string  `json:"transmission,omitempty" form:"transmission"`
	MinPrice     *float64 `json:"min_price,omitempty" form:"min_price"`
	MaxPrice     *float64 `json:"max_price,omitempty" form:"max_price"`
	IsActive     *bool   `json:"is_active,omitempty" form:"is_active"`
	Search       string  `json:"search,omitempty" form:"search"`
	common.PaginationParams
}