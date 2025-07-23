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

// VehicleModelRepository implements interfaces.VehicleModelRepository
type VehicleModelRepository struct {
	db *sql.DB
}

// NewVehicleModelRepository creates a new vehicle model repository
func NewVehicleModelRepository(db *sql.DB) interfaces.VehicleModelRepository {
	return &VehicleModelRepository{db: db}
}

// Create creates a new vehicle model
func (r *VehicleModelRepository) Create(ctx context.Context, model *master.VehicleModel) (*master.VehicleModel, error) {
	query := `
		INSERT INTO vehicle_models (model_code, model_name, brand_id, category_id, model_year, engine_capacity, fuel_type, transmission, seat_capacity, color, price, description, image_url, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING model_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		model.ModelCode,
		model.ModelName,
		model.BrandID,
		model.CategoryID,
		model.ModelYear,
		model.EngineCapacity,
		model.FuelType,
		model.Transmission,
		model.SeatCapacity,
		model.Color,
		model.Price,
		model.Description,
		model.ImageURL,
		model.CreatedBy,
	).Scan(&model.ModelID, &model.CreatedAt, &model.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle model: %w", err)
	}

	model.IsActive = true
	return model, nil
}

// GetByID retrieves a vehicle model by ID with related data
func (r *VehicleModelRepository) GetByID(ctx context.Context, id int) (*master.VehicleModel, error) {
	query := `
		SELECT vm.model_id, vm.model_code, vm.model_name, vm.brand_id, vm.category_id, vm.model_year, vm.engine_capacity, vm.fuel_type, vm.transmission, vm.seat_capacity, vm.color, vm.price, vm.description, vm.image_url, vm.is_active, vm.created_at, vm.updated_at, vm.created_by,
		       vb.brand_name, vc.category_name
		FROM vehicle_models vm
		JOIN vehicle_brands vb ON vm.brand_id = vb.brand_id
		JOIN vehicle_categories vc ON vm.category_id = vc.category_id
		WHERE vm.model_id = $1`

	model := &master.VehicleModel{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&model.ModelID,
		&model.ModelCode,
		&model.ModelName,
		&model.BrandID,
		&model.CategoryID,
		&model.ModelYear,
		&model.EngineCapacity,
		&model.FuelType,
		&model.Transmission,
		&model.SeatCapacity,
		&model.Color,
		&model.Price,
		&model.Description,
		&model.ImageURL,
		&model.IsActive,
		&model.CreatedAt,
		&model.UpdatedAt,
		&model.CreatedBy,
		&model.BrandName,
		&model.CategoryName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle model not found")
		}
		return nil, fmt.Errorf("failed to get vehicle model: %w", err)
	}

	return model, nil
}

// GetByCode retrieves a vehicle model by code with related data
func (r *VehicleModelRepository) GetByCode(ctx context.Context, code string) (*master.VehicleModel, error) {
	query := `
		SELECT vm.model_id, vm.model_code, vm.model_name, vm.brand_id, vm.category_id, vm.model_year, vm.engine_capacity, vm.fuel_type, vm.transmission, vm.seat_capacity, vm.color, vm.price, vm.description, vm.image_url, vm.is_active, vm.created_at, vm.updated_at, vm.created_by,
		       vb.brand_name, vc.category_name
		FROM vehicle_models vm
		JOIN vehicle_brands vb ON vm.brand_id = vb.brand_id
		JOIN vehicle_categories vc ON vm.category_id = vc.category_id
		WHERE vm.model_code = $1`

	model := &master.VehicleModel{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&model.ModelID,
		&model.ModelCode,
		&model.ModelName,
		&model.BrandID,
		&model.CategoryID,
		&model.ModelYear,
		&model.EngineCapacity,
		&model.FuelType,
		&model.Transmission,
		&model.SeatCapacity,
		&model.Color,
		&model.Price,
		&model.Description,
		&model.ImageURL,
		&model.IsActive,
		&model.CreatedAt,
		&model.UpdatedAt,
		&model.CreatedBy,
		&model.BrandName,
		&model.CategoryName,
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
func (r *VehicleModelRepository) Update(ctx context.Context, id int, model *master.VehicleModel) (*master.VehicleModel, error) {
	query := `
		UPDATE vehicle_models
		SET model_name = $1, brand_id = $2, category_id = $3, model_year = $4, engine_capacity = $5, fuel_type = $6, transmission = $7, seat_capacity = $8, color = $9, price = $10, description = $11, image_url = $12, is_active = $13, updated_at = NOW()
		WHERE model_id = $14
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		model.ModelName,
		model.BrandID,
		model.CategoryID,
		model.ModelYear,
		model.EngineCapacity,
		model.FuelType,
		model.Transmission,
		model.SeatCapacity,
		model.Color,
		model.Price,
		model.Description,
		model.ImageURL,
		model.IsActive,
		id,
	).Scan(&model.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle model not found")
		}
		return nil, fmt.Errorf("failed to update vehicle model: %w", err)
	}

	model.ModelID = id
	return model, nil
}

// Delete soft deletes a vehicle model
func (r *VehicleModelRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE vehicle_models SET is_active = FALSE, updated_at = NOW() WHERE model_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle model: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vehicle model not found")
	}

	return nil
}

// List retrieves vehicle models with filtering and pagination
func (r *VehicleModelRepository) List(ctx context.Context, params *master.VehicleModelFilterParams) (*common.PaginatedResponse, error) {
	params.Validate()

	// Build WHERE conditions
	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.BrandID != nil {
		conditions = append(conditions, fmt.Sprintf("vm.brand_id = $%d", argIndex))
		args = append(args, *params.BrandID)
		argIndex++
	}

	if params.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("vm.category_id = $%d", argIndex))
		args = append(args, *params.CategoryID)
		argIndex++
	}

	if params.ModelYear != nil {
		conditions = append(conditions, fmt.Sprintf("vm.model_year = $%d", argIndex))
		args = append(args, *params.ModelYear)
		argIndex++
	}

	if params.FuelType != "" {
		conditions = append(conditions, fmt.Sprintf("vm.fuel_type ILIKE $%d", argIndex))
		args = append(args, "%"+params.FuelType+"%")
		argIndex++
	}

	if params.Transmission != "" {
		conditions = append(conditions, fmt.Sprintf("vm.transmission ILIKE $%d", argIndex))
		args = append(args, "%"+params.Transmission+"%")
		argIndex++
	}

	if params.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("vm.price >= $%d", argIndex))
		args = append(args, *params.MinPrice)
		argIndex++
	}

	if params.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("vm.price <= $%d", argIndex))
		args = append(args, *params.MaxPrice)
		argIndex++
	}

	if params.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("vm.is_active = $%d", argIndex))
		args = append(args, *params.IsActive)
		argIndex++
	}

	if params.Search != "" {
		searchCondition := fmt.Sprintf("(vm.model_name ILIKE $%d OR vm.model_code ILIKE $%d OR vb.brand_name ILIKE $%d)", argIndex, argIndex, argIndex)
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
		FROM vehicle_models vm
		JOIN vehicle_brands vb ON vm.brand_id = vb.brand_id
		JOIN vehicle_categories vc ON vm.category_id = vc.category_id
		%s`, whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count vehicle models: %w", err)
	}

	// Build main query
	query := fmt.Sprintf(`
		SELECT vm.model_id, vm.model_code, vm.model_name, vb.brand_name, vc.category_name, vm.model_year, vm.fuel_type, vm.transmission, vm.seat_capacity, vm.price, vm.is_active, vm.created_at
		FROM vehicle_models vm
		JOIN vehicle_brands vb ON vm.brand_id = vb.brand_id
		JOIN vehicle_categories vc ON vm.category_id = vc.category_id
		%s
		ORDER BY vm.created_at DESC
		LIMIT $%d OFFSET $%d`,
		whereClause, argIndex, argIndex+1)

	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list vehicle models: %w", err)
	}
	defer rows.Close()

	var models []master.VehicleModelListItem
	for rows.Next() {
		var model master.VehicleModelListItem
		err := rows.Scan(
			&model.ModelID,
			&model.ModelCode,
			&model.ModelName,
			&model.BrandName,
			&model.CategoryName,
			&model.ModelYear,
			&model.FuelType,
			&model.Transmission,
			&model.SeatCapacity,
			&model.Price,
			&model.IsActive,
			&model.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vehicle model: %w", err)
		}
		models = append(models, model)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate vehicle models: %w", err)
	}

	return &common.PaginatedResponse{
		Data:       models,
		Total: int(total),
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: params.GetTotalPages(total),
		HasMore:    params.GetHasMore(total),
	}, nil
}

// GenerateCode generates a new vehicle model code
func (r *VehicleModelRepository) GenerateCode(ctx context.Context) (string, error) {
	query := `
		SELECT model_code
		FROM vehicle_models
		WHERE model_code LIKE 'VM-%'
		ORDER BY model_code DESC
		LIMIT 1`

	var lastCode sql.NullString
	err := r.db.QueryRowContext(ctx, query).Scan(&lastCode)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to get last vehicle model code: %w", err)
	}

	nextNumber := 1
	if lastCode.Valid {
		// Extract number from code (e.g., "VM-001" -> "001" -> 1)
		parts := strings.Split(lastCode.String, "-")
		if len(parts) == 2 {
			if num, err := strconv.Atoi(parts[1]); err == nil {
				nextNumber = num + 1
			}
		}
	}

	return fmt.Sprintf("VM-%03d", nextNumber), nil
}

// IsCodeExists checks if a vehicle model code already exists
func (r *VehicleModelRepository) IsCodeExists(ctx context.Context, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_models WHERE model_code = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, code).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check vehicle model code existence: %w", err)
	}

	return exists, nil
}