package interfaces

import (
	"context"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
)

// VehicleBrandRepository defines the interface for vehicle brand data operations
type VehicleBrandRepository interface {
	Create(ctx context.Context, brand *master.VehicleBrand) (*master.VehicleBrand, error)
	GetByID(ctx context.Context, id int) (*master.VehicleBrand, error)
	GetByCode(ctx context.Context, code string) (*master.VehicleBrand, error)
	Update(ctx context.Context, id int, brand *master.VehicleBrand) (*master.VehicleBrand, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params *master.VehicleBrandFilterParams) (*common.PaginatedResponse, error)
	GenerateCode(ctx context.Context) (string, error)
	IsCodeExists(ctx context.Context, code string) (bool, error)
	IsNameExists(ctx context.Context, name string, excludeID int) (bool, error)
}

// VehicleCategoryRepository defines the interface for vehicle category data operations
type VehicleCategoryRepository interface {
	Create(ctx context.Context, category *master.VehicleCategory) (*master.VehicleCategory, error)
	GetByID(ctx context.Context, id int) (*master.VehicleCategory, error)
	GetByCode(ctx context.Context, code string) (*master.VehicleCategory, error)
	Update(ctx context.Context, id int, category *master.VehicleCategory) (*master.VehicleCategory, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params *master.VehicleCategoryFilterParams) (*common.PaginatedResponse, error)
	GenerateCode(ctx context.Context) (string, error)
	IsCodeExists(ctx context.Context, code string) (bool, error)
	IsNameExists(ctx context.Context, name string, excludeID int) (bool, error)
}

// VehicleModelRepository defines the interface for vehicle model data operations
type VehicleModelRepository interface {
	Create(ctx context.Context, model *master.VehicleModel) (*master.VehicleModel, error)
	GetByID(ctx context.Context, id int) (*master.VehicleModel, error)
	GetByCode(ctx context.Context, code string) (*master.VehicleModel, error)
	Update(ctx context.Context, id int, model *master.VehicleModel) (*master.VehicleModel, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params *master.VehicleModelFilterParams) (*common.PaginatedResponse, error)
	GenerateCode(ctx context.Context) (string, error)
	IsCodeExists(ctx context.Context, code string) (bool, error)
}