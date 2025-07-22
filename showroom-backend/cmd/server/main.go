package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/config"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/database"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/handlers/admin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/handlers/auth"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/handlers/inventory"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/implementations"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/routes"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services"
	inventoryService "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services/inventory"
	masterService "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services/master"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/utils"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	if err := initializeDatabase(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Initialize dependencies
	dependencies := initializeDependencies(cfg)

	// Setup routes
	router := dependencies.router.SetupRoutes()

	// Configure server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting %s v%s", cfg.App.Name, cfg.App.Version)
		log.Printf("Server running on http://%s:%s", cfg.Server.Host, cfg.Server.Port)
		log.Printf("Environment: %s", cfg.Server.Env)
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// Dependencies holds all application dependencies
type Dependencies struct {
	// Repositories
	userRepo              interfaces.UserRepository
	sessionRepo           interfaces.UserSessionRepository
	customerRepo          interfaces.CustomerRepository
	supplierRepo          interfaces.SupplierRepository
	vehicleBrandRepo      interfaces.VehicleBrandRepository
	vehicleCategoryRepo   interfaces.VehicleCategoryRepository
	vehicleModelRepo      interfaces.VehicleModelRepository
	productCategoryRepo   interfaces.ProductCategoryRepository
	
	// Inventory Repositories
	productRepo            interfaces.ProductSparePartRepository
	purchaseOrderRepo      interfaces.PurchaseOrderRepository
	purchaseOrderDetailRepo interfaces.PurchaseOrderDetailRepository
	goodsReceiptRepo       interfaces.GoodsReceiptRepository
	goodsReceiptDetailRepo interfaces.GoodsReceiptDetailRepository
	stockMovementRepo      interfaces.StockMovementRepository
	stockAdjustmentRepo    interfaces.StockAdjustmentRepository
	supplierPaymentRepo    interfaces.SupplierPaymentRepository
	
	// Services
	authService           *services.AuthService
	userService           *services.UserService
	customerService       *masterService.CustomerService
	supplierService       *masterService.SupplierService
	vehicleBrandService   *masterService.VehicleBrandService
	vehicleCategoryService *masterService.VehicleCategoryService
	vehicleModelService   *masterService.VehicleModelService
	productCategoryService *masterService.ProductCategoryService
	
	// Inventory Services
	productService         *inventoryService.ProductService
	purchaseOrderService   *inventoryService.PurchaseOrderService
	goodsReceiptService    *inventoryService.GoodsReceiptService
	stockMovementService   *inventoryService.StockMovementService
	stockAdjustmentService *inventoryService.StockAdjustmentService
	supplierPaymentService *inventoryService.SupplierPaymentService
	
	// Handlers
	authHandler           *auth.Handler
	adminHandler          *admin.Handler
	customerHandler       *admin.CustomerHandler
	supplierHandler       *admin.SupplierHandler
	vehicleMasterHandler  *admin.VehicleMasterHandler
	productCategoryHandler *admin.ProductCategoryHandler
	
	// Inventory Handlers
	productHandler         *inventory.ProductHandler
	purchaseOrderHandler   *inventory.PurchaseOrderHandler
	goodsReceiptHandler    *inventory.GoodsReceiptHandler
	stockMovementHandler   *inventory.StockMovementHandler
	stockAdjustmentHandler *inventory.StockAdjustmentHandler
	supplierPaymentHandler *inventory.SupplierPaymentHandler
	
	// Utils
	jwtManager            *utils.JWTManager
	router                *routes.Router
}

// initializeDatabase sets up database connection and runs migrations
func initializeDatabase(cfg *config.Config) error {
	// Create database if it doesn't exist
	if err := database.CreateDatabase(cfg); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	if err := database.RunMigrations(database.GetDB()); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// initializeDependencies sets up all application dependencies
func initializeDependencies(cfg *config.Config) *Dependencies {
	db := database.GetDB()

	// Initialize repositories
	userRepo := implementations.NewUserRepository(db)
	sessionRepo := implementations.NewUserSessionRepository(db)
	customerRepo := implementations.NewCustomerRepository(db)
	supplierRepo := implementations.NewSupplierRepository(db)
	vehicleBrandRepo := implementations.NewVehicleBrandRepository(db)
	vehicleCategoryRepo := implementations.NewVehicleCategoryRepository(db)
	vehicleModelRepo := implementations.NewVehicleModelRepository(db)
	productCategoryRepo := implementations.NewProductCategoryRepository(db)

	// Initialize inventory repositories
	productRepo := implementations.NewProductSparePartRepository(db)
	purchaseOrderRepo := implementations.NewPurchaseOrderRepository(db)
	purchaseOrderDetailRepo := implementations.NewPurchaseOrderDetailRepository(db)
	goodsReceiptRepo := implementations.NewGoodsReceiptRepository(db)
	goodsReceiptDetailRepo := implementations.NewGoodsReceiptDetailRepository(db)
	stockMovementRepo := implementations.NewStockMovementRepository(db)
	stockAdjustmentRepo := implementations.NewStockAdjustmentRepository(db)
	supplierPaymentRepo := implementations.NewSupplierPaymentRepository(db)

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(cfg.JWT.SecretKey, cfg.JWT.GetExpiration())

	// Initialize services
	authService := services.NewAuthService(userRepo, sessionRepo, jwtManager)
	userService := services.NewUserService(userRepo, sessionRepo)
	customerService := masterService.NewCustomerService(customerRepo)
	supplierService := masterService.NewSupplierService(supplierRepo)
	vehicleBrandService := masterService.NewVehicleBrandService(vehicleBrandRepo)
	vehicleCategoryService := masterService.NewVehicleCategoryService(vehicleCategoryRepo)
	vehicleModelService := masterService.NewVehicleModelService(vehicleModelRepo, vehicleBrandRepo, vehicleCategoryRepo)
	productCategoryService := masterService.NewProductCategoryService(productCategoryRepo)

	// Initialize inventory services
	productService := inventoryService.NewProductService(productRepo, productCategoryRepo, stockMovementRepo)
	purchaseOrderService := inventoryService.NewPurchaseOrderService(purchaseOrderRepo, purchaseOrderDetailRepo, supplierRepo, productRepo)
	goodsReceiptService := inventoryService.NewGoodsReceiptService(goodsReceiptRepo, goodsReceiptDetailRepo, purchaseOrderRepo, purchaseOrderDetailRepo, stockMovementRepo)
	stockMovementService := inventoryService.NewStockMovementService(stockMovementRepo, productRepo)
	stockAdjustmentService := inventoryService.NewStockAdjustmentService(stockAdjustmentRepo, productRepo, stockMovementRepo, userRepo)
	supplierPaymentService := inventoryService.NewSupplierPaymentService(supplierPaymentRepo, supplierRepo, purchaseOrderRepo)

	// Initialize handlers
	authHandler := auth.NewHandler(authService)
	adminHandler := admin.NewHandler(userService)
	customerHandler := admin.NewCustomerHandler(customerService)
	supplierHandler := admin.NewSupplierHandler(supplierService)
	vehicleMasterHandler := admin.NewVehicleMasterHandler(vehicleBrandService, vehicleCategoryService, vehicleModelService)
	productCategoryHandler := admin.NewProductCategoryHandler(productCategoryService)

	// Initialize inventory handlers
	productHandler := inventory.NewProductHandler(productService)
	purchaseOrderHandler := inventory.NewPurchaseOrderHandler(purchaseOrderService)
	goodsReceiptHandler := inventory.NewGoodsReceiptHandler(goodsReceiptService)
	stockMovementHandler := inventory.NewStockMovementHandler(stockMovementService)
	stockAdjustmentHandler := inventory.NewStockAdjustmentHandler(stockAdjustmentService)
	supplierPaymentHandler := inventory.NewSupplierPaymentHandler(supplierPaymentService)

	// Initialize router
	router := routes.NewRouter(
		authHandler,
		adminHandler,
		customerHandler,
		supplierHandler,
		vehicleMasterHandler,
		productCategoryHandler,
		productHandler,
		purchaseOrderHandler,
		goodsReceiptHandler,
		stockMovementHandler,
		stockAdjustmentHandler,
		supplierPaymentHandler,
		jwtManager,
		sessionRepo,
		cfg,
	)

	return &Dependencies{
		userRepo:               userRepo,
		sessionRepo:            sessionRepo,
		customerRepo:           customerRepo,
		supplierRepo:           supplierRepo,
		vehicleBrandRepo:       vehicleBrandRepo,
		vehicleCategoryRepo:    vehicleCategoryRepo,
		vehicleModelRepo:       vehicleModelRepo,
		productCategoryRepo:    productCategoryRepo,
		productRepo:            productRepo,
		purchaseOrderRepo:      purchaseOrderRepo,
		purchaseOrderDetailRepo: purchaseOrderDetailRepo,
		goodsReceiptRepo:       goodsReceiptRepo,
		goodsReceiptDetailRepo: goodsReceiptDetailRepo,
		stockMovementRepo:      stockMovementRepo,
		stockAdjustmentRepo:    stockAdjustmentRepo,
		supplierPaymentRepo:    supplierPaymentRepo,
		authService:            authService,
		userService:            userService,
		customerService:        customerService,
		supplierService:        supplierService,
		vehicleBrandService:    vehicleBrandService,
		vehicleCategoryService: vehicleCategoryService,
		vehicleModelService:    vehicleModelService,
		productCategoryService: productCategoryService,
		productService:         productService,
		purchaseOrderService:   purchaseOrderService,
		goodsReceiptService:    goodsReceiptService,
		stockMovementService:   stockMovementService,
		stockAdjustmentService: stockAdjustmentService,
		supplierPaymentService: supplierPaymentService,
		authHandler:            authHandler,
		adminHandler:           adminHandler,
		customerHandler:        customerHandler,
		supplierHandler:        supplierHandler,
		vehicleMasterHandler:   vehicleMasterHandler,
		productCategoryHandler: productCategoryHandler,
		productHandler:         productHandler,
		purchaseOrderHandler:   purchaseOrderHandler,
		goodsReceiptHandler:    goodsReceiptHandler,
		stockMovementHandler:   stockMovementHandler,
		stockAdjustmentHandler: stockAdjustmentHandler,
		supplierPaymentHandler: supplierPaymentHandler,
		jwtManager:             jwtManager,
		router:                 router,
	}
}