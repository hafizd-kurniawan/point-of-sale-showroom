package inventory

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// PurchaseOrderType represents the type of purchase order
type PurchaseOrderType string

const (
	POTypeRegular  PurchaseOrderType = "regular"
	POTypeUrgent   PurchaseOrderType = "urgent"
	POTypeBlanket  PurchaseOrderType = "blanket"
	POTypeContract PurchaseOrderType = "contract"
)

// PurchaseOrderStatus represents the status of purchase order
type PurchaseOrderStatus string

const (
	POStatusDraft            PurchaseOrderStatus = "draft"
	POStatusSent             PurchaseOrderStatus = "sent"
	POStatusAcknowledged     PurchaseOrderStatus = "acknowledged"
	POStatusPartialReceived  PurchaseOrderStatus = "partial_received"
	POStatusReceived         PurchaseOrderStatus = "received"
	POStatusCompleted        PurchaseOrderStatus = "completed"
	POStatusCancelled        PurchaseOrderStatus = "cancelled"
)

// PaymentTerms represents payment terms for purchase orders
type PaymentTerms string

const (
	PaymentTermsCOD    PaymentTerms = "cod"
	PaymentTermsNet30  PaymentTerms = "net_30"
	PaymentTermsNet60  PaymentTerms = "net_60"
	PaymentTermsAdvance PaymentTerms = "advance"
)

// PurchaseOrderPart represents a purchase order in the system
type PurchaseOrderPart struct {
	PoID                 int                 `json:"po_id" db:"po_id"`
	PoNumber             string              `json:"po_number" db:"po_number"`
	SupplierID           int                 `json:"supplier_id" db:"supplier_id"`
	PoDate               time.Time           `json:"po_date" db:"po_date"`
	RequiredDate         *time.Time          `json:"required_date,omitempty" db:"required_date"`
	ExpectedDeliveryDate *time.Time          `json:"expected_delivery_date,omitempty" db:"expected_delivery_date"`
	PoType               PurchaseOrderType   `json:"po_type" db:"po_type"`
	Subtotal             float64             `json:"subtotal" db:"subtotal"`
	TaxAmount            float64             `json:"tax_amount" db:"tax_amount"`
	DiscountAmount       float64             `json:"discount_amount" db:"discount_amount"`
	ShippingCost         float64             `json:"shipping_cost" db:"shipping_cost"`
	TotalAmount          float64             `json:"total_amount" db:"total_amount"`
	Status               PurchaseOrderStatus `json:"status" db:"status"`
	PaymentTerms         PaymentTerms        `json:"payment_terms" db:"payment_terms"`
	PaymentDueDate       *time.Time          `json:"payment_due_date,omitempty" db:"payment_due_date"`
	CreatedBy            int                 `json:"created_by" db:"created_by"`
	ApprovedBy           *int                `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt           *time.Time          `json:"approved_at,omitempty" db:"approved_at"`
	DeliveryAddress      *string             `json:"delivery_address,omitempty" db:"delivery_address"`
	PoNotes              *string             `json:"po_notes,omitempty" db:"po_notes"`
	TermsAndConditions   *string             `json:"terms_and_conditions,omitempty" db:"terms_and_conditions"`
	CreatedAt            time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at" db:"updated_at"`

	// Related data
	SupplierName     string                   `json:"supplier_name,omitempty" db:"supplier_name"`
	CreatedByName    string                   `json:"created_by_name,omitempty" db:"created_by_name"`
	ApprovedByName   *string                  `json:"approved_by_name,omitempty" db:"approved_by_name"`
	Details          []PurchaseOrderDetail    `json:"details,omitempty"`
}

// PurchaseOrderPartListItem represents a simplified purchase order for list views
type PurchaseOrderPartListItem struct {
	PoID                 int                 `json:"po_id" db:"po_id"`
	PoNumber             string              `json:"po_number" db:"po_number"`
	SupplierName         string              `json:"supplier_name" db:"supplier_name"`
	PoDate               time.Time           `json:"po_date" db:"po_date"`
	ExpectedDeliveryDate *time.Time          `json:"expected_delivery_date,omitempty" db:"expected_delivery_date"`
	PoType               PurchaseOrderType   `json:"po_type" db:"po_type"`
	TotalAmount          float64             `json:"total_amount" db:"total_amount"`
	Status               PurchaseOrderStatus `json:"status" db:"status"`
	PaymentTerms         PaymentTerms        `json:"payment_terms" db:"payment_terms"`
	CreatedByName        string              `json:"created_by_name" db:"created_by_name"`
	CreatedAt            time.Time           `json:"created_at" db:"created_at"`
}

// PurchaseOrderPartCreateRequest represents a request to create a purchase order
type PurchaseOrderPartCreateRequest struct {
	SupplierID           int                 `json:"supplier_id" binding:"required"`
	RequiredDate         *time.Time          `json:"required_date,omitempty"`
	ExpectedDeliveryDate *time.Time          `json:"expected_delivery_date,omitempty"`
	PoType               PurchaseOrderType   `json:"po_type" binding:"required,oneof=regular urgent blanket contract"`
	PaymentTerms         PaymentTerms        `json:"payment_terms" binding:"required,oneof=cod net_30 net_60 advance"`
	DeliveryAddress      *string             `json:"delivery_address,omitempty" binding:"omitempty,max=500"`
	PoNotes              *string             `json:"po_notes,omitempty"`
	TermsAndConditions   *string             `json:"terms_and_conditions,omitempty"`
	Details              []PurchaseOrderDetailCreateRequest `json:"details" binding:"required,min=1"`
}

// PurchaseOrderPartUpdateRequest represents a request to update a purchase order
type PurchaseOrderPartUpdateRequest struct {
	RequiredDate         *time.Time          `json:"required_date,omitempty"`
	ExpectedDeliveryDate *time.Time          `json:"expected_delivery_date,omitempty"`
	PoType               *PurchaseOrderType  `json:"po_type,omitempty" binding:"omitempty,oneof=regular urgent blanket contract"`
	PaymentTerms         *PaymentTerms       `json:"payment_terms,omitempty" binding:"omitempty,oneof=cod net_30 net_60 advance"`
	DeliveryAddress      *string             `json:"delivery_address,omitempty" binding:"omitempty,max=500"`
	PoNotes              *string             `json:"po_notes,omitempty"`
	TermsAndConditions   *string             `json:"terms_and_conditions,omitempty"`
}

// PurchaseOrderPartFilterParams represents filtering parameters for purchase order queries
type PurchaseOrderPartFilterParams struct {
	SupplierID    *int                `json:"supplier_id,omitempty" form:"supplier_id"`
	PoType        *PurchaseOrderType  `json:"po_type,omitempty" form:"po_type"`
	Status        *PurchaseOrderStatus `json:"status,omitempty" form:"status"`
	PaymentTerms  *PaymentTerms       `json:"payment_terms,omitempty" form:"payment_terms"`
	CreatedBy     *int                `json:"created_by,omitempty" form:"created_by"`
	ApprovedBy    *int                `json:"approved_by,omitempty" form:"approved_by"`
	DateFrom      *time.Time          `json:"date_from,omitempty" form:"date_from"`
	DateTo        *time.Time          `json:"date_to,omitempty" form:"date_to"`
	Search        string              `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// LineStatus represents the status of purchase order line item
type LineStatus string

const (
	LineStatusPending   LineStatus = "pending"
	LineStatusPartial   LineStatus = "partial"
	LineStatusReceived  LineStatus = "received"
	LineStatusCancelled LineStatus = "cancelled"
)

// PurchaseOrderDetail represents a purchase order detail line
type PurchaseOrderDetail struct {
	PoDetailID       int        `json:"po_detail_id" db:"po_detail_id"`
	PoID             int        `json:"po_id" db:"po_id"`
	ProductID        int        `json:"product_id" db:"product_id"`
	ItemDescription  *string    `json:"item_description,omitempty" db:"item_description"`
	QuantityOrdered  int        `json:"quantity_ordered" db:"quantity_ordered"`
	QuantityReceived int        `json:"quantity_received" db:"quantity_received"`
	QuantityPending  int        `json:"quantity_pending" db:"quantity_pending"`
	UnitCost         float64    `json:"unit_cost" db:"unit_cost"`
	TotalCost        float64    `json:"total_cost" db:"total_cost"`
	ExpectedDate     *time.Time `json:"expected_date,omitempty" db:"expected_date"`
	ReceivedDate     *time.Time `json:"received_date,omitempty" db:"received_date"`
	LineStatus       LineStatus `json:"line_status" db:"line_status"`
	ItemNotes        *string    `json:"item_notes,omitempty" db:"item_notes"`

	// Related data
	ProductCode string  `json:"product_code,omitempty" db:"product_code"`
	ProductName string  `json:"product_name,omitempty" db:"product_name"`
	UnitMeasure string  `json:"unit_measure,omitempty" db:"unit_measure"`
}

// PurchaseOrderDetailCreateRequest represents a request to create a purchase order detail
type PurchaseOrderDetailCreateRequest struct {
	ProductID       int        `json:"product_id" binding:"required"`
	ItemDescription *string    `json:"item_description,omitempty" binding:"omitempty,max=500"`
	QuantityOrdered int        `json:"quantity_ordered" binding:"required,min=1"`
	UnitCost        float64    `json:"unit_cost" binding:"required,min=0"`
	ExpectedDate    *time.Time `json:"expected_date,omitempty"`
	ItemNotes       *string    `json:"item_notes,omitempty"`
}

// PurchaseOrderDetailUpdateRequest represents a request to update a purchase order detail
type PurchaseOrderDetailUpdateRequest struct {
	ItemDescription *string    `json:"item_description,omitempty" binding:"omitempty,max=500"`
	QuantityOrdered *int       `json:"quantity_ordered,omitempty" binding:"omitempty,min=1"`
	UnitCost        *float64   `json:"unit_cost,omitempty" binding:"omitempty,min=0"`
	ExpectedDate    *time.Time `json:"expected_date,omitempty"`
	ItemNotes       *string    `json:"item_notes,omitempty"`
}

// Methods for PurchaseOrderPart

// CanBeModified checks if the purchase order can be modified
func (po *PurchaseOrderPart) CanBeModified() bool {
	return po.Status == POStatusDraft
}

// CanBeSent checks if the purchase order can be sent to supplier
func (po *PurchaseOrderPart) CanBeSent() bool {
	return po.Status == POStatusDraft && len(po.Details) > 0
}

// CanBeCancelled checks if the purchase order can be cancelled
func (po *PurchaseOrderPart) CanBeCancelled() bool {
	return po.Status == POStatusDraft || po.Status == POStatusSent || po.Status == POStatusAcknowledged
}

// CalculateTotals calculates and updates the purchase order totals
func (po *PurchaseOrderPart) CalculateTotals() {
	po.Subtotal = 0
	for _, detail := range po.Details {
		po.Subtotal += detail.TotalCost
	}
	po.TotalAmount = po.Subtotal + po.TaxAmount + po.ShippingCost - po.DiscountAmount
}

// UpdatePaymentDueDate updates payment due date based on payment terms
func (po *PurchaseOrderPart) UpdatePaymentDueDate() {
	switch po.PaymentTerms {
	case PaymentTermsCOD:
		po.PaymentDueDate = &po.PoDate
	case PaymentTermsNet30:
		dueDate := po.PoDate.AddDate(0, 0, 30)
		po.PaymentDueDate = &dueDate
	case PaymentTermsNet60:
		dueDate := po.PoDate.AddDate(0, 0, 60)
		po.PaymentDueDate = &dueDate
	case PaymentTermsAdvance:
		// Payment before delivery
		po.PaymentDueDate = &po.PoDate
	}
}

// Methods for PurchaseOrderDetail

// CalculateTotalCost calculates and updates the total cost for the line item
func (pod *PurchaseOrderDetail) CalculateTotalCost() {
	pod.TotalCost = float64(pod.QuantityOrdered) * pod.UnitCost
}

// UpdatePendingQuantity updates the pending quantity based on ordered and received
func (pod *PurchaseOrderDetail) UpdatePendingQuantity() {
	pod.QuantityPending = pod.QuantityOrdered - pod.QuantityReceived
	if pod.QuantityPending < 0 {
		pod.QuantityPending = 0
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

// IsFullyReceived checks if the line item is fully received
func (pod *PurchaseOrderDetail) IsFullyReceived() bool {
	return pod.QuantityReceived >= pod.QuantityOrdered
}

// IsPartiallyReceived checks if the line item is partially received
func (pod *PurchaseOrderDetail) IsPartiallyReceived() bool {
	return pod.QuantityReceived > 0 && pod.QuantityReceived < pod.QuantityOrdered
}