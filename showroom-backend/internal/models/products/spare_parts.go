package products

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// ProductSparePart represents a spare part product in the system
type ProductSparePart struct {
	ProductID        int         `json:"product_id" db:"product_id"`
	ProductCode      string      `json:"product_code" db:"product_code"`
	ProductName      string      `json:"product_name" db:"product_name"`
	Description      *string     `json:"description,omitempty" db:"description"`
	BrandID          int         `json:"brand_id" db:"brand_id"`
	CategoryID       int         `json:"category_id" db:"category_id"`
	UnitMeasure      string      `json:"unit_measure" db:"unit_measure"`
	CostPrice        float64     `json:"cost_price" db:"cost_price"`
	SellingPrice     float64     `json:"selling_price" db:"selling_price"`
	MarkupPercentage float64     `json:"markup_percentage" db:"markup_percentage"`
	StockQuantity    int         `json:"stock_quantity" db:"stock_quantity"`
	MinStockLevel    int         `json:"min_stock_level" db:"min_stock_level"`
	MaxStockLevel    int         `json:"max_stock_level" db:"max_stock_level"`
	LocationRack     *string     `json:"location_rack,omitempty" db:"location_rack"`
	Barcode          *string     `json:"barcode,omitempty" db:"barcode"`
	Weight           *float64    `json:"weight,omitempty" db:"weight"`
	Dimensions       *string     `json:"dimensions,omitempty" db:"dimensions"`
	CreatedAt        time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at" db:"updated_at"`
	CreatedBy        int         `json:"created_by" db:"created_by"`
	IsActive         bool        `json:"is_active" db:"is_active"`
	ProductImage     *string     `json:"product_image,omitempty" db:"product_image"`
	Notes            *string     `json:"notes,omitempty" db:"notes"`
}

// ProductSparePartListItem represents a simplified spare part for list views
type ProductSparePartListItem struct {
	ProductID        int     `json:"product_id" db:"product_id"`
	ProductCode      string  `json:"product_code" db:"product_code"`
	ProductName      string  `json:"product_name" db:"product_name"`
	BrandID          int     `json:"brand_id" db:"brand_id"`
	CategoryID       int     `json:"category_id" db:"category_id"`
	UnitMeasure      string  `json:"unit_measure" db:"unit_measure"`
	CostPrice        float64 `json:"cost_price" db:"cost_price"`
	SellingPrice     float64 `json:"selling_price" db:"selling_price"`
	StockQuantity    int     `json:"stock_quantity" db:"stock_quantity"`
	MinStockLevel    int     `json:"min_stock_level" db:"min_stock_level"`
	LocationRack     *string `json:"location_rack,omitempty" db:"location_rack"`
	IsActive         bool    `json:"is_active" db:"is_active"`
}

// ProductSparePartCreateRequest represents a request to create a spare part
type ProductSparePartCreateRequest struct {
	ProductName      string   `json:"product_name" binding:"required,max=255"`
	Description      *string  `json:"description,omitempty"`
	BrandID          int      `json:"brand_id" binding:"required,min=1"`
	CategoryID       int      `json:"category_id" binding:"required,min=1"`
	UnitMeasure      string   `json:"unit_measure" binding:"required,max=50"`
	CostPrice        float64  `json:"cost_price" binding:"required,min=0"`
	SellingPrice     float64  `json:"selling_price" binding:"required,min=0"`
	MarkupPercentage float64  `json:"markup_percentage" binding:"min=0"`
	MinStockLevel    int      `json:"min_stock_level" binding:"min=0"`
	MaxStockLevel    int      `json:"max_stock_level" binding:"min=0"`
	LocationRack     *string  `json:"location_rack,omitempty" binding:"omitempty,max=100"`
	Barcode          *string  `json:"barcode,omitempty" binding:"omitempty,max=100"`
	Weight           *float64 `json:"weight,omitempty" binding:"omitempty,min=0"`
	Dimensions       *string  `json:"dimensions,omitempty" binding:"omitempty,max=100"`
	ProductImage     *string  `json:"product_image,omitempty" binding:"omitempty,max=500"`
	Notes            *string  `json:"notes,omitempty"`
}

// ProductSparePartUpdateRequest represents a request to update a spare part
type ProductSparePartUpdateRequest struct {
	ProductName      *string  `json:"product_name,omitempty" binding:"omitempty,max=255"`
	Description      *string  `json:"description,omitempty"`
	BrandID          *int     `json:"brand_id,omitempty" binding:"omitempty,min=1"`
	CategoryID       *int     `json:"category_id,omitempty" binding:"omitempty,min=1"`
	UnitMeasure      *string  `json:"unit_measure,omitempty" binding:"omitempty,max=50"`
	CostPrice        *float64 `json:"cost_price,omitempty" binding:"omitempty,min=0"`
	SellingPrice     *float64 `json:"selling_price,omitempty" binding:"omitempty,min=0"`
	MarkupPercentage *float64 `json:"markup_percentage,omitempty" binding:"omitempty,min=0"`
	MinStockLevel    *int     `json:"min_stock_level,omitempty" binding:"omitempty,min=0"`
	MaxStockLevel    *int     `json:"max_stock_level,omitempty" binding:"omitempty,min=0"`
	LocationRack     *string  `json:"location_rack,omitempty" binding:"omitempty,max=100"`
	Barcode          *string  `json:"barcode,omitempty" binding:"omitempty,max=100"`
	Weight           *float64 `json:"weight,omitempty" binding:"omitempty,min=0"`
	Dimensions       *string  `json:"dimensions,omitempty" binding:"omitempty,max=100"`
	ProductImage     *string  `json:"product_image,omitempty" binding:"omitempty,max=500"`
	Notes            *string  `json:"notes,omitempty"`
	IsActive         *bool    `json:"is_active,omitempty"`
}

// ProductSparePartFilterParams represents filtering parameters for spare part queries
type ProductSparePartFilterParams struct {
	BrandID         *int    `json:"brand_id,omitempty" form:"brand_id"`
	CategoryID      *int    `json:"category_id,omitempty" form:"category_id"`
	IsActive        *bool   `json:"is_active,omitempty" form:"is_active"`
	LowStock        *bool   `json:"low_stock,omitempty" form:"low_stock"`
	Search          string  `json:"search,omitempty" form:"search"`
	MinPrice        *float64 `json:"min_price,omitempty" form:"min_price"`
	MaxPrice        *float64 `json:"max_price,omitempty" form:"max_price"`
	common.PaginationParams
}

// IsLowStock checks if the product is below minimum stock level
func (p *ProductSparePart) IsLowStock() bool {
	return p.StockQuantity <= p.MinStockLevel
}

// CalculateMarkup calculates markup percentage based on cost and selling price
func (p *ProductSparePart) CalculateMarkup() float64 {
	if p.CostPrice == 0 {
		return 0
	}
	return ((p.SellingPrice - p.CostPrice) / p.CostPrice) * 100
}

// UpdateMarkupPercentage updates the markup percentage based on current prices
func (p *ProductSparePart) UpdateMarkupPercentage() {
	p.MarkupPercentage = p.CalculateMarkup()
}

// CalculateSellingPriceFromMarkup calculates selling price based on cost price and markup
func (p *ProductSparePart) CalculateSellingPriceFromMarkup() float64 {
	return p.CostPrice * (1 + p.MarkupPercentage/100)
}