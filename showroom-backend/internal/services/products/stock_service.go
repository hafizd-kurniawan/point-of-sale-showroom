package products

import (
	"context"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// StockService handles business logic for stock management
type StockService struct {
	stockMovementRepo interfaces.StockMovementRepository
	productRepo       interfaces.ProductSparePartRepository
}

// NewStockService creates a new stock service
func NewStockService(
	stockMovementRepo interfaces.StockMovementRepository,
	productRepo interfaces.ProductSparePartRepository,
) *StockService {
	return &StockService{
		stockMovementRepo: stockMovementRepo,
		productRepo:       productRepo,
	}
}

// GetStockMovements retrieves stock movements with pagination and filtering
func (s *StockService) GetStockMovements(ctx context.Context, params *products.StockMovementFilterParams) (*common.PaginatedResponse, error) {
	return s.stockMovementRepo.List(ctx, params)
}

// GetStockMovementByID retrieves a stock movement by ID
func (s *StockService) GetStockMovementByID(ctx context.Context, id int) (*products.StockMovement, error) {
	return s.stockMovementRepo.GetByID(ctx, id)
}

// GetStockMovementsByProduct retrieves stock movements for a specific product
func (s *StockService) GetStockMovementsByProduct(ctx context.Context, productID int, params *products.StockMovementFilterParams) (*common.PaginatedResponse, error) {
	// Validate product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	return s.stockMovementRepo.GetByProductID(ctx, productID, params)
}

// GetStockMovementsByReference retrieves stock movements by reference type and ID
func (s *StockService) GetStockMovementsByReference(ctx context.Context, referenceType products.ReferenceType, referenceID int) ([]products.StockMovement, error) {
	if !referenceType.IsValid() {
		return nil, fmt.Errorf("invalid reference type: %s", referenceType)
	}

	return s.stockMovementRepo.GetByReferenceID(ctx, referenceType, referenceID)
}

// GetProductMovementHistory retrieves movement history for a product
func (s *StockService) GetProductMovementHistory(ctx context.Context, productID int, limit int) ([]products.StockMovement, error) {
	// Validate product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if limit <= 0 || limit > 100 {
		limit = 50 // Default limit
	}

	return s.stockMovementRepo.GetMovementHistory(ctx, productID, limit)
}

// GetCurrentStock retrieves current stock quantity for a product
func (s *StockService) GetCurrentStock(ctx context.Context, productID int) (int, error) {
	// Validate product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return 0, fmt.Errorf("product not found: %w", err)
	}

	return s.stockMovementRepo.GetCurrentStock(ctx, productID)
}

// CreateManualStockMovement creates a manual stock movement with business validation
func (s *StockService) CreateManualStockMovement(ctx context.Context, req *products.StockMovementCreateRequest, createdBy int) (*products.StockMovement, error) {
	// Validate product exists
	_, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Validate movement type
	if !req.MovementType.IsValid() {
		return nil, fmt.Errorf("invalid movement type: %s", req.MovementType)
	}

	// Validate reference type
	if !req.ReferenceType.IsValid() {
		return nil, fmt.Errorf("invalid reference type: %s", req.ReferenceType)
	}

	// Get current stock
	currentStock, err := s.stockMovementRepo.GetCurrentStock(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current stock: %w", err)
	}

	// Validate stock quantity for outbound movements
	if (req.MovementType == products.MovementTypeOut || 
		req.MovementType == products.MovementTypeDamage || 
		req.MovementType == products.MovementTypeExpired ||
		(req.MovementType == products.MovementTypeAdjustment && req.QuantityMoved < 0)) {
		
		quantityToDeduct := req.QuantityMoved
		if req.MovementType == products.MovementTypeAdjustment && req.QuantityMoved < 0 {
			quantityToDeduct = -req.QuantityMoved
		}
		
		if currentStock < quantityToDeduct {
			return nil, fmt.Errorf("insufficient stock: current=%d, required=%d", currentStock, quantityToDeduct)
		}
	}

	// Calculate quantities
	var quantityAfter int
	switch req.MovementType {
	case products.MovementTypeIn, products.MovementTypeReturn:
		quantityAfter = currentStock + req.QuantityMoved
	case products.MovementTypeOut, products.MovementTypeDamage, products.MovementTypeExpired:
		quantityAfter = currentStock - req.QuantityMoved
	case products.MovementTypeAdjustment:
		quantityAfter = currentStock + req.QuantityMoved
	default:
		quantityAfter = currentStock + req.QuantityMoved
	}

	// Create movement record
	movement := &products.StockMovement{
		ProductID:      req.ProductID,
		MovementType:   req.MovementType,
		ReferenceType:  req.ReferenceType,
		ReferenceID:    req.ReferenceID,
		QuantityBefore: currentStock,
		QuantityMoved:  req.QuantityMoved,
		QuantityAfter:  quantityAfter,
		UnitCost:       req.UnitCost,
		LocationFrom:   req.LocationFrom,
		LocationTo:     req.LocationTo,
		MovementDate:   getMovementDate(req.MovementDate),
		ProcessedBy:    createdBy,
		MovementReason: req.MovementReason,
		Notes:          req.Notes,
	}

	// Calculate total value
	movement.CalculateTotalValue()

	// Validate quantities
	if !movement.ValidateQuantities() {
		return nil, fmt.Errorf("invalid quantity calculations")
	}

	// Create movement
	createdMovement, err := s.stockMovementRepo.Create(ctx, movement)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock movement: %w", err)
	}

	// Update product stock quantity
	err = s.productRepo.UpdateStock(ctx, req.ProductID, quantityAfter)
	if err != nil {
		return nil, fmt.Errorf("failed to update product stock: %w", err)
	}

	return createdMovement, nil
}

// BulkCreateMovements creates multiple stock movements in a transaction
func (s *StockService) BulkCreateMovements(ctx context.Context, movements []products.StockMovement) error {
	if len(movements) == 0 {
		return nil
	}

	// Validate all movements
	for i, movement := range movements {
		// Validate product exists
		_, err := s.productRepo.GetByID(ctx, movement.ProductID)
		if err != nil {
			return fmt.Errorf("product not found for movement %d: %w", i, err)
		}

		// Validate movement type
		if !movement.MovementType.IsValid() {
			return fmt.Errorf("invalid movement type for movement %d: %s", i, movement.MovementType)
		}

		// Validate reference type
		if !movement.ReferenceType.IsValid() {
			return fmt.Errorf("invalid reference type for movement %d: %s", i, movement.ReferenceType)
		}

		// Validate quantities
		if !movement.ValidateQuantities() {
			return fmt.Errorf("invalid quantity calculations for movement %d", i)
		}
	}

	return s.stockMovementRepo.BulkCreateMovements(ctx, movements)
}

// GetLowStockProducts retrieves products with low stock levels
func (s *StockService) GetLowStockProducts(ctx context.Context, params *products.ProductSparePartFilterParams) (*common.PaginatedResponse, error) {
	return s.productRepo.GetLowStockProducts(ctx, params)
}

// GetStockSummary retrieves stock summary for products
func (s *StockService) GetStockSummary(ctx context.Context, productID *int) (map[string]interface{}, error) {
	summary := make(map[string]interface{})

	if productID != nil {
		// Get specific product stock summary
		product, err := s.productRepo.GetByID(ctx, *productID)
		if err != nil {
			return nil, fmt.Errorf("product not found: %w", err)
		}

		currentStock, err := s.stockMovementRepo.GetCurrentStock(ctx, *productID)
		if err != nil {
			return nil, fmt.Errorf("failed to get current stock: %w", err)
		}

		summary["product_id"] = *productID
		summary["product_name"] = product.ProductName
		summary["product_code"] = product.ProductCode
		summary["current_stock"] = currentStock
		summary["min_stock_level"] = product.MinStockLevel
		summary["max_stock_level"] = product.MaxStockLevel
		summary["is_low_stock"] = currentStock <= product.MinStockLevel
		summary["stock_value"] = float64(currentStock) * product.CostPrice

		// Get recent movements
		recentMovements, err := s.stockMovementRepo.GetMovementHistory(ctx, *productID, 10)
		if err != nil {
			return nil, fmt.Errorf("failed to get recent movements: %w", err)
		}
		summary["recent_movements"] = recentMovements
	} else {
		// Get overall stock summary
		// This would require additional repository methods for aggregate queries
		summary["message"] = "Overall stock summary not implemented yet"
	}

	return summary, nil
}

// helper function to get movement date
func getMovementDate(date *time.Time) time.Time {
	if date != nil {
		return *date
	}
	return time.Now()
}