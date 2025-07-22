package inventory

import (
	"context"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// PurchaseOrderService provides business logic for purchase order operations
type PurchaseOrderService struct {
	poRepo        interfaces.PurchaseOrderRepository
	productRepo   interfaces.ProductSparePartRepository
	movementRepo  interfaces.StockMovementRepository
	codeGenerator *inventory.CodeGenerator
}

// NewPurchaseOrderService creates a new purchase order service
func NewPurchaseOrderService(
	poRepo interfaces.PurchaseOrderRepository,
	productRepo interfaces.ProductSparePartRepository,
	movementRepo interfaces.StockMovementRepository,
) *PurchaseOrderService {
	return &PurchaseOrderService{
		poRepo:        poRepo,
		productRepo:   productRepo,
		movementRepo:  movementRepo,
		codeGenerator: inventory.NewCodeGenerator(),
	}
}

// Create creates a new purchase order with auto-generated number
func (s *PurchaseOrderService) Create(ctx context.Context, req *inventory.PurchaseOrderPartCreateRequest, createdBy int) (*inventory.PurchaseOrderPart, error) {
	// Generate PO number
	poNumber, err := s.codeGenerator.GetNextPONumber(func() (int, error) {
		return s.poRepo.GetLastPOID(ctx)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate PO number: %w", err)
	}

	// Validate details
	if len(req.Details) == 0 {
		return nil, fmt.Errorf("purchase order must have at least one item")
	}

	// Create purchase order entity
	po := &inventory.PurchaseOrderPart{
		PoNumber:             poNumber,
		SupplierID:           req.SupplierID,
		PoDate:               time.Now(),
		RequiredDate:         req.RequiredDate,
		ExpectedDeliveryDate: req.ExpectedDeliveryDate,
		PoType:               req.PoType,
		Status:               inventory.POStatusDraft,
		PaymentTerms:         req.PaymentTerms,
		CreatedBy:            createdBy,
		DeliveryAddress:      req.DeliveryAddress,
		PoNotes:              req.PoNotes,
		TermsAndConditions:   req.TermsAndConditions,
	}

	// Calculate totals from details
	var subtotal float64
	for _, detail := range req.Details {
		// Validate product exists
		_, err := s.productRepo.GetByID(ctx, detail.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %d not found: %w", detail.ProductID, err)
		}

		totalCost := float64(detail.QuantityOrdered) * detail.UnitCost
		subtotal += totalCost
	}

	po.Subtotal = subtotal
	po.TotalAmount = subtotal + po.TaxAmount + po.ShippingCost - po.DiscountAmount

	// Update payment due date based on terms
	po.UpdatePaymentDueDate()

	// Create purchase order
	createdPO, err := s.poRepo.Create(ctx, po)
	if err != nil {
		return nil, fmt.Errorf("failed to create purchase order: %w", err)
	}

	return createdPO, nil
}

// GetByID retrieves a purchase order by ID
func (s *PurchaseOrderService) GetByID(ctx context.Context, id int) (*inventory.PurchaseOrderPart, error) {
	return s.poRepo.GetByID(ctx, id)
}

// GetByNumber retrieves a purchase order by number
func (s *PurchaseOrderService) GetByNumber(ctx context.Context, number string) (*inventory.PurchaseOrderPart, error) {
	return s.poRepo.GetByNumber(ctx, number)
}

// Update updates a purchase order
func (s *PurchaseOrderService) Update(ctx context.Context, id int, req *inventory.PurchaseOrderPartUpdateRequest) (*inventory.PurchaseOrderPart, error) {
	// Get existing PO
	existing, err := s.poRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("purchase order not found: %w", err)
	}

	// Check if PO can be modified
	if !existing.CanBeModified() {
		return nil, fmt.Errorf("purchase order cannot be modified in current status: %s", existing.Status)
	}

	// Update fields
	if req.RequiredDate != nil {
		existing.RequiredDate = req.RequiredDate
	}
	if req.ExpectedDeliveryDate != nil {
		existing.ExpectedDeliveryDate = req.ExpectedDeliveryDate
	}
	if req.PoType != nil {
		existing.PoType = *req.PoType
	}
	if req.PaymentTerms != nil {
		existing.PaymentTerms = *req.PaymentTerms
		existing.UpdatePaymentDueDate()
	}
	if req.DeliveryAddress != nil {
		existing.DeliveryAddress = req.DeliveryAddress
	}
	if req.PoNotes != nil {
		existing.PoNotes = req.PoNotes
	}
	if req.TermsAndConditions != nil {
		existing.TermsAndConditions = req.TermsAndConditions
	}

	return s.poRepo.Update(ctx, id, existing)
}

// Delete deletes a purchase order
func (s *PurchaseOrderService) Delete(ctx context.Context, id int) error {
	// Get existing PO
	existing, err := s.poRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("purchase order not found: %w", err)
	}

	// Check if PO can be deleted
	if !existing.CanBeModified() {
		return fmt.Errorf("purchase order cannot be deleted in current status: %s", existing.Status)
	}

	return s.poRepo.Delete(ctx, id)
}

// List retrieves purchase orders with filtering and pagination
func (s *PurchaseOrderService) List(ctx context.Context, params *inventory.PurchaseOrderPartFilterParams) ([]inventory.PurchaseOrderPartListItem, int, error) {
	return s.poRepo.List(ctx, params)
}

// Search searches purchase orders
func (s *PurchaseOrderService) Search(ctx context.Context, query string, page, limit int) ([]inventory.PurchaseOrderPartListItem, int, error) {
	return s.poRepo.Search(ctx, query, page, limit)
}

// SendToSupplier sends a purchase order to supplier
func (s *PurchaseOrderService) SendToSupplier(ctx context.Context, id int, userID int) error {
	// Get existing PO
	po, err := s.poRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("purchase order not found: %w", err)
	}

	// Check if PO can be sent
	if !po.CanBeSent() {
		return fmt.Errorf("purchase order cannot be sent in current status: %s", po.Status)
	}

	// Update status to sent
	return s.poRepo.UpdateStatus(ctx, id, inventory.POStatusSent)
}

// Approve approves a purchase order
func (s *PurchaseOrderService) Approve(ctx context.Context, id int, approvedBy int) error {
	// Get existing PO
	po, err := s.poRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("purchase order not found: %w", err)
	}

	// Check if PO can be approved
	if po.Status != inventory.POStatusDraft {
		return fmt.Errorf("only draft purchase orders can be approved")
	}

	// Update approval
	return s.poRepo.UpdateApproval(ctx, id, approvedBy)
}

// Cancel cancels a purchase order
func (s *PurchaseOrderService) Cancel(ctx context.Context, id int, userID int, reason string) error {
	// Get existing PO
	po, err := s.poRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("purchase order not found: %w", err)
	}

	// Check if PO can be cancelled
	if !po.CanBeCancelled() {
		return fmt.Errorf("purchase order cannot be cancelled in current status: %s", po.Status)
	}

	// Update status to cancelled
	return s.poRepo.UpdateStatus(ctx, id, inventory.POStatusCancelled)
}

// Acknowledge acknowledges receipt of PO by supplier
func (s *PurchaseOrderService) Acknowledge(ctx context.Context, id int) error {
	// Get existing PO
	po, err := s.poRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("purchase order not found: %w", err)
	}

	// Check if PO can be acknowledged
	if po.Status != inventory.POStatusSent {
		return fmt.Errorf("only sent purchase orders can be acknowledged")
	}

	// Update status to acknowledged
	return s.poRepo.UpdateStatus(ctx, id, inventory.POStatusAcknowledged)
}

// ReceiveGoods marks goods as received (partial or complete)
func (s *PurchaseOrderService) ReceiveGoods(ctx context.Context, id int, isPartial bool) error {
	// Get existing PO
	po, err := s.poRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("purchase order not found: %w", err)
	}

	// Check if PO can receive goods
	if po.Status != inventory.POStatusSent && po.Status != inventory.POStatusAcknowledged && po.Status != inventory.POStatusPartialReceived {
		return fmt.Errorf("cannot receive goods for purchase order in status: %s", po.Status)
	}

	// Update status based on whether it's partial or complete
	var newStatus inventory.PurchaseOrderStatus
	if isPartial {
		newStatus = inventory.POStatusPartialReceived
	} else {
		newStatus = inventory.POStatusReceived
	}

	return s.poRepo.UpdateStatus(ctx, id, newStatus)
}

// Complete marks a purchase order as complete
func (s *PurchaseOrderService) Complete(ctx context.Context, id int) error {
	// Get existing PO
	po, err := s.poRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("purchase order not found: %w", err)
	}

	// Check if PO can be completed
	if po.Status != inventory.POStatusReceived {
		return fmt.Errorf("only received purchase orders can be completed")
	}

	// Update status to completed
	return s.poRepo.UpdateStatus(ctx, id, inventory.POStatusCompleted)
}

// GetPendingApproval retrieves purchase orders pending approval
func (s *PurchaseOrderService) GetPendingApproval(ctx context.Context, page, limit int) ([]inventory.PurchaseOrderPartListItem, int, error) {
	return s.poRepo.GetPendingApproval(ctx, page, limit)
}

// GetReadyToSend retrieves purchase orders ready to send
func (s *PurchaseOrderService) GetReadyToSend(ctx context.Context, page, limit int) ([]inventory.PurchaseOrderPartListItem, int, error) {
	return s.poRepo.GetReadyToSend(ctx, page, limit)
}

// GetBySupplier retrieves purchase orders by supplier
func (s *PurchaseOrderService) GetBySupplier(ctx context.Context, supplierID int, page, limit int) ([]inventory.PurchaseOrderPartListItem, int, error) {
	return s.poRepo.GetBySupplier(ctx, supplierID, page, limit)
}

// GetByStatus retrieves purchase orders by status
func (s *PurchaseOrderService) GetByStatus(ctx context.Context, status inventory.PurchaseOrderStatus, page, limit int) ([]inventory.PurchaseOrderPartListItem, int, error) {
	return s.poRepo.GetByStatus(ctx, status, page, limit)
}

// ValidatePONumber validates PO number format
func (s *PurchaseOrderService) ValidatePONumber(number string) bool {
	return s.codeGenerator.ValidatePONumber(number)
}

// GetWorkflowActions returns available actions for a purchase order based on its status
func (s *PurchaseOrderService) GetWorkflowActions(ctx context.Context, id int) ([]string, error) {
	po, err := s.poRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("purchase order not found: %w", err)
	}

	var actions []string

	switch po.Status {
	case inventory.POStatusDraft:
		actions = append(actions, "edit", "approve", "send", "cancel", "delete")
	case inventory.POStatusSent:
		actions = append(actions, "acknowledge", "cancel")
	case inventory.POStatusAcknowledged:
		actions = append(actions, "receive_goods", "cancel")
	case inventory.POStatusPartialReceived:
		actions = append(actions, "receive_goods")
	case inventory.POStatusReceived:
		actions = append(actions, "complete")
	case inventory.POStatusCompleted:
		actions = []string{"view"} // Read-only
	case inventory.POStatusCancelled:
		actions = []string{"view"} // Read-only
	}

	return actions, nil
}

// GetPOSummary returns summary statistics for purchase orders
func (s *PurchaseOrderService) GetPOSummary(ctx context.Context, supplierID *int, dateFrom, dateTo *time.Time) (map[string]interface{}, error) {
	// This would typically involve complex aggregation queries
	// For now, return a basic structure
	summary := map[string]interface{}{
		"total_orders":    0,
		"total_amount":    0.0,
		"pending_orders":  0,
		"completed_orders": 0,
		"cancelled_orders": 0,
	}

	return summary, nil
}