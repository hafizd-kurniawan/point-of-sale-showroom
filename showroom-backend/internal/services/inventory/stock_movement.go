package inventory

import (
	"context"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// StockMovementService handles business logic for stock movements
type StockMovementService struct {
	stockMovementRepo interfaces.StockMovementRepository
	productRepo       interfaces.ProductSparePartRepository
}

// NewStockMovementService creates a new stock movement service
func NewStockMovementService(
	stockMovementRepo interfaces.StockMovementRepository,
	productRepo interfaces.ProductSparePartRepository,
) *StockMovementService {
	return &StockMovementService{
		stockMovementRepo: stockMovementRepo,
		productRepo:       productRepo,
	}
}

// CreateStockMovement creates a new stock movement and updates product stock
func (s *StockMovementService) CreateStockMovement(ctx context.Context, req *inventory.StockMovementCreateRequest, userID int) (*inventory.StockMovement, error) {
	// Validate product exists
	product, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Validate movement quantities
	if req.MovementType == inventory.MovementTypeOut && product.StockQuantity < req.QuantityMoved {
		return nil, fmt.Errorf("insufficient stock: available %d, requested %d", product.StockQuantity, req.QuantityMoved)
	}

	// Calculate before/after quantities
	quantityBefore := product.StockQuantity
	var quantityAfter int

	switch req.MovementType {
	case inventory.MovementTypeIn, inventory.MovementTypeReturn:
		quantityAfter = quantityBefore + req.QuantityMoved
	case inventory.MovementTypeOut, inventory.MovementTypeDamage, inventory.MovementTypeExpired:
		quantityAfter = quantityBefore - req.QuantityMoved
	case inventory.MovementTypeTransfer:
		// For transfer, we need location information
		if req.LocationFrom == nil || req.LocationTo == nil {
			return nil, fmt.Errorf("location information required for transfer movements")
		}
		quantityAfter = quantityBefore // No net change in total stock for transfer
	case inventory.MovementTypeAdjustment:
		// For adjustments, quantity moved can be positive or negative
		quantityAfter = quantityBefore + req.QuantityMoved
	default:
		return nil, fmt.Errorf("invalid movement type: %s", req.MovementType)
	}

	if quantityAfter < 0 {
		return nil, fmt.Errorf("movement would result in negative stock")
	}

	// Create stock movement
	movement := &inventory.StockMovement{
		ProductID:       req.ProductID,
		MovementType:    req.MovementType,
		ReferenceType:   req.ReferenceType,
		ReferenceID:     req.ReferenceID,
		QuantityBefore:  quantityBefore,
		QuantityMoved:   req.QuantityMoved,
		QuantityAfter:   quantityAfter,
		UnitCost:        req.UnitCost,
		LocationFrom:    req.LocationFrom,
		LocationTo:      req.LocationTo,
		MovementDate:    time.Now(),
		ProcessedBy:     userID,
		MovementReason:  req.MovementReason,
		Notes:           req.Notes,
	}

	// Calculate total value
	movement.CalculateTotalValue()

	// Create movement record
	createdMovement, err := s.stockMovementRepo.Create(ctx, movement)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock movement: %w", err)
	}

	// Update product stock (except for transfers which don't change total stock)
	if req.MovementType != inventory.MovementTypeTransfer {
		err = s.productRepo.UpdateStock(ctx, req.ProductID, quantityAfter)
		if err != nil {
			return nil, fmt.Errorf("failed to update product stock: %w", err)
		}
	}

	return createdMovement, nil
}

// GetStockMovement retrieves a stock movement by ID
func (s *StockMovementService) GetStockMovement(ctx context.Context, id int) (*inventory.StockMovement, error) {
	return s.stockMovementRepo.GetByID(ctx, id)
}

// ListStockMovements retrieves stock movements with filtering
func (s *StockMovementService) ListStockMovements(ctx context.Context, params *inventory.StockMovementFilterParams) ([]inventory.StockMovementListItem, int, error) {
	return s.stockMovementRepo.List(ctx, params)
}

// GetProductStockHistory retrieves stock movement history for a product
func (s *StockMovementService) GetProductStockHistory(ctx context.Context, productID int, page, limit int) ([]inventory.StockMovementListItem, int, error) {
	return s.stockMovementRepo.GetByProduct(ctx, productID, page, limit)
}

// GetAuditTrail retrieves audit trail for a product
func (s *StockMovementService) GetAuditTrail(ctx context.Context, productID int, page, limit int) ([]inventory.StockMovementListItem, int, error) {
	return s.stockMovementRepo.GetAuditTrail(ctx, productID, page, limit)
}

// GetMovementSummary retrieves movement summary for a product within date range
func (s *StockMovementService) GetMovementSummary(ctx context.Context, productID int, startDate, endDate string) (map[string]interface{}, error) {
	return s.stockMovementRepo.GetMovementSummary(ctx, productID, startDate, endDate)
}

// GetValueMovements retrieves total value of movements within date range
func (s *StockMovementService) GetValueMovements(ctx context.Context, startDate, endDate string) (float64, error) {
	return s.stockMovementRepo.GetValueMovements(ctx, startDate, endDate)
}

// CreateManualAdjustment creates a manual stock adjustment movement
func (s *StockMovementService) CreateManualAdjustment(ctx context.Context, productID int, adjustmentQuantity int, reason string, userID int) (*inventory.StockMovement, error) {
	req := &inventory.StockMovementCreateRequest{
		ProductID:       productID,
		MovementType:    inventory.MovementTypeAdjustment,
		ReferenceType:   inventory.ReferenceTypeAdjustment,
		ReferenceID:     0, // No specific reference for manual adjustments
		QuantityMoved:   adjustmentQuantity,
		UnitCost:        0, // Will be populated from product cost
		MovementReason:  reason,
		Notes:           &reason,
	}

	// Get product to get current cost
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	req.UnitCost = product.CostPrice

	return s.CreateStockMovement(ctx, req, userID)
}

// CreateTransfer creates a stock transfer between locations
func (s *StockMovementService) CreateTransfer(ctx context.Context, productID int, quantity int, fromLocation, toLocation string, reason string, userID int) (*inventory.StockMovement, error) {
	req := &inventory.StockMovementCreateRequest{
		ProductID:       productID,
		MovementType:    inventory.MovementTypeTransfer,
		ReferenceType:   inventory.ReferenceTypeTransfer,
		ReferenceID:     0,
		QuantityMoved:   quantity,
		UnitCost:        0, // Will be populated from product cost
		LocationFrom:    &fromLocation,
		LocationTo:      &toLocation,
		MovementReason:  reason,
		Notes:           &reason,
	}

	// Get product to get current cost
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	req.UnitCost = product.CostPrice

	return s.CreateStockMovement(ctx, req, userID)
}

// GetMovementsByType retrieves movements by type
func (s *StockMovementService) GetMovementsByType(ctx context.Context, movementType inventory.MovementType, page, limit int) ([]inventory.StockMovementListItem, int, error) {
	return s.stockMovementRepo.GetByType(ctx, movementType, page, limit)
}

// GetMovementsByReference retrieves movements by reference
func (s *StockMovementService) GetMovementsByReference(ctx context.Context, refType inventory.ReferenceType, refID int, page, limit int) ([]inventory.StockMovementListItem, int, error) {
	return s.stockMovementRepo.GetByReference(ctx, refType, refID, page, limit)
}

// GetMovementsByUser retrieves movements processed by a user
func (s *StockMovementService) GetMovementsByUser(ctx context.Context, userID int, page, limit int) ([]inventory.StockMovementListItem, int, error) {
	return s.stockMovementRepo.GetMovementsByUser(ctx, userID, page, limit)
}

// GetMovementsByDateRange retrieves movements within date range
func (s *StockMovementService) GetMovementsByDateRange(ctx context.Context, startDate, endDate string, page, limit int) ([]inventory.StockMovementListItem, int, error) {
	return s.stockMovementRepo.GetByDateRange(ctx, startDate, endDate, page, limit)
}