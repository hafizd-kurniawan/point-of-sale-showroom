package interfaces

import (
	"context"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
)

// SupplierRepository defines the interface for supplier data operations
type SupplierRepository interface {
	Create(ctx context.Context, supplier *master.Supplier) (*master.Supplier, error)
	GetByID(ctx context.Context, id int) (*master.Supplier, error)
	GetByCode(ctx context.Context, code string) (*master.Supplier, error)
	Update(ctx context.Context, id int, supplier *master.Supplier) (*master.Supplier, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params *master.SupplierFilterParams) (*common.PaginatedResponse, error)
	GenerateCode(ctx context.Context) (string, error)
	IsCodeExists(ctx context.Context, code string) (bool, error)
	IsEmailExists(ctx context.Context, email string, excludeID int) (bool, error)
	IsPhoneExists(ctx context.Context, phone string, excludeID int) (bool, error)
}