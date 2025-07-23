package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

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
	query := `
		INSERT INTO stock_adjustments (
			product_id, adjustment_type, quantity_system, quantity_physical,
			quantity_variance, cost_impact, adjustment_reason, notes, approved_by,
			adjustment_date, approved_at, supporting_documents_json, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING adjustment_id, created_at`

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
			notes, approved_by, adjustment_date, approved_at, supporting_documents_json,
			created_at, created_by
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
	query := `
		UPDATE stock_adjustments SET
			adjustment_type = $2, quantity_system = $3, quantity_physical = $4,
			quantity_variance = $5, cost_impact = $6, adjustment_reason = $7,
			notes = $8, adjustment_date = $9, supporting_documents_json = $10
		WHERE adjustment_id = $1 AND approved_by IS NULL`

	result, err := r.db.ExecContext(ctx, query,
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

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("stock adjustment not found or already approved")
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
		return fmt.Errorf("stock adjustment not found or already approved")
	}

	return nil
}

// List retrieves stock adjustments with pagination and filtering
func (r *StockAdjustmentRepository) List(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	var args []interface{}
	var conditions []string
	
	// Base query
	baseQuery := `
		FROM stock_adjustments sa
		LEFT JOIN products_spare_parts p ON sa.product_id = p.product_id
		WHERE 1=1`
	
	argIndex := 1

	// Apply filters
	if params.ProductID != nil {
		conditions = append(conditions, fmt.Sprintf("sa.product_id = $%d", argIndex))
		args = append(args, *params.ProductID)
		argIndex++
	}

	if params.AdjustmentType != nil {
		conditions = append(conditions, fmt.Sprintf("sa.adjustment_type = $%d", argIndex))
		args = append(args, *params.AdjustmentType)
		argIndex++
	}

	if params.ApprovedBy != nil {
		conditions = append(conditions, fmt.Sprintf("sa.approved_by = $%d", argIndex))
		args = append(args, *params.ApprovedBy)
		argIndex++
	}

	if params.CreatedBy != nil {
		conditions = append(conditions, fmt.Sprintf("sa.created_by = $%d", argIndex))
		args = append(args, *params.CreatedBy)
		argIndex++
	}

	if params.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("sa.adjustment_date >= $%d", argIndex))
		args = append(args, *params.DateFrom)
		argIndex++
	}

	if params.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("sa.adjustment_date <= $%d", argIndex))
		args = append(args, *params.DateTo)
		argIndex++
	}

	if params.IsApproved != nil {
		if *params.IsApproved {
			conditions = append(conditions, "sa.approved_by IS NOT NULL")
		} else {
			conditions = append(conditions, "sa.approved_by IS NULL")
		}
	}

	if params.HasVariance != nil {
		if *params.HasVariance {
			conditions = append(conditions, "sa.quantity_variance != 0")
		} else {
			conditions = append(conditions, "sa.quantity_variance = 0")
		}
	}

	if params.Search != "" {
		searchCondition := fmt.Sprintf(`(
			p.product_name ILIKE $%d OR
			p.product_code ILIKE $%d OR
			sa.adjustment_reason ILIKE $%d OR
			sa.notes ILIKE $%d
		)`, argIndex, argIndex, argIndex, argIndex)
		conditions = append(conditions, searchCondition)
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	// Combine conditions
	whereClause := baseQuery
	if len(conditions) > 0 {
		whereClause += " AND " + strings.Join(conditions, " AND ")
	}

	// Count total records
	countQuery := "SELECT COUNT(*) " + whereClause
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
		params.Limit = 20
	}

	offset := (params.Page - 1) * params.Limit

	// Get data with pagination
	dataQuery := `
		SELECT sa.adjustment_id, sa.product_id, sa.adjustment_type, sa.quantity_system,
			sa.quantity_physical, sa.quantity_variance, sa.cost_impact, sa.adjustment_reason,
			sa.adjustment_date, sa.approved_by, sa.created_by
		` + whereClause + `
		ORDER BY sa.adjustment_date DESC, sa.adjustment_id DESC
		LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)

	args = append(args, params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock adjustments: %w", err)
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

	totalPages := (total + int64(params.Limit) - 1) / int64(params.Limit)

	return &common.PaginatedResponse{
		Data:        items,
		Total: int(total),
		Page:        params.Page,
		Limit:       params.Limit,
		TotalPages:  int(totalPages),
		HasMore:     params.Page < int(totalPages),
		
	}, nil
}

// GetByProductID retrieves stock adjustments by product ID with pagination
func (r *StockAdjustmentRepository) GetByProductID(ctx context.Context, productID int, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	params.ProductID = &productID
	return r.List(ctx, params)
}

// GetPendingApproval retrieves stock adjustments pending approval with pagination
func (r *StockAdjustmentRepository) GetPendingApproval(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	isApproved := false
	params.IsApproved = &isApproved
	return r.List(ctx, params)
}

// Approve approves a stock adjustment
func (r *StockAdjustmentRepository) Approve(ctx context.Context, id int, approvedBy int) error {
	query := `
		UPDATE stock_adjustments SET
			approved_by = $2,
			approved_at = NOW()
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
		return fmt.Errorf("stock adjustment not found or already approved")
	}

	return nil
}

// GetVarianceReport retrieves variance report with pagination
func (r *StockAdjustmentRepository) GetVarianceReport(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error) {
	hasVariance := true
	params.HasVariance = &hasVariance
	return r.List(ctx, params)
}