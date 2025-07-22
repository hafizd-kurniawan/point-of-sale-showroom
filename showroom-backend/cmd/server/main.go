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
	masterHandler "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/handlers/admin/master"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/handlers/auth"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/implementations"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/routes"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services"
	masterServices "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services/master"
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
	// User repositories
	userRepo    interfaces.UserRepository
	sessionRepo interfaces.UserSessionRepository
	
	// Master data repositories
	customerRepo        interfaces.CustomerRepository
	supplierRepo        interfaces.SupplierRepository
	vehicleBrandRepo    interfaces.VehicleBrandRepository
	vehicleCategoryRepo interfaces.VehicleCategoryRepository
	vehicleModelRepo    interfaces.VehicleModelRepository
	productCategoryRepo interfaces.ProductCategoryRepository
	
	// Services
	authService         *services.AuthService
	userService         *services.UserService
	customerService     *masterServices.CustomerService
	supplierService     *masterServices.SupplierService
	vehicleBrandService *masterServices.VehicleBrandService
	
	// Handlers
	authHandler   *auth.Handler
	adminHandler  *admin.Handler
	masterHandler *masterHandler.Handler
	
	// Utils
	jwtManager *utils.JWTManager
	router     *routes.Router
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
	
	// Initialize master data repositories
	customerRepo := implementations.NewCustomerRepository(db)
	supplierRepo := implementations.NewSupplierRepository(db)
	vehicleBrandRepo := implementations.NewVehicleBrandRepository(db)
	vehicleCategoryRepo := implementations.NewVehicleCategoryRepository(db)
	vehicleModelRepo := implementations.NewVehicleModelRepository(db)
	productCategoryRepo := implementations.NewProductCategoryRepository(db)

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(cfg.JWT.SecretKey, cfg.JWT.GetExpiration())

	// Initialize services
	authService := services.NewAuthService(userRepo, sessionRepo, jwtManager)
	userService := services.NewUserService(userRepo, sessionRepo)
	
	// Initialize master data services
	customerService := masterServices.NewCustomerService(customerRepo)
	supplierService := masterServices.NewSupplierService(supplierRepo)
	vehicleBrandService := masterServices.NewVehicleBrandService(vehicleBrandRepo)

	// Initialize handlers
	authHandler := auth.NewHandler(authService)
	adminHandler := admin.NewHandler(userService)
	masterHandler := masterHandler.NewHandler(customerService, supplierService, vehicleBrandService)

	// Initialize router
	router := routes.NewRouter(authHandler, adminHandler, masterHandler, jwtManager, sessionRepo, cfg)

	return &Dependencies{
		userRepo:            userRepo,
		sessionRepo:         sessionRepo,
		customerRepo:        customerRepo,
		supplierRepo:        supplierRepo,
		vehicleBrandRepo:    vehicleBrandRepo,
		vehicleCategoryRepo: vehicleCategoryRepo,
		vehicleModelRepo:    vehicleModelRepo,
		productCategoryRepo: productCategoryRepo,
		authService:         authService,
		userService:         userService,
		customerService:     customerService,
		supplierService:     supplierService,
		vehicleBrandService: vehicleBrandService,
		authHandler:         authHandler,
		adminHandler:        adminHandler,
		masterHandler:       masterHandler,
		jwtManager:          jwtManager,
		router:              router,
	}
}