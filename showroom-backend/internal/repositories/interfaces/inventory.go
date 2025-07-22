package interfaces

import (
	"context"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
)

// ProductSparePartRepository defines the interface for product spare part data operations
type ProductSparePartRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, product *inventory.ProductSparePart) (*inventory.ProductSparePart, error)
	GetByID(ctx context.Context, id int) (*inventory.ProductSparePart, error)
	GetByCode(ctx context.Context, code string) (*inventory.ProductSparePart, error)
	GetByBarcode(ctx context.Context, barcode string) (*inventory.ProductSparePart, error)
	Update(ctx context.Context, id int, product *inventory.ProductSparePart) (*inventory.ProductSparePart, error)
	Delete(ctx context.Context, id int) error
	
	// List and filtering operations
	List(ctx context.Context, params *inventory.ProductSparePartFilterParams) ([]inventory.ProductSparePartListItem, int, error)
	GetLowStockProducts(ctx context.Context, page, limit int) ([]inventory.ProductSparePartListItem, int, error)
	GetByBrand(ctx context.Context, brandID int, page, limit int) ([]inventory.ProductSparePartListItem, int, error)
	GetByCategory(ctx context.Context, categoryID int, page, limit int) ([]inventory.ProductSparePartListItem, int, error)
	GetByLocation(ctx context.Context, location string, page, limit int) ([]inventory.ProductSparePartListItem, int, error)
	
	// Search operations
	Search(ctx context.Context, query string, page, limit int) ([]inventory.ProductSparePartListItem, int, error)
	
	// Existence checks
	ExistsByCode(ctx context.Context, code string) (bool, error)
	ExistsByBarcode(ctx context.Context, barcode string) (bool, error)
	ExistsByCodeExcludingID(ctx context.Context, code string, excludeID int) (bool, error)
	ExistsByBarcodeExcludingID(ctx context.Context, barcode string, excludeID int) (bool, error)
	
	// Code generation helpers
	GetLastProductID(ctx context.Context) (int, error)
	
	// Stock operations
	UpdateStock(ctx context.Context, productID int, newQuantity int) error
	GetStockQuantity(ctx context.Context, productID int) (int, error)
	
	// Inventory reporting
	GetInventoryValue(ctx context.Context) (float64, error)
	GetProductsByPriceRange(ctx context.Context, minPrice, maxPrice float64, page, limit int) ([]inventory.ProductSparePartListItem, int, error)
}

// PurchaseOrderRepository defines the interface for purchase order data operations
type PurchaseOrderRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, po *inventory.PurchaseOrderPart) (*inventory.PurchaseOrderPart, error)
	GetByID(ctx context.Context, id int) (*inventory.PurchaseOrderPart, error)
	GetByNumber(ctx context.Context, number string) (*inventory.PurchaseOrderPart, error)
	Update(ctx context.Context, id int, po *inventory.PurchaseOrderPart) (*inventory.PurchaseOrderPart, error)
	Delete(ctx context.Context, id int) error
	
	// List and filtering operations
	List(ctx context.Context, params *inventory.PurchaseOrderPartFilterParams) ([]inventory.PurchaseOrderPartListItem, int, error)
	GetBySupplier(ctx context.Context, supplierID int, page, limit int) ([]inventory.PurchaseOrderPartListItem, int, error)
	GetByStatus(ctx context.Context, status inventory.PurchaseOrderStatus, page, limit int) ([]inventory.PurchaseOrderPartListItem, int, error)
	GetByDateRange(ctx context.Context, startDate, endDate string, page, limit int) ([]inventory.PurchaseOrderPartListItem, int, error)
	
	// Search operations
	Search(ctx context.Context, query string, page, limit int) ([]inventory.PurchaseOrderPartListItem, int, error)
	
	// Existence checks
	ExistsByNumber(ctx context.Context, number string) (bool, error)
	ExistsByNumberExcludingID(ctx context.Context, number string, excludeID int) (bool, error)
	
	// Code generation helpers
	GetLastPOID(ctx context.Context) (int, error)
	
	// Status operations
	UpdateStatus(ctx context.Context, id int, status inventory.PurchaseOrderStatus) error
	UpdateApproval(ctx context.Context, id int, approvedBy int) error
	
	// Detail operations
	GetWithDetails(ctx context.Context, id int) (*inventory.PurchaseOrderPart, error)
	
	// Workflow operations
	GetPendingApproval(ctx context.Context, page, limit int) ([]inventory.PurchaseOrderPartListItem, int, error)
	GetReadyToSend(ctx context.Context, page, limit int) ([]inventory.PurchaseOrderPartListItem, int, error)
}

// PurchaseOrderDetailRepository defines the interface for purchase order detail operations
type PurchaseOrderDetailRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, detail *inventory.PurchaseOrderDetail) (*inventory.PurchaseOrderDetail, error)
	GetByID(ctx context.Context, id int) (*inventory.PurchaseOrderDetail, error)
	GetByPOID(ctx context.Context, poID int) ([]inventory.PurchaseOrderDetail, error)
	Update(ctx context.Context, id int, detail *inventory.PurchaseOrderDetail) (*inventory.PurchaseOrderDetail, error)
	Delete(ctx context.Context, id int) error
	
	// Batch operations
	CreateBatch(ctx context.Context, details []inventory.PurchaseOrderDetail) error
	UpdateQuantityReceived(ctx context.Context, id int, quantityReceived int) error
	UpdateLineStatus(ctx context.Context, id int, status inventory.LineStatus) error
	
	// Reporting
	GetPendingItems(ctx context.Context, page, limit int) ([]inventory.PurchaseOrderDetail, int, error)
	GetOverdueItems(ctx context.Context, page, limit int) ([]inventory.PurchaseOrderDetail, int, error)
}

// GoodsReceiptRepository defines the interface for goods receipt data operations
type GoodsReceiptRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, receipt *inventory.GoodsReceipt) (*inventory.GoodsReceipt, error)
	GetByID(ctx context.Context, id int) (*inventory.GoodsReceipt, error)
	GetByNumber(ctx context.Context, number string) (*inventory.GoodsReceipt, error)
	GetByPOID(ctx context.Context, poID int) ([]inventory.GoodsReceipt, error)
	Update(ctx context.Context, id int, receipt *inventory.GoodsReceipt) (*inventory.GoodsReceipt, error)
	Delete(ctx context.Context, id int) error
	
	// List and filtering operations
	List(ctx context.Context, params *inventory.GoodsReceiptFilterParams) ([]inventory.GoodsReceiptListItem, int, error)
	GetByStatus(ctx context.Context, status inventory.ReceiptStatus, page, limit int) ([]inventory.GoodsReceiptListItem, int, error)
	GetByDateRange(ctx context.Context, startDate, endDate string, page, limit int) ([]inventory.GoodsReceiptListItem, int, error)
	
	// Search operations
	Search(ctx context.Context, query string, page, limit int) ([]inventory.GoodsReceiptListItem, int, error)
	
	// Existence checks
	ExistsByNumber(ctx context.Context, number string) (bool, error)
	ExistsByNumberExcludingID(ctx context.Context, number string, excludeID int) (bool, error)
	
	// Code generation helpers
	GetLastReceiptID(ctx context.Context) (int, error)
	
	// Detail operations
	GetWithDetails(ctx context.Context, id int) (*inventory.GoodsReceipt, error)
	
	// Status operations
	UpdateStatus(ctx context.Context, id int, status inventory.ReceiptStatus) error
}

// GoodsReceiptDetailRepository defines the interface for goods receipt detail operations
type GoodsReceiptDetailRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, detail *inventory.GoodsReceiptDetail) (*inventory.GoodsReceiptDetail, error)
	GetByID(ctx context.Context, id int) (*inventory.GoodsReceiptDetail, error)
	GetByReceiptID(ctx context.Context, receiptID int) ([]inventory.GoodsReceiptDetail, error)
	Update(ctx context.Context, id int, detail *inventory.GoodsReceiptDetail) (*inventory.GoodsReceiptDetail, error)
	Delete(ctx context.Context, id int) error
	
	// Batch operations
	CreateBatch(ctx context.Context, details []inventory.GoodsReceiptDetail) error
	
	// Quality control operations
	GetRejectedItems(ctx context.Context, page, limit int) ([]inventory.GoodsReceiptDetail, int, error)
	GetDamagedItems(ctx context.Context, page, limit int) ([]inventory.GoodsReceiptDetail, int, error)
}

// StockMovementRepository defines the interface for stock movement data operations
type StockMovementRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, movement *inventory.StockMovement) (*inventory.StockMovement, error)
	GetByID(ctx context.Context, id int) (*inventory.StockMovement, error)
	
	// List and filtering operations
	List(ctx context.Context, params *inventory.StockMovementFilterParams) ([]inventory.StockMovementListItem, int, error)
	GetByProduct(ctx context.Context, productID int, page, limit int) ([]inventory.StockMovementListItem, int, error)
	GetByType(ctx context.Context, movementType inventory.MovementType, page, limit int) ([]inventory.StockMovementListItem, int, error)
	GetByReference(ctx context.Context, refType inventory.ReferenceType, refID int, page, limit int) ([]inventory.StockMovementListItem, int, error)
	GetByDateRange(ctx context.Context, startDate, endDate string, page, limit int) ([]inventory.StockMovementListItem, int, error)
	
	// Audit operations
	GetAuditTrail(ctx context.Context, productID int, page, limit int) ([]inventory.StockMovementListItem, int, error)
	GetMovementsByUser(ctx context.Context, userID int, page, limit int) ([]inventory.StockMovementListItem, int, error)
	
	// Reporting operations
	GetMovementSummary(ctx context.Context, productID int, startDate, endDate string) (map[string]interface{}, error)
	GetValueMovements(ctx context.Context, startDate, endDate string) (float64, error)
}

// StockAdjustmentRepository defines the interface for stock adjustment data operations
type StockAdjustmentRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, adjustment *inventory.StockAdjustment) (*inventory.StockAdjustment, error)
	GetByID(ctx context.Context, id int) (*inventory.StockAdjustment, error)
	Update(ctx context.Context, id int, adjustment *inventory.StockAdjustment) (*inventory.StockAdjustment, error)
	Delete(ctx context.Context, id int) error
	
	// List and filtering operations
	List(ctx context.Context, params *inventory.StockAdjustmentFilterParams) ([]inventory.StockAdjustmentListItem, int, error)
	GetByProduct(ctx context.Context, productID int, page, limit int) ([]inventory.StockAdjustmentListItem, int, error)
	GetByType(ctx context.Context, adjustmentType inventory.AdjustmentType, page, limit int) ([]inventory.StockAdjustmentListItem, int, error)
	GetPendingApproval(ctx context.Context, page, limit int) ([]inventory.StockAdjustmentListItem, int, error)
	
	// Approval operations
	UpdateApproval(ctx context.Context, id int, approvedBy int) error
	
	// Reporting operations
	GetAdjustmentSummary(ctx context.Context, startDate, endDate string) (map[string]interface{}, error)
}

// SupplierPaymentRepository defines the interface for supplier payment data operations
type SupplierPaymentRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, payment *inventory.SupplierPayment) (*inventory.SupplierPayment, error)
	GetByID(ctx context.Context, id int) (*inventory.SupplierPayment, error)
	GetByNumber(ctx context.Context, number string) (*inventory.SupplierPayment, error)
	Update(ctx context.Context, id int, payment *inventory.SupplierPayment) (*inventory.SupplierPayment, error)
	Delete(ctx context.Context, id int) error
	
	// List and filtering operations
	List(ctx context.Context, params *inventory.SupplierPaymentFilterParams) ([]inventory.SupplierPaymentListItem, int, error)
	GetBySupplier(ctx context.Context, supplierID int, page, limit int) ([]inventory.SupplierPaymentListItem, int, error)
	GetByPO(ctx context.Context, poID int, page, limit int) ([]inventory.SupplierPaymentListItem, int, error)
	GetByStatus(ctx context.Context, status inventory.PaymentStatus, page, limit int) ([]inventory.SupplierPaymentListItem, int, error)
	GetOverduePayments(ctx context.Context, page, limit int) ([]inventory.SupplierPaymentListItem, int, error)
	
	// Search operations
	Search(ctx context.Context, query string, page, limit int) ([]inventory.SupplierPaymentListItem, int, error)
	
	// Existence checks
	ExistsByNumber(ctx context.Context, number string) (bool, error)
	ExistsByNumberExcludingID(ctx context.Context, number string, excludeID int) (bool, error)
	
	// Code generation helpers
	GetLastPaymentID(ctx context.Context) (int, error)
	
	// Status operations
	UpdateStatus(ctx context.Context, id int, status inventory.PaymentStatus) error
	UpdateOverdueStatus(ctx context.Context) error
	
	// Financial reporting
	GetPaymentSummary(ctx context.Context, supplierID *int, startDate, endDate string) (*inventory.PaymentSummary, error)
	GetOutstandingBalance(ctx context.Context, supplierID int) (float64, error)
	GetTotalPaid(ctx context.Context, supplierID int, startDate, endDate string) (float64, error)
}