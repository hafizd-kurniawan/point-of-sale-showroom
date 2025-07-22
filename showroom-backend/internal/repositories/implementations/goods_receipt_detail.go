package implementations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

type goodsReceiptDetailRepository struct {
	db *sql.DB
}

// NewGoodsReceiptDetailRepository creates a new goods receipt detail repository
func NewGoodsReceiptDetailRepository(db *sql.DB) interfaces.GoodsReceiptDetailRepository {
	return &goodsReceiptDetailRepository{
		db: db,
	}
}

// Create creates a new goods receipt detail
func (r *goodsReceiptDetailRepository) Create(ctx context.Context, detail *inventory.GoodsReceiptDetail) (*inventory.GoodsReceiptDetail, error) {
	query := `
		INSERT INTO goods_receipt_details (
			receipt_id, po_detail_id, product_id, quantity_received,
			quantity_accepted, quantity_rejected, unit_cost, total_cost,
			condition_received, inspection_notes, rejection_reason,
			expiry_date, batch_number, serial_numbers_json
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING receipt_detail_id`

	err := r.db.QueryRowContext(ctx, query,
		detail.ReceiptID,
		detail.PoDetailID,
		detail.ProductID,
		detail.QuantityReceived,
		detail.QuantityAccepted,
		detail.QuantityRejected,
		detail.UnitCost,
		detail.TotalCost,
		detail.ConditionReceived,
		detail.InspectionNotes,
		detail.RejectionReason,
		detail.ExpiryDate,
		detail.BatchNumber,
		detail.SerialNumbersJSON,
	).Scan(&detail.ReceiptDetailID)

	if err != nil {
		return nil, fmt.Errorf("failed to create goods receipt detail: %w", err)
	}

	return detail, nil
}

// GetByID retrieves a goods receipt detail by ID
func (r *goodsReceiptDetailRepository) GetByID(ctx context.Context, id int) (*inventory.GoodsReceiptDetail, error) {
	query := `
		SELECT 
			grd.receipt_detail_id, grd.receipt_id, grd.po_detail_id, grd.product_id,
			grd.quantity_received, grd.quantity_accepted, grd.quantity_rejected,
			grd.unit_cost, grd.total_cost, grd.condition_received,
			grd.inspection_notes, grd.rejection_reason, grd.expiry_date,
			grd.batch_number, grd.serial_numbers_json
		FROM goods_receipt_details grd
		WHERE grd.receipt_detail_id = $1`

	detail := &inventory.GoodsReceiptDetail{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&detail.ReceiptDetailID,
		&detail.ReceiptID,
		&detail.PoDetailID,
		&detail.ProductID,
		&detail.QuantityReceived,
		&detail.QuantityAccepted,
		&detail.QuantityRejected,
		&detail.UnitCost,
		&detail.TotalCost,
		&detail.ConditionReceived,
		&detail.InspectionNotes,
		&detail.RejectionReason,
		&detail.ExpiryDate,
		&detail.BatchNumber,
		&detail.SerialNumbersJSON,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get goods receipt detail: %w", err)
	}

	return detail, nil
}

// GetByReceiptID retrieves all goods receipt details for a specific receipt
func (r *goodsReceiptDetailRepository) GetByReceiptID(ctx context.Context, receiptID int) ([]inventory.GoodsReceiptDetail, error) {
	query := `
		SELECT 
			grd.receipt_detail_id, grd.receipt_id, grd.po_detail_id, grd.product_id,
			grd.quantity_received, grd.quantity_accepted, grd.quantity_rejected,
			grd.unit_cost, grd.total_cost, grd.condition_received,
			grd.inspection_notes, grd.rejection_reason, grd.expiry_date,
			grd.batch_number, grd.serial_numbers_json
		FROM goods_receipt_details grd
		WHERE grd.receipt_id = $1
		ORDER BY grd.receipt_detail_id ASC`

	rows, err := r.db.QueryContext(ctx, query, receiptID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goods receipt details: %w", err)
	}
	defer rows.Close()

	var details []inventory.GoodsReceiptDetail
	for rows.Next() {
		detail := inventory.GoodsReceiptDetail{}

		err := rows.Scan(
			&detail.ReceiptDetailID,
			&detail.ReceiptID,
			&detail.PoDetailID,
			&detail.ProductID,
			&detail.QuantityReceived,
			&detail.QuantityAccepted,
			&detail.QuantityRejected,
			&detail.UnitCost,
			&detail.TotalCost,
			&detail.ConditionReceived,
			&detail.InspectionNotes,
			&detail.RejectionReason,
			&detail.ExpiryDate,
			&detail.BatchNumber,
			&detail.SerialNumbersJSON,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan goods receipt detail: %w", err)
		}

		details = append(details, detail)
	}

	return details, nil
}

// Update updates a goods receipt detail
func (r *goodsReceiptDetailRepository) Update(ctx context.Context, id int, detail *inventory.GoodsReceiptDetail) (*inventory.GoodsReceiptDetail, error) {
	query := `
		UPDATE goods_receipt_details 
		SET quantity_received = $2, quantity_accepted = $3, quantity_rejected = $4,
			unit_cost = $5, total_cost = $6, condition_received = $7,
			inspection_notes = $8, rejection_reason = $9, expiry_date = $10,
			batch_number = $11, serial_numbers_json = $12
		WHERE receipt_detail_id = $1`

	_, err := r.db.ExecContext(ctx, query,
		id,
		detail.QuantityReceived,
		detail.QuantityAccepted,
		detail.QuantityRejected,
		detail.UnitCost,
		detail.TotalCost,
		detail.ConditionReceived,
		detail.InspectionNotes,
		detail.RejectionReason,
		detail.ExpiryDate,
		detail.BatchNumber,
		detail.SerialNumbersJSON,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update goods receipt detail: %w", err)
	}

	detail.ReceiptDetailID = id
	return detail, nil
}

// Delete soft deletes a goods receipt detail
func (r *goodsReceiptDetailRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM goods_receipt_details WHERE receipt_detail_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete goods receipt detail: %w", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("goods receipt detail with ID %d not found", id)
	}

	return nil
}

// CreateBatch creates multiple goods receipt details in a transaction
func (r *goodsReceiptDetailRepository) CreateBatch(ctx context.Context, details []inventory.GoodsReceiptDetail) error {
	if len(details) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO goods_receipt_details (
			receipt_id, po_detail_id, product_id, quantity_received,
			quantity_accepted, quantity_rejected, unit_cost, total_cost,
			condition_received, inspection_notes, rejection_reason,
			expiry_date, batch_number, serial_numbers_json
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	for _, detail := range details {
		_, err := tx.ExecContext(ctx, query,
			detail.ReceiptID,
			detail.PoDetailID,
			detail.ProductID,
			detail.QuantityReceived,
			detail.QuantityAccepted,
			detail.QuantityRejected,
			detail.UnitCost,
			detail.TotalCost,
			detail.ConditionReceived,
			detail.InspectionNotes,
			detail.RejectionReason,
			detail.ExpiryDate,
			detail.BatchNumber,
			detail.SerialNumbersJSON,
		)

		if err != nil {
			return fmt.Errorf("failed to create goods receipt detail: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetRejectedItems retrieves goods receipt details with rejected quantities
func (r *goodsReceiptDetailRepository) GetRejectedItems(ctx context.Context, page, limit int) ([]inventory.GoodsReceiptDetail, int, error) {
	offset := (page - 1) * limit

	query := `
		SELECT 
			grd.receipt_detail_id, grd.receipt_id, grd.po_detail_id, grd.product_id,
			grd.quantity_received, grd.quantity_accepted, grd.quantity_rejected,
			grd.unit_cost, grd.total_cost, grd.condition_received,
			grd.inspection_notes, grd.rejection_reason, grd.expiry_date,
			grd.batch_number, grd.serial_numbers_json
		FROM goods_receipt_details grd
		WHERE grd.quantity_rejected > 0
		ORDER BY grd.receipt_detail_id DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get rejected items: %w", err)
	}
	defer rows.Close()

	var details []inventory.GoodsReceiptDetail
	for rows.Next() {
		detail := inventory.GoodsReceiptDetail{}

		err := rows.Scan(
			&detail.ReceiptDetailID,
			&detail.ReceiptID,
			&detail.PoDetailID,
			&detail.ProductID,
			&detail.QuantityReceived,
			&detail.QuantityAccepted,
			&detail.QuantityRejected,
			&detail.UnitCost,
			&detail.TotalCost,
			&detail.ConditionReceived,
			&detail.InspectionNotes,
			&detail.RejectionReason,
			&detail.ExpiryDate,
			&detail.BatchNumber,
			&detail.SerialNumbersJSON,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan rejected item: %w", err)
		}

		details = append(details, detail)
	}

	// Count total
	countQuery := `SELECT COUNT(*) FROM goods_receipt_details WHERE quantity_rejected > 0`
	var total int
	err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count rejected items: %w", err)
	}

	return details, total, nil
}

// GetDamagedItems retrieves goods receipt details with damaged condition
func (r *goodsReceiptDetailRepository) GetDamagedItems(ctx context.Context, page, limit int) ([]inventory.GoodsReceiptDetail, int, error) {
	offset := (page - 1) * limit

	query := `
		SELECT 
			grd.receipt_detail_id, grd.receipt_id, grd.po_detail_id, grd.product_id,
			grd.quantity_received, grd.quantity_accepted, grd.quantity_rejected,
			grd.unit_cost, grd.total_cost, grd.condition_received,
			grd.inspection_notes, grd.rejection_reason, grd.expiry_date,
			grd.batch_number, grd.serial_numbers_json
		FROM goods_receipt_details grd
		WHERE grd.condition_received = 'damaged'
		ORDER BY grd.receipt_detail_id DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get damaged items: %w", err)
	}
	defer rows.Close()

	var details []inventory.GoodsReceiptDetail
	for rows.Next() {
		detail := inventory.GoodsReceiptDetail{}

		err := rows.Scan(
			&detail.ReceiptDetailID,
			&detail.ReceiptID,
			&detail.PoDetailID,
			&detail.ProductID,
			&detail.QuantityReceived,
			&detail.QuantityAccepted,
			&detail.QuantityRejected,
			&detail.UnitCost,
			&detail.TotalCost,
			&detail.ConditionReceived,
			&detail.InspectionNotes,
			&detail.RejectionReason,
			&detail.ExpiryDate,
			&detail.BatchNumber,
			&detail.SerialNumbersJSON,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan damaged item: %w", err)
		}

		details = append(details, detail)
	}

	// Count total
	countQuery := `SELECT COUNT(*) FROM goods_receipt_details WHERE condition_received = 'damaged'`
	var total int
	err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count damaged items: %w", err)
	}

	return details, total, nil
}