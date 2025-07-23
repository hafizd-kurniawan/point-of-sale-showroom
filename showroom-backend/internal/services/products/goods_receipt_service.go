package products

import (
	"context"
	"fmt"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// GoodsReceiptService handles business logic for goods receipts
type GoodsReceiptService struct {
	goodsReceiptRepo       interfaces.GoodsReceiptRepository
	goodsReceiptDetailRepo interfaces.GoodsReceiptDetailRepository
	poDetailRepo           interfaces.PurchaseOrderDetailRepository
	stockMovementRepo      interfaces.StockMovementRepository
	productRepo            interfaces.ProductSparePartRepository
}

// NewGoodsReceiptService creates a new goods receipt service
func NewGoodsReceiptService(
	goodsReceiptRepo interfaces.GoodsReceiptRepository,
	goodsReceiptDetailRepo interfaces.GoodsReceiptDetailRepository,
	poDetailRepo interfaces.PurchaseOrderDetailRepository,
	stockMovementRepo interfaces.StockMovementRepository,
	productRepo interfaces.ProductSparePartRepository,
) *GoodsReceiptService {
	return &GoodsReceiptService{
		goodsReceiptRepo:       goodsReceiptRepo,
		goodsReceiptDetailRepo: goodsReceiptDetailRepo,
		poDetailRepo:           poDetailRepo,
		stockMovementRepo:      stockMovementRepo,
		productRepo:            productRepo,
	}
}

// CreateGoodsReceipt creates a new goods receipt
func (s *GoodsReceiptService) CreateGoodsReceipt(ctx context.Context, req *products.GoodsReceiptCreateRequest, receivedBy int) (*products.GoodsReceipt, error) {
	// Generate receipt number
	receiptNumber, err := s.goodsReceiptRepo.GenerateNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate receipt number: %w", err)
	}

	// Create goods receipt model
	receipt := &products.GoodsReceipt{
		POID:                  req.POID,
		ReceiptNumber:         receiptNumber,
		ReceiptDate:           req.ReceiptDate,
		ReceivedBy:            receivedBy,
		SupplierDeliveryNote:  req.SupplierDeliveryNote,
		SupplierInvoiceNumber: req.SupplierInvoiceNumber,
		ReceiptNotes:          req.ReceiptNotes,
		ReceiptDocumentsJSON:  req.ReceiptDocumentsJSON,
	}

	// Create the receipt
	createdReceipt, err := s.goodsReceiptRepo.Create(ctx, receipt)
	if err != nil {
		return nil, fmt.Errorf("failed to create goods receipt: %w", err)
	}

	return createdReceipt, nil
}

// GetGoodsReceipt retrieves a goods receipt by ID
func (s *GoodsReceiptService) GetGoodsReceipt(ctx context.Context, id int) (*products.GoodsReceipt, error) {
	receipt, err := s.goodsReceiptRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get goods receipt: %w", err)
	}

	return receipt, nil
}

// GetGoodsReceiptByNumber retrieves a goods receipt by receipt number
func (s *GoodsReceiptService) GetGoodsReceiptByNumber(ctx context.Context, number string) (*products.GoodsReceipt, error) {
	receipt, err := s.goodsReceiptRepo.GetByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get goods receipt by number: %w", err)
	}

	return receipt, nil
}

// UpdateGoodsReceipt updates a goods receipt
func (s *GoodsReceiptService) UpdateGoodsReceipt(ctx context.Context, id int, req *products.GoodsReceiptUpdateRequest) (*products.GoodsReceipt, error) {
	// Get existing receipt
	existing, err := s.goodsReceiptRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing receipt: %w", err)
	}

	// Update fields if provided
	if req.ReceiptDate != nil {
		existing.ReceiptDate = *req.ReceiptDate
	}
	if req.SupplierDeliveryNote != nil {
		existing.SupplierDeliveryNote = req.SupplierDeliveryNote
	}
	if req.SupplierInvoiceNumber != nil {
		existing.SupplierInvoiceNumber = req.SupplierInvoiceNumber
	}
	if req.ReceiptNotes != nil {
		existing.ReceiptNotes = req.ReceiptNotes
	}
	if req.DiscrepancyNotes != nil {
		existing.DiscrepancyNotes = req.DiscrepancyNotes
	}
	if req.ReceiptDocumentsJSON != nil {
		existing.ReceiptDocumentsJSON = req.ReceiptDocumentsJSON
	}

	// Update the receipt
	updatedReceipt, err := s.goodsReceiptRepo.Update(ctx, id, existing)
	if err != nil {
		return nil, fmt.Errorf("failed to update goods receipt: %w", err)
	}

	return updatedReceipt, nil
}

// DeleteGoodsReceipt deletes a goods receipt
func (s *GoodsReceiptService) DeleteGoodsReceipt(ctx context.Context, id int) error {
	err := s.goodsReceiptRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete goods receipt: %w", err)
	}

	return nil
}

// ListGoodsReceipts retrieves goods receipts with pagination
func (s *GoodsReceiptService) ListGoodsReceipts(ctx context.Context, params *products.GoodsReceiptFilterParams) (*common.PaginatedResponse, error) {
	receipts, err := s.goodsReceiptRepo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list goods receipts: %w", err)
	}

	return receipts, nil
}

// GetGoodsReceiptsByPO retrieves goods receipts for a specific PO
func (s *GoodsReceiptService) GetGoodsReceiptsByPO(ctx context.Context, poID int, params *products.GoodsReceiptFilterParams) (*common.PaginatedResponse, error) {
	receipts, err := s.goodsReceiptRepo.GetByPOID(ctx, poID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get goods receipts by PO: %w", err)
	}

	return receipts, nil
}

// AddReceiptDetail adds a detail line to a goods receipt
func (s *GoodsReceiptService) AddReceiptDetail(ctx context.Context, receiptID int, req *products.GoodsReceiptDetailCreateRequest) (*products.GoodsReceiptDetail, error) {
	// Validate receipt exists
	receipt, err := s.goodsReceiptRepo.GetByID(ctx, receiptID)
	if err != nil {
		return nil, fmt.Errorf("receipt not found: %w", err)
	}

	// Validate PO detail exists and belongs to the same PO
	poDetail, err := s.poDetailRepo.GetByID(ctx, req.PODetailID)
	if err != nil {
		return nil, fmt.Errorf("PO detail not found: %w", err)
	}

	if poDetail.POID != receipt.POID {
		return nil, fmt.Errorf("PO detail does not belong to the same PO as the receipt")
	}

	// Validate product matches
	if poDetail.ProductID != req.ProductID {
		return nil, fmt.Errorf("product ID does not match PO detail")
	}

	// Check if quantity received doesn't exceed what's remaining
	remainingQuantity := poDetail.QuantityOrdered - poDetail.QuantityReceived
	if req.QuantityReceived > remainingQuantity {
		return nil, fmt.Errorf("received quantity exceeds remaining quantity: received %d, remaining %d", req.QuantityReceived, remainingQuantity)
	}

	// Validate quantities
	if req.QuantityAccepted+req.QuantityRejected != req.QuantityReceived {
		return nil, fmt.Errorf("accepted + rejected quantities must equal received quantity")
	}

	// Create goods receipt detail model
	detail := &products.GoodsReceiptDetail{
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

	// Create the detail
	createdDetail, err := s.goodsReceiptDetailRepo.Create(ctx, detail)
	if err != nil {
		return nil, fmt.Errorf("failed to create receipt detail: %w", err)
	}

	// Update PO detail quantities
	err = s.poDetailRepo.UpdateQuantityReceived(ctx, req.PODetailID, req.QuantityReceived)
	if err != nil {
		return nil, fmt.Errorf("failed to update PO detail quantities: %w", err)
	}

	return createdDetail, nil
}

// GetReceiptDetails retrieves all details for a goods receipt
func (s *GoodsReceiptService) GetReceiptDetails(ctx context.Context, receiptID int) ([]products.GoodsReceiptDetail, error) {
	details, err := s.goodsReceiptDetailRepo.GetByReceiptID(ctx, receiptID)
	if err != nil {
		return nil, fmt.Errorf("failed to get receipt details: %w", err)
	}

	return details, nil
}

// ProcessGoodsReceipt processes a goods receipt and updates stock
func (s *GoodsReceiptService) ProcessGoodsReceipt(ctx context.Context, receiptID int, processedBy int) error {
	// Get receipt details
	details, err := s.goodsReceiptDetailRepo.GetByReceiptID(ctx, receiptID)
	if err != nil {
		return fmt.Errorf("failed to get receipt details: %w", err)
	}

	if len(details) == 0 {
		return fmt.Errorf("cannot process receipt without details")
	}

	var totalValue float64
	hasDiscrepancy := false

	// Process each detail and create stock movements
	for _, detail := range details {
		// Only create stock movement for accepted quantities
		if detail.QuantityAccepted > 0 {
			err = s.stockMovementRepo.CreateMovementForReceipt(
				ctx,
				detail.ProductID,
				detail.QuantityAccepted,
				detail.UnitCost,
				receiptID,
				processedBy,
			)
			if err != nil {
				return fmt.Errorf("failed to create stock movement for product %d: %w", detail.ProductID, err)
			}
		}

		// Check for discrepancies
		if detail.QuantityRejected > 0 || detail.ConditionReceived != "good" {
			hasDiscrepancy = true
		}

		totalValue += float64(detail.QuantityAccepted) * detail.UnitCost
	}

	// Update receipt status and total value
	var status products.ReceiptStatus
	if hasDiscrepancy {
		status = products.ReceiptStatusWithDiscrepancy
	} else {
		status = products.ReceiptStatusComplete
	}

	err = s.goodsReceiptRepo.UpdateStatus(ctx, receiptID, status)
	if err != nil {
		return fmt.Errorf("failed to update receipt status: %w", err)
	}

	// Update total received value
	receipt, err := s.goodsReceiptRepo.GetByID(ctx, receiptID)
	if err != nil {
		return fmt.Errorf("failed to get receipt: %w", err)
	}

	receipt.TotalReceivedValue = totalValue
	_, err = s.goodsReceiptRepo.Update(ctx, receiptID, receipt)
	if err != nil {
		return fmt.Errorf("failed to update receipt total value: %w", err)
	}

	return nil
}

// GetPendingReceiptItems gets items pending receipt for a PO
func (s *GoodsReceiptService) GetPendingReceiptItems(ctx context.Context, poID int) ([]products.PurchaseOrderDetail, error) {
	items, err := s.poDetailRepo.GetPendingReceiptItems(ctx, poID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending receipt items: %w", err)
	}

	return items, nil
}

// BulkReceiveItems creates receipt details for multiple items at once
func (s *GoodsReceiptService) BulkReceiveItems(ctx context.Context, receiptID int, details []products.GoodsReceiptDetailCreateRequest) error {
	// Validate receipt exists
	receipt, err := s.goodsReceiptRepo.GetByID(ctx, receiptID)
	if err != nil {
		return fmt.Errorf("receipt not found: %w", err)
	}

	// Convert requests to detail models
	var receiptDetails []products.GoodsReceiptDetail
	for _, req := range details {
		// Validate PO detail exists and belongs to the same PO
		poDetail, err := s.poDetailRepo.GetByID(ctx, req.PODetailID)
		if err != nil {
			return fmt.Errorf("PO detail %d not found: %w", req.PODetailID, err)
		}

		if poDetail.POID != receipt.POID {
			return fmt.Errorf("PO detail %d does not belong to the same PO as the receipt", req.PODetailID)
		}

		// Validate quantities
		if req.QuantityAccepted+req.QuantityRejected != req.QuantityReceived {
			return fmt.Errorf("accepted + rejected quantities must equal received quantity for PO detail %d", req.PODetailID)
		}

		detail := products.GoodsReceiptDetail{
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

		receiptDetails = append(receiptDetails, detail)
	}

	// Create all details
	err = s.goodsReceiptDetailRepo.BulkCreate(ctx, receiptDetails)
	if err != nil {
		return fmt.Errorf("failed to create receipt details: %w", err)
	}

	// Update PO detail quantities for each item
	for _, req := range details {
		err = s.poDetailRepo.UpdateQuantityReceived(ctx, req.PODetailID, req.QuantityReceived)
		if err != nil {
			return fmt.Errorf("failed to update PO detail %d quantities: %w", req.PODetailID, err)
		}
	}

	return nil
}