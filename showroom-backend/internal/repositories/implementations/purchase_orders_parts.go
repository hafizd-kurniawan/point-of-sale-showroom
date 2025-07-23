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

// PurchaseOrderPartsRepository implements interfaces.PurchaseOrderPartsRepository
type PurchaseOrderPartsRepository struct {
	db *sql.DB
}

// NewPurchaseOrderPartsRepository creates a new purchase order parts repository
func NewPurchaseOrderPartsRepository(db *sql.DB) interfaces.PurchaseOrderPartsRepository {
	return &PurchaseOrderPartsRepository{db: db}
}

// Create creates a new purchase order
func (r *PurchaseOrderPartsRepository) Create(ctx context.Context, po *products.PurchaseOrderParts) (*products.PurchaseOrderParts, error) {
	query := `
		INSERT INTO purchase_orders_parts (
			po_number, supplier_id, po_date, required_date, expected_delivery_date,
			po_type, subtotal, tax_amount, discount_amount, shipping_cost, total_amount,
			status, payment_terms, payment_due_date, created_by, delivery_address,
			po_notes, terms_and_conditions
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING po_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		po.PONumber,
		po.SupplierID,
		po.PODate,
		po.RequiredDate,
		po.ExpectedDeliveryDate,
		po.POType,
		po.Subtotal,
		po.TaxAmount,
		po.DiscountAmount,
		po.ShippingCost,
		po.TotalAmount,
		po.Status,
		po.PaymentTerms,
		po.PaymentDueDate,
		po.CreatedBy,
		po.DeliveryAddress,
		po.PONotes,
		po.TermsAndConditions,
	).Scan(&po.POID, &po.CreatedAt, &po.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create purchase order: %w", err)
	}

	return po, nil
}

// GetByID retrieves a purchase order by ID
func (r *PurchaseOrderPartsRepository) GetByID(ctx context.Context, id int) (*products.PurchaseOrderParts, error) {
	query := `
		SELECT po_id, po_number, supplier_id, po_date, required_date, expected_delivery_date,
			   po_type, subtotal, tax_amount, discount_amount, shipping_cost, total_amount,
			   status, payment_terms, payment_due_date, created_by, approved_by, approved_at,
			   delivery_address, po_notes, terms_and_conditions, created_at, updated_at
		FROM purchase_orders_parts
		WHERE po_id = $1`

	po := &products.PurchaseOrderParts{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&po.POID,
		&po.PONumber,
		&po.SupplierID,
		&po.PODate,
		&po.RequiredDate,
		&po.ExpectedDeliveryDate,
		&po.POType,
		&po.Subtotal,
		&po.TaxAmount,
		&po.DiscountAmount,
		&po.ShippingCost,
		&po.TotalAmount,
		&po.Status,
		&po.PaymentTerms,
		&po.PaymentDueDate,
		&po.CreatedBy,
		&po.ApprovedBy,
		&po.ApprovedAt,
		&po.DeliveryAddress,
		&po.PONotes,
		&po.TermsAndConditions,
		&po.CreatedAt,
		&po.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("purchase order with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get purchase order: %w", err)
	}

	return po, nil
}

// GetByNumber retrieves a purchase order by number
func (r *PurchaseOrderPartsRepository) GetByNumber(ctx context.Context, number string) (*products.PurchaseOrderParts, error) {
	query := `
		SELECT po_id, po_number, supplier_id, po_date, required_date, expected_delivery_date,
			   po_type, subtotal, tax_amount, discount_amount, shipping_cost, total_amount,
			   status, payment_terms, payment_due_date, created_by, approved_by, approved_at,
			   delivery_address, po_notes, terms_and_conditions, created_at, updated_at
		FROM purchase_orders_parts
		WHERE po_number = $1`

	po := &products.PurchaseOrderParts{}
	err := r.db.QueryRowContext(ctx, query, number).Scan(
		&po.POID,
		&po.PONumber,
		&po.SupplierID,
		&po.PODate,
		&po.RequiredDate,
		&po.ExpectedDeliveryDate,
		&po.POType,
		&po.Subtotal,
		&po.TaxAmount,
		&po.DiscountAmount,
		&po.ShippingCost,
		&po.TotalAmount,
		&po.Status,
		&po.PaymentTerms,
		&po.PaymentDueDate,
		&po.CreatedBy,
		&po.ApprovedBy,
		&po.ApprovedAt,
		&po.DeliveryAddress,
		&po.PONotes,
		&po.TermsAndConditions,
		&po.CreatedAt,
		&po.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("purchase order with number %s not found", number)
		}
		return nil, fmt.Errorf("failed to get purchase order: %w", err)
	}

	return po, nil
}

// Update updates a purchase order
func (r *PurchaseOrderPartsRepository) Update(ctx context.Context, id int, po *products.PurchaseOrderParts) (*products.PurchaseOrderParts, error) {
	query := `
		UPDATE purchase_orders_parts SET
			supplier_id = $2, required_date = $3, expected_delivery_date = $4,
			po_type = $5, tax_amount = $6, discount_amount = $7, shipping_cost = $8,
			payment_terms = $9, delivery_address = $10, po_notes = $11,
			terms_and_conditions = $12, updated_at = NOW()
		WHERE po_id = $1
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		id,
		po.SupplierID,
		po.RequiredDate,
		po.ExpectedDeliveryDate,
		po.POType,
		po.TaxAmount,
		po.DiscountAmount,
		po.ShippingCost,
		po.PaymentTerms,
		po.DeliveryAddress,
		po.PONotes,
		po.TermsAndConditions,
	).Scan(&po.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update purchase order: %w", err)
	}

	return r.GetByID(ctx, id)
}

// UpdateStatus updates the status of a purchase order
func (r *PurchaseOrderPartsRepository) UpdateStatus(ctx context.Context, id int, status products.POStatus) error {
	query := `UPDATE purchase_orders_parts SET status = $2, updated_at = NOW() WHERE po_id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update purchase order status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase order with ID %d not found", id)
	}

	return nil
}

// Delete deletes a purchase order
func (r *PurchaseOrderPartsRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM purchase_orders_parts WHERE po_id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete purchase order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase order with ID %d not found", id)
	}

	return nil
}

// List retrieves a paginated list of purchase orders
func (r *PurchaseOrderPartsRepository) List(ctx context.Context, params *products.PurchaseOrderPartsFilterParams) (*common.PaginatedResponse, error) {
	params.Validate()

	baseQuery := `
		SELECT po_id, po_number, supplier_id, po_date, required_date, expected_delivery_date,
			   po_type, total_amount, status, payment_terms, created_at
		FROM purchase_orders_parts`

	countQuery := `SELECT COUNT(*) FROM purchase_orders_parts`

	whereConditions, args := r.buildWhereConditions(params)
	if len(whereConditions) > 0 {
		whereClause := " WHERE " + strings.Join(whereConditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count purchase orders: %w", err)
	}

	// Add ordering and pagination
	baseQuery += ` ORDER BY po_date DESC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list purchase orders: %w", err)
	}
	defer rows.Close()

	var items []products.PurchaseOrderPartsListItem
	for rows.Next() {
		var item products.PurchaseOrderPartsListItem
		err := rows.Scan(
			&item.POID,
			&item.PONumber,
			&item.SupplierID,
			&item.PODate,
			&item.RequiredDate,
			&item.ExpectedDeliveryDate,
			&item.POType,
			&item.TotalAmount,
			&item.Status,
			&item.PaymentTerms,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan purchase order: %w", err)
		}
		items = append(items, item)
	}

	return &common.PaginatedResponse{
		Data:       items,
Total: int(total),
Page:       params.Page,
Limit:      params.Limit,
TotalPages: params.GetTotalPages(total),
HasMore:    params.GetHasMore(total),
	}, nil
}

// GetBySupplierID retrieves purchase orders by supplier ID
func (r *PurchaseOrderPartsRepository) GetBySupplierID(ctx context.Context, supplierID int, params *products.PurchaseOrderPartsFilterParams) (*common.PaginatedResponse, error) {
	params.SupplierID = &supplierID
	return r.List(ctx, params)
}

// GetByStatus retrieves purchase orders by status
func (r *PurchaseOrderPartsRepository) GetByStatus(ctx context.Context, status products.POStatus, params *products.PurchaseOrderPartsFilterParams) (*common.PaginatedResponse, error) {
	params.Status = &status
	return r.List(ctx, params)
}

// GetPendingApproval retrieves purchase orders pending approval
func (r *PurchaseOrderPartsRepository) GetPendingApproval(ctx context.Context, params *products.PurchaseOrderPartsFilterParams) (*common.PaginatedResponse, error) {
	params.Validate()

	baseQuery := `
		SELECT po_id, po_number, supplier_id, po_date, required_date, expected_delivery_date,
			   po_type, total_amount, status, payment_terms, created_at
		FROM purchase_orders_parts
		WHERE status = 'draft' AND approved_by IS NULL`

	countQuery := `SELECT COUNT(*) FROM purchase_orders_parts WHERE status = 'draft' AND approved_by IS NULL`

	whereConditions, args := r.buildWhereConditions(params)
	// Remove status condition since we already have it
	filteredConditions := []string{}
	filteredArgs := []interface{}{}
	argIndex := 0
	
	for _, condition := range whereConditions {
		if !strings.Contains(condition, "status") {
			filteredConditions = append(filteredConditions, condition)
			filteredArgs = append(filteredArgs, args[argIndex])
		}
		argIndex++
	}
	
	if len(filteredConditions) > 0 {
		whereClause := " AND " + strings.Join(filteredConditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
		args = filteredArgs
	} else {
		args = []interface{}{}
	}

	// Get total count
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count pending approval purchase orders: %w", err)
	}

	// Add ordering and pagination
	baseQuery += ` ORDER BY po_date ASC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending approval purchase orders: %w", err)
	}
	defer rows.Close()

	var items []products.PurchaseOrderPartsListItem
	for rows.Next() {
		var item products.PurchaseOrderPartsListItem
		err := rows.Scan(
			&item.POID,
			&item.PONumber,
			&item.SupplierID,
			&item.PODate,
			&item.RequiredDate,
			&item.ExpectedDeliveryDate,
			&item.POType,
			&item.TotalAmount,
			&item.Status,
			&item.PaymentTerms,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pending approval purchase order: %w", err)
		}
		items = append(items, item)
	}

	return &common.PaginatedResponse{
		Data:       items,
Total: int(total),
Page:       params.Page,
Limit:      params.Limit,
TotalPages: params.GetTotalPages(total),
HasMore:    params.GetHasMore(total),
	}, nil
}

// Approve approves a purchase order
func (r *PurchaseOrderPartsRepository) Approve(ctx context.Context, id int, approvedBy int) error {
	query := `
		UPDATE purchase_orders_parts 
		SET approved_by = $2, approved_at = NOW(), updated_at = NOW() 
		WHERE po_id = $1 AND status = 'draft' AND approved_by IS NULL`
	
	result, err := r.db.ExecContext(ctx, query, id, approvedBy)
	if err != nil {
		return fmt.Errorf("failed to approve purchase order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase order with ID %d not found or already approved", id)
	}

	return nil
}

// Cancel cancels a purchase order
func (r *PurchaseOrderPartsRepository) Cancel(ctx context.Context, id int) error {
	query := `
		UPDATE purchase_orders_parts 
		SET status = 'cancelled', updated_at = NOW() 
		WHERE po_id = $1 AND status IN ('draft', 'sent', 'acknowledged')`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to cancel purchase order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase order with ID %d not found or cannot be cancelled", id)
	}

	return nil
}

// GenerateNumber generates a new PO number
func (r *PurchaseOrderPartsRepository) GenerateNumber(ctx context.Context) (string, error) {
	currentYear := time.Now().Year()
	query := `
		SELECT COALESCE(MAX(CAST(SUBSTRING(po_number FROM LENGTH($1) + 1) AS INTEGER)), 0) + 1
		FROM purchase_orders_parts
		WHERE po_number ~ $2`

	prefix := fmt.Sprintf("PO-%d-", currentYear)
	pattern := fmt.Sprintf("^PO-%d-[0-9]+$", currentYear)

	var nextNumber int
	err := r.db.QueryRowContext(ctx, query, prefix, pattern).Scan(&nextNumber)
	if err != nil {
		return "", fmt.Errorf("failed to generate PO number: %w", err)
	}

	return fmt.Sprintf("PO-%d-%03d", currentYear, nextNumber), nil
}

// IsNumberExists checks if a PO number already exists
func (r *PurchaseOrderPartsRepository) IsNumberExists(ctx context.Context, number string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM purchase_orders_parts WHERE po_number = $1)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, number).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if PO number exists: %w", err)
	}

	return exists, nil
}

// CalculateTotals calculates and updates totals for a purchase order
func (r *PurchaseOrderPartsRepository) CalculateTotals(ctx context.Context, id int) (*products.PurchaseOrderParts, error) {
	// Calculate subtotal from line items
	subtotalQuery := `
		SELECT COALESCE(SUM(total_cost), 0)
		FROM purchase_order_details
		WHERE po_id = $1`

	var subtotal float64
	err := r.db.QueryRowContext(ctx, subtotalQuery, id).Scan(&subtotal)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate subtotal: %w", err)
	}

	// Get current PO to preserve tax, discount, and shipping
	po, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchase order: %w", err)
	}

	// Calculate total
	po.Subtotal = subtotal
	po.CalculateTotals()

	// Update the totals
	updateQuery := `
		UPDATE purchase_orders_parts 
		SET subtotal = $2, total_amount = $3, updated_at = NOW()
		WHERE po_id = $1`

	_, err = r.db.ExecContext(ctx, updateQuery, id, po.Subtotal, po.TotalAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to update totals: %w", err)
	}

	return po, nil
}

// buildWhereConditions builds WHERE conditions for queries
func (r *PurchaseOrderPartsRepository) buildWhereConditions(params *products.PurchaseOrderPartsFilterParams) ([]string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.SupplierID != nil {
		conditions = append(conditions, fmt.Sprintf("supplier_id = $%d", argIndex))
		args = append(args, *params.SupplierID)
		argIndex++
	}

	if params.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *params.Status)
		argIndex++
	}

	if params.POType != nil {
		conditions = append(conditions, fmt.Sprintf("po_type = $%d", argIndex))
		args = append(args, *params.POType)
		argIndex++
	}

	if params.PaymentTerms != nil {
		conditions = append(conditions, fmt.Sprintf("payment_terms = $%d", argIndex))
		args = append(args, *params.PaymentTerms)
		argIndex++
	}

	if params.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("po_date >= $%d", argIndex))
		args = append(args, *params.DateFrom)
		argIndex++
	}

	if params.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("po_date <= $%d", argIndex))
		args = append(args, *params.DateTo)
		argIndex++
	}

	if params.CreatedBy != nil {
		conditions = append(conditions, fmt.Sprintf("created_by = $%d", argIndex))
		args = append(args, *params.CreatedBy)
		argIndex++
	}

	if params.ApprovedBy != nil {
		conditions = append(conditions, fmt.Sprintf("approved_by = $%d", argIndex))
		args = append(args, *params.ApprovedBy)
		argIndex++
	}

	if params.MinAmount != nil {
		conditions = append(conditions, fmt.Sprintf("total_amount >= $%d", argIndex))
		args = append(args, *params.MinAmount)
		argIndex++
	}

	if params.MaxAmount != nil {
		conditions = append(conditions, fmt.Sprintf("total_amount <= $%d", argIndex))
		args = append(args, *params.MaxAmount)
		argIndex++
	}

	if params.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(po_number ILIKE $%d OR po_notes ILIKE $%d)", argIndex, argIndex))
		searchTerm := "%" + params.Search + "%"
		args = append(args, searchTerm)
		argIndex++
	}

	return conditions, args
}