package interfaces

import (
	"context"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/vehicle_purchase"
)

// VehiclePurchaseTransactionRepository defines the interface for vehicle purchase transaction data operations
type VehiclePurchaseTransactionRepository interface {
	Create(ctx context.Context, transaction *vehicle_purchase.VehiclePurchaseTransaction) (*vehicle_purchase.VehiclePurchaseTransaction, error)
	GetByID(ctx context.Context, id int) (*vehicle_purchase.VehiclePurchaseTransaction, error)
	GetByNumber(ctx context.Context, number string) (*vehicle_purchase.VehiclePurchaseTransaction, error)
	GetByVIN(ctx context.Context, vin string) (*vehicle_purchase.VehiclePurchaseTransaction, error)
	Update(ctx context.Context, id int, transaction *vehicle_purchase.VehiclePurchaseTransaction) (*vehicle_purchase.VehiclePurchaseTransaction, error)
	UpdateStatus(ctx context.Context, id int, status string, updatedBy int) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params *vehicle_purchase.VehiclePurchaseTransactionFilterParams) (*common.PaginatedResponse, error)
	GetByCustomerID(ctx context.Context, customerID int, params *vehicle_purchase.VehiclePurchaseTransactionFilterParams) (*common.PaginatedResponse, error)
	GetByStatus(ctx context.Context, status string, params *vehicle_purchase.VehiclePurchaseTransactionFilterParams) (*common.PaginatedResponse, error)
	GetPendingInspection(ctx context.Context, params *vehicle_purchase.VehiclePurchaseTransactionFilterParams) (*common.PaginatedResponse, error)
	GetPendingApproval(ctx context.Context, params *vehicle_purchase.VehiclePurchaseTransactionFilterParams) (*common.PaginatedResponse, error)
	CompleteInspection(ctx context.Context, id int, request *vehicle_purchase.TransactionInspectionRequest, inspectedBy int) error
	ProcessApproval(ctx context.Context, id int, request *vehicle_purchase.TransactionStatusApprovalRequest, approvedBy int) error
	GenerateNumber(ctx context.Context) (string, error)
	IsNumberExists(ctx context.Context, number string) (bool, error)
	IsVINExists(ctx context.Context, vin string, excludeID int) (bool, error)
	Search(ctx context.Context, query string, params *vehicle_purchase.VehiclePurchaseTransactionFilterParams) (*common.PaginatedResponse, error)
	GetDashboardStats(ctx context.Context) (*vehicle_purchase.TransactionDashboardStats, error)
}

// VehiclePurchasePaymentRepository defines the interface for vehicle purchase payment data operations
type VehiclePurchasePaymentRepository interface {
	Create(ctx context.Context, payment *vehicle_purchase.VehiclePurchasePayment) (*vehicle_purchase.VehiclePurchasePayment, error)
	GetByID(ctx context.Context, id int) (*vehicle_purchase.VehiclePurchasePayment, error)
	GetByNumber(ctx context.Context, number string) (*vehicle_purchase.VehiclePurchasePayment, error)
	Update(ctx context.Context, id int, payment *vehicle_purchase.VehiclePurchasePayment) (*vehicle_purchase.VehiclePurchasePayment, error)
	UpdateStatus(ctx context.Context, id int, status string, updatedBy int) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params *vehicle_purchase.VehiclePurchasePaymentFilterParams) (*common.PaginatedResponse, error)
	GetByTransactionID(ctx context.Context, transactionID int, params *vehicle_purchase.VehiclePurchasePaymentFilterParams) (*common.PaginatedResponse, error)
	GetByStatus(ctx context.Context, status string, params *vehicle_purchase.VehiclePurchasePaymentFilterParams) (*common.PaginatedResponse, error)
	GetPendingApproval(ctx context.Context, params *vehicle_purchase.VehiclePurchasePaymentFilterParams) (*common.PaginatedResponse, error)
	ProcessPayment(ctx context.Context, id int, request *vehicle_purchase.PaymentProcessRequest, processedBy int) error
	ProcessApproval(ctx context.Context, id int, request *vehicle_purchase.PaymentApprovalRequest, approvedBy int) error
	GenerateNumber(ctx context.Context) (string, error)
	IsNumberExists(ctx context.Context, number string) (bool, error)
	GetPaymentSummary(ctx context.Context, transactionID int) (*vehicle_purchase.PaymentSummary, error)
	GetOverduePayments(ctx context.Context, params *vehicle_purchase.VehiclePurchasePaymentFilterParams) (*common.PaginatedResponse, error)
	CalculatePaymentTotals(ctx context.Context, transactionID int) (*vehicle_purchase.PaymentSummary, error)
}