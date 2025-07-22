package inventory

import (
	"context"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// GoodsReceiptService handles business logic for goods receipts
type GoodsReceiptService struct {
	goodsReceiptRepo       interfaces.GoodsReceiptRepository
	goodsReceiptDetailRepo interfaces.GoodsReceiptDetailRepository
	purchaseOrderRepo      interfaces.PurchaseOrderRepository
	purchaseOrderDetailRepo interfaces.PurchaseOrderDetailRepository
	stockMovementRepo      interfaces.StockMovementRepository
	productRepo            interfaces.ProductSparePartRepository
	codeGenerator          *inventory.CodeGenerator
}

// NewGoodsReceiptService creates a new goods receipt service
func NewGoodsReceiptService(
	goodsReceiptRepo interfaces.GoodsReceiptRepository,
	goodsReceiptDetailRepo interfaces.GoodsReceiptDetailRepository,
	purchaseOrderRepo interfaces.PurchaseOrderRepository,
	purchaseOrderDetailRepo interfaces.PurchaseOrderDetailRepository,
	stockMovementRepo interfaces.StockMovementRepository,
	productRepo interfaces.ProductSparePartRepository,
	codeGenerator *inventory.CodeGenerator,
) *GoodsReceiptService {
	return &GoodsReceiptService{
		goodsReceiptRepo:        goodsReceiptRepo,
		goodsReceiptDetailRepo:  goodsReceiptDetailRepo,
		purchaseOrderRepo:       purchaseOrderRepo,
		purchaseOrderDetailRepo: purchaseOrderDetailRepo,
		stockMovementRepo:       stockMovementRepo,
		productRepo:             productRepo,
		codeGenerator:           codeGenerator,
	}
}

// CreateGoodsReceipt creates a new goods receipt with details
func (s *GoodsReceiptService) CreateGoodsReceipt(ctx context.Context, req *inventory.GoodsReceiptCreateRequest, userID int) (*inventory.GoodsReceipt, error) {
	// Validate purchase order exists
	po, err := s.purchaseOrderRepo.GetByID(ctx, req.PoID)
	if err != nil {
		return nil, fmt.Errorf("purchase order not found: %w", err)
	}

	if po.Status != inventory.POStatusSent && po.Status != inventory.POStatusAcknowledged {
		return nil, fmt.Errorf("purchase order is not in a receivable state")
	}

	// Generate receipt number  
	getLastReceiptID := func() (int, error) {
		// Get the last receipt ID from database
		return s.goodsReceiptRepo.GetLastReceiptID(ctx)
	}
	
	receiptNumber, err := s.codeGenerator.GetNextReceiptNumber(getLastReceiptID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate receipt number: %w", err)
	}

	// Calculate total received value
	var totalValue float64
	for _, detail := range req.Details {
		totalValue += detail.UnitCost * float64(detail.QuantityAccepted)
	}

	// Create goods receipt
	receipt := &inventory.GoodsReceipt{
		PoID:                  req.PoID,
		ReceiptNumber:         receiptNumber,
		ReceiptDate:           time.Now(),
		ReceivedBy:            userID,
		SupplierDeliveryNote:  req.SupplierDeliveryNote,
		SupplierInvoiceNumber: req.SupplierInvoiceNumber,
		TotalReceivedValue:    totalValue,
		ReceiptStatus:         inventory.ReceiptStatusComplete,
		ReceiptNotes:          req.ReceiptNotes,
		DiscrepancyNotes:      req.DiscrepancyNotes,
		ReceiptDocumentsJSON:  req.ReceiptDocumentsJSON,
	}

	// Check for discrepancies
	hasDiscrepancies := false
	for _, detail := range req.Details {
		if detail.QuantityRejected > 0 || detail.ConditionReceived != inventory.ConditionGood {
			hasDiscrepancies = true
			break
		}
	}

	if hasDiscrepancies {
		receipt.ReceiptStatus = inventory.ReceiptStatusWithDiscrepancy
	}

	// Save receipt
	createdReceipt, err := s.goodsReceiptRepo.Create(ctx, receipt)
	if err != nil {
		return nil, fmt.Errorf("failed to create goods receipt: %w", err)
	}

	// Create receipt details and update stock
	var details []inventory.GoodsReceiptDetail
	for _, detailReq := range req.Details {
		detail := inventory.GoodsReceiptDetail{
			ReceiptID:         createdReceipt.ReceiptID,
			PoDetailID:        detailReq.PoDetailID,
			ProductID:         detailReq.ProductID,
			QuantityReceived:  detailReq.QuantityReceived,
			QuantityAccepted:  detailReq.QuantityAccepted,
			QuantityRejected:  detailReq.QuantityRejected,
			UnitCost:          detailReq.UnitCost,
			TotalCost:         detailReq.UnitCost * float64(detailReq.QuantityReceived),
			ConditionReceived: detailReq.ConditionReceived,
			InspectionNotes:   detailReq.InspectionNotes,
			RejectionReason:   detailReq.RejectionReason,
			ExpiryDate:        detailReq.ExpiryDate,
			BatchNumber:       detailReq.BatchNumber,
			SerialNumbersJSON: detailReq.SerialNumbersJSON,
		}

		// Create stock movement for accepted items
		if detail.QuantityAccepted > 0 {
			referenceID := createdReceipt.ReceiptID
			movementReason := "Goods Receipt"
			notes := fmt.Sprintf("Receipt #%s", receiptNumber)
			
			movement := &inventory.StockMovement{
				ProductID:       detail.ProductID,
				MovementType:    inventory.MovementTypeIn,
				ReferenceType:   inventory.ReferenceTypePurchase,
				ReferenceID:     &referenceID,
				QuantityMoved:   detail.QuantityAccepted,
				UnitCost:        detail.UnitCost,
				MovementDate:    time.Now(),
				ProcessedBy:     userID,
				MovementReason:  &movementReason,
				Notes:           &notes,
			}

			// Get current stock to calculate before/after quantities
			product, err := s.productRepo.GetByID(ctx, detail.ProductID)
			if err != nil {
				return nil, fmt.Errorf("failed to get product: %w", err)
			}

			movement.QuantityBefore = product.StockQuantity
			movement.QuantityAfter = product.StockQuantity + detail.QuantityAccepted
			movement.TotalValue = float64(detail.QuantityAccepted) * detail.UnitCost

			// Create stock movement
			_, err = s.stockMovementRepo.Create(ctx, movement)
			if err != nil {
				return nil, fmt.Errorf("failed to create stock movement: %w", err)
			}

			// Update product stock
			err = s.productRepo.UpdateStock(ctx, detail.ProductID, movement.QuantityAfter)
			if err != nil {
				return nil, fmt.Errorf("failed to update product stock: %w", err)
			}
		}

		// Update purchase order detail quantity received
		err = s.purchaseOrderDetailRepo.UpdateQuantityReceived(ctx, detail.PoDetailID, detail.QuantityReceived)
		if err != nil {
			return nil, fmt.Errorf("failed to update PO detail: %w", err)
		}

		details = append(details, detail)
	}

	// Create receipt details
	err = s.goodsReceiptDetailRepo.CreateBatch(ctx, details)
	if err != nil {
		return nil, fmt.Errorf("failed to create receipt details: %w", err)
	}

	// Update purchase order status
	err = s.updatePurchaseOrderStatus(ctx, req.PoID)
	if err != nil {
		return nil, fmt.Errorf("failed to update PO status: %w", err)
	}

	return createdReceipt, nil
}

// GetGoodsReceipt retrieves a goods receipt by ID
func (s *GoodsReceiptService) GetGoodsReceipt(ctx context.Context, id int) (*inventory.GoodsReceipt, error) {
	receipt, err := s.goodsReceiptRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get receipt details
	details, err := s.goodsReceiptDetailRepo.GetByReceiptID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get receipt details: %w", err)
	}

	receipt.Details = details
	return receipt, nil
}

// ListGoodsReceipts retrieves goods receipts with filtering
func (s *GoodsReceiptService) ListGoodsReceipts(ctx context.Context, params *inventory.GoodsReceiptFilterParams) ([]inventory.GoodsReceiptListItem, int, error) {
	return s.goodsReceiptRepo.List(ctx, params)
}

// UpdateGoodsReceipt updates a goods receipt
func (s *GoodsReceiptService) UpdateGoodsReceipt(ctx context.Context, id int, req *inventory.GoodsReceiptUpdateRequest) (*inventory.GoodsReceipt, error) {
	receipt, err := s.goodsReceiptRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.SupplierDeliveryNote != nil {
		receipt.SupplierDeliveryNote = req.SupplierDeliveryNote
	}
	if req.SupplierInvoiceNumber != nil {
		receipt.SupplierInvoiceNumber = req.SupplierInvoiceNumber
	}
	if req.ReceiptStatus != nil {
		receipt.ReceiptStatus = *req.ReceiptStatus
	}
	if req.ReceiptNotes != nil {
		receipt.ReceiptNotes = req.ReceiptNotes
	}
	if req.DiscrepancyNotes != nil {
		receipt.DiscrepancyNotes = req.DiscrepancyNotes
	}
	if req.ReceiptDocumentsJSON != nil {
		receipt.ReceiptDocumentsJSON = req.ReceiptDocumentsJSON
	}

	return s.goodsReceiptRepo.Update(ctx, id, receipt)
}

// DeleteGoodsReceipt deletes a goods receipt
func (s *GoodsReceiptService) DeleteGoodsReceipt(ctx context.Context, id int) error {
	// Check if receipt exists
	_, err := s.goodsReceiptRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Note: In a real implementation, you might want to:
	// 1. Check if deletion is allowed (e.g., not already processed)
	// 2. Reverse stock movements
	// 3. Update purchase order details

	return s.goodsReceiptRepo.Delete(ctx, id)
}

// updatePurchaseOrderStatus updates the PO status based on received quantities
func (s *GoodsReceiptService) updatePurchaseOrderStatus(ctx context.Context, poID int) error {
	details, err := s.purchaseOrderDetailRepo.GetByPOID(ctx, poID)
	if err != nil {
		return err
	}

	allReceived := true
	anyReceived := false

	for _, detail := range details {
		if detail.QuantityReceived < detail.QuantityOrdered {
			allReceived = false
		}
		if detail.QuantityReceived > 0 {
			anyReceived = true
		}
	}

	var newStatus inventory.PurchaseOrderStatus
	if allReceived {
		newStatus = inventory.POStatusReceived
	} else if anyReceived {
		newStatus = inventory.POStatusPartialReceived
	} else {
		return nil // No change needed
	}

	return s.purchaseOrderRepo.UpdateStatus(ctx, poID, newStatus)
}