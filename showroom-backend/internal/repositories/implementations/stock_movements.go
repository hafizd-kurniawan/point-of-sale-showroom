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
			return nil, fmt.Errorf("stock movement not found")
		}
		return nil, fmt.Errorf("failed to get stock movement: %w", err)
	}

	return movement, nil
}

// List retrieves stock movements with pagination and filtering
func (r *StockMovementRepository) List(ctx context.Context, params *products.StockMovementFilterParams) (*common.PaginatedResponse, error) {
	var args []interface{}
	var conditions []string
	
	// Base query
	baseQuery := `
		FROM stock_movements sm
		LEFT JOIN products_spare_parts p ON sm.product_id = p.product_id
		WHERE 1=1`
	
	argIndex := 1

	// Apply filters
	if params.ProductID != nil {
		conditions = append(conditions, fmt.Sprintf("sm.product_id = $%d", argIndex))
		args = append(args, *params.ProductID)
		argIndex++
	}

	if params.MovementType != nil {
		conditions = append(conditions, fmt.Sprintf("sm.movement_type = $%d", argIndex))
		args = append(args, *params.MovementType)
		argIndex++
	}

	if params.ReferenceType != nil {
		conditions = append(conditions, fmt.Sprintf("sm.reference_type = $%d", argIndex))
		args = append(args, *params.ReferenceType)
		argIndex++
	}

	if params.ReferenceID != nil {
		conditions = append(conditions, fmt.Sprintf("sm.reference_id = $%d", argIndex))
		args = append(args, *params.ReferenceID)
		argIndex++
	}

	if params.ProcessedBy != nil {
		conditions = append(conditions, fmt.Sprintf("sm.processed_by = $%d", argIndex))
		args = append(args, *params.ProcessedBy)
		argIndex++
	}

	if params.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("sm.movement_date >= $%d", argIndex))
		args = append(args, *params.DateFrom)
		argIndex++
	}

	if params.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("sm.movement_date <= $%d", argIndex))
		args = append(args, *params.DateTo)
		argIndex++
	}

	if params.LocationFrom != nil {
		conditions = append(conditions, fmt.Sprintf("sm.location_from ILIKE $%d", argIndex))
		args = append(args, "%"+*params.LocationFrom+"%")
		argIndex++
	}

	if params.LocationTo != nil {
		conditions = append(conditions, fmt.Sprintf("sm.location_to ILIKE $%d", argIndex))
		args = append(args, "%"+*params.LocationTo+"%")
		argIndex++
	}

	if params.Search != "" {
		searchCondition := fmt.Sprintf(`(
			p.product_name ILIKE $%d OR
			p.product_code ILIKE $%d OR
			sm.movement_reason ILIKE $%d OR
			sm.notes ILIKE $%d
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
		return nil, fmt.Errorf("failed to count stock movements: %w", err)
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
		SELECT sm.movement_id, sm.product_id, sm.movement_type, sm.reference_type,
			sm.reference_id, sm.quantity_before, sm.quantity_moved, sm.quantity_after,
			sm.unit_cost, sm.total_value, sm.movement_date, sm.processed_by
		` + whereClause + `
		ORDER BY sm.movement_date DESC, sm.movement_id DESC
		LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)

	args = append(args, params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock movements: %w", err)
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

// GetByProductID retrieves stock movements by product ID with pagination
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

// CreateMovementForReceipt creates stock movement for goods receipt
func (r *StockMovementRepository) CreateMovementForReceipt(ctx context.Context, productID int, quantity int, unitCost float64, receiptID int, processedBy int) error {
	// Get current stock quantity first
	currentStock, err := r.GetCurrentStock(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get current stock: %w", err)
	}

	movement := &products.StockMovement{
		ProductID:      productID,
		MovementType:   products.MovementTypeIn,
		ReferenceType:  products.ReferenceTypePurchase,
		ReferenceID:    receiptID,
		QuantityBefore: currentStock,
		QuantityMoved:  quantity,
		QuantityAfter:  currentStock + quantity,
		UnitCost:       unitCost,
		TotalValue:     float64(quantity) * unitCost,
		MovementDate:   time.Now(),
		ProcessedBy:    processedBy,
		MovementReason: stringPtr("Goods receipt"),
	}

	_, err = r.Create(ctx, movement)
	return err
}

// CreateMovementForAdjustment creates stock movement for stock adjustment
func (r *StockMovementRepository) CreateMovementForAdjustment(ctx context.Context, productID int, quantityChange int, unitCost float64, adjustmentID int, processedBy int) error {
	// Get current stock quantity first
	currentStock, err := r.GetCurrentStock(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get current stock: %w", err)
	}

	movement := &products.StockMovement{
		ProductID:      productID,
		MovementType:   products.MovementTypeAdjustment,
		ReferenceType:  products.ReferenceTypeAdjustment,
		ReferenceID:    adjustmentID,
		QuantityBefore: currentStock,
		QuantityMoved:  quantityChange,
		QuantityAfter:  currentStock + quantityChange,
		UnitCost:       unitCost,
		TotalValue:     float64(quantityChange) * unitCost,
		MovementDate:   time.Now(),
		ProcessedBy:    processedBy,
		MovementReason: stringPtr("Stock adjustment"),
	}

	_, err = r.Create(ctx, movement)
	return err
}

// GetMovementHistory retrieves movement history for a product with limit
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
			return nil, fmt.Errorf("failed to scan stock movement: %w", err)
		}
		movements = append(movements, movement)
	}

	return movements, nil
}

// GetCurrentStock retrieves current stock quantity for a product
func (r *StockMovementRepository) GetCurrentStock(ctx context.Context, productID int) (int, error) {
	// Get from the latest stock movement if available
	query := `
		SELECT quantity_after 
		FROM stock_movements 
		WHERE product_id = $1 
		ORDER BY movement_date DESC, movement_id DESC 
		LIMIT 1`

	var currentStock int
	err := r.db.QueryRowContext(ctx, query, productID).Scan(&currentStock)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no movements, get from products table
			stockQuery := `SELECT COALESCE(stock_quantity, 0) FROM products_spare_parts WHERE product_id = $1`
			err = r.db.QueryRowContext(ctx, stockQuery, productID).Scan(&currentStock)
			if err != nil {
				if err == sql.ErrNoRows {
					return 0, fmt.Errorf("product not found")
				}
				return 0, fmt.Errorf("failed to get stock from product: %w", err)
			}
		} else {
			return 0, fmt.Errorf("failed to get current stock: %w", err)
		}
	}

	return currentStock, nil
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

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO stock_movements (
			product_id, movement_type, reference_type, reference_id, quantity_before,
			quantity_moved, quantity_after, unit_cost, total_value, location_from,
			location_to, movement_date, processed_by, movement_reason, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, movement := range movements {
		_, err = stmt.ExecContext(ctx,
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
		)
		if err != nil {
			return fmt.Errorf("failed to insert stock movement: %w", err)
		}
	}

	return tx.Commit()
}

// helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}