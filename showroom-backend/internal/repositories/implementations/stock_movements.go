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
	// Get current stock quantity first
	currentStock, err := r.GetCurrentStock(ctx, movement.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current stock: %w", err)
	}

	movement.QuantityBefore = currentStock

	// Calculate quantity after based on movement type
	if movement.MovementType == products.MovementTypeIn {
		movement.QuantityAfter = movement.QuantityBefore + movement.QuantityMoved
	} else {
		movement.QuantityAfter = movement.QuantityBefore - movement.QuantityMoved
	}

	// Calculate total value
	movement.TotalValue = float64(movement.QuantityMoved) * movement.UnitCost

	// Set movement date if not provided
	if movement.MovementDate.IsZero() {
		movement.MovementDate = time.Now()
	}

	query := `
		INSERT INTO stock_movements (
			product_id, movement_type, reference_type, reference_id,
			quantity_before, quantity_moved, quantity_after, unit_cost,
			total_value, location_from, location_to, movement_date,
			processed_by, movement_reason, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING movement_id, created_at`

	err = r.db.QueryRowContext(ctx, query,
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

	// Update product stock quantity
	updateStockQuery := `UPDATE products_spare_parts SET stock_quantity = $1 WHERE product_id = $2`
	_, err = r.db.ExecContext(ctx, updateStockQuery, movement.QuantityAfter, movement.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to update product stock: %w", err)
	}

	return movement, nil
}

// GetByID retrieves a stock movement by ID
func (r *StockMovementRepository) GetByID(ctx context.Context, id int) (*products.StockMovement, error) {
	query := `
		SELECT movement_id, product_id, movement_type, reference_type, reference_id,
			   quantity_before, quantity_moved, quantity_after, unit_cost, total_value,
			   location_from, location_to, movement_date, processed_by,
			   movement_reason, notes, created_at
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
			return nil, fmt.Errorf("stock movement not found")
		}
		return nil, fmt.Errorf("failed to get stock movement: %w", err)
	}

	return movement, nil
}

// List retrieves all stock movements with pagination
func (r *StockMovementRepository) List(ctx context.Context, params *products.StockMovementFilterParams) (*common.PaginatedResponse, error) {
	baseQuery := `
		FROM stock_movements sm 
		LEFT JOIN products_spare_parts psp ON sm.product_id = psp.product_id
		WHERE 1=1`
	
	args := []interface{}{}
	whereConditions := []string{}
	argIndex := 1

	// Add filters
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

	if params.DateFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sm.movement_date >= $%d", argIndex))
		args = append(args, *params.DateFrom)
		argIndex++
	}

	if params.DateTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sm.movement_date <= $%d", argIndex))
		args = append(args, *params.DateTo)
		argIndex++
	}

	if params.LocationFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sm.location_from = $%d", argIndex))
		args = append(args, *params.LocationFrom)
		argIndex++
	}

	if params.LocationTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sm.location_to = $%d", argIndex))
		args = append(args, *params.LocationTo)
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
		return nil, fmt.Errorf("failed to count stock movements: %w", err)
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
		sm.movement_id, sm.product_id, sm.movement_type, sm.reference_type,
		sm.reference_id, sm.quantity_before, sm.quantity_moved, sm.quantity_after,
		sm.unit_cost, sm.total_value, sm.movement_date, sm.processed_by`
	
	mainQuery := "SELECT " + selectFields + " " + baseQuery + 
		" ORDER BY sm.movement_date DESC, sm.movement_id DESC LIMIT $" + fmt.Sprintf("%d", argIndex) + 
		" OFFSET $" + fmt.Sprintf("%d", argIndex+1)
	
	args = append(args, params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, mainQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query stock movements: %w", err)
	}
	defer rows.Close()

	var movements []products.StockMovementListItem
	for rows.Next() {
		var movement products.StockMovementListItem
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
			&movement.MovementDate,
			&movement.ProcessedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock movement: %w", err)
		}
		movements = append(movements, movement)
	}

	totalPages := (total + int64(params.Limit) - 1) / int64(params.Limit)

	return &common.PaginatedResponse{
		Data:       movements,
		Total:      int(total),
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: int(totalPages),
		HasMore:    params.Page < int(totalPages),
	}, nil
}

// GetByProductID retrieves stock movements for a specific product
func (r *StockMovementRepository) GetByProductID(ctx context.Context, productID int, params *products.StockMovementFilterParams) (*common.PaginatedResponse, error) {
	// Set productID in params and call List
	params.ProductID = &productID
	return r.List(ctx, params)
}

// GetByReferenceID retrieves stock movements by reference
func (r *StockMovementRepository) GetByReferenceID(ctx context.Context, referenceType products.ReferenceType, referenceID int) ([]products.StockMovement, error) {
	query := `
		SELECT movement_id, product_id, movement_type, reference_type, reference_id,
			   quantity_before, quantity_moved, quantity_after, unit_cost, total_value,
			   location_from, location_to, movement_date, processed_by,
			   movement_reason, notes, created_at
		FROM stock_movements 
		WHERE reference_type = $1 AND reference_id = $2
		ORDER BY movement_date DESC, movement_id DESC`

	rows, err := r.db.QueryContext(ctx, query, referenceType, referenceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query stock movements: %w", err)
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
	movement := &products.StockMovement{
		ProductID:      productID,
		MovementType:   products.MovementTypeIn,
		ReferenceType:  products.ReferenceTypePurchase,
		ReferenceID:    receiptID,
		QuantityMoved:  quantity,
		UnitCost:       unitCost,
		ProcessedBy:    processedBy,
		MovementReason: stringPtr("Goods receipt"),
	}

	_, err := r.Create(ctx, movement)
	return err
}

// CreateMovementForAdjustment creates a stock movement for stock adjustment
func (r *StockMovementRepository) CreateMovementForAdjustment(ctx context.Context, productID int, quantityChange int, unitCost float64, adjustmentID int, processedBy int) error {
	var movementType products.MovementType
	if quantityChange >= 0 {
		movementType = products.MovementTypeIn
	} else {
		movementType = products.MovementTypeOut
		quantityChange = -quantityChange // Make positive for the movement
	}

	movement := &products.StockMovement{
		ProductID:      productID,
		MovementType:   movementType,
		ReferenceType:  products.ReferenceTypeAdjustment,
		ReferenceID:    adjustmentID,
		QuantityMoved:  quantityChange,
		UnitCost:       unitCost,
		ProcessedBy:    processedBy,
		MovementReason: stringPtr("Stock adjustment"),
	}

	_, err := r.Create(ctx, movement)
	return err
}

// GetMovementHistory gets recent stock movements for a product
func (r *StockMovementRepository) GetMovementHistory(ctx context.Context, productID int, limit int) ([]products.StockMovement, error) {
	query := `
		SELECT movement_id, product_id, movement_type, reference_type, reference_id,
			   quantity_before, quantity_moved, quantity_after, unit_cost, total_value,
			   location_from, location_to, movement_date, processed_by,
			   movement_reason, notes, created_at
		FROM stock_movements 
		WHERE product_id = $1
		ORDER BY movement_date DESC, movement_id DESC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, productID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query movement history: %w", err)
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

// GetCurrentStock gets the current stock quantity for a product
func (r *StockMovementRepository) GetCurrentStock(ctx context.Context, productID int) (int, error) {
	query := `SELECT COALESCE(stock_quantity, 0) FROM products_spare_parts WHERE product_id = $1`
	
	var quantity int
	err := r.db.QueryRowContext(ctx, query, productID).Scan(&quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get current stock: %w", err)
	}

	return quantity, nil
}

// BulkCreateMovements creates multiple stock movements
func (r *StockMovementRepository) BulkCreateMovements(ctx context.Context, movements []products.StockMovement) error {
	if len(movements) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, movement := range movements {
		_, err := r.Create(ctx, &movement)
		if err != nil {
			return fmt.Errorf("failed to create stock movement: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}