package services

import (
	"context"
	"fmt"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/vehicle_purchase"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// VehiclePurchaseService handles business logic for vehicle purchase operations
type VehiclePurchaseService struct {
	transactionRepo interfaces.VehiclePurchaseTransactionRepository
	paymentRepo     interfaces.VehiclePurchasePaymentRepository
	customerRepo    interfaces.CustomerRepository
	userRepo        interfaces.UserRepository
}

// NewVehiclePurchaseService creates a new vehicle purchase service
func NewVehiclePurchaseService(
	transactionRepo interfaces.VehiclePurchaseTransactionRepository,
	paymentRepo interfaces.VehiclePurchasePaymentRepository,
	customerRepo interfaces.CustomerRepository,
	userRepo interfaces.UserRepository,
) *VehiclePurchaseService {
	return &VehiclePurchaseService{
		transactionRepo: transactionRepo,
		paymentRepo:     paymentRepo,
		customerRepo:    customerRepo,
		userRepo:        userRepo,
	}
}

// Transaction Methods

// CreateTransaction creates a new vehicle purchase transaction
func (s *VehiclePurchaseService) CreateTransaction(ctx context.Context, req *vehicle_purchase.CreateVehiclePurchaseTransactionRequest, userID int) (*vehicle_purchase.VehiclePurchaseTransaction, error) {
	// Validate customer exists
	customer, err := s.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Check if VIN already exists
	if req.VinNumber != "" {
		exists, err := s.transactionRepo.IsVINExists(ctx, req.VinNumber, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to check VIN uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("VIN number already exists")
		}
	}

	// Generate transaction number
	transactionNumber, err := s.transactionRepo.GenerateNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate transaction number: %w", err)
	}

	// Create transaction
	transaction := &vehicle_purchase.VehiclePurchaseTransaction{
		TransactionNumber:   transactionNumber,
		CustomerID:         req.CustomerID,
		VehicleID:          req.VehicleID,
		VinNumber:          req.VinNumber,
		VehicleBrand:       req.VehicleBrand,
		VehicleModel:       req.VehicleModel,
		VehicleYear:        req.VehicleYear,
		VehicleColor:       req.VehicleColor,
		EngineNumber:       req.EngineNumber,
		RegistrationNumber: req.RegistrationNumber,
		PurchasePrice:      req.PurchasePrice,
		AgreedValue:        req.AgreedValue,
		OdometerReading:    req.OdometerReading,
		FuelType:           req.FuelType,
		Transmission:       req.Transmission,
		ConditionRating:    req.ConditionRating,
		PurchaseDate:       req.PurchaseDate,
		TransactionStatus:  "pending",
		EvaluationNotes:    req.EvaluationNotes,
		PurchaseNotes:      req.PurchaseNotes,
		DocumentsJSON:      req.DocumentsJSON,
		ProcessedBy:        userID,
		CustomerName:       customer.CustomerName,
	}

	return s.transactionRepo.Create(ctx, transaction)
}

// GetTransaction retrieves a transaction by ID
func (s *VehiclePurchaseService) GetTransaction(ctx context.Context, id int) (*vehicle_purchase.VehiclePurchaseTransaction, error) {
	return s.transactionRepo.GetByID(ctx, id)
}

// GetTransactionByNumber retrieves a transaction by number
func (s *VehiclePurchaseService) GetTransactionByNumber(ctx context.Context, number string) (*vehicle_purchase.VehiclePurchaseTransaction, error) {
	return s.transactionRepo.GetByNumber(ctx, number)
}

// GetTransactionByVIN retrieves a transaction by VIN
func (s *VehiclePurchaseService) GetTransactionByVIN(ctx context.Context, vin string) (*vehicle_purchase.VehiclePurchaseTransaction, error) {
	return s.transactionRepo.GetByVIN(ctx, vin)
}

// UpdateTransaction updates a transaction
func (s *VehiclePurchaseService) UpdateTransaction(ctx context.Context, id int, req *vehicle_purchase.UpdateVehiclePurchaseTransactionRequest) (*vehicle_purchase.VehiclePurchaseTransaction, error) {
	// Check if VIN already exists (excluding current transaction)
	if req.VinNumber != "" {
		exists, err := s.transactionRepo.IsVINExists(ctx, req.VinNumber, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check VIN uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("VIN number already exists")
		}
	}

	transaction := &vehicle_purchase.VehiclePurchaseTransaction{
		VinNumber:          req.VinNumber,
		VehicleColor:       req.VehicleColor,
		EngineNumber:       req.EngineNumber,
		RegistrationNumber: req.RegistrationNumber,
		PurchasePrice:      req.PurchasePrice,
		AgreedValue:        req.AgreedValue,
		OdometerReading:    req.OdometerReading,
		ConditionRating:    req.ConditionRating,
		TransactionStatus:  req.TransactionStatus,
		InspectionNotes:    req.InspectionNotes,
		EvaluationNotes:    req.EvaluationNotes,
		PurchaseNotes:      req.PurchaseNotes,
		DocumentsJSON:      req.DocumentsJSON,
	}

	return s.transactionRepo.Update(ctx, id, transaction)
}

// CompleteInspection completes vehicle inspection
func (s *VehiclePurchaseService) CompleteInspection(ctx context.Context, id int, req *vehicle_purchase.TransactionInspectionRequest, inspectedBy int) error {
	return s.transactionRepo.CompleteInspection(ctx, id, req, inspectedBy)
}

// ProcessApproval processes transaction approval
func (s *VehiclePurchaseService) ProcessApproval(ctx context.Context, id int, req *vehicle_purchase.TransactionStatusApprovalRequest, approvedBy int) error {
	return s.transactionRepo.ProcessApproval(ctx, id, req, approvedBy)
}

// GetDashboardStats gets dashboard statistics
func (s *VehiclePurchaseService) GetDashboardStats(ctx context.Context) (*vehicle_purchase.TransactionDashboardStats, error) {
	return s.transactionRepo.GetDashboardStats(ctx)
}

// Payment Methods

// CreatePayment creates a new payment
func (s *VehiclePurchaseService) CreatePayment(ctx context.Context, req *vehicle_purchase.CreateVehiclePurchasePaymentRequest, userID int) (*vehicle_purchase.VehiclePurchasePayment, error) {
	// Validate transaction exists
	_, err := s.transactionRepo.GetByID(ctx, req.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	// Generate payment number
	paymentNumber, err := s.paymentRepo.GenerateNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate payment number: %w", err)
	}

	payment := &vehicle_purchase.VehiclePurchasePayment{
		PaymentNumber:   paymentNumber,
		TransactionID:   req.TransactionID,
		PaymentAmount:   req.PaymentAmount,
		PaymentMethod:   req.PaymentMethod,
		PaymentReference: req.PaymentReference,
		BankName:        req.BankName,
		PaymentStatus:   "pending",
		PaymentNotes:    req.PaymentNotes,
		DueDate:         req.DueDate,
		ProcessedBy:     userID,
	}

	return s.paymentRepo.Create(ctx, payment)
}

// GetPayment retrieves a payment by ID
func (s *VehiclePurchaseService) GetPayment(ctx context.Context, id int) (*vehicle_purchase.VehiclePurchasePayment, error) {
	return s.paymentRepo.GetByID(ctx, id)
}

// GetPaymentByNumber retrieves a payment by number
func (s *VehiclePurchaseService) GetPaymentByNumber(ctx context.Context, number string) (*vehicle_purchase.VehiclePurchasePayment, error) {
	return s.paymentRepo.GetByNumber(ctx, number)
}

// ProcessPayment processes a payment
func (s *VehiclePurchaseService) ProcessPayment(ctx context.Context, id int, req *vehicle_purchase.PaymentProcessRequest, processedBy int) error {
	return s.paymentRepo.ProcessPayment(ctx, id, req, processedBy)
}

// ProcessPaymentApproval processes payment approval
func (s *VehiclePurchaseService) ProcessPaymentApproval(ctx context.Context, id int, req *vehicle_purchase.PaymentApprovalRequest, approvedBy int) error {
	return s.paymentRepo.ProcessApproval(ctx, id, req, approvedBy)
}

// GetPaymentSummary gets payment summary for a transaction
func (s *VehiclePurchaseService) GetPaymentSummary(ctx context.Context, transactionID int) (*vehicle_purchase.PaymentSummary, error) {
	return s.paymentRepo.GetPaymentSummary(ctx, transactionID)
}

// List Methods (stubs for now)
func (s *VehiclePurchaseService) ListTransactions(ctx context.Context, params *vehicle_purchase.VehiclePurchaseTransactionFilterParams) (*common.PaginatedResponse, error) {
	return s.transactionRepo.List(ctx, params)
}

func (s *VehiclePurchaseService) ListPayments(ctx context.Context, params *vehicle_purchase.VehiclePurchasePaymentFilterParams) (*common.PaginatedResponse, error) {
	return s.paymentRepo.List(ctx, params)
}

func (s *VehiclePurchaseService) GetTransactionsByCustomer(ctx context.Context, customerID int, params *vehicle_purchase.VehiclePurchaseTransactionFilterParams) (*common.PaginatedResponse, error) {
	return s.transactionRepo.GetByCustomerID(ctx, customerID, params)
}

func (s *VehiclePurchaseService) GetTransactionsByStatus(ctx context.Context, status string, params *vehicle_purchase.VehiclePurchaseTransactionFilterParams) (*common.PaginatedResponse, error) {
	return s.transactionRepo.GetByStatus(ctx, status, params)
}

func (s *VehiclePurchaseService) GetPendingInspections(ctx context.Context, params *vehicle_purchase.VehiclePurchaseTransactionFilterParams) (*common.PaginatedResponse, error) {
	return s.transactionRepo.GetPendingInspection(ctx, params)
}

func (s *VehiclePurchaseService) GetPendingApprovals(ctx context.Context, params *vehicle_purchase.VehiclePurchaseTransactionFilterParams) (*common.PaginatedResponse, error) {
	return s.transactionRepo.GetPendingApproval(ctx, params)
}

func (s *VehiclePurchaseService) GetPaymentsByTransaction(ctx context.Context, transactionID int, params *vehicle_purchase.VehiclePurchasePaymentFilterParams) (*common.PaginatedResponse, error) {
	return s.paymentRepo.GetByTransactionID(ctx, transactionID, params)
}

func (s *VehiclePurchaseService) GetPaymentsByStatus(ctx context.Context, status string, params *vehicle_purchase.VehiclePurchasePaymentFilterParams) (*common.PaginatedResponse, error) {
	return s.paymentRepo.GetByStatus(ctx, status, params)
}

func (s *VehiclePurchaseService) GetPendingPaymentApprovals(ctx context.Context, params *vehicle_purchase.VehiclePurchasePaymentFilterParams) (*common.PaginatedResponse, error) {
	return s.paymentRepo.GetPendingApproval(ctx, params)
}

func (s *VehiclePurchaseService) GetOverduePayments(ctx context.Context, params *vehicle_purchase.VehiclePurchasePaymentFilterParams) (*common.PaginatedResponse, error) {
	return s.paymentRepo.GetOverduePayments(ctx, params)
}