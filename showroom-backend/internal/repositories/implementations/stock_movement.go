package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

type stockMovementRepository struct {
	db *sql.DB
}

// NewStockMovementRepository creates a new stock movement repository
func NewStockMovementRepository(db *sql.DB) interfaces.StockMovementRepository {
	return &stockMovementRepository{db: db}
}

// Create creates a new stock movement
func (r *stockMovementRepository) Create(ctx context.Context, movement *inventory.StockMovement) (*inventory.StockMovement, error) {
	query := `
		INSERT INTO stock_movements (
			product_id, movement_type, reference_type, reference_id, quantity_before,
			quantity_moved, quantity_after, unit_cost, total_value, location_from,
			location_to, movement_date, processed_by, movement_reason, notes
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING movement_id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		movement.ProductID, movement.MovementType, movement.ReferenceType,
		movement.ReferenceID, movement.QuantityBefore, movement.QuantityMoved,
		movement.QuantityAfter, movement.UnitCost, movement.TotalValue,
		movement.LocationFrom, movement.LocationTo, movement.MovementDate,
		movement.ProcessedBy, movement.MovementReason, movement.Notes,
	).Scan(&movement.MovementID, &movement.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create stock movement: %w", err)
	}

	return movement, nil
}

// GetByID retrieves a stock movement by ID
func (r *stockMovementRepository) GetByID(ctx context.Context, id int) (*inventory.StockMovement, error) {
	query := `
		SELECT sm.movement_id, sm.product_id, sm.movement_type, sm.reference_type,
		       sm.reference_id, sm.quantity_before, sm.quantity_moved, sm.quantity_after,
		       sm.unit_cost, sm.total_value, sm.location_from, sm.location_to,
		       sm.movement_date, sm.processed_by, sm.movement_reason, sm.notes,
		       sm.created_at, p.product_code, p.product_name, p.unit_measure,
		       u.full_name as processed_by_name
		FROM stock_movements sm
		JOIN products_spare_parts p ON sm.product_id = p.product_id
		JOIN users u ON sm.processed_by = u.user_id
		WHERE sm.movement_id = $1`

	movement := &inventory.StockMovement{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&movement.MovementID, &movement.ProductID, &movement.MovementType,
		&movement.ReferenceType, &movement.ReferenceID, &movement.QuantityBefore,
		&movement.QuantityMoved, &movement.QuantityAfter, &movement.UnitCost,
		&movement.TotalValue, &movement.LocationFrom, &movement.LocationTo,
		&movement.MovementDate, &movement.ProcessedBy, &movement.MovementReason,
		&movement.Notes, &movement.CreatedAt, &movement.ProductCode,
		&movement.ProductName, &movement.UnitMeasure, &movement.ProcessedByName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("stock movement not found")
		}
		return nil, fmt.Errorf("failed to get stock movement: %w", err)
	}

	return movement, nil
}

// List retrieves stock movements with filtering and pagination
func (r *stockMovementRepository) List(ctx context.Context, params *inventory.StockMovementFilterParams) ([]inventory.StockMovementListItem, int, error) {
	whereConditions := []string{"1 = 1"}
	args := []interface{}{}
	argIndex := 1

	// Build WHERE conditions
	if params.ProductID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sm.product_id = $%d", argIndex))
		args = append(args, *params.ProductID)
		argIndex++
	}

	if params.MovementType != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sm.movement_type = $%d", argIndex))
		args = append(args, *params.MovementType)
		argIndex++
	}

	if params.ReferenceType != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sm.reference_type = $%d", argIndex))
		args = append(args, *params.ReferenceType)
		argIndex++
	}

	if params.ReferenceID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sm.reference_id = $%d", argIndex))
		args = append(args, *params.ReferenceID)
		argIndex++
	}

	if params.ProcessedBy != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sm.processed_by = $%d", argIndex))
		args = append(args, *params.ProcessedBy)
		argIndex++
	}

	if params.LocationFrom != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("sm.location_from ILIKE $%d", argIndex))
		args = append(args, "%"+params.LocationFrom+"%")
		argIndex++
	}

	if params.LocationTo != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("sm.location_to ILIKE $%d", argIndex))
		args = append(args, "%"+params.LocationTo+"%")
		argIndex++
	}

	if params.Search != "" {
		whereConditions = append(whereConditions, fmt.Sprintf(
			"(p.product_name ILIKE $%d OR p.product_code ILIKE $%d OR sm.movement_reason ILIKE $%d)",
			argIndex, argIndex, argIndex))
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM stock_movements sm
		JOIN products_spare_parts p ON sm.product_id = p.product_id
		WHERE %s`, whereClause)

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count stock movements: %w", err)
	}

	// Main query with pagination
	params.PaginationParams.Validate()
	offset := params.PaginationParams.GetOffset()
	limit := params.PaginationParams.Limit

	query := fmt.Sprintf(`
		SELECT sm.movement_id, p.product_code, p.product_name, sm.movement_type,
		       sm.reference_type, sm.reference_id, sm.quantity_moved, sm.unit_cost,
		       sm.total_value, sm.movement_date, u.full_name as processed_by_name,
		       sm.movement_reason
		FROM stock_movements sm
		JOIN products_spare_parts p ON sm.product_id = p.product_id
		JOIN users u ON sm.processed_by = u.user_id
		WHERE %s
		ORDER BY sm.movement_date DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list stock movements: %w", err)
	}
	defer rows.Close()

	var movements []inventory.StockMovementListItem
	for rows.Next() {
		var movement inventory.StockMovementListItem
		err := rows.Scan(
			&movement.MovementID, &movement.ProductCode, &movement.ProductName,
			&movement.MovementType, &movement.ReferenceType, &movement.ReferenceID,
			&movement.QuantityMoved, &movement.UnitCost, &movement.TotalValue,
			&movement.MovementDate, &movement.ProcessedByName, &movement.MovementReason,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan stock movement: %w", err)
		}
		movements = append(movements, movement)
	}

	return movements, total, nil
}

// GetByProduct retrieves stock movements by product
func (r *stockMovementRepository) GetByProduct(ctx context.Context, productID int, page, limit int) ([]inventory.StockMovementListItem, int, error) {
	params := &inventory.StockMovementFilterParams{
		ProductID: &productID,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// GetByType retrieves stock movements by type
func (r *stockMovementRepository) GetByType(ctx context.Context, movementType inventory.MovementType, page, limit int) ([]inventory.StockMovementListItem, int, error) {
	params := &inventory.StockMovementFilterParams{
		MovementType: &movementType,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// GetByReference retrieves stock movements by reference
func (r *stockMovementRepository) GetByReference(ctx context.Context, refType inventory.ReferenceType, refID int, page, limit int) ([]inventory.StockMovementListItem, int, error) {
	params := &inventory.StockMovementFilterParams{
		ReferenceType: &refType,
		ReferenceID:   &refID,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// GetByDateRange retrieves stock movements by date range
func (r *stockMovementRepository) GetByDateRange(ctx context.Context, startDate, endDate string, page, limit int) ([]inventory.StockMovementListItem, int, error) {
	// This is a simplified implementation - in practice, you'd parse the date strings
	params := &inventory.StockMovementFilterParams{}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// GetAuditTrail retrieves audit trail for a product
func (r *stockMovementRepository) GetAuditTrail(ctx context.Context, productID int, page, limit int) ([]inventory.StockMovementListItem, int, error) {
	return r.GetByProduct(ctx, productID, page, limit)
}

// GetMovementsByUser retrieves movements by user
func (r *stockMovementRepository) GetMovementsByUser(ctx context.Context, userID int, page, limit int) ([]inventory.StockMovementListItem, int, error) {
	params := &inventory.StockMovementFilterParams{
		ProcessedBy: &userID,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// GetMovementSummary retrieves movement summary for a product in date range
func (r *stockMovementRepository) GetMovementSummary(ctx context.Context, productID int, startDate, endDate string) (map[string]interface{}, error) {
	// This is a simplified implementation
	summary := map[string]interface{}{
		"total_in":  0,
		"total_out": 0,
		"net_change": 0,
		"total_value": 0.0,
	}
	return summary, nil
}

// GetValueMovements retrieves total value of movements in date range
func (r *stockMovementRepository) GetValueMovements(ctx context.Context, startDate, endDate string) (float64, error) {
	// This is a simplified implementation
	return 0.0, nil
}