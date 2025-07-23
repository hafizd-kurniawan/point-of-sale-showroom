package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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
	// Get current stock quantity from products table
	var currentStock int
	productQuery := `SELECT COALESCE(stock_quantity, 0) FROM products_spare_parts WHERE product_id = $1`
	err := r.db.QueryRowContext(ctx, productQuery, adjustment.ProductID).Scan(&currentStock)
	if err != nil {
		return nil, fmt.Errorf("failed to get current stock for product %d: %w", adjustment.ProductID, err)
	}

	// Set system quantity and calculate variance
	adjustment.QuantitySystem = currentStock

	// Get unit cost for cost impact calculation
	var unitCost float64
	costQuery := `SELECT COALESCE(cost_price, 0) FROM products_spare_parts WHERE product_id = $1`
	err = r.db.QueryRowContext(ctx, costQuery, adjustment.ProductID).Scan(&unitCost)
	if err != nil {
		return nil, fmt.Errorf("failed to get unit cost for product %d: %w", adjustment.ProductID, err)
	}

	adjustment.CalculateVariance(unitCost)

	// Set adjustment date if not provided
	if adjustment.AdjustmentDate.IsZero() {
		adjustment.AdjustmentDate = time.Now()
	}

	query := `
		INSERT INTO stock_adjustments (
			product_id, adjustment_type, quantity_system, quantity_physical,
			quantity_variance, cost_impact, adjustment_reason, notes, approved_by,
			adjustment_date, approved_at, supporting_documents_json, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
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
		adjustment.ApprovedBy,
		adjustment.AdjustmentDate,
		adjustment.ApprovedAt,
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
		SELECT adjustment_id, product_id, adjustment_type, quantity_system, quantity_physical,
			   quantity_variance, cost_impact, adjustment_reason, notes, approved_by,
			   adjustment_date, approved_at, supporting_documents_json, created_at, created_by
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
			return nil, fmt.Errorf("stock adjustment with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get stock adjustment: %w", err)
	}

	return adjustment, nil
}

// Update updates a stock adjustment
func (r *StockAdjustmentRepository) Update(ctx context.Context, id int, adjustment *products.StockAdjustment) (*products.StockAdjustment, error) {
	// Get current adjustment to preserve some fields
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Don't allow updates if already approved
	if current.IsApproved() {
		return nil, fmt.Errorf("cannot update approved stock adjustment")
	}

	// Get unit cost for recalculating cost impact
	var unitCost float64
	costQuery := `SELECT COALESCE(cost_price, 0) FROM products_spare_parts WHERE product_id = $1`
	err = r.db.QueryRowContext(ctx, costQuery, current.ProductID).Scan(&unitCost)
	if err != nil {
		return nil, fmt.Errorf("failed to get unit cost for product %d: %w", current.ProductID, err)
	}

	// Recalculate variance with new physical quantity
	adjustment.ProductID = current.ProductID
	adjustment.QuantitySystem = current.QuantitySystem
	adjustment.CalculateVariance(unitCost)

	query := `
		UPDATE stock_adjustments SET
			adjustment_type = $2, quantity_physical = $3, quantity_variance = $4,
			cost_impact = $5, adjustment_reason = $6, notes = $7,
			adjustment_date = $8, supporting_documents_json = $9
		WHERE adjustment_id = $1 AND approved_by IS NULL`

	result, err := r.db.ExecContext(ctx, query,
		id,
		adjustment.AdjustmentType,
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

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("stock adjustment with ID %d not found or already approved", id)
	}

	return r.GetByID(ctx, id)
}

// Delete deletes a stock adjustment (only if not approved)
func (r *StockAdjustmentRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM stock_adjustments WHERE adjustment_id = $1 AND approved_by IS NULL`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete stock adjustment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("stock adjustment with ID %d not found or already approved", id)
	}

	return nil
}

// List retrieves a paginated list of stock adjustments
func (r *StockAdjustmentRepository) List(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	params.Validate()

	baseQuery := `
		SELECT adjustment_id, product_id, adjustment_type, quantity_system, quantity_physical,
			   quantity_variance, cost_impact, adjustment_reason, adjustment_date,
			   approved_by, created_by
		FROM stock_adjustments`

	countQuery := `SELECT COUNT(*) FROM stock_adjustments`

	whereConditions, args := r.buildWhereConditions(params)
	if len(whereConditions) > 0 {
		whereClause := " WHERE " + strings.Join(whereConditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count stock adjustments: %w", err)
	}

	// Add ordering and pagination
	baseQuery += ` ORDER BY adjustment_date DESC, adjustment_id DESC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock adjustments: %w", err)
	}
	defer rows.Close()

	var items []products.StockAdjustmentListItem
	for rows.Next() {
		var item products.StockAdjustmentListItem
		err := rows.Scan(
			&item.AdjustmentID,
			&item.ProductID,
			&item.AdjustmentType,
			&item.QuantitySystem,
			&item.QuantityPhysical,
			&item.QuantityVariance,
			&item.CostImpact,
			&item.AdjustmentReason,
			&item.AdjustmentDate,
			&item.ApprovedBy,
			&item.CreatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock adjustment: %w", err)
		}
		items = append(items, item)
	}

	return &common.PaginatedResponse{
		Data:       items,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: params.GetTotalPages(total),
		HasMore:    params.GetHasMore(total),
	}, nil
}

// GetByProductID retrieves stock adjustments by product ID
func (r *StockAdjustmentRepository) GetByProductID(ctx context.Context, productID int, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	params.ProductID = &productID
	return r.List(ctx, params)
}

// GetPendingApproval retrieves stock adjustments pending approval
func (r *StockAdjustmentRepository) GetPendingApproval(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	params.Validate()

	baseQuery := `
		SELECT adjustment_id, product_id, adjustment_type, quantity_system, quantity_physical,
			   quantity_variance, cost_impact, adjustment_reason, adjustment_date,
			   approved_by, created_by
		FROM stock_adjustments
		WHERE approved_by IS NULL`

	countQuery := `SELECT COUNT(*) FROM stock_adjustments WHERE approved_by IS NULL`

	whereConditions, args := r.buildWhereConditions(params)
	// Remove approval condition since we already have it
	filteredConditions := []string{}
	filteredArgs := []interface{}{}
	argIndex := 0
	
	for _, condition := range whereConditions {
		if !strings.Contains(condition, "approved_by") {
			filteredConditions = append(filteredConditions, condition)
			filteredArgs = append(filteredArgs, args[argIndex])
		}
		argIndex++
	}
	
	if len(filteredConditions) > 0 {
		whereClause := " AND " + strings.Join(filteredConditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
		args = filteredArgs
	} else {
		args = []interface{}{}
	}

	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count pending approval adjustments: %w", err)
	}

	// Add ordering and pagination
	baseQuery += ` ORDER BY adjustment_date ASC, adjustment_id ASC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending approval adjustments: %w", err)
	}
	defer rows.Close()

	var items []products.StockAdjustmentListItem
	for rows.Next() {
		var item products.StockAdjustmentListItem
		err := rows.Scan(
			&item.AdjustmentID,
			&item.ProductID,
			&item.AdjustmentType,
			&item.QuantitySystem,
			&item.QuantityPhysical,
			&item.QuantityVariance,
			&item.CostImpact,
			&item.AdjustmentReason,
			&item.AdjustmentDate,
			&item.ApprovedBy,
			&item.CreatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pending approval adjustment: %w", err)
		}
		items = append(items, item)
	}

	return &common.PaginatedResponse{
		Data:       items,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: params.GetTotalPages(total),
		HasMore:    params.GetHasMore(total),
	}, nil
}

// Approve approves a stock adjustment
func (r *StockAdjustmentRepository) Approve(ctx context.Context, id int, approvedBy int) error {
	query := `
		UPDATE stock_adjustments 
		SET approved_by = $2, approved_at = NOW()
		WHERE adjustment_id = $1 AND approved_by IS NULL`
	
	result, err := r.db.ExecContext(ctx, query, id, approvedBy)
	if err != nil {
		return fmt.Errorf("failed to approve stock adjustment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("stock adjustment with ID %d not found or already approved", id)
	}

	return nil
}

// GetVarianceReport retrieves a variance report
func (r *StockAdjustmentRepository) GetVarianceReport(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	params.Validate()

	baseQuery := `
		SELECT adjustment_id, product_id, adjustment_type, quantity_system, quantity_physical,
			   quantity_variance, cost_impact, adjustment_reason, adjustment_date,
			   approved_by, created_by
		FROM stock_adjustments
		WHERE quantity_variance != 0`

	countQuery := `SELECT COUNT(*) FROM stock_adjustments WHERE quantity_variance != 0`

	whereConditions, args := r.buildWhereConditions(params)
	// Remove variance condition since we already have it
	filteredConditions := []string{}
	filteredArgs := []interface{}{}
	argIndex := 0
	
	for _, condition := range whereConditions {
		if !strings.Contains(condition, "quantity_variance") {
			filteredConditions = append(filteredConditions, condition)
			filteredArgs = append(filteredArgs, args[argIndex])
		}
		argIndex++
	}
	
	if len(filteredConditions) > 0 {
		whereClause := " AND " + strings.Join(filteredConditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
		args = filteredArgs
	} else {
		args = []interface{}{}
	}

	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count variance report: %w", err)
	}

	// Add ordering and pagination (order by variance amount desc)
	baseQuery += ` ORDER BY ABS(quantity_variance) DESC, adjustment_date DESC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get variance report: %w", err)
	}
	defer rows.Close()

	var items []products.StockAdjustmentListItem
	for rows.Next() {
		var item products.StockAdjustmentListItem
		err := rows.Scan(
			&item.AdjustmentID,
			&item.ProductID,
			&item.AdjustmentType,
			&item.QuantitySystem,
			&item.QuantityPhysical,
			&item.QuantityVariance,
			&item.CostImpact,
			&item.AdjustmentReason,
			&item.AdjustmentDate,
			&item.ApprovedBy,
			&item.CreatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan variance report item: %w", err)
		}
		items = append(items, item)
	}

	return &common.PaginatedResponse{
		Data:       items,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: params.GetTotalPages(total),
		HasMore:    params.GetHasMore(total),
	}, nil
}

// buildWhereConditions builds WHERE conditions for queries
func (r *StockAdjustmentRepository) buildWhereConditions(params *products.StockAdjustmentFilterParams) ([]string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.ProductID != nil {
		conditions = append(conditions, fmt.Sprintf("product_id = $%d", argIndex))
		args = append(args, *params.ProductID)
		argIndex++
	}

	if params.AdjustmentType != nil {
		conditions = append(conditions, fmt.Sprintf("adjustment_type = $%d", argIndex))
		args = append(args, *params.AdjustmentType)
		argIndex++
	}

	if params.ApprovedBy != nil {
		conditions = append(conditions, fmt.Sprintf("approved_by = $%d", argIndex))
		args = append(args, *params.ApprovedBy)
		argIndex++
	}

	if params.CreatedBy != nil {
		conditions = append(conditions, fmt.Sprintf("created_by = $%d", argIndex))
		args = append(args, *params.CreatedBy)
		argIndex++
	}

	if params.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("adjustment_date >= $%d", argIndex))
		args = append(args, *params.DateFrom)
		argIndex++
	}

	if params.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("adjustment_date <= $%d", argIndex))
		args = append(args, *params.DateTo)
		argIndex++
	}

	if params.IsApproved != nil {
		if *params.IsApproved {
			conditions = append(conditions, "approved_by IS NOT NULL")
		} else {
			conditions = append(conditions, "approved_by IS NULL")
		}
	}

	if params.HasVariance != nil && *params.HasVariance {
		conditions = append(conditions, "quantity_variance != 0")
	}

	if params.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(adjustment_reason ILIKE $%d OR notes ILIKE $%d)", argIndex, argIndex))
		searchTerm := "%" + params.Search + "%"
		args = append(args, searchTerm)
		argIndex++
	}

	return conditions, args
}