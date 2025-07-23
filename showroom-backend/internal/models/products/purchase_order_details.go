package products

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// LineStatus represents the status of a purchase order line item
type LineStatus string

const (
	LineStatusPending   LineStatus = "pending"
	LineStatusPartial   LineStatus = "partial"
	LineStatusReceived  LineStatus = "received"
	LineStatusCancelled LineStatus = "cancelled"
)

// IsValid checks if the line status is valid
func (s LineStatus) IsValid() bool {
	switch s {
	case LineStatusPending, LineStatusPartial, LineStatusReceived, LineStatusCancelled:
		return true
	default:
		return false
	}
}

// String returns the string representation of the line status
func (s LineStatus) String() string {
	return string(s)
}

// Value implements the driver.Valuer interface for LineStatus
func (s LineStatus) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan implements the sql.Scanner interface for LineStatus
func (s *LineStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch str := value.(type) {
	case string:
		*s = LineStatus(str)
	case []byte:
		*s = LineStatus(str)
	default:
		return fmt.Errorf("cannot scan %T into LineStatus", value)
	}
	return nil
}

// PurchaseOrderDetail represents a line item in a purchase order
type PurchaseOrderDetail struct {
	PODetailID       int         `json:"po_detail_id" db:"po_detail_id"`
	POID             int         `json:"po_id" db:"po_id"`
	ProductID        int         `json:"product_id" db:"product_id"`
	ItemDescription  *string     `json:"item_description,omitempty" db:"item_description"`
	QuantityOrdered  int         `json:"quantity_ordered" db:"quantity_ordered"`
	QuantityReceived int         `json:"quantity_received" db:"quantity_received"`
	QuantityPending  int         `json:"quantity_pending" db:"quantity_pending"`
	UnitCost         float64     `json:"unit_cost" db:"unit_cost"`
	TotalCost        float64     `json:"total_cost" db:"total_cost"`
	ExpectedDate     *time.Time  `json:"expected_date,omitempty" db:"expected_date"`
	ReceivedDate     *time.Time  `json:"received_date,omitempty" db:"received_date"`
	LineStatus       LineStatus  `json:"line_status" db:"line_status"`
	ItemNotes        *string     `json:"item_notes,omitempty" db:"item_notes"`
}

// PurchaseOrderDetailListItem represents a simplified purchase order detail for list views
type PurchaseOrderDetailListItem struct {
	PODetailID       int        `json:"po_detail_id" db:"po_detail_id"`
	POID             int        `json:"po_id" db:"po_id"`
	ProductID        int        `json:"product_id" db:"product_id"`
	ItemDescription  *string    `json:"item_description,omitempty" db:"item_description"`
	QuantityOrdered  int        `json:"quantity_ordered" db:"quantity_ordered"`
	QuantityReceived int        `json:"quantity_received" db:"quantity_received"`
	QuantityPending  int        `json:"quantity_pending" db:"quantity_pending"`
	UnitCost         float64    `json:"unit_cost" db:"unit_cost"`
	TotalCost        float64    `json:"total_cost" db:"total_cost"`
	LineStatus       LineStatus `json:"line_status" db:"line_status"`
}

// PurchaseOrderDetailCreateRequest represents a request to create a purchase order detail
type PurchaseOrderDetailCreateRequest struct {
	ProductID       int      `json:"product_id" binding:"required,min=1"`
	ItemDescription *string  `json:"item_description,omitempty" binding:"omitempty,max=500"`
	QuantityOrdered int      `json:"quantity_ordered" binding:"required,min=1"`
	UnitCost        float64  `json:"unit_cost" binding:"required,min=0"`
	ExpectedDate    *time.Time `json:"expected_date,omitempty"`
	ItemNotes       *string  `json:"item_notes,omitempty"`
}

// PurchaseOrderDetailUpdateRequest represents a request to update a purchase order detail
type PurchaseOrderDetailUpdateRequest struct {
	ProductID       *int       `json:"product_id,omitempty" binding:"omitempty,min=1"`
	ItemDescription *string    `json:"item_description,omitempty" binding:"omitempty,max=500"`
	QuantityOrdered *int       `json:"quantity_ordered,omitempty" binding:"omitempty,min=1"`
	UnitCost        *float64   `json:"unit_cost,omitempty" binding:"omitempty,min=0"`
	ExpectedDate    *time.Time `json:"expected_date,omitempty"`
	ItemNotes       *string    `json:"item_notes,omitempty"`
}

// PurchaseOrderDetailFilterParams represents filtering parameters for purchase order detail queries
type PurchaseOrderDetailFilterParams struct {
	POID       *int        `json:"po_id,omitempty" form:"po_id"`
	ProductID  *int        `json:"product_id,omitempty" form:"product_id"`
	LineStatus *LineStatus `json:"line_status,omitempty" form:"line_status"`
	Search     string      `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// ReceiveQuantity updates the received quantity and updates status accordingly
func (pod *PurchaseOrderDetail) ReceiveQuantity(quantity int) {
	if quantity <= 0 {
		return
	}
	
	// Ensure we don't receive more than ordered
	maxReceivable := pod.QuantityOrdered - pod.QuantityReceived
	if quantity > maxReceivable {
		quantity = maxReceivable
	}
	
	pod.QuantityReceived += quantity
	pod.QuantityPending = pod.QuantityOrdered - pod.QuantityReceived
	
	// Update status based on quantities
	pod.UpdateLineStatus()
	
	// Set received date if first receipt
	if pod.ReceivedDate == nil && pod.QuantityReceived > 0 {
		now := time.Now()
		pod.ReceivedDate = &now
	}
}

// UpdateLineStatus updates the line status based on quantities
func (pod *PurchaseOrderDetail) UpdateLineStatus() {
	if pod.QuantityReceived == 0 {
		pod.LineStatus = LineStatusPending
	} else if pod.QuantityReceived < pod.QuantityOrdered {
		pod.LineStatus = LineStatusPartial
	} else {
		pod.LineStatus = LineStatusReceived
	}
}

// CalculateTotalCost calculates total cost based on quantity and unit cost
func (pod *PurchaseOrderDetail) CalculateTotalCost() {
	pod.TotalCost = float64(pod.QuantityOrdered) * pod.UnitCost
}

// CanReceiveMore checks if more quantity can be received
func (pod *PurchaseOrderDetail) CanReceiveMore() bool {
	return pod.LineStatus != LineStatusReceived && pod.LineStatus != LineStatusCancelled && pod.QuantityPending > 0
}

// GetCompletionPercentage returns the completion percentage of the line item
func (pod *PurchaseOrderDetail) GetCompletionPercentage() float64 {
	if pod.QuantityOrdered == 0 {
		return 0
	}
	return (float64(pod.QuantityReceived) / float64(pod.QuantityOrdered)) * 100
}

// IsFullyReceived checks if the line item is fully received
func (pod *PurchaseOrderDetail) IsFullyReceived() bool {
	return pod.QuantityReceived >= pod.QuantityOrdered
}

// IsPartiallyReceived checks if the line item is partially received
func (pod *PurchaseOrderDetail) IsPartiallyReceived() bool {
	return pod.QuantityReceived > 0 && pod.QuantityReceived < pod.QuantityOrdered
}

// HasPendingQuantity checks if there is pending quantity to receive
func (pod *PurchaseOrderDetail) HasPendingQuantity() bool {
	return pod.QuantityPending > 0
}