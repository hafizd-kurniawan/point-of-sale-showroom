package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/vehicle_purchase"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// VehiclePurchaseTransactionRepository implements interfaces.VehiclePurchaseTransactionRepository
type VehiclePurchaseTransactionRepository struct {
	db *sql.DB
}

// NewVehiclePurchaseTransactionRepository creates a new vehicle purchase transaction repository
func NewVehiclePurchaseTransactionRepository(db *sql.DB) interfaces.VehiclePurchaseTransactionRepository {
	return &VehiclePurchaseTransactionRepository{db: db}
}

// Create creates a new vehicle purchase transaction
func (r *VehiclePurchaseTransactionRepository) Create(ctx context.Context, transaction *vehicle_purchase.VehiclePurchaseTransaction) (*vehicle_purchase.VehiclePurchaseTransaction, error) {
	query := `
		INSERT INTO vehicle_purchase_transactions (
			transaction_number, customer_id, vehicle_id, vin_number, vehicle_brand, vehicle_model,
			vehicle_year, vehicle_color, engine_number, registration_number, purchase_price,
			agreed_value, odometer_reading, fuel_type, transmission, condition_rating,
			purchase_date, transaction_status, inspection_notes, evaluation_notes,
			purchase_notes, documents_json, processed_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
		RETURNING transaction_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		transaction.TransactionNumber,
		transaction.CustomerID,
		transaction.VehicleID,
		transaction.VinNumber,
		transaction.VehicleBrand,
		transaction.VehicleModel,
		transaction.VehicleYear,
		transaction.VehicleColor,
		transaction.EngineNumber,
		transaction.RegistrationNumber,
		transaction.PurchasePrice,
		transaction.AgreedValue,
		transaction.OdometerReading,
		transaction.FuelType,
		transaction.Transmission,
		transaction.ConditionRating,
		transaction.PurchaseDate,
		transaction.TransactionStatus,
		transaction.InspectionNotes,
		transaction.EvaluationNotes,
		transaction.PurchaseNotes,
		transaction.DocumentsJSON,
		transaction.ProcessedBy,
	).Scan(&transaction.TransactionID, &transaction.CreatedAt, &transaction.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle purchase transaction: %w", err)
	}

	return transaction, nil
}

// GetByID retrieves a vehicle purchase transaction by ID
func (r *VehiclePurchaseTransactionRepository) GetByID(ctx context.Context, id int) (*vehicle_purchase.VehiclePurchaseTransaction, error) {
	query := `
		SELECT 
			vpt.transaction_id, vpt.transaction_number, vpt.customer_id, vpt.vehicle_id,
			vpt.vin_number, vpt.vehicle_brand, vpt.vehicle_model, vpt.vehicle_year,
			vpt.vehicle_color, vpt.engine_number, vpt.registration_number, vpt.purchase_price,
			vpt.agreed_value, vpt.odometer_reading, vpt.fuel_type, vpt.transmission,
			vpt.condition_rating, vpt.purchase_date, vpt.transaction_status, vpt.inspection_notes,
			vpt.evaluation_notes, vpt.purchase_notes, vpt.documents_json, vpt.processed_by,
			vpt.inspected_by, vpt.approved_by, vpt.approved_at, vpt.created_at, vpt.updated_at,
			c.customer_name,
			u1.full_name as processed_by_name,
			u2.full_name as inspected_by_name,
			u3.full_name as approved_by_name
		FROM vehicle_purchase_transactions vpt
		LEFT JOIN customers c ON vpt.customer_id = c.customer_id
		LEFT JOIN users u1 ON vpt.processed_by = u1.user_id
		LEFT JOIN users u2 ON vpt.inspected_by = u2.user_id
		LEFT JOIN users u3 ON vpt.approved_by = u3.user_id
		WHERE vpt.transaction_id = $1`

	transaction := &vehicle_purchase.VehiclePurchaseTransaction{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.TransactionID,
		&transaction.TransactionNumber,
		&transaction.CustomerID,
		&transaction.VehicleID,
		&transaction.VinNumber,
		&transaction.VehicleBrand,
		&transaction.VehicleModel,
		&transaction.VehicleYear,
		&transaction.VehicleColor,
		&transaction.EngineNumber,
		&transaction.RegistrationNumber,
		&transaction.PurchasePrice,
		&transaction.AgreedValue,
		&transaction.OdometerReading,
		&transaction.FuelType,
		&transaction.Transmission,
		&transaction.ConditionRating,
		&transaction.PurchaseDate,
		&transaction.TransactionStatus,
		&transaction.InspectionNotes,
		&transaction.EvaluationNotes,
		&transaction.PurchaseNotes,
		&transaction.DocumentsJSON,
		&transaction.ProcessedBy,
		&transaction.InspectedBy,
		&transaction.ApprovedBy,
		&transaction.ApprovedAt,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.CustomerName,
		&transaction.ProcessedByName,
		&transaction.InspectedByName,
		&transaction.ApprovedByName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle purchase transaction not found")
		}
		return nil, fmt.Errorf("failed to get vehicle purchase transaction: %w", err)
	}

	return transaction, nil
}

// GetByNumber retrieves a vehicle purchase transaction by transaction number
func (r *VehiclePurchaseTransactionRepository) GetByNumber(ctx context.Context, number string) (*vehicle_purchase.VehiclePurchaseTransaction, error) {
	query := `
		SELECT 
			vpt.transaction_id, vpt.transaction_number, vpt.customer_id, vpt.vehicle_id,
			vpt.vin_number, vpt.vehicle_brand, vpt.vehicle_model, vpt.vehicle_year,
			vpt.vehicle_color, vpt.engine_number, vpt.registration_number, vpt.purchase_price,
			vpt.agreed_value, vpt.odometer_reading, vpt.fuel_type, vpt.transmission,
			vpt.condition_rating, vpt.purchase_date, vpt.transaction_status, vpt.inspection_notes,
			vpt.evaluation_notes, vpt.purchase_notes, vpt.documents_json, vpt.processed_by,
			vpt.inspected_by, vpt.approved_by, vpt.approved_at, vpt.created_at, vpt.updated_at,
			c.customer_name,
			u1.full_name as processed_by_name,
			u2.full_name as inspected_by_name,
			u3.full_name as approved_by_name
		FROM vehicle_purchase_transactions vpt
		LEFT JOIN customers c ON vpt.customer_id = c.customer_id
		LEFT JOIN users u1 ON vpt.processed_by = u1.user_id
		LEFT JOIN users u2 ON vpt.inspected_by = u2.user_id
		LEFT JOIN users u3 ON vpt.approved_by = u3.user_id
		WHERE vpt.transaction_number = $1`

	transaction := &vehicle_purchase.VehiclePurchaseTransaction{}
	err := r.db.QueryRowContext(ctx, query, number).Scan(
		&transaction.TransactionID,
		&transaction.TransactionNumber,
		&transaction.CustomerID,
		&transaction.VehicleID,
		&transaction.VinNumber,
		&transaction.VehicleBrand,
		&transaction.VehicleModel,
		&transaction.VehicleYear,
		&transaction.VehicleColor,
		&transaction.EngineNumber,
		&transaction.RegistrationNumber,
		&transaction.PurchasePrice,
		&transaction.AgreedValue,
		&transaction.OdometerReading,
		&transaction.FuelType,
		&transaction.Transmission,
		&transaction.ConditionRating,
		&transaction.PurchaseDate,
		&transaction.TransactionStatus,
		&transaction.InspectionNotes,
		&transaction.EvaluationNotes,
		&transaction.PurchaseNotes,
		&transaction.DocumentsJSON,
		&transaction.ProcessedBy,
		&transaction.InspectedBy,
		&transaction.ApprovedBy,
		&transaction.ApprovedAt,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.CustomerName,
		&transaction.ProcessedByName,
		&transaction.InspectedByName,
		&transaction.ApprovedByName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle purchase transaction not found")
		}
		return nil, fmt.Errorf("failed to get vehicle purchase transaction: %w", err)
	}

	return transaction, nil
}

// GetByVIN retrieves a vehicle purchase transaction by VIN number
func (r *VehiclePurchaseTransactionRepository) GetByVIN(ctx context.Context, vin string) (*vehicle_purchase.VehiclePurchaseTransaction, error) {
	// Similar implementation as GetByNumber but with vin_number
	query := `
		SELECT 
			vpt.transaction_id, vpt.transaction_number, vpt.customer_id, vpt.vehicle_id,
			vpt.vin_number, vpt.vehicle_brand, vpt.vehicle_model, vpt.vehicle_year,
			vpt.vehicle_color, vpt.engine_number, vpt.registration_number, vpt.purchase_price,
			vpt.agreed_value, vpt.odometer_reading, vpt.fuel_type, vpt.transmission,
			vpt.condition_rating, vpt.purchase_date, vpt.transaction_status, vpt.inspection_notes,
			vpt.evaluation_notes, vpt.purchase_notes, vpt.documents_json, vpt.processed_by,
			vpt.inspected_by, vpt.approved_by, vpt.approved_at, vpt.created_at, vpt.updated_at,
			c.customer_name,
			u1.full_name as processed_by_name,
			u2.full_name as inspected_by_name,
			u3.full_name as approved_by_name
		FROM vehicle_purchase_transactions vpt
		LEFT JOIN customers c ON vpt.customer_id = c.customer_id
		LEFT JOIN users u1 ON vpt.processed_by = u1.user_id
		LEFT JOIN users u2 ON vpt.inspected_by = u2.user_id
		LEFT JOIN users u3 ON vpt.approved_by = u3.user_id
		WHERE vpt.vin_number = $1`

	transaction := &vehicle_purchase.VehiclePurchaseTransaction{}
	err := r.db.QueryRowContext(ctx, query, vin).Scan(
		&transaction.TransactionID,
		&transaction.TransactionNumber,
		&transaction.CustomerID,
		&transaction.VehicleID,
		&transaction.VinNumber,
		&transaction.VehicleBrand,
		&transaction.VehicleModel,
		&transaction.VehicleYear,
		&transaction.VehicleColor,
		&transaction.EngineNumber,
		&transaction.RegistrationNumber,
		&transaction.PurchasePrice,
		&transaction.AgreedValue,
		&transaction.OdometerReading,
		&transaction.FuelType,
		&transaction.Transmission,
		&transaction.ConditionRating,
		&transaction.PurchaseDate,
		&transaction.TransactionStatus,
		&transaction.InspectionNotes,
		&transaction.EvaluationNotes,
		&transaction.PurchaseNotes,
		&transaction.DocumentsJSON,
		&transaction.ProcessedBy,
		&transaction.InspectedBy,
		&transaction.ApprovedBy,
		&transaction.ApprovedAt,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.CustomerName,
		&transaction.ProcessedByName,
		&transaction.InspectedByName,
		&transaction.ApprovedByName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle purchase transaction not found")
		}
		return nil, fmt.Errorf("failed to get vehicle purchase transaction: %w", err)
	}

	return transaction, nil
}

// Update updates a vehicle purchase transaction
func (r *VehiclePurchaseTransactionRepository) Update(ctx context.Context, id int, transaction *vehicle_purchase.VehiclePurchaseTransaction) (*vehicle_purchase.VehiclePurchaseTransaction, error) {
	query := `
		UPDATE vehicle_purchase_transactions SET
			vin_number = $1, vehicle_color = $2, engine_number = $3, registration_number = $4,
			purchase_price = $5, agreed_value = $6, odometer_reading = $7, condition_rating = $8,
			transaction_status = $9, inspection_notes = $10, evaluation_notes = $11,
			purchase_notes = $12, documents_json = $13, updated_at = NOW()
		WHERE transaction_id = $14
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		transaction.VinNumber,
		transaction.VehicleColor,
		transaction.EngineNumber,
		transaction.RegistrationNumber,
		transaction.PurchasePrice,
		transaction.AgreedValue,
		transaction.OdometerReading,
		transaction.ConditionRating,
		transaction.TransactionStatus,
		transaction.InspectionNotes,
		transaction.EvaluationNotes,
		transaction.PurchaseNotes,
		transaction.DocumentsJSON,
		id,
	).Scan(&transaction.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update vehicle purchase transaction: %w", err)
	}

	transaction.TransactionID = id
	return transaction, nil
}

// UpdateStatus updates the status of a vehicle purchase transaction
func (r *VehiclePurchaseTransactionRepository) UpdateStatus(ctx context.Context, id int, status string, updatedBy int) error {
	query := `
		UPDATE vehicle_purchase_transactions 
		SET transaction_status = $1, updated_at = NOW()
		WHERE transaction_id = $2`

	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	return nil
}

// Delete soft deletes a vehicle purchase transaction (you may want to implement proper soft delete)
func (r *VehiclePurchaseTransactionRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM vehicle_purchase_transactions WHERE transaction_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle purchase transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vehicle purchase transaction not found")
	}

	return nil
}

// GenerateNumber generates a unique transaction number
func (r *VehiclePurchaseTransactionRepository) GenerateNumber(ctx context.Context) (string, error) {
	// Get the current year and month
	now := time.Now()
	prefix := fmt.Sprintf("TXN%d%02d", now.Year(), now.Month())

	// Get the next sequence number for this month
	query := `
		SELECT COALESCE(MAX(CAST(SUBSTRING(transaction_number FROM LENGTH($1) + 1) AS INTEGER)), 0) + 1
		FROM vehicle_purchase_transactions 
		WHERE transaction_number LIKE $1 || '%'`

	var nextNum int
	err := r.db.QueryRowContext(ctx, query, prefix).Scan(&nextNum)
	if err != nil {
		return "", fmt.Errorf("failed to generate transaction number: %w", err)
	}

	return fmt.Sprintf("%s%04d", prefix, nextNum), nil
}

// IsNumberExists checks if a transaction number already exists
func (r *VehiclePurchaseTransactionRepository) IsNumberExists(ctx context.Context, number string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_purchase_transactions WHERE transaction_number = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, number).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if transaction number exists: %w", err)
	}

	return exists, nil
}

// IsVINExists checks if a VIN already exists (excluding a specific transaction)
func (r *VehiclePurchaseTransactionRepository) IsVINExists(ctx context.Context, vin string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_purchase_transactions WHERE vin_number = $1 AND transaction_id != $2)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, vin, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if VIN exists: %w", err)
	}

	return exists, nil
}

// Stub implementations for other methods - these would need full implementation
func (r *VehiclePurchaseTransactionRepository) List(ctx context.Context, params *vehicle_purchase.VehiclePurchaseTransactionFilterParams) (*common.PaginatedResponse, error) {
	// Implementation would include complex filtering, pagination, sorting
	return nil, fmt.Errorf("not implemented yet")
}

func (r *VehiclePurchaseTransactionRepository) GetByCustomerID(ctx context.Context, customerID int, params *vehicle_purchase.VehiclePurchaseTransactionFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *VehiclePurchaseTransactionRepository) GetByStatus(ctx context.Context, status string, params *vehicle_purchase.VehiclePurchaseTransactionFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *VehiclePurchaseTransactionRepository) GetPendingInspection(ctx context.Context, params *vehicle_purchase.VehiclePurchaseTransactionFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *VehiclePurchaseTransactionRepository) GetPendingApproval(ctx context.Context, params *vehicle_purchase.VehiclePurchaseTransactionFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *VehiclePurchaseTransactionRepository) CompleteInspection(ctx context.Context, id int, request *vehicle_purchase.TransactionInspectionRequest, inspectedBy int) error {
	query := `
		UPDATE vehicle_purchase_transactions 
		SET condition_rating = $1, inspection_notes = $2, evaluation_notes = $3,
			transaction_status = 'inspection', inspected_by = $4, updated_at = NOW()
		WHERE transaction_id = $5`

	_, err := r.db.ExecContext(ctx, query,
		request.ConditionRating,
		request.InspectionNotes,
		request.EvaluationNotes,
		inspectedBy,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to complete inspection: %w", err)
	}

	return nil
}

func (r *VehiclePurchaseTransactionRepository) ProcessApproval(ctx context.Context, id int, request *vehicle_purchase.TransactionStatusApprovalRequest, approvedBy int) error {
	query := `
		UPDATE vehicle_purchase_transactions 
		SET transaction_status = $1, approved_by = $2, approved_at = NOW(), updated_at = NOW()
		WHERE transaction_id = $3`

	_, err := r.db.ExecContext(ctx, query, request.Status, approvedBy, id)
	if err != nil {
		return fmt.Errorf("failed to process approval: %w", err)
	}

	return nil
}

func (r *VehiclePurchaseTransactionRepository) Search(ctx context.Context, query string, params *vehicle_purchase.VehiclePurchaseTransactionFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *VehiclePurchaseTransactionRepository) GetDashboardStats(ctx context.Context) (*vehicle_purchase.TransactionDashboardStats, error) {
	return nil, fmt.Errorf("not implemented yet")
}