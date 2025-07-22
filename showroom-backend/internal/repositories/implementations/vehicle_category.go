package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/user"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

type vehicleCategoryRepository struct {
	db *sql.DB
}

// NewVehicleCategoryRepository creates a new vehicle category repository
func NewVehicleCategoryRepository(db *sql.DB) interfaces.VehicleCategoryRepository {
	return &vehicleCategoryRepository{db: db}
}

// Create creates a new vehicle category
func (r *vehicleCategoryRepository) Create(ctx context.Context, category *master.VehicleCategory) (*master.VehicleCategory, error) {
	query := `
		INSERT INTO vehicle_categories (category_code, category_name, description, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING category_id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		category.CategoryCode, category.CategoryName, category.Description, category.CreatedBy,
	).Scan(&category.CategoryID, &category.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle category: %w", err)
	}

	return category, nil
}

// GetByID retrieves a vehicle category by ID
func (r *vehicleCategoryRepository) GetByID(ctx context.Context, id int) (*master.VehicleCategory, error) {
	query := `
		SELECT c.category_id, c.category_code, c.category_name, c.description,
		       c.is_active, c.created_at, c.created_by,
		       u.user_id as creator_user_id, u.username as creator_username, u.full_name as creator_full_name
		FROM vehicle_categories c
		LEFT JOIN users u ON c.created_by = u.user_id
		WHERE c.category_id = $1 AND c.is_active = true`

	category := &master.VehicleCategory{}
	creator := &user.UserCreatorInfo{}
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&category.CategoryID, &category.CategoryCode, &category.CategoryName, &category.Description,
		&category.IsActive, &category.CreatedAt, &category.CreatedBy,
		&creator.UserID, &creator.Username, &creator.FullName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle category not found")
		}
		return nil, fmt.Errorf("failed to get vehicle category: %w", err)
	}

	category.Creator = creator
	return category, nil
}

// GetByCode retrieves a vehicle category by code
func (r *vehicleCategoryRepository) GetByCode(ctx context.Context, code string) (*master.VehicleCategory, error) {
	query := `
		SELECT c.category_id, c.category_code, c.category_name, c.description,
		       c.is_active, c.created_at, c.created_by
		FROM vehicle_categories c
		WHERE c.category_code = $1 AND c.is_active = true`

	category := &master.VehicleCategory{}
	
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&category.CategoryID, &category.CategoryCode, &category.CategoryName, &category.Description,
		&category.IsActive, &category.CreatedAt, &category.CreatedBy,
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
func (r *vehicleCategoryRepository) Update(ctx context.Context, id int, category *master.VehicleCategory) (*master.VehicleCategory, error) {
	query := `
		UPDATE vehicle_categories 
		SET category_name = $1, description = $2, is_active = $3
		WHERE category_id = $4 AND is_active = true
		RETURNING category_id`

	var categoryID int
	err := r.db.QueryRowContext(ctx, query,
		category.CategoryName, category.Description, category.IsActive, id,
	).Scan(&categoryID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle category not found")
		}
		return nil, fmt.Errorf("failed to update vehicle category: %w", err)
	}

	category.CategoryID = categoryID
	return category, nil
}

// Delete soft deletes a vehicle category
func (r *vehicleCategoryRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE vehicle_categories SET is_active = false WHERE category_id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vehicle category not found")
	}

	return nil
}

// ListActive retrieves all active vehicle categories
func (r *vehicleCategoryRepository) ListActive(ctx context.Context) ([]master.VehicleCategory, error) {
	query := `
		SELECT category_id, category_code, category_name, description, is_active, created_at
		FROM vehicle_categories
		WHERE is_active = true
		ORDER BY category_name ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list vehicle categories: %w", err)
	}
	defer rows.Close()

	var categories []master.VehicleCategory
	for rows.Next() {
		var category master.VehicleCategory
		err := rows.Scan(
			&category.CategoryID, &category.CategoryCode, &category.CategoryName, &category.Description,
			&category.IsActive, &category.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vehicle category: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// List retrieves vehicle categories with optional filtering by active status
func (r *vehicleCategoryRepository) List(ctx context.Context, isActive *bool) ([]master.VehicleCategory, error) {
	query := `
		SELECT category_id, category_code, category_name, description, is_active, created_at
		FROM vehicle_categories`
	
	var args []interface{}
	if isActive != nil {
		query += ` WHERE is_active = $1`
		args = append(args, *isActive)
	}
	
	query += ` ORDER BY category_name ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list vehicle categories: %w", err)
	}
	defer rows.Close()

	var categories []master.VehicleCategory
	for rows.Next() {
		var category master.VehicleCategory
		err := rows.Scan(
			&category.CategoryID, &category.CategoryCode, &category.CategoryName, &category.Description,
			&category.IsActive, &category.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vehicle category: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// ExistsByCode checks if a vehicle category with the given code exists
func (r *vehicleCategoryRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_categories WHERE category_code = $1)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, code).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check vehicle category code existence: %w", err)
	}
	
	return exists, nil
}

// ExistsByName checks if a vehicle category with the given name exists
func (r *vehicleCategoryRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_categories WHERE LOWER(category_name) = LOWER($1) AND is_active = true)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check vehicle category name existence: %w", err)
	}
	
	return exists, nil
}

// ExistsByCodeExcludingID checks if a vehicle category with the given code exists excluding a specific ID
func (r *vehicleCategoryRepository) ExistsByCodeExcludingID(ctx context.Context, code string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_categories WHERE category_code = $1 AND category_id != $2)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, code, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check vehicle category code existence: %w", err)
	}
	
	return exists, nil
}

// ExistsByNameExcludingID checks if a vehicle category with the given name exists excluding a specific ID
func (r *vehicleCategoryRepository) ExistsByNameExcludingID(ctx context.Context, name string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_categories WHERE LOWER(category_name) = LOWER($1) AND category_id != $2 AND is_active = true)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, name, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check vehicle category name existence: %w", err)
	}
	
	return exists, nil
}

// GetNextCategoryCode generates the next vehicle category code
func (r *vehicleCategoryRepository) GetNextCategoryCode(ctx context.Context) (string, error) {
	query := `
		SELECT category_code 
		FROM vehicle_categories 
		WHERE category_code LIKE 'CAT-%' 
		ORDER BY category_id DESC 
		LIMIT 1`
	
	var lastCode string
	err := r.db.QueryRowContext(ctx, query).Scan(&lastCode)
	
	if err != nil {
		if err == sql.ErrNoRows {
			// First category
			return "CAT-001", nil
		}
		return "", fmt.Errorf("failed to get last vehicle category code: %w", err)
	}
	
	// Extract number from code (e.g., "CAT-001" -> 1)
	parts := strings.Split(lastCode, "-")
	if len(parts) != 2 {
		return "CAT-001", nil
	}
	
	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return "CAT-001", nil
	}
	
	// Generate next code
	nextNum := num + 1
	return fmt.Sprintf("CAT-%03d", nextNum), nil
}