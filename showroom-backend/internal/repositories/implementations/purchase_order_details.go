package implementations

import (
	"context"
	"database/sql"
	"fmt"
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
			po_id, product_id, item_description, quantity_ordered, 
			unit_cost, total_cost, expected_date, line_status, item_notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING po_detail_id, quantity_received, quantity_pending, received_date`

	// Calculate total cost and set initial values
	detail.TotalCost = float64(detail.QuantityOrdered) * detail.UnitCost
	detail.QuantityReceived = 0
	detail.QuantityPending = detail.QuantityOrdered
	detail.LineStatus = products.LineStatusPending

	err := r.db.QueryRowContext(ctx, query,
		detail.POID,
		detail.ProductID,
		detail.ItemDescription,
		detail.QuantityOrdered,
		detail.UnitCost,
		detail.TotalCost,
		detail.ExpectedDate,
		detail.LineStatus,
		detail.ItemNotes,
	).Scan(&detail.PODetailID, &detail.QuantityReceived, &detail.QuantityPending, &detail.ReceivedDate)

	if err != nil {
		return nil, fmt.Errorf("failed to create purchase order detail: %w", err)
	}

	return detail, nil
}

// GetByID retrieves a purchase order detail by ID
func (r *PurchaseOrderDetailRepository) GetByID(ctx context.Context, id int) (*products.PurchaseOrderDetail, error) {
	query := `
		SELECT po_detail_id, po_id, product_id, item_description, quantity_ordered,
			   quantity_received, quantity_pending, unit_cost, total_cost,
			   expected_date, received_date, line_status, item_notes
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
		UPDATE purchase_order_details 
		SET product_id = $1, item_description = $2, quantity_ordered = $3,
			unit_cost = $4, total_cost = $5, expected_date = $6, item_notes = $7
		WHERE po_detail_id = $8`

	// Recalculate total cost
	detail.TotalCost = float64(detail.QuantityOrdered) * detail.UnitCost
	
	_, err := r.db.ExecContext(ctx, query,
		detail.ProductID,
		detail.ItemDescription,
		detail.QuantityOrdered,
		detail.UnitCost,
		detail.TotalCost,
		detail.ExpectedDate,
		detail.ItemNotes,
		id,
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
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase order detail not found")
	}

	return nil
}

// GetByPOID retrieves purchase order details by PO ID
func (r *PurchaseOrderDetailRepository) GetByPOID(ctx context.Context, poID int, params *products.PurchaseOrderDetailFilterParams) (*common.PaginatedResponse, error) {
	baseQuery := `
		FROM purchase_order_details pod 
		LEFT JOIN products_spare_parts psp ON pod.product_id = psp.product_id
		WHERE pod.po_id = $1`
	
	args := []interface{}{poID}
	whereConditions := []string{}
	argIndex := 2

	// Add filters
	if params.LineStatus != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("pod.line_status = $%d", argIndex))
		args = append(args, *params.LineStatus)
		argIndex++
	}

	if params.ProductID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("pod.product_id = $%d", argIndex))
		args = append(args, *params.ProductID)
		argIndex++
	}

	if params.Search != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("(psp.product_name ILIKE $%d OR pod.item_description ILIKE $%d)", argIndex, argIndex))
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

	// Main query
	selectFields := `
		pod.po_detail_id, pod.po_id, pod.product_id, pod.item_description,
		pod.quantity_ordered, pod.quantity_received, pod.quantity_pending,
		pod.unit_cost, pod.total_cost, pod.line_status`
	
	mainQuery := "SELECT " + selectFields + " " + baseQuery + 
		" ORDER BY pod.po_detail_id ASC LIMIT $" + fmt.Sprintf("%d", argIndex) + 
		" OFFSET $" + fmt.Sprintf("%d", argIndex+1)
	
	args = append(args, params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, mainQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query purchase order details: %w", err)
	}
	defer rows.Close()

	var details []products.PurchaseOrderDetailListItem
	for rows.Next() {
		var detail products.PurchaseOrderDetailListItem
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
			&detail.LineStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan purchase order detail: %w", err)
		}
		details = append(details, detail)
	}

	totalPages := (total + int64(params.Limit) - 1) / int64(params.Limit)

	return &common.PaginatedResponse{
		Data:       details,
		Total:      int(total),
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: int(totalPages),
		HasMore:    params.Page < int(totalPages),
	}, nil
}

// GetByProductID retrieves purchase order details by product ID
func (r *PurchaseOrderDetailRepository) GetByProductID(ctx context.Context, productID int, params *products.PurchaseOrderDetailFilterParams) (*common.PaginatedResponse, error) {
	baseQuery := `
		FROM purchase_order_details pod 
		LEFT JOIN purchase_orders_parts pop ON pod.po_id = pop.po_id
		WHERE pod.product_id = $1`
	
	args := []interface{}{productID}
	whereConditions := []string{}
	argIndex := 2

	// Add filters
	if params.LineStatus != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("pod.line_status = $%d", argIndex))
		args = append(args, *params.LineStatus)
		argIndex++
	}

	if params.POID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("pod.po_id = $%d", argIndex))
		args = append(args, *params.POID)
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

	// Main query
	selectFields := `
		pod.po_detail_id, pod.po_id, pod.product_id, pod.item_description,
		pod.quantity_ordered, pod.quantity_received, pod.quantity_pending,
		pod.unit_cost, pod.total_cost, pod.line_status`
	
	mainQuery := "SELECT " + selectFields + " " + baseQuery + 
		" ORDER BY pod.po_detail_id DESC LIMIT $" + fmt.Sprintf("%d", argIndex) + 
		" OFFSET $" + fmt.Sprintf("%d", argIndex+1)
	
	args = append(args, params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, mainQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query purchase order details: %w", err)
	}
	defer rows.Close()

	var details []products.PurchaseOrderDetailListItem
	for rows.Next() {
		var detail products.PurchaseOrderDetailListItem
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
			&detail.LineStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan purchase order detail: %w", err)
		}
		details = append(details, detail)
	}

	totalPages := (total + int64(params.Limit) - 1) / int64(params.Limit)

	return &common.PaginatedResponse{
		Data:       details,
		Total:      int(total),
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: int(totalPages),
		HasMore:    params.Page < int(totalPages),
	}, nil
}

// UpdateQuantityReceived updates the quantity received for a line item
func (r *PurchaseOrderDetailRepository) UpdateQuantityReceived(ctx context.Context, id int, quantityReceived int) error {
	query := `
		UPDATE purchase_order_details 
		SET quantity_received = quantity_received + $1,
			quantity_pending = quantity_ordered - (quantity_received + $1),
			line_status = CASE 
				WHEN (quantity_received + $1) >= quantity_ordered THEN 'received'
				WHEN (quantity_received + $1) > 0 THEN 'partial'
				ELSE 'pending'
			END,
			received_date = CASE 
				WHEN quantity_received = 0 AND $1 > 0 THEN NOW()
				ELSE received_date
			END
		WHERE po_detail_id = $2`

	_, err := r.db.ExecContext(ctx, query, quantityReceived, id)
	if err != nil {
		return fmt.Errorf("failed to update quantity received: %w", err)
	}

	return nil
}

// UpdateLineStatus updates the line status
func (r *PurchaseOrderDetailRepository) UpdateLineStatus(ctx context.Context, id int, status products.LineStatus) error {
	query := `UPDATE purchase_order_details SET line_status = $1 WHERE po_detail_id = $2`
	
	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update line status: %w", err)
	}

	return nil
}

// GetPendingReceiptItems gets items pending receipt for a PO
func (r *PurchaseOrderDetailRepository) GetPendingReceiptItems(ctx context.Context, poID int) ([]products.PurchaseOrderDetail, error) {
	query := `
		SELECT po_detail_id, po_id, product_id, item_description, quantity_ordered,
			   quantity_received, quantity_pending, unit_cost, total_cost,
			   expected_date, received_date, line_status, item_notes
		FROM purchase_order_details 
		WHERE po_id = $1 AND line_status IN ('pending', 'partial')
		ORDER BY po_detail_id ASC`

	rows, err := r.db.QueryContext(ctx, query, poID)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending receipt items: %w", err)
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

	query := `
		INSERT INTO purchase_order_details (
			po_id, product_id, item_description, quantity_ordered, 
			unit_cost, total_cost, expected_date, line_status, item_notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	for _, detail := range details {
		// Calculate total cost and set initial values
		detail.TotalCost = float64(detail.QuantityOrdered) * detail.UnitCost
		detail.LineStatus = products.LineStatusPending

		_, err := tx.ExecContext(ctx, query,
			detail.POID,
			detail.ProductID,
			detail.ItemDescription,
			detail.QuantityOrdered,
			detail.UnitCost,
			detail.TotalCost,
			detail.ExpectedDate,
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

	query := `
		UPDATE purchase_order_details 
		SET product_id = $1, item_description = $2, quantity_ordered = $3,
			unit_cost = $4, total_cost = $5, expected_date = $6, item_notes = $7
		WHERE po_detail_id = $8`

	for _, detail := range details {
		// Recalculate total cost
		detail.TotalCost = float64(detail.QuantityOrdered) * detail.UnitCost
		
		_, err := tx.ExecContext(ctx, query,
			detail.ProductID,
			detail.ItemDescription,
			detail.QuantityOrdered,
			detail.UnitCost,
			detail.TotalCost,
			detail.ExpectedDate,
			detail.ItemNotes,
			detail.PODetailID,
		)
		if err != nil {
			return fmt.Errorf("failed to update purchase order detail: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// CalculateSubtotal calculates subtotal for a PO
func (r *PurchaseOrderDetailRepository) CalculateSubtotal(ctx context.Context, poID int) (float64, error) {
	query := `SELECT COALESCE(SUM(total_cost), 0) FROM purchase_order_details WHERE po_id = $1`
	
	var subtotal float64
	err := r.db.QueryRowContext(ctx, query, poID).Scan(&subtotal)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate subtotal: %w", err)
	}

	return subtotal, nil
}