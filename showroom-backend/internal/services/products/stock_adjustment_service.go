package products

import (
	"context"
	"fmt"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
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
func (s *StockAdjustmentService) CreateStockAdjustment(ctx context.Context, req *products.StockAdjustmentCreateRequest, createdBy int) (*products.StockAdjustment, error) {
	// Validate product exists
	_, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Create stock adjustment from request
	adjustment := &products.StockAdjustment{
		ProductID:               req.ProductID,
		AdjustmentType:          req.AdjustmentType,
		QuantityPhysical:        req.QuantityPhysical,
		AdjustmentReason:        req.AdjustmentReason,
		Notes:                   req.Notes,
		SupportingDocumentsJSON: req.SupportingDocumentsJSON,
		CreatedBy:               createdBy,
	}

	// Set adjustment date if provided
	if req.AdjustmentDate != nil {
		adjustment.AdjustmentDate = *req.AdjustmentDate
	}

	// Create the adjustment (repository will calculate quantities and cost impact)
	createdAdjustment, err := s.stockAdjustmentRepo.Create(ctx, adjustment)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock adjustment: %w", err)
	}

	return createdAdjustment, nil
}

// GetStockAdjustment retrieves a stock adjustment by ID
func (s *StockAdjustmentService) GetStockAdjustment(ctx context.Context, id int) (*products.StockAdjustment, error) {
	adjustment, err := s.stockAdjustmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock adjustment: %w", err)
	}

	return adjustment, nil
}

// UpdateStockAdjustment updates a stock adjustment
func (s *StockAdjustmentService) UpdateStockAdjustment(ctx context.Context, id int, req *products.StockAdjustmentUpdateRequest) (*products.StockAdjustment, error) {
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
func (s *StockAdjustmentService) DeleteStockAdjustment(ctx context.Context, id int) error {
	err := s.stockAdjustmentRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete stock adjustment: %w", err)
	}

	return nil
}

// ListStockAdjustments retrieves stock adjustments with pagination
func (s *StockAdjustmentService) ListStockAdjustments(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	adjustments, err := s.stockAdjustmentRepo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock adjustments: %w", err)
	}

	return adjustments, nil
}

// GetProductStockAdjustments retrieves stock adjustments for a specific product
func (s *StockAdjustmentService) GetProductStockAdjustments(ctx context.Context, productID int, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	// Validate product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	adjustments, err := s.stockAdjustmentRepo.GetByProductID(ctx, productID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get product stock adjustments: %w", err)
	}

	return adjustments, nil
}

// GetPendingAdjustments retrieves stock adjustments pending approval
func (s *StockAdjustmentService) GetPendingAdjustments(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	adjustments, err := s.stockAdjustmentRepo.GetPendingApproval(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending adjustments: %w", err)
	}

	return adjustments, nil
}

// ApproveStockAdjustment approves a stock adjustment and creates corresponding stock movement
func (s *StockAdjustmentService) ApproveStockAdjustment(ctx context.Context, id int, approvedBy int) error {
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
		// Get product for cost price
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
func (s *StockAdjustmentService) GetVarianceReport(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	report, err := s.stockAdjustmentRepo.GetVarianceReport(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get variance report: %w", err)
	}

	return report, nil
}

// CreatePhysicalCountAdjustment creates multiple adjustments based on physical count results
func (s *StockAdjustmentService) CreatePhysicalCountAdjustments(ctx context.Context, counts []PhysicalCountItem, countedBy int, reason string) ([]products.StockAdjustment, error) {
	var adjustments []products.StockAdjustment

	for _, count := range counts {
		// Validate product exists
		_, err := s.productRepo.GetByID(ctx, count.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %d not found: %w", count.ProductID, err)
		}

		// Create adjustment
		adjustment := &products.StockAdjustment{
			ProductID:        count.ProductID,
			AdjustmentType:   products.AdjustmentTypePhysicalCount,
			QuantityPhysical: count.CountedQuantity,
			AdjustmentReason: reason,
			Notes:            count.Notes,
			CreatedBy:        countedBy,
		}

		createdAdjustment, err := s.stockAdjustmentRepo.Create(ctx, adjustment)
		if err != nil {
			return nil, fmt.Errorf("failed to create adjustment for product %d: %w", count.ProductID, err)
		}

		adjustments = append(adjustments, *createdAdjustment)
	}

	return adjustments, nil
}

// BulkApproveAdjustments approves multiple stock adjustments at once
func (s *StockAdjustmentService) BulkApproveAdjustments(ctx context.Context, adjustmentIDs []int, approvedBy int) error {
	for _, id := range adjustmentIDs {
		err := s.ApproveStockAdjustment(ctx, id, approvedBy)
		if err != nil {
			return fmt.Errorf("failed to approve adjustment %d: %w", id, err)
		}
	}

	return nil
}

// ValidateAdjustmentWorkflow validates if an adjustment can be approved based on business rules
func (s *StockAdjustmentService) ValidateAdjustmentWorkflow(ctx context.Context, id int) error {
	adjustment, err := s.stockAdjustmentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get adjustment: %w", err)
	}

	// Check if already approved
	if adjustment.ApprovedBy != nil {
		return fmt.Errorf("adjustment already approved")
	}

	// Check if supporting documents are required for certain types
	if adjustment.AdjustmentType == products.AdjustmentTypeTheft || 
	   adjustment.AdjustmentType == products.AdjustmentTypeDamage {
		if adjustment.SupportingDocumentsJSON == nil || *adjustment.SupportingDocumentsJSON == "" {
			return fmt.Errorf("supporting documents required for %s adjustments", adjustment.AdjustmentType)
		}
	}

	// Check if variance is significant (business rule example)
	product, err := s.productRepo.GetByID(ctx, adjustment.ProductID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	// If cost impact is high, might require additional approval
	costImpact := float64(adjustment.QuantityVariance) * product.CostPrice
	if costImpact > 10000 { // Example threshold
		return fmt.Errorf("high value adjustment requires manager approval")
	}

	return nil
}

// PhysicalCountItem represents an item counted during physical inventory
type PhysicalCountItem struct {
	ProductID       int     `json:"product_id"`
	CountedQuantity int     `json:"counted_quantity"`
	Notes           *string `json:"notes,omitempty"`
}