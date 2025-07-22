package inventory

import (
	"context"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// StockAdjustmentService handles business logic for stock adjustments
type StockAdjustmentService struct {
	stockAdjustmentRepo interfaces.StockAdjustmentRepository
	stockMovementRepo   interfaces.StockMovementRepository
	productRepo         interfaces.ProductSparePartRepository
}

// NewStockAdjustmentService creates a new stock adjustment service
func NewStockAdjustmentService(
	stockAdjustmentRepo interfaces.StockAdjustmentRepository,
	stockMovementRepo interfaces.StockMovementRepository,
	productRepo interfaces.ProductSparePartRepository,
) *StockAdjustmentService {
	return &StockAdjustmentService{
		stockAdjustmentRepo: stockAdjustmentRepo,
		stockMovementRepo:   stockMovementRepo,
		productRepo:         productRepo,
	}
}

// CreateStockAdjustment creates a new stock adjustment
func (s *StockAdjustmentService) CreateStockAdjustment(ctx context.Context, req *inventory.StockAdjustmentCreateRequest, userID int) (*inventory.StockAdjustment, error) {
	// Validate product exists
	product, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Create stock adjustment
	adjustment := &inventory.StockAdjustment{
		ProductID:                req.ProductID,
		AdjustmentType:           req.AdjustmentType,
		QuantitySystem:           product.StockQuantity,
		QuantityPhysical:         req.QuantityPhysical,
		AdjustmentReason:         req.AdjustmentReason,
		Notes:                    req.Notes,
		AdjustmentDate:           time.Now(),
		SupportingDocumentsJSON:  req.SupportingDocumentsJSON,
		CreatedBy:                userID,
	}

	// Calculate variance and cost impact
	adjustment.CalculateVariance()
	adjustment.CalculateCostImpact(product.CostPrice)

	// Create adjustment record
	createdAdjustment, err := s.stockAdjustmentRepo.Create(ctx, adjustment)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock adjustment: %w", err)
	}

	return createdAdjustment, nil
}

// GetStockAdjustment retrieves a stock adjustment by ID
func (s *StockAdjustmentService) GetStockAdjustment(ctx context.Context, id int) (*inventory.StockAdjustment, error) {
	return s.stockAdjustmentRepo.GetByID(ctx, id)
}

// ListStockAdjustments retrieves stock adjustments with filtering
func (s *StockAdjustmentService) ListStockAdjustments(ctx context.Context, params *inventory.StockAdjustmentFilterParams) ([]inventory.StockAdjustmentListItem, int, error) {
	return s.stockAdjustmentRepo.List(ctx, params)
}

// UpdateStockAdjustment updates a stock adjustment (only if not approved)
func (s *StockAdjustmentService) UpdateStockAdjustment(ctx context.Context, id int, req *inventory.StockAdjustmentUpdateRequest) (*inventory.StockAdjustment, error) {
	adjustment, err := s.stockAdjustmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if adjustment.ApprovedBy != nil {
		return nil, fmt.Errorf("cannot update approved stock adjustment")
	}

	// Update fields
	if req.AdjustmentType != nil {
		adjustment.AdjustmentType = *req.AdjustmentType
	}
	if req.QuantityPhysical != nil {
		adjustment.QuantityPhysical = *req.QuantityPhysical
	}
	if req.AdjustmentReason != nil {
		adjustment.AdjustmentReason = *req.AdjustmentReason
	}
	if req.Notes != nil {
		adjustment.Notes = req.Notes
	}
	if req.SupportingDocumentsJSON != nil {
		adjustment.SupportingDocumentsJSON = req.SupportingDocumentsJSON
	}

	// Recalculate variance and cost impact
	product, err := s.productRepo.GetByID(ctx, adjustment.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	adjustment.QuantitySystem = product.StockQuantity
	adjustment.CalculateVariance()
	adjustment.CalculateCostImpact(product.CostPrice)

	return s.stockAdjustmentRepo.Update(ctx, id, adjustment)
}

// ApproveStockAdjustment approves a stock adjustment and creates stock movement
func (s *StockAdjustmentService) ApproveStockAdjustment(ctx context.Context, id int, approverID int) (*inventory.StockAdjustment, error) {
	adjustment, err := s.stockAdjustmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if adjustment.ApprovedBy != nil {
		return nil, fmt.Errorf("stock adjustment already approved")
	}

	// Approve the adjustment
	err = s.stockAdjustmentRepo.UpdateApproval(ctx, id, approverID)
	if err != nil {
		return nil, fmt.Errorf("failed to approve adjustment: %w", err)
	}

	// Create stock movement if there's a variance
	if adjustment.QuantityVariance != 0 {
		err = s.createAdjustmentMovement(ctx, adjustment, approverID)
		if err != nil {
			return nil, fmt.Errorf("failed to create adjustment movement: %w", err)
		}
	}

	// Get updated adjustment
	return s.stockAdjustmentRepo.GetByID(ctx, id)
}

// DeleteStockAdjustment deletes a stock adjustment (only if not approved)
func (s *StockAdjustmentService) DeleteStockAdjustment(ctx context.Context, id int) error {
	return s.stockAdjustmentRepo.Delete(ctx, id)
}

// GetPendingAdjustments retrieves adjustments pending approval
func (s *StockAdjustmentService) GetPendingAdjustments(ctx context.Context, page, limit int) ([]inventory.StockAdjustmentListItem, int, error) {
	return s.stockAdjustmentRepo.GetPendingApproval(ctx, page, limit)
}

// GetAdjustmentsByProduct retrieves adjustments for a specific product
func (s *StockAdjustmentService) GetAdjustmentsByProduct(ctx context.Context, productID int, page, limit int) ([]inventory.StockAdjustmentListItem, int, error) {
	return s.stockAdjustmentRepo.GetByProduct(ctx, productID, page, limit)
}

// GetAdjustmentsByType retrieves adjustments by type
func (s *StockAdjustmentService) GetAdjustmentsByType(ctx context.Context, adjustmentType inventory.AdjustmentType, page, limit int) ([]inventory.StockAdjustmentListItem, int, error) {
	return s.stockAdjustmentRepo.GetByType(ctx, adjustmentType, page, limit)
}

// GetAdjustmentSummary retrieves adjustment summary for reporting
func (s *StockAdjustmentService) GetAdjustmentSummary(ctx context.Context, startDate, endDate string) (map[string]interface{}, error) {
	return s.stockAdjustmentRepo.GetAdjustmentSummary(ctx, startDate, endDate)
}

// CreatePhysicalCountAdjustment creates an adjustment based on physical count
func (s *StockAdjustmentService) CreatePhysicalCountAdjustment(ctx context.Context, productID int, physicalCount int, userID int, notes *string) (*inventory.StockAdjustment, error) {
	req := &inventory.StockAdjustmentCreateRequest{
		ProductID:        productID,
		AdjustmentType:   inventory.AdjustmentTypePhysicalCount,
		QuantityPhysical: physicalCount,
		AdjustmentReason: "Physical count adjustment",
		Notes:            notes,
	}

	return s.CreateStockAdjustment(ctx, req, userID)
}

// CreateDamageAdjustment creates an adjustment for damaged goods
func (s *StockAdjustmentService) CreateDamageAdjustment(ctx context.Context, productID int, damagedQuantity int, userID int, reason string, notes *string) (*inventory.StockAdjustment, error) {
	// Get current stock
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	req := &inventory.StockAdjustmentCreateRequest{
		ProductID:        productID,
		AdjustmentType:   inventory.AdjustmentTypeDamage,
		QuantityPhysical: product.StockQuantity - damagedQuantity,
		AdjustmentReason: reason,
		Notes:            notes,
	}

	return s.CreateStockAdjustment(ctx, req, userID)
}

// CreateWriteOffAdjustment creates an adjustment for write-off
func (s *StockAdjustmentService) CreateWriteOffAdjustment(ctx context.Context, productID int, writeOffQuantity int, userID int, reason string, notes *string) (*inventory.StockAdjustment, error) {
	// Get current stock
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	req := &inventory.StockAdjustmentCreateRequest{
		ProductID:        productID,
		AdjustmentType:   inventory.AdjustmentTypeWriteOff,
		QuantityPhysical: product.StockQuantity - writeOffQuantity,
		AdjustmentReason: reason,
		Notes:            notes,
	}

	return s.CreateStockAdjustment(ctx, req, userID)
}

// createAdjustmentMovement creates a stock movement for approved adjustment
func (s *StockAdjustmentService) createAdjustmentMovement(ctx context.Context, adjustment *inventory.StockAdjustment, approverID int) error {
	product, err := s.productRepo.GetByID(ctx, adjustment.ProductID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	var movementType inventory.MovementType
	quantityMoved := adjustment.QuantityVariance

	if quantityMoved > 0 {
		movementType = inventory.MovementTypeIn
	} else {
		movementType = inventory.MovementTypeOut
		quantityMoved = -quantityMoved // Make positive for movement
	}

	movement := &inventory.StockMovement{
		ProductID:       adjustment.ProductID,
		MovementType:    movementType,
		ReferenceType:   inventory.ReferenceTypeAdjustment,
		ReferenceID:     adjustment.AdjustmentID,
		QuantityBefore:  adjustment.QuantitySystem,
		QuantityMoved:   quantityMoved,
		QuantityAfter:   adjustment.QuantityPhysical,
		UnitCost:        product.CostPrice,
		MovementDate:    time.Now(),
		ProcessedBy:     approverID,
		MovementReason:  fmt.Sprintf("Stock Adjustment - %s", adjustment.AdjustmentType),
		Notes:           adjustment.Notes,
	}

	movement.CalculateTotalValue()

	// Create movement
	_, err = s.stockMovementRepo.Create(ctx, movement)
	if err != nil {
		return fmt.Errorf("failed to create movement: %w", err)
	}

	// Update product stock
	err = s.productRepo.UpdateStock(ctx, adjustment.ProductID, adjustment.QuantityPhysical)
	if err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	return nil
}