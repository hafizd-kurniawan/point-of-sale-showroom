package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

type purchaseOrderRepository struct {
	db *sql.DB
}

// NewPurchaseOrderRepository creates a new purchase order repository
func NewPurchaseOrderRepository(db *sql.DB) interfaces.PurchaseOrderRepository {
	return &purchaseOrderRepository{db: db}
}

// Create creates a new purchase order
func (r *purchaseOrderRepository) Create(ctx context.Context, po *inventory.PurchaseOrderPart) (*inventory.PurchaseOrderPart, error) {
	query := `
		INSERT INTO purchase_orders_parts (
			po_number, supplier_id, po_date, required_date, expected_delivery_date,
			po_type, subtotal, tax_amount, discount_amount, shipping_cost, total_amount,
			status, payment_terms, payment_due_date, created_by, delivery_address,
			po_notes, terms_and_conditions
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING po_id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		po.PoNumber, po.SupplierID, po.PoDate, po.RequiredDate, po.ExpectedDeliveryDate,
		po.PoType, po.Subtotal, po.TaxAmount, po.DiscountAmount, po.ShippingCost,
		po.TotalAmount, po.Status, po.PaymentTerms, po.PaymentDueDate, po.CreatedBy,
		po.DeliveryAddress, po.PoNotes, po.TermsAndConditions,
	).Scan(&po.PoID, &po.CreatedAt, &po.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create purchase order: %w", err)
	}

	return po, nil
}

// GetByID retrieves a purchase order by ID
func (r *purchaseOrderRepository) GetByID(ctx context.Context, id int) (*inventory.PurchaseOrderPart, error) {
	query := `
		SELECT p.po_id, p.po_number, p.supplier_id, p.po_date, p.required_date,
		       p.expected_delivery_date, p.po_type, p.subtotal, p.tax_amount,
		       p.discount_amount, p.shipping_cost, p.total_amount, p.status,
		       p.payment_terms, p.payment_due_date, p.created_by, p.approved_by,
		       p.approved_at, p.delivery_address, p.po_notes, p.terms_and_conditions,
		       p.created_at, p.updated_at,
		       s.supplier_name, c.full_name as created_by_name, a.full_name as approved_by_name
		FROM purchase_orders_parts p
		JOIN suppliers s ON p.supplier_id = s.supplier_id
		JOIN users c ON p.created_by = c.user_id
		LEFT JOIN users a ON p.approved_by = a.user_id
		WHERE p.po_id = $1`

	po := &inventory.PurchaseOrderPart{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&po.PoID, &po.PoNumber, &po.SupplierID, &po.PoDate, &po.RequiredDate,
		&po.ExpectedDeliveryDate, &po.PoType, &po.Subtotal, &po.TaxAmount,
		&po.DiscountAmount, &po.ShippingCost, &po.TotalAmount, &po.Status,
		&po.PaymentTerms, &po.PaymentDueDate, &po.CreatedBy, &po.ApprovedBy,
		&po.ApprovedAt, &po.DeliveryAddress, &po.PoNotes, &po.TermsAndConditions,
		&po.CreatedAt, &po.UpdatedAt, &po.SupplierName, &po.CreatedByName,
		&po.ApprovedByName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("purchase order not found")
		}
		return nil, fmt.Errorf("failed to get purchase order: %w", err)
	}

	return po, nil
}

// GetByNumber retrieves a purchase order by number
func (r *purchaseOrderRepository) GetByNumber(ctx context.Context, number string) (*inventory.PurchaseOrderPart, error) {
	query := `
		SELECT p.po_id, p.po_number, p.supplier_id, p.po_date, p.required_date,
		       p.expected_delivery_date, p.po_type, p.subtotal, p.tax_amount,
		       p.discount_amount, p.shipping_cost, p.total_amount, p.status,
		       p.payment_terms, p.payment_due_date, p.created_by, p.approved_by,
		       p.approved_at, p.delivery_address, p.po_notes, p.terms_and_conditions,
		       p.created_at, p.updated_at,
		       s.supplier_name, c.full_name as created_by_name, a.full_name as approved_by_name
		FROM purchase_orders_parts p
		JOIN suppliers s ON p.supplier_id = s.supplier_id
		JOIN users c ON p.created_by = c.user_id
		LEFT JOIN users a ON p.approved_by = a.user_id
		WHERE p.po_number = $1`

	po := &inventory.PurchaseOrderPart{}
	err := r.db.QueryRowContext(ctx, query, number).Scan(
		&po.PoID, &po.PoNumber, &po.SupplierID, &po.PoDate, &po.RequiredDate,
		&po.ExpectedDeliveryDate, &po.PoType, &po.Subtotal, &po.TaxAmount,
		&po.DiscountAmount, &po.ShippingCost, &po.TotalAmount, &po.Status,
		&po.PaymentTerms, &po.PaymentDueDate, &po.CreatedBy, &po.ApprovedBy,
		&po.ApprovedAt, &po.DeliveryAddress, &po.PoNotes, &po.TermsAndConditions,
		&po.CreatedAt, &po.UpdatedAt, &po.SupplierName, &po.CreatedByName,
		&po.ApprovedByName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("purchase order not found")
		}
		return nil, fmt.Errorf("failed to get purchase order: %w", err)
	}

	return po, nil
}

// Update updates a purchase order
func (r *purchaseOrderRepository) Update(ctx context.Context, id int, po *inventory.PurchaseOrderPart) (*inventory.PurchaseOrderPart, error) {
	query := `
		UPDATE purchase_orders_parts SET
			required_date = $2, expected_delivery_date = $3, po_type = $4,
			subtotal = $5, tax_amount = $6, discount_amount = $7, shipping_cost = $8,
			total_amount = $9, status = $10, payment_terms = $11, payment_due_date = $12,
			approved_by = $13, approved_at = $14, delivery_address = $15,
			po_notes = $16, terms_and_conditions = $17, updated_at = NOW()
		WHERE po_id = $1
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query, id,
		po.RequiredDate, po.ExpectedDeliveryDate, po.PoType, po.Subtotal,
		po.TaxAmount, po.DiscountAmount, po.ShippingCost, po.TotalAmount,
		po.Status, po.PaymentTerms, po.PaymentDueDate, po.ApprovedBy,
		po.ApprovedAt, po.DeliveryAddress, po.PoNotes, po.TermsAndConditions,
	).Scan(&po.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update purchase order: %w", err)
	}

	po.PoID = id
	return po, nil
}

// Delete deletes a purchase order
func (r *purchaseOrderRepository) Delete(ctx context.Context, id int) error {
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
		return fmt.Errorf("purchase order not found")
	}

	return nil
}

// List retrieves purchase orders with filtering and pagination
func (r *purchaseOrderRepository) List(ctx context.Context, params *inventory.PurchaseOrderPartFilterParams) ([]inventory.PurchaseOrderPartListItem, int, error) {
	whereConditions := []string{"1 = 1"}
	args := []interface{}{}
	argIndex := 1

	// Build WHERE conditions
	if params.SupplierID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.supplier_id = $%d", argIndex))
		args = append(args, *params.SupplierID)
		argIndex++
	}

	if params.Status != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.status = $%d", argIndex))
		args = append(args, *params.Status)
		argIndex++
	}

	if params.PoType != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.po_type = $%d", argIndex))
		args = append(args, *params.PoType)
		argIndex++
	}

	if params.PaymentTerms != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.payment_terms = $%d", argIndex))
		args = append(args, *params.PaymentTerms)
		argIndex++
	}

	if params.CreatedBy != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.created_by = $%d", argIndex))
		args = append(args, *params.CreatedBy)
		argIndex++
	}

	if params.ApprovedBy != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("p.approved_by = $%d", argIndex))
		args = append(args, *params.ApprovedBy)
		argIndex++
	}

	if params.Search != "" {
		whereConditions = append(whereConditions, fmt.Sprintf(
			"(p.po_number ILIKE $%d OR s.supplier_name ILIKE $%d OR p.po_notes ILIKE $%d)",
			argIndex, argIndex, argIndex))
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM purchase_orders_parts p
		JOIN suppliers s ON p.supplier_id = s.supplier_id
		WHERE %s`, whereClause)

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count purchase orders: %w", err)
	}

	// Main query with pagination
	params.PaginationParams.Validate()
	offset := params.PaginationParams.GetOffset()
	limit := params.PaginationParams.Limit

	query := fmt.Sprintf(`
		SELECT p.po_id, p.po_number, s.supplier_name, p.po_date, p.expected_delivery_date,
		       p.po_type, p.total_amount, p.status, p.payment_terms,
		       c.full_name as created_by_name, p.created_at
		FROM purchase_orders_parts p
		JOIN suppliers s ON p.supplier_id = s.supplier_id
		JOIN users c ON p.created_by = c.user_id
		WHERE %s
		ORDER BY p.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list purchase orders: %w", err)
	}
	defer rows.Close()

	var orders []inventory.PurchaseOrderPartListItem
	for rows.Next() {
		var order inventory.PurchaseOrderPartListItem
		err := rows.Scan(
			&order.PoID, &order.PoNumber, &order.SupplierName, &order.PoDate,
			&order.ExpectedDeliveryDate, &order.PoType, &order.TotalAmount,
			&order.Status, &order.PaymentTerms, &order.CreatedByName, &order.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan purchase order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, total, nil
}

// GetBySupplier retrieves purchase orders by supplier
func (r *purchaseOrderRepository) GetBySupplier(ctx context.Context, supplierID int, page, limit int) ([]inventory.PurchaseOrderPartListItem, int, error) {
	params := &inventory.PurchaseOrderPartFilterParams{
		SupplierID: &supplierID,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// GetByStatus retrieves purchase orders by status
func (r *purchaseOrderRepository) GetByStatus(ctx context.Context, status inventory.PurchaseOrderStatus, page, limit int) ([]inventory.PurchaseOrderPartListItem, int, error) {
	params := &inventory.PurchaseOrderPartFilterParams{
		Status: &status,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// GetByDateRange retrieves purchase orders by date range
func (r *purchaseOrderRepository) GetByDateRange(ctx context.Context, startDate, endDate string, page, limit int) ([]inventory.PurchaseOrderPartListItem, int, error) {
	// This is a simplified implementation - in practice, you'd parse the date strings
	params := &inventory.PurchaseOrderPartFilterParams{}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// Search searches purchase orders by query
func (r *purchaseOrderRepository) Search(ctx context.Context, query string, page, limit int) ([]inventory.PurchaseOrderPartListItem, int, error) {
	params := &inventory.PurchaseOrderPartFilterParams{
		Search: query,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// ExistsByNumber checks if a purchase order exists by number
func (r *purchaseOrderRepository) ExistsByNumber(ctx context.Context, number string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM purchase_orders_parts WHERE po_number = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, number).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check purchase order existence: %w", err)
	}
	return exists, nil
}

// ExistsByNumberExcludingID checks if a purchase order exists by number excluding a specific ID
func (r *purchaseOrderRepository) ExistsByNumberExcludingID(ctx context.Context, number string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM purchase_orders_parts WHERE po_number = $1 AND po_id != $2)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, number, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check purchase order existence: %w", err)
	}
	return exists, nil
}

// GetLastPOID gets the last purchase order ID for code generation
func (r *purchaseOrderRepository) GetLastPOID(ctx context.Context) (int, error) {
	query := `SELECT COALESCE(MAX(po_id), 0) FROM purchase_orders_parts`
	var lastID int
	err := r.db.QueryRowContext(ctx, query).Scan(&lastID)
	if err != nil {
		return 0, fmt.Errorf("failed to get last PO ID: %w", err)
	}
	return lastID, nil
}

// UpdateStatus updates purchase order status
func (r *purchaseOrderRepository) UpdateStatus(ctx context.Context, id int, status inventory.PurchaseOrderStatus) error {
	query := `UPDATE purchase_orders_parts SET status = $2, updated_at = NOW() WHERE po_id = $1`
	result, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase order not found")
	}

	return nil
}

// UpdateApproval updates purchase order approval
func (r *purchaseOrderRepository) UpdateApproval(ctx context.Context, id int, approvedBy int) error {
	query := `UPDATE purchase_orders_parts SET approved_by = $2, approved_at = NOW(), updated_at = NOW() WHERE po_id = $1`
	result, err := r.db.ExecContext(ctx, query, id, approvedBy)
	if err != nil {
		return fmt.Errorf("failed to update approval: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase order not found")
	}

	return nil
}

// GetWithDetails retrieves a purchase order with its details
func (r *purchaseOrderRepository) GetWithDetails(ctx context.Context, id int) (*inventory.PurchaseOrderPart, error) {
	// Get the main PO first
	po, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get the details - this would typically be done by a detail repository
	// For now, we'll return the PO without details
	return po, nil
}

// GetPendingApproval retrieves purchase orders pending approval
func (r *purchaseOrderRepository) GetPendingApproval(ctx context.Context, page, limit int) ([]inventory.PurchaseOrderPartListItem, int, error) {
	status := inventory.POStatusDraft
	return r.GetByStatus(ctx, status, page, limit)
}

// GetReadyToSend retrieves purchase orders ready to send
func (r *purchaseOrderRepository) GetReadyToSend(ctx context.Context, page, limit int) ([]inventory.PurchaseOrderPartListItem, int, error) {
	// This would be POs that are approved but not yet sent
	params := &inventory.PurchaseOrderPartFilterParams{}
	params.Page = page
	params.Limit = limit
	// Add logic to filter for approved but not sent orders
	return r.List(ctx, params)
}