package vehicle_purchase

import (
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// VehiclePurchaseTransaction represents a vehicle purchase transaction from a customer
type VehiclePurchaseTransaction struct {
	TransactionID      int       `json:"transaction_id" db:"transaction_id"`
	TransactionNumber  string    `json:"transaction_number" db:"transaction_number"`
	CustomerID         int       `json:"customer_id" db:"customer_id"`
	VehicleID          *int      `json:"vehicle_id,omitempty" db:"vehicle_id"`
	VinNumber          *string   `json:"vin_number,omitempty" db:"vin_number"`
	VehicleBrand       string    `json:"vehicle_brand" db:"vehicle_brand"`
	VehicleModel       string    `json:"vehicle_model" db:"vehicle_model"`
	VehicleYear        int       `json:"vehicle_year" db:"vehicle_year"`
	VehicleColor       string    `json:"vehicle_color" db:"vehicle_color"`
	EngineNumber       *string   `json:"engine_number,omitempty" db:"engine_number"`
	RegistrationNumber *string   `json:"registration_number,omitempty" db:"registration_number"`
	PurchasePrice      float64   `json:"purchase_price" db:"purchase_price"`
	AgreedValue        float64   `json:"agreed_value" db:"agreed_value"`
	OdometerReading    int       `json:"odometer_reading" db:"odometer_reading"`
	FuelType           string    `json:"fuel_type" db:"fuel_type"`
	Transmission       string    `json:"transmission" db:"transmission"`
	ConditionRating    *int      `json:"condition_rating,omitempty" db:"condition_rating"`
	PurchaseDate       time.Time `json:"purchase_date" db:"purchase_date"`
	TransactionStatus  string    `json:"transaction_status" db:"transaction_status"`
	InspectionNotes    *string   `json:"inspection_notes,omitempty" db:"inspection_notes"`
	EvaluationNotes    *string   `json:"evaluation_notes,omitempty" db:"evaluation_notes"`
	PurchaseNotes      *string   `json:"purchase_notes,omitempty" db:"purchase_notes"`
	DocumentsJSON      *string   `json:"documents_json,omitempty" db:"documents_json"`
	ProcessedBy        int       `json:"processed_by" db:"processed_by"`
	InspectedBy        *int      `json:"inspected_by,omitempty" db:"inspected_by"`
	ApprovedBy         *int      `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt         *time.Time `json:"approved_at,omitempty" db:"approved_at"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`

	// Related data for joins
	CustomerName       string `json:"customer_name,omitempty" db:"customer_name"`
	ProcessedByName    string `json:"processed_by_name,omitempty" db:"processed_by_name"`
	InspectedByName    string `json:"inspected_by_name,omitempty" db:"inspected_by_name"`
	ApprovedByName     string `json:"approved_by_name,omitempty" db:"approved_by_name"`
}

// VehiclePurchaseTransactionListItem represents a simplified transaction for list views
type VehiclePurchaseTransactionListItem struct {
	TransactionID     int       `json:"transaction_id" db:"transaction_id"`
	TransactionNumber string    `json:"transaction_number" db:"transaction_number"`
	CustomerName      string    `json:"customer_name" db:"customer_name"`
	VehicleBrand      string    `json:"vehicle_brand" db:"vehicle_brand"`
	VehicleModel      string    `json:"vehicle_model" db:"vehicle_model"`
	VehicleYear       int       `json:"vehicle_year" db:"vehicle_year"`
	PurchasePrice     float64   `json:"purchase_price" db:"purchase_price"`
	TransactionStatus string    `json:"transaction_status" db:"transaction_status"`
	PurchaseDate      time.Time `json:"purchase_date" db:"purchase_date"`
	ProcessedByName   string    `json:"processed_by_name" db:"processed_by_name"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

// VehiclePurchaseTransactionCreateRequest represents a request to create a vehicle purchase transaction
type VehiclePurchaseTransactionCreateRequest struct {
	CustomerID         int      `json:"customer_id" binding:"required"`
	VehicleID          *int     `json:"vehicle_id,omitempty"`
	VinNumber          *string  `json:"vin_number,omitempty" binding:"omitempty,max=50"`
	VehicleBrand       string   `json:"vehicle_brand" binding:"required,max=100"`
	VehicleModel       string   `json:"vehicle_model" binding:"required,max=100"`
	VehicleYear        int      `json:"vehicle_year" binding:"required,min=1900,max=2100"`
	VehicleColor       string   `json:"vehicle_color" binding:"required,max=50"`
	EngineNumber       *string  `json:"engine_number,omitempty" binding:"omitempty,max=100"`
	RegistrationNumber *string  `json:"registration_number,omitempty" binding:"omitempty,max=50"`
	PurchasePrice      float64  `json:"purchase_price" binding:"required,min=0"`
	AgreedValue        float64  `json:"agreed_value" binding:"required,min=0"`
	OdometerReading    int      `json:"odometer_reading" binding:"min=0"`
	FuelType           string   `json:"fuel_type" binding:"required,max=50"`
	Transmission       string   `json:"transmission" binding:"required,max=50"`
	ConditionRating    *int     `json:"condition_rating,omitempty" binding:"omitempty,min=1,max=10"`
	InspectionNotes    *string  `json:"inspection_notes,omitempty"`
	EvaluationNotes    *string  `json:"evaluation_notes,omitempty"`
	PurchaseNotes      *string  `json:"purchase_notes,omitempty"`
	DocumentsJSON      *string  `json:"documents_json,omitempty"`
}

// VehiclePurchaseTransactionUpdateRequest represents a request to update a vehicle purchase transaction
type VehiclePurchaseTransactionUpdateRequest struct {
	VinNumber          *string  `json:"vin_number,omitempty" binding:"omitempty,max=50"`
	VehicleColor       *string  `json:"vehicle_color,omitempty" binding:"omitempty,max=50"`
	EngineNumber       *string  `json:"engine_number,omitempty" binding:"omitempty,max=100"`
	RegistrationNumber *string  `json:"registration_number,omitempty" binding:"omitempty,max=50"`
	PurchasePrice      *float64 `json:"purchase_price,omitempty" binding:"omitempty,min=0"`
	AgreedValue        *float64 `json:"agreed_value,omitempty" binding:"omitempty,min=0"`
	OdometerReading    *int     `json:"odometer_reading,omitempty" binding:"omitempty,min=0"`
	ConditionRating    *int     `json:"condition_rating,omitempty" binding:"omitempty,min=1,max=10"`
	TransactionStatus  *string  `json:"transaction_status,omitempty" binding:"omitempty,oneof=pending inspection approved rejected completed cancelled"`
	InspectionNotes    *string  `json:"inspection_notes,omitempty"`
	EvaluationNotes    *string  `json:"evaluation_notes,omitempty"`
	PurchaseNotes      *string  `json:"purchase_notes,omitempty"`
	DocumentsJSON      *string  `json:"documents_json,omitempty"`
}

// VehiclePurchaseTransactionFilterParams represents filtering parameters for transaction queries
type VehiclePurchaseTransactionFilterParams struct {
	CustomerID        *int      `json:"customer_id,omitempty" form:"customer_id"`
	VehicleBrand      string    `json:"vehicle_brand,omitempty" form:"vehicle_brand"`
	VehicleModel      string    `json:"vehicle_model,omitempty" form:"vehicle_model"`
	VehicleYear       *int      `json:"vehicle_year,omitempty" form:"vehicle_year"`
	TransactionStatus string    `json:"transaction_status,omitempty" form:"transaction_status"`
	MinPurchasePrice  *float64  `json:"min_purchase_price,omitempty" form:"min_purchase_price"`
	MaxPurchasePrice  *float64  `json:"max_purchase_price,omitempty" form:"max_purchase_price"`
	ProcessedBy       *int      `json:"processed_by,omitempty" form:"processed_by"`
	InspectedBy       *int      `json:"inspected_by,omitempty" form:"inspected_by"`
	StartDate         *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate           *time.Time `json:"end_date,omitempty" form:"end_date"`
	Search            string    `json:"search,omitempty" form:"search"`
	common.PaginationParams
}

// TransactionStatusApprovalRequest represents a request to approve/reject a transaction
type TransactionStatusApprovalRequest struct {
	Status          string  `json:"status" binding:"required,oneof=approved rejected"`
	ApprovalNotes   *string `json:"approval_notes,omitempty"`
	ConditionRating *int    `json:"condition_rating,omitempty" binding:"omitempty,min=1,max=10"`
}

// TransactionInspectionRequest represents a request to complete inspection
type TransactionInspectionRequest struct {
	ConditionRating   int     `json:"condition_rating" binding:"required,min=1,max=10"`
	InspectionNotes   string  `json:"inspection_notes" binding:"required"`
	EvaluationNotes   *string `json:"evaluation_notes,omitempty"`
	RecommendedAction string  `json:"recommended_action" binding:"required,oneof=approve reject needs_repair"`
}

// TransactionDashboardStats represents dashboard statistics for transactions
type TransactionDashboardStats struct {
	TotalTransactions     int     `json:"total_transactions" db:"total_transactions"`
	PendingTransactions   int     `json:"pending_transactions" db:"pending_transactions"`
	CompletedTransactions int     `json:"completed_transactions" db:"completed_transactions"`
	TotalPurchaseValue    float64 `json:"total_purchase_value" db:"total_purchase_value"`
	MonthlyTransactions   int     `json:"monthly_transactions" db:"monthly_transactions"`
	MonthlyValue          float64 `json:"monthly_value" db:"monthly_value"`
	AvgTransactionValue   float64 `json:"avg_transaction_value" db:"avg_transaction_value"`
	PendingInspections    int     `json:"pending_inspections" db:"pending_inspections"`
	PendingApprovals      int     `json:"pending_approvals" db:"pending_approvals"`
}