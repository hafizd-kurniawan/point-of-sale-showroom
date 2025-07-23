package vehicle_purchase

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// VehiclePurchasePayment represents a payment for a vehicle purchase transaction
type VehiclePurchasePayment struct {
	PaymentID          int        `json:"payment_id" db:"payment_id"`
	TransactionID      int        `json:"transaction_id" db:"transaction_id"`
	PaymentNumber      string     `json:"payment_number" db:"payment_number"`
	PaymentMethod      string     `json:"payment_method" db:"payment_method"`
	PaymentAmount      float64    `json:"payment_amount" db:"payment_amount"`
	PaymentDate        time.Time  `json:"payment_date" db:"payment_date"`
	PaymentStatus      string     `json:"payment_status" db:"payment_status"`
	ReferenceNumber    *string    `json:"reference_number,omitempty" db:"reference_number"`
	BankAccount        *string    `json:"bank_account,omitempty" db:"bank_account"`
	PaymentDescription *string    `json:"payment_description,omitempty" db:"payment_description"`
	PaymentNotes       *string    `json:"payment_notes,omitempty" db:"payment_notes"`
	ProcessedBy        int        `json:"processed_by" db:"processed_by"`
	ApprovedBy         *int       `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt         *time.Time `json:"approved_at,omitempty" db:"approved_at"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`

	// Related data for joins
	TransactionNumber   string  `json:"transaction_number,omitempty" db:"transaction_number"`
	CustomerName        string  `json:"customer_name,omitempty" db:"customer_name"`
	VehicleBrand        string  `json:"vehicle_brand,omitempty" db:"vehicle_brand"`
	VehicleModel        string  `json:"vehicle_model,omitempty" db:"vehicle_model"`
	VehicleYear         int     `json:"vehicle_year,omitempty" db:"vehicle_year"`
	ProcessedByName     string  `json:"processed_by_name,omitempty" db:"processed_by_name"`
	ApprovedByName      string  `json:"approved_by_name,omitempty" db:"approved_by_name"`
}

// VehiclePurchasePaymentListItem represents a simplified payment for list views
type VehiclePurchasePaymentListItem struct {
	PaymentID         int       `json:"payment_id" db:"payment_id"`
	PaymentNumber     string    `json:"payment_number" db:"payment_number"`
	TransactionNumber string    `json:"transaction_number" db:"transaction_number"`
	CustomerName      string    `json:"customer_name" db:"customer_name"`
	VehicleBrand      string    `json:"vehicle_brand" db:"vehicle_brand"`
	VehicleModel      string    `json:"vehicle_model" db:"vehicle_model"`
	PaymentMethod     string    `json:"payment_method" db:"payment_method"`
	PaymentAmount     float64   `json:"payment_amount" db:"payment_amount"`
	PaymentStatus     string    `json:"payment_status" db:"payment_status"`
	PaymentDate       time.Time `json:"payment_date" db:"payment_date"`
	ProcessedByName   string    `json:"processed_by_name" db:"processed_by_name"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

// VehiclePurchasePaymentCreateRequest represents a request to create a vehicle purchase payment
type VehiclePurchasePaymentCreateRequest struct {
	TransactionID      int     `json:"transaction_id" binding:"required"`
	PaymentMethod      string  `json:"payment_method" binding:"required,oneof=cash transfer check financing"`
	PaymentAmount      float64 `json:"payment_amount" binding:"required,min=0"`
	ReferenceNumber    *string `json:"reference_number,omitempty" binding:"omitempty,max=100"`
	BankAccount        *string `json:"bank_account,omitempty" binding:"omitempty,max=100"`
	PaymentDescription *string `json:"payment_description,omitempty"`
	PaymentNotes       *string `json:"payment_notes,omitempty"`
}

// VehiclePurchasePaymentUpdateRequest represents a request to update a vehicle purchase payment
type VehiclePurchasePaymentUpdateRequest struct {
	PaymentMethod      *string  `json:"payment_method,omitempty" binding:"omitempty,oneof=cash transfer check financing"`
	PaymentAmount      *float64 `json:"payment_amount,omitempty" binding:"omitempty,min=0"`
	PaymentStatus      *string  `json:"payment_status,omitempty" binding:"omitempty,oneof=pending processing completed failed cancelled"`
	ReferenceNumber    *string  `json:"reference_number,omitempty" binding:"omitempty,max=100"`
	BankAccount        *string  `json:"bank_account,omitempty" binding:"omitempty,max=100"`
	PaymentDescription *string  `json:"payment_description,omitempty"`
	PaymentNotes       *string  `json:"payment_notes,omitempty"`
}

// VehiclePurchasePaymentFilterParams represents filtering parameters for payment queries
type VehiclePurchasePaymentFilterParams struct {
	TransactionID     *int       `json:"transaction_id,omitempty" form:"transaction_id"`
	PaymentMethod     string     `json:"payment_method,omitempty" form:"payment_method"`
	PaymentStatus     string     `json:"payment_status,omitempty" form:"payment_status"`
	MinPaymentAmount  *float64   `json:"min_payment_amount,omitempty" form:"min_payment_amount"`
	MaxPaymentAmount  *float64   `json:"max_payment_amount,omitempty" form:"max_payment_amount"`
	ProcessedBy       *int       `json:"processed_by,omitempty" form:"processed_by"`
	StartDate         *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate           *time.Time `json:"end_date,omitempty" form:"end_date"`
	Search            string     `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// PaymentApprovalRequest represents a request to approve/reject a payment
type PaymentApprovalRequest struct {
	Status        string  `json:"status" binding:"required,oneof=completed failed"`
	ApprovalNotes *string `json:"approval_notes,omitempty"`
}

// PaymentProcessRequest represents a request to process a payment
type PaymentProcessRequest struct {
	Status          string  `json:"status" binding:"required,oneof=processing completed failed"`
	ReferenceNumber *string `json:"reference_number,omitempty" binding:"omitempty,max=100"`
	ProcessingNotes *string `json:"processing_notes,omitempty"`
}

// PaymentSummary represents a summary of payments for a transaction
type PaymentSummary struct {
	TransactionID    int     `json:"transaction_id" db:"transaction_id"`
	TotalAmount      float64 `json:"total_amount" db:"total_amount"`
	PaidAmount       float64 `json:"paid_amount" db:"paid_amount"`
	PendingAmount    float64 `json:"pending_amount" db:"pending_amount"`
	PaymentCount     int     `json:"payment_count" db:"payment_count"`
	IsFullyPaid      bool    `json:"is_fully_paid" db:"is_fully_paid"`
	LastPaymentDate  *time.Time `json:"last_payment_date,omitempty" db:"last_payment_date"`
}