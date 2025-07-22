package inventory

import (
	"context"
	"fmt"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// ProductService provides business logic for product operations
type ProductService struct {
	productRepo   interfaces.ProductSparePartRepository
	movementRepo  interfaces.StockMovementRepository
	codeGenerator *inventory.CodeGenerator
}

// NewProductService creates a new product service
func NewProductService(
	productRepo interfaces.ProductSparePartRepository,
	movementRepo interfaces.StockMovementRepository,
) *ProductService {
	return &ProductService{
		productRepo:   productRepo,
		movementRepo:  movementRepo,
		codeGenerator: inventory.NewCodeGenerator(),
	}
}

// Create creates a new product with auto-generated code
func (s *ProductService) Create(ctx context.Context, req *inventory.ProductSparePartCreateRequest, createdBy int) (*inventory.ProductSparePart, error) {
	// Generate product code
	productCode, err := s.codeGenerator.GetNextProductCode(func() (int, error) {
		return s.productRepo.GetLastProductID(ctx)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate product code: %w", err)
	}

	// Check if barcode is unique (if provided)
	if req.Barcode != nil && *req.Barcode != "" {
		exists, err := s.productRepo.ExistsByBarcode(ctx, *req.Barcode)
		if err != nil {
			return nil, fmt.Errorf("failed to check barcode uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("barcode already exists")
		}
	}

	// Create product entity
	product := &inventory.ProductSparePart{
		ProductCode:      productCode,
		ProductName:      req.ProductName,
		Description:      req.Description,
		BrandID:          req.BrandID,
		CategoryID:       req.CategoryID,
		UnitMeasure:      req.UnitMeasure,
		CostPrice:        req.CostPrice,
		SellingPrice:     req.SellingPrice,
		MarkupPercentage: req.MarkupPercentage,
		StockQuantity:    req.StockQuantity,
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

	// Calculate markup if not provided
	if product.MarkupPercentage == nil {
		product.UpdateMarkupFromPrices()
	}

	// Create product in repository
	createdProduct, err := s.productRepo.Create(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Create initial stock movement if stock quantity > 0
	if req.StockQuantity > 0 {
		movement := &inventory.StockMovement{
			ProductID:      createdProduct.ProductID,
			MovementType:   inventory.MovementTypeIn,
			ReferenceType:  inventory.ReferenceTypeAdjustment,
			QuantityBefore: 0,
			QuantityMoved:  req.StockQuantity,
			QuantityAfter:  req.StockQuantity,
			UnitCost:       req.CostPrice,
			ProcessedBy:    createdBy,
			MovementReason: func() *string { s := "Initial stock"; return &s }(),
			Notes:          func() *string { s := "Initial stock entry for new product"; return &s }(),
		}
		movement.CalculateTotalValue()

		_, err = s.movementRepo.Create(ctx, movement)
		if err != nil {
			// Log error but don't fail the product creation
			// In production, you might want to use a transaction to ensure consistency
		}
	}

	return createdProduct, nil
}

// GetByID retrieves a product by ID
func (s *ProductService) GetByID(ctx context.Context, id int) (*inventory.ProductSparePart, error) {
	return s.productRepo.GetByID(ctx, id)
}

// GetByCode retrieves a product by code
func (s *ProductService) GetByCode(ctx context.Context, code string) (*inventory.ProductSparePart, error) {
	return s.productRepo.GetByCode(ctx, code)
}

// GetByBarcode retrieves a product by barcode
func (s *ProductService) GetByBarcode(ctx context.Context, barcode string) (*inventory.ProductSparePart, error) {
	return s.productRepo.GetByBarcode(ctx, barcode)
}

// Update updates a product
func (s *ProductService) Update(ctx context.Context, id int, req *inventory.ProductSparePartUpdateRequest) (*inventory.ProductSparePart, error) {
	// Get existing product
	existing, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Check barcode uniqueness if changing
	if req.Barcode != nil && *req.Barcode != "" {
		if existing.Barcode == nil || *existing.Barcode != *req.Barcode {
			exists, err := s.productRepo.ExistsByBarcodeExcludingID(ctx, *req.Barcode, id)
			if err != nil {
				return nil, fmt.Errorf("failed to check barcode uniqueness: %w", err)
			}
			if exists {
				return nil, fmt.Errorf("barcode already exists")
			}
		}
	}

	// Update fields
	if req.ProductName != nil {
		existing.ProductName = *req.ProductName
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.BrandID != nil {
		existing.BrandID = req.BrandID
	}
	if req.CategoryID != nil {
		existing.CategoryID = req.CategoryID
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
		existing.MarkupPercentage = req.MarkupPercentage
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
		existing.Barcode = req.Barcode
	}
	if req.Weight != nil {
		existing.Weight = req.Weight
	}
	if req.Dimensions != nil {
		existing.Dimensions = req.Dimensions
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	if req.ProductImage != nil {
		existing.ProductImage = req.ProductImage
	}
	if req.Notes != nil {
		existing.Notes = req.Notes
	}

	// Recalculate markup if prices changed
	if req.CostPrice != nil || req.SellingPrice != nil {
		existing.UpdateMarkupFromPrices()
	}

	// Handle stock quantity change
	if req.StockQuantity != nil && *req.StockQuantity != existing.StockQuantity {
		// This would typically require an adjustment through proper stock movement
		// For now, we'll just update the quantity
		existing.StockQuantity = *req.StockQuantity
	}

	return s.productRepo.Update(ctx, id, existing)
}

// Delete deletes a product
func (s *ProductService) Delete(ctx context.Context, id int) error {
	// Check if product has any stock movements or is referenced elsewhere
	// This would be a business rule check
	return s.productRepo.Delete(ctx, id)
}

// List retrieves products with filtering and pagination
func (s *ProductService) List(ctx context.Context, params *inventory.ProductSparePartFilterParams) ([]inventory.ProductSparePartListItem, int, error) {
	return s.productRepo.List(ctx, params)
}

// GetLowStockProducts retrieves products below minimum stock level
func (s *ProductService) GetLowStockProducts(ctx context.Context, page, limit int) ([]inventory.ProductSparePartListItem, int, error) {
	return s.productRepo.GetLowStockProducts(ctx, page, limit)
}

// Search searches products
func (s *ProductService) Search(ctx context.Context, query string, page, limit int) ([]inventory.ProductSparePartListItem, int, error) {
	return s.productRepo.Search(ctx, query, page, limit)
}

// UpdateStock updates product stock with movement tracking
func (s *ProductService) UpdateStock(ctx context.Context, productID int, newQuantity int, movementType inventory.MovementType, refType inventory.ReferenceType, refID *int, processedBy int, reason string) error {
	// Get current stock
	currentStock, err := s.productRepo.GetStockQuantity(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get current stock: %w", err)
	}

	// Calculate movement
	quantityMoved := newQuantity - currentStock

	// Get product for cost information
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Create stock movement
	movement := &inventory.StockMovement{
		ProductID:      productID,
		MovementType:   movementType,
		ReferenceType:  refType,
		ReferenceID:    refID,
		QuantityBefore: currentStock,
		QuantityMoved:  quantityMoved,
		QuantityAfter:  newQuantity,
		UnitCost:       product.CostPrice,
		ProcessedBy:    processedBy,
		MovementReason: &reason,
	}
	movement.CalculateTotalValue()

	// Create movement record
	_, err = s.movementRepo.Create(ctx, movement)
	if err != nil {
		return fmt.Errorf("failed to create stock movement: %w", err)
	}

	// Update product stock
	return s.productRepo.UpdateStock(ctx, productID, newQuantity)
}

// AdjustStock adjusts product stock with reason
func (s *ProductService) AdjustStock(ctx context.Context, productID int, newQuantity int, processedBy int, reason string, notes *string) error {
	return s.UpdateStock(ctx, productID, newQuantity, inventory.MovementTypeAdjustment, inventory.ReferenceTypeAdjustment, nil, processedBy, reason)
}

// GetInventoryValue calculates total inventory value
func (s *ProductService) GetInventoryValue(ctx context.Context) (float64, error) {
	return s.productRepo.GetInventoryValue(ctx)
}

// ValidateProductCode validates product code format
func (s *ProductService) ValidateProductCode(code string) bool {
	return s.codeGenerator.ValidateProductCode(code)
}

// CheckLowStockAlerts checks for products below minimum stock level
func (s *ProductService) CheckLowStockAlerts(ctx context.Context) ([]inventory.ProductSparePartListItem, error) {
	products, _, err := s.productRepo.GetLowStockProducts(ctx, 1, 100) // Get first 100 low stock items
	return products, err
}