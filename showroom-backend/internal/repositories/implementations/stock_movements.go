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

// StockMovementRepository implements interfaces.StockMovementRepository
type StockMovementRepository struct {
	db *sql.DB
}

// NewStockMovementRepository creates a new stock movement repository
func NewStockMovementRepository(db *sql.DB) interfaces.StockMovementRepository {
	return &StockMovementRepository{db: db}
}

// Create creates a new stock movement
func (r *StockMovementRepository) Create(ctx context.Context, movement *products.StockMovement) (*products.StockMovement, error) {
	query := `
		INSERT INTO stock_movements (
			product_id, movement_type, reference_type, reference_id, quantity_before,
			quantity_moved, quantity_after, unit_cost, total_value, location_from,
			location_to, movement_date, processed_by, movement_reason, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING movement_id, created_at`

	// Calculate total value if not set
	if movement.TotalValue == 0 {
		movement.CalculateTotalValue()
	}

	// Set movement date if not provided
	if movement.MovementDate.IsZero() {
		movement.MovementDate = time.Now()
	}

	err := r.db.QueryRowContext(ctx, query,
		movement.ProductID,
		movement.MovementType,
		movement.ReferenceType,
		movement.ReferenceID,
		movement.QuantityBefore,
		movement.QuantityMoved,
		movement.QuantityAfter,
		movement.UnitCost,
		movement.TotalValue,
		movement.LocationFrom,
		movement.LocationTo,
		movement.MovementDate,
		movement.ProcessedBy,
		movement.MovementReason,
		movement.Notes,
	).Scan(&movement.MovementID, &movement.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create stock movement: %w", err)
	}

	return movement, nil
}

// GetByID retrieves a stock movement by ID
func (r *StockMovementRepository) GetByID(ctx context.Context, id int) (*products.StockMovement, error) {
	query := `
		SELECT movement_id, product_id, movement_type, reference_type, reference_id,
			   quantity_before, quantity_moved, quantity_after, unit_cost, total_value,
			   location_from, location_to, movement_date, processed_by, movement_reason,
			   notes, created_at
		FROM stock_movements
		WHERE movement_id = $1`

	movement := &products.StockMovement{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&movement.MovementID,
		&movement.ProductID,
		&movement.MovementType,
		&movement.ReferenceType,
		&movement.ReferenceID,
		&movement.QuantityBefore,
		&movement.QuantityMoved,
		&movement.QuantityAfter,
		&movement.UnitCost,
		&movement.TotalValue,
		&movement.LocationFrom,
		&movement.LocationTo,
		&movement.MovementDate,
		&movement.ProcessedBy,
		&movement.MovementReason,
		&movement.Notes,
		&movement.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("stock movement with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get stock movement: %w", err)
	}

	return movement, nil
}

// List retrieves a paginated list of stock movements
func (r *StockMovementRepository) List(ctx context.Context, params *products.StockMovementFilterParams) (*common.PaginatedResponse, error) {
	params.Validate()

	baseQuery := `
		SELECT movement_id, product_id, movement_type, reference_type, reference_id,
			   quantity_before, quantity_moved, quantity_after, unit_cost, total_value,
			   movement_date, processed_by
		FROM stock_movements`

	countQuery := `SELECT COUNT(*) FROM stock_movements`

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
		return nil, fmt.Errorf("failed to count stock movements: %w", err)
	}

	// Add ordering and pagination
	baseQuery += ` ORDER BY movement_date DESC, movement_id DESC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock movements: %w", err)
	}
	defer rows.Close()

	var items []products.StockMovementListItem
	for rows.Next() {
		var item products.StockMovementListItem
		err := rows.Scan(
			&item.MovementID,
			&item.ProductID,
			&item.MovementType,
			&item.ReferenceType,
			&item.ReferenceID,
			&item.QuantityBefore,
			&item.QuantityMoved,
			&item.QuantityAfter,
			&item.UnitCost,
			&item.TotalValue,
			&item.MovementDate,
			&item.ProcessedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock movement: %w", err)
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

// GetByProductID retrieves stock movements by product ID
func (r *StockMovementRepository) GetByProductID(ctx context.Context, productID int, params *products.StockMovementFilterParams) (*common.PaginatedResponse, error) {
	params.ProductID = &productID
	return r.List(ctx, params)
}

// GetByReferenceID retrieves stock movements by reference type and ID
func (r *StockMovementRepository) GetByReferenceID(ctx context.Context, referenceType products.ReferenceType, referenceID int) ([]products.StockMovement, error) {
	query := `
		SELECT movement_id, product_id, movement_type, reference_type, reference_id,
			   quantity_before, quantity_moved, quantity_after, unit_cost, total_value,
			   location_from, location_to, movement_date, processed_by, movement_reason,
			   notes, created_at
		FROM stock_movements
		WHERE reference_type = $1 AND reference_id = $2
		ORDER BY movement_date DESC, movement_id DESC`

	rows, err := r.db.QueryContext(ctx, query, referenceType, referenceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock movements by reference: %w", err)
	}
	defer rows.Close()

	var movements []products.StockMovement
	for rows.Next() {
		var movement products.StockMovement
		err := rows.Scan(
			&movement.MovementID,
			&movement.ProductID,
			&movement.MovementType,
			&movement.ReferenceType,
			&movement.ReferenceID,
			&movement.QuantityBefore,
			&movement.QuantityMoved,
			&movement.QuantityAfter,
			&movement.UnitCost,
			&movement.TotalValue,
			&movement.LocationFrom,
			&movement.LocationTo,
			&movement.MovementDate,
			&movement.ProcessedBy,
			&movement.MovementReason,
			&movement.Notes,
			&movement.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock movement: %w", err)
		}
		movements = append(movements, movement)
	}

	return movements, nil
}

// CreateMovementForReceipt creates a stock movement for goods receipt
func (r *StockMovementRepository) CreateMovementForReceipt(ctx context.Context, productID int, quantity int, unitCost float64, receiptID int, processedBy int) error {
	// Get current stock quantity
	currentStock, err := r.GetCurrentStock(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get current stock: %w", err)
	}

	movement := &products.StockMovement{
		ProductID:       productID,
		MovementType:    products.MovementTypeIn,
		ReferenceType:   products.ReferenceTypePurchase,
		ReferenceID:     receiptID,
		QuantityBefore:  currentStock,
		QuantityMoved:   quantity,
		QuantityAfter:   currentStock + quantity,
		UnitCost:        unitCost,
		MovementDate:    time.Now(),
		ProcessedBy:     processedBy,
		MovementReason:  stringPtr("Goods receipt"),
	}

	movement.CalculateTotalValue()

	_, err = r.Create(ctx, movement)
	if err != nil {
		return fmt.Errorf("failed to create movement for receipt: %w", err)
	}

	return nil
}

// CreateMovementForAdjustment creates a stock movement for adjustment
func (r *StockMovementRepository) CreateMovementForAdjustment(ctx context.Context, productID int, quantityChange int, unitCost float64, adjustmentID int, processedBy int) error {
	// Get current stock quantity
	currentStock, err := r.GetCurrentStock(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get current stock: %w", err)
	}

	movement := &products.StockMovement{
		ProductID:       productID,
		MovementType:    products.MovementTypeAdjustment,
		ReferenceType:   products.ReferenceTypeAdjustment,
		ReferenceID:     adjustmentID,
		QuantityBefore:  currentStock,
		QuantityMoved:   quantityChange,
		QuantityAfter:   currentStock + quantityChange,
		UnitCost:        unitCost,
		MovementDate:    time.Now(),
		ProcessedBy:     processedBy,
		MovementReason:  stringPtr("Stock adjustment"),
	}

	movement.CalculateTotalValue()

	_, err = r.Create(ctx, movement)
	if err != nil {
		return fmt.Errorf("failed to create movement for adjustment: %w", err)
	}

	return nil
}

// GetMovementHistory retrieves movement history for a product
func (r *StockMovementRepository) GetMovementHistory(ctx context.Context, productID int, limit int) ([]products.StockMovement, error) {
	query := `
		SELECT movement_id, product_id, movement_type, reference_type, reference_id,
			   quantity_before, quantity_moved, quantity_after, unit_cost, total_value,
			   location_from, location_to, movement_date, processed_by, movement_reason,
			   notes, created_at
		FROM stock_movements
		WHERE product_id = $1
		ORDER BY movement_date DESC, movement_id DESC
		LIMIT $2`

	if limit <= 0 {
		limit = 50 // Default limit
	}

	rows, err := r.db.QueryContext(ctx, query, productID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get movement history: %w", err)
	}
	defer rows.Close()

	var movements []products.StockMovement
	for rows.Next() {
		var movement products.StockMovement
		err := rows.Scan(
			&movement.MovementID,
			&movement.ProductID,
			&movement.MovementType,
			&movement.ReferenceType,
			&movement.ReferenceID,
			&movement.QuantityBefore,
			&movement.QuantityMoved,
			&movement.QuantityAfter,
			&movement.UnitCost,
			&movement.TotalValue,
			&movement.LocationFrom,
			&movement.LocationTo,
			&movement.MovementDate,
			&movement.ProcessedBy,
			&movement.MovementReason,
			&movement.Notes,
			&movement.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movement history: %w", err)
		}
		movements = append(movements, movement)
	}

	return movements, nil
}

// GetCurrentStock gets the current stock quantity for a product
func (r *StockMovementRepository) GetCurrentStock(ctx context.Context, productID int) (int, error) {
	// Get from products table first
	query := `SELECT COALESCE(stock_quantity, 0) FROM products_spare_parts WHERE product_id = $1`
	
	var stock int
	err := r.db.QueryRowContext(ctx, query, productID).Scan(&stock)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get current stock: %w", err)
	}

	return stock, nil
}

// BulkCreateMovements creates multiple stock movements in a transaction
func (r *StockMovementRepository) BulkCreateMovements(ctx context.Context, movements []products.StockMovement) error {
	if len(movements) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO stock_movements (
			product_id, movement_type, reference_type, reference_id, quantity_before,
			quantity_moved, quantity_after, unit_cost, total_value, location_from,
			location_to, movement_date, processed_by, movement_reason, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	for i := range movements {
		// Calculate total value if not set
		if movements[i].TotalValue == 0 {
			movements[i].CalculateTotalValue()
		}

		// Set movement date if not provided
		if movements[i].MovementDate.IsZero() {
			movements[i].MovementDate = time.Now()
		}

		_, err = tx.ExecContext(ctx, query,
			movements[i].ProductID,
			movements[i].MovementType,
			movements[i].ReferenceType,
			movements[i].ReferenceID,
			movements[i].QuantityBefore,
			movements[i].QuantityMoved,
			movements[i].QuantityAfter,
			movements[i].UnitCost,
			movements[i].TotalValue,
			movements[i].LocationFrom,
			movements[i].LocationTo,
			movements[i].MovementDate,
			movements[i].ProcessedBy,
			movements[i].MovementReason,
			movements[i].Notes,
		)
		if err != nil {
			return fmt.Errorf("failed to create stock movement at index %d: %w", i, err)
		}
	}

	return tx.Commit()
}

// buildWhereConditions builds WHERE conditions for queries
func (r *StockMovementRepository) buildWhereConditions(params *products.StockMovementFilterParams) ([]string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.ProductID != nil {
		conditions = append(conditions, fmt.Sprintf("product_id = $%d", argIndex))
		args = append(args, *params.ProductID)
		argIndex++
	}

	if params.MovementType != nil {
		conditions = append(conditions, fmt.Sprintf("movement_type = $%d", argIndex))
		args = append(args, *params.MovementType)
		argIndex++
	}

	if params.ReferenceType != nil {
		conditions = append(conditions, fmt.Sprintf("reference_type = $%d", argIndex))
		args = append(args, *params.ReferenceType)
		argIndex++
	}

	if params.ReferenceID != nil {
		conditions = append(conditions, fmt.Sprintf("reference_id = $%d", argIndex))
		args = append(args, *params.ReferenceID)
		argIndex++
	}

	if params.ProcessedBy != nil {
		conditions = append(conditions, fmt.Sprintf("processed_by = $%d", argIndex))
		args = append(args, *params.ProcessedBy)
		argIndex++
	}

	if params.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("movement_date >= $%d", argIndex))
		args = append(args, *params.DateFrom)
		argIndex++
	}

	if params.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("movement_date <= $%d", argIndex))
		args = append(args, *params.DateTo)
		argIndex++
	}

	if params.LocationFrom != nil {
		conditions = append(conditions, fmt.Sprintf("location_from = $%d", argIndex))
		args = append(args, *params.LocationFrom)
		argIndex++
	}

	if params.LocationTo != nil {
		conditions = append(conditions, fmt.Sprintf("location_to = $%d", argIndex))
		args = append(args, *params.LocationTo)
		argIndex++
	}

	if params.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(movement_reason ILIKE $%d OR notes ILIKE $%d)", argIndex, argIndex))
		searchTerm := "%" + params.Search + "%"
		args = append(args, searchTerm)
		argIndex++
	}

	return conditions, args
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}