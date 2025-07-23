package transactions

import (
	"time"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

// PaymentMethod represents the payment method used
type PaymentMethod string

const (
	PaymentMethodCash     PaymentMethod = "cash"
	PaymentMethodTransfer PaymentMethod = "transfer"
	PaymentMethodCheck    PaymentMethod = "check"
)

// VehiclePurchaseTransaction represents a transaction for purchasing a vehicle from a customer
type VehiclePurchaseTransaction struct {
	PurchaseID              int               `json:"purchase_id" db:"purchase_id"`
	TransactionNumber       string            `json:"transaction_number" db:"transaction_number"`
	VehicleID               int               `json:"vehicle_id" db:"vehicle_id"`
	CustomerID              int               `json:"customer_id" db:"customer_id"`
	ProcessedBy             int               `json:"processed_by" db:"processed_by"`
	AgreedPrice             float64           `json:"agreed_price" db:"agreed_price"`
	NegotiatedFromPrice     *float64          `json:"negotiated_from_price,omitempty" db:"negotiated_from_price"`
	PurchaseDate            time.Time         `json:"purchase_date" db:"purchase_date"`
	PurchaseType            PurchaseType      `json:"purchase_type" db:"purchase_type"`
	PaymentMethod           PaymentMethod     `json:"payment_method" db:"payment_method"`
	PaymentReference        *string           `json:"payment_reference,omitempty" db:"payment_reference"`
	InitialDamageAssessment *string           `json:"initial_damage_assessment,omitempty" db:"initial_damage_assessment"`
	NegotiationNotes        *string           `json:"negotiation_notes,omitempty" db:"negotiation_notes"`
	Status                  TransactionStatus `json:"status" db:"status"`
	CreatedAt               time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time         `json:"updated_at" db:"updated_at"`
	DeletedAt               *time.Time        `json:"deleted_at,omitempty" db:"deleted_at"`
	CreatedBy               int               `json:"created_by" db:"created_by"`
	Notes                   *string           `json:"notes,omitempty" db:"notes"`

	// Related data
	CustomerName    string `json:"customer_name,omitempty" db:"customer_name"`
	ProcessedByName string `json:"processed_by_name,omitempty" db:"processed_by_name"`
	VehicleCode     string `json:"vehicle_code,omitempty" db:"vehicle_code"`
	ChassisNumber   string `json:"chassis_number,omitempty" db:"chassis_number"`
}

// VehiclePurchaseTransactionCreateRequest represents a request to create a vehicle purchase transaction
type VehiclePurchaseTransactionCreateRequest struct {
	VehicleID               int           `json:"vehicle_id" binding:"required,min=1"`
	CustomerID              int           `json:"customer_id" binding:"required,min=1"`
	AgreedPrice             float64       `json:"agreed_price" binding:"required,min=0"`
	NegotiatedFromPrice     *float64      `json:"negotiated_from_price,omitempty" binding:"omitempty,min=0"`
	PurchaseDate            time.Time     `json:"purchase_date" binding:"required"`
	PurchaseType            PurchaseType  `json:"purchase_type" binding:"required"`
	PaymentMethod           PaymentMethod `json:"payment_method" binding:"required"`
	PaymentReference        *string       `json:"payment_reference,omitempty" binding:"omitempty,max=100"`
	InitialDamageAssessment *string       `json:"initial_damage_assessment,omitempty"`
	NegotiationNotes        *string       `json:"negotiation_notes,omitempty"`
	Notes                   *string       `json:"notes,omitempty"`
}

// VehiclePurchasePayment represents a payment for a vehicle purchase
type VehiclePurchasePayment struct {
	PaymentID        int           `json:"payment_id" db:"payment_id"`
	PurchaseID       int           `json:"purchase_id" db:"purchase_id"`
	PaymentAmount    float64       `json:"payment_amount" db:"payment_amount"`
	PaymentDate      time.Time     `json:"payment_date" db:"payment_date"`
	PaymentMethod    PaymentMethod `json:"payment_method" db:"payment_method"`
	PaymentReference *string       `json:"payment_reference,omitempty" db:"payment_reference"`
	ReceiptNumber    string        `json:"receipt_number" db:"receipt_number"`
	ProcessedBy      int           `json:"processed_by" db:"processed_by"`
	CreatedAt        time.Time     `json:"created_at" db:"created_at"`
	Notes            *string       `json:"notes,omitempty" db:"notes"`

	// Related data
	ProcessedByName string `json:"processed_by_name,omitempty" db:"processed_by_name"`
}

// PaymentType represents the type of payment (cash, credit, etc.)
type PaymentType string

const (
	PaymentTypeCash        PaymentType = "cash"
	PaymentTypeCredit      PaymentType = "credit"
	PaymentTypeTradeIn     PaymentType = "trade_in"
	PaymentTypeCombination PaymentType = "combination"
)

// SalesTransactionStatus represents the status of a sales transaction
type SalesTransactionStatus string

const (
	SalesStatusQuotation      SalesTransactionStatus = "quotation"
	SalesStatusReserved       SalesTransactionStatus = "reserved"
	SalesStatusPendingPayment SalesTransactionStatus = "pending_payment"
	SalesStatusCompleted      SalesTransactionStatus = "completed"
	SalesStatusDelivered      SalesTransactionStatus = "delivered"
	SalesStatusCancelled      SalesTransactionStatus = "cancelled"
)

// WarrantyType represents the type of warranty
type WarrantyType string

const (
	WarrantyTypeEngine        WarrantyType = "engine"
	WarrantyTypeTransmission  WarrantyType = "transmission"
	WarrantyTypeElectrical    WarrantyType = "electrical"
	WarrantyTypeComprehensive WarrantyType = "comprehensive"
	WarrantyTypeLimited       WarrantyType = "limited"
)

// VehicleSalesTransaction represents a transaction for selling a vehicle to a customer
type VehicleSalesTransaction struct {
	SalesID              int                     `json:"sales_id" db:"sales_id"`
	TransactionNumber    string                  `json:"transaction_number" db:"transaction_number"`
	VehicleID            int                     `json:"vehicle_id" db:"vehicle_id"`
	CustomerID           int                     `json:"customer_id" db:"customer_id"`
	SalesPerson          int                     `json:"sales_person" db:"sales_person"`
	Cashier              *int                    `json:"cashier,omitempty" db:"cashier"`
	AskingPrice          float64                 `json:"asking_price" db:"asking_price"`
	NegotiatedPrice      *float64                `json:"negotiated_price,omitempty" db:"negotiated_price"`
	FinalSellingPrice    float64                 `json:"final_selling_price" db:"final_selling_price"`
	DiscountAmount       float64                 `json:"discount_amount" db:"discount_amount"`
	AdditionalFees       float64                 `json:"additional_fees" db:"additional_fees"`
	TotalAmount          float64                 `json:"total_amount" db:"total_amount"`
	PaymentType          PaymentType             `json:"payment_type" db:"payment_type"`
	DownPayment          float64                 `json:"down_payment" db:"down_payment"`
	TradeInValue         *float64                `json:"trade_in_value,omitempty" db:"trade_in_value"`
	TradeInVehicleID     *int                    `json:"trade_in_vehicle_id,omitempty" db:"trade_in_vehicle_id"`
	FinancingAmount      *float64                `json:"financing_amount,omitempty" db:"financing_amount"`
	SaleDate             time.Time               `json:"sale_date" db:"sale_date"`
	DeliveryDate         *time.Time              `json:"delivery_date,omitempty" db:"delivery_date"`
	WarrantyStartDate    *time.Time              `json:"warranty_start_date,omitempty" db:"warranty_start_date"`
	WarrantyEndDate      *time.Time              `json:"warranty_end_date,omitempty" db:"warranty_end_date"`
	WarrantyMonths       *int                    `json:"warranty_months,omitempty" db:"warranty_months"`
	WarrantyType         *WarrantyType           `json:"warranty_type,omitempty" db:"warranty_type"`
	Status               SalesTransactionStatus  `json:"status" db:"status"`
	ContractNumber       *string                 `json:"contract_number,omitempty" db:"contract_number"`
	SalesNotes           *string                 `json:"sales_notes,omitempty" db:"sales_notes"`
	DeliveryNotes        *string                 `json:"delivery_notes,omitempty" db:"delivery_notes"`
	ContractDocumentsJSON *string                `json:"contract_documents_json,omitempty" db:"contract_documents_json"`
	CreatedAt            time.Time               `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time               `json:"updated_at" db:"updated_at"`
	DeletedAt            *time.Time              `json:"deleted_at,omitempty" db:"deleted_at"`
	CreatedBy            int                     `json:"created_by" db:"created_by"`

	// Related data
	CustomerName     string `json:"customer_name,omitempty" db:"customer_name"`
	SalesPersonName  string `json:"sales_person_name,omitempty" db:"sales_person_name"`
	CashierName      string `json:"cashier_name,omitempty" db:"cashier_name"`
	VehicleCode      string `json:"vehicle_code,omitempty" db:"vehicle_code"`
	ChassisNumber    string `json:"chassis_number,omitempty" db:"chassis_number"`
}

// VehicleSalesTransactionCreateRequest represents a request to create a vehicle sales transaction
type VehicleSalesTransactionCreateRequest struct {
	VehicleID            int         `json:"vehicle_id" binding:"required,min=1"`
	CustomerID           int         `json:"customer_id" binding:"required,min=1"`
	AskingPrice          float64     `json:"asking_price" binding:"required,min=0"`
	NegotiatedPrice      *float64    `json:"negotiated_price,omitempty" binding:"omitempty,min=0"`
	FinalSellingPrice    float64     `json:"final_selling_price" binding:"required,min=0"`
	DiscountAmount       float64     `json:"discount_amount" binding:"min=0"`
	AdditionalFees       float64     `json:"additional_fees" binding:"min=0"`
	PaymentType          PaymentType `json:"payment_type" binding:"required"`
	DownPayment          float64     `json:"down_payment" binding:"min=0"`
	TradeInValue         *float64    `json:"trade_in_value,omitempty" binding:"omitempty,min=0"`
	TradeInVehicleID     *int        `json:"trade_in_vehicle_id,omitempty" binding:"omitempty,min=1"`
	FinancingAmount      *float64    `json:"financing_amount,omitempty" binding:"omitempty,min=0"`
	SaleDate             time.Time   `json:"sale_date" binding:"required"`
	DeliveryDate         *time.Time  `json:"delivery_date,omitempty"`
	WarrantyMonths       *int        `json:"warranty_months,omitempty" binding:"omitempty,min=1,max=60"`
	WarrantyType         *WarrantyType `json:"warranty_type,omitempty"`
	SalesNotes           *string     `json:"sales_notes,omitempty"`
	DeliveryNotes        *string     `json:"delivery_notes,omitempty"`
}

// SalesPayment represents a payment for a vehicle sale
type SalesPayment struct {
	PaymentID        int           `json:"payment_id" db:"payment_id"`
	SalesID          int           `json:"sales_id" db:"sales_id"`
	PaymentNumber    string        `json:"payment_number" db:"payment_number"`
	PaymentAmount    float64       `json:"payment_amount" db:"payment_amount"`
	PaymentMethod    PaymentMethod `json:"payment_method" db:"payment_method"`
	PaymentReference *string       `json:"payment_reference,omitempty" db:"payment_reference"`
	PaymentDate      time.Time     `json:"payment_date" db:"payment_date"`
	ReceivedBy       int           `json:"received_by" db:"received_by"`
	ReceiptNumber    string        `json:"receipt_number" db:"receipt_number"`
	Status           TransactionStatus `json:"status" db:"status"`
	PaymentNotes     *string       `json:"payment_notes,omitempty" db:"payment_notes"`
	CreatedAt        time.Time     `json:"created_at" db:"created_at"`

	// Related data
	ReceivedByName string `json:"received_by_name,omitempty" db:"received_by_name"`
}

// TransactionFilterParams represents filtering parameters for transaction queries
type TransactionFilterParams struct {
	CustomerID    *int                    `json:"customer_id,omitempty" form:"customer_id"`
	VehicleID     *int                    `json:"vehicle_id,omitempty" form:"vehicle_id"`
	Status        *TransactionStatus      `json:"status,omitempty" form:"status"`
	SalesStatus   *SalesTransactionStatus `json:"sales_status,omitempty" form:"sales_status"`
	PaymentMethod *PaymentMethod          `json:"payment_method,omitempty" form:"payment_method"`
	PaymentType   *PaymentType            `json:"payment_type,omitempty" form:"payment_type"`
	DateFrom      *time.Time              `json:"date_from,omitempty" form:"date_from"`
	DateTo        *time.Time              `json:"date_to,omitempty" form:"date_to"`
	MinAmount     *float64                `json:"min_amount,omitempty" form:"min_amount"`
	MaxAmount     *float64                `json:"max_amount,omitempty" form:"max_amount"`
	Search        string                  `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// IsValid validates the transaction status
func (s TransactionStatus) IsValid() bool {
	switch s {
	case TransactionStatusPending, TransactionStatusCompleted, TransactionStatusCancelled:
		return true
	default:
		return false
	}
}

// IsValid validates the payment method
func (p PaymentMethod) IsValid() bool {
	switch p {
	case PaymentMethodCash, PaymentMethodTransfer, PaymentMethodCheck:
		return true
	default:
		return false
	}
}

// IsValid validates the payment type
func (p PaymentType) IsValid() bool {
	switch p {
	case PaymentTypeCash, PaymentTypeCredit, PaymentTypeTradeIn, PaymentTypeCombination:
		return true
	default:
		return false
	}
}

// IsValid validates the sales transaction status
func (s SalesTransactionStatus) IsValid() bool {
	switch s {
	case SalesStatusQuotation, SalesStatusReserved, SalesStatusPendingPayment,
		 SalesStatusCompleted, SalesStatusDelivered, SalesStatusCancelled:
		return true
	default:
		return false
	}
}

// IsValid validates the warranty type
func (w WarrantyType) IsValid() bool {
	switch w {
	case WarrantyTypeEngine, WarrantyTypeTransmission, WarrantyTypeElectrical,
		 WarrantyTypeComprehensive, WarrantyTypeLimited:
		return true
	default:
		return false
	}
}

// CalculateTotalAmount calculates the total amount for a sales transaction
func (v *VehicleSalesTransaction) CalculateTotalAmount() {
	v.TotalAmount = v.FinalSellingPrice - v.DiscountAmount + v.AdditionalFees
	if v.TradeInValue != nil {
		v.TotalAmount -= *v.TradeInValue
	}
}

// IsCompleted checks if the transaction is completed
func (v *VehiclePurchaseTransaction) IsCompleted() bool {
	return v.Status == TransactionStatusCompleted
}

// CanBeCancelled checks if the transaction can be cancelled
func (v *VehiclePurchaseTransaction) CanBeCancelled() bool {
	return v.Status == TransactionStatusPending
}

// IsCompleted checks if the sales transaction is completed
func (v *VehicleSalesTransaction) IsCompleted() bool {
	return v.Status == SalesStatusCompleted || v.Status == SalesStatusDelivered
}

// RequiresFinancing checks if the sales transaction requires financing
func (v *VehicleSalesTransaction) RequiresFinancing() bool {
	return v.PaymentType == PaymentTypeCredit && v.FinancingAmount != nil && *v.FinancingAmount > 0
}