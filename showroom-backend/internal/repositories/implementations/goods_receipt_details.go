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
	// Validate quantities
	if !detail.ValidateQuantities() {
		return nil, fmt.Errorf("invalid quantities: received (%d) must equal accepted (%d) + rejected (%d)", 
			detail.QuantityReceived, detail.QuantityAccepted, detail.QuantityRejected)
	}

	// Calculate total cost
	detail.UpdateTotalCost()

	query := `
		INSERT INTO goods_receipt_details (
			receipt_id, po_detail_id, product_id, quantity_received, quantity_accepted,
			quantity_rejected, unit_cost, total_cost, condition_received,
			inspection_notes, rejection_reason, expiry_date, batch_number,
			serial_numbers_json
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING receipt_detail_id`

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
		SELECT receipt_detail_id, receipt_id, po_detail_id, product_id, quantity_received,
			   quantity_accepted, quantity_rejected, unit_cost, total_cost,
			   condition_received, inspection_notes, rejection_reason, expiry_date,
			   batch_number, serial_numbers_json
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
			return nil, fmt.Errorf("goods receipt detail with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get goods receipt detail: %w", err)
	}

	return detail, nil
}

// Update updates a goods receipt detail
func (r *GoodsReceiptDetailRepository) Update(ctx context.Context, id int, detail *products.GoodsReceiptDetail) (*products.GoodsReceiptDetail, error) {
	// Validate quantities
	if !detail.ValidateQuantities() {
		return nil, fmt.Errorf("invalid quantities: received (%d) must equal accepted (%d) + rejected (%d)", 
			detail.QuantityReceived, detail.QuantityAccepted, detail.QuantityRejected)
	}

	// Calculate total cost
	detail.UpdateTotalCost()

	query := `
		UPDATE goods_receipt_details SET
			quantity_received = $2, quantity_accepted = $3, quantity_rejected = $4,
			unit_cost = $5, total_cost = $6, condition_received = $7,
			inspection_notes = $8, rejection_reason = $9, expiry_date = $10,
			batch_number = $11, serial_numbers_json = $12
		WHERE receipt_detail_id = $1`

	result, err := r.db.ExecContext(ctx, query,
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

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("goods receipt detail with ID %d not found", id)
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
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("goods receipt detail with ID %d not found", id)
	}

	return nil
}

// GetByReceiptID retrieves goods receipt details by receipt ID
func (r *GoodsReceiptDetailRepository) GetByReceiptID(ctx context.Context, receiptID int) ([]products.GoodsReceiptDetail, error) {
	query := `
		SELECT receipt_detail_id, receipt_id, po_detail_id, product_id, quantity_received,
			   quantity_accepted, quantity_rejected, unit_cost, total_cost,
			   condition_received, inspection_notes, rejection_reason, expiry_date,
			   batch_number, serial_numbers_json
		FROM goods_receipt_details
		WHERE receipt_id = $1
		ORDER BY receipt_detail_id`

	rows, err := r.db.QueryContext(ctx, query, receiptID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goods receipt details by receipt ID: %w", err)
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

// GetByPODetailID retrieves goods receipt details by PO detail ID
func (r *GoodsReceiptDetailRepository) GetByPODetailID(ctx context.Context, poDetailID int) ([]products.GoodsReceiptDetail, error) {
	query := `
		SELECT receipt_detail_id, receipt_id, po_detail_id, product_id, quantity_received,
			   quantity_accepted, quantity_rejected, unit_cost, total_cost,
			   condition_received, inspection_notes, rejection_reason, expiry_date,
			   batch_number, serial_numbers_json
		FROM goods_receipt_details
		WHERE po_detail_id = $1
		ORDER BY receipt_detail_id`

	rows, err := r.db.QueryContext(ctx, query, poDetailID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goods receipt details by PO detail ID: %w", err)
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

// BulkCreate creates multiple goods receipt details in a transaction
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
			receipt_id, po_detail_id, product_id, quantity_received, quantity_accepted,
			quantity_rejected, unit_cost, total_cost, condition_received,
			inspection_notes, rejection_reason, expiry_date, batch_number,
			serial_numbers_json
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	for i := range details {
		// Validate quantities
		if !details[i].ValidateQuantities() {
			return fmt.Errorf("invalid quantities at index %d: received (%d) must equal accepted (%d) + rejected (%d)", 
				i, details[i].QuantityReceived, details[i].QuantityAccepted, details[i].QuantityRejected)
		}

		// Calculate total cost
		details[i].UpdateTotalCost()

		_, err = tx.ExecContext(ctx, query,
			details[i].ReceiptID,
			details[i].PODetailID,
			details[i].ProductID,
			details[i].QuantityReceived,
			details[i].QuantityAccepted,
			details[i].QuantityRejected,
			details[i].UnitCost,
			details[i].TotalCost,
			details[i].ConditionReceived,
			details[i].InspectionNotes,
			details[i].RejectionReason,
			details[i].ExpiryDate,
			details[i].BatchNumber,
			details[i].SerialNumbersJSON,
		)
		if err != nil {
			return fmt.Errorf("failed to create goods receipt detail at index %d: %w", i, err)
		}
	}

	return tx.Commit()
}

// UpdateQuantities updates the accepted and rejected quantities
func (r *GoodsReceiptDetailRepository) UpdateQuantities(ctx context.Context, id int, quantityAccepted, quantityRejected int) error {
	// Get current detail to validate
	detail, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Validate that the sum equals received quantity
	if detail.QuantityReceived != (quantityAccepted + quantityRejected) {
		return fmt.Errorf("invalid quantities: received (%d) must equal accepted (%d) + rejected (%d)", 
			detail.QuantityReceived, quantityAccepted, quantityRejected)
	}

	// Calculate new total cost
	totalCost := float64(quantityAccepted) * detail.UnitCost

	query := `
		UPDATE goods_receipt_details SET
			quantity_accepted = $2, quantity_rejected = $3, total_cost = $4
		WHERE receipt_detail_id = $1`

	result, err := r.db.ExecContext(ctx, query, id, quantityAccepted, quantityRejected, totalCost)
	if err != nil {
		return fmt.Errorf("failed to update quantities: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("goods receipt detail with ID %d not found", id)
	}

	return nil
}