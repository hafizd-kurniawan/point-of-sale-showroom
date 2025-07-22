package inventory

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// PaymentMethod represents the method of payment
type PaymentMethod string

const (
	PaymentMethodCash     PaymentMethod = "cash"
	PaymentMethodTransfer PaymentMethod = "transfer"
	PaymentMethodCheck    PaymentMethod = "check"
	PaymentMethodCredit   PaymentMethod = "credit"
)

// PaymentStatus represents the status of payment
type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPartial  PaymentStatus = "partial"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusOverdue  PaymentStatus = "overdue"
	PaymentStatusDisputed PaymentStatus = "disputed"
)

// SupplierPayment represents a supplier payment record in the system
type SupplierPayment struct {
	PaymentID         int           `json:"payment_id" db:"payment_id"`
	SupplierID        int           `json:"supplier_id" db:"supplier_id"`
	PoID              *int          `json:"po_id,omitempty" db:"po_id"`
	PaymentNumber     string        `json:"payment_number" db:"payment_number"`
	InvoiceAmount     float64       `json:"invoice_amount" db:"invoice_amount"`
	PaymentAmount     float64       `json:"payment_amount" db:"payment_amount"`
	DiscountTaken     float64       `json:"discount_taken" db:"discount_taken"`
	OutstandingAmount float64       `json:"outstanding_amount" db:"outstanding_amount"`
	InvoiceDate       *time.Time    `json:"invoice_date,omitempty" db:"invoice_date"`
	PaymentDate       time.Time     `json:"payment_date" db:"payment_date"`
	DueDate           *time.Time    `json:"due_date,omitempty" db:"due_date"`
	PaymentMethod     PaymentMethod `json:"payment_method" db:"payment_method"`
	PaymentReference  *string       `json:"payment_reference,omitempty" db:"payment_reference"`
	InvoiceNumber     *string       `json:"invoice_number,omitempty" db:"invoice_number"`
	PaymentStatus     PaymentStatus `json:"payment_status" db:"payment_status"`
	DaysOverdue       int           `json:"days_overdue" db:"days_overdue"`
	PenaltyAmount     float64       `json:"penalty_amount" db:"penalty_amount"`
	ProcessedBy       int           `json:"processed_by" db:"processed_by"`
	PaymentNotes      *string       `json:"payment_notes,omitempty" db:"payment_notes"`
	CreatedAt         time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at" db:"updated_at"`

	// Related data
	SupplierName    string  `json:"supplier_name,omitempty" db:"supplier_name"`
	PoNumber        *string `json:"po_number,omitempty" db:"po_number"`
	ProcessedByName string  `json:"processed_by_name,omitempty" db:"processed_by_name"`
}

// SupplierPaymentListItem represents a simplified supplier payment for list views
type SupplierPaymentListItem struct {
	PaymentID         int           `json:"payment_id" db:"payment_id"`
	PaymentNumber     string        `json:"payment_number" db:"payment_number"`
	SupplierName      string        `json:"supplier_name" db:"supplier_name"`
	PoNumber          *string       `json:"po_number,omitempty" db:"po_number"`
	InvoiceNumber     *string       `json:"invoice_number,omitempty" db:"invoice_number"`
	InvoiceAmount     float64       `json:"invoice_amount" db:"invoice_amount"`
	PaymentAmount     float64       `json:"payment_amount" db:"payment_amount"`
	OutstandingAmount float64       `json:"outstanding_amount" db:"outstanding_amount"`
	PaymentDate       time.Time     `json:"payment_date" db:"payment_date"`
	DueDate           *time.Time    `json:"due_date,omitempty" db:"due_date"`
	PaymentMethod     PaymentMethod `json:"payment_method" db:"payment_method"`
	PaymentStatus     PaymentStatus `json:"payment_status" db:"payment_status"`
	DaysOverdue       int           `json:"days_overdue" db:"days_overdue"`
	ProcessedByName   string        `json:"processed_by_name" db:"processed_by_name"`
}

// SupplierPaymentCreateRequest represents a request to create a supplier payment
type SupplierPaymentCreateRequest struct {
	SupplierID       int           `json:"supplier_id" binding:"required"`
	PoID             *int          `json:"po_id,omitempty"`
	InvoiceAmount    float64       `json:"invoice_amount" binding:"required,min=0"`
	PaymentAmount    float64       `json:"payment_amount" binding:"required,min=0"`
	DiscountTaken    float64       `json:"discount_taken" binding:"min=0"`
	InvoiceDate      *time.Time    `json:"invoice_date,omitempty"`
	DueDate          *time.Time    `json:"due_date,omitempty"`
	PaymentMethod    PaymentMethod `json:"payment_method" binding:"required,oneof=cash transfer check credit"`
	PaymentReference *string       `json:"payment_reference,omitempty" binding:"omitempty,max=100"`
	InvoiceNumber    *string       `json:"invoice_number,omitempty" binding:"omitempty,max=100"`
	PaymentNotes     *string       `json:"payment_notes,omitempty"`
}

// SupplierPaymentUpdateRequest represents a request to update a supplier payment
type SupplierPaymentUpdateRequest struct {
	PaymentAmount    *float64       `json:"payment_amount,omitempty" binding:"omitempty,min=0"`
	DiscountTaken    *float64       `json:"discount_taken,omitempty" binding:"omitempty,min=0"`
	DueDate          *time.Time     `json:"due_date,omitempty"`
	PaymentMethod    *PaymentMethod `json:"payment_method,omitempty" binding:"omitempty,oneof=cash transfer check credit"`
	PaymentReference *string        `json:"payment_reference,omitempty" binding:"omitempty,max=100"`
	PaymentStatus    *PaymentStatus `json:"payment_status,omitempty" binding:"omitempty,oneof=pending partial paid overdue disputed"`
	PaymentNotes     *string        `json:"payment_notes,omitempty"`
}

// SupplierPaymentFilterParams represents filtering parameters for supplier payment queries
type SupplierPaymentFilterParams struct {
	SupplierID    *int          `json:"supplier_id,omitempty" form:"supplier_id"`
	PoID          *int          `json:"po_id,omitempty" form:"po_id"`
	PaymentMethod *PaymentMethod `json:"payment_method,omitempty" form:"payment_method"`
	PaymentStatus *PaymentStatus `json:"payment_status,omitempty" form:"payment_status"`
	ProcessedBy   *int          `json:"processed_by,omitempty" form:"processed_by"`
	InvoiceNumber string        `json:"invoice_number,omitempty" form:"invoice_number"`
	IsOverdue     *bool         `json:"is_overdue,omitempty" form:"is_overdue"`
	DateFrom      *time.Time    `json:"date_from,omitempty" form:"date_from"`
	DateTo        *time.Time    `json:"date_to,omitempty" form:"date_to"`
	DueDateFrom   *time.Time    `json:"due_date_from,omitempty" form:"due_date_from"`
	DueDateTo     *time.Time    `json:"due_date_to,omitempty" form:"due_date_to"`
	Search        string        `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// PaymentSummary represents payment summary for a supplier or period
type PaymentSummary struct {
	TotalInvoiceAmount    float64 `json:"total_invoice_amount"`
	TotalPaymentAmount    float64 `json:"total_payment_amount"`
	TotalDiscountTaken    float64 `json:"total_discount_taken"`
	TotalOutstandingAmount float64 `json:"total_outstanding_amount"`
	TotalPenaltyAmount    float64 `json:"total_penalty_amount"`
	PendingPaymentsCount  int     `json:"pending_payments_count"`
	OverduePaymentsCount  int     `json:"overdue_payments_count"`
	AveragePaymentDays    float64 `json:"average_payment_days"`
}

// Methods for SupplierPayment

// CalculateOutstandingAmount calculates and updates the outstanding amount
func (sp *SupplierPayment) CalculateOutstandingAmount() {
	sp.OutstandingAmount = sp.InvoiceAmount - sp.PaymentAmount - sp.DiscountTaken
	if sp.OutstandingAmount < 0 {
		sp.OutstandingAmount = 0
	}
}

// CalculateDaysOverdue calculates and updates days overdue
func (sp *SupplierPayment) CalculateDaysOverdue() {
	if sp.DueDate == nil {
		sp.DaysOverdue = 0
		return
	}

	if sp.PaymentStatus == PaymentStatusPaid {
		sp.DaysOverdue = 0
		return
	}

	now := time.Now()
	if now.After(*sp.DueDate) {
		sp.DaysOverdue = int(now.Sub(*sp.DueDate).Hours() / 24)
	} else {
		sp.DaysOverdue = 0
	}
}

// UpdatePaymentStatus updates the payment status based on amounts and due date
func (sp *SupplierPayment) UpdatePaymentStatus() {
	sp.CalculateOutstandingAmount()
	sp.CalculateDaysOverdue()

	if sp.OutstandingAmount == 0 {
		sp.PaymentStatus = PaymentStatusPaid
	} else if sp.PaymentAmount > 0 {
		if sp.DaysOverdue > 0 {
			sp.PaymentStatus = PaymentStatusOverdue
		} else {
			sp.PaymentStatus = PaymentStatusPartial
		}
	} else {
		if sp.DaysOverdue > 0 {
			sp.PaymentStatus = PaymentStatusOverdue
		} else {
			sp.PaymentStatus = PaymentStatusPending
		}
	}
}

// CalculatePenalty calculates penalty amount based on days overdue
func (sp *SupplierPayment) CalculatePenalty(penaltyRate float64) {
	if sp.DaysOverdue > 0 && sp.OutstandingAmount > 0 {
		// Calculate penalty as percentage of outstanding amount per day
		sp.PenaltyAmount = sp.OutstandingAmount * (penaltyRate / 100) * float64(sp.DaysOverdue) / 365
	} else {
		sp.PenaltyAmount = 0
	}
}

// IsFullyPaid checks if the payment is fully paid
func (sp *SupplierPayment) IsFullyPaid() bool {
	return sp.PaymentStatus == PaymentStatusPaid && sp.OutstandingAmount == 0
}

// IsOverdue checks if the payment is overdue
func (sp *SupplierPayment) IsOverdue() bool {
	return sp.PaymentStatus == PaymentStatusOverdue || sp.DaysOverdue > 0
}

// IsPartiallyPaid checks if the payment is partially paid
func (sp *SupplierPayment) IsPartiallyPaid() bool {
	return sp.PaymentStatus == PaymentStatusPartial && sp.PaymentAmount > 0 && sp.OutstandingAmount > 0
}

// IsPending checks if the payment is pending
func (sp *SupplierPayment) IsPending() bool {
	return sp.PaymentStatus == PaymentStatusPending && sp.PaymentAmount == 0
}

// IsDisputed checks if the payment is disputed
func (sp *SupplierPayment) IsDisputed() bool {
	return sp.PaymentStatus == PaymentStatusDisputed
}

// CanBeModified checks if the payment can be modified
func (sp *SupplierPayment) CanBeModified() bool {
	return sp.PaymentStatus != PaymentStatusPaid && sp.PaymentStatus != PaymentStatusDisputed
}

// GetPaymentPercentage calculates the percentage of invoice amount paid
func (sp *SupplierPayment) GetPaymentPercentage() float64 {
	if sp.InvoiceAmount == 0 {
		return 0
	}
	return (sp.PaymentAmount / sp.InvoiceAmount) * 100
}

// GetDiscountPercentage calculates the percentage of discount taken
func (sp *SupplierPayment) GetDiscountPercentage() float64 {
	if sp.InvoiceAmount == 0 {
		return 0
	}
	return (sp.DiscountTaken / sp.InvoiceAmount) * 100
}

// GetTotalPaidAmount returns the total amount paid including discounts
func (sp *SupplierPayment) GetTotalPaidAmount() float64 {
	return sp.PaymentAmount + sp.DiscountTaken
}

// GetEffectiveAmount returns the invoice amount minus any penalties
func (sp *SupplierPayment) GetEffectiveAmount() float64 {
	return sp.InvoiceAmount + sp.PenaltyAmount
}