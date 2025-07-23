package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/config"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/database"
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

	db := database.GetDB()

	// Test Phase 4 & 5 database operations
	if err := testPhase4And5Database(db); err != nil {
		log.Fatalf("Phase 4 & 5 database test failed: %v", err)
	}

	log.Println("✅ All Phase 4 & 5 database tests passed!")
}

func testPhase4And5Database(db *sql.DB) error {
	ctx := context.Background()

	log.Println("🧪 Testing Phase 4 & 5 Database Implementation...")

	// Test Phase 4: Vehicle Purchase Flow
	if err := testVehiclePurchaseFlow(ctx, db); err != nil {
		return fmt.Errorf("vehicle purchase flow test failed: %w", err)
	}

	// Test Phase 5: Repair Management
	if err := testRepairManagement(ctx, db); err != nil {
		return fmt.Errorf("repair management test failed: %w", err)
	}

	return nil
}

func testVehiclePurchaseFlow(ctx context.Context, db *sql.DB) error {
	log.Println("📋 Testing Vehicle Purchase Flow...")

	// Test 1: Create vehicle purchase transaction
	log.Println("  - Testing vehicle purchase transaction creation...")
	transactionID, err := createTestTransaction(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create test transaction: %w", err)
	}
	log.Printf("    ✅ Created transaction ID: %d", transactionID)

	// Test 2: Create vehicle purchase payment
	log.Println("  - Testing vehicle purchase payment creation...")
	paymentID, err := createTestPayment(ctx, db, transactionID)
	if err != nil {
		return fmt.Errorf("failed to create test payment: %w", err)
	}
	log.Printf("    ✅ Created payment ID: %d", paymentID)

	// Test 3: Update transaction status
	log.Println("  - Testing transaction status update...")
	if err := updateTransactionStatus(ctx, db, transactionID); err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}
	log.Println("    ✅ Transaction status updated")

	return nil
}

func testRepairManagement(ctx context.Context, db *sql.DB) error {
	log.Println("🔧 Testing Repair Management...")

	// First, we need a transaction to link repairs to
	transactionID, err := createTestTransaction(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to create transaction for repair tests: %w", err)
	}

	// Test 1: Create vehicle damage
	log.Println("  - Testing vehicle damage creation...")
	damageID, err := createTestDamage(ctx, db, transactionID)
	if err != nil {
		return fmt.Errorf("failed to create test damage: %w", err)
	}
	log.Printf("    ✅ Created damage ID: %d", damageID)

	// Test 2: Create repair work order
	log.Println("  - Testing repair work order creation...")
	workOrderID, err := createTestWorkOrder(ctx, db, transactionID)
	if err != nil {
		return fmt.Errorf("failed to create test work order: %w", err)
	}
	log.Printf("    ✅ Created work order ID: %d", workOrderID)

	// Test 3: Create repair work detail
	log.Println("  - Testing repair work detail creation...")
	workDetailID, err := createTestWorkDetail(ctx, db, workOrderID, damageID)
	if err != nil {
		return fmt.Errorf("failed to create test work detail: %w", err)
	}
	log.Printf("    ✅ Created work detail ID: %d", workDetailID)

	// Test 4: Create quality inspection
	log.Println("  - Testing quality inspection creation...")
	inspectionID, err := createTestInspection(ctx, db, workOrderID)
	if err != nil {
		return fmt.Errorf("failed to create test inspection: %w", err)
	}
	log.Printf("    ✅ Created inspection ID: %d", inspectionID)

	return nil
}

func createTestTransaction(ctx context.Context, db *sql.DB) (int, error) {
	query := `
		INSERT INTO vehicle_purchase_transactions (
			transaction_number, customer_id, vehicle_brand, vehicle_model, vehicle_year,
			vehicle_color, purchase_price, agreed_value, odometer_reading, fuel_type,
			transmission, transaction_status, processed_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING transaction_id`

	var transactionID int
	err := db.QueryRowContext(ctx, query,
		"TXN202407TEST001",
		1, // Assuming customer ID 1 exists
		"Honda",
		"Civic",
		2020,
		"Blue",
		15000.00,
		15000.00,
		45000,
		"Gasoline",
		"Automatic",
		"pending",
		1, // Assuming user ID 1 exists
	).Scan(&transactionID)

	return transactionID, err
}

func createTestPayment(ctx context.Context, db *sql.DB, transactionID int) (int, error) {
	query := `
		INSERT INTO vehicle_purchase_payments (
			transaction_id, payment_number, payment_method, payment_amount,
			payment_status, processed_by
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING payment_id`

	var paymentID int
	err := db.QueryRowContext(ctx, query,
		transactionID,
		"PAY202407TEST001",
		"cash",
		15000.00,
		"pending",
		1, // Assuming user ID 1 exists
	).Scan(&paymentID)

	return paymentID, err
}

func updateTransactionStatus(ctx context.Context, db *sql.DB, transactionID int) error {
	query := `
		UPDATE vehicle_purchase_transactions 
		SET transaction_status = $1, updated_at = NOW()
		WHERE transaction_id = $2`

	_, err := db.ExecContext(ctx, query, "approved", transactionID)
	return err
}

func createTestDamage(ctx context.Context, db *sql.DB, transactionID int) (int, error) {
	query := `
		INSERT INTO vehicle_damages (
			transaction_id, damage_category, damage_type, damage_severity,
			damage_location, damage_description, estimated_cost, repair_priority,
			repair_required, identified_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING damage_id`

	var damageID int
	err := db.QueryRowContext(ctx, query,
		transactionID,
		"body",
		"Scratch",
		"minor",
		"Rear bumper",
		"Light scratch on rear bumper",
		150.00,
		2,
		true,
		1, // Assuming user ID 1 exists
	).Scan(&damageID)

	return damageID, err
}

func createTestWorkOrder(ctx context.Context, db *sql.DB, transactionID int) (int, error) {
	query := `
		INSERT INTO repair_work_orders (
			work_order_number, transaction_id, work_order_type, work_order_priority,
			estimated_cost, labor_hours_estimated, work_order_status, work_description,
			created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING work_order_id`

	var workOrderID int
	err := db.QueryRowContext(ctx, query,
		"WO202407TEST001",
		transactionID,
		"repair",
		3,
		500.00,
		8.0,
		"draft",
		"Repair body damages",
		1, // Assuming user ID 1 exists
	).Scan(&workOrderID)

	return workOrderID, err
}

func createTestWorkDetail(ctx context.Context, db *sql.DB, workOrderID, damageID int) (int, error) {
	query := `
		INSERT INTO repair_work_details (
			work_order_id, damage_id, task_sequence, task_description, task_type,
			estimated_hours, labor_rate, task_status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING work_detail_id`

	var workDetailID int
	err := db.QueryRowContext(ctx, query,
		workOrderID,
		damageID,
		1,
		"Sand and repaint rear bumper",
		"repair",
		4.0,
		50.00,
		"pending",
	).Scan(&workDetailID)

	return workDetailID, err
}

func createTestInspection(ctx context.Context, db *sql.DB, workOrderID int) (int, error) {
	query := `
		INSERT INTO quality_inspections (
			work_order_id, inspection_type, inspector_id, overall_rating,
			workmanship_rating, safety_rating, appearance_rating, functionality_rating,
			inspection_status, rework_required
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING inspection_id`

	var inspectionID int
	err := db.QueryRowContext(ctx, query,
		workOrderID,
		"post_repair",
		1, // Assuming user ID 1 exists
		9,
		9,
		10,
		8,
		9,
		"passed",
		false,
	).Scan(&inspectionID)

	return inspectionID, err
}