package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/config"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/handlers/admin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/handlers/auth"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/middleware"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/utils"
)

// Router handles all HTTP routes
type Router struct {
	authHandler            *auth.Handler
	adminHandler           *admin.Handler
	customerHandler        *admin.CustomerHandler
	supplierHandler        *admin.SupplierHandler
	vehicleMasterHandler   *admin.VehicleMasterHandler
	productCategoryHandler *admin.ProductCategoryHandler
	jwtManager             *utils.JWTManager
	sessionRepo            interfaces.UserSessionRepository
	config                 *config.Config
}

// NewRouter creates a new router
func NewRouter(
	authHandler *auth.Handler,
	adminHandler *admin.Handler,
	customerHandler *admin.CustomerHandler,
	supplierHandler *admin.SupplierHandler,
	vehicleMasterHandler *admin.VehicleMasterHandler,
	productCategoryHandler *admin.ProductCategoryHandler,
	jwtManager *utils.JWTManager,
	sessionRepo interfaces.UserSessionRepository,
	config *config.Config,
) *Router {
	return &Router{
		authHandler:            authHandler,
		adminHandler:           adminHandler,
		customerHandler:        customerHandler,
		supplierHandler:        supplierHandler,
		vehicleMasterHandler:   vehicleMasterHandler,
		productCategoryHandler: productCategoryHandler,
		jwtManager:             jwtManager,
		sessionRepo:            sessionRepo,
		config:                 config,
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