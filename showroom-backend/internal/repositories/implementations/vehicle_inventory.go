package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/transactions"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

type vehicleInventoryRepository struct {
	db *sql.DB
}

// NewVehicleInventoryRepository creates a new vehicle inventory repository
func NewVehicleInventoryRepository(db *sql.DB) interfaces.VehicleInventoryRepository {
	return &vehicleInventoryRepository{db: db}
}

// Create adds a new vehicle to inventory
func (r *vehicleInventoryRepository) Create(ctx context.Context, vehicle *transactions.VehicleInventory) (*transactions.VehicleInventory, error) {
	query := `
		INSERT INTO vehicles_inventory (
			vehicle_code, chassis_number, engine_number, license_plate, brand_id, category_id, model_id, 
			model_variant, year, color, mileage, fuel_type, transmission, engine_capacity, 
			purchase_price, estimated_selling_price, purchase_type, condition_grade, status, 
			purchase_date, purchased_from_customer_id, created_by, vehicle_images_json, 
			purchase_notes, condition_notes, has_complete_documents, documents_json
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, 
			$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27
		) RETURNING vehicle_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		vehicle.VehicleCode, vehicle.ChassisNumber, vehicle.EngineNumber, vehicle.LicensePlate,
		vehicle.BrandID, vehicle.CategoryID, vehicle.ModelID, vehicle.ModelVariant,
		vehicle.Year, vehicle.Color, vehicle.Mileage, vehicle.FuelType, vehicle.Transmission,
		vehicle.EngineCapacity, vehicle.PurchasePrice, vehicle.EstimatedSellingPrice,
		vehicle.PurchaseType, vehicle.ConditionGrade, vehicle.Status, vehicle.PurchaseDate,
		vehicle.PurchasedFromCustomerID, vehicle.CreatedBy, vehicle.VehicleImagesJSON,
		vehicle.PurchaseNotes, vehicle.ConditionNotes, vehicle.HasCompleteDocuments,
		vehicle.DocumentsJSON,
	).Scan(&vehicle.VehicleID, &vehicle.CreatedAt, &vehicle.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle inventory: %w", err)
	}

	return vehicle, nil
}

// GetByID retrieves a vehicle by its ID
func (r *vehicleInventoryRepository) GetByID(ctx context.Context, vehicleID int) (*transactions.VehicleInventory, error) {
	query := `
		SELECT 
			vi.vehicle_id, vi.vehicle_code, vi.chassis_number, vi.engine_number, vi.license_plate,
			vi.brand_id, vi.category_id, vi.model_id, vi.model_variant, vi.year, vi.color,
			vi.mileage, vi.fuel_type, vi.transmission, vi.engine_capacity, vi.purchase_price,
			vi.estimated_selling_price, vi.final_selling_price, vi.purchase_type, vi.condition_grade,
			vi.status, vi.purchase_date, vi.ready_to_sell_date, vi.sold_date, vi.created_at,
			vi.updated_at, vi.deleted_at, vi.purchased_from_customer_id, vi.created_by,
			vi.vehicle_images_json, vi.purchase_notes, vi.condition_notes, vi.has_complete_documents,
			vi.documents_json,
			vb.brand_name, vc.category_name, vm.model_name,
			COALESCE(c.customer_name, '') as customer_name,
			u.full_name as created_by_name
		FROM vehicles_inventory vi
		LEFT JOIN vehicle_brands vb ON vi.brand_id = vb.brand_id
		LEFT JOIN vehicle_categories vc ON vi.category_id = vc.category_id
		LEFT JOIN vehicle_models vm ON vi.model_id = vm.model_id
		LEFT JOIN customers c ON vi.purchased_from_customer_id = c.customer_id
		LEFT JOIN users u ON vi.created_by = u.user_id
		WHERE vi.vehicle_id = $1 AND vi.deleted_at IS NULL`

	vehicle := &transactions.VehicleInventory{}
	err := r.db.QueryRowContext(ctx, query, vehicleID).Scan(
		&vehicle.VehicleID, &vehicle.VehicleCode, &vehicle.ChassisNumber, &vehicle.EngineNumber,
		&vehicle.LicensePlate, &vehicle.BrandID, &vehicle.CategoryID, &vehicle.ModelID,
		&vehicle.ModelVariant, &vehicle.Year, &vehicle.Color, &vehicle.Mileage,
		&vehicle.FuelType, &vehicle.Transmission, &vehicle.EngineCapacity,
		&vehicle.PurchasePrice, &vehicle.EstimatedSellingPrice, &vehicle.FinalSellingPrice,
		&vehicle.PurchaseType, &vehicle.ConditionGrade, &vehicle.Status,
		&vehicle.PurchaseDate, &vehicle.ReadyToSellDate, &vehicle.SoldDate,
		&vehicle.CreatedAt, &vehicle.UpdatedAt, &vehicle.DeletedAt,
		&vehicle.PurchasedFromCustomerID, &vehicle.CreatedBy,
		&vehicle.VehicleImagesJSON, &vehicle.PurchaseNotes, &vehicle.ConditionNotes,
		&vehicle.HasCompleteDocuments, &vehicle.DocumentsJSON,
		&vehicle.BrandName, &vehicle.CategoryName, &vehicle.ModelName,
		&vehicle.CustomerName, &vehicle.CreatedByName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get vehicle by ID: %w", err)
	}

	return vehicle, nil
}

// GetByVehicleCode retrieves a vehicle by its vehicle code
func (r *vehicleInventoryRepository) GetByVehicleCode(ctx context.Context, vehicleCode string) (*transactions.VehicleInventory, error) {
	query := `
		SELECT 
			vi.vehicle_id, vi.vehicle_code, vi.chassis_number, vi.engine_number, vi.license_plate,
			vi.brand_id, vi.category_id, vi.model_id, vi.model_variant, vi.year, vi.color,
			vi.mileage, vi.fuel_type, vi.transmission, vi.engine_capacity, vi.purchase_price,
			vi.estimated_selling_price, vi.final_selling_price, vi.purchase_type, vi.condition_grade,
			vi.status, vi.purchase_date, vi.ready_to_sell_date, vi.sold_date, vi.created_at,
			vi.updated_at, vi.deleted_at, vi.purchased_from_customer_id, vi.created_by,
			vi.vehicle_images_json, vi.purchase_notes, vi.condition_notes, vi.has_complete_documents,
			vi.documents_json,
			vb.brand_name, vc.category_name, vm.model_name,
			COALESCE(c.customer_name, '') as customer_name,
			u.full_name as created_by_name
		FROM vehicles_inventory vi
		LEFT JOIN vehicle_brands vb ON vi.brand_id = vb.brand_id
		LEFT JOIN vehicle_categories vc ON vi.category_id = vc.category_id
		LEFT JOIN vehicle_models vm ON vi.model_id = vm.model_id
		LEFT JOIN customers c ON vi.purchased_from_customer_id = c.customer_id
		LEFT JOIN users u ON vi.created_by = u.user_id
		WHERE vi.vehicle_code = $1 AND vi.deleted_at IS NULL`

	vehicle := &transactions.VehicleInventory{}
	err := r.db.QueryRowContext(ctx, query, vehicleCode).Scan(
		&vehicle.VehicleID, &vehicle.VehicleCode, &vehicle.ChassisNumber, &vehicle.EngineNumber,
		&vehicle.LicensePlate, &vehicle.BrandID, &vehicle.CategoryID, &vehicle.ModelID,
		&vehicle.ModelVariant, &vehicle.Year, &vehicle.Color, &vehicle.Mileage,
		&vehicle.FuelType, &vehicle.Transmission, &vehicle.EngineCapacity,
		&vehicle.PurchasePrice, &vehicle.EstimatedSellingPrice, &vehicle.FinalSellingPrice,
		&vehicle.PurchaseType, &vehicle.ConditionGrade, &vehicle.Status,
		&vehicle.PurchaseDate, &vehicle.ReadyToSellDate, &vehicle.SoldDate,
		&vehicle.CreatedAt, &vehicle.UpdatedAt, &vehicle.DeletedAt,
		&vehicle.PurchasedFromCustomerID, &vehicle.CreatedBy,
		&vehicle.VehicleImagesJSON, &vehicle.PurchaseNotes, &vehicle.ConditionNotes,
		&vehicle.HasCompleteDocuments, &vehicle.DocumentsJSON,
		&vehicle.BrandName, &vehicle.CategoryName, &vehicle.ModelName,
		&vehicle.CustomerName, &vehicle.CreatedByName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get vehicle by code: %w", err)
	}

	return vehicle, nil
}

// GetByChassisNumber retrieves a vehicle by its chassis number
func (r *vehicleInventoryRepository) GetByChassisNumber(ctx context.Context, chassisNumber string) (*transactions.VehicleInventory, error) {
	query := `
		SELECT 
			vi.vehicle_id, vi.vehicle_code, vi.chassis_number, vi.engine_number, vi.license_plate,
			vi.brand_id, vi.category_id, vi.model_id, vi.model_variant, vi.year, vi.color,
			vi.mileage, vi.fuel_type, vi.transmission, vi.engine_capacity, vi.purchase_price,
			vi.estimated_selling_price, vi.final_selling_price, vi.purchase_type, vi.condition_grade,
			vi.status, vi.purchase_date, vi.ready_to_sell_date, vi.sold_date, vi.created_at,
			vi.updated_at, vi.deleted_at, vi.purchased_from_customer_id, vi.created_by,
			vi.vehicle_images_json, vi.purchase_notes, vi.condition_notes, vi.has_complete_documents,
			vi.documents_json,
			vb.brand_name, vc.category_name, vm.model_name,
			COALESCE(c.customer_name, '') as customer_name,
			u.full_name as created_by_name
		FROM vehicles_inventory vi
		LEFT JOIN vehicle_brands vb ON vi.brand_id = vb.brand_id
		LEFT JOIN vehicle_categories vc ON vi.category_id = vc.category_id
		LEFT JOIN vehicle_models vm ON vi.model_id = vm.model_id
		LEFT JOIN customers c ON vi.purchased_from_customer_id = c.customer_id
		LEFT JOIN users u ON vi.created_by = u.user_id
		WHERE vi.chassis_number = $1 AND vi.deleted_at IS NULL`

	vehicle := &transactions.VehicleInventory{}
	err := r.db.QueryRowContext(ctx, query, chassisNumber).Scan(
		&vehicle.VehicleID, &vehicle.VehicleCode, &vehicle.ChassisNumber, &vehicle.EngineNumber,
		&vehicle.LicensePlate, &vehicle.BrandID, &vehicle.CategoryID, &vehicle.ModelID,
		&vehicle.ModelVariant, &vehicle.Year, &vehicle.Color, &vehicle.Mileage,
		&vehicle.FuelType, &vehicle.Transmission, &vehicle.EngineCapacity,
		&vehicle.PurchasePrice, &vehicle.EstimatedSellingPrice, &vehicle.FinalSellingPrice,
		&vehicle.PurchaseType, &vehicle.ConditionGrade, &vehicle.Status,
		&vehicle.PurchaseDate, &vehicle.ReadyToSellDate, &vehicle.SoldDate,
		&vehicle.CreatedAt, &vehicle.UpdatedAt, &vehicle.DeletedAt,
		&vehicle.PurchasedFromCustomerID, &vehicle.CreatedBy,
		&vehicle.VehicleImagesJSON, &vehicle.PurchaseNotes, &vehicle.ConditionNotes,
		&vehicle.HasCompleteDocuments, &vehicle.DocumentsJSON,
		&vehicle.BrandName, &vehicle.CategoryName, &vehicle.ModelName,
		&vehicle.CustomerName, &vehicle.CreatedByName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get vehicle by chassis number: %w", err)
	}

	return vehicle, nil
}

// List retrieves vehicles with filtering and pagination
func (r *vehicleInventoryRepository) List(ctx context.Context, filter *transactions.VehicleInventoryFilterParams) ([]*transactions.VehicleInventoryListItem, *common.PaginationMeta, error) {
	// Base query
	baseQuery := `
		FROM vehicles_inventory vi
		LEFT JOIN vehicle_brands vb ON vi.brand_id = vb.brand_id
		LEFT JOIN vehicle_categories vc ON vi.category_id = vc.category_id
		LEFT JOIN vehicle_models vm ON vi.model_id = vm.model_id
		WHERE vi.deleted_at IS NULL`

	var conditions []string
	var args []interface{}
	argIndex := 1

	// Build WHERE conditions
	if filter.BrandID != nil {
		conditions = append(conditions, fmt.Sprintf("vi.brand_id = $%d", argIndex))
		args = append(args, *filter.BrandID)
		argIndex++
	}

	if filter.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("vi.category_id = $%d", argIndex))
		args = append(args, *filter.CategoryID)
		argIndex++
	}

	if filter.ModelID != nil {
		conditions = append(conditions, fmt.Sprintf("vi.model_id = $%d", argIndex))
		args = append(args, *filter.ModelID)
		argIndex++
	}

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("vi.status = $%d", argIndex))
		args = append(args, string(*filter.Status))
		argIndex++
	}

	if filter.PurchaseType != nil {
		conditions = append(conditions, fmt.Sprintf("vi.purchase_type = $%d", argIndex))
		args = append(args, string(*filter.PurchaseType))
		argIndex++
	}

	if filter.ConditionGrade != nil {
		conditions = append(conditions, fmt.Sprintf("vi.condition_grade = $%d", argIndex))
		args = append(args, string(*filter.ConditionGrade))
		argIndex++
	}

	if filter.MinYear != nil {
		conditions = append(conditions, fmt.Sprintf("vi.year >= $%d", argIndex))
		args = append(args, *filter.MinYear)
		argIndex++
	}

	if filter.MaxYear != nil {
		conditions = append(conditions, fmt.Sprintf("vi.year <= $%d", argIndex))
		args = append(args, *filter.MaxYear)
		argIndex++
	}

	if filter.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("vi.purchase_price >= $%d", argIndex))
		args = append(args, *filter.MinPrice)
		argIndex++
	}

	if filter.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("vi.purchase_price <= $%d", argIndex))
		args = append(args, *filter.MaxPrice)
		argIndex++
	}

	if filter.Search != "" {
		searchCondition := fmt.Sprintf(`(
			vi.vehicle_code ILIKE $%d OR 
			vi.chassis_number ILIKE $%d OR 
			vi.engine_number ILIKE $%d OR 
			vi.license_plate ILIKE $%d OR
			vb.brand_name ILIKE $%d OR
			vm.model_name ILIKE $%d
		)`, argIndex, argIndex, argIndex, argIndex, argIndex, argIndex)
		conditions = append(conditions, searchCondition)
		searchPattern := "%" + filter.Search + "%"
		args = append(args, searchPattern)
		argIndex++
	}

	// Add conditions to base query
	whereClause := baseQuery
	if len(conditions) > 0 {
		whereClause += " AND " + strings.Join(conditions, " AND ")
	}

	// Count total records
	countQuery := "SELECT COUNT(*) " + whereClause
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count vehicles: %w", err)
	}

	// Validate pagination
	filter.PaginationParams.Validate()

	// Build final query with pagination
	finalQuery := `
		SELECT 
			vi.vehicle_id, vi.vehicle_code, vi.chassis_number, vi.license_plate,
			vb.brand_name, vm.model_name, vi.year, vi.color, vi.mileage,
			vi.purchase_price, vi.estimated_selling_price, vi.condition_grade,
			vi.status, vi.purchase_date, vi.created_at
		` + whereClause + `
		ORDER BY vi.created_at DESC
		LIMIT $` + fmt.Sprintf("%d", argIndex) + ` OFFSET $` + fmt.Sprintf("%d", argIndex+1)

	args = append(args, filter.Limit, filter.GetOffset())

	rows, err := r.db.QueryContext(ctx, finalQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []*transactions.VehicleInventoryListItem
	for rows.Next() {
		vehicle := &transactions.VehicleInventoryListItem{}
		err := rows.Scan(
			&vehicle.VehicleID, &vehicle.VehicleCode, &vehicle.ChassisNumber,
			&vehicle.LicensePlate, &vehicle.BrandName, &vehicle.ModelName,
			&vehicle.Year, &vehicle.Color, &vehicle.Mileage, &vehicle.PurchasePrice,
			&vehicle.EstimatedSellingPrice, &vehicle.ConditionGrade, &vehicle.Status,
			&vehicle.PurchaseDate, &vehicle.CreatedAt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan vehicle: %w", err)
		}
		vehicles = append(vehicles, vehicle)
	}

	// Create pagination meta
	meta := &common.PaginationMeta{
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: filter.GetTotalPages(total),
		HasMore:    filter.GetHasMore(total),
	}

	return vehicles, meta, nil
}

// Update updates a vehicle's information
func (r *vehicleInventoryRepository) Update(ctx context.Context, vehicleID int, updateReq *transactions.VehicleInventoryUpdateRequest) (*transactions.VehicleInventory, error) {
	var setParts []string
	var args []interface{}
	argIndex := 1

	if updateReq.LicensePlate != nil {
		setParts = append(setParts, fmt.Sprintf("license_plate = $%d", argIndex))
		args = append(args, *updateReq.LicensePlate)
		argIndex++
	}

	if updateReq.ModelVariant != nil {
		setParts = append(setParts, fmt.Sprintf("model_variant = $%d", argIndex))
		args = append(args, *updateReq.ModelVariant)
		argIndex++
	}

	if updateReq.Color != nil {
		setParts = append(setParts, fmt.Sprintf("color = $%d", argIndex))
		args = append(args, *updateReq.Color)
		argIndex++
	}

	if updateReq.Mileage != nil {
		setParts = append(setParts, fmt.Sprintf("mileage = $%d", argIndex))
		args = append(args, *updateReq.Mileage)
		argIndex++
	}

	if updateReq.EstimatedSellingPrice != nil {
		setParts = append(setParts, fmt.Sprintf("estimated_selling_price = $%d", argIndex))
		args = append(args, *updateReq.EstimatedSellingPrice)
		argIndex++
	}

	if updateReq.FinalSellingPrice != nil {
		setParts = append(setParts, fmt.Sprintf("final_selling_price = $%d", argIndex))
		args = append(args, *updateReq.FinalSellingPrice)
		argIndex++
	}

	if updateReq.ConditionGrade != nil {
		setParts = append(setParts, fmt.Sprintf("condition_grade = $%d", argIndex))
		args = append(args, string(*updateReq.ConditionGrade))
		argIndex++
	}

	if updateReq.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, string(*updateReq.Status))
		argIndex++
	}

	if updateReq.ReadyToSellDate != nil {
		setParts = append(setParts, fmt.Sprintf("ready_to_sell_date = $%d", argIndex))
		args = append(args, *updateReq.ReadyToSellDate)
		argIndex++
	}

	if updateReq.PurchaseNotes != nil {
		setParts = append(setParts, fmt.Sprintf("purchase_notes = $%d", argIndex))
		args = append(args, *updateReq.PurchaseNotes)
		argIndex++
	}

	if updateReq.ConditionNotes != nil {
		setParts = append(setParts, fmt.Sprintf("condition_notes = $%d", argIndex))
		args = append(args, *updateReq.ConditionNotes)
		argIndex++
	}

	if updateReq.HasCompleteDocuments != nil {
		setParts = append(setParts, fmt.Sprintf("has_complete_documents = $%d", argIndex))
		args = append(args, *updateReq.HasCompleteDocuments)
		argIndex++
	}

	if len(setParts) == 0 {
		return r.GetByID(ctx, vehicleID)
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	query := fmt.Sprintf(`
		UPDATE vehicles_inventory 
		SET %s 
		WHERE vehicle_id = $%d AND deleted_at IS NULL`,
		strings.Join(setParts, ", "), argIndex)

	args = append(args, vehicleID)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update vehicle: %w", err)
	}

	return r.GetByID(ctx, vehicleID)
}

// UpdateStatus updates a vehicle's status
func (r *vehicleInventoryRepository) UpdateStatus(ctx context.Context, vehicleID int, status transactions.VehicleStatus, updatedBy int) error {
	query := `
		UPDATE vehicles_inventory 
		SET status = $1, updated_at = $2
		WHERE vehicle_id = $3 AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, string(status), time.Now(), vehicleID)
	if err != nil {
		return fmt.Errorf("failed to update vehicle status: %w", err)
	}

	return nil
}

// SoftDelete soft deletes a vehicle (if allowed)
func (r *vehicleInventoryRepository) SoftDelete(ctx context.Context, vehicleID int, deletedBy int) error {
	query := `
		UPDATE vehicles_inventory 
		SET deleted_at = $1, updated_at = $1
		WHERE vehicle_id = $2 AND deleted_at IS NULL AND status != 'sold'`

	result, err := r.db.ExecContext(ctx, query, time.Now(), vehicleID)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vehicle not found or cannot be deleted")
	}

	return nil
}

// GetAvailableForSale retrieves vehicles available for sale
func (r *vehicleInventoryRepository) GetAvailableForSale(ctx context.Context, filter *transactions.VehicleInventoryFilterParams) ([]*transactions.VehicleInventoryListItem, *common.PaginationMeta, error) {
	// Set status filter to ready_to_sell
	readyStatus := transactions.VehicleStatusReadyToSell
	filter.Status = &readyStatus
	return r.List(ctx, filter)
}

// GetNeedingRepair retrieves vehicles that need repair
func (r *vehicleInventoryRepository) GetNeedingRepair(ctx context.Context, filter *transactions.VehicleInventoryFilterParams) ([]*transactions.VehicleInventoryListItem, *common.PaginationMeta, error) {
	// Set status filter to pending repair approval or in repair
	pendingStatus := transactions.VehicleStatusPendingRepairApproval
	filter.Status = &pendingStatus
	return r.List(ctx, filter)
}

// GenerateVehicleCode generates a unique vehicle code
func (r *vehicleInventoryRepository) GenerateVehicleCode(ctx context.Context) (string, error) {
	query := `
		SELECT COALESCE(MAX(CAST(SUBSTRING(vehicle_code FROM 4) AS INTEGER)), 0) + 1
		FROM vehicles_inventory 
		WHERE vehicle_code ~ '^VH-[0-9]+$'`

	var nextNumber int
	err := r.db.QueryRowContext(ctx, query).Scan(&nextNumber)
	if err != nil {
		return "", fmt.Errorf("failed to generate vehicle code: %w", err)
	}

	return fmt.Sprintf("VH-%03d", nextNumber), nil
}

// CheckChassisNumberExists checks if chassis number already exists
func (r *vehicleInventoryRepository) CheckChassisNumberExists(ctx context.Context, chassisNumber string, excludeVehicleID ...int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicles_inventory WHERE chassis_number = $1 AND deleted_at IS NULL`
	args := []interface{}{chassisNumber}

	if len(excludeVehicleID) > 0 {
		query += ` AND vehicle_id != $2`
		args = append(args, excludeVehicleID[0])
	}

	query += `)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check chassis number: %w", err)
	}

	return exists, nil
}

// CheckEngineNumberExists checks if engine number already exists
func (r *vehicleInventoryRepository) CheckEngineNumberExists(ctx context.Context, engineNumber string, excludeVehicleID ...int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicles_inventory WHERE engine_number = $1 AND deleted_at IS NULL`
	args := []interface{}{engineNumber}

	if len(excludeVehicleID) > 0 {
		query += ` AND vehicle_id != $2`
		args = append(args, excludeVehicleID[0])
	}

	query += `)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check engine number: %w", err)
	}

	return exists, nil
}