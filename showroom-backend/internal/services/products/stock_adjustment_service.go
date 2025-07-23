package products

import (
	"context"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// StockAdjustmentService handles business logic for stock adjustment management
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

// CreateStockAdjustment creates a new stock adjustment with business validation
func (s *StockAdjustmentService) CreateStockAdjustment(ctx context.Context, req *products.StockAdjustmentCreateRequest, createdBy int) (*products.StockAdjustment, error) {
	// Validate product exists
	product, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Validate adjustment type
	if !req.AdjustmentType.IsValid() {
		return nil, fmt.Errorf("invalid adjustment type: %s", req.AdjustmentType)
	}

	// Get system stock quantity
	systemStock, err := s.stockMovementRepo.GetCurrentStock(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current stock: %w", err)
	}

	// Calculate variance
	quantityVariance := req.QuantityPhysical - systemStock

	// Calculate cost impact
	costImpact := float64(quantityVariance) * product.CostPrice

	// Set adjustment date
	adjustmentDate := time.Now()
	if req.AdjustmentDate != nil {
		adjustmentDate = *req.AdjustmentDate
	}

	// Create adjustment record
	adjustment := &products.StockAdjustment{
		ProductID:               req.ProductID,
		AdjustmentType:          req.AdjustmentType,
		QuantitySystem:          systemStock,
		QuantityPhysical:        req.QuantityPhysical,
		QuantityVariance:        quantityVariance,
		CostImpact:              costImpact,
		AdjustmentReason:        req.AdjustmentReason,
		Notes:                   req.Notes,
		AdjustmentDate:          adjustmentDate,
		SupportingDocumentsJSON: req.SupportingDocumentsJSON,
		CreatedBy:               createdBy,
	}

	// Validate business rules
	if err := s.validateAdjustmentBusinessRules(adjustment); err != nil {
		return nil, fmt.Errorf("business validation failed: %w", err)
	}

	createdAdjustment, err := s.stockAdjustmentRepo.Create(ctx, adjustment)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock adjustment: %w", err)
	}

	// For certain adjustment types, auto-approve if within tolerance
	if s.shouldAutoApprove(adjustment) {
		err = s.ApproveStockAdjustment(ctx, createdAdjustment.AdjustmentID, createdBy)
		if err != nil {
			return nil, fmt.Errorf("failed to auto-approve adjustment: %w", err)
		}
		// Refresh adjustment to get updated approval status
		return s.stockAdjustmentRepo.GetByID(ctx, createdAdjustment.AdjustmentID)
	}

	return createdAdjustment, nil
}

// GetStockAdjustments retrieves stock adjustments with pagination and filtering
func (s *StockAdjustmentService) GetStockAdjustments(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	return s.stockAdjustmentRepo.List(ctx, params)
}

// GetStockAdjustmentByID retrieves a stock adjustment by ID
func (s *StockAdjustmentService) GetStockAdjustmentByID(ctx context.Context, id int) (*products.StockAdjustment, error) {
	return s.stockAdjustmentRepo.GetByID(ctx, id)
}

// GetStockAdjustmentsByProduct retrieves stock adjustments for a specific product
func (s *StockAdjustmentService) GetStockAdjustmentsByProduct(ctx context.Context, productID int, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	// Validate product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	return s.stockAdjustmentRepo.GetByProductID(ctx, productID, params)
}

// GetPendingApprovalAdjustments retrieves stock adjustments pending approval
func (s *StockAdjustmentService) GetPendingApprovalAdjustments(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	return s.stockAdjustmentRepo.GetPendingApproval(ctx, params)
}

// GetVarianceReport retrieves variance report
func (s *StockAdjustmentService) GetVarianceReport(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	return s.stockAdjustmentRepo.GetVarianceReport(ctx, params)
}

// UpdateStockAdjustment updates a stock adjustment (only if not approved)
func (s *StockAdjustmentService) UpdateStockAdjustment(ctx context.Context, id int, req *products.StockAdjustmentUpdateRequest) (*products.StockAdjustment, error) {
	// Get existing adjustment
	adjustment, err := s.stockAdjustmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("stock adjustment not found: %w", err)
	}

	// Check if already approved
	if adjustment.IsApproved() {
		return nil, fmt.Errorf("cannot update approved stock adjustment")
	}

	// Get product for cost calculation
	product, err := s.productRepo.GetByID(ctx, adjustment.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Update fields if provided
	if req.AdjustmentType != nil {
		if !req.AdjustmentType.IsValid() {
			return nil, fmt.Errorf("invalid adjustment type: %s", *req.AdjustmentType)
		}
		adjustment.AdjustmentType = *req.AdjustmentType
	}
	
	if req.QuantityPhysical != nil {
		adjustment.QuantityPhysical = *req.QuantityPhysical
		// Recalculate variance and cost impact
		adjustment.CalculateVariance(product.CostPrice)
	}
	
	if req.AdjustmentReason != nil {
		adjustment.AdjustmentReason = *req.AdjustmentReason
	}
	
	if req.Notes != nil {
		adjustment.Notes = req.Notes
	}
	
	if req.AdjustmentDate != nil {
		adjustment.AdjustmentDate = *req.AdjustmentDate
	}
	
	if req.SupportingDocumentsJSON != nil {
		adjustment.SupportingDocumentsJSON = req.SupportingDocumentsJSON
	}

	// Validate business rules
	if err := s.validateAdjustmentBusinessRules(adjustment); err != nil {
		return nil, fmt.Errorf("business validation failed: %w", err)
	}

	return s.stockAdjustmentRepo.Update(ctx, id, adjustment)
}

// ApproveStockAdjustment approves a stock adjustment and creates stock movement
func (s *StockAdjustmentService) ApproveStockAdjustment(ctx context.Context, id int, approvedBy int) error {
	// Get adjustment
	adjustment, err := s.stockAdjustmentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("stock adjustment not found: %w", err)
	}

	// Check if already approved
	if adjustment.IsApproved() {
		return fmt.Errorf("stock adjustment already approved")
	}

	// Get product for unit cost
	product, err := s.productRepo.GetByID(ctx, adjustment.ProductID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	// Approve the adjustment
	err = s.stockAdjustmentRepo.Approve(ctx, id, approvedBy)
	if err != nil {
		return fmt.Errorf("failed to approve adjustment: %w", err)
	}

	// Create stock movement if there's a variance
	if adjustment.QuantityVariance != 0 {
		err = s.stockMovementRepo.CreateMovementForAdjustment(
			ctx,
			adjustment.ProductID,
			adjustment.QuantityVariance,
			product.CostPrice,
			id,
			approvedBy,
		)
		if err != nil {
			return fmt.Errorf("failed to create stock movement: %w", err)
		}

		// Update product stock quantity
		newStock := adjustment.QuantitySystem + adjustment.QuantityVariance
		err = s.productRepo.UpdateStock(ctx, adjustment.ProductID, newStock)
		if err != nil {
			return fmt.Errorf("failed to update product stock: %w", err)
		}
	}

	return nil
}

// DeleteStockAdjustment deletes a stock adjustment (only if not approved)
func (s *StockAdjustmentService) DeleteStockAdjustment(ctx context.Context, id int) error {
	// Get adjustment
	adjustment, err := s.stockAdjustmentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("stock adjustment not found: %w", err)
	}

	// Check if already approved
	if adjustment.IsApproved() {
		return fmt.Errorf("cannot delete approved stock adjustment")
	}

	return s.stockAdjustmentRepo.Delete(ctx, id)
}

// ProcessPhysicalCountAdjustment processes physical count results and creates adjustments
func (s *StockAdjustmentService) ProcessPhysicalCountAdjustment(ctx context.Context, counts []struct {
	ProductID        int    `json:"product_id"`
	QuantityPhysical int    `json:"quantity_physical"`
	Notes            string `json:"notes,omitempty"`
}, createdBy int) ([]int, error) {
	var adjustmentIDs []int

	for _, count := range counts {
		// Get current system stock
		systemStock, err := s.stockMovementRepo.GetCurrentStock(ctx, count.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get current stock for product %d: %w", count.ProductID, err)
		}

		// Only create adjustment if there's a variance
		if count.QuantityPhysical != systemStock {
			req := &products.StockAdjustmentCreateRequest{
				ProductID:        count.ProductID,
				AdjustmentType:   products.AdjustmentTypePhysicalCount,
				QuantityPhysical: count.QuantityPhysical,
				AdjustmentReason: "Physical count variance",
				Notes:            &count.Notes,
			}

			adjustment, err := s.CreateStockAdjustment(ctx, req, createdBy)
			if err != nil {
				return nil, fmt.Errorf("failed to create adjustment for product %d: %w", count.ProductID, err)
			}

			adjustmentIDs = append(adjustmentIDs, adjustment.AdjustmentID)
		}
	}

	return adjustmentIDs, nil
}

// GetAdjustmentSummary retrieves adjustment summary for analysis
func (s *StockAdjustmentService) GetAdjustmentSummary(ctx context.Context, productID *int, dateFrom, dateTo *time.Time) (map[string]interface{}, error) {
	// This would require additional repository methods for aggregate queries
	// For now, return a placeholder
	summary := map[string]interface{}{
		"message": "Adjustment summary not fully implemented",
	}

	if productID != nil {
		// Validate product exists
		_, err := s.productRepo.GetByID(ctx, *productID)
		if err != nil {
			return nil, fmt.Errorf("product not found: %w", err)
		}
		summary["product_id"] = *productID
	}

	return summary, nil
}

// validateAdjustmentBusinessRules validates business rules for stock adjustments
func (s *StockAdjustmentService) validateAdjustmentBusinessRules(adjustment *products.StockAdjustment) error {
	// Validate negative stock for certain adjustment types
	if adjustment.QuantityPhysical < 0 {
		return fmt.Errorf("physical quantity cannot be negative")
	}

	// Validate significant variances require reason
	if absInt(adjustment.QuantityVariance) > 10 && adjustment.AdjustmentReason == "" {
		return fmt.Errorf("adjustment reason required for variances greater than 10 units")
	}

	// Validate write-off adjustments always have negative variance
	if adjustment.AdjustmentType == products.AdjustmentTypeWriteOff && adjustment.QuantityVariance >= 0 {
		return fmt.Errorf("write-off adjustments must reduce stock quantity")
	}

	// Validate damage/theft adjustments have negative variance
	if (adjustment.AdjustmentType == products.AdjustmentTypeDamage || 
		adjustment.AdjustmentType == products.AdjustmentTypeTheft) && 
		adjustment.QuantityVariance >= 0 {
		return fmt.Errorf("damage/theft adjustments must reduce stock quantity")
	}

	return nil
}

// shouldAutoApprove determines if an adjustment should be auto-approved
func (s *StockAdjustmentService) shouldAutoApprove(adjustment *products.StockAdjustment) bool {
	// Auto-approve small physical count variances
	if adjustment.AdjustmentType == products.AdjustmentTypePhysicalCount &&
		absInt(adjustment.QuantityVariance) <= 2 &&
		absInt(int(adjustment.CostImpact)) <= 100 {
		return true
	}

	// Auto-approve corrections within tolerance
	if adjustment.AdjustmentType == products.AdjustmentTypeCorrection &&
		absInt(adjustment.QuantityVariance) <= 1 {
		return true
	}

	return false
}

// helper function to get absolute value for integers
func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}