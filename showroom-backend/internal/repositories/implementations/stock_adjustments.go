package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// StockAdjustmentRepository implements interfaces.StockAdjustmentRepository
type StockAdjustmentRepository struct {
	db *sql.DB
}

// NewStockAdjustmentRepository creates a new stock adjustment repository
func NewStockAdjustmentRepository(db *sql.DB) interfaces.StockAdjustmentRepository {
	return &StockAdjustmentRepository{db: db}
}

// Create creates a new stock adjustment
func (r *StockAdjustmentRepository) Create(ctx context.Context, adjustment *products.StockAdjustment) (*products.StockAdjustment, error) {
	// Get current system quantity
	var systemQuantity int
	err := r.db.QueryRowContext(ctx, `SELECT COALESCE(stock_quantity, 0) FROM products_spare_parts WHERE product_id = $1`, adjustment.ProductID).Scan(&systemQuantity)
	if err != nil {
		return nil, fmt.Errorf("failed to get system quantity: %w", err)
	}

	// Calculate variance and cost impact
	adjustment.QuantitySystem = systemQuantity
	adjustment.QuantityVariance = adjustment.QuantityPhysical - adjustment.QuantitySystem

	// Get cost price for cost impact calculation
	var costPrice float64
	err = r.db.QueryRowContext(ctx, `SELECT COALESCE(cost_price, 0) FROM products_spare_parts WHERE product_id = $1`, adjustment.ProductID).Scan(&costPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to get cost price: %w", err)
	}

	adjustment.CostImpact = float64(adjustment.QuantityVariance) * costPrice

	// Set adjustment date if not provided
	if adjustment.AdjustmentDate.IsZero() {
		adjustment.AdjustmentDate = time.Now()
	}

	query := `
		INSERT INTO stock_adjustments (
			product_id, adjustment_type, quantity_system, quantity_physical,
			quantity_variance, cost_impact, adjustment_reason, notes,
			adjustment_date, supporting_documents_json, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING adjustment_id, created_at`

	err = r.db.QueryRowContext(ctx, query,
		adjustment.ProductID,
		adjustment.AdjustmentType,
		adjustment.QuantitySystem,
		adjustment.QuantityPhysical,
		adjustment.QuantityVariance,
		adjustment.CostImpact,
		adjustment.AdjustmentReason,
		adjustment.Notes,
		adjustment.AdjustmentDate,
		adjustment.SupportingDocumentsJSON,
		adjustment.CreatedBy,
	).Scan(&adjustment.AdjustmentID, &adjustment.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create stock adjustment: %w", err)
	}

	return adjustment, nil
}

// GetByID retrieves a stock adjustment by ID
func (r *StockAdjustmentRepository) GetByID(ctx context.Context, id int) (*products.StockAdjustment, error) {
	query := `
		SELECT adjustment_id, product_id, adjustment_type, quantity_system,
			   quantity_physical, quantity_variance, cost_impact, adjustment_reason,
			   notes, approved_by, adjustment_date, approved_at,
			   supporting_documents_json, created_at, created_by
		FROM stock_adjustments 
		WHERE adjustment_id = $1`

	adjustment := &products.StockAdjustment{}
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("stock adjustment not found")
		}
		return nil, fmt.Errorf("failed to get stock adjustment: %w", err)
	}

	return adjustment, nil
}

// Update updates a stock adjustment
func (r *StockAdjustmentRepository) Update(ctx context.Context, id int, adjustment *products.StockAdjustment) (*products.StockAdjustment, error) {
	// Recalculate variance and cost impact if quantities changed
	adjustment.QuantityVariance = adjustment.QuantityPhysical - adjustment.QuantitySystem

	// Get cost price for cost impact calculation
	var costPrice float64
	err := r.db.QueryRowContext(ctx, `SELECT COALESCE(cost_price, 0) FROM products_spare_parts WHERE product_id = $1`, adjustment.ProductID).Scan(&costPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to get cost price: %w", err)
	}

	adjustment.CostImpact = float64(adjustment.QuantityVariance) * costPrice

	query := `
		UPDATE stock_adjustments 
		SET adjustment_type = $1, quantity_physical = $2, quantity_variance = $3,
			cost_impact = $4, adjustment_reason = $5, notes = $6,
			adjustment_date = $7, supporting_documents_json = $8
		WHERE adjustment_id = $9`

	_, err = r.db.ExecContext(ctx, query,
		adjustment.AdjustmentType,
		adjustment.QuantityPhysical,
		adjustment.QuantityVariance,
		adjustment.CostImpact,
		adjustment.AdjustmentReason,
		adjustment.Notes,
		adjustment.AdjustmentDate,
		adjustment.SupportingDocumentsJSON,
		id,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update stock adjustment: %w", err)
	}

	return r.GetByID(ctx, id)
}

// Delete deletes a stock adjustment
func (r *StockAdjustmentRepository) Delete(ctx context.Context, id int) error {
	// Check if adjustment is approved
	var approvedBy *int
	err := r.db.QueryRowContext(ctx, `SELECT approved_by FROM stock_adjustments WHERE adjustment_id = $1`, id).Scan(&approvedBy)
	if err != nil {
		return fmt.Errorf("failed to check adjustment status: %w", err)
	}

	if approvedBy != nil {
		return fmt.Errorf("cannot delete approved stock adjustment")
	}

	query := `DELETE FROM stock_adjustments WHERE adjustment_id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete stock adjustment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("stock adjustment not found")
	}

	return nil
}

// List retrieves all stock adjustments with pagination
func (r *StockAdjustmentRepository) List(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	baseQuery := `
		FROM stock_adjustments sa 
		LEFT JOIN products_spare_parts psp ON sa.product_id = psp.product_id
		WHERE 1=1`
	
	args := []interface{}{}
	whereConditions := []string{}
	argIndex := 1

	// Add filters
	if params.ProductID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sa.product_id = $%d", argIndex))
		args = append(args, *params.ProductID)
		argIndex++
	}

	if params.AdjustmentType != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sa.adjustment_type = $%d", argIndex))
		args = append(args, *params.AdjustmentType)
		argIndex++
	}

	if params.CreatedBy != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sa.created_by = $%d", argIndex))
		args = append(args, *params.CreatedBy)
		argIndex++
	}

	if params.ApprovedBy != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sa.approved_by = $%d", argIndex))
		args = append(args, *params.ApprovedBy)
		argIndex++
	}

	if params.DateFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sa.adjustment_date >= $%d", argIndex))
		args = append(args, *params.DateFrom)
		argIndex++
	}

	if params.DateTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sa.adjustment_date <= $%d", argIndex))
		args = append(args, *params.DateTo)
		argIndex++
	}

	if params.HasVariance != nil && *params.HasVariance {
		whereConditions = append(whereConditions, "sa.quantity_variance != 0")
	}

	if params.Search != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("(psp.product_name ILIKE $%d OR sa.adjustment_reason ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	if len(whereConditions) > 0 {
		baseQuery += " AND " + strings.Join(whereConditions, " AND ")
	}

	// Count query
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count stock adjustments: %w", err)
	}

	// Calculate pagination
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	offset := (params.Page - 1) * params.Limit

	// Main query
	selectFields := `
		sa.adjustment_id, sa.product_id, sa.adjustment_type, sa.quantity_system,
		sa.quantity_physical, sa.quantity_variance, sa.cost_impact,
		sa.adjustment_reason, sa.adjustment_date, sa.approved_by, sa.created_by`
	
	mainQuery := "SELECT " + selectFields + " " + baseQuery + 
		" ORDER BY sa.adjustment_date DESC, sa.adjustment_id DESC LIMIT $" + fmt.Sprintf("%d", argIndex) + 
		" OFFSET $" + fmt.Sprintf("%d", argIndex+1)
	
	args = append(args, params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, mainQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query stock adjustments: %w", err)
	}
	defer rows.Close()

	var adjustments []products.StockAdjustmentListItem
	for rows.Next() {
		var adjustment products.StockAdjustmentListItem
		err := rows.Scan(
			&adjustment.AdjustmentID,
			&adjustment.ProductID,
			&adjustment.AdjustmentType,
			&adjustment.QuantitySystem,
			&adjustment.QuantityPhysical,
			&adjustment.QuantityVariance,
			&adjustment.CostImpact,
			&adjustment.AdjustmentReason,
			&adjustment.AdjustmentDate,
			&adjustment.ApprovedBy,
			&adjustment.CreatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock adjustment: %w", err)
		}
		adjustments = append(adjustments, adjustment)
	}

	totalPages := (total + int64(params.Limit) - 1) / int64(params.Limit)

	return &common.PaginatedResponse{
		Data:       adjustments,
		Total:      int(total),
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: int(totalPages),
		HasMore:    params.Page < int(totalPages),
	}, nil
}

// GetByProductID retrieves stock adjustments for a specific product
func (r *StockAdjustmentRepository) GetByProductID(ctx context.Context, productID int, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	// Set productID in params and call List
	params.ProductID = &productID
	return r.List(ctx, params)
}

// GetPendingApproval retrieves stock adjustments pending approval
func (r *StockAdjustmentRepository) GetPendingApproval(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	baseQuery := `
		FROM stock_adjustments sa 
		LEFT JOIN products_spare_parts psp ON sa.product_id = psp.product_id
		WHERE sa.approved_by IS NULL`
	
	args := []interface{}{}
	whereConditions := []string{}
	argIndex := 1

	// Add additional filters
	if params.ProductID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sa.product_id = $%d", argIndex))
		args = append(args, *params.ProductID)
		argIndex++
	}

	if params.AdjustmentType != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sa.adjustment_type = $%d", argIndex))
		args = append(args, *params.AdjustmentType)
		argIndex++
	}

	if params.CreatedBy != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sa.created_by = $%d", argIndex))
		args = append(args, *params.CreatedBy)
		argIndex++
	}

	if len(whereConditions) > 0 {
		baseQuery += " AND " + strings.Join(whereConditions, " AND ")
	}

	// Count query
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count pending adjustments: %w", err)
	}

	// Calculate pagination
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	offset := (params.Page - 1) * params.Limit

	// Main query
	selectFields := `
		sa.adjustment_id, sa.product_id, sa.adjustment_type, sa.quantity_system,
		sa.quantity_physical, sa.quantity_variance, sa.cost_impact,
		sa.adjustment_reason, sa.adjustment_date, sa.approved_by, sa.created_by`
	
	mainQuery := "SELECT " + selectFields + " " + baseQuery + 
		" ORDER BY sa.created_at ASC LIMIT $" + fmt.Sprintf("%d", argIndex) + 
		" OFFSET $" + fmt.Sprintf("%d", argIndex+1)
	
	args = append(args, params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, mainQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending adjustments: %w", err)
	}
	defer rows.Close()

	var adjustments []products.StockAdjustmentListItem
	for rows.Next() {
		var adjustment products.StockAdjustmentListItem
		err := rows.Scan(
			&adjustment.AdjustmentID,
			&adjustment.ProductID,
			&adjustment.AdjustmentType,
			&adjustment.QuantitySystem,
			&adjustment.QuantityPhysical,
			&adjustment.QuantityVariance,
			&adjustment.CostImpact,
			&adjustment.AdjustmentReason,
			&adjustment.AdjustmentDate,
			&adjustment.ApprovedBy,
			&adjustment.CreatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock adjustment: %w", err)
		}
		adjustments = append(adjustments, adjustment)
	}

	totalPages := (total + int64(params.Limit) - 1) / int64(params.Limit)

	return &common.PaginatedResponse{
		Data:       adjustments,
		Total:      int(total),
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: int(totalPages),
		HasMore:    params.Page < int(totalPages),
	}, nil
}

// Approve approves a stock adjustment and creates stock movement
func (r *StockAdjustmentRepository) Approve(ctx context.Context, id int, approvedBy int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update approval status
	now := time.Now()
	query := `UPDATE stock_adjustments SET approved_by = $1, approved_at = $2 WHERE adjustment_id = $3`
	_, err = tx.ExecContext(ctx, query, approvedBy, now, id)
	if err != nil {
		return fmt.Errorf("failed to approve adjustment: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetVarianceReport gets variance report
func (r *StockAdjustmentRepository) GetVarianceReport(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	baseQuery := `
		FROM stock_adjustments sa 
		LEFT JOIN products_spare_parts psp ON sa.product_id = psp.product_id
		WHERE sa.quantity_variance != 0`
	
	args := []interface{}{}
	whereConditions := []string{}
	argIndex := 1

	// Add filters
	if params.ProductID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sa.product_id = $%d", argIndex))
		args = append(args, *params.ProductID)
		argIndex++
	}

	if params.AdjustmentType != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sa.adjustment_type = $%d", argIndex))
		args = append(args, *params.AdjustmentType)
		argIndex++
	}

	if params.DateFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sa.adjustment_date >= $%d", argIndex))
		args = append(args, *params.DateFrom)
		argIndex++
	}

	if params.DateTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sa.adjustment_date <= $%d", argIndex))
		args = append(args, *params.DateTo)
		argIndex++
	}

	if len(whereConditions) > 0 {
		baseQuery += " AND " + strings.Join(whereConditions, " AND ")
	}

	// Count query
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count variance adjustments: %w", err)
	}

	// Calculate pagination
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	offset := (params.Page - 1) * params.Limit

	// Main query
	selectFields := `
		sa.adjustment_id, sa.product_id, sa.adjustment_type, sa.quantity_system,
		sa.quantity_physical, sa.quantity_variance, sa.cost_impact,
		sa.adjustment_reason, sa.adjustment_date, sa.approved_by, sa.created_by`
	
	mainQuery := "SELECT " + selectFields + " " + baseQuery + 
		" ORDER BY ABS(sa.cost_impact) DESC, sa.adjustment_date DESC LIMIT $" + fmt.Sprintf("%d", argIndex) + 
		" OFFSET $" + fmt.Sprintf("%d", argIndex+1)
	
	args = append(args, params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, mainQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query variance adjustments: %w", err)
	}
	defer rows.Close()

	var adjustments []products.StockAdjustmentListItem
	for rows.Next() {
		var adjustment products.StockAdjustmentListItem
		err := rows.Scan(
			&adjustment.AdjustmentID,
			&adjustment.ProductID,
			&adjustment.AdjustmentType,
			&adjustment.QuantitySystem,
			&adjustment.QuantityPhysical,
			&adjustment.QuantityVariance,
			&adjustment.CostImpact,
			&adjustment.AdjustmentReason,
			&adjustment.AdjustmentDate,
			&adjustment.ApprovedBy,
			&adjustment.CreatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock adjustment: %w", err)
		}
		adjustments = append(adjustments, adjustment)
	}

	totalPages := (total + int64(params.Limit) - 1) / int64(params.Limit)

	return &common.PaginatedResponse{
		Data:       adjustments,
		Total:      int(total),
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: int(totalPages),
		HasMore:    params.Page < int(totalPages),
	}, nil
}