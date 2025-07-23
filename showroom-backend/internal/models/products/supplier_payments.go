package products

import (
	"database/sql/driver"
	"fmt"
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

// IsValid checks if the payment method is valid
func (m PaymentMethod) IsValid() bool {
	switch m {
	case PaymentMethodCash, PaymentMethodTransfer, PaymentMethodCheck, PaymentMethodCredit:
		return true
	default:
		return false
	}
}

// String returns the string representation of the payment method
func (m PaymentMethod) String() string {
	return string(m)
}

// Value implements the driver.Valuer interface for PaymentMethod
func (m PaymentMethod) Value() (driver.Value, error) {
	return string(m), nil
}

// Scan implements the sql.Scanner interface for PaymentMethod
func (m *PaymentMethod) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch s := value.(type) {
	case string:
		*m = PaymentMethod(s)
	case []byte:
		*m = PaymentMethod(s)
	default:
		return fmt.Errorf("cannot scan %T into PaymentMethod", value)
	}
	return nil
}

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPartial  PaymentStatus = "partial"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusOverdue  PaymentStatus = "overdue"
	PaymentStatusDisputed PaymentStatus = "disputed"
)

// IsValid checks if the payment status is valid
func (s PaymentStatus) IsValid() bool {
	switch s {
	case PaymentStatusPending, PaymentStatusPartial, PaymentStatusPaid, PaymentStatusOverdue, PaymentStatusDisputed:
		return true
	default:
		return false
	}
}

// String returns the string representation of the payment status
func (s PaymentStatus) String() string {
	return string(s)
}

// Value implements the driver.Valuer interface for PaymentStatus
func (s PaymentStatus) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan implements the sql.Scanner interface for PaymentStatus
func (s *PaymentStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch str := value.(type) {
	case string:
		*s = PaymentStatus(str)
	case []byte:
		*s = PaymentStatus(str)
	default:
		return fmt.Errorf("cannot scan %T into PaymentStatus", value)
	}
	return nil
}

// SupplierPayment represents a payment to a supplier
type SupplierPayment struct {
	PaymentID         int           `json:"payment_id" db:"payment_id"`
	SupplierID        int           `json:"supplier_id" db:"supplier_id"`
	POID              *int          `json:"po_id,omitempty" db:"po_id"`
	PaymentNumber     string        `json:"payment_number" db:"payment_number"`
	InvoiceAmount     float64       `json:"invoice_amount" db:"invoice_amount"`
	PaymentAmount     float64       `json:"payment_amount" db:"payment_amount"`
	DiscountTaken     float64       `json:"discount_taken" db:"discount_taken"`
	OutstandingAmount float64       `json:"outstanding_amount" db:"outstanding_amount"`
	InvoiceDate       time.Time     `json:"invoice_date" db:"invoice_date"`
	PaymentDate       time.Time     `json:"payment_date" db:"payment_date"`
	DueDate           time.Time     `json:"due_date" db:"due_date"`
	PaymentMethod     PaymentMethod `json:"payment_method" db:"payment_method"`
	PaymentReference  *string       `json:"payment_reference,omitempty" db:"payment_reference"`
	InvoiceNumber     string        `json:"invoice_number" db:"invoice_number"`
	PaymentStatus     PaymentStatus `json:"payment_status" db:"payment_status"`
	DaysOverdue       int           `json:"days_overdue" db:"days_overdue"`
	PenaltyAmount     float64       `json:"penalty_amount" db:"penalty_amount"`
	ProcessedBy       int           `json:"processed_by" db:"processed_by"`
	PaymentNotes      *string       `json:"payment_notes,omitempty" db:"payment_notes"`
	CreatedAt         time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at" db:"updated_at"`
}

// SupplierPaymentListItem represents a simplified supplier payment for list views
type SupplierPaymentListItem struct {
	PaymentID         int           `json:"payment_id" db:"payment_id"`
	SupplierID        int           `json:"supplier_id" db:"supplier_id"`
	POID              *int          `json:"po_id,omitempty" db:"po_id"`
	PaymentNumber     string        `json:"payment_number" db:"payment_number"`
	InvoiceAmount     float64       `json:"invoice_amount" db:"invoice_amount"`
	PaymentAmount     float64       `json:"payment_amount" db:"payment_amount"`
	OutstandingAmount float64       `json:"outstanding_amount" db:"outstanding_amount"`
	InvoiceDate       time.Time     `json:"invoice_date" db:"invoice_date"`
	PaymentDate       time.Time     `json:"payment_date" db:"payment_date"`
	DueDate           time.Time     `json:"due_date" db:"due_date"`
	PaymentMethod     PaymentMethod `json:"payment_method" db:"payment_method"`
	InvoiceNumber     string        `json:"invoice_number" db:"invoice_number"`
	PaymentStatus     PaymentStatus `json:"payment_status" db:"payment_status"`
	DaysOverdue       int           `json:"days_overdue" db:"days_overdue"`
}

// SupplierPaymentCreateRequest represents a request to create a supplier payment
type SupplierPaymentCreateRequest struct {
	SupplierID       int           `json:"supplier_id" binding:"required,min=1"`
	POID             *int          `json:"po_id,omitempty" binding:"omitempty,min=1"`
	InvoiceAmount    float64       `json:"invoice_amount" binding:"required,min=0"`
	PaymentAmount    float64       `json:"payment_amount" binding:"required,min=0"`
	DiscountTaken    float64       `json:"discount_taken" binding:"min=0"`
	InvoiceDate      time.Time     `json:"invoice_date" binding:"required"`
	PaymentDate      time.Time     `json:"payment_date" binding:"required"`
	DueDate          time.Time     `json:"due_date" binding:"required"`
	PaymentMethod    PaymentMethod `json:"payment_method" binding:"required"`
	PaymentReference *string       `json:"payment_reference,omitempty" binding:"omitempty,max=100"`
	InvoiceNumber    string        `json:"invoice_number" binding:"required,max=100"`
	PaymentNotes     *string       `json:"payment_notes,omitempty"`
}

// SupplierPaymentUpdateRequest represents a request to update a supplier payment
type SupplierPaymentUpdateRequest struct {
	PaymentAmount     *float64       `json:"payment_amount,omitempty" binding:"omitempty,min=0"`
	DiscountTaken     *float64       `json:"discount_taken,omitempty" binding:"omitempty,min=0"`
	PaymentDate       *time.Time     `json:"payment_date,omitempty"`
	DueDate           *time.Time     `json:"due_date,omitempty"`
	PaymentMethod     *PaymentMethod `json:"payment_method,omitempty"`
	PaymentReference  *string        `json:"payment_reference,omitempty" binding:"omitempty,max=100"`
	PaymentStatus     *PaymentStatus `json:"payment_status,omitempty"`
	PenaltyAmount     *float64       `json:"penalty_amount,omitempty" binding:"omitempty,min=0"`
	PaymentNotes      *string        `json:"payment_notes,omitempty"`
}

// SupplierPaymentFilterParams represents filtering parameters for supplier payment queries
type SupplierPaymentFilterParams struct {
	SupplierID     *int           `json:"supplier_id,omitempty" form:"supplier_id"`
	POID           *int           `json:"po_id,omitempty" form:"po_id"`
	PaymentStatus  *PaymentStatus `json:"payment_status,omitempty" form:"payment_status"`
	PaymentMethod  *PaymentMethod `json:"payment_method,omitempty" form:"payment_method"`
	ProcessedBy    *int           `json:"processed_by,omitempty" form:"processed_by"`
	DateFrom       *time.Time     `json:"date_from,omitempty" form:"date_from"`
	DateTo         *time.Time     `json:"date_to,omitempty" form:"date_to"`
	IsOverdue      *bool          `json:"is_overdue,omitempty" form:"is_overdue"`
	MinAmount      *float64       `json:"min_amount,omitempty" form:"min_amount"`
	MaxAmount      *float64       `json:"max_amount,omitempty" form:"max_amount"`
	Search         string         `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// CalculateOutstandingAmount calculates the outstanding amount
func (sp *SupplierPayment) CalculateOutstandingAmount() {
	sp.OutstandingAmount = sp.InvoiceAmount - sp.PaymentAmount - sp.DiscountTaken
}

// UpdateDaysOverdue calculates and updates days overdue
func (sp *SupplierPayment) UpdateDaysOverdue() {
	if sp.PaymentStatus == PaymentStatusPaid {
		sp.DaysOverdue = 0
		return
	}
	
	now := time.Now()
	if now.After(sp.DueDate) {
		sp.DaysOverdue = int(now.Sub(sp.DueDate).Hours() / 24)
	} else {
		sp.DaysOverdue = 0
	}
}

// UpdatePaymentStatus updates the payment status based on amounts and dates
func (sp *SupplierPayment) UpdatePaymentStatus() {
	sp.CalculateOutstandingAmount()
	sp.UpdateDaysOverdue()
	
	// If disputed, don't change status automatically
	if sp.PaymentStatus == PaymentStatusDisputed {
		return
	}
	
	if sp.OutstandingAmount <= 0 {
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

// IsFullyPaid checks if the payment is fully paid
func (sp *SupplierPayment) IsFullyPaid() bool {
	return sp.PaymentStatus == PaymentStatusPaid || sp.OutstandingAmount <= 0
}

// IsOverdue checks if the payment is overdue
func (sp *SupplierPayment) IsOverdue() bool {
	return sp.PaymentStatus == PaymentStatusOverdue || sp.DaysOverdue > 0
}

// CanAddPayment checks if additional payment can be added
func (sp *SupplierPayment) CanAddPayment() bool {
	return !sp.IsFullyPaid() && sp.PaymentStatus != PaymentStatusDisputed
}

// AddPayment adds a payment amount and updates status
func (sp *SupplierPayment) AddPayment(amount float64, method PaymentMethod, reference *string) {
	if !sp.CanAddPayment() || amount <= 0 {
		return
	}
	
	// Ensure we don't overpay
	maxPayable := sp.InvoiceAmount - sp.PaymentAmount - sp.DiscountTaken
	if amount > maxPayable {
		amount = maxPayable
	}
	
	sp.PaymentAmount += amount
	sp.PaymentMethod = method
	sp.PaymentReference = reference
	sp.PaymentDate = time.Now()
	sp.UpdatePaymentStatus()
}

// CalculatePenalty calculates penalty based on days overdue and rate
func (sp *SupplierPayment) CalculatePenalty(dailyPenaltyRate float64) {
	if sp.DaysOverdue <= 0 || sp.IsFullyPaid() {
		sp.PenaltyAmount = 0
		return
	}
	
	sp.PenaltyAmount = sp.OutstandingAmount * (dailyPenaltyRate / 100) * float64(sp.DaysOverdue)
}

// GetPaymentProgressPercentage returns the payment progress as percentage
func (sp *SupplierPayment) GetPaymentProgressPercentage() float64 {
	if sp.InvoiceAmount == 0 {
		return 0
	}
	return ((sp.PaymentAmount + sp.DiscountTaken) / sp.InvoiceAmount) * 100
}

// GetPaymentDescription returns a human-readable description of the payment
func (sp *SupplierPayment) GetPaymentDescription() string {
	switch sp.PaymentStatus {
	case PaymentStatusPending:
		return "Payment pending"
	case PaymentStatusPartial:
		percentage := sp.GetPaymentProgressPercentage()
		return fmt.Sprintf("Partially paid (%.1f%%)", percentage)
	case PaymentStatusPaid:
		return "Fully paid"
	case PaymentStatusOverdue:
		return fmt.Sprintf("Overdue by %d days", sp.DaysOverdue)
	case PaymentStatusDisputed:
		return "Payment disputed"
	default:
		return string(sp.PaymentStatus)
	}
}

// GetRemainingAmount returns the amount still to be paid
func (sp *SupplierPayment) GetRemainingAmount() float64 {
	remaining := sp.InvoiceAmount - sp.PaymentAmount - sp.DiscountTaken
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetTotalAmountDue returns total amount due including penalties
func (sp *SupplierPayment) GetTotalAmountDue() float64 {
	return sp.GetRemainingAmount() + sp.PenaltyAmount
}