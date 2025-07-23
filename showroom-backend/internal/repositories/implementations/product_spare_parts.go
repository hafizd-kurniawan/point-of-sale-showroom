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

// ProductSparePartRepository implements interfaces.ProductSparePartRepository
type ProductSparePartRepository struct {
	db *sql.DB
}

// NewProductSparePartRepository creates a new product spare part repository
func NewProductSparePartRepository(db *sql.DB) interfaces.ProductSparePartRepository {
	return &ProductSparePartRepository{db: db}
}

// Create creates a new product spare part
func (r *ProductSparePartRepository) Create(ctx context.Context, product *products.ProductSparePart) (*products.ProductSparePart, error) {
	query := `
		INSERT INTO products_spare_parts (
			product_code, product_name, description, brand_id, category_id, unit_measure,
			cost_price, selling_price, markup_percentage, stock_quantity, min_stock_level,
			max_stock_level, location_rack, barcode, weight, dimensions, created_by,
			is_active, product_image, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
		RETURNING product_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		product.ProductCode,
		product.ProductName,
		product.Description,
		product.BrandID,
		product.CategoryID,
		product.UnitMeasure,
		product.CostPrice,
		product.SellingPrice,
		product.MarkupPercentage,
		product.StockQuantity,
		product.MinStockLevel,
		product.MaxStockLevel,
		product.LocationRack,
		product.Barcode,
		product.Weight,
		product.Dimensions,
		product.CreatedBy,
		product.IsActive,
		product.ProductImage,
		product.Notes,
	).Scan(&product.ProductID, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create product spare part: %w", err)
	}

	return product, nil
}

// GetByID retrieves a product spare part by ID
func (r *ProductSparePartRepository) GetByID(ctx context.Context, id int) (*products.ProductSparePart, error) {
	query := `
		SELECT product_id, product_code, product_name, description, brand_id, category_id,
			   unit_measure, cost_price, selling_price, markup_percentage, stock_quantity,
			   min_stock_level, max_stock_level, location_rack, barcode, weight, dimensions,
			   created_at, updated_at, created_by, is_active, product_image, notes
		FROM products_spare_parts
		WHERE product_id = $1`

	product := &products.ProductSparePart{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ProductID,
		&product.ProductCode,
		&product.ProductName,
		&product.Description,
		&product.BrandID,
		&product.CategoryID,
		&product.UnitMeasure,
		&product.CostPrice,
		&product.SellingPrice,
		&product.MarkupPercentage,
		&product.StockQuantity,
		&product.MinStockLevel,
		&product.MaxStockLevel,
		&product.LocationRack,
		&product.Barcode,
		&product.Weight,
		&product.Dimensions,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.CreatedBy,
		&product.IsActive,
		&product.ProductImage,
		&product.Notes,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product spare part with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get product spare part: %w", err)
	}

	return product, nil
}

// GetByCode retrieves a product spare part by code
func (r *ProductSparePartRepository) GetByCode(ctx context.Context, code string) (*products.ProductSparePart, error) {
	query := `
		SELECT product_id, product_code, product_name, description, brand_id, category_id,
			   unit_measure, cost_price, selling_price, markup_percentage, stock_quantity,
			   min_stock_level, max_stock_level, location_rack, barcode, weight, dimensions,
			   created_at, updated_at, created_by, is_active, product_image, notes
		FROM products_spare_parts
		WHERE product_code = $1`

	product := &products.ProductSparePart{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&product.ProductID,
		&product.ProductCode,
		&product.ProductName,
		&product.Description,
		&product.BrandID,
		&product.CategoryID,
		&product.UnitMeasure,
		&product.CostPrice,
		&product.SellingPrice,
		&product.MarkupPercentage,
		&product.StockQuantity,
		&product.MinStockLevel,
		&product.MaxStockLevel,
		&product.LocationRack,
		&product.Barcode,
		&product.Weight,
		&product.Dimensions,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.CreatedBy,
		&product.IsActive,
		&product.ProductImage,
		&product.Notes,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product spare part with code %s not found", code)
		}
		return nil, fmt.Errorf("failed to get product spare part: %w", err)
	}

	return product, nil
}

// GetByBarcode retrieves a product spare part by barcode
func (r *ProductSparePartRepository) GetByBarcode(ctx context.Context, barcode string) (*products.ProductSparePart, error) {
	query := `
		SELECT product_id, product_code, product_name, description, brand_id, category_id,
			   unit_measure, cost_price, selling_price, markup_percentage, stock_quantity,
			   min_stock_level, max_stock_level, location_rack, barcode, weight, dimensions,
			   created_at, updated_at, created_by, is_active, product_image, notes
		FROM products_spare_parts
		WHERE barcode = $1`

	product := &products.ProductSparePart{}
	err := r.db.QueryRowContext(ctx, query, barcode).Scan(
		&product.ProductID,
		&product.ProductCode,
		&product.ProductName,
		&product.Description,
		&product.BrandID,
		&product.CategoryID,
		&product.UnitMeasure,
		&product.CostPrice,
		&product.SellingPrice,
		&product.MarkupPercentage,
		&product.StockQuantity,
		&product.MinStockLevel,
		&product.MaxStockLevel,
		&product.LocationRack,
		&product.Barcode,
		&product.Weight,
		&product.Dimensions,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.CreatedBy,
		&product.IsActive,
		&product.ProductImage,
		&product.Notes,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product spare part with barcode %s not found", barcode)
		}
		return nil, fmt.Errorf("failed to get product spare part: %w", err)
	}

	return product, nil
}

// Update updates a product spare part
func (r *ProductSparePartRepository) Update(ctx context.Context, id int, product *products.ProductSparePart) (*products.ProductSparePart, error) {
	query := `
		UPDATE products_spare_parts SET
			product_name = $2, description = $3, brand_id = $4, category_id = $5,
			unit_measure = $6, cost_price = $7, selling_price = $8, markup_percentage = $9,
			min_stock_level = $10, max_stock_level = $11, location_rack = $12, barcode = $13,
			weight = $14, dimensions = $15, is_active = $16, product_image = $17, notes = $18,
			updated_at = NOW()
		WHERE product_id = $1
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		id,
		product.ProductName,
		product.Description,
		product.BrandID,
		product.CategoryID,
		product.UnitMeasure,
		product.CostPrice,
		product.SellingPrice,
		product.MarkupPercentage,
		product.MinStockLevel,
		product.MaxStockLevel,
		product.LocationRack,
		product.Barcode,
		product.Weight,
		product.Dimensions,
		product.IsActive,
		product.ProductImage,
		product.Notes,
	).Scan(&product.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update product spare part: %w", err)
	}

	return r.GetByID(ctx, id)
}

// Delete deletes a product spare part
func (r *ProductSparePartRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM products_spare_parts WHERE product_id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product spare part: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product spare part with ID %d not found", id)
	}

	return nil
}

// List retrieves a paginated list of product spare parts
func (r *ProductSparePartRepository) List(ctx context.Context, params *products.ProductSparePartFilterParams) (*common.PaginatedResponse, error) {
	params.Validate()

	baseQuery := `
		SELECT product_id, product_code, product_name, brand_id, category_id,
			   unit_measure, cost_price, selling_price, stock_quantity,
			   min_stock_level, location_rack, is_active
		FROM products_spare_parts`

	countQuery := `SELECT COUNT(*) FROM products_spare_parts`

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
		return nil, fmt.Errorf("failed to count product spare parts: %w", err)
	}

	// Add ordering and pagination
	baseQuery += ` ORDER BY product_name ASC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list product spare parts: %w", err)
	}
	defer rows.Close()

	var items []products.ProductSparePartListItem
	for rows.Next() {
		var item products.ProductSparePartListItem
		err := rows.Scan(
			&item.ProductID,
			&item.ProductCode,
			&item.ProductName,
			&item.BrandID,
			&item.CategoryID,
			&item.UnitMeasure,
			&item.CostPrice,
			&item.SellingPrice,
			&item.StockQuantity,
			&item.MinStockLevel,
			&item.LocationRack,
			&item.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product spare part: %w", err)
		}
		items = append(items, item)
	}

	return &common.PaginatedResponse{
		Data:       items,
Total: int(total),
Page:       params.Page,
Limit:      params.Limit,
TotalPages: params.GetTotalPages(total),
HasMore:    params.GetHasMore(total),
	}, nil
}

// GetLowStockProducts retrieves products below minimum stock level
func (r *ProductSparePartRepository) GetLowStockProducts(ctx context.Context, params *products.ProductSparePartFilterParams) (*common.PaginatedResponse, error) {
	params.Validate()

	baseQuery := `
		SELECT product_id, product_code, product_name, brand_id, category_id,
			   unit_measure, cost_price, selling_price, stock_quantity,
			   min_stock_level, location_rack, is_active
		FROM products_spare_parts
		WHERE stock_quantity <= min_stock_level AND is_active = true`

	countQuery := `SELECT COUNT(*) FROM products_spare_parts WHERE stock_quantity <= min_stock_level AND is_active = true`

	whereConditions, args := r.buildWhereConditions(params)
	if len(whereConditions) > 0 {
		// Skip is_active condition since we already have it
		filteredConditions := []string{}
		filteredArgs := []interface{}{}
		argIndex := 0
		
		for _, condition := range whereConditions {
			if !strings.Contains(condition, "is_active") {
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
	} else {
		args = []interface{}{}
	}

	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count low stock products: %w", err)
	}

	// Add ordering and pagination
	baseQuery += ` ORDER BY (stock_quantity - min_stock_level) ASC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock products: %w", err)
	}
	defer rows.Close()

	var items []products.ProductSparePartListItem
	for rows.Next() {
		var item products.ProductSparePartListItem
		err := rows.Scan(
			&item.ProductID,
			&item.ProductCode,
			&item.ProductName,
			&item.BrandID,
			&item.CategoryID,
			&item.UnitMeasure,
			&item.CostPrice,
			&item.SellingPrice,
			&item.StockQuantity,
			&item.MinStockLevel,
			&item.LocationRack,
			&item.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan low stock product: %w", err)
		}
		items = append(items, item)
	}

	return &common.PaginatedResponse{
		Data:       items,
Total: int(total),
Page:       params.Page,
Limit:      params.Limit,
TotalPages: params.GetTotalPages(total),
HasMore:    params.GetHasMore(total),
	}, nil
}

// UpdateStock updates the stock quantity of a product
func (r *ProductSparePartRepository) UpdateStock(ctx context.Context, id int, newQuantity int) error {
	query := `UPDATE products_spare_parts SET stock_quantity = $2, updated_at = NOW() WHERE product_id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id, newQuantity)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %d not found", id)
	}

	return nil
}

// UpdateStockWithMovement updates stock and creates movement record in a transaction
func (r *ProductSparePartRepository) UpdateStockWithMovement(ctx context.Context, id int, quantityChange int, movementDetails *products.StockMovement) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get current stock
	var currentStock int
	err = tx.QueryRowContext(ctx, "SELECT stock_quantity FROM products_spare_parts WHERE product_id = $1", id).Scan(&currentStock)
	if err != nil {
		return fmt.Errorf("failed to get current stock: %w", err)
	}

	// Calculate new stock
	newStock := currentStock + quantityChange
	if newStock < 0 {
		return fmt.Errorf("insufficient stock: current %d, requested change %d", currentStock, quantityChange)
	}

	// Update stock
	_, err = tx.ExecContext(ctx, "UPDATE products_spare_parts SET stock_quantity = $2, updated_at = NOW() WHERE product_id = $1", id, newStock)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// Create movement record
	movementDetails.ProductID = id
	movementDetails.QuantityBefore = currentStock
	movementDetails.QuantityMoved = quantityChange
	movementDetails.QuantityAfter = newStock
	movementDetails.CreatedAt = time.Now()

	movementQuery := `
		INSERT INTO stock_movements (
			product_id, movement_type, reference_type, reference_id, quantity_before,
			quantity_moved, quantity_after, unit_cost, total_value, location_from,
			location_to, movement_date, processed_by, movement_reason, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	_, err = tx.ExecContext(ctx, movementQuery,
		movementDetails.ProductID,
		movementDetails.MovementType,
		movementDetails.ReferenceType,
		movementDetails.ReferenceID,
		movementDetails.QuantityBefore,
		movementDetails.QuantityMoved,
		movementDetails.QuantityAfter,
		movementDetails.UnitCost,
		movementDetails.TotalValue,
		movementDetails.LocationFrom,
		movementDetails.LocationTo,
		movementDetails.MovementDate,
		movementDetails.ProcessedBy,
		movementDetails.MovementReason,
		movementDetails.Notes,
	)
	if err != nil {
		return fmt.Errorf("failed to create stock movement: %w", err)
	}

	return tx.Commit()
}

// GenerateCode generates a new product code
func (r *ProductSparePartRepository) GenerateCode(ctx context.Context) (string, error) {
	query := `
		SELECT COALESCE(MAX(CAST(SUBSTRING(product_code FROM 5) AS INTEGER)), 0) + 1
		FROM products_spare_parts
		WHERE product_code ~ '^PRD-[0-9]+$'`

	var nextNumber int
	err := r.db.QueryRowContext(ctx, query).Scan(&nextNumber)
	if err != nil {
		return "", fmt.Errorf("failed to generate product code: %w", err)
	}

	return fmt.Sprintf("PRD-%03d", nextNumber), nil
}

// IsCodeExists checks if a product code already exists
func (r *ProductSparePartRepository) IsCodeExists(ctx context.Context, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM products_spare_parts WHERE product_code = $1)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, code).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if code exists: %w", err)
	}

	return exists, nil
}

// IsBarcodeExists checks if a barcode already exists (excluding specific ID)
func (r *ProductSparePartRepository) IsBarcodeExists(ctx context.Context, barcode string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM products_spare_parts WHERE barcode = $1 AND product_id != $2)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, barcode, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if barcode exists: %w", err)
	}

	return exists, nil
}

// GetByBrandID retrieves products by brand ID
func (r *ProductSparePartRepository) GetByBrandID(ctx context.Context, brandID int, params *products.ProductSparePartFilterParams) (*common.PaginatedResponse, error) {
	params.BrandID = &brandID
	return r.List(ctx, params)
}

// GetByCategoryID retrieves products by category ID
func (r *ProductSparePartRepository) GetByCategoryID(ctx context.Context, categoryID int, params *products.ProductSparePartFilterParams) (*common.PaginatedResponse, error) {
	params.CategoryID = &categoryID
	return r.List(ctx, params)
}

// Search searches products by query
func (r *ProductSparePartRepository) Search(ctx context.Context, query string, params *products.ProductSparePartFilterParams) (*common.PaginatedResponse, error) {
	params.Search = query
	return r.List(ctx, params)
}

// buildWhereConditions builds WHERE conditions for queries
func (r *ProductSparePartRepository) buildWhereConditions(params *products.ProductSparePartFilterParams) ([]string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.BrandID != nil {
		conditions = append(conditions, fmt.Sprintf("brand_id = $%d", argIndex))
		args = append(args, *params.BrandID)
		argIndex++
	}

	if params.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", argIndex))
		args = append(args, *params.CategoryID)
		argIndex++
	}

	if params.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *params.IsActive)
		argIndex++
	}

	if params.LowStock != nil && *params.LowStock {
		conditions = append(conditions, "stock_quantity <= min_stock_level")
	}

	if params.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("selling_price >= $%d", argIndex))
		args = append(args, *params.MinPrice)
		argIndex++
	}

	if params.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("selling_price <= $%d", argIndex))
		args = append(args, *params.MaxPrice)
		argIndex++
	}

	if params.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(product_name ILIKE $%d OR product_code ILIKE $%d OR barcode ILIKE $%d)", argIndex, argIndex, argIndex))
		searchTerm := "%" + params.Search + "%"
		args = append(args, searchTerm)
		argIndex++
	}

	return conditions, args
}