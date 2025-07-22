package implementations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

type stockAdjustmentRepository struct {
	db *sql.DB
}

// NewStockAdjustmentRepository creates a new stock adjustment repository
func NewStockAdjustmentRepository(db *sql.DB) interfaces.StockAdjustmentRepository {
	return &stockAdjustmentRepository{
		db: db,
	}
}

// Create creates a new stock adjustment
func (r *stockAdjustmentRepository) Create(ctx context.Context, adjustment *inventory.StockAdjustment) (*inventory.StockAdjustment, error) {
	query := `
		INSERT INTO stock_adjustments (
			product_id, adjustment_type, quantity_system, quantity_physical,
			quantity_variance, cost_impact, adjustment_reason, notes,
			approved_by, adjustment_date, approved_at, supporting_documents_json,
			created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING adjustment_id`

	err := r.db.QueryRowContext(ctx, query,
		adjustment.ProductID,
		adjustment.AdjustmentType,
		adjustment.QuantitySystem,
		adjustment.QuantityPhysical,
		adjustment.QuantityVariance,
		adjustment.CostImpact,
		adjustment.AdjustmentReason,
		adjustment.Notes,
		adjustment.ApprovedBy,
		adjustment.AdjustmentDate,
		adjustment.ApprovedAt,
		adjustment.SupportingDocumentsJSON,
		adjustment.CreatedBy,
	).Scan(&adjustment.AdjustmentID)

	if err != nil {
		return nil, fmt.Errorf("failed to create stock adjustment: %w", err)
	}

	return adjustment, nil
}

// GetByID retrieves a stock adjustment by ID
func (r *stockAdjustmentRepository) GetByID(ctx context.Context, id int) (*inventory.StockAdjustment, error) {
	query := `
		SELECT 
			sa.adjustment_id, sa.product_id, sa.adjustment_type,
			sa.quantity_system, sa.quantity_physical, sa.quantity_variance,
			sa.cost_impact, sa.adjustment_reason, sa.notes,
			sa.approved_by, sa.adjustment_date, sa.approved_at,
			sa.supporting_documents_json, sa.created_at, sa.created_by
		FROM stock_adjustments sa
		WHERE sa.adjustment_id = $1`

	adjustment := &inventory.StockAdjustment{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&adjustment.AdjustmentID,
		&adjustment.ProductID,
		&adjustment.AdjustmentType,
		&adjustment.QuantitySystem,
		&adjustment.QuantityPhysical,
		&adjustment.QuantityVariance,
		&adjustment.CostImpact,
		&adjustment.AdjustmentReason,
		&adjustment.Notes,
		&adjustment.ApprovedBy,
		&adjustment.AdjustmentDate,
		&adjustment.ApprovedAt,
		&adjustment.SupportingDocumentsJSON,
		&adjustment.CreatedAt,
		&adjustment.CreatedBy,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get stock adjustment: %w", err)
	}

	return adjustment, nil
}

// Update updates a stock adjustment
func (r *stockAdjustmentRepository) Update(ctx context.Context, id int, adjustment *inventory.StockAdjustment) (*inventory.StockAdjustment, error) {
	query := `
		UPDATE stock_adjustments 
		SET adjustment_type = $2, quantity_system = $3, quantity_physical = $4,
			quantity_variance = $5, cost_impact = $6, adjustment_reason = $7,
			notes = $8, adjustment_date = $9, supporting_documents_json = $10
		WHERE adjustment_id = $1 AND approved_by IS NULL`

	_, err := r.db.ExecContext(ctx, query,
		id,
		adjustment.AdjustmentType,
		adjustment.QuantitySystem,
		adjustment.QuantityPhysical,
		adjustment.QuantityVariance,
		adjustment.CostImpact,
		adjustment.AdjustmentReason,
		adjustment.Notes,
		adjustment.AdjustmentDate,
		adjustment.SupportingDocumentsJSON,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update stock adjustment: %w", err)
	}

	adjustment.AdjustmentID = id
	return adjustment, nil
}

// Delete soft deletes a stock adjustment
func (r *stockAdjustmentRepository) Delete(ctx context.Context, id int) error {
	// Check if adjustment is approved
	var approvedBy sql.NullInt64
	checkQuery := `SELECT approved_by FROM stock_adjustments WHERE adjustment_id = $1`
	err := r.db.QueryRowContext(ctx, checkQuery, id).Scan(&approvedBy)
	if err != nil {
		return fmt.Errorf("failed to check adjustment status: %w", err)
	}

	if approvedBy.Valid {
		return fmt.Errorf("cannot delete approved stock adjustment")
	}

	query := `DELETE FROM stock_adjustments WHERE adjustment_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete stock adjustment: %w", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("stock adjustment with ID %d not found", id)
	}

	return nil
}

// List retrieves stock adjustments with filtering and pagination
func (r *stockAdjustmentRepository) List(ctx context.Context, params *inventory.StockAdjustmentFilterParams) ([]inventory.StockAdjustmentListItem, int, error) {
	query := `
		SELECT 
			sa.adjustment_id, '' as product_code, '' as product_name, sa.adjustment_type,
			sa.quantity_variance, sa.cost_impact, sa.adjustment_reason,
			sa.adjustment_date, '' as created_by_name, NULL as approved_by_name, sa.approved_at
		FROM stock_adjustments sa
		ORDER BY sa.adjustment_date DESC
		LIMIT $1 OFFSET $2`

	offset := (params.Page - 1) * params.Limit

	rows, err := r.db.QueryContext(ctx, query, params.Limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get stock adjustments: %w", err)
	}
	defer rows.Close()

	var items []inventory.StockAdjustmentListItem
	for rows.Next() {
		item := inventory.StockAdjustmentListItem{}

		err := rows.Scan(
			&item.AdjustmentID,
			&item.ProductCode, // This will be empty in simple version
			&item.ProductName, // This will be empty in simple version
			&item.AdjustmentType,
			&item.QuantityVariance,
			&item.CostImpact,
			&item.AdjustmentReason,
			&item.AdjustmentDate,
			&item.CreatedByName, // This will be empty in simple version
			&item.ApprovedByName,
			&item.ApprovedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan stock adjustment list item: %w", err)
		}

		items = append(items, item)
	}

	// Count total
	countQuery := `SELECT COUNT(*) FROM stock_adjustments`
	var total int
	err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count stock adjustments: %w", err)
	}

	return items, total, nil
}

// GetByProduct retrieves stock adjustments for a specific product
func (r *stockAdjustmentRepository) GetByProduct(ctx context.Context, productID int, page, limit int) ([]inventory.StockAdjustmentListItem, int, error) {
	params := &inventory.StockAdjustmentFilterParams{
		ProductID: &productID,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// GetByType retrieves stock adjustments by adjustment type
func (r *stockAdjustmentRepository) GetByType(ctx context.Context, adjustmentType inventory.AdjustmentType, page, limit int) ([]inventory.StockAdjustmentListItem, int, error) {
	params := &inventory.StockAdjustmentFilterParams{
		AdjustmentType: &adjustmentType,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// GetPendingApproval retrieves stock adjustments pending approval
func (r *stockAdjustmentRepository) GetPendingApproval(ctx context.Context, page, limit int) ([]inventory.StockAdjustmentListItem, int, error) {
	offset := (page - 1) * limit

	query := `
		SELECT 
			sa.adjustment_id, '' as product_code, '' as product_name, sa.adjustment_type,
			sa.quantity_variance, sa.cost_impact, sa.adjustment_reason,
			sa.adjustment_date, '' as created_by_name, NULL as approved_by_name, sa.approved_at
		FROM stock_adjustments sa
		WHERE sa.approved_by IS NULL
		ORDER BY sa.adjustment_date DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get pending adjustments: %w", err)
	}
	defer rows.Close()

	var items []inventory.StockAdjustmentListItem
	for rows.Next() {
		item := inventory.StockAdjustmentListItem{}

		err := rows.Scan(
			&item.AdjustmentID,
			&item.ProductCode,
			&item.ProductName,
			&item.AdjustmentType,
			&item.QuantityVariance,
			&item.CostImpact,
			&item.AdjustmentReason,
			&item.AdjustmentDate,
			&item.CreatedByName,
			&item.ApprovedByName,
			&item.ApprovedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pending adjustment: %w", err)
		}

		items = append(items, item)
	}

	// Count total
	countQuery := `SELECT COUNT(*) FROM stock_adjustments WHERE approved_by IS NULL`
	var total int
	err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pending adjustments: %w", err)
	}

	return items, total, nil
}

// UpdateApproval approves a stock adjustment
func (r *stockAdjustmentRepository) UpdateApproval(ctx context.Context, id int, approvedBy int) error {
	query := `
		UPDATE stock_adjustments 
		SET approved_by = $2, approved_at = CURRENT_TIMESTAMP
		WHERE adjustment_id = $1 AND approved_by IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, approvedBy)
	if err != nil {
		return fmt.Errorf("failed to approve stock adjustment: %w", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("stock adjustment with ID %d not found or already approved", id)
	}

	return nil
}

// GetAdjustmentSummary retrieves adjustment summary for reporting
func (r *stockAdjustmentRepository) GetAdjustmentSummary(ctx context.Context, startDate, endDate string) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_adjustments,
			COUNT(CASE WHEN approved_by IS NOT NULL THEN 1 END) as approved_adjustments,
			COUNT(CASE WHEN approved_by IS NULL THEN 1 END) as pending_adjustments,
			SUM(CASE WHEN quantity_variance > 0 THEN quantity_variance ELSE 0 END) as positive_variance,
			SUM(CASE WHEN quantity_variance < 0 THEN ABS(quantity_variance) ELSE 0 END) as negative_variance,
			SUM(CASE WHEN cost_impact > 0 THEN cost_impact ELSE 0 END) as positive_cost_impact,
			SUM(CASE WHEN cost_impact < 0 THEN ABS(cost_impact) ELSE 0 END) as negative_cost_impact
		FROM stock_adjustments sa
		WHERE sa.adjustment_date >= $1::date 
		AND sa.adjustment_date <= $2::date`

	var (
		totalAdjustments, approvedAdjustments, pendingAdjustments                 int
		positiveVariance, negativeVariance                                        int
		positiveCostImpact, negativeCostImpact                                    float64
	)

	err := r.db.QueryRowContext(ctx, query, startDate, endDate).Scan(
		&totalAdjustments, &approvedAdjustments, &pendingAdjustments,
		&positiveVariance, &negativeVariance,
		&positiveCostImpact, &negativeCostImpact,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get adjustment summary: %w", err)
	}

	summary := map[string]interface{}{
		"total_adjustments":    totalAdjustments,
		"approved_adjustments": approvedAdjustments,
		"pending_adjustments":  pendingAdjustments,
		"variance_summary": map[string]interface{}{
			"positive_variance": positiveVariance,
			"negative_variance": negativeVariance,
			"net_variance":      positiveVariance - negativeVariance,
		},
		"cost_impact_summary": map[string]interface{}{
			"positive_cost_impact": positiveCostImpact,
			"negative_cost_impact": negativeCostImpact,
			"net_cost_impact":      positiveCostImpact - negativeCostImpact,
		},
	}

	return summary, nil
}