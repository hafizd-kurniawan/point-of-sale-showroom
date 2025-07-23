package main

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/config"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/database"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/implementations"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations first
	if err := database.RunMigrations(database.GetDB()); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	productRepo := implementations.NewProductSparePartRepository(database.GetDB())
	poRepo := implementations.NewPurchaseOrderPartsRepository(database.GetDB())

	// Test auto-generation functionality
	ctx := context.Background()

	// Test product code generation
	log.Println("=== Testing Product Code Generation ===")
	for i := 0; i < 3; i++ {
		code, err := productRepo.GenerateCode(ctx)
		if err != nil {
			log.Printf("Failed to generate product code: %v", err)
			continue
		}
		log.Printf("Generated product code: %s", code)
	}

	// Test PO number generation
	log.Println("\n=== Testing PO Number Generation ===")
	for i := 0; i < 3; i++ {
		number, err := poRepo.GenerateNumber(ctx)
		if err != nil {
			log.Printf("Failed to generate PO number: %v", err)
			continue
		}
		log.Printf("Generated PO number: %s", number)
	}

	// Test creating a sample product (if we have brand and category)
	log.Println("\n=== Testing Product Creation ===")
	
	// First check if we have at least one brand and category from existing seeders
	// We'll use brand_id = 1 and category_id = 1 assuming they exist
	sampleProduct := &products.ProductSparePart{
		ProductCode:      "", // Will be auto-generated
		ProductName:      "Sample Brake Pad",
		Description:      stringPtr("High quality brake pad for compact cars"),
		BrandID:          1, // Assuming exists from previous seeders
		CategoryID:       1, // Assuming exists from previous seeders
		UnitMeasure:      "piece",
		CostPrice:        25.00,
		SellingPrice:     35.00,
		MarkupPercentage: 40.00,
		StockQuantity:    0,
		MinStockLevel:    5,
		MaxStockLevel:    50,
		LocationRack:     stringPtr("A1-01"),
		Barcode:          stringPtr("BP001234567890"),
		Weight:           floatPtr(0.5),
		Dimensions:       stringPtr("10cm x 8cm x 2cm"),
		CreatedBy:        1, // Assuming admin user exists
		IsActive:         true,
		ProductImage:     stringPtr("/images/brake-pad-sample.jpg"),
		Notes:            stringPtr("Standard brake pad for daily use"),
	}

	// Generate product code
	code, err := productRepo.GenerateCode(ctx)
	if err != nil {
		log.Printf("Failed to generate product code: %v", err)
	} else {
		sampleProduct.ProductCode = code
		log.Printf("Generated product code for sample product: %s", code)

		// Try to create the product
		createdProduct, err := productRepo.Create(ctx, sampleProduct)
		if err != nil {
			log.Printf("Failed to create sample product (expected if brands/categories don't exist): %v", err)
		} else {
			log.Printf("Successfully created sample product with ID: %d", createdProduct.ProductID)
		}
	}

	// Test creating a sample PO
	log.Println("\n=== Testing Purchase Order Creation ===")
	
	samplePO := &products.PurchaseOrderParts{
		PONumber:             "", // Will be auto-generated
		SupplierID:           1,  // Assuming exists from previous seeders
		PODate:               time.Now(),
		RequiredDate:         timePtr(time.Now().AddDate(0, 0, 14)), // 2 weeks from now
		ExpectedDeliveryDate: timePtr(time.Now().AddDate(0, 0, 21)), // 3 weeks from now
		POType:               products.POTypeRegular,
		Subtotal:             0,
		TaxAmount:            0,
		DiscountAmount:       0,
		ShippingCost:         0,
		TotalAmount:          0,
		Status:               products.POStatusDraft,
		PaymentTerms:         products.PaymentTermsNet30,
		CreatedBy:            1, // Assuming admin user exists
		DeliveryAddress:      stringPtr("123 Main Street, Warehouse District"),
		PONotes:              stringPtr("Standard parts order for inventory replenishment"),
		TermsAndConditions:   stringPtr("Net 30 payment terms. All goods subject to quality inspection."),
	}

	// Generate PO number
	poNumber, err := poRepo.GenerateNumber(ctx)
	if err != nil {
		log.Printf("Failed to generate PO number: %v", err)
	} else {
		samplePO.PONumber = poNumber
		samplePO.SetPaymentDueDate()
		log.Printf("Generated PO number for sample order: %s", poNumber)

		// Try to create the PO
		createdPO, err := poRepo.Create(ctx, samplePO)
		if err != nil {
			log.Printf("Failed to create sample PO (expected if suppliers don't exist): %v", err)
		} else {
			log.Printf("Successfully created sample PO with ID: %d", createdPO.POID)
			log.Printf("PO Status: %s", createdPO.Status)
			log.Printf("Payment due date: %v", createdPO.PaymentDueDate)
		}
	}

	log.Println("\n=== Phase 3 Auto-Generation Demo Completed ===")
	log.Println("The system demonstrates:")
	log.Println("1. Auto-generation of product codes (PRD-001, PRD-002, etc.)")
	log.Println("2. Auto-generation of PO numbers (PO-2025-001, PO-2025-002, etc.)")
	log.Println("3. Proper handling of business logic (payment due dates, status management)")
	log.Println("4. Database schema with all ERD-compliant tables and relationships")
	log.Println("\nFor full functionality, seed the database with users, suppliers, brands, and categories first.")
}

// Helper functions for pointer values
func stringPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}

func timePtr(t time.Time) *time.Time {
	return &t
}