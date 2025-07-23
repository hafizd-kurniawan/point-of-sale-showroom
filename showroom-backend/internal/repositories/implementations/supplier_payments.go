package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// SupplierPaymentRepository implements interfaces.SupplierPaymentRepository
type SupplierPaymentRepository struct {
	db *sql.DB
}

// NewSupplierPaymentRepository creates a new supplier payment repository
func NewSupplierPaymentRepository(db *sql.DB) interfaces.SupplierPaymentRepository {
	return &SupplierPaymentRepository{db: db}
}

// Create creates a new supplier payment
func (r *SupplierPaymentRepository) Create(ctx context.Context, payment *products.SupplierPayment) (*products.SupplierPayment, error) {
	// Calculate outstanding amount and days overdue
	payment.OutstandingAmount = payment.InvoiceAmount - payment.PaymentAmount - payment.DiscountTaken
	
	// Calculate days overdue
	now := time.Now()
	if now.After(payment.DueDate) {
		payment.DaysOverdue = int(now.Sub(payment.DueDate).Hours() / 24)
	} else {
		payment.DaysOverdue = 0
	}

	// Set payment status based on amounts
	if payment.PaymentAmount == 0 {
		if payment.DaysOverdue > 0 {
			payment.PaymentStatus = products.PaymentStatusOverdue
		} else {
			payment.PaymentStatus = products.PaymentStatusPending
		}
	} else if payment.OutstandingAmount > 0 {
		payment.PaymentStatus = products.PaymentStatusPartial
	} else {
		payment.PaymentStatus = products.PaymentStatusPaid
	}

	query := `
		INSERT INTO supplier_payments (
			supplier_id, po_id, payment_number, invoice_amount, payment_amount,
			discount_taken, outstanding_amount, invoice_date, payment_date,
			due_date, payment_method, payment_reference, invoice_number,
			payment_status, days_overdue, penalty_amount, processed_by, payment_notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING payment_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		payment.SupplierID,
		payment.POID,
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
	).Scan(&payment.PaymentID, &payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create supplier payment: %w", err)
	}

	return payment, nil
}

// GetByID retrieves a supplier payment by ID
func (r *SupplierPaymentRepository) GetByID(ctx context.Context, id int) (*products.SupplierPayment, error) {
	query := `
		SELECT payment_id, supplier_id, po_id, payment_number, invoice_amount,
			   payment_amount, discount_taken, outstanding_amount, invoice_date,
			   payment_date, due_date, payment_method, payment_reference,
			   invoice_number, payment_status, days_overdue, penalty_amount,
			   processed_by, payment_notes, created_at, updated_at
		FROM supplier_payments 
		WHERE payment_id = $1`

	payment := &products.SupplierPayment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.PaymentID,
		&payment.SupplierID,
		&payment.POID,
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("supplier payment not found")
		}
		return nil, fmt.Errorf("failed to get supplier payment: %w", err)
	}

	return payment, nil
}

// GetByNumber retrieves a supplier payment by payment number
func (r *SupplierPaymentRepository) GetByNumber(ctx context.Context, number string) (*products.SupplierPayment, error) {
	query := `
		SELECT payment_id, supplier_id, po_id, payment_number, invoice_amount,
			   payment_amount, discount_taken, outstanding_amount, invoice_date,
			   payment_date, due_date, payment_method, payment_reference,
			   invoice_number, payment_status, days_overdue, penalty_amount,
			   processed_by, payment_notes, created_at, updated_at
		FROM supplier_payments 
		WHERE payment_number = $1`

	payment := &products.SupplierPayment{}
	err := r.db.QueryRowContext(ctx, query, number).Scan(
		&payment.PaymentID,
		&payment.SupplierID,
		&payment.POID,
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("supplier payment not found")
		}
		return nil, fmt.Errorf("failed to get supplier payment: %w", err)
	}

	return payment, nil
}

// Update updates a supplier payment
func (r *SupplierPaymentRepository) Update(ctx context.Context, id int, payment *products.SupplierPayment) (*products.SupplierPayment, error) {
	// Recalculate outstanding amount and status
	payment.OutstandingAmount = payment.InvoiceAmount - payment.PaymentAmount - payment.DiscountTaken

	// Calculate days overdue
	now := time.Now()
	if now.After(payment.DueDate) {
		payment.DaysOverdue = int(now.Sub(payment.DueDate).Hours() / 24)
	} else {
		payment.DaysOverdue = 0
	}

	// Update payment status based on amounts
	if payment.PaymentAmount == 0 {
		if payment.DaysOverdue > 0 {
			payment.PaymentStatus = products.PaymentStatusOverdue
		} else {
			payment.PaymentStatus = products.PaymentStatusPending
		}
	} else if payment.OutstandingAmount > 0 {
		payment.PaymentStatus = products.PaymentStatusPartial
	} else {
		payment.PaymentStatus = products.PaymentStatusPaid
	}

	query := `
		UPDATE supplier_payments 
		SET invoice_amount = $1, payment_amount = $2, discount_taken = $3,
			outstanding_amount = $4, invoice_date = $5, payment_date = $6,
			due_date = $7, payment_method = $8, payment_reference = $9,
			invoice_number = $10, payment_status = $11, days_overdue = $12,
			penalty_amount = $13, payment_notes = $14, updated_at = NOW()
		WHERE payment_id = $15`

	_, err := r.db.ExecContext(ctx, query,
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
		payment.PaymentNotes,
		id,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update supplier payment: %w", err)
	}

	return r.GetByID(ctx, id)
}

// Delete deletes a supplier payment
func (r *SupplierPaymentRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM supplier_payments WHERE payment_id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete supplier payment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("supplier payment not found")
	}

	return nil
}

// List retrieves all supplier payments with pagination
func (r *SupplierPaymentRepository) List(ctx context.Context, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error) {
	baseQuery := `
		FROM supplier_payments sp 
		LEFT JOIN supplier s ON sp.supplier_id = s.supplier_id
		WHERE 1=1`
	
	args := []interface{}{}
	whereConditions := []string{}
	argIndex := 1

	// Add filters
	if params.SupplierID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sp.supplier_id = $%d", argIndex))
		args = append(args, *params.SupplierID)
		argIndex++
	}

	if params.PaymentStatus != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sp.payment_status = $%d", argIndex))
		args = append(args, *params.PaymentStatus)
		argIndex++
	}

	if params.PaymentMethod != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sp.payment_method = $%d", argIndex))
		args = append(args, *params.PaymentMethod)
		argIndex++
	}

	if params.POID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sp.po_id = $%d", argIndex))
		args = append(args, *params.POID)
		argIndex++
	}

	if params.DateFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sp.payment_date >= $%d", argIndex))
		args = append(args, *params.DateFrom)
		argIndex++
	}

	if params.DateTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("sp.payment_date <= $%d", argIndex))
		args = append(args, *params.DateTo)
		argIndex++
	}

	if params.IsOverdue != nil && *params.IsOverdue {
		whereConditions = append(whereConditions, "sp.payment_status = 'overdue'")
	}

	if params.Search != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("(sp.payment_number ILIKE $%d OR sp.invoice_number ILIKE $%d OR s.supplier_name ILIKE $%d)", argIndex, argIndex, argIndex))
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	if len(whereConditions) > 0 {
		baseQuery += " AND " + strings.Join(whereConditions, " AND ")
	}

	// Count query
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count supplier payments: %w", err)
	}

	// Calculate pagination
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	offset := (params.Page - 1) * params.Limit

	// Main query
	selectFields := `
		sp.payment_id, sp.supplier_id, sp.po_id, sp.payment_number,
		sp.invoice_amount, sp.payment_amount, sp.outstanding_amount,
		sp.invoice_date, sp.payment_date, sp.due_date, sp.payment_method,
		sp.invoice_number, sp.payment_status, sp.days_overdue`
	
	mainQuery := "SELECT " + selectFields + " " + baseQuery + 
		" ORDER BY sp.payment_date DESC, sp.payment_id DESC LIMIT $" + fmt.Sprintf("%d", argIndex) + 
		" OFFSET $" + fmt.Sprintf("%d", argIndex+1)
	
	args = append(args, params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, mainQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query supplier payments: %w", err)
	}
	defer rows.Close()

	var payments []products.SupplierPaymentListItem
	for rows.Next() {
		var payment products.SupplierPaymentListItem
		err := rows.Scan(
			&payment.PaymentID,
			&payment.SupplierID,
			&payment.POID,
			&payment.PaymentNumber,
			&payment.InvoiceAmount,
			&payment.PaymentAmount,
			&payment.OutstandingAmount,
			&payment.InvoiceDate,
			&payment.PaymentDate,
			&payment.DueDate,
			&payment.PaymentMethod,
			&payment.InvoiceNumber,
			&payment.PaymentStatus,
			&payment.DaysOverdue,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan supplier payment: %w", err)
		}
		payments = append(payments, payment)
	}

	totalPages := (total + int64(params.Limit) - 1) / int64(params.Limit)

	return &common.PaginatedResponse{
		Data:       payments,
		Total:      int(total),
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: int(totalPages),
		HasMore:    params.Page < int(totalPages),
	}, nil
}

// GetBySupplierID retrieves supplier payments for a specific supplier
func (r *SupplierPaymentRepository) GetBySupplierID(ctx context.Context, supplierID int, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error) {
	// Set supplierID in params and call List
	params.SupplierID = &supplierID
	return r.List(ctx, params)
}

// GetByPOID retrieves supplier payments for a specific PO
func (r *SupplierPaymentRepository) GetByPOID(ctx context.Context, poID int, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error) {
	// Set POID in params and call List
	params.POID = &poID
	return r.List(ctx, params)
}

// GetOverduePayments retrieves overdue payments
func (r *SupplierPaymentRepository) GetOverduePayments(ctx context.Context, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error) {
	isOverdue := true
	params.IsOverdue = &isOverdue
	return r.List(ctx, params)
}

// UpdatePaymentStatus updates the payment status
func (r *SupplierPaymentRepository) UpdatePaymentStatus(ctx context.Context, id int, status products.PaymentStatus) error {
	query := `UPDATE supplier_payments SET payment_status = $1, updated_at = NOW() WHERE payment_id = $2`
	
	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("supplier payment not found")
	}

	return nil
}

// AddPayment adds a payment amount to an existing payment record
func (r *SupplierPaymentRepository) AddPayment(ctx context.Context, id int, amount float64, method products.PaymentMethod, reference *string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get current payment
	payment, err := r.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	// Update payment amount and outstanding amount
	newPaymentAmount := payment.PaymentAmount + amount
	newOutstandingAmount := payment.InvoiceAmount - newPaymentAmount - payment.DiscountTaken

	// Update status based on new amounts
	var newStatus products.PaymentStatus
	if newOutstandingAmount <= 0 {
		newStatus = products.PaymentStatusPaid
	} else {
		newStatus = products.PaymentStatusPartial
	}

	query := `
		UPDATE supplier_payments 
		SET payment_amount = $1, outstanding_amount = $2, payment_status = $3,
			payment_method = $4, payment_reference = $5, payment_date = NOW(), updated_at = NOW()
		WHERE payment_id = $6`

	_, err = tx.ExecContext(ctx, query,
		newPaymentAmount,
		newOutstandingAmount,
		newStatus,
		method,
		reference,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GenerateNumber generates a new payment number
func (r *SupplierPaymentRepository) GenerateNumber(ctx context.Context) (string, error) {
	// Generate payment number with format PAY-YYYYMMDD-XXXX
	now := time.Now()
	dateStr := now.Format("20060102")
	
	query := `
		SELECT COALESCE(MAX(
			CAST(SUBSTRING(payment_number FROM 'PAY-\d{8}-(\d+)') AS INTEGER)
		), 0) + 1
		FROM supplier_payments 
		WHERE payment_number LIKE $1`
	
	prefix := "PAY-" + dateStr + "-%"
	var nextNumber int
	err := r.db.QueryRowContext(ctx, query, prefix).Scan(&nextNumber)
	if err != nil {
		return "", fmt.Errorf("failed to generate payment number: %w", err)
	}

	return fmt.Sprintf("PAY-%s-%04d", dateStr, nextNumber), nil
}

// IsNumberExists checks if a payment number already exists
func (r *SupplierPaymentRepository) IsNumberExists(ctx context.Context, number string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM supplier_payments WHERE payment_number = $1)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, number).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check payment number existence: %w", err)
	}

	return exists, nil
}

// UpdateOverdueStatus updates overdue status for all payments
func (r *SupplierPaymentRepository) UpdateOverdueStatus(ctx context.Context) error {
	query := `
		UPDATE supplier_payments 
		SET payment_status = 'overdue',
			days_overdue = EXTRACT(DAY FROM (NOW() - due_date))::int,
			updated_at = NOW()
		WHERE due_date < NOW() 
		AND payment_status IN ('pending', 'partial')
		AND outstanding_amount > 0`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to update overdue status: %w", err)
	}

	return nil
}

// GetPaymentSummary gets payment summary for a supplier or all suppliers
func (r *SupplierPaymentRepository) GetPaymentSummary(ctx context.Context, supplierID *int) (map[string]interface{}, error) {
	baseQuery := `
		SELECT 
			COUNT(*) as total_payments,
			COALESCE(SUM(invoice_amount), 0) as total_invoiced,
			COALESCE(SUM(payment_amount), 0) as total_paid,
			COALESCE(SUM(outstanding_amount), 0) as total_outstanding,
			COUNT(CASE WHEN payment_status = 'overdue' THEN 1 END) as overdue_count,
			COALESCE(SUM(CASE WHEN payment_status = 'overdue' THEN outstanding_amount ELSE 0 END), 0) as overdue_amount
		FROM supplier_payments 
		WHERE 1=1`

	args := []interface{}{}
	if supplierID != nil {
		baseQuery += " AND supplier_id = $1"
		args = append(args, *supplierID)
	}

	var totalPayments int
	var totalInvoiced, totalPaid, totalOutstanding, overdueAmount float64
	var overdueCount int

	err := r.db.QueryRowContext(ctx, baseQuery, args...).Scan(
		&totalPayments,
		&totalInvoiced,
		&totalPaid,
		&totalOutstanding,
		&overdueCount,
		&overdueAmount,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment summary: %w", err)
	}

	return map[string]interface{}{
		"total_payments":     totalPayments,
		"total_invoiced":     totalInvoiced,
		"total_paid":         totalPaid,
		"total_outstanding":  totalOutstanding,
		"overdue_count":      overdueCount,
		"overdue_amount":     overdueAmount,
		"payment_rate":       totalPaid / totalInvoiced * 100,
	}, nil
}