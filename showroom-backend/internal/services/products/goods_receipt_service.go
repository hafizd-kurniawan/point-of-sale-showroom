package products

import (
	"context"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// GoodsReceiptService handles business logic for goods receipt management
type GoodsReceiptService struct {
	goodsReceiptRepo       interfaces.GoodsReceiptRepository
	goodsReceiptDetailRepo interfaces.GoodsReceiptDetailRepository
	purchaseOrderDetailRepo interfaces.PurchaseOrderDetailRepository
	stockMovementRepo      interfaces.StockMovementRepository
	productRepo            interfaces.ProductSparePartRepository
}

// NewGoodsReceiptService creates a new goods receipt service
func NewGoodsReceiptService(
	goodsReceiptRepo interfaces.GoodsReceiptRepository,
	goodsReceiptDetailRepo interfaces.GoodsReceiptDetailRepository,
	purchaseOrderDetailRepo interfaces.PurchaseOrderDetailRepository,
	stockMovementRepo interfaces.StockMovementRepository,
	productRepo interfaces.ProductSparePartRepository,
) *GoodsReceiptService {
	return &GoodsReceiptService{
		goodsReceiptRepo:       goodsReceiptRepo,
		goodsReceiptDetailRepo: goodsReceiptDetailRepo,
		purchaseOrderDetailRepo: purchaseOrderDetailRepo,
		stockMovementRepo:      stockMovementRepo,
		productRepo:            productRepo,
	}
}

// CreateGoodsReceipt creates a new goods receipt with business validation
func (s *GoodsReceiptService) CreateGoodsReceipt(ctx context.Context, req *products.GoodsReceiptCreateRequest, receivedBy int) (*products.GoodsReceipt, error) {
	// Validate PO exists and get pending items
	pendingItems, err := s.purchaseOrderDetailRepo.GetPendingReceiptItems(ctx, req.POID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending PO items: %w", err)
	}

	if len(pendingItems) == 0 {
		return nil, fmt.Errorf("no pending items found for purchase order")
	}

	// Generate receipt number
	receiptNumber, err := s.goodsReceiptRepo.GenerateNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate receipt number: %w", err)
	}

	// Create goods receipt record
	receipt := &products.GoodsReceipt{
		POID:                  req.POID,
		ReceiptNumber:         receiptNumber,
		ReceiptDate:           req.ReceiptDate,
		ReceivedBy:            receivedBy,
		SupplierDeliveryNote:  req.SupplierDeliveryNote,
		SupplierInvoiceNumber: req.SupplierInvoiceNumber,
		TotalReceivedValue:    0, // Will be calculated
		ReceiptStatus:         products.ReceiptStatusPartial,
		ReceiptNotes:          req.ReceiptNotes,
		ReceiptDocumentsJSON:  req.ReceiptDocumentsJSON,
	}

	createdReceipt, err := s.goodsReceiptRepo.Create(ctx, receipt)
	if err != nil {
		return nil, fmt.Errorf("failed to create goods receipt: %w", err)
	}

	return createdReceipt, nil
}

// AddGoodsReceiptDetails adds receipt details and processes stock movements
func (s *GoodsReceiptService) AddGoodsReceiptDetails(ctx context.Context, receiptID int, details []products.GoodsReceiptDetailCreateRequest) error {
	// Validate receipt exists
	receipt, err := s.goodsReceiptRepo.GetByID(ctx, receiptID)
	if err != nil {
		return fmt.Errorf("goods receipt not found: %w", err)
	}

	var receiptDetails []products.GoodsReceiptDetail
	var totalReceivedValue float64

	for _, detailReq := range details {
		// Validate PO detail exists
		poDetail, err := s.purchaseOrderDetailRepo.GetByID(ctx, detailReq.PODetailID)
		if err != nil {
			return fmt.Errorf("PO detail not found: %w", err)
		}

		// Validate product exists
		_, err = s.productRepo.GetByID(ctx, detailReq.ProductID)
		if err != nil {
			return fmt.Errorf("product not found: %w", err)
		}

		// Validate quantities
		if detailReq.QuantityReceived != (detailReq.QuantityAccepted + detailReq.QuantityRejected) {
			return fmt.Errorf("quantity received must equal accepted + rejected for PO detail %d", detailReq.PODetailID)
		}

		// Validate we don't receive more than ordered
		remainingToReceive := poDetail.QuantityOrdered - poDetail.QuantityReceived
		if detailReq.QuantityReceived > remainingToReceive {
			return fmt.Errorf("cannot receive more than pending quantity for PO detail %d", detailReq.PODetailID)
		}

		// Create receipt detail
		detail := &products.GoodsReceiptDetail{
			ReceiptID:         receiptID,
			PODetailID:        detailReq.PODetailID,
			ProductID:         detailReq.ProductID,
			QuantityReceived:  detailReq.QuantityReceived,
			QuantityAccepted:  detailReq.QuantityAccepted,
			QuantityRejected:  detailReq.QuantityRejected,
			UnitCost:          detailReq.UnitCost,
			ConditionReceived: detailReq.ConditionReceived,
			InspectionNotes:   detailReq.InspectionNotes,
			RejectionReason:   detailReq.RejectionReason,
			ExpiryDate:        detailReq.ExpiryDate,
			BatchNumber:       detailReq.BatchNumber,
			SerialNumbersJSON: detailReq.SerialNumbersJSON,
		}

		// Calculate total cost
		detail.UpdateTotalCost()
		totalReceivedValue += detail.TotalCost

		receiptDetails = append(receiptDetails, *detail)

		// Update PO detail quantities
		err = s.purchaseOrderDetailRepo.UpdateQuantityReceived(ctx, detailReq.PODetailID, detailReq.QuantityReceived)
		if err != nil {
			return fmt.Errorf("failed to update PO detail quantities: %w", err)
		}

		// Create stock movement for accepted quantity only
		if detailReq.QuantityAccepted > 0 {
			err = s.stockMovementRepo.CreateMovementForReceipt(ctx, detailReq.ProductID, detailReq.QuantityAccepted, detailReq.UnitCost, receiptID, receipt.ReceivedBy)
			if err != nil {
				return fmt.Errorf("failed to create stock movement: %w", err)
			}

			// Update product stock
			currentStock, err := s.stockMovementRepo.GetCurrentStock(ctx, detailReq.ProductID)
			if err != nil {
				return fmt.Errorf("failed to get current stock: %w", err)
			}

			err = s.productRepo.UpdateStock(ctx, detailReq.ProductID, currentStock+detailReq.QuantityAccepted)
			if err != nil {
				return fmt.Errorf("failed to update product stock: %w", err)
			}
		}
	}

	// Bulk create receipt details
	err = s.goodsReceiptDetailRepo.BulkCreate(ctx, receiptDetails)
	if err != nil {
		return fmt.Errorf("failed to create receipt details: %w", err)
	}

	// Update receipt total and status
	hasDiscrepancy := false
	for _, detail := range receiptDetails {
		if detail.HasDiscrepancy() {
			hasDiscrepancy = true
			break
		}
	}

	// Check if receipt is complete (all PO items received)
	pendingItems, err := s.purchaseOrderDetailRepo.GetPendingReceiptItems(ctx, receipt.POID)
	if err != nil {
		return fmt.Errorf("failed to check pending items: %w", err)
	}
	isComplete := len(pendingItems) == 0

	// Update receipt status and total
	receipt.TotalReceivedValue = totalReceivedValue
	receipt.UpdateStatus(hasDiscrepancy, isComplete)

	_, err = s.goodsReceiptRepo.Update(ctx, receiptID, receipt)
	if err != nil {
		return fmt.Errorf("failed to update receipt status: %w", err)
	}

	return nil
}

// GetGoodsReceipts retrieves goods receipts with pagination and filtering
func (s *GoodsReceiptService) GetGoodsReceipts(ctx context.Context, params *products.GoodsReceiptFilterParams) (*common.PaginatedResponse, error) {
	return s.goodsReceiptRepo.List(ctx, params)
}

// GetGoodsReceiptByID retrieves a goods receipt by ID with details
func (s *GoodsReceiptService) GetGoodsReceiptByID(ctx context.Context, id int) (*products.GoodsReceipt, []products.GoodsReceiptDetail, error) {
	receipt, err := s.goodsReceiptRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("goods receipt not found: %w", err)
	}

	details, err := s.goodsReceiptDetailRepo.GetByReceiptID(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get receipt details: %w", err)
	}

	return receipt, details, nil
}

// GetGoodsReceiptByNumber retrieves a goods receipt by receipt number
func (s *GoodsReceiptService) GetGoodsReceiptByNumber(ctx context.Context, number string) (*products.GoodsReceipt, error) {
	return s.goodsReceiptRepo.GetByNumber(ctx, number)
}

// GetGoodsReceiptsByPO retrieves goods receipts for a purchase order
func (s *GoodsReceiptService) GetGoodsReceiptsByPO(ctx context.Context, poID int, params *products.GoodsReceiptFilterParams) (*common.PaginatedResponse, error) {
	return s.goodsReceiptRepo.GetByPOID(ctx, poID, params)
}

// UpdateGoodsReceipt updates a goods receipt
func (s *GoodsReceiptService) UpdateGoodsReceipt(ctx context.Context, id int, req *products.GoodsReceiptUpdateRequest) (*products.GoodsReceipt, error) {
	// Get existing receipt
	receipt, err := s.goodsReceiptRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("goods receipt not found: %w", err)
	}

	// Update fields if provided
	if req.ReceiptDate != nil {
		receipt.ReceiptDate = *req.ReceiptDate
	}
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

// UpdateReceiptDetailQuantities updates quantities for a receipt detail
func (s *GoodsReceiptService) UpdateReceiptDetailQuantities(ctx context.Context, detailID int, quantityAccepted, quantityRejected int) error {
	// Get existing detail
	detail, err := s.goodsReceiptDetailRepo.GetByID(ctx, detailID)
	if err != nil {
		return fmt.Errorf("receipt detail not found: %w", err)
	}

	// Validate quantities
	if detail.QuantityReceived != (quantityAccepted + quantityRejected) {
		return fmt.Errorf("quantity received must equal accepted + rejected")
	}

	// Calculate quantity difference for stock adjustment
	quantityDiff := quantityAccepted - detail.QuantityAccepted

	// Update receipt detail
	err = s.goodsReceiptDetailRepo.UpdateQuantities(ctx, detailID, quantityAccepted, quantityRejected)
	if err != nil {
		return fmt.Errorf("failed to update receipt detail quantities: %w", err)
	}

	// Update stock if there's a difference
	if quantityDiff != 0 {
		// Get receipt for processed_by info
		receipt, err := s.goodsReceiptRepo.GetByID(ctx, detail.ReceiptID)
		if err != nil {
			return fmt.Errorf("failed to get receipt: %w", err)
		}

		// Create stock movement for the difference
		currentStock, err := s.stockMovementRepo.GetCurrentStock(ctx, detail.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get current stock: %w", err)
		}

		movementType := products.MovementTypeIn
		if quantityDiff < 0 {
			movementType = products.MovementTypeOut
		}

		movement := &products.StockMovement{
			ProductID:      detail.ProductID,
			MovementType:   movementType,
			ReferenceType:  products.ReferenceTypePurchase,
			ReferenceID:    detail.ReceiptID,
			QuantityBefore: currentStock,
			QuantityMoved:  abs(quantityDiff),
			QuantityAfter:  currentStock + quantityDiff,
			UnitCost:       detail.UnitCost,
			MovementDate:   time.Now(),
			ProcessedBy:    receipt.ReceivedBy,
			MovementReason: stringPtr("Receipt quantity adjustment"),
		}

		movement.CalculateTotalValue()

		_, err = s.stockMovementRepo.Create(ctx, movement)
		if err != nil {
			return fmt.Errorf("failed to create adjustment movement: %w", err)
		}

		// Update product stock
		err = s.productRepo.UpdateStock(ctx, detail.ProductID, currentStock+quantityDiff)
		if err != nil {
			return fmt.Errorf("failed to update product stock: %w", err)
		}
	}

	return nil
}

// DeleteGoodsReceipt deletes a goods receipt and reverses stock movements
func (s *GoodsReceiptService) DeleteGoodsReceipt(ctx context.Context, id int) error {
	// Get receipt and details
	receipt, details, err := s.GetGoodsReceiptByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get receipt: %w", err)
	}

	// Reverse stock movements and update PO details
	for _, detail := range details {
		if detail.QuantityAccepted > 0 {
			// Create reverse stock movement
			currentStock, err := s.stockMovementRepo.GetCurrentStock(ctx, detail.ProductID)
			if err != nil {
				return fmt.Errorf("failed to get current stock: %w", err)
			}

			if currentStock < detail.QuantityAccepted {
				return fmt.Errorf("insufficient stock to reverse receipt for product %d", detail.ProductID)
			}

			movement := &products.StockMovement{
				ProductID:      detail.ProductID,
				MovementType:   products.MovementTypeOut,
				ReferenceType:  products.ReferenceTypePurchase,
				ReferenceID:    id,
				QuantityBefore: currentStock,
				QuantityMoved:  detail.QuantityAccepted,
				QuantityAfter:  currentStock - detail.QuantityAccepted,
				UnitCost:       detail.UnitCost,
				MovementDate:   time.Now(),
				ProcessedBy:    receipt.ReceivedBy,
				MovementReason: stringPtr("Receipt deletion reversal"),
			}

			movement.CalculateTotalValue()

			_, err = s.stockMovementRepo.Create(ctx, movement)
			if err != nil {
				return fmt.Errorf("failed to create reversal movement: %w", err)
			}

			// Update product stock
			err = s.productRepo.UpdateStock(ctx, detail.ProductID, currentStock-detail.QuantityAccepted)
			if err != nil {
				return fmt.Errorf("failed to update product stock: %w", err)
			}
		}

		// Reverse PO detail quantity received
		err = s.purchaseOrderDetailRepo.UpdateQuantityReceived(ctx, detail.PODetailID, -detail.QuantityReceived)
		if err != nil {
			return fmt.Errorf("failed to reverse PO detail quantities: %w", err)
		}
	}

	// Delete the receipt (details will be cascade deleted)
	return s.goodsReceiptRepo.Delete(ctx, id)
}

// helper function to get absolute value
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}