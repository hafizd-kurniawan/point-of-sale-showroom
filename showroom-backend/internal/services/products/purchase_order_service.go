package products

import (
	"context"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// PurchaseOrderService handles business logic for purchase orders
type PurchaseOrderService struct {
	poRepo         interfaces.PurchaseOrderPartsRepository
	poDetailRepo   interfaces.PurchaseOrderDetailRepository
	productRepo    interfaces.ProductSparePartRepository
	receiptRepo    interfaces.GoodsReceiptRepository
	stockRepo      interfaces.StockMovementRepository
}

// NewPurchaseOrderService creates a new purchase order service
func NewPurchaseOrderService(
	poRepo interfaces.PurchaseOrderPartsRepository,
	poDetailRepo interfaces.PurchaseOrderDetailRepository,
	productRepo interfaces.ProductSparePartRepository,
	receiptRepo interfaces.GoodsReceiptRepository,
	stockRepo interfaces.StockMovementRepository,
) *PurchaseOrderService {
	return &PurchaseOrderService{
		poRepo:       poRepo,
		poDetailRepo: poDetailRepo,
		productRepo:  productRepo,
		receiptRepo:  receiptRepo,
		stockRepo:    stockRepo,
	}
}

// CreatePurchaseOrder creates a new purchase order with auto-generated number
func (s *PurchaseOrderService) CreatePurchaseOrder(ctx context.Context, req *products.PurchaseOrderPartsCreateRequest, createdBy int) (*products.PurchaseOrderParts, error) {
	// Generate PO number
	poNumber, err := s.poRepo.GenerateNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PO number: %w", err)
	}

	// Create PO model
	po := &products.PurchaseOrderParts{
		PONumber:             poNumber,
		SupplierID:           req.SupplierID,
		PODate:               time.Now(),
		RequiredDate:         req.RequiredDate,
		ExpectedDeliveryDate: req.ExpectedDeliveryDate,
		POType:               req.POType,
		Subtotal:             0,
		TaxAmount:            0,
		DiscountAmount:       0,
		ShippingCost:         0,
		TotalAmount:          0,
		Status:               products.POStatusDraft,
		PaymentTerms:         req.PaymentTerms,
		CreatedBy:            createdBy,
		DeliveryAddress:      req.DeliveryAddress,
		PONotes:              req.PONotes,
		TermsAndConditions:   req.TermsAndConditions,
	}

	// Set payment due date based on terms
	po.SetPaymentDueDate()

	// Create PO
	createdPO, err := s.poRepo.Create(ctx, po)
	if err != nil {
		return nil, fmt.Errorf("failed to create purchase order: %w", err)
	}

	return createdPO, nil
}

// AddLineItem adds a line item to a purchase order
func (s *PurchaseOrderService) AddLineItem(ctx context.Context, poID int, req *products.PurchaseOrderDetailCreateRequest) (*products.PurchaseOrderDetail, error) {
	// Verify PO exists and is editable
	po, err := s.poRepo.GetByID(ctx, poID)
	if err != nil {
		return nil, fmt.Errorf("purchase order not found: %w", err)
	}

	if !po.CanEdit() {
		return nil, fmt.Errorf("purchase order cannot be edited in current status: %s", po.Status)
	}

	// Verify product exists
	product, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Create line item
	detail := &products.PurchaseOrderDetail{
		POID:            poID,
		ProductID:       req.ProductID,
		ItemDescription: req.ItemDescription,
		QuantityOrdered: req.QuantityOrdered,
		QuantityReceived: 0,
		QuantityPending: req.QuantityOrdered,
		UnitCost:        req.UnitCost,
		ExpectedDate:    req.ExpectedDate,
		LineStatus:      products.LineStatusPending,
		ItemNotes:       req.ItemNotes,
	}

	// Calculate total cost
	detail.CalculateTotalCost()

	// Create detail
	createdDetail, err := s.poDetailRepo.Create(ctx, detail)
	if err != nil {
		return nil, fmt.Errorf("failed to create line item: %w", err)
	}

	// Recalculate PO totals
	_, err = s.poRepo.CalculateTotals(ctx, poID)
	if err != nil {
		return nil, fmt.Errorf("failed to recalculate totals: %w", err)
	}

	// Use product name as description if not provided
	if createdDetail.ItemDescription == nil {
		createdDetail.ItemDescription = &product.ProductName
	}

	return createdDetail, nil
}

// GetPurchaseOrder retrieves a purchase order by ID
func (s *PurchaseOrderService) GetPurchaseOrder(ctx context.Context, id int) (*products.PurchaseOrderParts, error) {
	po, err := s.poRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchase order: %w", err)
	}
	return po, nil
}

// ListPurchaseOrders retrieves a paginated list of purchase orders
func (s *PurchaseOrderService) ListPurchaseOrders(ctx context.Context, params *products.PurchaseOrderPartsFilterParams) (*common.PaginatedResponse, error) {
	response, err := s.poRepo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list purchase orders: %w", err)
	}
	return response, nil
}

// GetPurchaseOrderDetails retrieves line items for a purchase order
func (s *PurchaseOrderService) GetPurchaseOrderDetails(ctx context.Context, poID int, params *products.PurchaseOrderDetailFilterParams) (*common.PaginatedResponse, error) {
	// Verify PO exists
	_, err := s.poRepo.GetByID(ctx, poID)
	if err != nil {
		return nil, fmt.Errorf("purchase order not found: %w", err)
	}

	response, err := s.poDetailRepo.GetByPOID(ctx, poID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchase order details: %w", err)
	}
	return response, nil
}

// ApprovePurchaseOrder approves a purchase order
func (s *PurchaseOrderService) ApprovePurchaseOrder(ctx context.Context, id int, approvedBy int) error {
	// Get PO
	po, err := s.poRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("purchase order not found: %w", err)
	}

	if !po.CanApprove() {
		return fmt.Errorf("purchase order cannot be approved in current status or already approved")
	}

	// Approve PO
	err = s.poRepo.Approve(ctx, id, approvedBy)
	if err != nil {
		return fmt.Errorf("failed to approve purchase order: %w", err)
	}

	return nil
}

// SendPurchaseOrder sends a purchase order to supplier
func (s *PurchaseOrderService) SendPurchaseOrder(ctx context.Context, id int) error {
	// Get PO
	po, err := s.poRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("purchase order not found: %w", err)
	}

	if !po.CanSend() {
		return fmt.Errorf("purchase order cannot be sent: not approved or wrong status")
	}

	// Update status to sent
	err = s.poRepo.UpdateStatus(ctx, id, products.POStatusSent)
	if err != nil {
		return fmt.Errorf("failed to update purchase order status: %w", err)
	}

	// TODO: Add integration with supplier system or email notification

	return nil
}

// CancelPurchaseOrder cancels a purchase order
func (s *PurchaseOrderService) CancelPurchaseOrder(ctx context.Context, id int) error {
	// Get PO
	po, err := s.poRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("purchase order not found: %w", err)
	}

	if !po.CanCancel() {
		return fmt.Errorf("purchase order cannot be cancelled in current status: %s", po.Status)
	}

	// Cancel PO
	err = s.poRepo.Cancel(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to cancel purchase order: %w", err)
	}

	return nil
}

// CreateGoodsReceipt creates a goods receipt for a purchase order
func (s *PurchaseOrderService) CreateGoodsReceipt(ctx context.Context, req *products.GoodsReceiptCreateRequest, receivedBy int) (*products.GoodsReceipt, error) {
	// Verify PO exists and can receive goods
	po, err := s.poRepo.GetByID(ctx, req.POID)
	if err != nil {
		return nil, fmt.Errorf("purchase order not found: %w", err)
	}

	if po.Status != products.POStatusSent && po.Status != products.POStatusAcknowledged && po.Status != products.POStatusPartialReceived {
		return nil, fmt.Errorf("purchase order status does not allow goods receipt: %s", po.Status)
	}

	// Generate receipt number
	receiptNumber, err := s.receiptRepo.GenerateNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate receipt number: %w", err)
	}

	// Create receipt
	receipt := &products.GoodsReceipt{
		POID:                  req.POID,
		ReceiptNumber:         receiptNumber,
		ReceiptDate:           req.ReceiptDate,
		ReceivedBy:            receivedBy,
		SupplierDeliveryNote:  req.SupplierDeliveryNote,
		SupplierInvoiceNumber: req.SupplierInvoiceNumber,
		TotalReceivedValue:    0,
		ReceiptStatus:         products.ReceiptStatusPartial,
		ReceiptNotes:          req.ReceiptNotes,
		ReceiptDocumentsJSON:  req.ReceiptDocumentsJSON,
	}

	createdReceipt, err := s.receiptRepo.Create(ctx, receipt)
	if err != nil {
		return nil, fmt.Errorf("failed to create goods receipt: %w", err)
	}

	return createdReceipt, nil
}

// ProcessGoodsReceiptItem processes individual items in a goods receipt
func (s *PurchaseOrderService) ProcessGoodsReceiptItem(ctx context.Context, receiptID int, req *products.GoodsReceiptDetailCreateRequest) error {
	// Verify receipt exists
	receipt, err := s.receiptRepo.GetByID(ctx, receiptID)
	if err != nil {
		return fmt.Errorf("goods receipt not found: %w", err)
	}

	// Verify PO detail exists
	poDetail, err := s.poDetailRepo.GetByID(ctx, req.PODetailID)
	if err != nil {
		return fmt.Errorf("purchase order detail not found: %w", err)
	}

	if poDetail.POID != receipt.POID {
		return fmt.Errorf("purchase order detail does not belong to the same PO as receipt")
	}

	// Validate quantities
	if req.QuantityReceived != (req.QuantityAccepted + req.QuantityRejected) {
		return fmt.Errorf("quantity received must equal accepted + rejected quantities")
	}

	// Create receipt detail
	receiptDetail := &products.GoodsReceiptDetail{
		ReceiptID:         receiptID,
		PODetailID:        req.PODetailID,
		ProductID:         req.ProductID,
		QuantityReceived:  req.QuantityReceived,
		QuantityAccepted:  req.QuantityAccepted,
		QuantityRejected:  req.QuantityRejected,
		UnitCost:          req.UnitCost,
		ConditionReceived: req.ConditionReceived,
		InspectionNotes:   req.InspectionNotes,
		RejectionReason:   req.RejectionReason,
		ExpiryDate:        req.ExpiryDate,
		BatchNumber:       req.BatchNumber,
		SerialNumbersJSON: req.SerialNumbersJSON,
	}

	receiptDetail.UpdateTotalCost()

	// TODO: Create receipt detail in repository (interface method needed)
	// err = s.receiptDetailRepo.Create(ctx, receiptDetail)

	// Update PO detail quantities
	poDetail.ReceiveQuantity(req.QuantityReceived)
	_, err = s.poDetailRepo.Update(ctx, req.PODetailID, poDetail)
	if err != nil {
		return fmt.Errorf("failed to update PO detail: %w", err)
	}

	// Update product stock if goods are accepted
	if req.QuantityAccepted > 0 {
		movement := &products.StockMovement{
			ProductID:     req.ProductID,
			MovementType:  products.MovementTypeIn,
			ReferenceType: products.ReferenceTypePurchase,
			ReferenceID:   receiptID,
			UnitCost:      req.UnitCost,
			MovementDate:  receipt.ReceiptDate,
			ProcessedBy:   receipt.ReceivedBy,
			LocationTo:    nil, // Could be set based on product location
		}

		reason := fmt.Sprintf("Goods receipt from PO %d", receipt.POID)
		movement.MovementReason = &reason

		err = s.productRepo.UpdateStockWithMovement(ctx, req.ProductID, req.QuantityAccepted, movement)
		if err != nil {
			return fmt.Errorf("failed to update product stock: %w", err)
		}
	}

	// Update PO status based on all line items
	err = s.updatePOStatusAfterReceipt(ctx, receipt.POID)
	if err != nil {
		return fmt.Errorf("failed to update PO status: %w", err)
	}

	return nil
}

// updatePOStatusAfterReceipt updates PO status based on receipt progress
func (s *PurchaseOrderService) updatePOStatusAfterReceipt(ctx context.Context, poID int) error {
	// Get all pending line items
	pendingItems, err := s.poDetailRepo.GetPendingReceiptItems(ctx, poID)
	if err != nil {
		return fmt.Errorf("failed to get pending items: %w", err)
	}

	var newStatus products.POStatus
	if len(pendingItems) == 0 {
		// All items received
		newStatus = products.POStatusReceived
	} else {
		// Some items still pending
		newStatus = products.POStatusPartialReceived
	}

	err = s.poRepo.UpdateStatus(ctx, poID, newStatus)
	if err != nil {
		return fmt.Errorf("failed to update PO status: %w", err)
	}

	return nil
}

// GetPendingApprovalPOs retrieves purchase orders pending approval
func (s *PurchaseOrderService) GetPendingApprovalPOs(ctx context.Context, params *products.PurchaseOrderPartsFilterParams) (*common.PaginatedResponse, error) {
	response, err := s.poRepo.GetPendingApproval(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending approval POs: %w", err)
	}
	return response, nil
}

// UpdatePurchaseOrder updates a purchase order (only if in draft status)
func (s *PurchaseOrderService) UpdatePurchaseOrder(ctx context.Context, id int, req *products.PurchaseOrderPartsUpdateRequest) (*products.PurchaseOrderParts, error) {
	// Get existing PO
	po, err := s.poRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("purchase order not found: %w", err)
	}

	if !po.CanEdit() {
		return nil, fmt.Errorf("purchase order cannot be edited in current status: %s", po.Status)
	}

	// Update fields
	if req.SupplierID != nil {
		po.SupplierID = *req.SupplierID
	}
	if req.RequiredDate != nil {
		po.RequiredDate = req.RequiredDate
	}
	if req.ExpectedDeliveryDate != nil {
		po.ExpectedDeliveryDate = req.ExpectedDeliveryDate
	}
	if req.POType != nil {
		po.POType = *req.POType
	}
	if req.TaxAmount != nil {
		po.TaxAmount = *req.TaxAmount
	}
	if req.DiscountAmount != nil {
		po.DiscountAmount = *req.DiscountAmount
	}
	if req.ShippingCost != nil {
		po.ShippingCost = *req.ShippingCost
	}
	if req.PaymentTerms != nil {
		po.PaymentTerms = *req.PaymentTerms
		po.SetPaymentDueDate() // Recalculate due date
	}
	if req.DeliveryAddress != nil {
		po.DeliveryAddress = req.DeliveryAddress
	}
	if req.PONotes != nil {
		po.PONotes = req.PONotes
	}
	if req.TermsAndConditions != nil {
		po.TermsAndConditions = req.TermsAndConditions
	}

	// Recalculate totals
	po.CalculateTotals()

	// Update PO
	updatedPO, err := s.poRepo.Update(ctx, id, po)
	if err != nil {
		return nil, fmt.Errorf("failed to update purchase order: %w", err)
	}

	return updatedPO, nil
}