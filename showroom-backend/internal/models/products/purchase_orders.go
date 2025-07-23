package products

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// POType represents the type of purchase order
type POType string

const (
	POTypeRegular  POType = "regular"
	POTypeUrgent   POType = "urgent"
	POTypeBlanket  POType = "blanket"
	POTypeContract POType = "contract"
)

// IsValid checks if the PO type is valid
func (t POType) IsValid() bool {
	switch t {
	case POTypeRegular, POTypeUrgent, POTypeBlanket, POTypeContract:
		return true
	default:
		return false
	}
}

// String returns the string representation of the PO type
func (t POType) String() string {
	return string(t)
}

// Value implements the driver.Valuer interface for POType
func (t POType) Value() (driver.Value, error) {
	return string(t), nil
}

// Scan implements the sql.Scanner interface for POType
func (t *POType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch s := value.(type) {
	case string:
		*t = POType(s)
	case []byte:
		*t = POType(s)
	default:
		return fmt.Errorf("cannot scan %T into POType", value)
	}
	return nil
}

// POStatus represents the status of a purchase order
type POStatus string

const (
	POStatusDraft           POStatus = "draft"
	POStatusSent            POStatus = "sent"
	POStatusAcknowledged    POStatus = "acknowledged"
	POStatusPartialReceived POStatus = "partial_received"
	POStatusReceived        POStatus = "received"
	POStatusCompleted       POStatus = "completed"
	POStatusCancelled       POStatus = "cancelled"
)

// IsValid checks if the PO status is valid
func (s POStatus) IsValid() bool {
	switch s {
	case POStatusDraft, POStatusSent, POStatusAcknowledged, POStatusPartialReceived, POStatusReceived, POStatusCompleted, POStatusCancelled:
		return true
	default:
		return false
	}
}

// String returns the string representation of the PO status
func (s POStatus) String() string {
	return string(s)
}

// Value implements the driver.Valuer interface for POStatus
func (s POStatus) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan implements the sql.Scanner interface for POStatus
func (s *POStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch str := value.(type) {
	case string:
		*s = POStatus(str)
	case []byte:
		*s = POStatus(str)
	default:
		return fmt.Errorf("cannot scan %T into POStatus", value)
	}
	return nil
}

// PaymentTerms represents payment terms for purchase orders
type PaymentTerms string

const (
	PaymentTermsCOD    PaymentTerms = "cod"
	PaymentTermsNet30  PaymentTerms = "net_30"
	PaymentTermsNet60  PaymentTerms = "net_60"
	PaymentTermsAdvance PaymentTerms = "advance"
)

// IsValid checks if the payment terms are valid
func (p PaymentTerms) IsValid() bool {
	switch p {
	case PaymentTermsCOD, PaymentTermsNet30, PaymentTermsNet60, PaymentTermsAdvance:
		return true
	default:
		return false
	}
}

// String returns the string representation of the payment terms
func (p PaymentTerms) String() string {
	return string(p)
}

// Value implements the driver.Valuer interface for PaymentTerms
func (p PaymentTerms) Value() (driver.Value, error) {
	return string(p), nil
}

// Scan implements the sql.Scanner interface for PaymentTerms
func (p *PaymentTerms) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch s := value.(type) {
	case string:
		*p = PaymentTerms(s)
	case []byte:
		*p = PaymentTerms(s)
	default:
		return fmt.Errorf("cannot scan %T into PaymentTerms", value)
	}
	return nil
}

// PurchaseOrderParts represents a purchase order for spare parts
type PurchaseOrderParts struct {
	POID                 int           `json:"po_id" db:"po_id"`
	PONumber             string        `json:"po_number" db:"po_number"`
	SupplierID           int           `json:"supplier_id" db:"supplier_id"`
	PODate               time.Time     `json:"po_date" db:"po_date"`
	RequiredDate         *time.Time    `json:"required_date,omitempty" db:"required_date"`
	ExpectedDeliveryDate *time.Time    `json:"expected_delivery_date,omitempty" db:"expected_delivery_date"`
	POType               POType        `json:"po_type" db:"po_type"`
	Subtotal             float64       `json:"subtotal" db:"subtotal"`
	TaxAmount            float64       `json:"tax_amount" db:"tax_amount"`
	DiscountAmount       float64       `json:"discount_amount" db:"discount_amount"`
	ShippingCost         float64       `json:"shipping_cost" db:"shipping_cost"`
	TotalAmount          float64       `json:"total_amount" db:"total_amount"`
	Status               POStatus      `json:"status" db:"status"`
	PaymentTerms         PaymentTerms  `json:"payment_terms" db:"payment_terms"`
	PaymentDueDate       *time.Time    `json:"payment_due_date,omitempty" db:"payment_due_date"`
	CreatedBy            int           `json:"created_by" db:"created_by"`
	ApprovedBy           *int          `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt           *time.Time    `json:"approved_at,omitempty" db:"approved_at"`
	DeliveryAddress      *string       `json:"delivery_address,omitempty" db:"delivery_address"`
	PONotes              *string       `json:"po_notes,omitempty" db:"po_notes"`
	TermsAndConditions   *string       `json:"terms_and_conditions,omitempty" db:"terms_and_conditions"`
	CreatedAt            time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at" db:"updated_at"`
}

// PurchaseOrderPartsListItem represents a simplified purchase order for list views
type PurchaseOrderPartsListItem struct {
	POID                 int          `json:"po_id" db:"po_id"`
	PONumber             string       `json:"po_number" db:"po_number"`
	SupplierID           int          `json:"supplier_id" db:"supplier_id"`
	PODate               time.Time    `json:"po_date" db:"po_date"`
	RequiredDate         *time.Time   `json:"required_date,omitempty" db:"required_date"`
	ExpectedDeliveryDate *time.Time   `json:"expected_delivery_date,omitempty" db:"expected_delivery_date"`
	POType               POType       `json:"po_type" db:"po_type"`
	TotalAmount          float64      `json:"total_amount" db:"total_amount"`
	Status               POStatus     `json:"status" db:"status"`
	PaymentTerms         PaymentTerms `json:"payment_terms" db:"payment_terms"`
	CreatedAt            time.Time    `json:"created_at" db:"created_at"`
}

// PurchaseOrderPartsCreateRequest represents a request to create a purchase order
type PurchaseOrderPartsCreateRequest struct {
	SupplierID           int           `json:"supplier_id" binding:"required,min=1"`
	RequiredDate         *time.Time    `json:"required_date,omitempty"`
	ExpectedDeliveryDate *time.Time    `json:"expected_delivery_date,omitempty"`
	POType               POType        `json:"po_type" binding:"required"`
	PaymentTerms         PaymentTerms  `json:"payment_terms" binding:"required"`
	DeliveryAddress      *string       `json:"delivery_address,omitempty" binding:"omitempty,max=500"`
	PONotes              *string       `json:"po_notes,omitempty"`
	TermsAndConditions   *string       `json:"terms_and_conditions,omitempty"`
}

// PurchaseOrderPartsUpdateRequest represents a request to update a purchase order
type PurchaseOrderPartsUpdateRequest struct {
	SupplierID           *int          `json:"supplier_id,omitempty" binding:"omitempty,min=1"`
	RequiredDate         *time.Time    `json:"required_date,omitempty"`
	ExpectedDeliveryDate *time.Time    `json:"expected_delivery_date,omitempty"`
	POType               *POType       `json:"po_type,omitempty"`
	TaxAmount            *float64      `json:"tax_amount,omitempty" binding:"omitempty,min=0"`
	DiscountAmount       *float64      `json:"discount_amount,omitempty" binding:"omitempty,min=0"`
	ShippingCost         *float64      `json:"shipping_cost,omitempty" binding:"omitempty,min=0"`
	PaymentTerms         *PaymentTerms `json:"payment_terms,omitempty"`
	DeliveryAddress      *string       `json:"delivery_address,omitempty" binding:"omitempty,max=500"`
	PONotes              *string       `json:"po_notes,omitempty"`
	TermsAndConditions   *string       `json:"terms_and_conditions,omitempty"`
}

// PurchaseOrderPartsFilterParams represents filtering parameters for purchase order queries
type PurchaseOrderPartsFilterParams struct {
	SupplierID       *int         `json:"supplier_id,omitempty" form:"supplier_id"`
	Status           *POStatus    `json:"status,omitempty" form:"status"`
	POType           *POType      `json:"po_type,omitempty" form:"po_type"`
	PaymentTerms     *PaymentTerms `json:"payment_terms,omitempty" form:"payment_terms"`
	DateFrom         *time.Time   `json:"date_from,omitempty" form:"date_from"`
	DateTo           *time.Time   `json:"date_to,omitempty" form:"date_to"`
	CreatedBy        *int         `json:"created_by,omitempty" form:"created_by"`
	ApprovedBy       *int         `json:"approved_by,omitempty" form:"approved_by"`
	MinAmount        *float64     `json:"min_amount,omitempty" form:"min_amount"`
	MaxAmount        *float64     `json:"max_amount,omitempty" form:"max_amount"`
	Search           string       `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// CanEdit checks if the purchase order can be edited
func (po *PurchaseOrderParts) CanEdit() bool {
	return po.Status == POStatusDraft
}

// CanCancel checks if the purchase order can be cancelled
func (po *PurchaseOrderParts) CanCancel() bool {
	return po.Status == POStatusDraft || po.Status == POStatusSent || po.Status == POStatusAcknowledged
}

// CanApprove checks if the purchase order can be approved
func (po *PurchaseOrderParts) CanApprove() bool {
	return po.Status == POStatusDraft && po.ApprovedBy == nil
}

// CanSend checks if the purchase order can be sent to supplier
func (po *PurchaseOrderParts) CanSend() bool {
	return po.Status == POStatusDraft && po.ApprovedBy != nil
}

// IsApproved checks if the purchase order is approved
func (po *PurchaseOrderParts) IsApproved() bool {
	return po.ApprovedBy != nil && po.ApprovedAt != nil
}

// CalculateTotals calculates subtotal and total amount
func (po *PurchaseOrderParts) CalculateTotals() {
	po.TotalAmount = po.Subtotal + po.TaxAmount + po.ShippingCost - po.DiscountAmount
}

// SetPaymentDueDate sets the payment due date based on payment terms
func (po *PurchaseOrderParts) SetPaymentDueDate() {
	if po.PaymentDueDate != nil {
		return // Already set
	}
	
	baseDate := po.PODate
	if po.RequiredDate != nil {
		baseDate = *po.RequiredDate
	}
	
	switch po.PaymentTerms {
	case PaymentTermsCOD:
		po.PaymentDueDate = &baseDate
	case PaymentTermsNet30:
		dueDate := baseDate.AddDate(0, 0, 30)
		po.PaymentDueDate = &dueDate
	case PaymentTermsNet60:
		dueDate := baseDate.AddDate(0, 0, 60)
		po.PaymentDueDate = &dueDate
	case PaymentTermsAdvance:
		// Due immediately for advance payment
		po.PaymentDueDate = &baseDate
	}
}