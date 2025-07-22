package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/config"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/database"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/implementations"
	masterService "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services/master"
)

// Helper function to get pointer to string
func stringPtr(s string) *string {
	return &s
}

// Helper function to get pointer to float64
func float64Ptr(f float64) *float64 {
	return &f
}

// Helper function to get pointer to int
func intPtr(i int) *int {
	return &i
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

	// Seed master data
	if err := seedMasterData(); err != nil {
		log.Fatalf("Failed to seed master data: %v", err)
	}

	log.Println("Master data seeding completed successfully!")
}

func initializeDatabase(cfg *config.Config) error {
	if err := database.Connect(cfg); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	return nil
}

func checkAdminUserExists() bool {
	db := database.GetDB()
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin' AND is_active = true").Scan(&count)
	return err == nil && count > 0
}

func getAdminUserID() (int, error) {
	db := database.GetDB()
	var adminID int
	err := db.QueryRow("SELECT user_id FROM users WHERE role = 'admin' AND is_active = true LIMIT 1").Scan(&adminID)
	if err != nil {
		return 0, fmt.Errorf("failed to get admin user ID: %w", err)
	}
	return adminID, nil
}

func seedMasterData() error {
	ctx := context.Background()
	db := database.GetDB()

	// Get admin user ID for created_by field
	adminID, err := getAdminUserID()
	if err != nil {
		return err
	}

	// Initialize repositories and services
	customerRepo := implementations.NewCustomerRepository(db)
	supplierRepo := implementations.NewSupplierRepository(db)
	vehicleBrandRepo := implementations.NewVehicleBrandRepository(db)
	vehicleCategoryRepo := implementations.NewVehicleCategoryRepository(db)
	vehicleModelRepo := implementations.NewVehicleModelRepository(db)
	productCategoryRepo := implementations.NewProductCategoryRepository(db)

	customerService := masterService.NewCustomerService(customerRepo)
	supplierService := masterService.NewSupplierService(supplierRepo)
	vehicleBrandService := masterService.NewVehicleBrandService(vehicleBrandRepo)
	vehicleCategoryService := masterService.NewVehicleCategoryService(vehicleCategoryRepo)
	vehicleModelService := masterService.NewVehicleModelService(vehicleModelRepo, vehicleBrandRepo, vehicleCategoryRepo)
	productCategoryService := masterService.NewProductCategoryService(productCategoryRepo)

	// Seed customers
	log.Println("Seeding customers...")
	customers := []master.CustomerCreateRequest{
		{
			CustomerName: "John Doe",
			CustomerType: master.CustomerTypeIndividual,
			Phone:        "081234567890",
			Email:        stringPtr("john.doe@email.com"),
			Address:      "Jl. Sudirman No. 123",
			City:         "Jakarta",
			PostalCode:   stringPtr("12345"),
			Notes:        stringPtr("Regular customer"),
		},
		{
			CustomerName: "PT. Auto Prima",
			CustomerType: master.CustomerTypeCorporate,
			Phone:        "021-5551234",
			Email:        stringPtr("info@autoprima.com"),
			Address:      "Jl. Gatot Subroto No. 456",
			City:         "Jakarta",
			PostalCode:   stringPtr("12940"),
			TaxNumber:    stringPtr("01.234.567.8-901.000"),
			ContactPerson: stringPtr("Budi Santoso"),
			Notes:        stringPtr("Corporate fleet customer"),
		},
		{
			CustomerName: "Jane Smith",
			CustomerType: master.CustomerTypeIndividual,
			Phone:        "082345678901",
			Email:        stringPtr("jane.smith@email.com"),
			Address:      "Jl. Thamrin No. 789",
			City:         "Bandung",
			PostalCode:   stringPtr("40111"),
		},
	}

	for _, customer := range customers {
		_, err := customerService.CreateCustomer(ctx, &customer, adminID)
		if err != nil {
			log.Printf("Warning: Failed to create customer %s: %v", customer.CustomerName, err)
		}
	}

	// Seed suppliers
	log.Println("Seeding suppliers...")
	suppliers := []master.SupplierCreateRequest{
		{
			SupplierName: "PT. Spare Part Indonesia",
			SupplierType: master.SupplierTypeParts,
			Phone:        "021-7778888",
			Email:        stringPtr("sales@sparepart.co.id"),
			Address:      "Jl. Raya Industri No. 100",
			City:         "Tangerang",
			PostalCode:   stringPtr("15111"),
			TaxNumber:    stringPtr("02.345.678.9-012.000"),
			ContactPerson: "Ahmad Yusuf",
			BankAccount:  stringPtr("BCA 1234567890"),
			PaymentTerms: stringPtr("Net 30 days"),
			Notes:        stringPtr("Main spare parts supplier"),
		},
		{
			SupplierName: "CV. Mobil Berkah",
			SupplierType: master.SupplierTypeVehicle,
			Phone:        "021-9990000",
			Email:        stringPtr("info@mobilberkah.com"),
			Address:      "Jl. Boulevard No. 200",
			City:         "Jakarta",
			PostalCode:   stringPtr("13450"),
			TaxNumber:    stringPtr("03.456.789.0-123.000"),
			ContactPerson: "Siti Nurhaliza",
			BankAccount:  stringPtr("Mandiri 0987654321"),
			PaymentTerms: stringPtr("COD"),
		},
		{
			SupplierName: "PT. Universal Auto",
			SupplierType: master.SupplierTypeBoth,
			Phone:        "021-1112222",
			Email:        stringPtr("procurement@universalauto.id"),
			Address:      "Jl. Margonda Raya No. 300",
			City:         "Depok",
			PostalCode:   stringPtr("16424"),
			TaxNumber:    stringPtr("04.567.890.1-234.000"),
			ContactPerson: "Andi Wijaya",
			BankAccount:  stringPtr("BNI 1122334455"),
			PaymentTerms: stringPtr("Net 15 days"),
		},
	}

	for _, supplier := range suppliers {
		_, err := supplierService.CreateSupplier(ctx, &supplier, adminID)
		if err != nil {
			log.Printf("Warning: Failed to create supplier %s: %v", supplier.SupplierName, err)
		}
	}

	// Seed vehicle brands
	log.Println("Seeding vehicle brands...")
	brands := []master.VehicleBrandCreateRequest{
		{
			BrandName:     "Toyota",
			CountryOrigin: "Japan",
			Description:   stringPtr("World's leading automotive manufacturer"),
		},
		{
			BrandName:     "Honda",
			CountryOrigin: "Japan",
			Description:   stringPtr("Reliable and fuel-efficient vehicles"),
		},
		{
			BrandName:     "Mitsubishi",
			CountryOrigin: "Japan",
			Description:   stringPtr("Innovative automotive technology"),
		},
		{
			BrandName:     "Suzuki",
			CountryOrigin: "Japan",
			Description:   stringPtr("Compact and efficient vehicles"),
		},
		{
			BrandName:     "Daihatsu",
			CountryOrigin: "Japan",
			Description:   stringPtr("Small cars specialist"),
		},
	}

	var createdBrands []*master.VehicleBrand
	for _, brand := range brands {
		created, err := vehicleBrandService.CreateVehicleBrand(ctx, &brand, adminID)
		if err != nil {
			log.Printf("Warning: Failed to create brand %s: %v", brand.BrandName, err)
		} else {
			createdBrands = append(createdBrands, created)
		}
	}

	// Seed vehicle categories
	log.Println("Seeding vehicle categories...")
	categories := []master.VehicleCategoryCreateRequest{
		{
			CategoryName: "Sedan",
			Description:  stringPtr("Four-door passenger car with separate trunk"),
		},
		{
			CategoryName: "SUV",
			Description:  stringPtr("Sport Utility Vehicle for multiple purposes"),
		},
		{
			CategoryName: "Hatchback",
			Description:  stringPtr("Compact car with rear door and no separate trunk"),
		},
		{
			CategoryName: "MPV",
			Description:  stringPtr("Multi-Purpose Vehicle for large families"),
		},
		{
			CategoryName: "Pickup",
			Description:  stringPtr("Light truck with open cargo area"),
		},
	}

	var createdCategories []*master.VehicleCategory
	for _, category := range categories {
		created, err := vehicleCategoryService.CreateVehicleCategory(ctx, &category, adminID)
		if err != nil {
			log.Printf("Warning: Failed to create category %s: %v", category.CategoryName, err)
		} else {
			createdCategories = append(createdCategories, created)
		}
	}

	// Seed vehicle models (if we have brands and categories)
	if len(createdBrands) > 0 && len(createdCategories) > 0 {
		log.Println("Seeding vehicle models...")
		models := []master.VehicleModelCreateRequest{
			{
				ModelName:      "Avanza",
				BrandID:        createdBrands[0].BrandID, // Toyota
				CategoryID:     createdCategories[3].CategoryID, // MPV
				ModelYear:      2023,
				EngineCapacity: float64Ptr(1.3),
				FuelType:       "Gasoline",
				Transmission:   "Manual",
				SeatCapacity:   7,
				Color:          "White",
				Price:          220000000,
				Description:    stringPtr("Popular family MPV"),
			},
			{
				ModelName:      "Civic",
				BrandID:        createdBrands[1].BrandID, // Honda
				CategoryID:     createdCategories[0].CategoryID, // Sedan
				ModelYear:      2023,
				EngineCapacity: float64Ptr(1.5),
				FuelType:       "Gasoline",
				Transmission:   "CVT",
				SeatCapacity:   5,
				Color:          "Black",
				Price:          550000000,
				Description:    stringPtr("Premium sedan with advanced features"),
			},
			{
				ModelName:      "Xpander",
				BrandID:        createdBrands[2].BrandID, // Mitsubishi
				CategoryID:     createdCategories[3].CategoryID, // MPV
				ModelYear:      2023,
				EngineCapacity: float64Ptr(1.5),
				FuelType:       "Gasoline",
				Transmission:   "CVT",
				SeatCapacity:   7,
				Color:          "Silver",
				Price:          280000000,
				Description:    stringPtr("Modern and stylish MPV"),
			},
		}

		for _, model := range models {
			_, err := vehicleModelService.CreateVehicleModel(ctx, &model, adminID)
			if err != nil {
				log.Printf("Warning: Failed to create model %s: %v", model.ModelName, err)
			}
		}
	}

	// Seed product categories
	log.Println("Seeding product categories...")
	rootCategories := []master.ProductCategoryCreateRequest{
		{
			CategoryName: "Engine Parts",
			Description:  stringPtr("Engine components and accessories"),
		},
		{
			CategoryName: "Body Parts",
			Description:  stringPtr("Exterior and interior body components"),
		},
		{
			CategoryName: "Electrical Parts",
			Description:  stringPtr("Electrical and electronic components"),
		},
		{
			CategoryName: "Brake System",
			Description:  stringPtr("Brake components and accessories"),
		},
	}

	var createdRootCategories []*master.ProductCategory
	for _, category := range rootCategories {
		created, err := productCategoryService.CreateProductCategory(ctx, &category, adminID)
		if err != nil {
			log.Printf("Warning: Failed to create product category %s: %v", category.CategoryName, err)
		} else {
			createdRootCategories = append(createdRootCategories, created)
		}
	}

	// Seed sub-categories
	if len(createdRootCategories) > 0 {
		log.Println("Seeding product sub-categories...")
		subCategories := []master.ProductCategoryCreateRequest{
			{
				CategoryName: "Oil Filter",
				Description:  stringPtr("Engine oil filtration components"),
				ParentID:     &createdRootCategories[0].CategoryID, // Engine Parts
			},
			{
				CategoryName: "Air Filter",
				Description:  stringPtr("Air intake filtration components"),
				ParentID:     &createdRootCategories[0].CategoryID, // Engine Parts
			},
			{
				CategoryName: "Bumper",
				Description:  stringPtr("Front and rear bumper components"),
				ParentID:     &createdRootCategories[1].CategoryID, // Body Parts
			},
			{
				CategoryName: "Headlight",
				Description:  stringPtr("Front lighting components"),
				ParentID:     &createdRootCategories[2].CategoryID, // Electrical Parts
			},
			{
				CategoryName: "Brake Pad",
				Description:  stringPtr("Brake friction components"),
				ParentID:     &createdRootCategories[3].CategoryID, // Brake System
			},
		}

		for _, category := range subCategories {
			_, err := productCategoryService.CreateProductCategory(ctx, &category, adminID)
			if err != nil {
				log.Printf("Warning: Failed to create sub-category %s: %v", category.CategoryName, err)
			}
		}
	}

	log.Println("Master data seeding completed!")
	return nil
}