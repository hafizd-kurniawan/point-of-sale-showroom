package interfaces

import (
	"context"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
)

// ProductSparePartRepository defines the interface for product spare part data operations
type ProductSparePartRepository interface {
	Create(ctx context.Context, product *products.ProductSparePart) (*products.ProductSparePart, error)
	GetByID(ctx context.Context, id int) (*products.ProductSparePart, error)
	GetByCode(ctx context.Context, code string) (*products.ProductSparePart, error)
	GetByBarcode(ctx context.Context, barcode string) (*products.ProductSparePart, error)
	Update(ctx context.Context, id int, product *products.ProductSparePart) (*products.ProductSparePart, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params *products.ProductSparePartFilterParams) (*common.PaginatedResponse, error)
	GetLowStockProducts(ctx context.Context, params *products.ProductSparePartFilterParams) (*common.PaginatedResponse, error)
	UpdateStock(ctx context.Context, id int, newQuantity int) error
	UpdateStockWithMovement(ctx context.Context, id int, quantityChange int, movementDetails *products.StockMovement) error
	GenerateCode(ctx context.Context) (string, error)
	IsCodeExists(ctx context.Context, code string) (bool, error)
	IsBarcodeExists(ctx context.Context, barcode string, excludeID int) (bool, error)
	GetByBrandID(ctx context.Context, brandID int, params *products.ProductSparePartFilterParams) (*common.PaginatedResponse, error)
	GetByCategoryID(ctx context.Context, categoryID int, params *products.ProductSparePartFilterParams) (*common.PaginatedResponse, error)
	Search(ctx context.Context, query string, params *products.ProductSparePartFilterParams) (*common.PaginatedResponse, error)
}

// PurchaseOrderPartsRepository defines the interface for purchase order data operations
type PurchaseOrderPartsRepository interface {
	Create(ctx context.Context, po *products.PurchaseOrderParts) (*products.PurchaseOrderParts, error)
	GetByID(ctx context.Context, id int) (*products.PurchaseOrderParts, error)
	GetByNumber(ctx context.Context, number string) (*products.PurchaseOrderParts, error)
	Update(ctx context.Context, id int, po *products.PurchaseOrderParts) (*products.PurchaseOrderParts, error)
	UpdateStatus(ctx context.Context, id int, status products.POStatus) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params *products.PurchaseOrderPartsFilterParams) (*common.PaginatedResponse, error)
	GetBySupplierID(ctx context.Context, supplierID int, params *products.PurchaseOrderPartsFilterParams) (*common.PaginatedResponse, error)
	GetByStatus(ctx context.Context, status products.POStatus, params *products.PurchaseOrderPartsFilterParams) (*common.PaginatedResponse, error)
	GetPendingApproval(ctx context.Context, params *products.PurchaseOrderPartsFilterParams) (*common.PaginatedResponse, error)
	Approve(ctx context.Context, id int, approvedBy int) error
	Cancel(ctx context.Context, id int) error
	GenerateNumber(ctx context.Context) (string, error)
	IsNumberExists(ctx context.Context, number string) (bool, error)
	CalculateTotals(ctx context.Context, id int) (*products.PurchaseOrderParts, error)
}

// PurchaseOrderDetailRepository defines the interface for purchase order detail data operations
type PurchaseOrderDetailRepository interface {
	Create(ctx context.Context, detail *products.PurchaseOrderDetail) (*products.PurchaseOrderDetail, error)
	GetByID(ctx context.Context, id int) (*products.PurchaseOrderDetail, error)
	Update(ctx context.Context, id int, detail *products.PurchaseOrderDetail) (*products.PurchaseOrderDetail, error)
	Delete(ctx context.Context, id int) error
	GetByPOID(ctx context.Context, poID int, params *products.PurchaseOrderDetailFilterParams) (*common.PaginatedResponse, error)
	GetByProductID(ctx context.Context, productID int, params *products.PurchaseOrderDetailFilterParams) (*common.PaginatedResponse, error)
	UpdateQuantityReceived(ctx context.Context, id int, quantityReceived int) error
	UpdateLineStatus(ctx context.Context, id int, status products.LineStatus) error
	GetPendingReceiptItems(ctx context.Context, poID int) ([]products.PurchaseOrderDetail, error)
	BulkCreate(ctx context.Context, details []products.PurchaseOrderDetail) error
	BulkUpdate(ctx context.Context, details []products.PurchaseOrderDetail) error
	CalculateSubtotal(ctx context.Context, poID int) (float64, error)
}

// GoodsReceiptRepository defines the interface for goods receipt data operations
type GoodsReceiptRepository interface {
	Create(ctx context.Context, receipt *products.GoodsReceipt) (*products.GoodsReceipt, error)
	GetByID(ctx context.Context, id int) (*products.GoodsReceipt, error)
	GetByNumber(ctx context.Context, number string) (*products.GoodsReceipt, error)
	Update(ctx context.Context, id int, receipt *products.GoodsReceipt) (*products.GoodsReceipt, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params *products.GoodsReceiptFilterParams) (*common.PaginatedResponse, error)
	GetByPOID(ctx context.Context, poID int, params *products.GoodsReceiptFilterParams) (*common.PaginatedResponse, error)
	GenerateNumber(ctx context.Context) (string, error)
	IsNumberExists(ctx context.Context, number string) (bool, error)
	UpdateStatus(ctx context.Context, id int, status products.ReceiptStatus) error
}

// GoodsReceiptDetailRepository defines the interface for goods receipt detail data operations
type GoodsReceiptDetailRepository interface {
	Create(ctx context.Context, detail *products.GoodsReceiptDetail) (*products.GoodsReceiptDetail, error)
	GetByID(ctx context.Context, id int) (*products.GoodsReceiptDetail, error)
	Update(ctx context.Context, id int, detail *products.GoodsReceiptDetail) (*products.GoodsReceiptDetail, error)
	Delete(ctx context.Context, id int) error
	GetByReceiptID(ctx context.Context, receiptID int) ([]products.GoodsReceiptDetail, error)
	GetByPODetailID(ctx context.Context, poDetailID int) ([]products.GoodsReceiptDetail, error)
	BulkCreate(ctx context.Context, details []products.GoodsReceiptDetail) error
	UpdateQuantities(ctx context.Context, id int, quantityAccepted, quantityRejected int) error
}

// StockMovementRepository defines the interface for stock movement data operations
type StockMovementRepository interface {
	Create(ctx context.Context, movement *products.StockMovement) (*products.StockMovement, error)
	GetByID(ctx context.Context, id int) (*products.StockMovement, error)
	List(ctx context.Context, params *products.StockMovementFilterParams) (*common.PaginatedResponse, error)
	GetByProductID(ctx context.Context, productID int, params *products.StockMovementFilterParams) (*common.PaginatedResponse, error)
	GetByReferenceID(ctx context.Context, referenceType products.ReferenceType, referenceID int) ([]products.StockMovement, error)
	CreateMovementForReceipt(ctx context.Context, productID int, quantity int, unitCost float64, receiptID int, processedBy int) error
	CreateMovementForAdjustment(ctx context.Context, productID int, quantityChange int, unitCost float64, adjustmentID int, processedBy int) error
	GetMovementHistory(ctx context.Context, productID int, limit int) ([]products.StockMovement, error)
	GetCurrentStock(ctx context.Context, productID int) (int, error)
	BulkCreateMovements(ctx context.Context, movements []products.StockMovement) error
}

// StockAdjustmentRepository defines the interface for stock adjustment data operations
type StockAdjustmentRepository interface {
	Create(ctx context.Context, adjustment *products.StockAdjustment) (*products.StockAdjustment, error)
	GetByID(ctx context.Context, id int) (*products.StockAdjustment, error)
	Update(ctx context.Context, id int, adjustment *products.StockAdjustment) (*products.StockAdjustment, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error)
	GetByProductID(ctx context.Context, productID int, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error)
	GetPendingApproval(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error)
	Approve(ctx context.Context, id int, approvedBy int) error
	GetVarianceReport(ctx context.Context, params *products.StockAdjustmentFilterParams) (*common.PaginatedResponse, error)
}

// SupplierPaymentRepository defines the interface for supplier payment data operations
type SupplierPaymentRepository interface {
	Create(ctx context.Context, payment *products.SupplierPayment) (*products.SupplierPayment, error)
	GetByID(ctx context.Context, id int) (*products.SupplierPayment, error)
	GetByNumber(ctx context.Context, number string) (*products.SupplierPayment, error)
	Update(ctx context.Context, id int, payment *products.SupplierPayment) (*products.SupplierPayment, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error)
	GetBySupplierID(ctx context.Context, supplierID int, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error)
	GetByPOID(ctx context.Context, poID int, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error)
	GetOverduePayments(ctx context.Context, params *products.SupplierPaymentFilterParams) (*common.PaginatedResponse, error)
	UpdatePaymentStatus(ctx context.Context, id int, status products.PaymentStatus) error
	AddPayment(ctx context.Context, id int, amount float64, method products.PaymentMethod, reference *string) error
	GenerateNumber(ctx context.Context) (string, error)
	IsNumberExists(ctx context.Context, number string) (bool, error)
	UpdateOverdueStatus(ctx context.Context) error
	GetPaymentSummary(ctx context.Context, supplierID *int) (map[string]interface{}, error)
}