package implementations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

type supplierPaymentRepository struct {
	db *sql.DB
}

// NewSupplierPaymentRepository creates a new supplier payment repository
func NewSupplierPaymentRepository(db *sql.DB) interfaces.SupplierPaymentRepository {
	return &supplierPaymentRepository{
		db: db,
	}
}

// Create creates a new supplier payment
func (r *supplierPaymentRepository) Create(ctx context.Context, payment *inventory.SupplierPayment) (*inventory.SupplierPayment, error) {
	query := `
		INSERT INTO supplier_payments (
			supplier_id, po_id, payment_number, invoice_amount, payment_amount,
			discount_taken, outstanding_amount, invoice_date, payment_date,
			due_date, payment_method, payment_reference, invoice_number,
			payment_status, days_overdue, penalty_amount, processed_by, payment_notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING payment_id`

	err := r.db.QueryRowContext(ctx, query,
		payment.SupplierID,
		payment.PoID,
		payment.PaymentNumber,
		payment.InvoiceAmount,
		payment.PaymentAmount,
		payment.DiscountTaken,
		payment.OutstandingAmount,
		payment.InvoiceDate,
		payment.PaymentDate,
		payment.DueDate,
		payment.PaymentMethod,
		payment.PaymentReference,
		payment.InvoiceNumber,
		payment.PaymentStatus,
		payment.DaysOverdue,
		payment.PenaltyAmount,
		payment.ProcessedBy,
		payment.PaymentNotes,
	).Scan(&payment.PaymentID)

	if err != nil {
		return nil, fmt.Errorf("failed to create supplier payment: %w", err)
	}

	return payment, nil
}

// GetByID retrieves a supplier payment by ID
func (r *supplierPaymentRepository) GetByID(ctx context.Context, id int) (*inventory.SupplierPayment, error) {
	query := `
		SELECT 
			sp.payment_id, sp.supplier_id, sp.po_id, sp.payment_number,
			sp.invoice_amount, sp.payment_amount, sp.discount_taken,
			sp.outstanding_amount, sp.invoice_date, sp.payment_date,
			sp.due_date, sp.payment_method, sp.payment_reference,
			sp.invoice_number, sp.payment_status, sp.days_overdue,
			sp.penalty_amount, sp.processed_by, sp.payment_notes,
			sp.created_at, sp.updated_at
		FROM supplier_payments sp
		WHERE sp.payment_id = $1`

	payment := &inventory.SupplierPayment{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.PaymentID,
		&payment.SupplierID,
		&payment.PoID,
		&payment.PaymentNumber,
		&payment.InvoiceAmount,
		&payment.PaymentAmount,
		&payment.DiscountTaken,
		&payment.OutstandingAmount,
		&payment.InvoiceDate,
		&payment.PaymentDate,
		&payment.DueDate,
		&payment.PaymentMethod,
		&payment.PaymentReference,
		&payment.InvoiceNumber,
		&payment.PaymentStatus,
		&payment.DaysOverdue,
		&payment.PenaltyAmount,
		&payment.ProcessedBy,
		&payment.PaymentNotes,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get supplier payment: %w", err)
	}

	return payment, nil
}

// GetByNumber retrieves a supplier payment by payment number
func (r *supplierPaymentRepository) GetByNumber(ctx context.Context, number string) (*inventory.SupplierPayment, error) {
	query := `
		SELECT 
			sp.payment_id, sp.supplier_id, sp.po_id, sp.payment_number,
			sp.invoice_amount, sp.payment_amount, sp.discount_taken,
			sp.outstanding_amount, sp.invoice_date, sp.payment_date,
			sp.due_date, sp.payment_method, sp.payment_reference,
			sp.invoice_number, sp.payment_status, sp.days_overdue,
			sp.penalty_amount, sp.processed_by, sp.payment_notes,
			sp.created_at, sp.updated_at
		FROM supplier_payments sp
		WHERE sp.payment_number = $1`

	payment := &inventory.SupplierPayment{}

	err := r.db.QueryRowContext(ctx, query, number).Scan(
		&payment.PaymentID,
		&payment.SupplierID,
		&payment.PoID,
		&payment.PaymentNumber,
		&payment.InvoiceAmount,
		&payment.PaymentAmount,
		&payment.DiscountTaken,
		&payment.OutstandingAmount,
		&payment.InvoiceDate,
		&payment.PaymentDate,
		&payment.DueDate,
		&payment.PaymentMethod,
		&payment.PaymentReference,
		&payment.InvoiceNumber,
		&payment.PaymentStatus,
		&payment.DaysOverdue,
		&payment.PenaltyAmount,
		&payment.ProcessedBy,
		&payment.PaymentNotes,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get supplier payment by number: %w", err)
	}

	return payment, nil
}

// Update updates a supplier payment
func (r *supplierPaymentRepository) Update(ctx context.Context, id int, payment *inventory.SupplierPayment) (*inventory.SupplierPayment, error) {
	query := `
		UPDATE supplier_payments 
		SET invoice_amount = $2, payment_amount = $3, discount_taken = $4,
			outstanding_amount = $5, invoice_date = $6, payment_date = $7,
			due_date = $8, payment_method = $9, payment_reference = $10,
			invoice_number = $11, payment_status = $12, days_overdue = $13,
			penalty_amount = $14, processed_by = $15, payment_notes = $16
		WHERE payment_id = $1`

	_, err := r.db.ExecContext(ctx, query,
		id,
		payment.InvoiceAmount,
		payment.PaymentAmount,
		payment.DiscountTaken,
		payment.OutstandingAmount,
		payment.InvoiceDate,
		payment.PaymentDate,
		payment.DueDate,
		payment.PaymentMethod,
		payment.PaymentReference,
		payment.InvoiceNumber,
		payment.PaymentStatus,
		payment.DaysOverdue,
		payment.PenaltyAmount,
		payment.ProcessedBy,
		payment.PaymentNotes,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update supplier payment: %w", err)
	}

	payment.PaymentID = id
	return payment, nil
}

// Delete soft deletes a supplier payment
func (r *supplierPaymentRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM supplier_payments WHERE payment_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete supplier payment: %w", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("supplier payment with ID %d not found", id)
	}

	return nil
}

// List retrieves supplier payments with filtering and pagination
func (r *supplierPaymentRepository) List(ctx context.Context, params *inventory.SupplierPaymentFilterParams) ([]inventory.SupplierPaymentListItem, int, error) {
	query := `
		SELECT 
			sp.payment_id, sp.payment_number, '' as supplier_name,
			NULL as po_number, sp.invoice_number, sp.invoice_amount, 
			sp.payment_amount, sp.outstanding_amount, sp.payment_date,
			sp.due_date, sp.payment_method, sp.payment_status,
			sp.days_overdue, '' as processed_by_name
		FROM supplier_payments sp
		ORDER BY sp.payment_date DESC
		LIMIT $1 OFFSET $2`

	offset := (params.Page - 1) * params.Limit

	rows, err := r.db.QueryContext(ctx, query, params.Limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get supplier payments: %w", err)
	}
	defer rows.Close()

	var items []inventory.SupplierPaymentListItem
	for rows.Next() {
		item := inventory.SupplierPaymentListItem{}

		err := rows.Scan(
			&item.PaymentID,
			&item.PaymentNumber,
			&item.SupplierName,
			&item.PoNumber,
			&item.InvoiceNumber,
			&item.InvoiceAmount,
			&item.PaymentAmount,
			&item.OutstandingAmount,
			&item.PaymentDate,
			&item.DueDate,
			&item.PaymentMethod,
			&item.PaymentStatus,
			&item.DaysOverdue,
			&item.ProcessedByName,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan supplier payment list item: %w", err)
		}

		items = append(items, item)
	}

	// Count total
	countQuery := `SELECT COUNT(*) FROM supplier_payments`
	var total int
	err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count supplier payments: %w", err)
	}

	return items, total, nil
}

// GetBySupplier retrieves supplier payments for a specific supplier
func (r *supplierPaymentRepository) GetBySupplier(ctx context.Context, supplierID int, page, limit int) ([]inventory.SupplierPaymentListItem, int, error) {
	params := &inventory.SupplierPaymentFilterParams{
		SupplierID: &supplierID,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// GetByPO retrieves supplier payments for a specific purchase order
func (r *supplierPaymentRepository) GetByPO(ctx context.Context, poID int, page, limit int) ([]inventory.SupplierPaymentListItem, int, error) {
	params := &inventory.SupplierPaymentFilterParams{
		PoID: &poID,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// GetByStatus retrieves supplier payments by status
func (r *supplierPaymentRepository) GetByStatus(ctx context.Context, status inventory.PaymentStatus, page, limit int) ([]inventory.SupplierPaymentListItem, int, error) {
	params := &inventory.SupplierPaymentFilterParams{
		PaymentStatus: &status,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// GetOverduePayments retrieves overdue payments
func (r *supplierPaymentRepository) GetOverduePayments(ctx context.Context, page, limit int) ([]inventory.SupplierPaymentListItem, int, error) {
	overdueOnly := true
	params := &inventory.SupplierPaymentFilterParams{
		IsOverdue: &overdueOnly,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// Search searches supplier payments
func (r *supplierPaymentRepository) Search(ctx context.Context, query string, page, limit int) ([]inventory.SupplierPaymentListItem, int, error) {
	params := &inventory.SupplierPaymentFilterParams{
		Search: query,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// ExistsByNumber checks if a supplier payment exists by payment number
func (r *supplierPaymentRepository) ExistsByNumber(ctx context.Context, number string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM supplier_payments WHERE payment_number = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, number).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check supplier payment existence: %w", err)
	}
	return exists, nil
}

// ExistsByNumberExcludingID checks if a supplier payment exists by payment number excluding a specific ID
func (r *supplierPaymentRepository) ExistsByNumberExcludingID(ctx context.Context, number string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM supplier_payments WHERE payment_number = $1 AND payment_id != $2)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, number, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check supplier payment existence: %w", err)
	}
	return exists, nil
}

// GetLastPaymentID gets the last payment ID for code generation
func (r *supplierPaymentRepository) GetLastPaymentID(ctx context.Context) (int, error) {
	query := `SELECT COALESCE(MAX(payment_id), 0) FROM supplier_payments`
	var lastID int
	err := r.db.QueryRowContext(ctx, query).Scan(&lastID)
	if err != nil {
		return 0, fmt.Errorf("failed to get last payment ID: %w", err)
	}
	return lastID, nil
}

// UpdateStatus updates the payment status
func (r *supplierPaymentRepository) UpdateStatus(ctx context.Context, id int, status inventory.PaymentStatus) error {
	query := `
		UPDATE supplier_payments 
		SET payment_status = $2
		WHERE payment_id = $1`

	result, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("supplier payment with ID %d not found", id)
	}

	return nil
}

// UpdateOverdueStatus updates overdue status for all payments
func (r *supplierPaymentRepository) UpdateOverdueStatus(ctx context.Context) error {
	query := `
		UPDATE supplier_payments 
		SET days_overdue = CASE 
			WHEN due_date < CURRENT_DATE AND payment_status IN ('pending', 'partial') 
			THEN EXTRACT(DAY FROM CURRENT_DATE - due_date)::int
			ELSE 0 
		END,
		payment_status = CASE 
			WHEN due_date < CURRENT_DATE AND payment_status = 'pending' 
			THEN 'overdue'
			ELSE payment_status 
		END`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to update overdue status: %w", err)
	}

	return nil
}

// GetPaymentSummary retrieves payment summary for reporting
func (r *supplierPaymentRepository) GetPaymentSummary(ctx context.Context, supplierID *int, startDate, endDate string) (*inventory.PaymentSummary, error) {
	query := `
		SELECT 
			COALESCE(SUM(invoice_amount), 0) as total_invoice_amount,
			COALESCE(SUM(payment_amount), 0) as total_payment_amount,
			COALESCE(SUM(discount_taken), 0) as total_discount_taken,
			COALESCE(SUM(outstanding_amount), 0) as total_outstanding_amount,
			COALESCE(SUM(penalty_amount), 0) as total_penalty_amount,
			COUNT(CASE WHEN payment_status = 'pending' THEN 1 END) as pending_payments_count,
			COUNT(CASE WHEN payment_status = 'overdue' THEN 1 END) as overdue_payments_count,
			0.0 as average_payment_days
		FROM supplier_payments 
		WHERE payment_date >= $1::date 
		AND payment_date <= $2::date`

	args := []interface{}{startDate, endDate}
	if supplierID != nil {
		query += " AND supplier_id = $3"
		args = append(args, *supplierID)
	}

	summary := &inventory.PaymentSummary{}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&summary.TotalInvoiceAmount,
		&summary.TotalPaymentAmount,
		&summary.TotalDiscountTaken,
		&summary.TotalOutstandingAmount,
		&summary.TotalPenaltyAmount,
		&summary.PendingPaymentsCount,
		&summary.OverduePaymentsCount,
		&summary.AveragePaymentDays,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get payment summary: %w", err)
	}

	return summary, nil
}

// GetOutstandingBalance retrieves outstanding balance for a supplier
func (r *supplierPaymentRepository) GetOutstandingBalance(ctx context.Context, supplierID int) (float64, error) {
	query := `
		SELECT COALESCE(SUM(outstanding_amount), 0) 
		FROM supplier_payments 
		WHERE supplier_id = $1 AND payment_status IN ('pending', 'partial', 'overdue')`

	var balance float64
	err := r.db.QueryRowContext(ctx, query, supplierID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("failed to get outstanding balance: %w", err)
	}

	return balance, nil
}

// GetTotalPaid retrieves total paid amount for a supplier within date range
func (r *supplierPaymentRepository) GetTotalPaid(ctx context.Context, supplierID int, startDate, endDate string) (float64, error) {
	query := `
		SELECT COALESCE(SUM(payment_amount), 0) 
		FROM supplier_payments 
		WHERE supplier_id = $1 
		AND payment_status = 'paid' 
		AND payment_date >= $2::date 
		AND payment_date <= $3::date`

	var totalPaid float64
	err := r.db.QueryRowContext(ctx, query, supplierID, startDate, endDate).Scan(&totalPaid)
	if err != nil {
		return 0, fmt.Errorf("failed to get total paid: %w", err)
	}

	return totalPaid, nil
}