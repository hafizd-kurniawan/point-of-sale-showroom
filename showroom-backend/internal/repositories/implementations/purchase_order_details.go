package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// PurchaseOrderDetailRepository implements interfaces.PurchaseOrderDetailRepository
type PurchaseOrderDetailRepository struct {
	db *sql.DB
}

// NewPurchaseOrderDetailRepository creates a new purchase order detail repository
func NewPurchaseOrderDetailRepository(db *sql.DB) interfaces.PurchaseOrderDetailRepository {
	return &PurchaseOrderDetailRepository{db: db}
}

// Create creates a new purchase order detail
func (r *PurchaseOrderDetailRepository) Create(ctx context.Context, detail *products.PurchaseOrderDetail) (*products.PurchaseOrderDetail, error) {
	query := `
		INSERT INTO purchase_order_details (
			po_id, product_id, item_description, quantity_ordered, quantity_received,
			quantity_pending, unit_cost, total_cost, expected_date, received_date,
			line_status, item_notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING po_detail_id`

	err := r.db.QueryRowContext(ctx, query,
		detail.POID,
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
	).Scan(&detail.PODetailID)

	if err != nil {
		return nil, fmt.Errorf("failed to create purchase order detail: %w", err)
	}

	return detail, nil
}

// GetByID retrieves a purchase order detail by ID
func (r *PurchaseOrderDetailRepository) GetByID(ctx context.Context, id int) (*products.PurchaseOrderDetail, error) {
	query := `
		SELECT po_detail_id, po_id, product_id, item_description, quantity_ordered,
			quantity_received, quantity_pending, unit_cost, total_cost, expected_date,
			received_date, line_status, item_notes
		FROM purchase_order_details
		WHERE po_detail_id = $1`

	detail := &products.PurchaseOrderDetail{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&detail.PODetailID,
		&detail.POID,
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
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("purchase order detail not found")
		}
		return nil, fmt.Errorf("failed to get purchase order detail: %w", err)
	}

	return detail, nil
}

// Update updates a purchase order detail
func (r *PurchaseOrderDetailRepository) Update(ctx context.Context, id int, detail *products.PurchaseOrderDetail) (*products.PurchaseOrderDetail, error) {
	query := `
		UPDATE purchase_order_details SET
			product_id = $2, item_description = $3, quantity_ordered = $4,
			quantity_received = $5, quantity_pending = $6, unit_cost = $7,
			total_cost = $8, expected_date = $9, received_date = $10,
			line_status = $11, item_notes = $12
		WHERE po_detail_id = $1`

	_, err := r.db.ExecContext(ctx, query,
		id,
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
		return nil, fmt.Errorf("failed to update purchase order detail: %w", err)
	}

	return r.GetByID(ctx, id)
}

// Delete deletes a purchase order detail
func (r *PurchaseOrderDetailRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM purchase_order_details WHERE po_detail_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete purchase order detail: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase order detail not found")
	}

	return nil
}

// GetByPOID retrieves purchase order details by purchase order ID with pagination
func (r *PurchaseOrderDetailRepository) GetByPOID(ctx context.Context, poID int, params *products.PurchaseOrderDetailFilterParams) (*common.PaginatedResponse, error) {
	var args []interface{}
	var conditions []string
	
	// Base query
	baseQuery := `
		FROM purchase_order_details pod
		LEFT JOIN products_spare_parts p ON pod.product_id = p.product_id
		WHERE pod.po_id = $1`
	
	args = append(args, poID)
	argIndex := 2

	// Apply filters
	if params.LineStatus != nil {
		conditions = append(conditions, fmt.Sprintf("pod.line_status = $%d", argIndex))
		args = append(args, *params.LineStatus)
		argIndex++
	}

	if params.Search != "" {
		searchCondition := fmt.Sprintf(`(
			pod.item_description ILIKE $%d OR
			p.product_name ILIKE $%d OR
			p.product_code ILIKE $%d
		)`, argIndex, argIndex, argIndex)
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
		return nil, fmt.Errorf("failed to count purchase order details: %w", err)
	}

	// Calculate pagination
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	offset := (params.Page - 1) * params.Limit

	// Get data with pagination
	dataQuery := `
		SELECT pod.po_detail_id, pod.po_id, pod.product_id, pod.item_description,
			pod.quantity_ordered, pod.quantity_received, pod.quantity_pending,
			pod.unit_cost, pod.total_cost, pod.line_status
		` + whereClause + `
		ORDER BY pod.po_detail_id
		LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)

	args = append(args, params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchase order details: %w", err)
	}
	defer rows.Close()

	var items []products.PurchaseOrderDetailListItem
	for rows.Next() {
		var item products.PurchaseOrderDetailListItem
		err := rows.Scan(
			&item.PODetailID,
			&item.POID,
			&item.ProductID,
			&item.ItemDescription,
			&item.QuantityOrdered,
			&item.QuantityReceived,
			&item.QuantityPending,
			&item.UnitCost,
			&item.TotalCost,
			&item.LineStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan purchase order detail: %w", err)
		}
		items = append(items, item)
	}

	totalPages := (total + int64(params.Limit) - 1) / int64(params.Limit)

	return &common.PaginatedResponse{
		Data:       items,
		Total:      int(total),
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: int(totalPages),
		HasMore:    params.Page < int(totalPages),
	}, nil
}

// GetByProductID retrieves purchase order details by product ID with pagination
func (r *PurchaseOrderDetailRepository) GetByProductID(ctx context.Context, productID int, params *products.PurchaseOrderDetailFilterParams) (*common.PaginatedResponse, error) {
	var args []interface{}
	var conditions []string
	
	// Base query
	baseQuery := `
		FROM purchase_order_details pod
		LEFT JOIN purchase_orders_parts po ON pod.po_id = po.po_id
		WHERE pod.product_id = $1`
	
	args = append(args, productID)
	argIndex := 2

	// Apply filters
	if params.POID != nil {
		conditions = append(conditions, fmt.Sprintf("pod.po_id = $%d", argIndex))
		args = append(args, *params.POID)
		argIndex++
	}

	if params.LineStatus != nil {
		conditions = append(conditions, fmt.Sprintf("pod.line_status = $%d", argIndex))
		args = append(args, *params.LineStatus)
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
		return nil, fmt.Errorf("failed to count purchase order details: %w", err)
	}

	// Calculate pagination
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	offset := (params.Page - 1) * params.Limit

	// Get data with pagination
	dataQuery := `
		SELECT pod.po_detail_id, pod.po_id, pod.product_id, pod.item_description,
			pod.quantity_ordered, pod.quantity_received, pod.quantity_pending,
			pod.unit_cost, pod.total_cost, pod.line_status
		` + whereClause + `
		ORDER BY pod.po_detail_id DESC
		LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)

	args = append(args, params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchase order details: %w", err)
	}
	defer rows.Close()

	var items []products.PurchaseOrderDetailListItem
	for rows.Next() {
		var item products.PurchaseOrderDetailListItem
		err := rows.Scan(
			&item.PODetailID,
			&item.POID,
			&item.ProductID,
			&item.ItemDescription,
			&item.QuantityOrdered,
			&item.QuantityReceived,
			&item.QuantityPending,
			&item.UnitCost,
			&item.TotalCost,
			&item.LineStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan purchase order detail: %w", err)
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

// UpdateQuantityReceived updates the quantity received for a purchase order detail
func (r *PurchaseOrderDetailRepository) UpdateQuantityReceived(ctx context.Context, id int, quantityReceived int) error {
	query := `
		UPDATE purchase_order_details SET
			quantity_received = quantity_received + $2,
			quantity_pending = quantity_ordered - (quantity_received + $2),
			line_status = CASE 
				WHEN (quantity_received + $2) >= quantity_ordered THEN 'received'
				WHEN (quantity_received + $2) > 0 THEN 'partial'
				ELSE 'pending'
			END,
			received_date = CASE 
				WHEN quantity_received = 0 AND $2 > 0 THEN NOW()
				ELSE received_date
			END
		WHERE po_detail_id = $1`

	result, err := r.db.ExecContext(ctx, query, id, quantityReceived)
	if err != nil {
		return fmt.Errorf("failed to update quantity received: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase order detail not found")
	}

	return nil
}

// UpdateLineStatus updates the line status for a purchase order detail
func (r *PurchaseOrderDetailRepository) UpdateLineStatus(ctx context.Context, id int, status products.LineStatus) error {
	query := `UPDATE purchase_order_details SET line_status = $2 WHERE po_detail_id = $1`

	result, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update line status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase order detail not found")
	}

	return nil
}

// GetPendingReceiptItems retrieves pending receipt items for a purchase order
func (r *PurchaseOrderDetailRepository) GetPendingReceiptItems(ctx context.Context, poID int) ([]products.PurchaseOrderDetail, error) {
	query := `
		SELECT po_detail_id, po_id, product_id, item_description, quantity_ordered,
			quantity_received, quantity_pending, unit_cost, total_cost, expected_date,
			received_date, line_status, item_notes
		FROM purchase_order_details
		WHERE po_id = $1 AND line_status IN ('pending', 'partial')
		ORDER BY po_detail_id`

	rows, err := r.db.QueryContext(ctx, query, poID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending receipt items: %w", err)
	}
	defer rows.Close()

	var details []products.PurchaseOrderDetail
	for rows.Next() {
		var detail products.PurchaseOrderDetail
		err := rows.Scan(
			&detail.PODetailID,
			&detail.POID,
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
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan purchase order detail: %w", err)
		}
		details = append(details, detail)
	}

	return details, nil
}

// BulkCreate creates multiple purchase order details
func (r *PurchaseOrderDetailRepository) BulkCreate(ctx context.Context, details []products.PurchaseOrderDetail) error {
	if len(details) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO purchase_order_details (
			po_id, product_id, item_description, quantity_ordered, quantity_received,
			quantity_pending, unit_cost, total_cost, expected_date, received_date,
			line_status, item_notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, detail := range details {
		_, err = stmt.ExecContext(ctx,
			detail.POID,
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
			return fmt.Errorf("failed to insert purchase order detail: %w", err)
		}
	}

	return tx.Commit()
}

// BulkUpdate updates multiple purchase order details
func (r *PurchaseOrderDetailRepository) BulkUpdate(ctx context.Context, details []products.PurchaseOrderDetail) error {
	if len(details) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		UPDATE purchase_order_details SET
			product_id = $2, item_description = $3, quantity_ordered = $4,
			quantity_received = $5, quantity_pending = $6, unit_cost = $7,
			total_cost = $8, expected_date = $9, received_date = $10,
			line_status = $11, item_notes = $12
		WHERE po_detail_id = $1`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, detail := range details {
		_, err = stmt.ExecContext(ctx,
			detail.PODetailID,
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
			return fmt.Errorf("failed to update purchase order detail: %w", err)
		}
	}

	return tx.Commit()
}

// CalculateSubtotal calculates the subtotal for a purchase order
func (r *PurchaseOrderDetailRepository) CalculateSubtotal(ctx context.Context, poID int) (float64, error) {
	query := `SELECT COALESCE(SUM(total_cost), 0) FROM purchase_order_details WHERE po_id = $1`

	var subtotal float64
	err := r.db.QueryRowContext(ctx, query, poID).Scan(&subtotal)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate subtotal: %w", err)
	}

	return subtotal, nil
}