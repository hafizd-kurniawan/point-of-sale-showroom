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

	// Calculate derived fields
	detail.QuantityPending = detail.QuantityOrdered - detail.QuantityReceived
	detail.CalculateTotalCost()
	detail.UpdateLineStatus()

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
			return nil, fmt.Errorf("purchase order detail with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get purchase order detail: %w", err)
	}

	return detail, nil
}

// Update updates a purchase order detail
func (r *PurchaseOrderDetailRepository) Update(ctx context.Context, id int, detail *products.PurchaseOrderDetail) (*products.PurchaseOrderDetail, error) {
	// Calculate derived fields
	detail.QuantityPending = detail.QuantityOrdered - detail.QuantityReceived
	detail.CalculateTotalCost()
	detail.UpdateLineStatus()

	query := `
		UPDATE purchase_order_details SET
			product_id = $2, item_description = $3, quantity_ordered = $4,
			quantity_received = $5, quantity_pending = $6, unit_cost = $7,
			total_cost = $8, expected_date = $9, received_date = $10,
			line_status = $11, item_notes = $12
		WHERE po_detail_id = $1`

	result, err := r.db.ExecContext(ctx, query,
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

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("purchase order detail with ID %d not found", id)
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
		return fmt.Errorf("purchase order detail with ID %d not found", id)
	}

	return nil
}

// GetByPOID retrieves purchase order details by PO ID
func (r *PurchaseOrderDetailRepository) GetByPOID(ctx context.Context, poID int, params *products.PurchaseOrderDetailFilterParams) (*common.PaginatedResponse, error) {
	params.POID = &poID
	return r.list(ctx, params)
}

// GetByProductID retrieves purchase order details by product ID
func (r *PurchaseOrderDetailRepository) GetByProductID(ctx context.Context, productID int, params *products.PurchaseOrderDetailFilterParams) (*common.PaginatedResponse, error) {
	params.ProductID = &productID
	return r.list(ctx, params)
}

// UpdateQuantityReceived updates the quantity received for a detail
func (r *PurchaseOrderDetailRepository) UpdateQuantityReceived(ctx context.Context, id int, quantityReceived int) error {
	// Get current detail to calculate new values
	detail, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	detail.ReceiveQuantity(quantityReceived - detail.QuantityReceived)

	query := `
		UPDATE purchase_order_details SET
			quantity_received = $2, quantity_pending = $3, line_status = $4, received_date = $5
		WHERE po_detail_id = $1`

	result, err := r.db.ExecContext(ctx, query,
		id,
		detail.QuantityReceived,
		detail.QuantityPending,
		detail.LineStatus,
		detail.ReceivedDate,
	)

	if err != nil {
		return fmt.Errorf("failed to update quantity received: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("purchase order detail with ID %d not found", id)
	}

	return nil
}

// UpdateLineStatus updates the line status
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
		return fmt.Errorf("purchase order detail with ID %d not found", id)
	}

	return nil
}

// GetPendingReceiptItems retrieves items pending receipt for a PO
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
			return nil, fmt.Errorf("failed to scan pending receipt item: %w", err)
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
			po_id, product_id, item_description, quantity_ordered, quantity_received,
			quantity_pending, unit_cost, total_cost, expected_date, received_date,
			line_status, item_notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	for i := range details {
		// Calculate derived fields
		details[i].QuantityPending = details[i].QuantityOrdered - details[i].QuantityReceived
		details[i].CalculateTotalCost()
		details[i].UpdateLineStatus()

		_, err = tx.ExecContext(ctx, query,
			details[i].POID,
			details[i].ProductID,
			details[i].ItemDescription,
			details[i].QuantityOrdered,
			details[i].QuantityReceived,
			details[i].QuantityPending,
			details[i].UnitCost,
			details[i].TotalCost,
			details[i].ExpectedDate,
			details[i].ReceivedDate,
			details[i].LineStatus,
			details[i].ItemNotes,
		)
		if err != nil {
			return fmt.Errorf("failed to create purchase order detail at index %d: %w", i, err)
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

	query := `
		UPDATE purchase_order_details SET
			product_id = $2, item_description = $3, quantity_ordered = $4,
			quantity_received = $5, quantity_pending = $6, unit_cost = $7,
			total_cost = $8, expected_date = $9, received_date = $10,
			line_status = $11, item_notes = $12
		WHERE po_detail_id = $1`

	for i := range details {
		// Calculate derived fields
		details[i].QuantityPending = details[i].QuantityOrdered - details[i].QuantityReceived
		details[i].CalculateTotalCost()
		details[i].UpdateLineStatus()

		result, err := tx.ExecContext(ctx, query,
			details[i].PODetailID,
			details[i].ProductID,
			details[i].ItemDescription,
			details[i].QuantityOrdered,
			details[i].QuantityReceived,
			details[i].QuantityPending,
			details[i].UnitCost,
			details[i].TotalCost,
			details[i].ExpectedDate,
			details[i].ReceivedDate,
			details[i].LineStatus,
			details[i].ItemNotes,
		)
		if err != nil {
			return fmt.Errorf("failed to update purchase order detail at index %d: %w", i, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected for detail at index %d: %w", i, err)
		}

		if rowsAffected == 0 {
			return fmt.Errorf("purchase order detail with ID %d not found at index %d", details[i].PODetailID, i)
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

// list is a helper method for paginated listing
func (r *PurchaseOrderDetailRepository) list(ctx context.Context, params *products.PurchaseOrderDetailFilterParams) (*common.PaginatedResponse, error) {
	params.Validate()

	baseQuery := `
		SELECT po_detail_id, po_id, product_id, item_description, quantity_ordered,
			   quantity_received, quantity_pending, unit_cost, total_cost, line_status
		FROM purchase_order_details`

	countQuery := `SELECT COUNT(*) FROM purchase_order_details`

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
		return nil, fmt.Errorf("failed to count purchase order details: %w", err)
	}

	// Add ordering and pagination
	baseQuery += ` ORDER BY po_detail_id ASC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list purchase order details: %w", err)
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

	return &common.PaginatedResponse{
		Data:       items,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: params.GetTotalPages(total),
		HasMore:    params.GetHasMore(total),
	}, nil
}

// buildWhereConditions builds WHERE conditions for queries
func (r *PurchaseOrderDetailRepository) buildWhereConditions(params *products.PurchaseOrderDetailFilterParams) ([]string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.POID != nil {
		conditions = append(conditions, fmt.Sprintf("po_id = $%d", argIndex))
		args = append(args, *params.POID)
		argIndex++
	}

	if params.ProductID != nil {
		conditions = append(conditions, fmt.Sprintf("product_id = $%d", argIndex))
		args = append(args, *params.ProductID)
		argIndex++
	}

	if params.LineStatus != nil {
		conditions = append(conditions, fmt.Sprintf("line_status = $%d", argIndex))
		args = append(args, *params.LineStatus)
		argIndex++
	}

	if params.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(item_description ILIKE $%d OR item_notes ILIKE $%d)", argIndex, argIndex))
		searchTerm := "%" + params.Search + "%"
		args = append(args, searchTerm)
		argIndex++
	}

	return conditions, args
}