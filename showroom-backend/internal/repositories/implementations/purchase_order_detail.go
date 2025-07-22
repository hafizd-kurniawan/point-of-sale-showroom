package implementations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

type purchaseOrderDetailRepository struct {
	db *sql.DB
}

// NewPurchaseOrderDetailRepository creates a new purchase order detail repository
func NewPurchaseOrderDetailRepository(db *sql.DB) interfaces.PurchaseOrderDetailRepository {
	return &purchaseOrderDetailRepository{
		db: db,
	}
}

// Create creates a new purchase order detail
func (r *purchaseOrderDetailRepository) Create(ctx context.Context, detail *inventory.PurchaseOrderDetail) (*inventory.PurchaseOrderDetail, error) {
	query := `
		INSERT INTO purchase_order_details (
			po_id, product_id, item_description, quantity_ordered, 
			quantity_received, quantity_pending, unit_cost, total_cost,
			expected_date, received_date, line_status, item_notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING po_detail_id`

	err := r.db.QueryRowContext(ctx, query,
		detail.PoID,
		detail.ProductID,
		detail.ItemDescription,
		detail.QuantityOrdered,
		detail.QuantityReceived,
		detail.QuantityPending,
		detail.UnitCost,
		detail.TotalCost,
		detail.ExpectedDate,
		detail.ReceivedDate,
		detail.LineStatus,
		detail.ItemNotes,
	).Scan(&detail.PoDetailID)

	if err != nil {
		return nil, fmt.Errorf("failed to create purchase order detail: %w", err)
	}

	return detail, nil
}

// GetByID retrieves a purchase order detail by ID
func (r *purchaseOrderDetailRepository) GetByID(ctx context.Context, id int) (*inventory.PurchaseOrderDetail, error) {
	query := `
		SELECT 
			pod.po_detail_id, pod.po_id, pod.product_id, pod.item_description,
			pod.quantity_ordered, pod.quantity_received, pod.quantity_pending,
			pod.unit_cost, pod.total_cost, pod.expected_date, pod.received_date,
			pod.line_status, pod.item_notes,
			COALESCE(p.product_code, '') as product_code, 
			COALESCE(p.product_name, '') as product_name
		FROM purchase_order_details pod
		LEFT JOIN products_spare_parts p ON pod.product_id = p.product_id
		LEFT JOIN purchase_orders_parts po ON pod.po_id = po.po_id
		WHERE pod.po_detail_id = $1`

	detail := &inventory.PurchaseOrderDetail{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&detail.PoDetailID,
		&detail.PoID,
		&detail.ProductID,
		&detail.ItemDescription,
		&detail.QuantityOrdered,
		&detail.QuantityReceived,
		&detail.QuantityPending,
		&detail.UnitCost,
		&detail.TotalCost,
		&detail.ExpectedDate,
		&detail.ReceivedDate,
		&detail.LineStatus,
		&detail.ItemNotes,
		&detail.ProductCode,
		&detail.ProductName,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get purchase order detail: %w", err)
	}

	return detail, nil
}

// GetByPOID retrieves all purchase order details for a specific PO
func (r *purchaseOrderDetailRepository) GetByPOID(ctx context.Context, poID int) ([]inventory.PurchaseOrderDetail, error) {
	query := `
		SELECT 
			pod.po_detail_id, pod.po_id, pod.product_id, pod.item_description,
			pod.quantity_ordered, pod.quantity_received, pod.quantity_pending,
			pod.unit_cost, pod.total_cost, pod.expected_date, pod.received_date,
			pod.line_status, pod.item_notes,
			COALESCE(p.product_code, '') as product_code, 
			COALESCE(p.product_name, '') as product_name, 
			COALESCE(p.unit_measure, '') as unit_measure
		FROM purchase_order_details pod
		LEFT JOIN products_spare_parts p ON pod.product_id = p.product_id
		LEFT JOIN purchase_orders_parts po ON pod.po_id = po.po_id
		WHERE pod.po_id = $1
		ORDER BY pod.po_detail_id ASC`

	rows, err := r.db.QueryContext(ctx, query, poID)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchase order details: %w", err)
	}
	defer rows.Close()

	var details []inventory.PurchaseOrderDetail
	for rows.Next() {
		detail := inventory.PurchaseOrderDetail{}

		err := rows.Scan(
			&detail.PoDetailID,
			&detail.PoID,
			&detail.ProductID,
			&detail.ItemDescription,
			&detail.QuantityOrdered,
			&detail.QuantityReceived,
			&detail.QuantityPending,
			&detail.UnitCost,
			&detail.TotalCost,
			&detail.ExpectedDate,
			&detail.ReceivedDate,
			&detail.LineStatus,
			&detail.ItemNotes,
			&detail.ProductCode,
			&detail.ProductName,
			&detail.UnitMeasure,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan purchase order detail: %w", err)
		}

		details = append(details, detail)
	}

	return details, nil
}

// Update updates a purchase order detail
func (r *purchaseOrderDetailRepository) Update(ctx context.Context, id int, detail *inventory.PurchaseOrderDetail) (*inventory.PurchaseOrderDetail, error) {
	query := `
		UPDATE purchase_order_details 
		SET item_description = $2, quantity_ordered = $3, quantity_received = $4,
			quantity_pending = $5, unit_cost = $6, total_cost = $7, 
			expected_date = $8, received_date = $9, line_status = $10, 
			item_notes = $11
		WHERE po_detail_id = $1`

	_, err := r.db.ExecContext(ctx, query,
		id,
		detail.ItemDescription,
		detail.QuantityOrdered,
		detail.QuantityReceived,
		detail.QuantityPending,
		detail.UnitCost,
		detail.TotalCost,
		detail.ExpectedDate,
		detail.ReceivedDate,
		detail.LineStatus,
		detail.ItemNotes,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update purchase order detail: %w", err)
	}

	detail.PoDetailID = id
	return detail, nil
}

// Delete soft deletes a purchase order detail
func (r *purchaseOrderDetailRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM purchase_order_details WHERE po_detail_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete purchase order detail: %w", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("purchase order detail with ID %d not found", id)
	}

	return nil
}

// CreateBatch creates multiple purchase order details in a transaction
func (r *purchaseOrderDetailRepository) CreateBatch(ctx context.Context, details []inventory.PurchaseOrderDetail) error {
	if len(details) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO purchase_order_details (
			po_id, product_id, item_description, quantity_ordered, 
			quantity_received, quantity_pending, unit_cost, total_cost,
			expected_date, received_date, line_status, item_notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	for _, detail := range details {
		_, err := tx.ExecContext(ctx, query,
			detail.PoID,
			detail.ProductID,
			detail.ItemDescription,
			detail.QuantityOrdered,
			detail.QuantityReceived,
			detail.QuantityPending,
			detail.UnitCost,
			detail.TotalCost,
			detail.ExpectedDate,
			detail.ReceivedDate,
			detail.LineStatus,
			detail.ItemNotes,
		)

		if err != nil {
			return fmt.Errorf("failed to create purchase order detail: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateQuantityReceived updates the quantity received for a purchase order detail
func (r *purchaseOrderDetailRepository) UpdateQuantityReceived(ctx context.Context, id int, quantityReceived int) error {
	query := `
		UPDATE purchase_order_details 
		SET quantity_received = $2, 
			quantity_pending = quantity_ordered - $2,
			received_date = CASE 
				WHEN $2 > 0 AND received_date IS NULL THEN CURRENT_TIMESTAMP 
				ELSE received_date 
			END,
			line_status = CASE
				WHEN $2 = 0 THEN 'pending'
				WHEN $2 < quantity_ordered THEN 'partial'
				WHEN $2 = quantity_ordered THEN 'received'
				ELSE line_status
			END
		WHERE po_detail_id = $1`

	result, err := r.db.ExecContext(ctx, query, id, quantityReceived)
	if err != nil {
		return fmt.Errorf("failed to update quantity received: %w", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("purchase order detail with ID %d not found", id)
	}

	return nil
}

// UpdateLineStatus updates the line status for a purchase order detail
func (r *purchaseOrderDetailRepository) UpdateLineStatus(ctx context.Context, id int, status inventory.LineStatus) error {
	query := `
		UPDATE purchase_order_details 
		SET line_status = $2
		WHERE po_detail_id = $1`

	result, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update line status: %w", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("purchase order detail with ID %d not found", id)
	}

	return nil
}

// GetPendingItems retrieves purchase order details with pending status
func (r *purchaseOrderDetailRepository) GetPendingItems(ctx context.Context, page, limit int) ([]inventory.PurchaseOrderDetail, int, error) {
	offset := (page - 1) * limit

	// Count total records
	countQuery := `
		SELECT COUNT(*) 
		FROM purchase_order_details pod
		JOIN purchase_orders_parts po ON pod.po_id = po.po_id
		WHERE pod.line_status = 'pending' 
		AND po.status NOT IN ('cancelled', 'completed')`

	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pending items: %w", err)
	}

	// Get paginated results
	query := `
		SELECT 
			pod.po_detail_id, pod.po_id, pod.product_id, pod.item_description,
			pod.quantity_ordered, pod.quantity_received, pod.quantity_pending,
			pod.unit_cost, pod.total_cost, pod.expected_date, pod.received_date,
			pod.line_status, pod.item_notes,
			COALESCE(p.product_code, '') as product_code, 
			COALESCE(p.product_name, '') as product_name
		FROM purchase_order_details pod
		JOIN purchase_orders_parts po ON pod.po_id = po.po_id
		LEFT JOIN products_spare_parts p ON pod.product_id = p.product_id
		WHERE pod.line_status = 'pending' 
		AND po.status NOT IN ('cancelled', 'completed')
		ORDER BY pod.expected_date ASC, pod.po_detail_id ASC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get pending items: %w", err)
	}
	defer rows.Close()

	var details []inventory.PurchaseOrderDetail
	for rows.Next() {
		detail := inventory.PurchaseOrderDetail{}

		err := rows.Scan(
			&detail.PoDetailID,
			&detail.PoID,
			&detail.ProductID,
			&detail.ItemDescription,
			&detail.QuantityOrdered,
			&detail.QuantityReceived,
			&detail.QuantityPending,
			&detail.UnitCost,
			&detail.TotalCost,
			&detail.ExpectedDate,
			&detail.ReceivedDate,
			&detail.LineStatus,
			&detail.ItemNotes,
			&detail.ProductCode,
			&detail.ProductName,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pending item: %w", err)
		}

		details = append(details, detail)
	}

	return details, total, nil
}

// GetOverdueItems retrieves purchase order details that are overdue
func (r *purchaseOrderDetailRepository) GetOverdueItems(ctx context.Context, page, limit int) ([]inventory.PurchaseOrderDetail, int, error) {
	offset := (page - 1) * limit

	// Count total records
	countQuery := `
		SELECT COUNT(*) 
		FROM purchase_order_details pod
		JOIN purchase_orders_parts po ON pod.po_id = po.po_id
		WHERE pod.expected_date < CURRENT_DATE 
		AND pod.line_status IN ('pending', 'partial')
		AND po.status NOT IN ('cancelled', 'completed')`

	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count overdue items: %w", err)
	}

	// Get paginated results
	query := `
		SELECT 
			pod.po_detail_id, pod.po_id, pod.product_id, pod.item_description,
			pod.quantity_ordered, pod.quantity_received, pod.quantity_pending,
			pod.unit_cost, pod.total_cost, pod.expected_date, pod.received_date,
			pod.line_status, pod.item_notes,
			COALESCE(p.product_code, '') as product_code, 
			COALESCE(p.product_name, '') as product_name,
			(CURRENT_DATE - pod.expected_date) as days_overdue
		FROM purchase_order_details pod
		JOIN purchase_orders_parts po ON pod.po_id = po.po_id
		LEFT JOIN products_spare_parts p ON pod.product_id = p.product_id
		WHERE pod.expected_date < CURRENT_DATE 
		AND pod.line_status IN ('pending', 'partial')
		AND po.status NOT IN ('cancelled', 'completed')
		ORDER BY pod.expected_date ASC, pod.po_detail_id ASC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get overdue items: %w", err)
	}
	defer rows.Close()

	var details []inventory.PurchaseOrderDetail
	for rows.Next() {
		detail := inventory.PurchaseOrderDetail{}
		var daysOverdue int

		err := rows.Scan(
			&detail.PoDetailID,
			&detail.PoID,
			&detail.ProductID,
			&detail.ItemDescription,
			&detail.QuantityOrdered,
			&detail.QuantityReceived,
			&detail.QuantityPending,
			&detail.UnitCost,
			&detail.TotalCost,
			&detail.ExpectedDate,
			&detail.ReceivedDate,
			&detail.LineStatus,
			&detail.ItemNotes,
			&detail.ProductCode,
			&detail.ProductName,
			&daysOverdue,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan overdue item: %w", err)
		}

		details = append(details, detail)
	}

	return details, total, nil
}