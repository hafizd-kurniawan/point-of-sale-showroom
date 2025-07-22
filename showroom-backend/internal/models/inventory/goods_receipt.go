package inventory

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// ReceiptStatus represents the status of goods receipt
type ReceiptStatus string

const (
	ReceiptStatusPartial         ReceiptStatus = "partial"
	ReceiptStatusComplete        ReceiptStatus = "complete"
	ReceiptStatusWithDiscrepancy ReceiptStatus = "with_discrepancy"
)

// ConditionReceived represents the condition of received goods
type ConditionReceived string

const (
	ConditionGood      ConditionReceived = "good"
	ConditionDamaged   ConditionReceived = "damaged"
	ConditionExpired   ConditionReceived = "expired"
	ConditionWrongItem ConditionReceived = "wrong_item"
)

// GoodsReceipt represents a goods receipt in the system
type GoodsReceipt struct {
	ReceiptID              int           `json:"receipt_id" db:"receipt_id"`
	PoID                   int           `json:"po_id" db:"po_id"`
	ReceiptNumber          string        `json:"receipt_number" db:"receipt_number"`
	ReceiptDate            time.Time     `json:"receipt_date" db:"receipt_date"`
	ReceivedBy             int           `json:"received_by" db:"received_by"`
	SupplierDeliveryNote   *string       `json:"supplier_delivery_note,omitempty" db:"supplier_delivery_note"`
	SupplierInvoiceNumber  *string       `json:"supplier_invoice_number,omitempty" db:"supplier_invoice_number"`
	TotalReceivedValue     float64       `json:"total_received_value" db:"total_received_value"`
	ReceiptStatus          ReceiptStatus `json:"receipt_status" db:"receipt_status"`
	ReceiptNotes           *string       `json:"receipt_notes,omitempty" db:"receipt_notes"`
	DiscrepancyNotes       *string       `json:"discrepancy_notes,omitempty" db:"discrepancy_notes"`
	ReceiptDocumentsJSON   *string       `json:"receipt_documents_json,omitempty" db:"receipt_documents_json"`
	CreatedAt              time.Time     `json:"created_at" db:"created_at"`

	// Related data
	PoNumber         string                   `json:"po_number,omitempty" db:"po_number"`
	SupplierName     string                   `json:"supplier_name,omitempty" db:"supplier_name"`
	ReceivedByName   string                   `json:"received_by_name,omitempty" db:"received_by_name"`
	Details          []GoodsReceiptDetail     `json:"details,omitempty"`
}

// GoodsReceiptListItem represents a simplified goods receipt for list views
type GoodsReceiptListItem struct {
	ReceiptID             int           `json:"receipt_id" db:"receipt_id"`
	ReceiptNumber         string        `json:"receipt_number" db:"receipt_number"`
	PoNumber              string        `json:"po_number" db:"po_number"`
	SupplierName          string        `json:"supplier_name" db:"supplier_name"`
	ReceiptDate           time.Time     `json:"receipt_date" db:"receipt_date"`
	TotalReceivedValue    float64       `json:"total_received_value" db:"total_received_value"`
	ReceiptStatus         ReceiptStatus `json:"receipt_status" db:"receipt_status"`
	ReceivedByName        string        `json:"received_by_name" db:"received_by_name"`
	CreatedAt             time.Time     `json:"created_at" db:"created_at"`
}

// GoodsReceiptCreateRequest represents a request to create a goods receipt
type GoodsReceiptCreateRequest struct {
	PoID                  int                                `json:"po_id" binding:"required"`
	SupplierDeliveryNote  *string                            `json:"supplier_delivery_note,omitempty" binding:"omitempty,max=100"`
	SupplierInvoiceNumber *string                            `json:"supplier_invoice_number,omitempty" binding:"omitempty,max=100"`
	ReceiptNotes          *string                            `json:"receipt_notes,omitempty"`
	DiscrepancyNotes      *string                            `json:"discrepancy_notes,omitempty"`
	ReceiptDocumentsJSON  *string                            `json:"receipt_documents_json,omitempty"`
	Details               []GoodsReceiptDetailCreateRequest  `json:"details" binding:"required,min=1"`
}

// GoodsReceiptUpdateRequest represents a request to update a goods receipt
type GoodsReceiptUpdateRequest struct {
	SupplierDeliveryNote  *string         `json:"supplier_delivery_note,omitempty" binding:"omitempty,max=100"`
	SupplierInvoiceNumber *string         `json:"supplier_invoice_number,omitempty" binding:"omitempty,max=100"`
	ReceiptStatus         *ReceiptStatus  `json:"receipt_status,omitempty" binding:"omitempty,oneof=partial complete with_discrepancy"`
	ReceiptNotes          *string         `json:"receipt_notes,omitempty"`
	DiscrepancyNotes      *string         `json:"discrepancy_notes,omitempty"`
	ReceiptDocumentsJSON  *string         `json:"receipt_documents_json,omitempty"`
}

// GoodsReceiptFilterParams represents filtering parameters for goods receipt queries
type GoodsReceiptFilterParams struct {
	PoID          *int          `json:"po_id,omitempty" form:"po_id"`
	ReceiptStatus *ReceiptStatus `json:"receipt_status,omitempty" form:"receipt_status"`
	ReceivedBy    *int          `json:"received_by,omitempty" form:"received_by"`
	SupplierID    *int          `json:"supplier_id,omitempty" form:"supplier_id"`
	DateFrom      *time.Time    `json:"date_from,omitempty" form:"date_from"`
	DateTo        *time.Time    `json:"date_to,omitempty" form:"date_to"`
	Search        string        `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// GoodsReceiptDetail represents a goods receipt detail line
type GoodsReceiptDetail struct {
	ReceiptDetailID   int               `json:"receipt_detail_id" db:"receipt_detail_id"`
	ReceiptID         int               `json:"receipt_id" db:"receipt_id"`
	PoDetailID        int               `json:"po_detail_id" db:"po_detail_id"`
	ProductID         int               `json:"product_id" db:"product_id"`
	QuantityReceived  int               `json:"quantity_received" db:"quantity_received"`
	QuantityAccepted  int               `json:"quantity_accepted" db:"quantity_accepted"`
	QuantityRejected  int               `json:"quantity_rejected" db:"quantity_rejected"`
	UnitCost          float64           `json:"unit_cost" db:"unit_cost"`
	TotalCost         float64           `json:"total_cost" db:"total_cost"`
	ConditionReceived ConditionReceived `json:"condition_received" db:"condition_received"`
	InspectionNotes   *string           `json:"inspection_notes,omitempty" db:"inspection_notes"`
	RejectionReason   *string           `json:"rejection_reason,omitempty" db:"rejection_reason"`
	ExpiryDate        *time.Time        `json:"expiry_date,omitempty" db:"expiry_date"`
	BatchNumber       *string           `json:"batch_number,omitempty" db:"batch_number"`
	SerialNumbersJSON *string           `json:"serial_numbers_json,omitempty" db:"serial_numbers_json"`

	// Related data
	ProductCode      string  `json:"product_code,omitempty" db:"product_code"`
	ProductName      string  `json:"product_name,omitempty" db:"product_name"`
	UnitMeasure      string  `json:"unit_measure,omitempty" db:"unit_measure"`
	QuantityOrdered  int     `json:"quantity_ordered,omitempty" db:"quantity_ordered"`
}

// GoodsReceiptDetailCreateRequest represents a request to create a goods receipt detail
type GoodsReceiptDetailCreateRequest struct {
	PoDetailID        int                `json:"po_detail_id" binding:"required"`
	ProductID         int                `json:"product_id" binding:"required"`
	QuantityReceived  int                `json:"quantity_received" binding:"required,min=0"`
	QuantityAccepted  int                `json:"quantity_accepted" binding:"required,min=0"`
	QuantityRejected  int                `json:"quantity_rejected" binding:"min=0"`
	UnitCost          float64            `json:"unit_cost" binding:"required,min=0"`
	ConditionReceived ConditionReceived  `json:"condition_received" binding:"required,oneof=good damaged expired wrong_item"`
	InspectionNotes   *string            `json:"inspection_notes,omitempty"`
	RejectionReason   *string            `json:"rejection_reason,omitempty"`
	ExpiryDate        *time.Time         `json:"expiry_date,omitempty"`
	BatchNumber       *string            `json:"batch_number,omitempty" binding:"omitempty,max=100"`
	SerialNumbersJSON *string            `json:"serial_numbers_json,omitempty"`
}

// GoodsReceiptDetailUpdateRequest represents a request to update a goods receipt detail
type GoodsReceiptDetailUpdateRequest struct {
	QuantityReceived  *int               `json:"quantity_received,omitempty" binding:"omitempty,min=0"`
	QuantityAccepted  *int               `json:"quantity_accepted,omitempty" binding:"omitempty,min=0"`
	QuantityRejected  *int               `json:"quantity_rejected,omitempty" binding:"omitempty,min=0"`
	UnitCost          *float64           `json:"unit_cost,omitempty" binding:"omitempty,min=0"`
	ConditionReceived *ConditionReceived `json:"condition_received,omitempty" binding:"omitempty,oneof=good damaged expired wrong_item"`
	InspectionNotes   *string            `json:"inspection_notes,omitempty"`
	RejectionReason   *string            `json:"rejection_reason,omitempty"`
	ExpiryDate        *time.Time         `json:"expiry_date,omitempty"`
	BatchNumber       *string            `json:"batch_number,omitempty" binding:"omitempty,max=100"`
	SerialNumbersJSON *string            `json:"serial_numbers_json,omitempty"`
}

// Methods for GoodsReceipt

// CalculateTotalValue calculates and updates the total received value
func (gr *GoodsReceipt) CalculateTotalValue() {
	gr.TotalReceivedValue = 0
	for _, detail := range gr.Details {
		gr.TotalReceivedValue += detail.TotalCost
	}
}

// UpdateStatus updates the receipt status based on details
func (gr *GoodsReceipt) UpdateStatus() {
	hasDiscrepancy := false
	allComplete := true

	for _, detail := range gr.Details {
		if detail.QuantityRejected > 0 || detail.ConditionReceived != ConditionGood {
			hasDiscrepancy = true
		}
		if detail.QuantityReceived < detail.QuantityOrdered {
			allComplete = false
		}
	}

	if hasDiscrepancy {
		gr.ReceiptStatus = ReceiptStatusWithDiscrepancy
	} else if allComplete {
		gr.ReceiptStatus = ReceiptStatusComplete
	} else {
		gr.ReceiptStatus = ReceiptStatusPartial
	}
}

// HasDiscrepancies checks if the receipt has any discrepancies
func (gr *GoodsReceipt) HasDiscrepancies() bool {
	return gr.ReceiptStatus == ReceiptStatusWithDiscrepancy
}

// IsComplete checks if the receipt is complete
func (gr *GoodsReceipt) IsComplete() bool {
	return gr.ReceiptStatus == ReceiptStatusComplete
}

// Methods for GoodsReceiptDetail

// CalculateTotalCost calculates and updates the total cost for the detail
func (grd *GoodsReceiptDetail) CalculateTotalCost() {
	grd.TotalCost = float64(grd.QuantityAccepted) * grd.UnitCost
}

// ValidateQuantities validates that quantities are consistent
func (grd *GoodsReceiptDetail) ValidateQuantities() bool {
	return grd.QuantityReceived == (grd.QuantityAccepted + grd.QuantityRejected)
}

// HasRejections checks if the detail has any rejected items
func (grd *GoodsReceiptDetail) HasRejections() bool {
	return grd.QuantityRejected > 0
}

// IsGoodCondition checks if the received items are in good condition
func (grd *GoodsReceiptDetail) IsGoodCondition() bool {
	return grd.ConditionReceived == ConditionGood
}

// RequiresRejectionReason checks if rejection reason is required
func (grd *GoodsReceiptDetail) RequiresRejectionReason() bool {
	return grd.QuantityRejected > 0 || grd.ConditionReceived != ConditionGood
}