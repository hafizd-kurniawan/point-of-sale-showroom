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

type vehicleBrandRepository struct {
	db *sql.DB
}

// NewVehicleBrandRepository creates a new vehicle brand repository
func NewVehicleBrandRepository(db *sql.DB) interfaces.VehicleBrandRepository {
	return &vehicleBrandRepository{db: db}
}

// Create creates a new vehicle brand
func (r *vehicleBrandRepository) Create(ctx context.Context, brand *master.VehicleBrand) (*master.VehicleBrand, error) {
	query := `
		INSERT INTO vehicle_brands (brand_code, brand_name, country_origin, logo_image, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING brand_id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		brand.BrandCode, brand.BrandName, brand.CountryOrigin, brand.LogoImage, brand.CreatedBy,
	).Scan(&brand.BrandID, &brand.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle brand: %w", err)
	}

	return brand, nil
}

// GetByID retrieves a vehicle brand by ID
func (r *vehicleBrandRepository) GetByID(ctx context.Context, id int) (*master.VehicleBrand, error) {
	query := `
		SELECT b.brand_id, b.brand_code, b.brand_name, b.country_origin, b.logo_image,
		       b.is_active, b.created_at, b.created_by,
		       u.user_id as creator_user_id, u.username as creator_username, u.full_name as creator_full_name
		FROM vehicle_brands b
		LEFT JOIN users u ON b.created_by = u.user_id
		WHERE b.brand_id = $1 AND b.is_active = true`

	brand := &master.VehicleBrand{}
	creator := &user.UserCreatorInfo{}
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&brand.BrandID, &brand.BrandCode, &brand.BrandName, &brand.CountryOrigin,
		&brand.LogoImage, &brand.IsActive, &brand.CreatedAt, &brand.CreatedBy,
		&creator.UserID, &creator.Username, &creator.FullName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle brand not found")
		}
		return nil, fmt.Errorf("failed to get vehicle brand: %w", err)
	}

	brand.Creator = creator
	return brand, nil
}

// GetByCode retrieves a vehicle brand by code
func (r *vehicleBrandRepository) GetByCode(ctx context.Context, code string) (*master.VehicleBrand, error) {
	query := `
		SELECT b.brand_id, b.brand_code, b.brand_name, b.country_origin, b.logo_image,
		       b.is_active, b.created_at, b.created_by
		FROM vehicle_brands b
		WHERE b.brand_code = $1 AND b.is_active = true`

	brand := &master.VehicleBrand{}
	
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&brand.BrandID, &brand.BrandCode, &brand.BrandName, &brand.CountryOrigin,
		&brand.LogoImage, &brand.IsActive, &brand.CreatedAt, &brand.CreatedBy,
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
func (r *vehicleBrandRepository) Update(ctx context.Context, id int, brand *master.VehicleBrand) (*master.VehicleBrand, error) {
	query := `
		UPDATE vehicle_brands 
		SET brand_name = $1, country_origin = $2, logo_image = $3, is_active = $4
		WHERE brand_id = $5 AND is_active = true
		RETURNING brand_id`

	var brandID int
	err := r.db.QueryRowContext(ctx, query,
		brand.BrandName, brand.CountryOrigin, brand.LogoImage, brand.IsActive, id,
	).Scan(&brandID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle brand not found")
		}
		return nil, fmt.Errorf("failed to update vehicle brand: %w", err)
	}

	brand.BrandID = brandID
	return brand, nil
}

// Delete soft deletes a vehicle brand
func (r *vehicleBrandRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE vehicle_brands SET is_active = false WHERE brand_id = $1 AND is_active = true`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle brand: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vehicle brand not found")
	}

	return nil
}

// ListActive retrieves all active vehicle brands
func (r *vehicleBrandRepository) ListActive(ctx context.Context) ([]master.VehicleBrand, error) {
	query := `
		SELECT brand_id, brand_code, brand_name, country_origin, logo_image, is_active, created_at
		FROM vehicle_brands
		WHERE is_active = true
		ORDER BY brand_name ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list vehicle brands: %w", err)
	}
	defer rows.Close()

	var brands []master.VehicleBrand
	for rows.Next() {
		var brand master.VehicleBrand
		err := rows.Scan(
			&brand.BrandID, &brand.BrandCode, &brand.BrandName, &brand.CountryOrigin,
			&brand.LogoImage, &brand.IsActive, &brand.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vehicle brand: %w", err)
		}
		brands = append(brands, brand)
	}

	return brands, nil
}

// List retrieves vehicle brands with optional filtering by active status
func (r *vehicleBrandRepository) List(ctx context.Context, isActive *bool) ([]master.VehicleBrand, error) {
	query := `
		SELECT brand_id, brand_code, brand_name, country_origin, logo_image, is_active, created_at
		FROM vehicle_brands`
	
	var args []interface{}
	if isActive != nil {
		query += ` WHERE is_active = $1`
		args = append(args, *isActive)
	}
	
	query += ` ORDER BY brand_name ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list vehicle brands: %w", err)
	}
	defer rows.Close()

	var brands []master.VehicleBrand
	for rows.Next() {
		var brand master.VehicleBrand
		err := rows.Scan(
			&brand.BrandID, &brand.BrandCode, &brand.BrandName, &brand.CountryOrigin,
			&brand.LogoImage, &brand.IsActive, &brand.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vehicle brand: %w", err)
		}
		brands = append(brands, brand)
	}

	return brands, nil
}

// ExistsByCode checks if a vehicle brand with the given code exists
func (r *vehicleBrandRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_brands WHERE brand_code = $1)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, code).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check vehicle brand code existence: %w", err)
	}
	
	return exists, nil
}

// ExistsByName checks if a vehicle brand with the given name exists
func (r *vehicleBrandRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_brands WHERE LOWER(brand_name) = LOWER($1) AND is_active = true)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check vehicle brand name existence: %w", err)
	}
	
	return exists, nil
}

// ExistsByCodeExcludingID checks if a vehicle brand with the given code exists excluding a specific ID
func (r *vehicleBrandRepository) ExistsByCodeExcludingID(ctx context.Context, code string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_brands WHERE brand_code = $1 AND brand_id != $2)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, code, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check vehicle brand code existence: %w", err)
	}
	
	return exists, nil
}

// ExistsByNameExcludingID checks if a vehicle brand with the given name exists excluding a specific ID
func (r *vehicleBrandRepository) ExistsByNameExcludingID(ctx context.Context, name string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_brands WHERE LOWER(brand_name) = LOWER($1) AND brand_id != $2 AND is_active = true)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, name, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check vehicle brand name existence: %w", err)
	}
	
	return exists, nil
}

// GetNextBrandCode generates the next vehicle brand code
func (r *vehicleBrandRepository) GetNextBrandCode(ctx context.Context) (string, error) {
	query := `
		SELECT brand_code 
		FROM vehicle_brands 
		WHERE brand_code LIKE 'BRAND-%' 
		ORDER BY brand_id DESC 
		LIMIT 1`
	
	var lastCode string
	err := r.db.QueryRowContext(ctx, query).Scan(&lastCode)
	
	if err != nil {
		if err == sql.ErrNoRows {
			// First brand
			return "BRAND-001", nil
		}
		return "", fmt.Errorf("failed to get last vehicle brand code: %w", err)
	}
	
	// Extract number from code (e.g., "BRAND-001" -> 1)
	parts := strings.Split(lastCode, "-")
	if len(parts) != 2 {
		return "BRAND-001", nil
	}
	
	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return "BRAND-001", nil
	}
	
	// Generate next code
	nextNum := num + 1
	return fmt.Sprintf("BRAND-%03d", nextNum), nil
}