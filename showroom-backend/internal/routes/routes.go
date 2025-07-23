package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/config"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/handlers/admin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/handlers/auth"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/handlers/products"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/handlers/repair"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/handlers/vehicle_purchase"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/middleware"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/utils"
)

// Router handles all HTTP routes
type Router struct {
	authHandler               *auth.Handler
	adminHandler              *admin.Handler
	customerHandler           *admin.CustomerHandler
	supplierHandler           *admin.SupplierHandler
	vehicleMasterHandler      *admin.VehicleMasterHandler
	productCategoryHandler    *admin.ProductCategoryHandler
	productHandler            *admin.ProductHandler
	purchaseOrderHandler      *admin.PurchaseOrderHandler
	poDetailHandler           *products.PurchaseOrderDetailHandler
	goodsReceiptHandler       *products.GoodsReceiptHandler
	stockMovementHandler      *products.StockMovementHandler
	stockAdjustmentHandler    *products.StockAdjustmentHandler
	supplierPaymentHandler    *products.SupplierPaymentHandler
	
	// Phase 4 & 5 handlers
	vehiclePurchaseTransactionHandler *vehicle_purchase.TransactionHandler
	vehiclePurchasePaymentHandler     *vehicle_purchase.PaymentHandler
	repairDamageHandler               *repair.DamageHandler
	repairWorkOrderHandler            *repair.WorkOrderHandler
	repairWorkDetailHandler           *repair.WorkDetailHandler
	repairPartsUsageHandler           *repair.PartsUsageHandler
	qualityInspectionHandler          *repair.QualityInspectionHandler
	
	jwtManager                *utils.JWTManager
	sessionRepo               interfaces.UserSessionRepository
	config                    *config.Config
}

// NewRouter creates a new router
func NewRouter(
	authHandler *auth.Handler,
	adminHandler *admin.Handler,
	customerHandler *admin.CustomerHandler,
	supplierHandler *admin.SupplierHandler,
	vehicleMasterHandler *admin.VehicleMasterHandler,
	productCategoryHandler *admin.ProductCategoryHandler,
	productHandler *admin.ProductHandler,
	purchaseOrderHandler *admin.PurchaseOrderHandler,
	poDetailHandler *products.PurchaseOrderDetailHandler,
	goodsReceiptHandler *products.GoodsReceiptHandler,
	stockMovementHandler *products.StockMovementHandler,
	stockAdjustmentHandler *products.StockAdjustmentHandler,
	supplierPaymentHandler *products.SupplierPaymentHandler,
	
	// Phase 4 & 5 handlers
	vehiclePurchaseTransactionHandler *vehicle_purchase.TransactionHandler,
	vehiclePurchasePaymentHandler *vehicle_purchase.PaymentHandler,
	repairDamageHandler *repair.DamageHandler,
	repairWorkOrderHandler *repair.WorkOrderHandler,
	repairWorkDetailHandler *repair.WorkDetailHandler,
	repairPartsUsageHandler *repair.PartsUsageHandler,
	qualityInspectionHandler *repair.QualityInspectionHandler,
	
	jwtManager *utils.JWTManager,
	sessionRepo interfaces.UserSessionRepository,
	config *config.Config,
) *Router {
	return &Router{
		authHandler:               authHandler,
		adminHandler:              adminHandler,
		customerHandler:           customerHandler,
		supplierHandler:           supplierHandler,
		vehicleMasterHandler:      vehicleMasterHandler,
		productCategoryHandler:    productCategoryHandler,
		productHandler:            productHandler,
		purchaseOrderHandler:      purchaseOrderHandler,
		poDetailHandler:           poDetailHandler,
		goodsReceiptHandler:       goodsReceiptHandler,
		stockMovementHandler:      stockMovementHandler,
		stockAdjustmentHandler:    stockAdjustmentHandler,
		supplierPaymentHandler:    supplierPaymentHandler,
		
		// Phase 4 & 5 handlers
		vehiclePurchaseTransactionHandler: vehiclePurchaseTransactionHandler,
		vehiclePurchasePaymentHandler:     vehiclePurchasePaymentHandler,
		repairDamageHandler:               repairDamageHandler,
		repairWorkOrderHandler:            repairWorkOrderHandler,
		repairWorkDetailHandler:           repairWorkDetailHandler,
		repairPartsUsageHandler:           repairPartsUsageHandler,
		qualityInspectionHandler:          qualityInspectionHandler,
		
		jwtManager:                jwtManager,
		sessionRepo:               sessionRepo,
		config:                    config,
	}
}

// SetupRoutes configures all routes
func (r *Router) SetupRoutes() *gin.Engine {
	// Set Gin mode based on environment
	if r.config.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())

	// Health check endpoint (no auth required)
	router.GET("/api/v1/health", r.healthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")

	// Authentication routes (no auth required)
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/login", r.authHandler.Login)
		
		// Protected auth routes
		authProtected := authGroup.Use(middleware.AuthMiddleware(r.jwtManager, r.sessionRepo))
		{
			authProtected.POST("/logout", r.authHandler.Logout)
			authProtected.GET("/me", r.authHandler.Me)
			authProtected.GET("/profile", r.authHandler.Profile)
			authProtected.POST("/change-password", r.authHandler.ChangePassword)
			authProtected.POST("/refresh", r.authHandler.RefreshToken)
		}
	}

	// Admin routes (admin role required)
	adminGroup := v1.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware(r.jwtManager, r.sessionRepo))
	adminGroup.Use(middleware.RequireRole("admin"))
	{
		// User management
		userGroup := adminGroup.Group("/users")
		{
			userGroup.POST("", r.adminHandler.CreateUser)
			userGroup.GET("", r.adminHandler.GetUsers)
			userGroup.GET("/:id", r.adminHandler.GetUser)
			userGroup.PUT("/:id", r.adminHandler.UpdateUser)
			userGroup.DELETE("/:id", r.adminHandler.DeleteUser)
			userGroup.GET("/role/:role", r.adminHandler.GetUsersByRole)
			userGroup.GET("/:id/sessions", r.adminHandler.GetUserSessions)
			userGroup.DELETE("/:id/sessions", r.adminHandler.RevokeUserSessions)
		}

		// Customer management
		customerGroup := adminGroup.Group("/customers")
		{
			customerGroup.POST("", r.customerHandler.CreateCustomer)
			customerGroup.GET("", r.customerHandler.GetCustomers)
			customerGroup.GET("/:id", r.customerHandler.GetCustomer)
			customerGroup.PUT("/:id", r.customerHandler.UpdateCustomer)
			customerGroup.DELETE("/:id", r.customerHandler.DeleteCustomer)
		}

		// Supplier management
		supplierGroup := adminGroup.Group("/suppliers")
		{
			supplierGroup.POST("", r.supplierHandler.CreateSupplier)
			supplierGroup.GET("", r.supplierHandler.GetSuppliers)
			supplierGroup.GET("/:id", r.supplierHandler.GetSupplier)
			supplierGroup.PUT("/:id", r.supplierHandler.UpdateSupplier)
			supplierGroup.DELETE("/:id", r.supplierHandler.DeleteSupplier)
		}

		// Vehicle brand management
		vehicleBrandGroup := adminGroup.Group("/vehicle-brands")
		{
			vehicleBrandGroup.POST("", r.vehicleMasterHandler.CreateVehicleBrand)
			vehicleBrandGroup.GET("", r.vehicleMasterHandler.GetVehicleBrands)
			vehicleBrandGroup.GET("/:id", r.vehicleMasterHandler.GetVehicleBrand)
			vehicleBrandGroup.PUT("/:id", r.vehicleMasterHandler.UpdateVehicleBrand)
			vehicleBrandGroup.DELETE("/:id", r.vehicleMasterHandler.DeleteVehicleBrand)
		}

		// Vehicle category management
		vehicleCategoryGroup := adminGroup.Group("/vehicle-categories")
		{
			vehicleCategoryGroup.POST("", r.vehicleMasterHandler.CreateVehicleCategory)
			vehicleCategoryGroup.GET("", r.vehicleMasterHandler.GetVehicleCategories)
			vehicleCategoryGroup.GET("/:id", r.vehicleMasterHandler.GetVehicleCategory)
			vehicleCategoryGroup.PUT("/:id", r.vehicleMasterHandler.UpdateVehicleCategory)
			vehicleCategoryGroup.DELETE("/:id", r.vehicleMasterHandler.DeleteVehicleCategory)
		}

		// Vehicle model management
		vehicleModelGroup := adminGroup.Group("/vehicle-models")
		{
			vehicleModelGroup.POST("", r.vehicleMasterHandler.CreateVehicleModel)
			vehicleModelGroup.GET("", r.vehicleMasterHandler.GetVehicleModels)
			vehicleModelGroup.GET("/:id", r.vehicleMasterHandler.GetVehicleModel)
			vehicleModelGroup.PUT("/:id", r.vehicleMasterHandler.UpdateVehicleModel)
			vehicleModelGroup.DELETE("/:id", r.vehicleMasterHandler.DeleteVehicleModel)
		}

		// Product category management
		productCategoryGroup := adminGroup.Group("/product-categories")
		{
			productCategoryGroup.POST("", r.productCategoryHandler.CreateProductCategory)
			productCategoryGroup.GET("", r.productCategoryHandler.GetProductCategories)
			productCategoryGroup.GET("/:id", r.productCategoryHandler.GetProductCategory)
			productCategoryGroup.PUT("/:id", r.productCategoryHandler.UpdateProductCategory)
			productCategoryGroup.DELETE("/:id", r.productCategoryHandler.DeleteProductCategory)
			productCategoryGroup.GET("/tree", r.productCategoryHandler.GetProductCategoryTree)
			productCategoryGroup.GET("/:id/children", r.productCategoryHandler.GetProductCategoryChildren)
		}

		// Product management
		productGroup := adminGroup.Group("/products")
		{
			productGroup.POST("", r.productHandler.CreateProduct)
			productGroup.GET("", r.productHandler.GetProducts)
			productGroup.GET("/:id", r.productHandler.GetProduct)
			productGroup.PUT("/:id", r.productHandler.UpdateProduct)
			productGroup.DELETE("/:id", r.productHandler.DeleteProduct)
			productGroup.GET("/low-stock", r.productHandler.GetLowStockProducts)
			productGroup.GET("/:id/stock-movements", r.stockMovementHandler.GetProductStockMovements)
			productGroup.GET("/:id/stock-history", r.stockMovementHandler.GetProductStockHistory)
			productGroup.GET("/:id/current-stock", r.stockMovementHandler.GetCurrentStock)
			productGroup.GET("/:id/adjustments", r.stockAdjustmentHandler.GetProductStockAdjustments)
		}

		// Purchase Order management
		purchaseOrderGroup := adminGroup.Group("/purchase-orders")
		{
			purchaseOrderGroup.POST("", r.purchaseOrderHandler.CreatePurchaseOrder)
			purchaseOrderGroup.GET("", r.purchaseOrderHandler.GetPurchaseOrders)
			purchaseOrderGroup.GET("/:id", r.purchaseOrderHandler.GetPurchaseOrder)
			purchaseOrderGroup.PUT("/:id", r.purchaseOrderHandler.UpdatePurchaseOrder)
			purchaseOrderGroup.POST("/:id/approve", r.purchaseOrderHandler.ApprovePurchaseOrder)
			purchaseOrderGroup.POST("/:id/cancel", r.purchaseOrderHandler.CancelPurchaseOrder)
			purchaseOrderGroup.GET("/pending-approval", r.purchaseOrderHandler.GetPendingApproval)
			
			// Purchase Order Details
			purchaseOrderGroup.POST("/:id/details", r.poDetailHandler.CreatePODetail)
			purchaseOrderGroup.GET("/:id/details", r.poDetailHandler.GetPODetails)
			purchaseOrderGroup.GET("/:id/pending-receipt-items", r.poDetailHandler.GetPendingReceiptItems)
			purchaseOrderGroup.POST("/:id/bulk-details", r.poDetailHandler.BulkCreatePODetails)
		}

		// Purchase Order Details management
		poDetailGroup := adminGroup.Group("/purchase-order-details")
		{
			poDetailGroup.GET("/:id", r.poDetailHandler.GetPODetail)
			poDetailGroup.PUT("/:id", r.poDetailHandler.UpdatePODetail)
			poDetailGroup.DELETE("/:id", r.poDetailHandler.DeletePODetail)
		}

		// Goods Receipt management
		goodsReceiptGroup := adminGroup.Group("/goods-receipts")
		{
			goodsReceiptGroup.POST("", r.goodsReceiptHandler.CreateGoodsReceipt)
			goodsReceiptGroup.GET("", r.goodsReceiptHandler.ListGoodsReceipts)
			goodsReceiptGroup.GET("/:id", r.goodsReceiptHandler.GetGoodsReceipt)
			goodsReceiptGroup.PUT("/:id", r.goodsReceiptHandler.UpdateGoodsReceipt)
			goodsReceiptGroup.DELETE("/:id", r.goodsReceiptHandler.DeleteGoodsReceipt)
			goodsReceiptGroup.POST("/:id/process", r.goodsReceiptHandler.ProcessGoodsReceipt)
			goodsReceiptGroup.POST("/:id/details", r.goodsReceiptHandler.AddReceiptDetail)
			goodsReceiptGroup.GET("/:id/details", r.goodsReceiptHandler.GetReceiptDetails)
			goodsReceiptGroup.POST("/:id/bulk-receive", r.goodsReceiptHandler.BulkReceiveItems)
		}

		// Stock Movement management
		stockMovementGroup := adminGroup.Group("/stock-movements")
		{
			stockMovementGroup.POST("", r.stockMovementHandler.CreateStockMovement)
			stockMovementGroup.GET("", r.stockMovementHandler.ListStockMovements)
			stockMovementGroup.GET("/:id", r.stockMovementHandler.GetStockMovement)
			stockMovementGroup.POST("/transfer", r.stockMovementHandler.TransferStock)
		}

		// Stock Adjustment management
		stockAdjustmentGroup := adminGroup.Group("/stock-adjustments")
		{
			stockAdjustmentGroup.POST("", r.stockAdjustmentHandler.CreateStockAdjustment)
			stockAdjustmentGroup.GET("", r.stockAdjustmentHandler.ListStockAdjustments)
			stockAdjustmentGroup.GET("/:id", r.stockAdjustmentHandler.GetStockAdjustment)
			stockAdjustmentGroup.PUT("/:id", r.stockAdjustmentHandler.UpdateStockAdjustment)
			stockAdjustmentGroup.DELETE("/:id", r.stockAdjustmentHandler.DeleteStockAdjustment)
			stockAdjustmentGroup.GET("/pending", r.stockAdjustmentHandler.GetPendingAdjustments)
			stockAdjustmentGroup.POST("/:id/approve", r.stockAdjustmentHandler.ApproveStockAdjustment)
			stockAdjustmentGroup.GET("/variance-report", r.stockAdjustmentHandler.GetVarianceReport)
			stockAdjustmentGroup.POST("/physical-count", r.stockAdjustmentHandler.CreatePhysicalCountAdjustments)
			stockAdjustmentGroup.POST("/bulk-approve", r.stockAdjustmentHandler.BulkApproveAdjustments)
		}

		// Supplier Payment management
		supplierPaymentGroup := adminGroup.Group("/supplier-payments")
		{
			supplierPaymentGroup.POST("", r.supplierPaymentHandler.CreateSupplierPayment)
			supplierPaymentGroup.GET("", r.supplierPaymentHandler.ListSupplierPayments)
			supplierPaymentGroup.GET("/:id", r.supplierPaymentHandler.GetSupplierPayment)
			supplierPaymentGroup.PUT("/:id", r.supplierPaymentHandler.UpdateSupplierPayment)
			supplierPaymentGroup.DELETE("/:id", r.supplierPaymentHandler.DeleteSupplierPayment)
			supplierPaymentGroup.POST("/:id/process", r.supplierPaymentHandler.ProcessPayment)
			supplierPaymentGroup.PUT("/:id/status", r.supplierPaymentHandler.UpdatePaymentStatus)
			supplierPaymentGroup.GET("/overdue", r.supplierPaymentHandler.GetOverduePayments)
			supplierPaymentGroup.GET("/summary", r.supplierPaymentHandler.GetPaymentSummary)
			supplierPaymentGroup.POST("/update-overdue", r.supplierPaymentHandler.UpdateOverduePayments)
			supplierPaymentGroup.POST("/calculate-terms", r.supplierPaymentHandler.CalculatePaymentTerms)
		}
	}

	// Phase 4: Vehicle Purchase routes (admin and sales role required)
	vehiclePurchaseGroup := v1.Group("/vehicle-purchases")
	vehiclePurchaseGroup.Use(middleware.AuthMiddleware(r.jwtManager, r.sessionRepo))
	vehiclePurchaseGroup.Use(middleware.RequireRole("admin", "sales"))
	{
		// Vehicle Purchase Transactions
		transactionGroup := vehiclePurchaseGroup.Group("/transactions")
		{
			transactionGroup.POST("", r.vehiclePurchaseTransactionHandler.CreateTransaction)
			transactionGroup.GET("", r.vehiclePurchaseTransactionHandler.GetTransactions)
			transactionGroup.GET("/:id", r.vehiclePurchaseTransactionHandler.GetTransaction)
			transactionGroup.GET("/number/:number", r.vehiclePurchaseTransactionHandler.GetTransactionByNumber)
			transactionGroup.GET("/vin/:vin", r.vehiclePurchaseTransactionHandler.GetTransactionByVIN)
			transactionGroup.PUT("/:id", r.vehiclePurchaseTransactionHandler.UpdateTransaction)
			transactionGroup.POST("/:id/inspect", r.vehiclePurchaseTransactionHandler.CompleteInspection)
			transactionGroup.POST("/:id/approve", r.vehiclePurchaseTransactionHandler.ProcessApproval)
			transactionGroup.GET("/pending-inspection", r.vehiclePurchaseTransactionHandler.GetPendingInspections)
			transactionGroup.GET("/pending-approval", r.vehiclePurchaseTransactionHandler.GetPendingApprovals)
			transactionGroup.GET("/:transaction_id/payments", r.vehiclePurchasePaymentHandler.GetPaymentsByTransaction)
			transactionGroup.GET("/:transaction_id/payment-summary", r.vehiclePurchasePaymentHandler.GetPaymentSummary)
		}

		// Vehicle Purchase Payments
		paymentGroup := vehiclePurchaseGroup.Group("/payments")
		{
			paymentGroup.POST("", r.vehiclePurchasePaymentHandler.CreatePayment)
			paymentGroup.GET("", r.vehiclePurchasePaymentHandler.GetPayments)
			paymentGroup.GET("/:id", r.vehiclePurchasePaymentHandler.GetPayment)
			paymentGroup.GET("/number/:number", r.vehiclePurchasePaymentHandler.GetPaymentByNumber)
			paymentGroup.POST("/:id/process", r.vehiclePurchasePaymentHandler.ProcessPayment)
			paymentGroup.POST("/:id/approve", r.vehiclePurchasePaymentHandler.ProcessPaymentApproval)
			paymentGroup.GET("/pending-approval", r.vehiclePurchasePaymentHandler.GetPendingApprovals)
			paymentGroup.GET("/overdue", r.vehiclePurchasePaymentHandler.GetOverduePayments)
		}

		// Dashboard
		vehiclePurchaseGroup.GET("/dashboard", r.vehiclePurchaseTransactionHandler.GetDashboardStats)
	}

	// Phase 5: Repair Management routes (admin, mechanic, and inspector roles)
	repairGroup := v1.Group("/repairs")
	repairGroup.Use(middleware.AuthMiddleware(r.jwtManager, r.sessionRepo))
	repairGroup.Use(middleware.RequireRole("admin", "mechanic", "inspector"))
	{
		// Vehicle Damages
		damageGroup := repairGroup.Group("/damages")
		{
			damageGroup.POST("", r.repairDamageHandler.CreateDamage)
			damageGroup.GET("", r.repairDamageHandler.GetDamages)
			damageGroup.GET("/:id", r.repairDamageHandler.GetDamage)
			damageGroup.PUT("/:id", r.repairDamageHandler.UpdateDamage)
			damageGroup.POST("/:id/assess", r.repairDamageHandler.AssessDamage)
			damageGroup.GET("/high-priority", r.repairDamageHandler.GetHighPriorityDamages)
		}

		// Work Orders
		workOrderGroup := repairGroup.Group("/work-orders")
		{
			workOrderGroup.POST("", r.repairWorkOrderHandler.CreateWorkOrder)
			workOrderGroup.GET("", r.repairWorkOrderHandler.GetWorkOrders)
			workOrderGroup.GET("/:id", r.repairWorkOrderHandler.GetWorkOrder)
			workOrderGroup.GET("/number/:number", r.repairWorkOrderHandler.GetWorkOrderByNumber)
			workOrderGroup.PUT("/:id", r.repairWorkOrderHandler.UpdateWorkOrder)
			workOrderGroup.POST("/:id/assign", r.repairWorkOrderHandler.AssignMechanic)
			workOrderGroup.POST("/:id/approve", r.repairWorkOrderHandler.ProcessApproval)
			workOrderGroup.GET("/pending-approval", r.repairWorkOrderHandler.GetPendingApprovals)
			workOrderGroup.GET("/mechanic/:mechanic_id", r.repairWorkOrderHandler.GetWorkOrdersByMechanic)
			
			// Work Order related endpoints
			workOrderGroup.GET("/:work_order_id/details", r.repairWorkDetailHandler.GetWorkDetailsByWorkOrder)
			workOrderGroup.GET("/:work_order_id/detail-summary", r.repairWorkDetailHandler.GetWorkDetailSummary)
			workOrderGroup.GET("/:work_order_id/parts-summary", r.repairPartsUsageHandler.GetPartsUsageSummary)
			workOrderGroup.GET("/:work_order_id/inventory-impact", r.repairPartsUsageHandler.GetInventoryImpact)
			workOrderGroup.GET("/:work_order_id/inspections", r.qualityInspectionHandler.GetInspectionsByWorkOrder)
			workOrderGroup.GET("/:work_order_id/quality-metrics", r.qualityInspectionHandler.GetQualityMetrics)
		}

		// Work Details
		workDetailGroup := repairGroup.Group("/work-details")
		{
			workDetailGroup.POST("", r.repairWorkDetailHandler.CreateWorkDetail)
			workDetailGroup.GET("", r.repairWorkDetailHandler.GetWorkDetails)
			workDetailGroup.GET("/:id", r.repairWorkDetailHandler.GetWorkDetail)
			workDetailGroup.PUT("/:id", r.repairWorkDetailHandler.UpdateWorkDetail)
			workDetailGroup.POST("/:id/progress", r.repairWorkDetailHandler.UpdateProgress)
			workDetailGroup.POST("/:id/assign", r.repairWorkDetailHandler.AssignMechanic)
			workDetailGroup.POST("/:id/quality-check", r.repairWorkDetailHandler.PerformQualityCheck)
			
			// Work Detail related endpoints
			workDetailGroup.GET("/:work_detail_id/parts-usage", r.repairPartsUsageHandler.GetPartsUsageByWorkDetail)
			workDetailGroup.POST("/:work_detail_id/issue-parts", r.repairPartsUsageHandler.IssuePartsForRepair)
		}

		// Parts Usage
		partsUsageGroup := repairGroup.Group("/parts-usage")
		{
			partsUsageGroup.POST("", r.repairPartsUsageHandler.CreatePartsUsage)
			partsUsageGroup.GET("", r.repairPartsUsageHandler.GetPartsUsages)
			partsUsageGroup.GET("/:id", r.repairPartsUsageHandler.GetPartsUsage)
			partsUsageGroup.POST("/:id/approve", r.repairPartsUsageHandler.ProcessApproval)
		}

		// Quality Inspections
		inspectionGroup := repairGroup.Group("/inspections")
		{
			inspectionGroup.POST("", r.qualityInspectionHandler.CreateInspection)
			inspectionGroup.GET("", r.qualityInspectionHandler.GetInspections)
			inspectionGroup.GET("/:id", r.qualityInspectionHandler.GetInspection)
			inspectionGroup.POST("/:id/sign-off", r.qualityInspectionHandler.SignOffInspection)
			inspectionGroup.POST("/:id/schedule-rework", r.qualityInspectionHandler.ScheduleRework)
			inspectionGroup.POST("/schedule", r.qualityInspectionHandler.ScheduleInspection)
			inspectionGroup.GET("/failed", r.qualityInspectionHandler.GetFailedInspections)
			inspectionGroup.GET("/rework-required", r.qualityInspectionHandler.GetReworkRequired)
			inspectionGroup.GET("/dashboard", r.qualityInspectionHandler.GetInspectionDashboard)
		}

		// Transaction related repair endpoints
		transactionRepairGroup := repairGroup.Group("/transactions")
		{
			transactionRepairGroup.GET("/:transaction_id/damage-summary", r.repairDamageHandler.GetDamageSummary)
			transactionRepairGroup.GET("/:transaction_id/work-order-summary", r.repairWorkOrderHandler.GetWorkOrderSummary)
		}
	}

	return router
}

// healthCheck handles health check requests
func (r *Router) healthCheck(c *gin.Context) {
	response := &common.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   r.config.App.Version,
	}
	c.JSON(http.StatusOK, response)
}