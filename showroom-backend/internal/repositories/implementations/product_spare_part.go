package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

type productSparePartRepository struct {
	db *sql.DB
}

// NewProductSparePartRepository creates a new product spare part repository
func NewProductSparePartRepository(db *sql.DB) interfaces.ProductSparePartRepository {
	return &productSparePartRepository{db: db}
}

// Create creates a new product spare part
func (r *productSparePartRepository) Create(ctx context.Context, product *inventory.ProductSparePart) (*inventory.ProductSparePart, error) {
	query := `
		INSERT INTO products_spare_parts (
			product_code, product_name, description, brand_id, category_id, unit_measure,
			cost_price, selling_price, markup_percentage, stock_quantity, min_stock_level,
			max_stock_level, location_rack, barcode, weight, dimensions, created_by,
			is_active, product_image, notes
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
		RETURNING product_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		product.ProductCode, product.ProductName, product.Description, product.BrandID,
		product.CategoryID, product.UnitMeasure, product.CostPrice, product.SellingPrice,
		product.MarkupPercentage, product.StockQuantity, product.MinStockLevel,
		product.MaxStockLevel, product.LocationRack, product.Barcode, product.Weight,
		product.Dimensions, product.CreatedBy, product.IsActive, product.ProductImage,
		product.Notes,
	).Scan(&product.ProductID, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

// GetByID retrieves a product by ID
func (r *productSparePartRepository) GetByID(ctx context.Context, id int) (*inventory.ProductSparePart, error) {
	query := `
		SELECT p.product_id, p.product_code, p.product_name, p.description, p.brand_id,
		       p.category_id, p.unit_measure, p.cost_price, p.selling_price, p.markup_percentage,
		       p.stock_quantity, p.min_stock_level, p.max_stock_level, p.location_rack,
		       p.barcode, p.weight, p.dimensions, p.created_at, p.updated_at, p.created_by,
		       p.is_active, p.product_image, p.notes,
		       vb.brand_name, pc.category_name
		FROM products_spare_parts p
		LEFT JOIN vehicle_brands vb ON p.brand_id = vb.brand_id
		LEFT JOIN product_categories pc ON p.category_id = pc.category_id
		WHERE p.product_id = $1`

	product := &inventory.ProductSparePart{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ProductID, &product.ProductCode, &product.ProductName, &product.Description,
		&product.BrandID, &product.CategoryID, &product.UnitMeasure, &product.CostPrice,
		&product.SellingPrice, &product.MarkupPercentage, &product.StockQuantity,
		&product.MinStockLevel, &product.MaxStockLevel, &product.LocationRack,
		&product.Barcode, &product.Weight, &product.Dimensions, &product.CreatedAt,
		&product.UpdatedAt, &product.CreatedBy, &product.IsActive, &product.ProductImage,
		&product.Notes, &product.BrandName, &product.CategoryName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// GetByCode retrieves a product by code
func (r *productSparePartRepository) GetByCode(ctx context.Context, code string) (*inventory.ProductSparePart, error) {
	query := `
		SELECT p.product_id, p.product_code, p.product_name, p.description, p.brand_id,
		       p.category_id, p.unit_measure, p.cost_price, p.selling_price, p.markup_percentage,
		       p.stock_quantity, p.min_stock_level, p.max_stock_level, p.location_rack,
		       p.barcode, p.weight, p.dimensions, p.created_at, p.updated_at, p.created_by,
		       p.is_active, p.product_image, p.notes,
		       vb.brand_name, pc.category_name
		FROM products_spare_parts p
		LEFT JOIN vehicle_brands vb ON p.brand_id = vb.brand_id
		LEFT JOIN product_categories pc ON p.category_id = pc.category_id
		WHERE p.product_code = $1`

	product := &inventory.ProductSparePart{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&product.ProductID, &product.ProductCode, &product.ProductName, &product.Description,
		&product.BrandID, &product.CategoryID, &product.UnitMeasure, &product.CostPrice,
		&product.SellingPrice, &product.MarkupPercentage, &product.StockQuantity,
		&product.MinStockLevel, &product.MaxStockLevel, &product.LocationRack,
		&product.Barcode, &product.Weight, &product.Dimensions, &product.CreatedAt,
		&product.UpdatedAt, &product.CreatedBy, &product.IsActive, &product.ProductImage,
		&product.Notes, &product.BrandName, &product.CategoryName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// GetByBarcode retrieves a product by barcode
func (r *productSparePartRepository) GetByBarcode(ctx context.Context, barcode string) (*inventory.ProductSparePart, error) {
	query := `
		SELECT p.product_id, p.product_code, p.product_name, p.description, p.brand_id,
		       p.category_id, p.unit_measure, p.cost_price, p.selling_price, p.markup_percentage,
		       p.stock_quantity, p.min_stock_level, p.max_stock_level, p.location_rack,
		       p.barcode, p.weight, p.dimensions, p.created_at, p.updated_at, p.created_by,
		       p.is_active, p.product_image, p.notes,
		       vb.brand_name, pc.category_name
		FROM products_spare_parts p
		LEFT JOIN vehicle_brands vb ON p.brand_id = vb.brand_id
		LEFT JOIN product_categories pc ON p.category_id = pc.category_id
		WHERE p.barcode = $1`

	product := &inventory.ProductSparePart{}
	err := r.db.QueryRowContext(ctx, query, barcode).Scan(
		&product.ProductID, &product.ProductCode, &product.ProductName, &product.Description,
		&product.BrandID, &product.CategoryID, &product.UnitMeasure, &product.CostPrice,
		&product.SellingPrice, &product.MarkupPercentage, &product.StockQuantity,
		&product.MinStockLevel, &product.MaxStockLevel, &product.LocationRack,
		&product.Barcode, &product.Weight, &product.Dimensions, &product.CreatedAt,
		&product.UpdatedAt, &product.CreatedBy, &product.IsActive, &product.ProductImage,
		&product.Notes, &product.BrandName, &product.CategoryName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// Update updates a product
func (r *productSparePartRepository) Update(ctx context.Context, id int, product *inventory.ProductSparePart) (*inventory.ProductSparePart, error) {
	query := `
		UPDATE products_spare_parts SET
			product_name = $2, description = $3, brand_id = $4, category_id = $5,
			unit_measure = $6, cost_price = $7, selling_price = $8, markup_percentage = $9,
			stock_quantity = $10, min_stock_level = $11, max_stock_level = $12,
			location_rack = $13, barcode = $14, weight = $15, dimensions = $16,
			is_active = $17, product_image = $18, notes = $19, updated_at = NOW()
		WHERE product_id = $1
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query, id,
		product.ProductName, product.Description, product.BrandID, product.CategoryID,
		product.UnitMeasure, product.CostPrice, product.SellingPrice, product.MarkupPercentage,
		product.StockQuantity, product.MinStockLevel, product.MaxStockLevel,
		product.LocationRack, product.Barcode, product.Weight, product.Dimensions,
		product.IsActive, product.ProductImage, product.Notes,
	).Scan(&product.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	product.ProductID = id
	return product, nil
}

// Delete deletes a product
func (r *productSparePartRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM products_spare_parts WHERE product_id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// List retrieves products with filtering and pagination
func (r *productSparePartRepository) List(ctx context.Context, params *inventory.ProductSparePartFilterParams) ([]inventory.ProductSparePartListItem, int, error) {
	whereConditions := []string{"1 = 1"}
	args := []interface{}{}
	argIndex := 1

	// Build WHERE conditions
	if params.BrandID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.brand_id = $%d", argIndex))
		args = append(args, *params.BrandID)
		argIndex++
	}

	if params.CategoryID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.category_id = $%d", argIndex))
		args = append(args, *params.CategoryID)
		argIndex++
	}

	if params.IsActive != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.is_active = $%d", argIndex))
		args = append(args, *params.IsActive)
		argIndex++
	}

	if params.LocationRack != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("p.location_rack ILIKE $%d", argIndex))
		args = append(args, "%"+params.LocationRack+"%")
		argIndex++
	}

	if params.Barcode != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("p.barcode = $%d", argIndex))
		args = append(args, params.Barcode)
		argIndex++
	}

	if params.LowStock != nil && *params.LowStock {
		whereConditions = append(whereConditions, "p.stock_quantity <= p.min_stock_level")
	}

	if params.MinPrice != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.selling_price >= $%d", argIndex))
		args = append(args, *params.MinPrice)
		argIndex++
	}

	if params.MaxPrice != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.selling_price <= $%d", argIndex))
		args = append(args, *params.MaxPrice)
		argIndex++
	}

	if params.Search != "" {
		whereConditions = append(whereConditions, fmt.Sprintf(
			"(p.product_name ILIKE $%d OR p.product_code ILIKE $%d OR p.description ILIKE $%d)",
			argIndex, argIndex, argIndex))
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM products_spare_parts p
		LEFT JOIN vehicle_brands vb ON p.brand_id = vb.brand_id
		LEFT JOIN product_categories pc ON p.category_id = pc.category_id
		WHERE %s`, whereClause)

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Main query with pagination
	params.PaginationParams.Validate()
	offset := params.PaginationParams.GetOffset()
	limit := params.PaginationParams.Limit

	query := fmt.Sprintf(`
		SELECT p.product_id, p.product_code, p.product_name, vb.brand_name, pc.category_name,
		       p.unit_measure, p.cost_price, p.selling_price, p.stock_quantity,
		       p.min_stock_level, p.location_rack, p.is_active, p.created_at
		FROM products_spare_parts p
		LEFT JOIN vehicle_brands vb ON p.brand_id = vb.brand_id
		LEFT JOIN product_categories pc ON p.category_id = pc.category_id
		WHERE %s
		ORDER BY p.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []inventory.ProductSparePartListItem
	for rows.Next() {
		var product inventory.ProductSparePartListItem
		err := rows.Scan(
			&product.ProductID, &product.ProductCode, &product.ProductName,
			&product.BrandName, &product.CategoryName, &product.UnitMeasure,
			&product.CostPrice, &product.SellingPrice, &product.StockQuantity,
			&product.MinStockLevel, &product.LocationRack, &product.IsActive,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, total, nil
}

// GetLowStockProducts retrieves products below minimum stock level
func (r *productSparePartRepository) GetLowStockProducts(ctx context.Context, page, limit int) ([]inventory.ProductSparePartListItem, int, error) {
	params := &inventory.ProductSparePartFilterParams{
		LowStock: func() *bool { b := true; return &b }(),
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// GetByBrand retrieves products by brand
func (r *productSparePartRepository) GetByBrand(ctx context.Context, brandID int, page, limit int) ([]inventory.ProductSparePartListItem, int, error) {
	params := &inventory.ProductSparePartFilterParams{
		BrandID: &brandID,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// GetByCategory retrieves products by category
func (r *productSparePartRepository) GetByCategory(ctx context.Context, categoryID int, page, limit int) ([]inventory.ProductSparePartListItem, int, error) {
	params := &inventory.ProductSparePartFilterParams{
		CategoryID: &categoryID,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// GetByLocation retrieves products by location
func (r *productSparePartRepository) GetByLocation(ctx context.Context, location string, page, limit int) ([]inventory.ProductSparePartListItem, int, error) {
	params := &inventory.ProductSparePartFilterParams{
		LocationRack: location,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// Search searches products by query
func (r *productSparePartRepository) Search(ctx context.Context, query string, page, limit int) ([]inventory.ProductSparePartListItem, int, error) {
	params := &inventory.ProductSparePartFilterParams{
		Search: query,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// ExistsByCode checks if a product exists by code
func (r *productSparePartRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM products_spare_parts WHERE product_code = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, code).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check product existence: %w", err)
	}
	return exists, nil
}

// ExistsByBarcode checks if a product exists by barcode
func (r *productSparePartRepository) ExistsByBarcode(ctx context.Context, barcode string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM products_spare_parts WHERE barcode = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, barcode).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check product existence: %w", err)
	}
	return exists, nil
}

// ExistsByCodeExcludingID checks if a product exists by code excluding a specific ID
func (r *productSparePartRepository) ExistsByCodeExcludingID(ctx context.Context, code string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM products_spare_parts WHERE product_code = $1 AND product_id != $2)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, code, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check product existence: %w", err)
	}
	return exists, nil
}

// ExistsByBarcodeExcludingID checks if a product exists by barcode excluding a specific ID
func (r *productSparePartRepository) ExistsByBarcodeExcludingID(ctx context.Context, barcode string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM products_spare_parts WHERE barcode = $1 AND product_id != $2)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, barcode, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check product existence: %w", err)
	}
	return exists, nil
}

// GetLastProductID gets the last product ID for code generation
func (r *productSparePartRepository) GetLastProductID(ctx context.Context) (int, error) {
	query := `SELECT COALESCE(MAX(product_id), 0) FROM products_spare_parts`
	var lastID int
	err := r.db.QueryRowContext(ctx, query).Scan(&lastID)
	if err != nil {
		return 0, fmt.Errorf("failed to get last product ID: %w", err)
	}
	return lastID, nil
}

// UpdateStock updates product stock quantity
func (r *productSparePartRepository) UpdateStock(ctx context.Context, productID int, newQuantity int) error {
	query := `UPDATE products_spare_parts SET stock_quantity = $2, updated_at = NOW() WHERE product_id = $1`
	result, err := r.db.ExecContext(ctx, query, productID, newQuantity)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// GetStockQuantity gets current stock quantity
func (r *productSparePartRepository) GetStockQuantity(ctx context.Context, productID int) (int, error) {
	query := `SELECT stock_quantity FROM products_spare_parts WHERE product_id = $1`
	var quantity int
	err := r.db.QueryRowContext(ctx, query, productID).Scan(&quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("product not found")
		}
		return 0, fmt.Errorf("failed to get stock quantity: %w", err)
	}
	return quantity, nil
}

// GetInventoryValue calculates total inventory value
func (r *productSparePartRepository) GetInventoryValue(ctx context.Context) (float64, error) {
	query := `SELECT COALESCE(SUM(stock_quantity * cost_price), 0) FROM products_spare_parts WHERE is_active = true`
	var value float64
	err := r.db.QueryRowContext(ctx, query).Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate inventory value: %w", err)
	}
	return value, nil
}

// GetProductsByPriceRange retrieves products within price range
func (r *productSparePartRepository) GetProductsByPriceRange(ctx context.Context, minPrice, maxPrice float64, page, limit int) ([]inventory.ProductSparePartListItem, int, error) {
	params := &inventory.ProductSparePartFilterParams{
		MinPrice: &minPrice,
		MaxPrice: &maxPrice,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}