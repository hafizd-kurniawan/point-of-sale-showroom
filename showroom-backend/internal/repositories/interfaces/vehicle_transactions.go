package interfaces

import (
	"context"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/transactions"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
)

// VehicleInventoryRepository defines the interface for vehicle inventory repository
type VehicleInventoryRepository interface {
	// Create adds a new vehicle to inventory
	Create(ctx context.Context, vehicle *transactions.VehicleInventory) (*transactions.VehicleInventory, error)
	
	// GetByID retrieves a vehicle by its ID
	GetByID(ctx context.Context, vehicleID int) (*transactions.VehicleInventory, error)
	
	// GetByVehicleCode retrieves a vehicle by its vehicle code
	GetByVehicleCode(ctx context.Context, vehicleCode string) (*transactions.VehicleInventory, error)
	
	// GetByChassisNumber retrieves a vehicle by its chassis number
	GetByChassisNumber(ctx context.Context, chassisNumber string) (*transactions.VehicleInventory, error)
	
	// List retrieves vehicles with filtering and pagination
	List(ctx context.Context, filter *transactions.VehicleInventoryFilterParams) ([]*transactions.VehicleInventoryListItem, *common.PaginationMeta, error)
	
	// Update updates a vehicle's information
	Update(ctx context.Context, vehicleID int, updateReq *transactions.VehicleInventoryUpdateRequest) (*transactions.VehicleInventory, error)
	
	// UpdateStatus updates a vehicle's status
	UpdateStatus(ctx context.Context, vehicleID int, status transactions.VehicleStatus, updatedBy int) error
	
	// SoftDelete soft deletes a vehicle (if allowed)
	SoftDelete(ctx context.Context, vehicleID int, deletedBy int) error
	
	// GetAvailableForSale retrieves vehicles available for sale
	GetAvailableForSale(ctx context.Context, filter *transactions.VehicleInventoryFilterParams) ([]*transactions.VehicleInventoryListItem, *common.PaginationMeta, error)
	
	// GetNeedingRepair retrieves vehicles that need repair
	GetNeedingRepair(ctx context.Context, filter *transactions.VehicleInventoryFilterParams) ([]*transactions.VehicleInventoryListItem, *common.PaginationMeta, error)
	
	// GenerateVehicleCode generates a unique vehicle code
	GenerateVehicleCode(ctx context.Context) (string, error)
	
	// CheckChassisNumberExists checks if chassis number already exists
	CheckChassisNumberExists(ctx context.Context, chassisNumber string, excludeVehicleID ...int) (bool, error)
	
	// CheckEngineNumberExists checks if engine number already exists
	CheckEngineNumberExists(ctx context.Context, engineNumber string, excludeVehicleID ...int) (bool, error)
}

// VehicleTransactionRepository defines the interface for vehicle transaction repository
type VehicleTransactionRepository interface {
	// Purchase transactions
	CreatePurchaseTransaction(ctx context.Context, purchase *transactions.VehiclePurchaseTransaction) (*transactions.VehiclePurchaseTransaction, error)
	GetPurchaseTransaction(ctx context.Context, purchaseID int) (*transactions.VehiclePurchaseTransaction, error)
	ListPurchaseTransactions(ctx context.Context, filter *transactions.TransactionFilterParams) ([]*transactions.VehiclePurchaseTransaction, *common.PaginationMeta, error)
	UpdatePurchaseTransactionStatus(ctx context.Context, purchaseID int, status transactions.TransactionStatus, updatedBy int) error
	
	// Purchase payments
	CreatePurchasePayment(ctx context.Context, payment *transactions.VehiclePurchasePayment) (*transactions.VehiclePurchasePayment, error)
	GetPurchasePaymentsByTransactionID(ctx context.Context, purchaseID int) ([]*transactions.VehiclePurchasePayment, error)
	
	// Sales transactions
	CreateSalesTransaction(ctx context.Context, sales *transactions.VehicleSalesTransaction) (*transactions.VehicleSalesTransaction, error)
	GetSalesTransaction(ctx context.Context, salesID int) (*transactions.VehicleSalesTransaction, error)
	ListSalesTransactions(ctx context.Context, filter *transactions.TransactionFilterParams) ([]*transactions.VehicleSalesTransaction, *common.PaginationMeta, error)
	UpdateSalesTransactionStatus(ctx context.Context, salesID int, status transactions.SalesTransactionStatus, updatedBy int) error
	CompleteSalesTransaction(ctx context.Context, salesID int, completedBy int) error
	
	// Sales payments
	CreateSalesPayment(ctx context.Context, payment *transactions.SalesPayment) (*transactions.SalesPayment, error)
	GetSalesPaymentsByTransactionID(ctx context.Context, salesID int) ([]*transactions.SalesPayment, error)
	
	// Transaction number generation
	GeneratePurchaseTransactionNumber(ctx context.Context) (string, error)
	GenerateSalesTransactionNumber(ctx context.Context) (string, error)
	GeneratePaymentNumber(ctx context.Context) (string, error)
	GenerateReceiptNumber(ctx context.Context) (string, error)
}

// FinancialTransactionRepository defines the interface for financial transaction repository
type FinancialTransactionRepository interface {
	// Financial transactions
	Create(ctx context.Context, transaction *transactions.FinancialTransaction) (*transactions.FinancialTransaction, error)
	GetByID(ctx context.Context, transactionID int) (*transactions.FinancialTransaction, error)
	List(ctx context.Context, filter *transactions.FinancialFilterParams) ([]*transactions.FinancialTransaction, *common.PaginationMeta, error)
	Update(ctx context.Context, transactionID int, transaction *transactions.FinancialTransaction) (*transactions.FinancialTransaction, error)
	UpdateStatus(ctx context.Context, transactionID int, status transactions.FinancialTransactionStatus, updatedBy int) error
	SoftDelete(ctx context.Context, transactionID int, deletedBy int) error
	
	// Business expenses
	CreateExpense(ctx context.Context, expense *transactions.BusinessExpense) (*transactions.BusinessExpense, error)
	GetExpenseByID(ctx context.Context, expenseID int) (*transactions.BusinessExpense, error)
	ListExpenses(ctx context.Context, filter *transactions.FinancialFilterParams) ([]*transactions.BusinessExpense, *common.PaginationMeta, error)
	UpdateExpenseStatus(ctx context.Context, expenseID int, status transactions.BusinessExpenseStatus, updatedBy int, notes *string) error
	ApproveExpense(ctx context.Context, expenseID int, approvedBy int, notes *string) error
	RejectExpense(ctx context.Context, expenseID int, rejectedBy int, reason string) error
	
	// Cash flow
	CreateCashFlow(ctx context.Context, cashflow *transactions.CashFlowDaily) (*transactions.CashFlowDaily, error)
	GetCashFlowByDate(ctx context.Context, date string) (*transactions.CashFlowDaily, error)
	ListCashFlow(ctx context.Context, dateFrom, dateTo string) ([]*transactions.CashFlowDaily, error)
	ReconcileCashFlow(ctx context.Context, cashflowID int, reconciledBy int, notes *string) error
	
	// Generate transaction numbers
	GenerateTransactionNumber(ctx context.Context) (string, error)
	GenerateExpenseNumber(ctx context.Context) (string, error)
	
	// Financial reporting helpers
	GetIncomeByDateRange(ctx context.Context, dateFrom, dateTo string) (float64, error)
	GetExpenseByDateRange(ctx context.Context, dateFrom, dateTo string) (float64, error)
	GetBalanceByAccountType(ctx context.Context, accountType transactions.AccountType) (float64, error)
}