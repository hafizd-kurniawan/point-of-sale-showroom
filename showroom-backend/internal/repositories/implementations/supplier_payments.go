package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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
	query := `
		INSERT INTO supplier_payments (
			supplier_id, po_id, payment_number, invoice_amount, payment_amount,
			discount_taken, outstanding_amount, invoice_date, payment_date, due_date,
			payment_method, payment_reference, invoice_number, payment_status,
			days_overdue, penalty_amount, processed_by, payment_notes
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
	query := `
		UPDATE supplier_payments SET
			payment_amount = $2, discount_taken = $3, outstanding_amount = $4,
			payment_date = $5, due_date = $6, payment_method = $7,
			payment_reference = $8, payment_status = $9, days_overdue = $10,
			penalty_amount = $11, payment_notes = $12, updated_at = NOW()
		WHERE payment_id = $1`

	_, err := r.db.ExecContext(ctx, query,
		id,
		payment.PaymentAmount,
		payment.DiscountTaken,
		payment.OutstandingAmount,
		payment.PaymentDate,
		payment.DueDate,
		payment.PaymentMethod,
		payment.PaymentReference,
		payment.PaymentStatus,
		payment.DaysOverdue,
		payment.PenaltyAmount,
		payment.PaymentNotes,
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
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("supplier payment not found")
	}

	return nil
}

// List retrieves supplier payments with pagination and filtering
func (r *SupplierPaymentRepository) List(ctx context.Context, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error) {
	var args []interface{}
	var conditions []string
	
	// Base query
	baseQuery := `
		FROM supplier_payments sp
		LEFT JOIN suppliers s ON sp.supplier_id = s.supplier_id
		LEFT JOIN purchase_orders_parts po ON sp.po_id = po.po_id
		WHERE 1=1`
	
	argIndex := 1

	// Apply filters
	if params.SupplierID != nil {
		conditions = append(conditions, fmt.Sprintf("sp.supplier_id = $%d", argIndex))
		args = append(args, *params.SupplierID)
		argIndex++
	}

	if params.POID != nil {
		conditions = append(conditions, fmt.Sprintf("sp.po_id = $%d", argIndex))
		args = append(args, *params.POID)
		argIndex++
	}

	if params.PaymentStatus != nil {
		conditions = append(conditions, fmt.Sprintf("sp.payment_status = $%d", argIndex))
		args = append(args, *params.PaymentStatus)
		argIndex++
	}

	if params.PaymentMethod != nil {
		conditions = append(conditions, fmt.Sprintf("sp.payment_method = $%d", argIndex))
		args = append(args, *params.PaymentMethod)
		argIndex++
	}

	if params.ProcessedBy != nil {
		conditions = append(conditions, fmt.Sprintf("sp.processed_by = $%d", argIndex))
		args = append(args, *params.ProcessedBy)
		argIndex++
	}

	if params.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("sp.payment_date >= $%d", argIndex))
		args = append(args, *params.DateFrom)
		argIndex++
	}

	if params.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("sp.payment_date <= $%d", argIndex))
		args = append(args, *params.DateTo)
		argIndex++
	}

	if params.IsOverdue != nil {
		if *params.IsOverdue {
			conditions = append(conditions, "sp.payment_status = 'overdue' OR sp.days_overdue > 0")
		} else {
			conditions = append(conditions, "sp.payment_status != 'overdue' AND sp.days_overdue = 0")
		}
	}

	if params.MinAmount != nil {
		conditions = append(conditions, fmt.Sprintf("sp.invoice_amount >= $%d", argIndex))
		args = append(args, *params.MinAmount)
		argIndex++
	}

	if params.MaxAmount != nil {
		conditions = append(conditions, fmt.Sprintf("sp.invoice_amount <= $%d", argIndex))
		args = append(args, *params.MaxAmount)
		argIndex++
	}

	if params.Search != "" {
		searchCondition := fmt.Sprintf(`(
			sp.payment_number ILIKE $%d OR
			sp.invoice_number ILIKE $%d OR
			sp.payment_reference ILIKE $%d OR
			s.company_name ILIKE $%d OR
			po.po_number ILIKE $%d
		)`, argIndex, argIndex, argIndex, argIndex, argIndex)
		conditions = append(conditions, searchCondition)
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	// Combine conditions
	whereClause := baseQuery
	if len(conditions) > 0 {
		whereClause += " AND " + strings.Join(conditions, " AND ")
	}

	// Count total records
	countQuery := "SELECT COUNT(*) " + whereClause
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
		params.Limit = 20
	}

	offset := (params.Page - 1) * params.Limit

	// Get data with pagination
	dataQuery := `
		SELECT sp.payment_id, sp.supplier_id, sp.po_id, sp.payment_number,
			sp.invoice_amount, sp.payment_amount, sp.outstanding_amount,
			sp.invoice_date, sp.payment_date, sp.due_date, sp.payment_method,
			sp.invoice_number, sp.payment_status, sp.days_overdue
		` + whereClause + `
		ORDER BY sp.payment_date DESC, sp.payment_id DESC
		LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)

	args = append(args, params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get supplier payments: %w", err)
	}
	defer rows.Close()

	var items []products.SupplierPaymentListItem
	for rows.Next() {
		var item products.SupplierPaymentListItem
		err := rows.Scan(
			&item.PaymentID,
			&item.SupplierID,
			&item.POID,
			&item.PaymentNumber,
			&item.InvoiceAmount,
			&item.PaymentAmount,
			&item.OutstandingAmount,
			&item.InvoiceDate,
			&item.PaymentDate,
			&item.DueDate,
			&item.PaymentMethod,
			&item.InvoiceNumber,
			&item.PaymentStatus,
			&item.DaysOverdue,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan supplier payment: %w", err)
		}
		items = append(items, item)
	}

	totalPages := (total + int64(params.Limit) - 1) / int64(params.Limit)

	return &common.PaginatedResponse{
		Data:        items,
		Total: int(total),
		Page:        params.Page,
		Limit:       params.Limit,
		TotalPages:  int(totalPages),
		HasMore:     params.Page < int(totalPages),
		
	}, nil
}

// GetBySupplierID retrieves supplier payments by supplier ID with pagination
func (r *SupplierPaymentRepository) GetBySupplierID(ctx context.Context, supplierID int, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error) {
	params.SupplierID = &supplierID
	return r.List(ctx, params)
}

// GetByPOID retrieves supplier payments by purchase order ID with pagination
func (r *SupplierPaymentRepository) GetByPOID(ctx context.Context, poID int, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error) {
	params.POID = &poID
	return r.List(ctx, params)
}

// GetOverduePayments retrieves overdue payments with pagination
func (r *SupplierPaymentRepository) GetOverduePayments(ctx context.Context, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error) {
	isOverdue := true
	params.IsOverdue = &isOverdue
	return r.List(ctx, params)
}

// UpdatePaymentStatus updates the payment status
func (r *SupplierPaymentRepository) UpdatePaymentStatus(ctx context.Context, id int, status products.PaymentStatus) error {
	query := `
		UPDATE supplier_payments SET 
			payment_status = $2, 
			updated_at = NOW()
		WHERE payment_id = $1`

	result, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("supplier payment not found")
	}

	return nil
}

// AddPayment adds a payment amount to an existing payment record
func (r *SupplierPaymentRepository) AddPayment(ctx context.Context, id int, amount float64, method products.PaymentMethod, reference *string) error {
	query := `
		UPDATE supplier_payments SET
			payment_amount = payment_amount + $2,
			outstanding_amount = invoice_amount - (payment_amount + $2) - discount_taken,
			payment_method = $3,
			payment_reference = $4,
			payment_date = NOW(),
			payment_status = CASE 
				WHEN (invoice_amount - (payment_amount + $2) - discount_taken) <= 0 THEN 'paid'
				WHEN payment_amount + $2 > 0 THEN 'partial'
				ELSE payment_status
			END,
			updated_at = NOW()
		WHERE payment_id = $1`

	result, err := r.db.ExecContext(ctx, query, id, amount, method, reference)
	if err != nil {
		return fmt.Errorf("failed to add payment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("supplier payment not found")
	}

	return nil
}

// GenerateNumber generates a unique payment number
func (r *SupplierPaymentRepository) GenerateNumber(ctx context.Context) (string, error) {
	currentYear := time.Now().Year()
	prefix := fmt.Sprintf("SP%d", currentYear)

	query := `
		SELECT COALESCE(MAX(CAST(SUBSTRING(payment_number FROM LENGTH($1) + 1) AS INTEGER)), 0) + 1
		FROM supplier_payments
		WHERE payment_number LIKE $1 || '%'`

	var nextNumber int
	err := r.db.QueryRowContext(ctx, query, prefix).Scan(&nextNumber)
	if err != nil {
		return "", fmt.Errorf("failed to generate payment number: %w", err)
	}

	return fmt.Sprintf("%s%04d", prefix, nextNumber), nil
}

// IsNumberExists checks if a payment number already exists
func (r *SupplierPaymentRepository) IsNumberExists(ctx context.Context, number string) (bool, error) {
	query := `SELECT COUNT(*) FROM supplier_payments WHERE payment_number = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, number).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check payment number existence: %w", err)
	}

	return count > 0, nil
}

// UpdateOverdueStatus updates overdue status for all payments
func (r *SupplierPaymentRepository) UpdateOverdueStatus(ctx context.Context) error {
	query := `
		UPDATE supplier_payments SET
			days_overdue = CASE 
				WHEN payment_status = 'paid' THEN 0
				WHEN due_date < CURRENT_DATE THEN EXTRACT(DAYS FROM (CURRENT_DATE - due_date))::INTEGER
				ELSE 0
			END,
			payment_status = CASE 
				WHEN payment_status = 'paid' THEN 'paid'
				WHEN payment_status = 'disputed' THEN 'disputed'
				WHEN due_date < CURRENT_DATE AND outstanding_amount > 0 THEN 'overdue'
				WHEN payment_amount > 0 AND outstanding_amount > 0 THEN 'partial'
				WHEN outstanding_amount <= 0 THEN 'paid'
				ELSE 'pending'
			END,
			updated_at = NOW()
		WHERE payment_status != 'disputed'`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to update overdue status: %w", err)
	}

	return nil
}

// GetPaymentSummary retrieves payment summary for a supplier or all suppliers
func (r *SupplierPaymentRepository) GetPaymentSummary(ctx context.Context, supplierID *int) (map[string]interface{}, error) {
	var query string
	var args []interface{}

	if supplierID != nil {
		query = `
			SELECT 
				COUNT(*) as total_payments,
				COALESCE(SUM(invoice_amount), 0) as total_invoiced,
				COALESCE(SUM(payment_amount), 0) as total_paid,
				COALESCE(SUM(outstanding_amount), 0) as total_outstanding,
				COALESCE(SUM(penalty_amount), 0) as total_penalties,
				COUNT(CASE WHEN payment_status = 'overdue' THEN 1 END) as overdue_count,
				COUNT(CASE WHEN payment_status = 'paid' THEN 1 END) as paid_count,
				COUNT(CASE WHEN payment_status = 'pending' THEN 1 END) as pending_count,
				COUNT(CASE WHEN payment_status = 'partial' THEN 1 END) as partial_count
			FROM supplier_payments 
			WHERE supplier_id = $1`
		args = append(args, *supplierID)
	} else {
		query = `
			SELECT 
				COUNT(*) as total_payments,
				COALESCE(SUM(invoice_amount), 0) as total_invoiced,
				COALESCE(SUM(payment_amount), 0) as total_paid,
				COALESCE(SUM(outstanding_amount), 0) as total_outstanding,
				COALESCE(SUM(penalty_amount), 0) as total_penalties,
				COUNT(CASE WHEN payment_status = 'overdue' THEN 1 END) as overdue_count,
				COUNT(CASE WHEN payment_status = 'paid' THEN 1 END) as paid_count,
				COUNT(CASE WHEN payment_status = 'pending' THEN 1 END) as pending_count,
				COUNT(CASE WHEN payment_status = 'partial' THEN 1 END) as partial_count
			FROM supplier_payments`
	}

	var totalPayments, overdueCount, paidCount, pendingCount, partialCount int
	var totalInvoiced, totalPaid, totalOutstanding, totalPenalties float64

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&totalPayments,
		&totalInvoiced,
		&totalPaid,
		&totalOutstanding,
		&totalPenalties,
		&overdueCount,
		&paidCount,
		&pendingCount,
		&partialCount,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get payment summary: %w", err)
	}

	summary := map[string]interface{}{
		"total_payments":     totalPayments,
		"total_invoiced":     totalInvoiced,
		"total_paid":         totalPaid,
		"total_outstanding":  totalOutstanding,
		"total_penalties":    totalPenalties,
		"overdue_count":      overdueCount,
		"paid_count":         paidCount,
		"pending_count":      pendingCount,
		"partial_count":      partialCount,
		"payment_rate":       0.0,
	}

	// Calculate payment rate
	if totalInvoiced > 0 {
		summary["payment_rate"] = (totalPaid / totalInvoiced) * 100
	}

	return summary, nil
}