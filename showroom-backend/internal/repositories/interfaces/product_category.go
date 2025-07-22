package interfaces

import (
	"context"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
)

// ProductCategoryRepository defines the interface for product category data operations
type ProductCategoryRepository interface {
	Create(ctx context.Context, category *master.ProductCategory) (*master.ProductCategory, error)
	GetByID(ctx context.Context, id int) (*master.ProductCategory, error)
	GetByCode(ctx context.Context, code string) (*master.ProductCategory, error)
	Update(ctx context.Context, id int, category *master.ProductCategory) (*master.ProductCategory, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params *master.ProductCategoryFilterParams) (*common.PaginatedResponse, error)
	GetTree(ctx context.Context) ([]master.ProductCategoryTree, error)
	GetChildren(ctx context.Context, parentID int) ([]master.ProductCategory, error)
	GenerateCode(ctx context.Context) (string, error)
	IsCodeExists(ctx context.Context, code string) (bool, error)
	IsNameExists(ctx context.Context, name string, parentID *int, excludeID int) (bool, error)
	UpdatePath(ctx context.Context, categoryID int, path string) error
	UpdateLevel(ctx context.Context, categoryID int, level int) error
}