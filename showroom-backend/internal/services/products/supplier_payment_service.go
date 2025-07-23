package products

import (
	"context"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// SupplierPaymentService handles business logic for supplier payments
type SupplierPaymentService struct {
	supplierPaymentRepo interfaces.SupplierPaymentRepository
	poRepo              interfaces.PurchaseOrderPartsRepository
}

// NewSupplierPaymentService creates a new supplier payment service
func NewSupplierPaymentService(
	supplierPaymentRepo interfaces.SupplierPaymentRepository,
	poRepo interfaces.PurchaseOrderPartsRepository,
) *SupplierPaymentService {
	return &SupplierPaymentService{
		supplierPaymentRepo: supplierPaymentRepo,
		poRepo:              poRepo,
	}
}

// CreateSupplierPayment creates a new supplier payment
func (s *SupplierPaymentService) CreateSupplierPayment(ctx context.Context, req *products.SupplierPaymentCreateRequest, processedBy int) (*products.SupplierPayment, error) {
	// Generate payment number
	paymentNumber, err := s.supplierPaymentRepo.GenerateNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate payment number: %w", err)
	}

	// Validate PO if provided
	if req.POID != nil {
		_, err := s.poRepo.GetByID(ctx, *req.POID)
		if err != nil {
			return nil, fmt.Errorf("purchase order not found: %w", err)
		}
	}

	// Validate amounts
	if req.PaymentAmount > req.InvoiceAmount {
		return nil, fmt.Errorf("payment amount cannot exceed invoice amount")
	}

	if req.DiscountTaken < 0 {
		return nil, fmt.Errorf("discount taken cannot be negative")
	}

	// Create supplier payment model
	payment := &products.SupplierPayment{
		SupplierID:       req.SupplierID,
		POID:             req.POID,
		PaymentNumber:    paymentNumber,
		InvoiceAmount:    req.InvoiceAmount,
		PaymentAmount:    req.PaymentAmount,
		DiscountTaken:    req.DiscountTaken,
		InvoiceDate:      req.InvoiceDate,
		PaymentDate:      req.PaymentDate,
		DueDate:          req.DueDate,
		PaymentMethod:    req.PaymentMethod,
		PaymentReference: req.PaymentReference,
		InvoiceNumber:    req.InvoiceNumber,
		ProcessedBy:      processedBy,
		PaymentNotes:     req.PaymentNotes,
	}

	// Create the payment
	createdPayment, err := s.supplierPaymentRepo.Create(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to create supplier payment: %w", err)
	}

	return createdPayment, nil
}

// GetSupplierPayment retrieves a supplier payment by ID
func (s *SupplierPaymentService) GetSupplierPayment(ctx context.Context, id int) (*products.SupplierPayment, error) {
	payment, err := s.supplierPaymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get supplier payment: %w", err)
	}

	return payment, nil
}

// GetSupplierPaymentByNumber retrieves a supplier payment by payment number
func (s *SupplierPaymentService) GetSupplierPaymentByNumber(ctx context.Context, number string) (*products.SupplierPayment, error) {
	payment, err := s.supplierPaymentRepo.GetByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get supplier payment by number: %w", err)
	}

	return payment, nil
}

// UpdateSupplierPayment updates a supplier payment
func (s *SupplierPaymentService) UpdateSupplierPayment(ctx context.Context, id int, req *products.SupplierPaymentUpdateRequest) (*products.SupplierPayment, error) {
	// Get existing payment
	existing, err := s.supplierPaymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing payment: %w", err)
	}

	// Check if payment is already fully paid
	if existing.PaymentStatus == products.PaymentStatusPaid {
		return nil, fmt.Errorf("cannot update fully paid supplier payment")
	}

	// Update fields if provided
	if req.PaymentAmount != nil {
		existing.PaymentAmount = *req.PaymentAmount
	}
	if req.DiscountTaken != nil {
		existing.DiscountTaken = *req.DiscountTaken
	}
	if req.PaymentDate != nil {
		existing.PaymentDate = *req.PaymentDate
	}
	if req.DueDate != nil {
		existing.DueDate = *req.DueDate
	}
	if req.PaymentMethod != nil {
		existing.PaymentMethod = *req.PaymentMethod
	}
	if req.PaymentReference != nil {
		existing.PaymentReference = req.PaymentReference
	}
	if req.PaymentStatus != nil {
		existing.PaymentStatus = *req.PaymentStatus
	}
	if req.PenaltyAmount != nil {
		existing.PenaltyAmount = *req.PenaltyAmount
	}
	if req.PaymentNotes != nil {
		existing.PaymentNotes = req.PaymentNotes
	}

	// Validate amounts
	if existing.PaymentAmount > existing.InvoiceAmount {
		return nil, fmt.Errorf("payment amount cannot exceed invoice amount")
	}

	// Update the payment
	updatedPayment, err := s.supplierPaymentRepo.Update(ctx, id, existing)
	if err != nil {
		return nil, fmt.Errorf("failed to update supplier payment: %w", err)
	}

	return updatedPayment, nil
}

// DeleteSupplierPayment deletes a supplier payment
func (s *SupplierPaymentService) DeleteSupplierPayment(ctx context.Context, id int) error {
	// Check if payment can be deleted
	payment, err := s.supplierPaymentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	if payment.PaymentStatus == products.PaymentStatusPaid {
		return fmt.Errorf("cannot delete paid supplier payment")
	}

	err = s.supplierPaymentRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete supplier payment: %w", err)
	}

	return nil
}

// ListSupplierPayments retrieves supplier payments with pagination
func (s *SupplierPaymentService) ListSupplierPayments(ctx context.Context, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error) {
	payments, err := s.supplierPaymentRepo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list supplier payments: %w", err)
	}

	return payments, nil
}

// GetSupplierPayments retrieves payments for a specific supplier
func (s *SupplierPaymentService) GetSupplierPayments(ctx context.Context, supplierID int, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error) {
	payments, err := s.supplierPaymentRepo.GetBySupplierID(ctx, supplierID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get supplier payments: %w", err)
	}

	return payments, nil
}

// GetPOPayments retrieves payments for a specific purchase order
func (s *SupplierPaymentService) GetPOPayments(ctx context.Context, poID int, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error) {
	payments, err := s.supplierPaymentRepo.GetByPOID(ctx, poID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get PO payments: %w", err)
	}

	return payments, nil
}

// GetOverduePayments retrieves overdue payments
func (s *SupplierPaymentService) GetOverduePayments(ctx context.Context, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error) {
	payments, err := s.supplierPaymentRepo.GetOverduePayments(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue payments: %w", err)
	}

	return payments, nil
}

// ProcessPayment processes a payment for an invoice
func (s *SupplierPaymentService) ProcessPayment(ctx context.Context, id int, req *PaymentProcessRequest) error {
	// Get existing payment
	payment, err := s.supplierPaymentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	// Check if payment is already fully paid
	if payment.PaymentStatus == products.PaymentStatusPaid {
		return fmt.Errorf("payment is already fully paid")
	}

	// Validate payment amount
	if req.Amount <= 0 {
		return fmt.Errorf("payment amount must be positive")
	}

	remainingAmount := payment.OutstandingAmount
	if req.Amount > remainingAmount {
		return fmt.Errorf("payment amount exceeds outstanding amount: max %f, requested %f", remainingAmount, req.Amount)
	}

	// Add the payment
	err = s.supplierPaymentRepo.AddPayment(ctx, id, req.Amount, req.Method, req.Reference)
	if err != nil {
		return fmt.Errorf("failed to process payment: %w", err)
	}

	return nil
}

// UpdatePaymentStatus updates the payment status
func (s *SupplierPaymentService) UpdatePaymentStatus(ctx context.Context, id int, status products.PaymentStatus) error {
	err := s.supplierPaymentRepo.UpdatePaymentStatus(ctx, id, status)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	return nil
}

// UpdateOverduePayments updates overdue status for all payments
func (s *SupplierPaymentService) UpdateOverduePayments(ctx context.Context) error {
	err := s.supplierPaymentRepo.UpdateOverdueStatus(ctx)
	if err != nil {
		return fmt.Errorf("failed to update overdue payments: %w", err)
	}

	return nil
}

// GetPaymentSummary gets payment summary for a supplier or all suppliers
func (s *SupplierPaymentService) GetPaymentSummary(ctx context.Context, supplierID *int) (map[string]interface{}, error) {
	summary, err := s.supplierPaymentRepo.GetPaymentSummary(ctx, supplierID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment summary: %w", err)
	}

	return summary, nil
}

// CreatePaymentFromPO creates a payment record from a purchase order
func (s *SupplierPaymentService) CreatePaymentFromPO(ctx context.Context, poID int, req *POPaymentCreateRequest, processedBy int) (*products.SupplierPayment, error) {
	// Get PO details
	po, err := s.poRepo.GetByID(ctx, poID)
	if err != nil {
		return nil, fmt.Errorf("purchase order not found: %w", err)
	}

	// Calculate totals if not provided
	if req.InvoiceAmount <= 0 {
		// Get calculated totals from PO
		poWithTotals, err := s.poRepo.CalculateTotals(ctx, poID)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate PO totals: %w", err)
		}
		req.InvoiceAmount = poWithTotals.TotalAmount
	}

	// Create payment request
	paymentReq := &products.SupplierPaymentCreateRequest{
		SupplierID:       po.SupplierID,
		POID:             &poID,
		InvoiceAmount:    req.InvoiceAmount,
		PaymentAmount:    req.PaymentAmount,
		DiscountTaken:    req.DiscountTaken,
		InvoiceDate:      req.InvoiceDate,
		PaymentDate:      req.PaymentDate,
		DueDate:          req.DueDate,
		PaymentMethod:    req.PaymentMethod,
		PaymentReference: req.PaymentReference,
		InvoiceNumber:    req.InvoiceNumber,
		PaymentNotes:     req.PaymentNotes,
	}

	// Create the payment
	payment, err := s.CreateSupplierPayment(ctx, paymentReq, processedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment from PO: %w", err)
	}

	return payment, nil
}

// CalculatePaymentTerms calculates payment terms based on business rules
func (s *SupplierPaymentService) CalculatePaymentTerms(invoiceDate time.Time, termsDays int) PaymentTerms {
	dueDate := invoiceDate.AddDate(0, 0, termsDays)
	
	// Calculate early payment discount (2% if paid within 10 days)
	earlyPaymentDate := invoiceDate.AddDate(0, 0, 10)
	earlyPaymentDiscount := 0.02 // 2%

	return PaymentTerms{
		DueDate:              dueDate,
		EarlyPaymentDate:     earlyPaymentDate,
		EarlyPaymentDiscount: earlyPaymentDiscount,
		TermsDays:            termsDays,
	}
}

// ValidatePaymentTerms validates payment terms against business rules
func (s *SupplierPaymentService) ValidatePaymentTerms(ctx context.Context, payment *products.SupplierPayment) error {
	// Check if payment date is not in the future
	if payment.PaymentDate.After(time.Now()) {
		return fmt.Errorf("payment date cannot be in the future")
	}

	// Check if invoice date is not after payment date
	if payment.InvoiceDate.After(payment.PaymentDate) {
		return fmt.Errorf("invoice date cannot be after payment date")
	}

	// Check if due date is reasonable (not more than 1 year from invoice date)
	maxDueDate := payment.InvoiceDate.AddDate(1, 0, 0)
	if payment.DueDate.After(maxDueDate) {
		return fmt.Errorf("due date is too far in the future")
	}

	return nil
}

// PaymentProcessRequest represents a request to process a payment
type PaymentProcessRequest struct {
	Amount    float64                   `json:"amount" binding:"required,min=0"`
	Method    products.PaymentMethod    `json:"method" binding:"required"`
	Reference *string                   `json:"reference,omitempty"`
}

// POPaymentCreateRequest represents a request to create payment from PO
type POPaymentCreateRequest struct {
	InvoiceAmount     float64                `json:"invoice_amount" binding:"min=0"`
	PaymentAmount     float64                `json:"payment_amount" binding:"required,min=0"`
	DiscountTaken     float64                `json:"discount_taken" binding:"min=0"`
	InvoiceDate       time.Time              `json:"invoice_date" binding:"required"`
	PaymentDate       time.Time              `json:"payment_date" binding:"required"`
	DueDate           time.Time              `json:"due_date" binding:"required"`
	PaymentMethod     products.PaymentMethod `json:"payment_method" binding:"required"`
	PaymentReference  *string                `json:"payment_reference,omitempty"`
	InvoiceNumber     string                 `json:"invoice_number" binding:"required"`
	PaymentNotes      *string                `json:"payment_notes,omitempty"`
}

// PaymentTerms represents payment terms calculation result
type PaymentTerms struct {
	DueDate              time.Time `json:"due_date"`
	EarlyPaymentDate     time.Time `json:"early_payment_date"`
	EarlyPaymentDiscount float64   `json:"early_payment_discount"`
	TermsDays            int       `json:"terms_days"`
}