package products

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// ReceiptStatus represents the status of a goods receipt
type ReceiptStatus string

const (
	ReceiptStatusPartial        ReceiptStatus = "partial"
	ReceiptStatusComplete       ReceiptStatus = "complete"
	ReceiptStatusWithDiscrepancy ReceiptStatus = "with_discrepancy"
)

// IsValid checks if the receipt status is valid
func (s ReceiptStatus) IsValid() bool {
	switch s {
	case ReceiptStatusPartial, ReceiptStatusComplete, ReceiptStatusWithDiscrepancy:
		return true
	default:
		return false
	}
}

// String returns the string representation of the receipt status
func (s ReceiptStatus) String() string {
	return string(s)
}

// Value implements the driver.Valuer interface for ReceiptStatus
func (s ReceiptStatus) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan implements the sql.Scanner interface for ReceiptStatus
func (s *ReceiptStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch str := value.(type) {
	case string:
		*s = ReceiptStatus(str)
	case []byte:
		*s = ReceiptStatus(str)
	default:
		return fmt.Errorf("cannot scan %T into ReceiptStatus", value)
	}
	return nil
}

// GoodsReceipt represents a goods receipt from a supplier
type GoodsReceipt struct {
	ReceiptID             int           `json:"receipt_id" db:"receipt_id"`
	POID                  int           `json:"po_id" db:"po_id"`
	ReceiptNumber         string        `json:"receipt_number" db:"receipt_number"`
	ReceiptDate           time.Time     `json:"receipt_date" db:"receipt_date"`
	ReceivedBy            int           `json:"received_by" db:"received_by"`
	SupplierDeliveryNote  *string       `json:"supplier_delivery_note,omitempty" db:"supplier_delivery_note"`
	SupplierInvoiceNumber *string       `json:"supplier_invoice_number,omitempty" db:"supplier_invoice_number"`
	TotalReceivedValue    float64       `json:"total_received_value" db:"total_received_value"`
	ReceiptStatus         ReceiptStatus `json:"receipt_status" db:"receipt_status"`
	ReceiptNotes          *string       `json:"receipt_notes,omitempty" db:"receipt_notes"`
	DiscrepancyNotes      *string       `json:"discrepancy_notes,omitempty" db:"discrepancy_notes"`
	ReceiptDocumentsJSON  *string       `json:"receipt_documents_json,omitempty" db:"receipt_documents_json"`
	CreatedAt             time.Time     `json:"created_at" db:"created_at"`
}

// GoodsReceiptListItem represents a simplified goods receipt for list views
type GoodsReceiptListItem struct {
	ReceiptID             int           `json:"receipt_id" db:"receipt_id"`
	POID                  int           `json:"po_id" db:"po_id"`
	ReceiptNumber         string        `json:"receipt_number" db:"receipt_number"`
	ReceiptDate           time.Time     `json:"receipt_date" db:"receipt_date"`
	ReceivedBy            int           `json:"received_by" db:"received_by"`
	SupplierDeliveryNote  *string       `json:"supplier_delivery_note,omitempty" db:"supplier_delivery_note"`
	SupplierInvoiceNumber *string       `json:"supplier_invoice_number,omitempty" db:"supplier_invoice_number"`
	TotalReceivedValue    float64       `json:"total_received_value" db:"total_received_value"`
	ReceiptStatus         ReceiptStatus `json:"receipt_status" db:"receipt_status"`
	CreatedAt             time.Time     `json:"created_at" db:"created_at"`
}

// GoodsReceiptCreateRequest represents a request to create a goods receipt
type GoodsReceiptCreateRequest struct {
	POID                  int       `json:"po_id" binding:"required,min=1"`
	ReceiptDate           time.Time `json:"receipt_date" binding:"required"`
	SupplierDeliveryNote  *string   `json:"supplier_delivery_note,omitempty" binding:"omitempty,max=100"`
	SupplierInvoiceNumber *string   `json:"supplier_invoice_number,omitempty" binding:"omitempty,max=100"`
	ReceiptNotes          *string   `json:"receipt_notes,omitempty"`
	ReceiptDocumentsJSON  *string   `json:"receipt_documents_json,omitempty"`
}

// GoodsReceiptUpdateRequest represents a request to update a goods receipt
type GoodsReceiptUpdateRequest struct {
	ReceiptDate           *time.Time     `json:"receipt_date,omitempty"`
	SupplierDeliveryNote  *string        `json:"supplier_delivery_note,omitempty" binding:"omitempty,max=100"`
	SupplierInvoiceNumber *string        `json:"supplier_invoice_number,omitempty" binding:"omitempty,max=100"`
	ReceiptStatus         *ReceiptStatus `json:"receipt_status,omitempty"`
	ReceiptNotes          *string        `json:"receipt_notes,omitempty"`
	DiscrepancyNotes      *string        `json:"discrepancy_notes,omitempty"`
	ReceiptDocumentsJSON  *string        `json:"receipt_documents_json,omitempty"`
}

// GoodsReceiptFilterParams represents filtering parameters for goods receipt queries
type GoodsReceiptFilterParams struct {
	POID           *int           `json:"po_id,omitempty" form:"po_id"`
	ReceivedBy     *int           `json:"received_by,omitempty" form:"received_by"`
	ReceiptStatus  *ReceiptStatus `json:"receipt_status,omitempty" form:"receipt_status"`
	DateFrom       *time.Time     `json:"date_from,omitempty" form:"date_from"`
	DateTo         *time.Time     `json:"date_to,omitempty" form:"date_to"`
	Search         string         `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// ConditionReceived represents the condition of received goods
type ConditionReceived string

const (
	ConditionGood      ConditionReceived = "good"
	ConditionDamaged   ConditionReceived = "damaged"
	ConditionExpired   ConditionReceived = "expired"
	ConditionWrongItem ConditionReceived = "wrong_item"
)

// IsValid checks if the condition is valid
func (c ConditionReceived) IsValid() bool {
	switch c {
	case ConditionGood, ConditionDamaged, ConditionExpired, ConditionWrongItem:
		return true
	default:
		return false
	}
}

// String returns the string representation of the condition
func (c ConditionReceived) String() string {
	return string(c)
}

// Value implements the driver.Valuer interface for ConditionReceived
func (c ConditionReceived) Value() (driver.Value, error) {
	return string(c), nil
}

// Scan implements the sql.Scanner interface for ConditionReceived
func (c *ConditionReceived) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch s := value.(type) {
	case string:
		*c = ConditionReceived(s)
	case []byte:
		*c = ConditionReceived(s)
	default:
		return fmt.Errorf("cannot scan %T into ConditionReceived", value)
	}
	return nil
}

// GoodsReceiptDetail represents a line item in a goods receipt
type GoodsReceiptDetail struct {
	ReceiptDetailID   int               `json:"receipt_detail_id" db:"receipt_detail_id"`
	ReceiptID         int               `json:"receipt_id" db:"receipt_id"`
	PODetailID        int               `json:"po_detail_id" db:"po_detail_id"`
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
}

// GoodsReceiptDetailCreateRequest represents a request to create a goods receipt detail
type GoodsReceiptDetailCreateRequest struct {
	PODetailID        int                `json:"po_detail_id" binding:"required,min=1"`
	ProductID         int                `json:"product_id" binding:"required,min=1"`
	QuantityReceived  int                `json:"quantity_received" binding:"required,min=0"`
	QuantityAccepted  int                `json:"quantity_accepted" binding:"required,min=0"`
	QuantityRejected  int                `json:"quantity_rejected" binding:"min=0"`
	UnitCost          float64            `json:"unit_cost" binding:"required,min=0"`
	ConditionReceived ConditionReceived  `json:"condition_received" binding:"required"`
	InspectionNotes   *string            `json:"inspection_notes,omitempty"`
	RejectionReason   *string            `json:"rejection_reason,omitempty"`
	ExpiryDate        *time.Time         `json:"expiry_date,omitempty"`
	BatchNumber       *string            `json:"batch_number,omitempty" binding:"omitempty,max=100"`
	SerialNumbersJSON *string            `json:"serial_numbers_json,omitempty"`
}

// UpdateTotalCost calculates and updates the total cost
func (grd *GoodsReceiptDetail) UpdateTotalCost() {
	grd.TotalCost = float64(grd.QuantityAccepted) * grd.UnitCost
}

// ValidateQuantities ensures quantity relationships are correct
func (grd *GoodsReceiptDetail) ValidateQuantities() bool {
	return grd.QuantityReceived == (grd.QuantityAccepted + grd.QuantityRejected)
}

// HasDiscrepancy checks if there are any discrepancies in the receipt
func (grd *GoodsReceiptDetail) HasDiscrepancy() bool {
	return grd.QuantityRejected > 0 || grd.ConditionReceived != ConditionGood
}

// IsFullyAccepted checks if all received quantity is accepted
func (grd *GoodsReceiptDetail) IsFullyAccepted() bool {
	return grd.QuantityReceived > 0 && grd.QuantityRejected == 0 && grd.ConditionReceived == ConditionGood
}

// UpdateStatus updates the receipt status based on details
func (gr *GoodsReceipt) UpdateStatus(hasDiscrepancy bool, isComplete bool) {
	if hasDiscrepancy {
		gr.ReceiptStatus = ReceiptStatusWithDiscrepancy
	} else if isComplete {
		gr.ReceiptStatus = ReceiptStatusComplete
	} else {
		gr.ReceiptStatus = ReceiptStatusPartial
	}
}

// CalculateTotalValue calculates the total received value
func (gr *GoodsReceipt) CalculateTotalValue(details []GoodsReceiptDetail) {
	total := 0.0
	for _, detail := range details {
		total += detail.TotalCost
	}
	gr.TotalReceivedValue = total
}