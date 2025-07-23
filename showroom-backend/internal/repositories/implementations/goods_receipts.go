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

// GoodsReceiptRepository implements interfaces.GoodsReceiptRepository
type GoodsReceiptRepository struct {
	db *sql.DB
}

// NewGoodsReceiptRepository creates a new goods receipt repository
func NewGoodsReceiptRepository(db *sql.DB) interfaces.GoodsReceiptRepository {
	return &GoodsReceiptRepository{db: db}
}

// Create creates a new goods receipt
func (r *GoodsReceiptRepository) Create(ctx context.Context, receipt *products.GoodsReceipt) (*products.GoodsReceipt, error) {
	// Generate receipt number if not provided
	if receipt.ReceiptNumber == "" {
		number, err := r.GenerateNumber(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to generate receipt number: %w", err)
		}
		receipt.ReceiptNumber = number
	}

	// Set default status
	if receipt.ReceiptStatus == "" {
		receipt.ReceiptStatus = products.ReceiptStatusPartial
	}

	query := `
		INSERT INTO goods_receipts (
			po_id, receipt_number, receipt_date, received_by, supplier_delivery_note,
			supplier_invoice_number, total_received_value, receipt_status,
			receipt_notes, discrepancy_notes, receipt_documents_json
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING receipt_id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		receipt.POID,
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
	).Scan(&receipt.ReceiptID, &receipt.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create goods receipt: %w", err)
	}

	return receipt, nil
}

// GetByID retrieves a goods receipt by ID
func (r *GoodsReceiptRepository) GetByID(ctx context.Context, id int) (*products.GoodsReceipt, error) {
	query := `
		SELECT receipt_id, po_id, receipt_number, receipt_date, received_by,
			   supplier_delivery_note, supplier_invoice_number, total_received_value,
			   receipt_status, receipt_notes, discrepancy_notes, receipt_documents_json,
			   created_at
		FROM goods_receipts
		WHERE receipt_id = $1`

	receipt := &products.GoodsReceipt{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&receipt.ReceiptID,
		&receipt.POID,
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("goods receipt with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get goods receipt: %w", err)
	}

	return receipt, nil
}

// GetByNumber retrieves a goods receipt by number
func (r *GoodsReceiptRepository) GetByNumber(ctx context.Context, number string) (*products.GoodsReceipt, error) {
	query := `
		SELECT receipt_id, po_id, receipt_number, receipt_date, received_by,
			   supplier_delivery_note, supplier_invoice_number, total_received_value,
			   receipt_status, receipt_notes, discrepancy_notes, receipt_documents_json,
			   created_at
		FROM goods_receipts
		WHERE receipt_number = $1`

	receipt := &products.GoodsReceipt{}
	err := r.db.QueryRowContext(ctx, query, number).Scan(
		&receipt.ReceiptID,
		&receipt.POID,
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("goods receipt with number %s not found", number)
		}
		return nil, fmt.Errorf("failed to get goods receipt: %w", err)
	}

	return receipt, nil
}

// Update updates a goods receipt
func (r *GoodsReceiptRepository) Update(ctx context.Context, id int, receipt *products.GoodsReceipt) (*products.GoodsReceipt, error) {
	query := `
		UPDATE goods_receipts SET
			receipt_date = $2, supplier_delivery_note = $3, supplier_invoice_number = $4,
			receipt_status = $5, receipt_notes = $6, discrepancy_notes = $7,
			receipt_documents_json = $8
		WHERE receipt_id = $1`

	result, err := r.db.ExecContext(ctx, query,
		id,
		receipt.ReceiptDate,
		receipt.SupplierDeliveryNote,
		receipt.SupplierInvoiceNumber,
		receipt.ReceiptStatus,
		receipt.ReceiptNotes,
		receipt.DiscrepancyNotes,
		receipt.ReceiptDocumentsJSON,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update goods receipt: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("goods receipt with ID %d not found", id)
	}

	return r.GetByID(ctx, id)
}

// Delete deletes a goods receipt
func (r *GoodsReceiptRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM goods_receipts WHERE receipt_id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete goods receipt: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("goods receipt with ID %d not found", id)
	}

	return nil
}

// List retrieves a paginated list of goods receipts
func (r *GoodsReceiptRepository) List(ctx context.Context, params *products.GoodsReceiptFilterParams) (*common.PaginatedResponse, error) {
	params.Validate()

	baseQuery := `
		SELECT receipt_id, po_id, receipt_number, receipt_date, received_by,
			   supplier_delivery_note, supplier_invoice_number, total_received_value,
			   receipt_status, created_at
		FROM goods_receipts`

	countQuery := `SELECT COUNT(*) FROM goods_receipts`

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
		return nil, fmt.Errorf("failed to count goods receipts: %w", err)
	}

	// Add ordering and pagination
	baseQuery += ` ORDER BY receipt_date DESC, receipt_id DESC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	args = append(args, params.Limit, params.GetOffset())

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list goods receipts: %w", err)
	}
	defer rows.Close()

	var items []products.GoodsReceiptListItem
	for rows.Next() {
		var item products.GoodsReceiptListItem
		err := rows.Scan(
			&item.ReceiptID,
			&item.POID,
			&item.ReceiptNumber,
			&item.ReceiptDate,
			&item.ReceivedBy,
			&item.SupplierDeliveryNote,
			&item.SupplierInvoiceNumber,
			&item.TotalReceivedValue,
			&item.ReceiptStatus,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan goods receipt: %w", err)
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

// GetByPOID retrieves goods receipts by purchase order ID
func (r *GoodsReceiptRepository) GetByPOID(ctx context.Context, poID int, params *products.GoodsReceiptFilterParams) (*common.PaginatedResponse, error) {
	params.POID = &poID
	return r.List(ctx, params)
}

// GenerateNumber generates a new receipt number
func (r *GoodsReceiptRepository) GenerateNumber(ctx context.Context) (string, error) {
	currentYear := time.Now().Year()
	query := `
		SELECT COALESCE(MAX(CAST(SUBSTRING(receipt_number FROM LENGTH($1) + 1) AS INTEGER)), 0) + 1
		FROM goods_receipts
		WHERE receipt_number ~ $2`

	prefix := fmt.Sprintf("GR-%d-", currentYear)
	pattern := fmt.Sprintf("^GR-%d-[0-9]+$", currentYear)

	var nextNumber int
	err := r.db.QueryRowContext(ctx, query, prefix, pattern).Scan(&nextNumber)
	if err != nil {
		return "", fmt.Errorf("failed to generate receipt number: %w", err)
	}

	return fmt.Sprintf("GR-%d-%03d", currentYear, nextNumber), nil
}

// IsNumberExists checks if a receipt number already exists
func (r *GoodsReceiptRepository) IsNumberExists(ctx context.Context, number string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM goods_receipts WHERE receipt_number = $1)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, number).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if receipt number exists: %w", err)
	}

	return exists, nil
}

// UpdateStatus updates the receipt status
func (r *GoodsReceiptRepository) UpdateStatus(ctx context.Context, id int, status products.ReceiptStatus) error {
	query := `UPDATE goods_receipts SET receipt_status = $2 WHERE receipt_id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update receipt status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("goods receipt with ID %d not found", id)
	}

	return nil
}

// buildWhereConditions builds WHERE conditions for queries
func (r *GoodsReceiptRepository) buildWhereConditions(params *products.GoodsReceiptFilterParams) ([]string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if params.POID != nil {
		conditions = append(conditions, fmt.Sprintf("po_id = $%d", argIndex))
		args = append(args, *params.POID)
		argIndex++
	}

	if params.ReceivedBy != nil {
		conditions = append(conditions, fmt.Sprintf("received_by = $%d", argIndex))
		args = append(args, *params.ReceivedBy)
		argIndex++
	}

	if params.ReceiptStatus != nil {
		conditions = append(conditions, fmt.Sprintf("receipt_status = $%d", argIndex))
		args = append(args, *params.ReceiptStatus)
		argIndex++
	}

	if params.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("receipt_date >= $%d", argIndex))
		args = append(args, *params.DateFrom)
		argIndex++
	}

	if params.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("receipt_date <= $%d", argIndex))
		args = append(args, *params.DateTo)
		argIndex++
	}

	if params.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(receipt_number ILIKE $%d OR supplier_delivery_note ILIKE $%d OR supplier_invoice_number ILIKE $%d OR receipt_notes ILIKE $%d)", argIndex, argIndex, argIndex, argIndex))
		searchTerm := "%" + params.Search + "%"
		args = append(args, searchTerm)
		argIndex++
	}

	return conditions, args
}