package implementations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// GoodsReceiptDetailRepository implements interfaces.GoodsReceiptDetailRepository
type GoodsReceiptDetailRepository struct {
	db *sql.DB
}

// NewGoodsReceiptDetailRepository creates a new goods receipt detail repository
func NewGoodsReceiptDetailRepository(db *sql.DB) interfaces.GoodsReceiptDetailRepository {
	return &GoodsReceiptDetailRepository{db: db}
}

// Create creates a new goods receipt detail
func (r *GoodsReceiptDetailRepository) Create(ctx context.Context, detail *products.GoodsReceiptDetail) (*products.GoodsReceiptDetail, error) {
	query := `
		INSERT INTO goods_receipt_details (
			receipt_id, po_detail_id, product_id, quantity_received,
			quantity_accepted, quantity_rejected, unit_cost, total_cost,
			condition_received, inspection_notes, rejection_reason,
			expiry_date, batch_number, serial_numbers_json
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING receipt_detail_id`

	// Calculate total cost
	detail.TotalCost = float64(detail.QuantityReceived) * detail.UnitCost

	err := r.db.QueryRowContext(ctx, query,
		detail.ReceiptID,
		detail.PODetailID,
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
func (r *GoodsReceiptDetailRepository) GetByID(ctx context.Context, id int) (*products.GoodsReceiptDetail, error) {
	query := `
		SELECT receipt_detail_id, receipt_id, po_detail_id, product_id,
			   quantity_received, quantity_accepted, quantity_rejected,
			   unit_cost, total_cost, condition_received, inspection_notes,
			   rejection_reason, expiry_date, batch_number, serial_numbers_json
		FROM goods_receipt_details 
		WHERE receipt_detail_id = $1`

	detail := &products.GoodsReceiptDetail{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&detail.ReceiptDetailID,
		&detail.ReceiptID,
		&detail.PODetailID,
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("goods receipt detail not found")
		}
		return nil, fmt.Errorf("failed to get goods receipt detail: %w", err)
	}

	return detail, nil
}

// Update updates a goods receipt detail
func (r *GoodsReceiptDetailRepository) Update(ctx context.Context, id int, detail *products.GoodsReceiptDetail) (*products.GoodsReceiptDetail, error) {
	query := `
		UPDATE goods_receipt_details 
		SET quantity_received = $1, quantity_accepted = $2, quantity_rejected = $3,
			unit_cost = $4, total_cost = $5, condition_received = $6,
			inspection_notes = $7, rejection_reason = $8, expiry_date = $9,
			batch_number = $10, serial_numbers_json = $11
		WHERE receipt_detail_id = $12`

	// Recalculate total cost
	detail.TotalCost = float64(detail.QuantityReceived) * detail.UnitCost
	
	_, err := r.db.ExecContext(ctx, query,
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
		id,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update goods receipt detail: %w", err)
	}

	return r.GetByID(ctx, id)
}

// Delete deletes a goods receipt detail
func (r *GoodsReceiptDetailRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM goods_receipt_details WHERE receipt_detail_id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete goods receipt detail: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("goods receipt detail not found")
	}

	return nil
}

// GetByReceiptID retrieves all details for a goods receipt
func (r *GoodsReceiptDetailRepository) GetByReceiptID(ctx context.Context, receiptID int) ([]products.GoodsReceiptDetail, error) {
	query := `
		SELECT receipt_detail_id, receipt_id, po_detail_id, product_id,
			   quantity_received, quantity_accepted, quantity_rejected,
			   unit_cost, total_cost, condition_received, inspection_notes,
			   rejection_reason, expiry_date, batch_number, serial_numbers_json
		FROM goods_receipt_details 
		WHERE receipt_id = $1
		ORDER BY receipt_detail_id ASC`

	rows, err := r.db.QueryContext(ctx, query, receiptID)
	if err != nil {
		return nil, fmt.Errorf("failed to query goods receipt details: %w", err)
	}
	defer rows.Close()

	var details []products.GoodsReceiptDetail
	for rows.Next() {
		var detail products.GoodsReceiptDetail
		err := rows.Scan(
			&detail.ReceiptDetailID,
			&detail.ReceiptID,
			&detail.PODetailID,
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

// GetByPODetailID retrieves all receipt details for a PO detail
func (r *GoodsReceiptDetailRepository) GetByPODetailID(ctx context.Context, poDetailID int) ([]products.GoodsReceiptDetail, error) {
	query := `
		SELECT receipt_detail_id, receipt_id, po_detail_id, product_id,
			   quantity_received, quantity_accepted, quantity_rejected,
			   unit_cost, total_cost, condition_received, inspection_notes,
			   rejection_reason, expiry_date, batch_number, serial_numbers_json
		FROM goods_receipt_details 
		WHERE po_detail_id = $1
		ORDER BY receipt_detail_id DESC`

	rows, err := r.db.QueryContext(ctx, query, poDetailID)
	if err != nil {
		return nil, fmt.Errorf("failed to query goods receipt details: %w", err)
	}
	defer rows.Close()

	var details []products.GoodsReceiptDetail
	for rows.Next() {
		var detail products.GoodsReceiptDetail
		err := rows.Scan(
			&detail.ReceiptDetailID,
			&detail.ReceiptID,
			&detail.PODetailID,
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

// BulkCreate creates multiple goods receipt details
func (r *GoodsReceiptDetailRepository) BulkCreate(ctx context.Context, details []products.GoodsReceiptDetail) error {
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
		// Calculate total cost
		detail.TotalCost = float64(detail.QuantityReceived) * detail.UnitCost

		_, err := tx.ExecContext(ctx, query,
			detail.ReceiptID,
			detail.PODetailID,
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

// UpdateQuantities updates quantities for accepted and rejected items
func (r *GoodsReceiptDetailRepository) UpdateQuantities(ctx context.Context, id int, quantityAccepted, quantityRejected int) error {
	query := `
		UPDATE goods_receipt_details 
		SET quantity_accepted = $1, quantity_rejected = $2,
			total_cost = quantity_received * unit_cost
		WHERE receipt_detail_id = $3`

	result, err := r.db.ExecContext(ctx, query, quantityAccepted, quantityRejected, id)
	if err != nil {
		return fmt.Errorf("failed to update quantities: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("goods receipt detail not found")
	}

	return nil
}