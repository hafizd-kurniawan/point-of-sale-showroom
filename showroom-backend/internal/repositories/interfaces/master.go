package interfaces

import (
	"context"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
)

// CustomerRepository defines the interface for customer data operations
type CustomerRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, customer *master.Customer) (*master.Customer, error)
	GetByID(ctx context.Context, id int) (*master.Customer, error)
	GetByCode(ctx context.Context, code string) (*master.Customer, error)
	Update(ctx context.Context, id int, customer *master.Customer) (*master.Customer, error)
	Delete(ctx context.Context, id int) error
	
	// List and filtering operations
	List(ctx context.Context, params *master.CustomerFilterParams) ([]master.CustomerListItem, int, error)
	
	// Existence checks
	ExistsByCode(ctx context.Context, code string) (bool, error)
	ExistsByIDCardNumber(ctx context.Context, idCard string) (bool, error)
	ExistsByCodeExcludingID(ctx context.Context, code string, excludeID int) (bool, error)
	ExistsByIDCardNumberExcludingID(ctx context.Context, idCard string, excludeID int) (bool, error)
	
	// Code generation
	GetNextCustomerCode(ctx context.Context) (string, error)
}

// SupplierRepository defines the interface for supplier data operations
type SupplierRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, supplier *master.Supplier) (*master.Supplier, error)
	GetByID(ctx context.Context, id int) (*master.Supplier, error)
	GetByCode(ctx context.Context, code string) (*master.Supplier, error)
	Update(ctx context.Context, id int, supplier *master.Supplier) (*master.Supplier, error)
	Delete(ctx context.Context, id int) error
	
	// List and filtering operations
	List(ctx context.Context, params *master.SupplierFilterParams) ([]master.SupplierListItem, int, error)
	
	// Existence checks
	ExistsByCode(ctx context.Context, code string) (bool, error)
	ExistsByCodeExcludingID(ctx context.Context, code string, excludeID int) (bool, error)
	
	// Code generation
	GetNextSupplierCode(ctx context.Context) (string, error)
}

// VehicleBrandRepository defines the interface for vehicle brand data operations
type VehicleBrandRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, brand *master.VehicleBrand) (*master.VehicleBrand, error)
	GetByID(ctx context.Context, id int) (*master.VehicleBrand, error)
	GetByCode(ctx context.Context, code string) (*master.VehicleBrand, error)
	Update(ctx context.Context, id int, brand *master.VehicleBrand) (*master.VehicleBrand, error)
	Delete(ctx context.Context, id int) error
	
	// List operations
	ListActive(ctx context.Context) ([]master.VehicleBrand, error)
	List(ctx context.Context, isActive *bool) ([]master.VehicleBrand, error)
	
	// Existence checks
	ExistsByCode(ctx context.Context, code string) (bool, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsByCodeExcludingID(ctx context.Context, code string, excludeID int) (bool, error)
	ExistsByNameExcludingID(ctx context.Context, name string, excludeID int) (bool, error)
	
	// Code generation
	GetNextBrandCode(ctx context.Context) (string, error)
}

// VehicleCategoryRepository defines the interface for vehicle category data operations
type VehicleCategoryRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, category *master.VehicleCategory) (*master.VehicleCategory, error)
	GetByID(ctx context.Context, id int) (*master.VehicleCategory, error)
	GetByCode(ctx context.Context, code string) (*master.VehicleCategory, error)
	Update(ctx context.Context, id int, category *master.VehicleCategory) (*master.VehicleCategory, error)
	Delete(ctx context.Context, id int) error
	
	// List operations
	ListActive(ctx context.Context) ([]master.VehicleCategory, error)
	List(ctx context.Context, isActive *bool) ([]master.VehicleCategory, error)
	
	// Existence checks
	ExistsByCode(ctx context.Context, code string) (bool, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsByCodeExcludingID(ctx context.Context, code string, excludeID int) (bool, error)
	ExistsByNameExcludingID(ctx context.Context, name string, excludeID int) (bool, error)
	
	// Code generation
	GetNextCategoryCode(ctx context.Context) (string, error)
}

// VehicleModelRepository defines the interface for vehicle model data operations
type VehicleModelRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, model *master.VehicleModel) (*master.VehicleModel, error)
	GetByID(ctx context.Context, id int) (*master.VehicleModel, error)
	GetByCode(ctx context.Context, code string) (*master.VehicleModel, error)
	Update(ctx context.Context, id int, model *master.VehicleModel) (*master.VehicleModel, error)
	Delete(ctx context.Context, id int) error
	
	// List and filtering operations
	List(ctx context.Context, params *master.VehicleModelFilterParams) ([]master.VehicleModelListItem, int, error)
	
	// Existence checks
	ExistsByCode(ctx context.Context, code string) (bool, error)
	ExistsByCodeExcludingID(ctx context.Context, code string, excludeID int) (bool, error)
	
	// Validation checks
	BrandExists(ctx context.Context, brandID int) (bool, error)
	CategoryExists(ctx context.Context, categoryID int) (bool, error)
	
	// Code generation
	GetNextModelCode(ctx context.Context) (string, error)
}

// ProductCategoryRepository defines the interface for product category data operations
type ProductCategoryRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, category *master.ProductCategory) (*master.ProductCategory, error)
	GetByID(ctx context.Context, id int) (*master.ProductCategory, error)
	GetByCode(ctx context.Context, code string) (*master.ProductCategory, error)
	Update(ctx context.Context, id int, category *master.ProductCategory) (*master.ProductCategory, error)
	Delete(ctx context.Context, id int) error
	
	// List operations
	ListActive(ctx context.Context) ([]master.ProductCategory, error)
	List(ctx context.Context, isActive *bool) ([]master.ProductCategory, error)
	
	// Existence checks
	ExistsByCode(ctx context.Context, code string) (bool, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsByCodeExcludingID(ctx context.Context, code string, excludeID int) (bool, error)
	ExistsByNameExcludingID(ctx context.Context, name string, excludeID int) (bool, error)
	
	// Code generation
	GetNextCategoryCode(ctx context.Context) (string, error)
}