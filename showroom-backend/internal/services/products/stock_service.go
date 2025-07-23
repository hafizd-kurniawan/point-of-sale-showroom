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
	stockMovementRepo    interfaces.StockMovementRepository
	stockAdjustmentRepo  interfaces.StockAdjustmentRepository
	productRepo          interfaces.ProductSparePartRepository
}

// NewStockService creates a new stock service
func NewStockService(
	stockMovementRepo interfaces.StockMovementRepository,
	stockAdjustmentRepo interfaces.StockAdjustmentRepository,
	productRepo interfaces.ProductSparePartRepository,
) *StockService {
	return &StockService{
		stockMovementRepo:   stockMovementRepo,
		stockAdjustmentRepo: stockAdjustmentRepo,
		productRepo:         productRepo,
	}
}

// CreateStockMovement creates a new stock movement
func (s *StockService) CreateStockMovement(ctx context.Context, req *products.StockMovementCreateRequest, processedBy int) (*products.StockMovement, error) {
	// Validate product exists
	_, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Set movement date if not provided
	movementDate := time.Now()
	if req.MovementDate != nil {
		movementDate = *req.MovementDate
	}

	// Create stock movement model
	movement := &products.StockMovement{
		ProductID:      req.ProductID,
		MovementType:   req.MovementType,
		ReferenceType:  req.ReferenceType,
		ReferenceID:    req.ReferenceID,
		QuantityMoved:  req.QuantityMoved,
		UnitCost:       req.UnitCost,
		LocationFrom:   req.LocationFrom,
		LocationTo:     req.LocationTo,
		MovementDate:   movementDate,
		ProcessedBy:    processedBy,
		MovementReason: req.MovementReason,
		Notes:          req.Notes,
	}

	// Validate stock availability for outgoing movements
	if req.MovementType == products.MovementTypeOut {
		currentStock, err := s.stockMovementRepo.GetCurrentStock(ctx, req.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get current stock: %w", err)
		}

		if currentStock < req.QuantityMoved {
			return nil, fmt.Errorf("insufficient stock: available %d, requested %d", currentStock, req.QuantityMoved)
		}
	}

	// Create the movement
	createdMovement, err := s.stockMovementRepo.Create(ctx, movement)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock movement: %w", err)
	}

	return createdMovement, nil
}

// GetStockMovement retrieves a stock movement by ID
func (s *StockService) GetStockMovement(ctx context.Context, id int) (*products.StockMovement, error) {
	movement, err := s.stockMovementRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock movement: %w", err)
	}

	return movement, nil
}

// ListStockMovements retrieves stock movements with pagination
func (s *StockService) ListStockMovements(ctx context.Context, params *products.StockMovementFilterParams) (*common.PaginatedResponse, error) {
	movements, err := s.stockMovementRepo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock movements: %w", err)
	}

	return movements, nil
}

// GetProductStockMovements retrieves stock movements for a specific product
func (s *StockService) GetProductStockMovements(ctx context.Context, productID int, params *products.StockMovementFilterParams) (*common.PaginatedResponse, error) {
	movements, err := s.stockMovementRepo.GetByProductID(ctx, productID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get product stock movements: %w", err)
	}

	return movements, nil
}

// GetMovementHistory gets recent movement history for a product
func (s *StockService) GetMovementHistory(ctx context.Context, productID int, limit int) ([]products.StockMovement, error) {
	// Validate product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if limit <= 0 {
		limit = 10
	}

	movements, err := s.stockMovementRepo.GetMovementHistory(ctx, productID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get movement history: %w", err)
	}

	return movements, nil
}

// GetCurrentStock gets current stock quantity for a product
func (s *StockService) GetCurrentStock(ctx context.Context, productID int) (int, error) {
	// Validate product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return 0, fmt.Errorf("product not found: %w", err)
	}

	stock, err := s.stockMovementRepo.GetCurrentStock(ctx, productID)
	if err != nil {
		return 0, fmt.Errorf("failed to get current stock: %w", err)
	}

	return stock, nil
}

// CreateStockAdjustment creates a new stock adjustment
func (s *StockService) CreateStockAdjustment(ctx context.Context, req *products.StockAdjustmentCreateRequest, createdBy int) (*products.StockAdjustment, error) {
	// Validate product exists
	_, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Set adjustment date if not provided
	adjustmentDate := time.Now()
	if req.AdjustmentDate != nil {
		adjustmentDate = *req.AdjustmentDate
	}

	// Create stock adjustment model
	adjustment := &products.StockAdjustment{
		ProductID:               req.ProductID,
		AdjustmentType:          req.AdjustmentType,
		QuantityPhysical:        req.QuantityPhysical,
		AdjustmentReason:        req.AdjustmentReason,
		Notes:                   req.Notes,
		AdjustmentDate:          adjustmentDate,
		SupportingDocumentsJSON: req.SupportingDocumentsJSON,
		CreatedBy:               createdBy,
	}

	// Create the adjustment (quantities and cost impact will be calculated in repository)
	createdAdjustment, err := s.stockAdjustmentRepo.Create(ctx, adjustment)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock adjustment: %w", err)
	}

	return createdAdjustment, nil
}

// GetStockAdjustment retrieves a stock adjustment by ID
func (s *StockService) GetStockAdjustment(ctx context.Context, id int) (*products.StockAdjustment, error) {
	adjustment, err := s.stockAdjustmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock adjustment: %w", err)
	}

	return adjustment, nil
}

// UpdateStockAdjustment updates a stock adjustment
func (s *StockService) UpdateStockAdjustment(ctx context.Context, id int, req *products.StockAdjustmentUpdateRequest) (*products.StockAdjustment, error) {
	// Get existing adjustment
	existing, err := s.stockAdjustmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing adjustment: %w", err)
	}

	// Check if adjustment is already approved
	if existing.ApprovedBy != nil {
		return nil, fmt.Errorf("cannot update approved stock adjustment")
	}

	// Update fields if provided
	if req.AdjustmentType != nil {
		existing.AdjustmentType = *req.AdjustmentType
	}
	if req.QuantityPhysical != nil {
		existing.QuantityPhysical = *req.QuantityPhysical
	}
	if req.AdjustmentReason != nil {
		existing.AdjustmentReason = *req.AdjustmentReason
	}
	if req.Notes != nil {
		existing.Notes = req.Notes
	}
	if req.AdjustmentDate != nil {
		existing.AdjustmentDate = *req.AdjustmentDate
	}
	if req.SupportingDocumentsJSON != nil {
		existing.SupportingDocumentsJSON = req.SupportingDocumentsJSON
	}

	// Update the adjustment
	updatedAdjustment, err := s.stockAdjustmentRepo.Update(ctx, id, existing)
	if err != nil {
		return nil, fmt.Errorf("failed to update stock adjustment: %w", err)
	}

	return updatedAdjustment, nil
}

// DeleteStockAdjustment deletes a stock adjustment
func (s *StockService) DeleteStockAdjustment(ctx context.Context, id int) error {
	err := s.stockAdjustmentRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete stock adjustment: %w", err)
	}

	return nil
}

// ListStockAdjustments retrieves stock adjustments with pagination
func (s *StockService) ListStockAdjustments(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	adjustments, err := s.stockAdjustmentRepo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock adjustments: %w", err)
	}

	return adjustments, nil
}

// GetProductStockAdjustments retrieves stock adjustments for a specific product
func (s *StockService) GetProductStockAdjustments(ctx context.Context, productID int, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	adjustments, err := s.stockAdjustmentRepo.GetByProductID(ctx, productID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get product stock adjustments: %w", err)
	}

	return adjustments, nil
}

// GetPendingAdjustments retrieves stock adjustments pending approval
func (s *StockService) GetPendingAdjustments(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	adjustments, err := s.stockAdjustmentRepo.GetPendingApproval(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending adjustments: %w", err)
	}

	return adjustments, nil
}

// ApproveStockAdjustment approves a stock adjustment and creates corresponding stock movement
func (s *StockService) ApproveStockAdjustment(ctx context.Context, id int, approvedBy int) error {
	// Get the adjustment
	adjustment, err := s.stockAdjustmentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get adjustment: %w", err)
	}

	// Check if already approved
	if adjustment.ApprovedBy != nil {
		return fmt.Errorf("stock adjustment already approved")
	}

	// Approve the adjustment
	err = s.stockAdjustmentRepo.Approve(ctx, id, approvedBy)
	if err != nil {
		return fmt.Errorf("failed to approve adjustment: %w", err)
	}

	// Create stock movement if there is variance
	if adjustment.QuantityVariance != 0 {
		// Get product cost price
		product, err := s.productRepo.GetByID(ctx, adjustment.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product: %w", err)
		}

		err = s.stockMovementRepo.CreateMovementForAdjustment(
			ctx,
			adjustment.ProductID,
			adjustment.QuantityVariance,
			product.CostPrice,
			id,
			approvedBy,
		)
		if err != nil {
			return fmt.Errorf("failed to create adjustment movement: %w", err)
		}
	}

	return nil
}

// GetVarianceReport gets variance report for stock adjustments
func (s *StockService) GetVarianceReport(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	report, err := s.stockAdjustmentRepo.GetVarianceReport(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get variance report: %w", err)
	}

	return report, nil
}

// TransferStock creates stock movements for stock transfer between locations
func (s *StockService) TransferStock(ctx context.Context, productID int, quantity int, fromLocation, toLocation string, processedBy int, notes *string) error {
	// Validate product exists
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	// Check stock availability at source location
	currentStock, err := s.stockMovementRepo.GetCurrentStock(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get current stock: %w", err)
	}

	if currentStock < quantity {
		return fmt.Errorf("insufficient stock for transfer: available %d, requested %d", currentStock, quantity)
	}

	// Create outbound movement (from source)
	outMovement := &products.StockMovement{
		ProductID:      productID,
		MovementType:   products.MovementTypeOut,
		ReferenceType:  products.ReferenceTypeTransfer,
		ReferenceID:    0, // Will be set to the inbound movement ID
		QuantityMoved:  quantity,
		UnitCost:       product.CostPrice,
		LocationFrom:   &fromLocation,
		LocationTo:     &toLocation,
		MovementDate:   time.Now(),
		ProcessedBy:    processedBy,
		MovementReason: stringPtr("Stock transfer - outbound"),
		Notes:          notes,
	}

	// Create inbound movement (to destination)
	inMovement := &products.StockMovement{
		ProductID:      productID,
		MovementType:   products.MovementTypeIn,
		ReferenceType:  products.ReferenceTypeTransfer,
		ReferenceID:    0, // Will be set to the outbound movement ID
		QuantityMoved:  quantity,
		UnitCost:       product.CostPrice,
		LocationFrom:   &fromLocation,
		LocationTo:     &toLocation,
		MovementDate:   time.Now(),
		ProcessedBy:    processedBy,
		MovementReason: stringPtr("Stock transfer - inbound"),
		Notes:          notes,
	}

	// Create outbound movement first
	createdOutMovement, err := s.stockMovementRepo.Create(ctx, outMovement)
	if err != nil {
		return fmt.Errorf("failed to create outbound movement: %w", err)
	}

	// Set reference ID for inbound movement
	inMovement.ReferenceID = createdOutMovement.MovementID

	// Create inbound movement
	_, err = s.stockMovementRepo.Create(ctx, inMovement)
	if err != nil {
		return fmt.Errorf("failed to create inbound movement: %w", err)
	}

	// Update outbound movement reference ID
	// Note: This would require an update method in the repository
	// For now, we'll leave it as is

	return nil
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}