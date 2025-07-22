package inventory

import (
	"context"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// SupplierPaymentService handles business logic for supplier payments
type SupplierPaymentService struct {
	supplierPaymentRepo interfaces.SupplierPaymentRepository
	purchaseOrderRepo   interfaces.PurchaseOrderRepository
	codeGenerator       *inventory.CodeGenerator
}

// NewSupplierPaymentService creates a new supplier payment service
func NewSupplierPaymentService(
	supplierPaymentRepo interfaces.SupplierPaymentRepository,
	purchaseOrderRepo interfaces.PurchaseOrderRepository,
	codeGenerator *inventory.CodeGenerator,
) *SupplierPaymentService {
	return &SupplierPaymentService{
		supplierPaymentRepo: supplierPaymentRepo,
		purchaseOrderRepo:   purchaseOrderRepo,
		codeGenerator:       codeGenerator,
	}
}

// CreateSupplierPayment creates a new supplier payment
func (s *SupplierPaymentService) CreateSupplierPayment(ctx context.Context, req *inventory.SupplierPaymentCreateRequest, userID int) (*inventory.SupplierPayment, error) {
	// Validate purchase order if provided
	if req.PoID != nil {
		_, err := s.purchaseOrderRepo.GetByID(ctx, *req.PoID)
		if err != nil {
			return nil, fmt.Errorf("purchase order not found: %w", err)
		}
	}

	// Generate payment number
	paymentNumber, err := s.codeGenerator.GenerateSupplierPaymentNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate payment number: %w", err)
	}

	// Calculate outstanding amount
	outstandingAmount := req.InvoiceAmount - req.PaymentAmount - req.DiscountTaken
	if outstandingAmount < 0 {
		outstandingAmount = 0
	}

	// Determine payment status
	var paymentStatus inventory.PaymentStatus
	if req.PaymentAmount == 0 {
		paymentStatus = inventory.PaymentStatusPending
	} else if outstandingAmount > 0 {
		paymentStatus = inventory.PaymentStatusPartial
	} else {
		paymentStatus = inventory.PaymentStatusPaid
	}

	// Create supplier payment
	payment := &inventory.SupplierPayment{
		SupplierID:        req.SupplierID,
		PoID:              req.PoID,
		PaymentNumber:     paymentNumber,
		InvoiceAmount:     req.InvoiceAmount,
		PaymentAmount:     req.PaymentAmount,
		DiscountTaken:     req.DiscountTaken,
		OutstandingAmount: outstandingAmount,
		InvoiceDate:       req.InvoiceDate,
		PaymentDate:       time.Now(),
		DueDate:           req.DueDate,
		PaymentMethod:     req.PaymentMethod,
		PaymentReference:  req.PaymentReference,
		InvoiceNumber:     req.InvoiceNumber,
		PaymentStatus:     paymentStatus,
		DaysOverdue:       0,
		PenaltyAmount:     0,
		ProcessedBy:       userID,
		PaymentNotes:      req.PaymentNotes,
	}

	// Calculate days overdue if due date is provided
	if payment.DueDate != nil {
		payment.CalculateDaysOverdue()
		if payment.DaysOverdue > 0 && payment.PaymentStatus == inventory.PaymentStatusPending {
			payment.PaymentStatus = inventory.PaymentStatusOverdue
		}
	}

	// Create payment record
	createdPayment, err := s.supplierPaymentRepo.Create(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to create supplier payment: %w", err)
	}

	return createdPayment, nil
}

// GetSupplierPayment retrieves a supplier payment by ID
func (s *SupplierPaymentService) GetSupplierPayment(ctx context.Context, id int) (*inventory.SupplierPayment, error) {
	return s.supplierPaymentRepo.GetByID(ctx, id)
}

// ListSupplierPayments retrieves supplier payments with filtering
func (s *SupplierPaymentService) ListSupplierPayments(ctx context.Context, params *inventory.SupplierPaymentFilterParams) ([]inventory.SupplierPaymentListItem, int, error) {
	return s.supplierPaymentRepo.List(ctx, params)
}

// UpdateSupplierPayment updates a supplier payment
func (s *SupplierPaymentService) UpdateSupplierPayment(ctx context.Context, id int, req *inventory.SupplierPaymentUpdateRequest) (*inventory.SupplierPayment, error) {
	payment, err := s.supplierPaymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.PaymentAmount != nil {
		payment.PaymentAmount = *req.PaymentAmount
	}
	if req.DiscountTaken != nil {
		payment.DiscountTaken = *req.DiscountTaken
	}
	if req.DueDate != nil {
		payment.DueDate = req.DueDate
	}
	if req.PaymentMethod != nil {
		payment.PaymentMethod = *req.PaymentMethod
	}
	if req.PaymentReference != nil {
		payment.PaymentReference = req.PaymentReference
	}
	if req.PaymentStatus != nil {
		payment.PaymentStatus = *req.PaymentStatus
	}
	if req.PaymentNotes != nil {
		payment.PaymentNotes = req.PaymentNotes
	}

	// Recalculate outstanding amount and status
	payment.CalculateOutstandingAmount()
	if payment.PaymentAmount >= payment.InvoiceAmount-payment.DiscountTaken {
		payment.PaymentStatus = inventory.PaymentStatusPaid
	} else if payment.PaymentAmount > 0 {
		payment.PaymentStatus = inventory.PaymentStatusPartial
	}

	// Recalculate days overdue
	payment.CalculateDaysOverdue()
	if payment.DaysOverdue > 0 && payment.PaymentStatus != inventory.PaymentStatusPaid {
		payment.PaymentStatus = inventory.PaymentStatusOverdue
	}

	return s.supplierPaymentRepo.Update(ctx, id, payment)
}

// DeleteSupplierPayment deletes a supplier payment
func (s *SupplierPaymentService) DeleteSupplierPayment(ctx context.Context, id int) error {
	// Check if payment exists
	_, err := s.supplierPaymentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.supplierPaymentRepo.Delete(ctx, id)
}

// GetPaymentsBySupplier retrieves payments for a specific supplier
func (s *SupplierPaymentService) GetPaymentsBySupplier(ctx context.Context, supplierID int, page, limit int) ([]inventory.SupplierPaymentListItem, int, error) {
	return s.supplierPaymentRepo.GetBySupplier(ctx, supplierID, page, limit)
}

// GetPaymentsByPO retrieves payments for a specific purchase order
func (s *SupplierPaymentService) GetPaymentsByPO(ctx context.Context, poID int, page, limit int) ([]inventory.SupplierPaymentListItem, int, error) {
	return s.supplierPaymentRepo.GetByPO(ctx, poID, page, limit)
}

// GetPaymentsByStatus retrieves payments by status
func (s *SupplierPaymentService) GetPaymentsByStatus(ctx context.Context, status inventory.PaymentStatus, page, limit int) ([]inventory.SupplierPaymentListItem, int, error) {
	return s.supplierPaymentRepo.GetByStatus(ctx, status, page, limit)
}

// GetOverduePayments retrieves overdue payments
func (s *SupplierPaymentService) GetOverduePayments(ctx context.Context, page, limit int) ([]inventory.SupplierPaymentListItem, int, error) {
	return s.supplierPaymentRepo.GetOverduePayments(ctx, page, limit)
}

// SearchPayments searches payments
func (s *SupplierPaymentService) SearchPayments(ctx context.Context, query string, page, limit int) ([]inventory.SupplierPaymentListItem, int, error) {
	return s.supplierPaymentRepo.Search(ctx, query, page, limit)
}

// ProcessPayment processes a payment for an invoice
func (s *SupplierPaymentService) ProcessPayment(ctx context.Context, supplierID int, invoiceAmount, paymentAmount, discountTaken float64, paymentMethod inventory.PaymentMethod, reference *string, userID int) (*inventory.SupplierPayment, error) {
	req := &inventory.SupplierPaymentCreateRequest{
		SupplierID:       supplierID,
		InvoiceAmount:    invoiceAmount,
		PaymentAmount:    paymentAmount,
		DiscountTaken:    discountTaken,
		PaymentMethod:    paymentMethod,
		PaymentReference: reference,
		InvoiceDate:      nil,
		DueDate:          nil,
		PaymentNotes:     nil,
	}

	return s.CreateSupplierPayment(ctx, req, userID)
}

// MarkPaymentPaid marks a payment as paid
func (s *SupplierPaymentService) MarkPaymentPaid(ctx context.Context, id int) error {
	payment, err := s.supplierPaymentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Set payment amount to invoice amount minus discount
	payment.PaymentAmount = payment.InvoiceAmount - payment.DiscountTaken
	payment.OutstandingAmount = 0
	payment.PaymentStatus = inventory.PaymentStatusPaid
	payment.PaymentDate = time.Now()

	_, err = s.supplierPaymentRepo.Update(ctx, id, payment)
	return err
}

// UpdateOverdueStatus updates overdue status for all payments
func (s *SupplierPaymentService) UpdateOverdueStatus(ctx context.Context) error {
	return s.supplierPaymentRepo.UpdateOverdueStatus(ctx)
}

// GetPaymentSummary retrieves payment summary for reporting
func (s *SupplierPaymentService) GetPaymentSummary(ctx context.Context, supplierID *int, startDate, endDate string) (*inventory.PaymentSummary, error) {
	return s.supplierPaymentRepo.GetPaymentSummary(ctx, supplierID, startDate, endDate)
}

// GetOutstandingBalance retrieves outstanding balance for a supplier
func (s *SupplierPaymentService) GetOutstandingBalance(ctx context.Context, supplierID int) (float64, error) {
	return s.supplierPaymentRepo.GetOutstandingBalance(ctx, supplierID)
}

// GetTotalPaid retrieves total paid amount for a supplier within date range
func (s *SupplierPaymentService) GetTotalPaid(ctx context.Context, supplierID int, startDate, endDate string) (float64, error) {
	return s.supplierPaymentRepo.GetTotalPaid(ctx, supplierID, startDate, endDate)
}

// CreatePOPayment creates a payment for a specific purchase order
func (s *SupplierPaymentService) CreatePOPayment(ctx context.Context, poID int, paymentAmount, discountTaken float64, paymentMethod inventory.PaymentMethod, reference *string, userID int) (*inventory.SupplierPayment, error) {
	po, err := s.purchaseOrderRepo.GetByID(ctx, poID)
	if err != nil {
		return nil, fmt.Errorf("purchase order not found: %w", err)
	}

	req := &inventory.SupplierPaymentCreateRequest{
		SupplierID:       po.SupplierID,
		PoID:             &poID,
		InvoiceAmount:    po.TotalAmount,
		PaymentAmount:    paymentAmount,
		DiscountTaken:    discountTaken,
		PaymentMethod:    paymentMethod,
		PaymentReference: reference,
		InvoiceDate:      &po.PoDate,
		DueDate:          po.PaymentDueDate,
		PaymentNotes:     nil,
	}

	return s.CreateSupplierPayment(ctx, req, userID)
}