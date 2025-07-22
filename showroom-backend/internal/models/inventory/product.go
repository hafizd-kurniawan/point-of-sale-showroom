package inventory

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// ProductSparePart represents a spare part product in the system
type ProductSparePart struct {
	ProductID        int       `json:"product_id" db:"product_id"`
	ProductCode      string    `json:"product_code" db:"product_code"`
	ProductName      string    `json:"product_name" db:"product_name"`
	Description      *string   `json:"description,omitempty" db:"description"`
	BrandID          *int      `json:"brand_id,omitempty" db:"brand_id"`
	CategoryID       *int      `json:"category_id,omitempty" db:"category_id"`
	UnitMeasure      string    `json:"unit_measure" db:"unit_measure"`
	CostPrice        float64   `json:"cost_price" db:"cost_price"`
	SellingPrice     float64   `json:"selling_price" db:"selling_price"`
	MarkupPercentage *float64  `json:"markup_percentage,omitempty" db:"markup_percentage"`
	StockQuantity    int       `json:"stock_quantity" db:"stock_quantity"`
	MinStockLevel    int       `json:"min_stock_level" db:"min_stock_level"`
	MaxStockLevel    int       `json:"max_stock_level" db:"max_stock_level"`
	LocationRack     *string   `json:"location_rack,omitempty" db:"location_rack"`
	Barcode          *string   `json:"barcode,omitempty" db:"barcode"`
	Weight           *float64  `json:"weight,omitempty" db:"weight"`
	Dimensions       *string   `json:"dimensions,omitempty" db:"dimensions"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy        int       `json:"created_by" db:"created_by"`
	IsActive         bool      `json:"is_active" db:"is_active"`
	ProductImage     *string   `json:"product_image,omitempty" db:"product_image"`
	Notes            *string   `json:"notes,omitempty" db:"notes"`

	// Related data
	BrandName    *string `json:"brand_name,omitempty" db:"brand_name"`
	CategoryName *string `json:"category_name,omitempty" db:"category_name"`
}

// ProductSparePartListItem represents a simplified product for list views
type ProductSparePartListItem struct {
	ProductID     int      `json:"product_id" db:"product_id"`
	ProductCode   string   `json:"product_code" db:"product_code"`
	ProductName   string   `json:"product_name" db:"product_name"`
	BrandName     *string  `json:"brand_name,omitempty" db:"brand_name"`
	CategoryName  *string  `json:"category_name,omitempty" db:"category_name"`
	UnitMeasure   string   `json:"unit_measure" db:"unit_measure"`
	CostPrice     float64  `json:"cost_price" db:"cost_price"`
	SellingPrice  float64  `json:"selling_price" db:"selling_price"`
	StockQuantity int      `json:"stock_quantity" db:"stock_quantity"`
	MinStockLevel int      `json:"min_stock_level" db:"min_stock_level"`
	LocationRack  *string  `json:"location_rack,omitempty" db:"location_rack"`
	IsActive      bool     `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// ProductSparePartCreateRequest represents a request to create a product
type ProductSparePartCreateRequest struct {
	ProductName      string   `json:"product_name" binding:"required,max=255"`
	Description      *string  `json:"description,omitempty"`
	BrandID          *int     `json:"brand_id,omitempty"`
	CategoryID       *int     `json:"category_id,omitempty"`
	UnitMeasure      string   `json:"unit_measure" binding:"required,max=50"`
	CostPrice        float64  `json:"cost_price" binding:"required,min=0"`
	SellingPrice     float64  `json:"selling_price" binding:"required,min=0"`
	MarkupPercentage *float64 `json:"markup_percentage,omitempty" binding:"omitempty,min=0"`
	StockQuantity    int      `json:"stock_quantity" binding:"min=0"`
	MinStockLevel    int      `json:"min_stock_level" binding:"min=0"`
	MaxStockLevel    int      `json:"max_stock_level" binding:"min=0"`
	LocationRack     *string  `json:"location_rack,omitempty" binding:"omitempty,max=100"`
	Barcode          *string  `json:"barcode,omitempty" binding:"omitempty,max=100"`
	Weight           *float64 `json:"weight,omitempty" binding:"omitempty,min=0"`
	Dimensions       *string  `json:"dimensions,omitempty" binding:"omitempty,max=100"`
	ProductImage     *string  `json:"product_image,omitempty" binding:"omitempty,max=500"`
	Notes            *string  `json:"notes,omitempty"`
}

// ProductSparePartUpdateRequest represents a request to update a product
type ProductSparePartUpdateRequest struct {
	ProductName      *string  `json:"product_name,omitempty" binding:"omitempty,max=255"`
	Description      *string  `json:"description,omitempty"`
	BrandID          *int     `json:"brand_id,omitempty"`
	CategoryID       *int     `json:"category_id,omitempty"`
	UnitMeasure      *string  `json:"unit_measure,omitempty" binding:"omitempty,max=50"`
	CostPrice        *float64 `json:"cost_price,omitempty" binding:"omitempty,min=0"`
	SellingPrice     *float64 `json:"selling_price,omitempty" binding:"omitempty,min=0"`
	MarkupPercentage *float64 `json:"markup_percentage,omitempty" binding:"omitempty,min=0"`
	StockQuantity    *int     `json:"stock_quantity,omitempty" binding:"omitempty,min=0"`
	MinStockLevel    *int     `json:"min_stock_level,omitempty" binding:"omitempty,min=0"`
	MaxStockLevel    *int     `json:"max_stock_level,omitempty" binding:"omitempty,min=0"`
	LocationRack     *string  `json:"location_rack,omitempty" binding:"omitempty,max=100"`
	Barcode          *string  `json:"barcode,omitempty" binding:"omitempty,max=100"`
	Weight           *float64 `json:"weight,omitempty" binding:"omitempty,min=0"`
	Dimensions       *string  `json:"dimensions,omitempty" binding:"omitempty,max=100"`
	IsActive         *bool    `json:"is_active,omitempty"`
	ProductImage     *string  `json:"product_image,omitempty" binding:"omitempty,max=500"`
	Notes            *string  `json:"notes,omitempty"`
}

// ProductSparePartFilterParams represents filtering parameters for product queries
type ProductSparePartFilterParams struct {
	BrandID       *int    `json:"brand_id,omitempty" form:"brand_id"`
	CategoryID    *int    `json:"category_id,omitempty" form:"category_id"`
	MinPrice      *float64 `json:"min_price,omitempty" form:"min_price"`
	MaxPrice      *float64 `json:"max_price,omitempty" form:"max_price"`
	LocationRack  string  `json:"location_rack,omitempty" form:"location_rack"`
	Barcode       string  `json:"barcode,omitempty" form:"barcode"`
	LowStock      *bool   `json:"low_stock,omitempty" form:"low_stock"`
	IsActive      *bool   `json:"is_active,omitempty" form:"is_active"`
	Search        string  `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// IsLowStock checks if the product is below minimum stock level
func (p *ProductSparePart) IsLowStock() bool {
	return p.StockQuantity <= p.MinStockLevel
}

// IsOverStock checks if the product is above maximum stock level
func (p *ProductSparePart) IsOverStock() bool {
	return p.StockQuantity >= p.MaxStockLevel
}

// CalculateMarkup calculates markup percentage from cost and selling price
func (p *ProductSparePart) CalculateMarkup() float64 {
	if p.CostPrice == 0 {
		return 0
	}
	return ((p.SellingPrice - p.CostPrice) / p.CostPrice) * 100
}

// UpdateMarkupFromPrices updates markup percentage based on current prices
func (p *ProductSparePart) UpdateMarkupFromPrices() {
	markup := p.CalculateMarkup()
	p.MarkupPercentage = &markup
}

// CalculateSellingPriceFromMarkup calculates selling price from cost price and markup percentage
func (p *ProductSparePart) CalculateSellingPriceFromMarkup() float64 {
	if p.MarkupPercentage == nil || p.CostPrice == 0 {
		return p.SellingPrice
	}
	return p.CostPrice * (1 + (*p.MarkupPercentage / 100))
}