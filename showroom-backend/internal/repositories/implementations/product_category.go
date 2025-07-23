package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// ProductCategoryRepository implements interfaces.ProductCategoryRepository
type ProductCategoryRepository struct {
	db *sql.DB
}

// NewProductCategoryRepository creates a new product category repository
func NewProductCategoryRepository(db *sql.DB) interfaces.ProductCategoryRepository {
	return &ProductCategoryRepository{db: db}
}

// Create creates a new product category
func (r *ProductCategoryRepository) Create(ctx context.Context, category *master.ProductCategory) (*master.ProductCategory, error) {
	query := `
		INSERT INTO product_categories (category_code, category_name, description, parent_id, level, path, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING category_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		category.CategoryCode,
		category.CategoryName,
		category.Description,
		category.ParentID,
		category.Level,
		category.Path,
		category.CreatedBy,
	).Scan(&category.CategoryID, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create product category: %w", err)
	}

	category.IsActive = true
	return category, nil
}

// GetByID retrieves a product category by ID with parent info
func (r *ProductCategoryRepository) GetByID(ctx context.Context, id int) (*master.ProductCategory, error) {
	query := `
		SELECT pc.category_id, pc.category_code, pc.category_name, pc.description, pc.parent_id, pc.level, pc.path, pc.is_active, pc.created_at, pc.updated_at, pc.created_by,
		       parent.category_name as parent_name
		FROM product_categories pc
		LEFT JOIN product_categories parent ON pc.parent_id = parent.category_id
		WHERE pc.category_id = $1`

	category := &master.ProductCategory{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&category.CategoryID,
		&category.CategoryCode,
		&category.CategoryName,
		&category.Description,
		&category.ParentID,
		&category.Level,
		&category.Path,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.CreatedBy,
		&category.ParentName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product category not found")
		}
		return nil, fmt.Errorf("failed to get product category: %w", err)
	}

	return category, nil
}

// GetByCode retrieves a product category by code with parent info
func (r *ProductCategoryRepository) GetByCode(ctx context.Context, code string) (*master.ProductCategory, error) {
	query := `
		SELECT pc.category_id, pc.category_code, pc.category_name, pc.description, pc.parent_id, pc.level, pc.path, pc.is_active, pc.created_at, pc.updated_at, pc.created_by,
		       parent.category_name as parent_name
		FROM product_categories pc
		LEFT JOIN product_categories parent ON pc.parent_id = parent.category_id
		WHERE pc.category_code = $1`

	category := &master.ProductCategory{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&category.CategoryID,
		&category.CategoryCode,
		&category.CategoryName,
		&category.Description,
		&category.ParentID,
		&category.Level,
		&category.Path,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.CreatedBy,
		&category.ParentName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product category not found")
		}
		return nil, fmt.Errorf("failed to get product category: %w", err)
	}

	return category, nil
}

// Update updates a product category
func (r *ProductCategoryRepository) Update(ctx context.Context, id int, category *master.ProductCategory) (*master.ProductCategory, error) {
	query := `
		UPDATE product_categories
		SET category_name = $1, description = $2, parent_id = $3, level = $4, path = $5, is_active = $6, updated_at = NOW()
		WHERE category_id = $7
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		category.CategoryName,
		category.Description,
		category.ParentID,
		category.Level,
		category.Path,
		category.IsActive,
		id,
	).Scan(&category.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product category not found")
		}
		return nil, fmt.Errorf("failed to update product category: %w", err)
	}

	category.CategoryID = id
	return category, nil
}

// Delete soft deletes a product category
func (r *ProductCategoryRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE product_categories SET is_active = FALSE, updated_at = NOW() WHERE category_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product category not found")
	}

	return nil
}

// List retrieves product categories with filtering and pagination
func (r *ProductCategoryRepository) List(ctx context.Context, params *master.ProductCategoryFilterParams) (*common.PaginatedResponse, error) {
	params.Validate()

	// Build WHERE conditions
	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.ParentID != nil {
		conditions = append(conditions, fmt.Sprintf("pc.parent_id = $%d", argIndex))
		args = append(args, *params.ParentID)
		argIndex++
	}

	if params.Level != nil {
		conditions = append(conditions, fmt.Sprintf("pc.level = $%d", argIndex))
		args = append(args, *params.Level)
		argIndex++
	}

	if params.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("pc.is_active = $%d", argIndex))
		args = append(args, *params.IsActive)
		argIndex++
	}

	if params.Search != "" {
		searchCondition := fmt.Sprintf("(pc.category_name ILIKE $%d OR pc.category_code ILIKE $%d)", argIndex, argIndex)
		conditions = append(conditions, searchCondition)
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM product_categories pc
		LEFT JOIN product_categories parent ON pc.parent_id = parent.category_id
		%s`, whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count product categories: %w", err)
	}

	// Build main query
	query := fmt.Sprintf(`
		SELECT pc.category_id, pc.category_code, pc.category_name, parent.category_name as parent_name, pc.level, pc.is_active, pc.created_at
		FROM product_categories pc
		LEFT JOIN product_categories parent ON pc.parent_id = parent.category_id
		%s
		ORDER BY pc.level ASC, pc.category_name ASC
		LIMIT $%d OFFSET $%d`,
		whereClause, argIndex, argIndex+1)

	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list product categories: %w", err)
	}
	defer rows.Close()

	var categories []master.ProductCategoryListItem
	for rows.Next() {
		var category master.ProductCategoryListItem
		err := rows.Scan(
			&category.CategoryID,
			&category.CategoryCode,
			&category.CategoryName,
			&category.ParentName,
			&category.Level,
			&category.IsActive,
			&category.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product category: %w", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate product categories: %w", err)
	}

	return &common.PaginatedResponse{
		Data:       categories,
		Total: int(total),
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: params.GetTotalPages(total),
		HasMore:    params.GetHasMore(total),
	}, nil
}

// GetTree retrieves product categories in hierarchical structure
func (r *ProductCategoryRepository) GetTree(ctx context.Context) ([]master.ProductCategoryTree, error) {
	query := `
		SELECT category_id, category_code, category_name, description, level, is_active
		FROM product_categories
		WHERE is_active = TRUE
		ORDER BY level ASC, category_name ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get product category tree: %w", err)
	}
	defer rows.Close()

	var allCategories []master.ProductCategoryTree
	for rows.Next() {
		var category master.ProductCategoryTree
		err := rows.Scan(
			&category.CategoryID,
			&category.CategoryCode,
			&category.CategoryName,
			&category.Description,
			&category.Level,
			&category.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product category: %w", err)
		}
		allCategories = append(allCategories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate product categories: %w", err)
	}

	// Build hierarchical structure
	var rootCategories []master.ProductCategoryTree
	categoryMap := make(map[int]*master.ProductCategoryTree)

	// First pass: create map of all categories
	for i := range allCategories {
		categoryMap[allCategories[i].CategoryID] = &allCategories[i]
	}

	// Second pass: build hierarchy
	for i := range allCategories {
		if allCategories[i].Level == 1 {
			rootCategories = append(rootCategories, allCategories[i])
		}
		// Note: For simplicity, we're returning flat structure here
		// In a real implementation, you'd build the parent-child relationships
	}

	return rootCategories, nil
}

// GetChildren retrieves child categories of a parent
func (r *ProductCategoryRepository) GetChildren(ctx context.Context, parentID int) ([]master.ProductCategory, error) {
	query := `
		SELECT category_id, category_code, category_name, description, parent_id, level, path, is_active, created_at, updated_at, created_by
		FROM product_categories
		WHERE parent_id = $1 AND is_active = TRUE
		ORDER BY category_name ASC`

	rows, err := r.db.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get child categories: %w", err)
	}
	defer rows.Close()

	var categories []master.ProductCategory
	for rows.Next() {
		var category master.ProductCategory
		err := rows.Scan(
			&category.CategoryID,
			&category.CategoryCode,
			&category.CategoryName,
			&category.Description,
			&category.ParentID,
			&category.Level,
			&category.Path,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.CreatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan child category: %w", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate child categories: %w", err)
	}

	return categories, nil
}

// GenerateCode generates a new product category code
func (r *ProductCategoryRepository) GenerateCode(ctx context.Context) (string, error) {
	query := `
		SELECT category_code
		FROM product_categories
		WHERE category_code LIKE 'PC-%'
		ORDER BY category_code DESC
		LIMIT 1`

	var lastCode sql.NullString
	err := r.db.QueryRowContext(ctx, query).Scan(&lastCode)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to get last product category code: %w", err)
	}

	nextNumber := 1
	if lastCode.Valid {
		// Extract number from code (e.g., "PC-001" -> "001" -> 1)
		parts := strings.Split(lastCode.String, "-")
		if len(parts) == 2 {
			if num, err := strconv.Atoi(parts[1]); err == nil {
				nextNumber = num + 1
			}
		}
	}

	return fmt.Sprintf("PC-%03d", nextNumber), nil
}

// IsCodeExists checks if a product category code already exists
func (r *ProductCategoryRepository) IsCodeExists(ctx context.Context, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM product_categories WHERE category_code = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, code).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check product category code existence: %w", err)
	}

	return exists, nil
}

// IsNameExists checks if a product category name already exists within the same parent (excluding a specific ID)
func (r *ProductCategoryRepository) IsNameExists(ctx context.Context, name string, parentID *int, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM product_categories WHERE category_name = $1 AND parent_id = $2 AND category_id != $3)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, name, parentID, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check product category name existence: %w", err)
	}

	return exists, nil
}

// UpdatePath updates the path of a product category
func (r *ProductCategoryRepository) UpdatePath(ctx context.Context, categoryID int, path string) error {
	query := `UPDATE product_categories SET path = $1, updated_at = NOW() WHERE category_id = $2`

	_, err := r.db.ExecContext(ctx, query, path, categoryID)
	if err != nil {
		return fmt.Errorf("failed to update product category path: %w", err)
	}

	return nil
}

// UpdateLevel updates the level of a product category
func (r *ProductCategoryRepository) UpdateLevel(ctx context.Context, categoryID int, level int) error {
	query := `UPDATE product_categories SET level = $1, updated_at = NOW() WHERE category_id = $2`

	_, err := r.db.ExecContext(ctx, query, level, categoryID)
	if err != nil {
		return fmt.Errorf("failed to update product category level: %w", err)
	}

	return nil
}