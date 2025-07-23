package products

import (
	"context"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// SupplierPaymentService handles business logic for supplier payment management
type SupplierPaymentService struct {
	supplierPaymentRepo interfaces.SupplierPaymentRepository
}

// NewSupplierPaymentService creates a new supplier payment service
func NewSupplierPaymentService(
	supplierPaymentRepo interfaces.SupplierPaymentRepository,
) *SupplierPaymentService {
	return &SupplierPaymentService{
		supplierPaymentRepo: supplierPaymentRepo,
	}
}

// CreateSupplierPayment creates a new supplier payment with business validation
func (s *SupplierPaymentService) CreateSupplierPayment(ctx context.Context, req *products.SupplierPaymentCreateRequest, processedBy int) (*products.SupplierPayment, error) {
	// Validate payment method
	if !req.PaymentMethod.IsValid() {
		return nil, fmt.Errorf("invalid payment method: %s", req.PaymentMethod)
	}

	// Validate dates
	if req.PaymentDate.Before(req.InvoiceDate) {
		return nil, fmt.Errorf("payment date cannot be before invoice date")
	}

	if req.DueDate.Before(req.InvoiceDate) {
		return nil, fmt.Errorf("due date cannot be before invoice date")
	}

	// Validate amounts
	if req.InvoiceAmount <= 0 {
		return nil, fmt.Errorf("invoice amount must be positive")
	}

	if req.PaymentAmount < 0 {
		return nil, fmt.Errorf("payment amount cannot be negative")
	}

	if req.DiscountTaken < 0 {
		return nil, fmt.Errorf("discount taken cannot be negative")
	}

	if (req.PaymentAmount + req.DiscountTaken) > req.InvoiceAmount {
		return nil, fmt.Errorf("payment amount plus discount cannot exceed invoice amount")
	}

	// Generate payment number
	paymentNumber, err := s.supplierPaymentRepo.GenerateNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate payment number: %w", err)
	}

	// Calculate outstanding amount and days overdue
	outstandingAmount := req.InvoiceAmount - req.PaymentAmount - req.DiscountTaken
	
	var daysOverdue int
	if time.Now().After(req.DueDate) && outstandingAmount > 0 {
		daysOverdue = int(time.Now().Sub(req.DueDate).Hours() / 24)
	}

	// Determine payment status
	var paymentStatus products.PaymentStatus
	if outstandingAmount <= 0 {
		paymentStatus = products.PaymentStatusPaid
	} else if req.PaymentAmount > 0 {
		if daysOverdue > 0 {
			paymentStatus = products.PaymentStatusOverdue
		} else {
			paymentStatus = products.PaymentStatusPartial
		}
	} else {
		if daysOverdue > 0 {
			paymentStatus = products.PaymentStatusOverdue
		} else {
			paymentStatus = products.PaymentStatusPending
		}
	}

	// Create payment record
	payment := &products.SupplierPayment{
		SupplierID:        req.SupplierID,
		POID:              req.POID,
		PaymentNumber:     paymentNumber,
		InvoiceAmount:     req.InvoiceAmount,
		PaymentAmount:     req.PaymentAmount,
		DiscountTaken:     req.DiscountTaken,
		OutstandingAmount: outstandingAmount,
		InvoiceDate:       req.InvoiceDate,
		PaymentDate:       req.PaymentDate,
		DueDate:           req.DueDate,
		PaymentMethod:     req.PaymentMethod,
		PaymentReference:  req.PaymentReference,
		InvoiceNumber:     req.InvoiceNumber,
		PaymentStatus:     paymentStatus,
		DaysOverdue:       daysOverdue,
		PenaltyAmount:     0, // Will be calculated if needed
		ProcessedBy:       processedBy,
		PaymentNotes:      req.PaymentNotes,
	}

	// Calculate penalty if overdue
	if daysOverdue > 0 {
		payment.CalculatePenalty(0.1) // 0.1% daily penalty rate
	}

	createdPayment, err := s.supplierPaymentRepo.Create(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to create supplier payment: %w", err)
	}

	return createdPayment, nil
}

// GetSupplierPayments retrieves supplier payments with pagination and filtering
func (s *SupplierPaymentService) GetSupplierPayments(ctx context.Context, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error) {
	return s.supplierPaymentRepo.List(ctx, params)
}

// GetSupplierPaymentByID retrieves a supplier payment by ID
func (s *SupplierPaymentService) GetSupplierPaymentByID(ctx context.Context, id int) (*products.SupplierPayment, error) {
	return s.supplierPaymentRepo.GetByID(ctx, id)
}

// GetSupplierPaymentByNumber retrieves a supplier payment by payment number
func (s *SupplierPaymentService) GetSupplierPaymentByNumber(ctx context.Context, number string) (*products.SupplierPayment, error) {
	return s.supplierPaymentRepo.GetByNumber(ctx, number)
}

// GetSupplierPaymentsBySupplier retrieves payments for a specific supplier
func (s *SupplierPaymentService) GetSupplierPaymentsBySupplier(ctx context.Context, supplierID int, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error) {
	return s.supplierPaymentRepo.GetBySupplierID(ctx, supplierID, params)
}

// GetSupplierPaymentsByPO retrieves payments for a specific purchase order
func (s *SupplierPaymentService) GetSupplierPaymentsByPO(ctx context.Context, poID int, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error) {
	return s.supplierPaymentRepo.GetByPOID(ctx, poID, params)
}

// GetOverduePayments retrieves overdue payments
func (s *SupplierPaymentService) GetOverduePayments(ctx context.Context, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error) {
	return s.supplierPaymentRepo.GetOverduePayments(ctx, params)
}

// UpdateSupplierPayment updates a supplier payment
func (s *SupplierPaymentService) UpdateSupplierPayment(ctx context.Context, id int, req *products.SupplierPaymentUpdateRequest) (*products.SupplierPayment, error) {
	// Get existing payment
	payment, err := s.supplierPaymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("supplier payment not found: %w", err)
	}

	// Check if payment is not fully paid (prevent modification of completed payments)
	if payment.IsFullyPaid() {
		return nil, fmt.Errorf("cannot modify fully paid invoices")
	}

	// Update fields if provided
	if req.PaymentAmount != nil {
		if *req.PaymentAmount < 0 {
			return nil, fmt.Errorf("payment amount cannot be negative")
		}
		payment.PaymentAmount = *req.PaymentAmount
	}

	if req.DiscountTaken != nil {
		if *req.DiscountTaken < 0 {
			return nil, fmt.Errorf("discount taken cannot be negative")
		}
		payment.DiscountTaken = *req.DiscountTaken
	}

	if req.PaymentDate != nil {
		payment.PaymentDate = *req.PaymentDate
	}

	if req.DueDate != nil {
		payment.DueDate = *req.DueDate
	}

	if req.PaymentMethod != nil {
		if !req.PaymentMethod.IsValid() {
			return nil, fmt.Errorf("invalid payment method: %s", *req.PaymentMethod)
		}
		payment.PaymentMethod = *req.PaymentMethod
	}

	if req.PaymentReference != nil {
		payment.PaymentReference = req.PaymentReference
	}

	if req.PaymentStatus != nil {
		if !req.PaymentStatus.IsValid() {
			return nil, fmt.Errorf("invalid payment status: %s", *req.PaymentStatus)
		}
		payment.PaymentStatus = *req.PaymentStatus
	}

	if req.PenaltyAmount != nil {
		if *req.PenaltyAmount < 0 {
			return nil, fmt.Errorf("penalty amount cannot be negative")
		}
		payment.PenaltyAmount = *req.PenaltyAmount
	}

	if req.PaymentNotes != nil {
		payment.PaymentNotes = req.PaymentNotes
	}

	// Validate amounts don't exceed invoice amount
	if (payment.PaymentAmount + payment.DiscountTaken) > payment.InvoiceAmount {
		return nil, fmt.Errorf("payment amount plus discount cannot exceed invoice amount")
	}

	// Recalculate outstanding amount and status
	payment.UpdatePaymentStatus()

	// Update penalty if overdue
	if payment.DaysOverdue > 0 {
		payment.CalculatePenalty(0.1) // 0.1% daily penalty rate
	}

	return s.supplierPaymentRepo.Update(ctx, id, payment)
}

// AddPayment adds a payment to an existing payment record
func (s *SupplierPaymentService) AddPayment(ctx context.Context, id int, amount float64, method products.PaymentMethod, reference *string) error {
	// Validate amount
	if amount <= 0 {
		return fmt.Errorf("payment amount must be positive")
	}

	// Validate payment method
	if !method.IsValid() {
		return fmt.Errorf("invalid payment method: %s", method)
	}

	// Get existing payment
	payment, err := s.supplierPaymentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("supplier payment not found: %w", err)
	}

	// Check if payment can accept additional amounts
	if !payment.CanAddPayment() {
		return fmt.Errorf("cannot add payment to this invoice")
	}

	// Check if amount exceeds outstanding
	if amount > payment.GetRemainingAmount() {
		return fmt.Errorf("payment amount exceeds outstanding amount")
	}

	return s.supplierPaymentRepo.AddPayment(ctx, id, amount, method, reference)
}

// UpdatePaymentStatus updates the payment status
func (s *SupplierPaymentService) UpdatePaymentStatus(ctx context.Context, id int, status products.PaymentStatus) error {
	// Validate status
	if !status.IsValid() {
		return fmt.Errorf("invalid payment status: %s", status)
	}

	return s.supplierPaymentRepo.UpdatePaymentStatus(ctx, id, status)
}

// ProcessOverduePayments updates overdue status for all payments
func (s *SupplierPaymentService) ProcessOverduePayments(ctx context.Context) error {
	return s.supplierPaymentRepo.UpdateOverdueStatus(ctx)
}

// GetPaymentSummary retrieves payment summary for analysis
func (s *SupplierPaymentService) GetPaymentSummary(ctx context.Context, supplierID *int) (map[string]interface{}, error) {
	return s.supplierPaymentRepo.GetPaymentSummary(ctx, supplierID)
}

// DeleteSupplierPayment deletes a supplier payment (with business rules)
func (s *SupplierPaymentService) DeleteSupplierPayment(ctx context.Context, id int) error {
	// Get payment
	payment, err := s.supplierPaymentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("supplier payment not found: %w", err)
	}

	// Prevent deletion of payments with positive amounts (has financial impact)
	if payment.PaymentAmount > 0 {
		return fmt.Errorf("cannot delete payments with recorded payment amounts")
	}

	// Only allow deletion of pending payments
	if payment.PaymentStatus != products.PaymentStatusPending {
		return fmt.Errorf("can only delete pending payments")
	}

	return s.supplierPaymentRepo.Delete(ctx, id)
}

// CalculateAgeingReport calculates ageing report for supplier payments
func (s *SupplierPaymentService) CalculateAgeingReport(ctx context.Context, supplierID *int) (map[string]interface{}, error) {
	// This would require custom queries for ageing buckets
	// For now, return payment summary as base
	summary, err := s.GetPaymentSummary(ctx, supplierID)
	if err != nil {
		return nil, err
	}

	// Add ageing buckets (would need custom implementation)
	summary["ageing"] = map[string]float64{
		"current":    0.0, // 0-30 days
		"30_days":    0.0, // 31-60 days
		"60_days":    0.0, // 61-90 days
		"90_days":    0.0, // 91+ days
	}

	summary["note"] = "Ageing report implementation requires custom queries"

	return summary, nil
}

// BulkUpdatePaymentStatus updates payment status for multiple payments
func (s *SupplierPaymentService) BulkUpdatePaymentStatus(ctx context.Context, paymentIDs []int, status products.PaymentStatus) error {
	// Validate status
	if !status.IsValid() {
		return fmt.Errorf("invalid payment status: %s", status)
	}

	// Update each payment
	for _, id := range paymentIDs {
		err := s.supplierPaymentRepo.UpdatePaymentStatus(ctx, id, status)
		if err != nil {
			return fmt.Errorf("failed to update payment %d: %w", id, err)
		}
	}

	return nil
}

// SchedulePaymentReminders would schedule reminders for due payments
func (s *SupplierPaymentService) SchedulePaymentReminders(ctx context.Context, daysBefore int) ([]int, error) {
	// This would integrate with a notification system
	// For now, return empty list
	var reminderPaymentIDs []int

	// Implementation would:
	// 1. Find payments due in X days
	// 2. Create reminder notifications
	// 3. Return list of payment IDs for which reminders were scheduled

	return reminderPaymentIDs, nil
}

// ValidatePaymentData validates payment data for batch imports
func (s *SupplierPaymentService) ValidatePaymentData(ctx context.Context, payments []products.SupplierPaymentCreateRequest) []error {
	var errors []error

	for i, payment := range payments {
		// Check payment method
		if !payment.PaymentMethod.IsValid() {
			errors = append(errors, fmt.Errorf("row %d: invalid payment method %s", i+1, payment.PaymentMethod))
		}

		// Check amounts
		if payment.InvoiceAmount <= 0 {
			errors = append(errors, fmt.Errorf("row %d: invoice amount must be positive", i+1))
		}

		if payment.PaymentAmount < 0 {
			errors = append(errors, fmt.Errorf("row %d: payment amount cannot be negative", i+1))
		}

		if (payment.PaymentAmount + payment.DiscountTaken) > payment.InvoiceAmount {
			errors = append(errors, fmt.Errorf("row %d: payment plus discount exceeds invoice amount", i+1))
		}

		// Check dates
		if payment.PaymentDate.Before(payment.InvoiceDate) {
			errors = append(errors, fmt.Errorf("row %d: payment date cannot be before invoice date", i+1))
		}

		// Check for duplicate invoice numbers (would need repository support)
		// This is a placeholder for business rule validation
	}

	return errors
}