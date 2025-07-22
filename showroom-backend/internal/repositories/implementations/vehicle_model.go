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

type vehicleModelRepository struct {
	db *sql.DB
}

// NewVehicleModelRepository creates a new vehicle model repository
func NewVehicleModelRepository(db *sql.DB) interfaces.VehicleModelRepository {
	return &vehicleModelRepository{db: db}
}

// Create creates a new vehicle model
func (r *vehicleModelRepository) Create(ctx context.Context, model *master.VehicleModel) (*master.VehicleModel, error) {
	query := `
		INSERT INTO vehicle_models (brand_id, category_id, model_code, model_name, year_start, year_end,
		                           fuel_type, transmission, engine_capacity, specifications_json, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING model_id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		model.BrandID, model.CategoryID, model.ModelCode, model.ModelName, model.YearStart,
		model.YearEnd, model.FuelType, model.Transmission, model.EngineCapacity,
		model.SpecificationsJSON, model.CreatedBy,
	).Scan(&model.ModelID, &model.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle model: %w", err)
	}

	return model, nil
}

// GetByID retrieves a vehicle model by ID
func (r *vehicleModelRepository) GetByID(ctx context.Context, id int) (*master.VehicleModel, error) {
	query := `
		SELECT m.model_id, m.brand_id, m.category_id, m.model_code, m.model_name, m.year_start, m.year_end,
		       m.fuel_type, m.transmission, m.engine_capacity, m.specifications_json, m.is_active,
		       m.created_at, m.created_by,
		       u.user_id as creator_user_id, u.username as creator_username, u.full_name as creator_full_name,
		       b.brand_id as brand_id, b.brand_name, b.country_origin,
		       c.category_id as category_id, c.category_name
		FROM vehicle_models m
		LEFT JOIN users u ON m.created_by = u.user_id
		LEFT JOIN vehicle_brands b ON m.brand_id = b.brand_id
		LEFT JOIN vehicle_categories c ON m.category_id = c.category_id
		WHERE m.model_id = $1 AND m.is_active = true`

	model := &master.VehicleModel{}
	creator := &user.UserCreatorInfo{}
	brand := &master.VehicleBrandInfo{}
	category := &master.VehicleCategoryInfo{}
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&model.ModelID, &model.BrandID, &model.CategoryID, &model.ModelCode, &model.ModelName,
		&model.YearStart, &model.YearEnd, &model.FuelType, &model.Transmission,
		&model.EngineCapacity, &model.SpecificationsJSON, &model.IsActive, &model.CreatedAt,
		&model.CreatedBy,
		&creator.UserID, &creator.Username, &creator.FullName,
		&brand.BrandID, &brand.BrandName, &brand.CountryOrigin,
		&category.CategoryID, &category.CategoryName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle model not found")
		}
		return nil, fmt.Errorf("failed to get vehicle model: %w", err)
	}

	model.Creator = creator
	model.Brand = brand
	model.Category = category
	return model, nil
}

// GetByCode retrieves a vehicle model by code
func (r *vehicleModelRepository) GetByCode(ctx context.Context, code string) (*master.VehicleModel, error) {
	query := `
		SELECT m.model_id, m.brand_id, m.category_id, m.model_code, m.model_name, m.year_start, m.year_end,
		       m.fuel_type, m.transmission, m.engine_capacity, m.specifications_json, m.is_active,
		       m.created_at, m.created_by
		FROM vehicle_models m
		WHERE m.model_code = $1 AND m.is_active = true`

	model := &master.VehicleModel{}
	
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&model.ModelID, &model.BrandID, &model.CategoryID, &model.ModelCode, &model.ModelName,
		&model.YearStart, &model.YearEnd, &model.FuelType, &model.Transmission,
		&model.EngineCapacity, &model.SpecificationsJSON, &model.IsActive, &model.CreatedAt,
		&model.CreatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle model not found")
		}
		return nil, fmt.Errorf("failed to get vehicle model: %w", err)
	}

	return model, nil
}

// Update updates a vehicle model
func (r *vehicleModelRepository) Update(ctx context.Context, id int, model *master.VehicleModel) (*master.VehicleModel, error) {
	query := `
		UPDATE vehicle_models 
		SET brand_id = $1, category_id = $2, model_name = $3, year_start = $4, year_end = $5,
		    fuel_type = $6, transmission = $7, engine_capacity = $8, specifications_json = $9,
		    is_active = $10
		WHERE model_id = $11 AND is_active = true
		RETURNING model_id`

	var modelID int
	err := r.db.QueryRowContext(ctx, query,
		model.BrandID, model.CategoryID, model.ModelName, model.YearStart, model.YearEnd,
		model.FuelType, model.Transmission, model.EngineCapacity, model.SpecificationsJSON,
		model.IsActive, id,
	).Scan(&modelID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle model not found")
		}
		return nil, fmt.Errorf("failed to update vehicle model: %w", err)
	}

	model.ModelID = modelID
	return model, nil
}

// Delete soft deletes a vehicle model
func (r *vehicleModelRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE vehicle_models SET is_active = false WHERE model_id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle model: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vehicle model not found")
	}

	return nil
}

// List retrieves vehicle models with filtering and pagination
func (r *vehicleModelRepository) List(ctx context.Context, params *master.VehicleModelFilterParams) ([]master.VehicleModelListItem, int, error) {
	params.Validate()

	// Build WHERE clause
	var conditions []string
	var args []interface{}
	argIndex := 1

	conditions = append(conditions, "m.is_active = true")

	if params.BrandID != nil {
		conditions = append(conditions, fmt.Sprintf("m.brand_id = $%d", argIndex))
		args = append(args, *params.BrandID)
		argIndex++
	}

	if params.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("m.category_id = $%d", argIndex))
		args = append(args, *params.CategoryID)
		argIndex++
	}

	if params.FuelType != nil {
		conditions = append(conditions, fmt.Sprintf("m.fuel_type = $%d", argIndex))
		args = append(args, *params.FuelType)
		argIndex++
	}

	if params.Transmission != nil {
		conditions = append(conditions, fmt.Sprintf("m.transmission = $%d", argIndex))
		args = append(args, *params.Transmission)
		argIndex++
	}

	if params.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("m.is_active = $%d", argIndex))
		args = append(args, *params.IsActive)
		argIndex++
	}

	if params.YearStart != nil {
		conditions = append(conditions, fmt.Sprintf("m.year_start >= $%d", argIndex))
		args = append(args, *params.YearStart)
		argIndex++
	}

	if params.Search != "" {
		searchCondition := fmt.Sprintf(`(
			LOWER(m.model_name) LIKE LOWER($%d) OR 
			LOWER(b.brand_name) LIKE LOWER($%d)
		)`, argIndex, argIndex)
		conditions = append(conditions, searchCondition)
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM vehicle_models m
		LEFT JOIN vehicle_brands b ON m.brand_id = b.brand_id
		WHERE %s`, whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count vehicle models: %w", err)
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT m.model_id, m.model_code, m.model_name, m.brand_id, b.brand_name,
		       m.category_id, c.category_name, m.year_start, m.year_end, m.fuel_type,
		       m.transmission, m.engine_capacity, m.is_active, m.created_at
		FROM vehicle_models m
		LEFT JOIN vehicle_brands b ON m.brand_id = b.brand_id
		LEFT JOIN vehicle_categories c ON m.category_id = c.category_id
		WHERE %s
		ORDER BY m.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list vehicle models: %w", err)
	}
	defer rows.Close()

	var models []master.VehicleModelListItem
	for rows.Next() {
		var model master.VehicleModelListItem
		err := rows.Scan(
			&model.ModelID, &model.ModelCode, &model.ModelName, &model.BrandID, &model.BrandName,
			&model.CategoryID, &model.CategoryName, &model.YearStart, &model.YearEnd,
			&model.FuelType, &model.Transmission, &model.EngineCapacity, &model.IsActive,
			&model.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan vehicle model: %w", err)
		}
		models = append(models, model)
	}

	return models, total, nil
}

// ExistsByCode checks if a vehicle model with the given code exists
func (r *vehicleModelRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_models WHERE model_code = $1)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, code).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check vehicle model code existence: %w", err)
	}
	
	return exists, nil
}

// ExistsByCodeExcludingID checks if a vehicle model with the given code exists excluding a specific ID
func (r *vehicleModelRepository) ExistsByCodeExcludingID(ctx context.Context, code string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_models WHERE model_code = $1 AND model_id != $2)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, code, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check vehicle model code existence: %w", err)
	}
	
	return exists, nil
}

// BrandExists checks if a brand with the given ID exists
func (r *vehicleModelRepository) BrandExists(ctx context.Context, brandID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_brands WHERE brand_id = $1 AND is_active = true)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, brandID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check brand existence: %w", err)
	}
	
	return exists, nil
}

// CategoryExists checks if a category with the given ID exists
func (r *vehicleModelRepository) CategoryExists(ctx context.Context, categoryID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_categories WHERE category_id = $1 AND is_active = true)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, categoryID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check category existence: %w", err)
	}
	
	return exists, nil
}

// GetNextModelCode generates the next vehicle model code
func (r *vehicleModelRepository) GetNextModelCode(ctx context.Context) (string, error) {
	query := `
		SELECT model_code 
		FROM vehicle_models 
		WHERE model_code LIKE 'MODEL-%' 
		ORDER BY model_id DESC 
		LIMIT 1`
	
	var lastCode string
	err := r.db.QueryRowContext(ctx, query).Scan(&lastCode)
	
	if err != nil {
		if err == sql.ErrNoRows {
			// First model
			return "MODEL-001", nil
		}
		return "", fmt.Errorf("failed to get last vehicle model code: %w", err)
	}
	
	// Extract number from code (e.g., "MODEL-001" -> 1)
	parts := strings.Split(lastCode, "-")
	if len(parts) != 2 {
		return "MODEL-001", nil
	}
	
	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return "MODEL-001", nil
	}
	
	// Generate next code
	nextNum := num + 1
	return fmt.Sprintf("MODEL-%03d", nextNum), nil
}