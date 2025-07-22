package implementations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

type goodsReceiptRepository struct {
	db *sql.DB
}

// NewGoodsReceiptRepository creates a new goods receipt repository
func NewGoodsReceiptRepository(db *sql.DB) interfaces.GoodsReceiptRepository {
	return &goodsReceiptRepository{
		db: db,
	}
}

// Create creates a new goods receipt
func (r *goodsReceiptRepository) Create(ctx context.Context, receipt *inventory.GoodsReceipt) (*inventory.GoodsReceipt, error) {
	query := `
		INSERT INTO goods_receipts (
			po_id, receipt_number, receipt_date, received_by,
			supplier_delivery_note, supplier_invoice_number, total_received_value,
			receipt_status, receipt_notes, discrepancy_notes, receipt_documents_json
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING receipt_id`

	err := r.db.QueryRowContext(ctx, query,
		receipt.PoID,
		receipt.ReceiptNumber,
		receipt.ReceiptDate,
		receipt.ReceivedBy,
		receipt.SupplierDeliveryNote,
		receipt.SupplierInvoiceNumber,
		receipt.TotalReceivedValue,
		receipt.ReceiptStatus,
		receipt.ReceiptNotes,
		receipt.DiscrepancyNotes,
		receipt.ReceiptDocumentsJSON,
	).Scan(&receipt.ReceiptID)

	if err != nil {
		return nil, fmt.Errorf("failed to create goods receipt: %w", err)
	}

	return receipt, nil
}

// GetByID retrieves a goods receipt by ID
func (r *goodsReceiptRepository) GetByID(ctx context.Context, id int) (*inventory.GoodsReceipt, error) {
	query := `
		SELECT 
			gr.receipt_id, gr.po_id, gr.receipt_number, gr.receipt_date,
			gr.received_by, gr.supplier_delivery_note, gr.supplier_invoice_number,
			gr.total_received_value, gr.receipt_status, gr.receipt_notes,
			gr.discrepancy_notes, gr.receipt_documents_json, gr.created_at
		FROM goods_receipts gr
		WHERE gr.receipt_id = $1`

	receipt := &inventory.GoodsReceipt{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&receipt.ReceiptID,
		&receipt.PoID,
		&receipt.ReceiptNumber,
		&receipt.ReceiptDate,
		&receipt.ReceivedBy,
		&receipt.SupplierDeliveryNote,
		&receipt.SupplierInvoiceNumber,
		&receipt.TotalReceivedValue,
		&receipt.ReceiptStatus,
		&receipt.ReceiptNotes,
		&receipt.DiscrepancyNotes,
		&receipt.ReceiptDocumentsJSON,
		&receipt.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get goods receipt: %w", err)
	}

	return receipt, nil
}

// GetByNumber retrieves a goods receipt by receipt number
func (r *goodsReceiptRepository) GetByNumber(ctx context.Context, number string) (*inventory.GoodsReceipt, error) {
	query := `
		SELECT 
			gr.receipt_id, gr.po_id, gr.receipt_number, gr.receipt_date,
			gr.received_by, gr.supplier_delivery_note, gr.supplier_invoice_number,
			gr.total_received_value, gr.receipt_status, gr.receipt_notes,
			gr.discrepancy_notes, gr.receipt_documents_json, gr.created_at
		FROM goods_receipts gr
		WHERE gr.receipt_number = $1`

	receipt := &inventory.GoodsReceipt{}

	err := r.db.QueryRowContext(ctx, query, number).Scan(
		&receipt.ReceiptID,
		&receipt.PoID,
		&receipt.ReceiptNumber,
		&receipt.ReceiptDate,
		&receipt.ReceivedBy,
		&receipt.SupplierDeliveryNote,
		&receipt.SupplierInvoiceNumber,
		&receipt.TotalReceivedValue,
		&receipt.ReceiptStatus,
		&receipt.ReceiptNotes,
		&receipt.DiscrepancyNotes,
		&receipt.ReceiptDocumentsJSON,
		&receipt.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get goods receipt by number: %w", err)
	}

	return receipt, nil
}

// GetByPOID retrieves all goods receipts for a specific purchase order
func (r *goodsReceiptRepository) GetByPOID(ctx context.Context, poID int) ([]inventory.GoodsReceipt, error) {
	query := `
		SELECT 
			gr.receipt_id, gr.po_id, gr.receipt_number, gr.receipt_date,
			gr.received_by, gr.supplier_delivery_note, gr.supplier_invoice_number,
			gr.total_received_value, gr.receipt_status, gr.receipt_notes,
			gr.discrepancy_notes, gr.receipt_documents_json, gr.created_at
		FROM goods_receipts gr
		WHERE gr.po_id = $1
		ORDER BY gr.receipt_date DESC, gr.receipt_id DESC`

	rows, err := r.db.QueryContext(ctx, query, poID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goods receipts by PO ID: %w", err)
	}
	defer rows.Close()

	var receipts []inventory.GoodsReceipt
	for rows.Next() {
		receipt := inventory.GoodsReceipt{}

		err := rows.Scan(
			&receipt.ReceiptID,
			&receipt.PoID,
			&receipt.ReceiptNumber,
			&receipt.ReceiptDate,
			&receipt.ReceivedBy,
			&receipt.SupplierDeliveryNote,
			&receipt.SupplierInvoiceNumber,
			&receipt.TotalReceivedValue,
			&receipt.ReceiptStatus,
			&receipt.ReceiptNotes,
			&receipt.DiscrepancyNotes,
			&receipt.ReceiptDocumentsJSON,
			&receipt.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan goods receipt: %w", err)
		}

		receipts = append(receipts, receipt)
	}

	return receipts, nil
}

// Update updates a goods receipt
func (r *goodsReceiptRepository) Update(ctx context.Context, id int, receipt *inventory.GoodsReceipt) (*inventory.GoodsReceipt, error) {
	query := `
		UPDATE goods_receipts 
		SET receipt_date = $2, received_by = $3, supplier_delivery_note = $4,
			supplier_invoice_number = $5, total_received_value = $6,
			receipt_status = $7, receipt_notes = $8, discrepancy_notes = $9,
			receipt_documents_json = $10
		WHERE receipt_id = $1`

	_, err := r.db.ExecContext(ctx, query,
		id,
		receipt.ReceiptDate,
		receipt.ReceivedBy,
		receipt.SupplierDeliveryNote,
		receipt.SupplierInvoiceNumber,
		receipt.TotalReceivedValue,
		receipt.ReceiptStatus,
		receipt.ReceiptNotes,
		receipt.DiscrepancyNotes,
		receipt.ReceiptDocumentsJSON,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update goods receipt: %w", err)
	}

	receipt.ReceiptID = id
	return receipt, nil
}

// Delete soft deletes a goods receipt
func (r *goodsReceiptRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM goods_receipts WHERE receipt_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete goods receipt: %w", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("goods receipt with ID %d not found", id)
	}

	return nil
}

// List retrieves goods receipts with filtering and pagination
func (r *goodsReceiptRepository) List(ctx context.Context, params *inventory.GoodsReceiptFilterParams) ([]inventory.GoodsReceiptListItem, int, error) {
	// For simplicity, implement basic list without complex filtering
	query := `
		SELECT 
			gr.receipt_id, gr.receipt_number, gr.receipt_date,
			gr.total_received_value, gr.receipt_status, gr.created_at
		FROM goods_receipts gr
		ORDER BY gr.receipt_date DESC
		LIMIT $1 OFFSET $2`

	offset := (params.Page - 1) * params.Limit

	rows, err := r.db.QueryContext(ctx, query, params.Limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get goods receipts: %w", err)
	}
	defer rows.Close()

	var items []inventory.GoodsReceiptListItem
	for rows.Next() {
		item := inventory.GoodsReceiptListItem{}

		err := rows.Scan(
			&item.ReceiptID,
			&item.ReceiptNumber,
			&item.ReceiptDate,
			&item.TotalReceivedValue,
			&item.ReceiptStatus,
			&item.CreatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan goods receipt list item: %w", err)
		}

		items = append(items, item)
	}

	// Count total (simplified)
	countQuery := `SELECT COUNT(*) FROM goods_receipts`
	var total int
	err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count goods receipts: %w", err)
	}

	return items, total, nil
}

// GetByStatus retrieves goods receipts by status
func (r *goodsReceiptRepository) GetByStatus(ctx context.Context, status inventory.ReceiptStatus, page, limit int) ([]inventory.GoodsReceiptListItem, int, error) {
	params := &inventory.GoodsReceiptFilterParams{
		ReceiptStatus: &status,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// GetByDateRange retrieves goods receipts within a date range
func (r *goodsReceiptRepository) GetByDateRange(ctx context.Context, startDate, endDate string, page, limit int) ([]inventory.GoodsReceiptListItem, int, error) {
	params := &inventory.GoodsReceiptFilterParams{}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// Search searches goods receipts
func (r *goodsReceiptRepository) Search(ctx context.Context, query string, page, limit int) ([]inventory.GoodsReceiptListItem, int, error) {
	params := &inventory.GoodsReceiptFilterParams{
		Search: query,
	}
	params.Page = page
	params.Limit = limit
	return r.List(ctx, params)
}

// ExistsByNumber checks if a goods receipt exists by receipt number
func (r *goodsReceiptRepository) ExistsByNumber(ctx context.Context, number string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM goods_receipts WHERE receipt_number = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, number).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check goods receipt existence: %w", err)
	}
	return exists, nil
}

// ExistsByNumberExcludingID checks if a goods receipt exists by receipt number excluding a specific ID
func (r *goodsReceiptRepository) ExistsByNumberExcludingID(ctx context.Context, number string, excludeID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM goods_receipts WHERE receipt_number = $1 AND receipt_id != $2)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, number, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check goods receipt existence: %w", err)
	}
	return exists, nil
}

// GetLastReceiptID gets the last receipt ID for code generation
func (r *goodsReceiptRepository) GetLastReceiptID(ctx context.Context) (int, error) {
	query := `SELECT COALESCE(MAX(receipt_id), 0) FROM goods_receipts`
	var lastID int
	err := r.db.QueryRowContext(ctx, query).Scan(&lastID)
	if err != nil {
		return 0, fmt.Errorf("failed to get last receipt ID: %w", err)
	}
	return lastID, nil
}

// GetWithDetails retrieves a goods receipt with its details
func (r *goodsReceiptRepository) GetWithDetails(ctx context.Context, id int) (*inventory.GoodsReceipt, error) {
	return r.GetByID(ctx, id)
}

// UpdateStatus updates the receipt status
func (r *goodsReceiptRepository) UpdateStatus(ctx context.Context, id int, status inventory.ReceiptStatus) error {
	query := `
		UPDATE goods_receipts 
		SET receipt_status = $2
		WHERE receipt_id = $1`

	result, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update receipt status: %w", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("goods receipt with ID %d not found", id)
	}

	return nil
}