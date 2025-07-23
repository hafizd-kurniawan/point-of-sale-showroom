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
	query := `
		INSERT INTO goods_receipts (
			po_id, receipt_number, receipt_date, received_by, supplier_delivery_note,
			supplier_invoice_number, total_received_value, receipt_status, receipt_notes,
			discrepancy_notes, receipt_documents_json
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
			return nil, fmt.Errorf("goods receipt not found")
		}
		return nil, fmt.Errorf("failed to get goods receipt: %w", err)
	}

	return receipt, nil
}

// GetByNumber retrieves a goods receipt by receipt number
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
			return nil, fmt.Errorf("goods receipt not found")
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
			total_received_value = $5, receipt_status = $6, receipt_notes = $7,
			discrepancy_notes = $8, receipt_documents_json = $9
		WHERE receipt_id = $1`

	_, err := r.db.ExecContext(ctx, query,
		id,
		receipt.ReceiptDate,
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
		return fmt.Errorf("goods receipt not found")
	}

	return nil
}

// List retrieves goods receipts with pagination and filtering
func (r *GoodsReceiptRepository) List(ctx context.Context, params *products.GoodsReceiptFilterParams) (*common.PaginatedResponse, error) {
	var args []interface{}
	var conditions []string
	
	// Base query
	baseQuery := `
		FROM goods_receipts gr
		LEFT JOIN purchase_orders_parts po ON gr.po_id = po.po_id
		LEFT JOIN suppliers s ON po.supplier_id = s.supplier_id
		WHERE 1=1`
	
	argIndex := 1

	// Apply filters
	if params.POID != nil {
		conditions = append(conditions, fmt.Sprintf("gr.po_id = $%d", argIndex))
		args = append(args, *params.POID)
		argIndex++
	}

	if params.ReceivedBy != nil {
		conditions = append(conditions, fmt.Sprintf("gr.received_by = $%d", argIndex))
		args = append(args, *params.ReceivedBy)
		argIndex++
	}

	if params.ReceiptStatus != nil {
		conditions = append(conditions, fmt.Sprintf("gr.receipt_status = $%d", argIndex))
		args = append(args, *params.ReceiptStatus)
		argIndex++
	}

	if params.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("gr.receipt_date >= $%d", argIndex))
		args = append(args, *params.DateFrom)
		argIndex++
	}

	if params.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("gr.receipt_date <= $%d", argIndex))
		args = append(args, *params.DateTo)
		argIndex++
	}

	if params.Search != "" {
		searchCondition := fmt.Sprintf(`(
			gr.receipt_number ILIKE $%d OR
			gr.supplier_delivery_note ILIKE $%d OR
			gr.supplier_invoice_number ILIKE $%d OR
			po.po_number ILIKE $%d OR
			s.company_name ILIKE $%d
		)`, argIndex, argIndex, argIndex, argIndex, argIndex)
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
		return nil, fmt.Errorf("failed to count goods receipts: %w", err)
	}

	// Calculate pagination
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 20
	}

	offset := (params.Page - 1) * params.Limit

	// Get data with pagination
	dataQuery := `
		SELECT gr.receipt_id, gr.po_id, gr.receipt_number, gr.receipt_date,
			gr.received_by, gr.supplier_delivery_note, gr.supplier_invoice_number,
			gr.total_received_value, gr.receipt_status, gr.created_at
		` + whereClause + `
		ORDER BY gr.receipt_date DESC, gr.receipt_id DESC
		LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)

	args = append(args, params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get goods receipts: %w", err)
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

// GetByPOID retrieves goods receipts by purchase order ID with pagination
func (r *GoodsReceiptRepository) GetByPOID(ctx context.Context, poID int, params *products.GoodsReceiptFilterParams) (*common.PaginatedResponse, error) {
	params.POID = &poID
	return r.List(ctx, params)
}

// GenerateNumber generates a unique receipt number
func (r *GoodsReceiptRepository) GenerateNumber(ctx context.Context) (string, error) {
	currentYear := time.Now().Year()
	prefix := fmt.Sprintf("GR%d", currentYear)

	query := `
		SELECT COALESCE(MAX(CAST(SUBSTRING(receipt_number FROM LENGTH($1) + 1) AS INTEGER)), 0) + 1
		FROM goods_receipts
		WHERE receipt_number LIKE $1 || '%'`

	var nextNumber int
	err := r.db.QueryRowContext(ctx, query, prefix).Scan(&nextNumber)
	if err != nil {
		return "", fmt.Errorf("failed to generate receipt number: %w", err)
	}

	return fmt.Sprintf("%s%04d", prefix, nextNumber), nil
}

// IsNumberExists checks if a receipt number already exists
func (r *GoodsReceiptRepository) IsNumberExists(ctx context.Context, number string) (bool, error) {
	query := `SELECT COUNT(*) FROM goods_receipts WHERE receipt_number = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, number).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check receipt number existence: %w", err)
	}

	return count > 0, nil
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
		return fmt.Errorf("goods receipt not found")
	}

	return nil
}