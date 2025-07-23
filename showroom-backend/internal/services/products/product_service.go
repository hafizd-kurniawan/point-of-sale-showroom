package products

import (
	"context"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// ProductService handles business logic for products
type ProductService struct {
	productRepo interfaces.ProductSparePartRepository
	stockRepo   interfaces.StockMovementRepository
}

// NewProductService creates a new product service
func NewProductService(
	productRepo interfaces.ProductSparePartRepository,
	stockRepo interfaces.StockMovementRepository,
) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		stockRepo:   stockRepo,
	}
}

// CreateProduct creates a new product with auto-generated code
func (s *ProductService) CreateProduct(ctx context.Context, req *products.ProductSparePartCreateRequest, createdBy int) (*products.ProductSparePart, error) {
	// Generate product code
	productCode, err := s.productRepo.GenerateCode(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate product code: %w", err)
	}

	// Calculate markup percentage if not provided
	markupPercentage := req.MarkupPercentage
	if req.CostPrice > 0 && req.SellingPrice > 0 {
		markupPercentage = ((req.SellingPrice - req.CostPrice) / req.CostPrice) * 100
	}

	// Create product model
	product := &products.ProductSparePart{
		ProductCode:      productCode,
		ProductName:      req.ProductName,
		Description:      req.Description,
		BrandID:          req.BrandID,
		CategoryID:       req.CategoryID,
		UnitMeasure:      req.UnitMeasure,
		CostPrice:        req.CostPrice,
		SellingPrice:     req.SellingPrice,
		MarkupPercentage: markupPercentage,
		StockQuantity:    0, // Start with zero stock
		MinStockLevel:    req.MinStockLevel,
		MaxStockLevel:    req.MaxStockLevel,
		LocationRack:     req.LocationRack,
		Barcode:          req.Barcode,
		Weight:           req.Weight,
		Dimensions:       req.Dimensions,
		CreatedBy:        createdBy,
		IsActive:         true,
		ProductImage:     req.ProductImage,
		Notes:            req.Notes,
	}

	// Validate barcode uniqueness if provided
	if req.Barcode != nil && *req.Barcode != "" {
		exists, err := s.productRepo.IsBarcodeExists(ctx, *req.Barcode, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to check barcode uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("barcode %s already exists", *req.Barcode)
		}
	}

	// Create product
	createdProduct, err := s.productRepo.Create(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return createdProduct, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, id int, req *products.ProductSparePartUpdateRequest) (*products.ProductSparePart, error) {
	// Get existing product
	existing, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing product: %w", err)
	}

	// Update fields if provided
	if req.ProductName != nil {
		existing.ProductName = *req.ProductName
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.BrandID != nil {
		existing.BrandID = *req.BrandID
	}
	if req.CategoryID != nil {
		existing.CategoryID = *req.CategoryID
	}
	if req.UnitMeasure != nil {
		existing.UnitMeasure = *req.UnitMeasure
	}
	if req.CostPrice != nil {
		existing.CostPrice = *req.CostPrice
	}
	if req.SellingPrice != nil {
		existing.SellingPrice = *req.SellingPrice
	}
	if req.MarkupPercentage != nil {
		existing.MarkupPercentage = *req.MarkupPercentage
	}
	if req.MinStockLevel != nil {
		existing.MinStockLevel = *req.MinStockLevel
	}
	if req.MaxStockLevel != nil {
		existing.MaxStockLevel = *req.MaxStockLevel
	}
	if req.LocationRack != nil {
		existing.LocationRack = req.LocationRack
	}
	if req.Barcode != nil {
		// Validate barcode uniqueness if changed
		if *req.Barcode != "" {
			exists, err := s.productRepo.IsBarcodeExists(ctx, *req.Barcode, id)
			if err != nil {
				return nil, fmt.Errorf("failed to check barcode uniqueness: %w", err)
			}
			if exists {
				return nil, fmt.Errorf("barcode %s already exists", *req.Barcode)
			}
		}
		existing.Barcode = req.Barcode
	}
	if req.Weight != nil {
		existing.Weight = req.Weight
	}
	if req.Dimensions != nil {
		existing.Dimensions = req.Dimensions
	}
	if req.ProductImage != nil {
		existing.ProductImage = req.ProductImage
	}
	if req.Notes != nil {
		existing.Notes = req.Notes
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	// Recalculate markup if prices changed
	if req.CostPrice != nil || req.SellingPrice != nil {
		existing.UpdateMarkupPercentage()
	}

	// Update product
	updatedProduct, err := s.productRepo.Update(ctx, id, existing)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return updatedProduct, nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(ctx context.Context, id int) (*products.ProductSparePart, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}

// GetProductByCode retrieves a product by code
func (s *ProductService) GetProductByCode(ctx context.Context, code string) (*products.ProductSparePart, error) {
	product, err := s.productRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by code: %w", err)
	}
	return product, nil
}

// GetProductByBarcode retrieves a product by barcode
func (s *ProductService) GetProductByBarcode(ctx context.Context, barcode string) (*products.ProductSparePart, error) {
	product, err := s.productRepo.GetByBarcode(ctx, barcode)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by barcode: %w", err)
	}
	return product, nil
}

// ListProducts retrieves a paginated list of products
func (s *ProductService) ListProducts(ctx context.Context, params *products.ProductSparePartFilterParams) (*common.PaginatedResponse, error) {
	response, err := s.productRepo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	return response, nil
}

// GetLowStockProducts retrieves products below minimum stock level
func (s *ProductService) GetLowStockProducts(ctx context.Context, params *products.ProductSparePartFilterParams) (*common.PaginatedResponse, error) {
	response, err := s.productRepo.GetLowStockProducts(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock products: %w", err)
	}
	return response, nil
}

// SearchProducts searches products by query
func (s *ProductService) SearchProducts(ctx context.Context, query string, params *products.ProductSparePartFilterParams) (*common.PaginatedResponse, error) {
	response, err := s.productRepo.Search(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}
	return response, nil
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	// Check if product exists
	_, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	// TODO: Add business logic to check if product can be deleted
	// For example, check if it's referenced in any purchase orders or sales

	err = s.productRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

// AdjustStock adjusts the stock quantity of a product with audit trail
func (s *ProductService) AdjustStock(ctx context.Context, productID int, newQuantity int, reason string, adjustedBy int) error {
	// Get current product
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	quantityChange := newQuantity - product.StockQuantity

	// Create stock movement
	movement := &products.StockMovement{
		ProductID:      productID,
		MovementType:   products.MovementTypeAdjustment,
		ReferenceType:  products.ReferenceTypeAdjustment,
		ReferenceID:    productID, // Using product ID as reference for manual adjustments
		QuantityMoved:  quantityChange,
		UnitCost:       product.CostPrice,
		MovementDate:   time.Now(),
		ProcessedBy:    adjustedBy,
		MovementReason: &reason,
	}

	movement.CalculateTotalValue()

	// Update stock with movement
	err = s.productRepo.UpdateStockWithMovement(ctx, productID, quantityChange, movement)
	if err != nil {
		return fmt.Errorf("failed to adjust stock: %w", err)
	}

	return nil
}

// GetStockMovementHistory retrieves stock movement history for a product
func (s *ProductService) GetStockMovementHistory(ctx context.Context, productID int, limit int) ([]products.StockMovement, error) {
	// Verify product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	movements, err := s.stockRepo.GetMovementHistory(ctx, productID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock movement history: %w", err)
	}

	return movements, nil
}

// ValidateProduct validates business rules for a product
func (s *ProductService) ValidateProduct(product *products.ProductSparePart) error {
	if product.SellingPrice < product.CostPrice {
		return fmt.Errorf("selling price cannot be less than cost price")
	}

	if product.MinStockLevel < 0 {
		return fmt.Errorf("minimum stock level cannot be negative")
	}

	if product.MaxStockLevel < product.MinStockLevel {
		return fmt.Errorf("maximum stock level cannot be less than minimum stock level")
	}

	if product.StockQuantity < 0 {
		return fmt.Errorf("stock quantity cannot be negative")
	}

	return nil
}