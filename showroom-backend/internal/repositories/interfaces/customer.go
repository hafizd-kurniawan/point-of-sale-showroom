package interfaces

import (
	"context"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
)

// CustomerRepository defines the interface for customer data operations
type CustomerRepository interface {
	Create(ctx context.Context, customer *master.Customer) (*master.Customer, error)
	GetByID(ctx context.Context, id int) (*master.Customer, error)
	GetByCode(ctx context.Context, code string) (*master.Customer, error)
	Update(ctx context.Context, id int, customer *master.Customer) (*master.Customer, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, params *master.CustomerFilterParams) (*common.PaginatedResponse, error)
	GenerateCode(ctx context.Context) (string, error)
	IsCodeExists(ctx context.Context, code string) (bool, error)
	IsEmailExists(ctx context.Context, email string, excludeID int) (bool, error)
	IsPhoneExists(ctx context.Context, phone string, excludeID int) (bool, error)
}