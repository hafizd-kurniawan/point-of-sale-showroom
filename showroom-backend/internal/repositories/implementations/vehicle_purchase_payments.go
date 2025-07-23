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

// VehiclePurchasePaymentRepository implements interfaces.VehiclePurchasePaymentRepository
type VehiclePurchasePaymentRepository struct {
	db *sql.DB
}

// NewVehiclePurchasePaymentRepository creates a new vehicle purchase payment repository
func NewVehiclePurchasePaymentRepository(db *sql.DB) interfaces.VehiclePurchasePaymentRepository {
	return &VehiclePurchasePaymentRepository{db: db}
}

// Create creates a new vehicle purchase payment
func (r *VehiclePurchasePaymentRepository) Create(ctx context.Context, payment *vehicle_purchase.VehiclePurchasePayment) (*vehicle_purchase.VehiclePurchasePayment, error) {
	query := `
		INSERT INTO vehicle_purchase_payments (
			payment_number, transaction_id, payment_amount, payment_method, 
			payment_reference, bank_name, payment_status, payment_notes,
			due_date, processed_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING payment_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		payment.PaymentNumber,
		payment.TransactionID,
		payment.PaymentAmount,
		payment.PaymentMethod,
		payment.PaymentReference,
		payment.BankName,
		payment.PaymentStatus,
		payment.PaymentNotes,
		payment.DueDate,
		payment.ProcessedBy,
	).Scan(&payment.PaymentID, &payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle purchase payment: %w", err)
	}

	return payment, nil
}

// GetByID retrieves a vehicle purchase payment by ID
func (r *VehiclePurchasePaymentRepository) GetByID(ctx context.Context, id int) (*vehicle_purchase.VehiclePurchasePayment, error) {
	query := `
		SELECT 
			vpp.payment_id, vpp.payment_number, vpp.transaction_id, vpp.payment_amount,
			vpp.payment_method, vpp.payment_reference, vpp.bank_name, vpp.payment_status,
			vpp.payment_notes, vpp.due_date, vpp.processed_by, vpp.approved_by,
			vpp.processed_at, vpp.approved_at, vpp.created_at, vpp.updated_at,
			vpt.transaction_number,
			u1.full_name as processed_by_name,
			u2.full_name as approved_by_name
		FROM vehicle_purchase_payments vpp
		LEFT JOIN vehicle_purchase_transactions vpt ON vpp.transaction_id = vpt.transaction_id
		LEFT JOIN users u1 ON vpp.processed_by = u1.user_id
		LEFT JOIN users u2 ON vpp.approved_by = u2.user_id
		WHERE vpp.payment_id = $1`

	payment := &vehicle_purchase.VehiclePurchasePayment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.PaymentID,
		&payment.PaymentNumber,
		&payment.TransactionID,
		&payment.PaymentAmount,
		&payment.PaymentMethod,
		&payment.PaymentReference,
		&payment.BankName,
		&payment.PaymentStatus,
		&payment.PaymentNotes,
		&payment.DueDate,
		&payment.ProcessedBy,
		&payment.ApprovedBy,
		&payment.ProcessedAt,
		&payment.ApprovedAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
		&payment.TransactionNumber,
		&payment.ProcessedByName,
		&payment.ApprovedByName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle purchase payment not found")
		}
		return nil, fmt.Errorf("failed to get vehicle purchase payment: %w", err)
	}

	return payment, nil
}

// GetByNumber retrieves a vehicle purchase payment by payment number
func (r *VehiclePurchasePaymentRepository) GetByNumber(ctx context.Context, number string) (*vehicle_purchase.VehiclePurchasePayment, error) {
	query := `
		SELECT 
			vpp.payment_id, vpp.payment_number, vpp.transaction_id, vpp.payment_amount,
			vpp.payment_method, vpp.payment_reference, vpp.bank_name, vpp.payment_status,
			vpp.payment_notes, vpp.due_date, vpp.processed_by, vpp.approved_by,
			vpp.processed_at, vpp.approved_at, vpp.created_at, vpp.updated_at,
			vpt.transaction_number,
			u1.full_name as processed_by_name,
			u2.full_name as approved_by_name
		FROM vehicle_purchase_payments vpp
		LEFT JOIN vehicle_purchase_transactions vpt ON vpp.transaction_id = vpt.transaction_id
		LEFT JOIN users u1 ON vpp.processed_by = u1.user_id
		LEFT JOIN users u2 ON vpp.approved_by = u2.user_id
		WHERE vpp.payment_number = $1`

	payment := &vehicle_purchase.VehiclePurchasePayment{}
	err := r.db.QueryRowContext(ctx, query, number).Scan(
		&payment.PaymentID,
		&payment.PaymentNumber,
		&payment.TransactionID,
		&payment.PaymentAmount,
		&payment.PaymentMethod,
		&payment.PaymentReference,
		&payment.BankName,
		&payment.PaymentStatus,
		&payment.PaymentNotes,
		&payment.DueDate,
		&payment.ProcessedBy,
		&payment.ApprovedBy,
		&payment.ProcessedAt,
		&payment.ApprovedAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
		&payment.TransactionNumber,
		&payment.ProcessedByName,
		&payment.ApprovedByName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehicle purchase payment not found")
		}
		return nil, fmt.Errorf("failed to get vehicle purchase payment: %w", err)
	}

	return payment, nil
}

// Update updates a vehicle purchase payment
func (r *VehiclePurchasePaymentRepository) Update(ctx context.Context, id int, payment *vehicle_purchase.VehiclePurchasePayment) (*vehicle_purchase.VehiclePurchasePayment, error) {
	query := `
		UPDATE vehicle_purchase_payments SET
			payment_amount = $1, payment_method = $2, payment_reference = $3,
			bank_name = $4, payment_status = $5, payment_notes = $6,
			due_date = $7, updated_at = NOW()
		WHERE payment_id = $8
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		payment.PaymentAmount,
		payment.PaymentMethod,
		payment.PaymentReference,
		payment.BankName,
		payment.PaymentStatus,
		payment.PaymentNotes,
		payment.DueDate,
		id,
	).Scan(&payment.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update vehicle purchase payment: %w", err)
	}

	payment.PaymentID = id
	return payment, nil
}

// UpdateStatus updates the status of a vehicle purchase payment
func (r *VehiclePurchasePaymentRepository) UpdateStatus(ctx context.Context, id int, status string, updatedBy int) error {
	query := `
		UPDATE vehicle_purchase_payments 
		SET payment_status = $1, updated_at = NOW()
		WHERE payment_id = $2`

	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	return nil
}

// Delete soft deletes a vehicle purchase payment
func (r *VehiclePurchasePaymentRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM vehicle_purchase_payments WHERE payment_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle purchase payment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vehicle purchase payment not found")
	}

	return nil
}

// GenerateNumber generates a unique payment number
func (r *VehiclePurchasePaymentRepository) GenerateNumber(ctx context.Context) (string, error) {
	now := time.Now()
	prefix := fmt.Sprintf("PAY%d%02d", now.Year(), now.Month())

	query := `
		SELECT COALESCE(MAX(CAST(SUBSTRING(payment_number FROM LENGTH($1) + 1) AS INTEGER)), 0) + 1
		FROM vehicle_purchase_payments 
		WHERE payment_number LIKE $1 || '%'`

	var nextNum int
	err := r.db.QueryRowContext(ctx, query, prefix).Scan(&nextNum)
	if err != nil {
		return "", fmt.Errorf("failed to generate payment number: %w", err)
	}

	return fmt.Sprintf("%s%04d", prefix, nextNum), nil
}

// IsNumberExists checks if a payment number already exists
func (r *VehiclePurchasePaymentRepository) IsNumberExists(ctx context.Context, number string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicle_purchase_payments WHERE payment_number = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, number).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if payment number exists: %w", err)
	}

	return exists, nil
}

// ProcessPayment processes a payment
func (r *VehiclePurchasePaymentRepository) ProcessPayment(ctx context.Context, id int, request *vehicle_purchase.PaymentProcessRequest, processedBy int) error {
	query := `
		UPDATE vehicle_purchase_payments 
		SET payment_status = $1, processed_by = $2, processed_at = NOW(), updated_at = NOW()
		WHERE payment_id = $3`

	_, err := r.db.ExecContext(ctx, query, request.Status, processedBy, id)
	if err != nil {
		return fmt.Errorf("failed to process payment: %w", err)
	}

	return nil
}

// ProcessApproval processes payment approval
func (r *VehiclePurchasePaymentRepository) ProcessApproval(ctx context.Context, id int, request *vehicle_purchase.PaymentApprovalRequest, approvedBy int) error {
	query := `
		UPDATE vehicle_purchase_payments 
		SET payment_status = $1, approved_by = $2, approved_at = NOW(), updated_at = NOW()
		WHERE payment_id = $3`

	_, err := r.db.ExecContext(ctx, query, request.Status, approvedBy, id)
	if err != nil {
		return fmt.Errorf("failed to process approval: %w", err)
	}

	return nil
}

// GetPaymentSummary gets payment summary for a transaction
func (r *VehiclePurchasePaymentRepository) GetPaymentSummary(ctx context.Context, transactionID int) (*vehicle_purchase.PaymentSummary, error) {
	query := `
		SELECT 
			COUNT(*) as total_payments,
			COALESCE(SUM(payment_amount), 0) as total_amount,
			COALESCE(SUM(CASE WHEN payment_status = 'completed' THEN payment_amount ELSE 0 END), 0) as paid_amount,
			COALESCE(SUM(CASE WHEN payment_status = 'pending' THEN payment_amount ELSE 0 END), 0) as pending_amount
		FROM vehicle_purchase_payments 
		WHERE transaction_id = $1`

	summary := &vehicle_purchase.PaymentSummary{}
	err := r.db.QueryRowContext(ctx, query, transactionID).Scan(
		&summary.TotalPayments,
		&summary.TotalAmount,
		&summary.PaidAmount,
		&summary.PendingAmount,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get payment summary: %w", err)
	}

	summary.TransactionID = transactionID
	summary.RemainingAmount = summary.TotalAmount - summary.PaidAmount

	return summary, nil
}

// CalculatePaymentTotals calculates payment totals for a transaction
func (r *VehiclePurchasePaymentRepository) CalculatePaymentTotals(ctx context.Context, transactionID int) (*vehicle_purchase.PaymentSummary, error) {
	return r.GetPaymentSummary(ctx, transactionID)
}

// Stub implementations for list methods
func (r *VehiclePurchasePaymentRepository) List(ctx context.Context, params *vehicle_purchase.VehiclePurchasePaymentFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *VehiclePurchasePaymentRepository) GetByTransactionID(ctx context.Context, transactionID int, params *vehicle_purchase.VehiclePurchasePaymentFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *VehiclePurchasePaymentRepository) GetByStatus(ctx context.Context, status string, params *vehicle_purchase.VehiclePurchasePaymentFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *VehiclePurchasePaymentRepository) GetPendingApproval(ctx context.Context, params *vehicle_purchase.VehiclePurchasePaymentFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r *VehiclePurchasePaymentRepository) GetOverduePayments(ctx context.Context, params *vehicle_purchase.VehiclePurchasePaymentFilterParams) (*common.PaginatedResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}