package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/config"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/database"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/implementations"
	inventoryService "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services/inventory"
	masterService "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services/master"
)

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}

func intPtr(i int) *int {
	return &i
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func main() {
	// Load environment variables
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	if err := initializeDatabase(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Check if admin user exists
	if !checkAdminUserExists() {
		log.Fatal("Admin user not found. Please run the application first to create an admin user.")
	}

	// Seed inventory data
	log.Println("Starting inventory data seeding...")
	if err := seedInventoryData(); err != nil {
		log.Fatalf("Failed to seed inventory data: %v", err)
	}

	log.Println("Inventory data seeding completed successfully!")
}

// initializeDatabase sets up database connection
func initializeDatabase(cfg *config.Config) error {
	if err := database.Connect(cfg); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	return nil
}

// checkAdminUserExists checks if an admin user exists
func checkAdminUserExists() bool {
	db := database.GetDB()
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&count).Error
	return err == nil && count > 0
}

// seedInventoryData seeds all inventory-related data
func seedInventoryData() error {
	ctx := context.Background()
	db := database.GetDB()

	// Initialize repositories
	supplierRepo := implementations.NewSupplierRepository(db)
	productCategoryRepo := implementations.NewProductCategoryRepository(db)
	productRepo := implementations.NewProductSparePartRepository(db)
	purchaseOrderRepo := implementations.NewPurchaseOrderRepository(db)
	purchaseOrderDetailRepo := implementations.NewPurchaseOrderDetailRepository(db)
	goodsReceiptRepo := implementations.NewGoodsReceiptRepository(db)
	goodsReceiptDetailRepo := implementations.NewGoodsReceiptDetailRepository(db)
	stockMovementRepo := implementations.NewStockMovementRepository(db)
	stockAdjustmentRepo := implementations.NewStockAdjustmentRepository(db)
	supplierPaymentRepo := implementations.NewSupplierPaymentRepository(db)
	userRepo := implementations.NewUserRepository(db)

	// Initialize services
	supplierService := masterService.NewSupplierService(supplierRepo)
	productCategoryService := masterService.NewProductCategoryService(productCategoryRepo)
	productService := inventoryService.NewProductService(productRepo, productCategoryRepo, stockMovementRepo)
	purchaseOrderService := inventoryService.NewPurchaseOrderService(purchaseOrderRepo, purchaseOrderDetailRepo, supplierRepo, productRepo)
	goodsReceiptService := inventoryService.NewGoodsReceiptService(goodsReceiptRepo, goodsReceiptDetailRepo, purchaseOrderRepo, purchaseOrderDetailRepo, stockMovementRepo)
	stockMovementService := inventoryService.NewStockMovementService(stockMovementRepo, productRepo)
	stockAdjustmentService := inventoryService.NewStockAdjustmentService(stockAdjustmentRepo, productRepo, stockMovementRepo, userRepo)
	supplierPaymentService := inventoryService.NewSupplierPaymentService(supplierPaymentRepo, supplierRepo, purchaseOrderRepo)

	// Seed data
	log.Println("Seeding product catalog...")
	if err := seedProductCatalog(ctx, productService, productCategoryService); err != nil {
		return fmt.Errorf("failed to seed products: %w", err)
	}

	log.Println("Seeding purchase orders...")
	if err := seedPurchaseOrders(ctx, purchaseOrderService, supplierService, productService); err != nil {
		return fmt.Errorf("failed to seed purchase orders: %w", err)
	}

	log.Println("Seeding goods receipts...")
	if err := seedGoodsReceipts(ctx, goodsReceiptService, purchaseOrderService); err != nil {
		return fmt.Errorf("failed to seed goods receipts: %w", err)
	}

	log.Println("Seeding stock movements...")
	if err := seedStockMovements(ctx, stockMovementService, productService); err != nil {
		return fmt.Errorf("failed to seed stock movements: %w", err)
	}

	log.Println("Seeding stock adjustments...")
	if err := seedStockAdjustments(ctx, stockAdjustmentService, productService); err != nil {
		return fmt.Errorf("failed to seed stock adjustments: %w", err)
	}

	log.Println("Seeding supplier payments...")
	if err := seedSupplierPayments(ctx, supplierPaymentService, purchaseOrderService); err != nil {
		return fmt.Errorf("failed to seed supplier payments: %w", err)
	}

	return nil
}

// seedProductCatalog creates a variety of spare parts products
func seedProductCatalog(ctx context.Context, productService *inventoryService.ProductService, categoryService *masterService.ProductCategoryService) error {
	// Get existing categories
	categories, _, err := categoryService.GetProductCategories(ctx, nil, 1, 100)
	if err != nil {
		return err
	}

	if len(categories) == 0 {
		return fmt.Errorf("no product categories found. Please seed master data first")
	}

	// Sample products for automotive spare parts
	products := []inventory.ProductSparePartCreateRequest{
		{
			ProductName:        "Oil Filter - Engine",
			Description:        stringPtr("High-quality engine oil filter for various car models"),
			CategoryID:         &categories[0].CategoryID,
			UnitMeasure:        "pieces",
			CostPrice:          25000,
			SellingPrice:       45000,
			MarkupPercentage:   float64Ptr(80),
			StockQuantity:      100,
			MinStockLevel:      intPtr(20),
			MaxStockLevel:      intPtr(200),
			LocationRack:       stringPtr("A1-01"),
			Weight:             float64Ptr(0.5),
			Dimensions:         stringPtr("10x10x5 cm"),
			Notes:              stringPtr("Compatible with Toyota, Honda, Nissan"),
		},
		{
			ProductName:        "Brake Pad Set - Front",
			Description:        stringPtr("Premium ceramic brake pads for front wheels"),
			CategoryID:         &categories[0].CategoryID,
			UnitMeasure:        "set",
			CostPrice:          150000,
			SellingPrice:       250000,
			MarkupPercentage:   float64Ptr(66.67),
			StockQuantity:      50,
			MinStockLevel:      intPtr(10),
			MaxStockLevel:      intPtr(100),
			LocationRack:       stringPtr("B2-03"),
			Weight:             float64Ptr(2.5),
			Dimensions:         stringPtr("25x15x3 cm"),
			Notes:              stringPtr("For sedans and SUVs"),
		},
		{
			ProductName:        "Spark Plug - Iridium",
			Description:        stringPtr("Long-lasting iridium spark plugs"),
			CategoryID:         &categories[0].CategoryID,
			UnitMeasure:        "pieces",
			CostPrice:          75000,
			SellingPrice:       120000,
			MarkupPercentage:   float64Ptr(60),
			StockQuantity:      200,
			MinStockLevel:      intPtr(50),
			MaxStockLevel:      intPtr(300),
			LocationRack:       stringPtr("C1-05"),
			Weight:             float64Ptr(0.1),
			Dimensions:         stringPtr("5x5x10 cm"),
			Notes:              stringPtr("Universal fit for most engines"),
		},
		{
			ProductName:        "Air Filter - Cabin",
			Description:        stringPtr("HEPA cabin air filter for clean air inside vehicle"),
			CategoryID:         &categories[0].CategoryID,
			UnitMeasure:        "pieces",
			CostPrice:          35000,
			SellingPrice:       60000,
			MarkupPercentage:   float64Ptr(71.43),
			StockQuantity:      75,
			MinStockLevel:      intPtr(15),
			MaxStockLevel:      intPtr(150),
			LocationRack:       stringPtr("A2-02"),
			Weight:             float64Ptr(0.3),
			Dimensions:         stringPtr("20x15x3 cm"),
			Notes:              stringPtr("Replaces every 6 months"),
		},
		{
			ProductName:        "Tire - All Season 195/65R15",
			Description:        stringPtr("All-season radial tire for passenger cars"),
			CategoryID:         &categories[0].CategoryID,
			UnitMeasure:        "pieces",
			CostPrice:          800000,
			SellingPrice:       1200000,
			MarkupPercentage:   float64Ptr(50),
			StockQuantity:      20,
			MinStockLevel:      intPtr(8),
			MaxStockLevel:      intPtr(40),
			LocationRack:       stringPtr("D1-01"),
			Weight:             float64Ptr(15.0),
			Dimensions:         stringPtr("65x65x20 cm"),
			Notes:              stringPtr("DOT certified, 5-year warranty"),
		},
	}

	// Create products
	for _, productReq := range products {
		_, err := productService.Create(ctx, &productReq, 1) // Assuming admin user ID = 1
		if err != nil {
			log.Printf("Failed to create product %s: %v", productReq.ProductName, err)
			continue
		}
		log.Printf("Created product: %s", productReq.ProductName)
	}

	return nil
}

// seedPurchaseOrders creates purchase orders with various states
func seedPurchaseOrders(ctx context.Context, poService *inventoryService.PurchaseOrderService, supplierService *masterService.SupplierService, productService *inventoryService.ProductService) error {
	// Get suppliers
	suppliers, _, err := supplierService.GetSuppliers(ctx, nil, 1, 100)
	if err != nil {
		return err
	}

	if len(suppliers) == 0 {
		return fmt.Errorf("no suppliers found. Please seed master data first")
	}

	// Get products
	products, _, err := productService.GetProducts(ctx, nil, 1, 100)
	if err != nil {
		return err
	}

	if len(products) == 0 {
		return fmt.Errorf("no products found")
	}

	// Sample purchase orders with different statuses
	poStatuses := []inventory.POStatus{
		inventory.POStatusDraft,
		inventory.POStatusSent,
		inventory.POStatusAcknowledged,
		inventory.POStatusPartialReceived,
		inventory.POStatusReceived,
		inventory.POStatusCompleted,
	}

	for i := 0; i < 10; i++ {
		supplier := suppliers[rand.Intn(len(suppliers))]
		status := poStatuses[rand.Intn(len(poStatuses))]

		// Create purchase order
		poReq := inventory.PurchaseOrderCreateRequest{
			SupplierID:             supplier.SupplierID,
			RequiredDate:           timePtr(time.Now().AddDate(0, 0, rand.Intn(30)+7)), // 7-37 days from now
			ExpectedDeliveryDate:   timePtr(time.Now().AddDate(0, 0, rand.Intn(45)+14)), // 14-59 days from now
			POType:                 inventory.POTypeRegular,
			PaymentTerms:           inventory.PaymentTermsNet30,
			DeliveryAddress:        stringPtr("Warehouse A, Jl. Industrial No. 123, Jakarta"),
			PONotes:                stringPtr("Standard parts order for monthly stock replenishment"),
			TermsAndConditions:     stringPtr("Standard 30-day payment terms. Quality inspection required upon delivery."),
		}

		// Add random products to the order
		numProducts := rand.Intn(3) + 1 // 1-3 products per order
		for j := 0; j < numProducts; j++ {
			product := products[rand.Intn(len(products))]
			quantity := rand.Intn(50) + 10 // 10-59 pieces

			detail := inventory.PurchaseOrderDetailCreateRequest{
				ProductID:       product.ProductID,
				ItemDescription: product.ProductName,
				QuantityOrdered: quantity,
				UnitCost:        product.CostPrice,
				ExpectedDate:    poReq.ExpectedDeliveryDate,
				ItemNotes:       stringPtr("Standard quality required"),
			}
			poReq.Details = append(poReq.Details, detail)
		}

		po, err := poService.CreatePurchaseOrder(ctx, &poReq, 1) // Admin user
		if err != nil {
			log.Printf("Failed to create purchase order: %v", err)
			continue
		}

		// Update status if not draft
		if status != inventory.POStatusDraft {
			// Simulate sending the PO
			_, err = poService.SendPurchaseOrder(ctx, po.PoID, 1)
			if err != nil {
				log.Printf("Failed to send PO %d: %v", po.PoID, err)
			}
		}

		log.Printf("Created purchase order: %s with status %s", po.PoNumber, status)
	}

	return nil
}

// seedGoodsReceipts creates goods receipts for some purchase orders
func seedGoodsReceipts(ctx context.Context, grService *inventoryService.GoodsReceiptService, poService *inventoryService.PurchaseOrderService) error {
	// Get sent/acknowledged purchase orders
	// This is a simplified approach - in real implementation, you'd query with filters
	
	// Create some sample goods receipts
	for i := 0; i < 5; i++ {
		// This is simplified - normally you'd create receipts for actual POs
		grReq := inventory.GoodsReceiptCreateRequest{
			PoID:                    i + 1, // Assuming PO IDs 1-5 exist
			SupplierDeliveryNote:    stringPtr(fmt.Sprintf("DN-%d", rand.Intn(9999)+1000)),
			SupplierInvoiceNumber:   stringPtr(fmt.Sprintf("INV-%d", rand.Intn(9999)+1000)),
			ReceiptNotes:            stringPtr("All items received in good condition"),
			Details: []inventory.GoodsReceiptDetailCreateRequest{
				{
					PoDetailID:        i + 1, // Simplified
					ProductID:         i + 1, // Simplified
					QuantityReceived:  rand.Intn(40) + 10,
					QuantityAccepted:  rand.Intn(40) + 10,
					QuantityRejected:  0,
					UnitCost:          float64(rand.Intn(100000) + 10000),
					ConditionReceived: inventory.ConditionGood,
					InspectionNotes:   stringPtr("Quality inspection passed"),
					ExpiryDate:        timePtr(time.Now().AddDate(2, 0, 0)), // 2 years from now
					BatchNumber:       stringPtr(fmt.Sprintf("BATCH-%d", rand.Intn(9999))),
				},
			},
		}

		_, err := grService.CreateGoodsReceipt(ctx, &grReq, 1)
		if err != nil {
			log.Printf("Failed to create goods receipt: %v", err)
			continue
		}

		log.Printf("Created goods receipt for PO %d", grReq.PoID)
	}

	return nil
}

// seedStockMovements creates various stock movement records
func seedStockMovements(ctx context.Context, smService *inventoryService.StockMovementService, productService *inventoryService.ProductService) error {
	// Get products
	products, _, err := productService.GetProducts(ctx, nil, 1, 100)
	if err != nil {
		return err
	}

	movementTypes := []inventory.MovementType{
		inventory.MovementTypeIn,
		inventory.MovementTypeOut,
		inventory.MovementTypeTransfer,
		inventory.MovementTypeAdjustment,
	}

	for i := 0; i < 20; i++ {
		product := products[rand.Intn(len(products))]
		movementType := movementTypes[rand.Intn(len(movementTypes))]

		switch movementType {
		case inventory.MovementTypeTransfer:
			_, err = smService.CreateTransfer(ctx, product.ProductID, rand.Intn(20)+1, "A1-01", "B2-03", "Stock relocation", 1)
		case inventory.MovementTypeAdjustment:
			_, err = smService.CreateManualAdjustment(ctx, product.ProductID, rand.Intn(10)+1, "Manual stock correction", 1)
		}

		if err != nil {
			log.Printf("Failed to create stock movement: %v", err)
			continue
		}

		log.Printf("Created stock movement for product %s", product.ProductName)
	}

	return nil
}

// seedStockAdjustments creates stock adjustment records
func seedStockAdjustments(ctx context.Context, saService *inventoryService.StockAdjustmentService, productService *inventoryService.ProductService) error {
	// Get products
	products, _, err := productService.GetProducts(ctx, nil, 1, 100)
	if err != nil {
		return err
	}

	adjustmentTypes := []inventory.AdjustmentType{
		inventory.AdjustmentTypePhysicalCount,
		inventory.AdjustmentTypeDamage,
		inventory.AdjustmentTypeExpired,
		inventory.AdjustmentTypeCorrection,
	}

	for i := 0; i < 8; i++ {
		product := products[rand.Intn(len(products))]
		adjustmentType := adjustmentTypes[rand.Intn(len(adjustmentTypes))]

		switch adjustmentType {
		case inventory.AdjustmentTypePhysicalCount:
			_, err = saService.CreatePhysicalCountAdjustment(ctx, product.ProductID, rand.Intn(50)+10, 1, stringPtr("Physical count discrepancy"))
		case inventory.AdjustmentTypeDamage:
			_, err = saService.CreateDamageAdjustment(ctx, product.ProductID, rand.Intn(5)+1, 1, "Damaged during handling", stringPtr("Items damaged during warehouse operations"))
		case inventory.AdjustmentTypeCorrection:
			_, err = saService.CreateWriteOffAdjustment(ctx, product.ProductID, rand.Intn(3)+1, 1, "System error correction", stringPtr("Correcting system calculation error"))
		}

		if err != nil {
			log.Printf("Failed to create stock adjustment: %v", err)
			continue
		}

		log.Printf("Created stock adjustment for product %s", product.ProductName)
	}

	return nil
}

// seedSupplierPayments creates payment records for purchase orders
func seedSupplierPayments(ctx context.Context, spService *inventoryService.SupplierPaymentService, poService *inventoryService.PurchaseOrderService) error {
	paymentStatuses := []inventory.PaymentStatus{
		inventory.PaymentStatusPending,
		inventory.PaymentStatusPartial,
		inventory.PaymentStatusPaid,
	}

	paymentMethods := []inventory.PaymentMethod{
		inventory.PaymentMethodTransfer,
		inventory.PaymentMethodCheck,
		inventory.PaymentMethodCash,
	}

	for i := 0; i < 6; i++ {
		status := paymentStatuses[rand.Intn(len(paymentStatuses))]
		method := paymentMethods[rand.Intn(len(paymentMethods))]

		invoiceAmount := float64(rand.Intn(5000000) + 500000) // 500k - 5.5M
		var paymentAmount float64
		
		switch status {
		case inventory.PaymentStatusPaid:
			paymentAmount = invoiceAmount
		case inventory.PaymentStatusPartial:
			paymentAmount = invoiceAmount * 0.5
		default:
			paymentAmount = 0
		}

		paymentReq := inventory.SupplierPaymentCreateRequest{
			SupplierID:       intPtr(rand.Intn(3) + 1), // Assuming supplier IDs 1-3
			PoID:             intPtr(i + 1),
			InvoiceAmount:    invoiceAmount,
			PaymentAmount:    paymentAmount,
			DiscountTaken:    invoiceAmount * 0.02, // 2% discount
			InvoiceDate:      time.Now().AddDate(0, 0, -rand.Intn(30)),
			PaymentDate:      timePtr(time.Now().AddDate(0, 0, rand.Intn(7))),
			DueDate:          time.Now().AddDate(0, 0, 30),
			PaymentMethod:    method,
			PaymentReference: stringPtr(fmt.Sprintf("TXN-%d", rand.Intn(999999)+100000)),
			InvoiceNumber:    fmt.Sprintf("INV-%d", rand.Intn(9999)+1000),
			PaymentNotes:     stringPtr("Payment processed successfully"),
		}

		_, err := spService.CreateSupplierPayment(ctx, &paymentReq)
		if err != nil {
			log.Printf("Failed to create supplier payment: %v", err)
			continue
		}

		log.Printf("Created supplier payment with status %s", status)
	}

	return nil
}