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

// VehicleBrandRepository implements interfaces.VehicleBrandRepository
type VehicleBrandRepository struct {
	db *sql.DB
}

// NewVehicleBrandRepository creates a new vehicle brand repository
func NewVehicleBrandRepository(db *sql.DB) interfaces.VehicleBrandRepository {
	return &VehicleBrandRepository{db: db}
}

// Create creates a new vehicle brand
func (r *VehicleBrandRepository) Create(ctx context.Context, brand *master.VehicleBrand) (*master.VehicleBrand, error) {
	query := `
		INSERT INTO vehicle_brands (brand_code, brand_name, country_origin, description, logo_url, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING brand_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		brand.BrandCode,
		brand.BrandName,
		brand.CountryOrigin,
		brand.Description,
		brand.LogoURL,
		brand.CreatedBy,
	).Scan(&brand.BrandID, &brand.CreatedAt, &brand.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle brand: %w", err)
	}

	brand.IsActive = true
	return brand, nil
}

// GetByID retrieves a vehicle brand by ID
func (r *VehicleBrandRepository) GetByID(ctx context.Context, id int) (*master.VehicleBrand, error) {
	query := `
		SELECT brand_id, brand_code, brand_name, country_origin, description, logo_url, is_active, created_at, updated_at, created_by
		FROM vehicle_brands
		WHERE brand_id = $1`

	brand := &master.VehicleBrand{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&brand.BrandID,
		&brand.BrandCode,
		&brand.BrandName,
		&brand.CountryOrigin,
		&brand.Description,
		&brand.LogoURL,
		&brand.IsActive,
		&brand.CreatedAt,
		&brand.UpdatedAt,
		&brand.CreatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle brand not found")
		}
		return nil, fmt.Errorf("failed to get vehicle brand: %w", err)
	}

	return brand, nil
}

// GetByCode retrieves a vehicle brand by code
func (r *VehicleBrandRepository) GetByCode(ctx context.Context, code string) (*master.VehicleBrand, error) {
	query := `
		SELECT brand_id, brand_code, brand_name, country_origin, description, logo_url, is_active, created_at, updated_at, created_by
		FROM vehicle_brands
		WHERE brand_code = $1`

	brand := &master.VehicleBrand{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&brand.BrandID,
		&brand.BrandCode,
		&brand.BrandName,
		&brand.CountryOrigin,
		&brand.Description,
		&brand.LogoURL,
		&brand.IsActive,
		&brand.CreatedAt,
		&brand.UpdatedAt,
		&brand.CreatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle brand not found")
		}
		return nil, fmt.Errorf("failed to get vehicle brand: %w", err)
	}

	return brand, nil
}

// Update updates a vehicle brand
func (r *VehicleBrandRepository) Update(ctx context.Context, id int, brand *master.VehicleBrand) (*master.VehicleBrand, error) {
	query := `
		UPDATE vehicle_brands
		SET brand_name = $1, country_origin = $2, description = $3, logo_url = $4, is_active = $5, updated_at = NOW()
		WHERE brand_id = $6
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		brand.BrandName,
		brand.CountryOrigin,
		brand.Description,
		brand.LogoURL,
		brand.IsActive,
		id,
	).Scan(&brand.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle brand not found")
		}
		return nil, fmt.Errorf("failed to update vehicle brand: %w", err)
	}

	brand.BrandID = id
	return brand, nil
}

// Delete soft deletes a vehicle brand
func (r *VehicleBrandRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE vehicle_brands SET is_active = FALSE, updated_at = NOW() WHERE brand_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle brand: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vehicle brand not found")
	}

	return nil
}

// List retrieves vehicle brands with filtering and pagination
func (r *VehicleBrandRepository) List(ctx context.Context, params *master.VehicleBrandFilterParams) (*common.PaginatedResponse, error) {
	params.Validate()

	// Build WHERE conditions
	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *params.IsActive)
		argIndex++
	}

	if params.CountryOrigin != "" {
		conditions = append(conditions, fmt.Sprintf("country_origin ILIKE $%d", argIndex))
		args = append(args, "%"+params.CountryOrigin+"%")
		argIndex++
	}

	if params.Search != "" {
		searchCondition := fmt.Sprintf("(brand_name ILIKE $%d OR brand_code ILIKE $%d)", argIndex, argIndex)
		conditions = append(conditions, searchCondition)
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM vehicle_brands %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count vehicle brands: %w", err)
	}

	// Build main query
	query := fmt.Sprintf(`
		SELECT brand_id, brand_code, brand_name, country_origin, is_active, created_at
		FROM vehicle_brands
		%s
		ORDER BY brand_name ASC
		LIMIT $%d OFFSET $%d`,
		whereClause, argIndex, argIndex+1)

	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list vehicle brands: %w", err)
	}
	defer rows.Close()

	var brands []master.VehicleBrandListItem
	for rows.Next() {
		var brand master.VehicleBrandListItem
		err := rows.Scan(
			&brand.BrandID,
			&brand.BrandCode,
			&brand.BrandName,
			&brand.CountryOrigin,
			&brand.IsActive,
			&brand.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vehicle brand: %w", err)
		}
		brands = append(brands, brand)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate vehicle brands: %w", err)
	}

	return &common.PaginatedResponse{
		Data:       brands,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: params.GetTotalPages(total),
		HasMore:    params.GetHasMore(total),
	}, nil
}

// GenerateCode generates a new vehicle brand code
func (r *VehicleBrandRepository) GenerateCode(ctx context.Context) (string, error) {
	query := `
		SELECT brand_code
		FROM vehicle_brands
		WHERE brand_code LIKE 'VB-%'
		ORDER BY brand_code DESC
		LIMIT 1`

	var lastCode sql.NullString
	err := r.db.QueryRowContext(ctx, query).Scan(&lastCode)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to get last vehicle brand code: %w", err)
	}

	nextNumber := 1
	if lastCode.Valid {
		// Extract number from code (e.g., "VB-001" -> "001" -> 1)
		parts := strings.Split(lastCode.String, "-")
		if len(parts) == 2 {
			if num, err := strconv.Atoi(parts[1]); err == nil {
				nextNumber = num + 1
			}
		}
	}

	return fmt.Sprintf("VB-%03d", nextNumber), nil
}

// IsCodeExists checks if a vehicle brand code already exists
func (r *VehicleBrandRepository) IsCodeExists(ctx context.Context, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_brands WHERE brand_code = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, code).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check vehicle brand code existence: %w", err)
	}

	return exists, nil
}

// IsNameExists checks if a vehicle brand name already exists (excluding a specific ID)
func (r *VehicleBrandRepository) IsNameExists(ctx context.Context, name string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_brands WHERE brand_name = $1 AND brand_id != $2)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, name, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check vehicle brand name existence: %w", err)
	}

	return exists, nil
}

// VehicleCategoryRepository implements interfaces.VehicleCategoryRepository
type VehicleCategoryRepository struct {
	db *sql.DB
}

// NewVehicleCategoryRepository creates a new vehicle category repository
func NewVehicleCategoryRepository(db *sql.DB) interfaces.VehicleCategoryRepository {
	return &VehicleCategoryRepository{db: db}
}

// Create creates a new vehicle category
func (r *VehicleCategoryRepository) Create(ctx context.Context, category *master.VehicleCategory) (*master.VehicleCategory, error) {
	query := `
		INSERT INTO vehicle_categories (category_code, category_name, description, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING category_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		category.CategoryCode,
		category.CategoryName,
		category.Description,
		category.CreatedBy,
	).Scan(&category.CategoryID, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle category: %w", err)
	}

	category.IsActive = true
	return category, nil
}

// GetByID retrieves a vehicle category by ID
func (r *VehicleCategoryRepository) GetByID(ctx context.Context, id int) (*master.VehicleCategory, error) {
	query := `
		SELECT category_id, category_code, category_name, description, is_active, created_at, updated_at, created_by
		FROM vehicle_categories
		WHERE category_id = $1`

	category := &master.VehicleCategory{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&category.CategoryID,
		&category.CategoryCode,
		&category.CategoryName,
		&category.Description,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.CreatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle category not found")
		}
		return nil, fmt.Errorf("failed to get vehicle category: %w", err)
	}

	return category, nil
}

// GetByCode retrieves a vehicle category by code
func (r *VehicleCategoryRepository) GetByCode(ctx context.Context, code string) (*master.VehicleCategory, error) {
	query := `
		SELECT category_id, category_code, category_name, description, is_active, created_at, updated_at, created_by
		FROM vehicle_categories
		WHERE category_code = $1`

	category := &master.VehicleCategory{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&category.CategoryID,
		&category.CategoryCode,
		&category.CategoryName,
		&category.Description,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.CreatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle category not found")
		}
		return nil, fmt.Errorf("failed to get vehicle category: %w", err)
	}

	return category, nil
}

// Update updates a vehicle category
func (r *VehicleCategoryRepository) Update(ctx context.Context, id int, category *master.VehicleCategory) (*master.VehicleCategory, error) {
	query := `
		UPDATE vehicle_categories
		SET category_name = $1, description = $2, is_active = $3, updated_at = NOW()
		WHERE category_id = $4
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		category.CategoryName,
		category.Description,
		category.IsActive,
		id,
	).Scan(&category.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle category not found")
		}
		return nil, fmt.Errorf("failed to update vehicle category: %w", err)
	}

	category.CategoryID = id
	return category, nil
}

// Delete soft deletes a vehicle category
func (r *VehicleCategoryRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE vehicle_categories SET is_active = FALSE, updated_at = NOW() WHERE category_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vehicle category not found")
	}

	return nil
}

// List retrieves vehicle categories with filtering and pagination
func (r *VehicleCategoryRepository) List(ctx context.Context, params *master.VehicleCategoryFilterParams) (*common.PaginatedResponse, error) {
	params.Validate()

	// Build WHERE conditions
	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *params.IsActive)
		argIndex++
	}

	if params.Search != "" {
		searchCondition := fmt.Sprintf("(category_name ILIKE $%d OR category_code ILIKE $%d)", argIndex, argIndex)
		conditions = append(conditions, searchCondition)
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM vehicle_categories %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count vehicle categories: %w", err)
	}

	// Build main query
	query := fmt.Sprintf(`
		SELECT category_id, category_code, category_name, is_active, created_at
		FROM vehicle_categories
		%s
		ORDER BY category_name ASC
		LIMIT $%d OFFSET $%d`,
		whereClause, argIndex, argIndex+1)

	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list vehicle categories: %w", err)
	}
	defer rows.Close()

	var categories []master.VehicleCategoryListItem
	for rows.Next() {
		var category master.VehicleCategoryListItem
		err := rows.Scan(
			&category.CategoryID,
			&category.CategoryCode,
			&category.CategoryName,
			&category.IsActive,
			&category.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vehicle category: %w", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate vehicle categories: %w", err)
	}

	return &common.PaginatedResponse{
		Data:       categories,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: params.GetTotalPages(total),
		HasMore:    params.GetHasMore(total),
	}, nil
}

// GenerateCode generates a new vehicle category code
func (r *VehicleCategoryRepository) GenerateCode(ctx context.Context) (string, error) {
	query := `
		SELECT category_code
		FROM vehicle_categories
		WHERE category_code LIKE 'VC-%'
		ORDER BY category_code DESC
		LIMIT 1`

	var lastCode sql.NullString
	err := r.db.QueryRowContext(ctx, query).Scan(&lastCode)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to get last vehicle category code: %w", err)
	}

	nextNumber := 1
	if lastCode.Valid {
		// Extract number from code (e.g., "VC-001" -> "001" -> 1)
		parts := strings.Split(lastCode.String, "-")
		if len(parts) == 2 {
			if num, err := strconv.Atoi(parts[1]); err == nil {
				nextNumber = num + 1
			}
		}
	}

	return fmt.Sprintf("VC-%03d", nextNumber), nil
}

// IsCodeExists checks if a vehicle category code already exists
func (r *VehicleCategoryRepository) IsCodeExists(ctx context.Context, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_categories WHERE category_code = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, code).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check vehicle category code existence: %w", err)
	}

	return exists, nil
}

// IsNameExists checks if a vehicle category name already exists (excluding a specific ID)
func (r *VehicleCategoryRepository) IsNameExists(ctx context.Context, name string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_categories WHERE category_name = $1 AND category_id != $2)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, name, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check vehicle category name existence: %w", err)
	}

	return exists, nil
}