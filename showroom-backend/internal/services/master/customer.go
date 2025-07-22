package master

import (
	"context"
	"fmt"
	"time"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// CustomerService handles customer business logic
type CustomerService struct {
	customerRepo interfaces.CustomerRepository
}

// NewCustomerService creates a new customer service
func NewCustomerService(customerRepo interfaces.CustomerRepository) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
	}
}

// Create creates a new customer
func (s *CustomerService) Create(ctx context.Context, req *master.CustomerCreateRequest, createdBy int) (*master.Customer, error) {
	// Validate customer type
	if !req.CustomerType.IsValid() {
		return nil, fmt.Errorf("invalid customer type")
	}

	// Check if ID card number already exists (if provided)
	if req.IDCardNumber != nil && *req.IDCardNumber != "" {
		exists, err := s.customerRepo.ExistsByIDCardNumber(ctx, *req.IDCardNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to check ID card number: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("customer with this ID card number already exists")
		}
	}

	// Generate customer code
	code, err := s.customerRepo.GetNextCustomerCode(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate customer code: %w", err)
	}

	// Parse birth date if provided
	var birthDate *time.Time
	if req.BirthDate != nil && *req.BirthDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			return nil, fmt.Errorf("invalid birth date format, expected YYYY-MM-DD")
		}
		birthDate = &parsed
	}

	// Create customer
	customer := &master.Customer{
		CustomerCode:   code,
		CustomerName:   req.CustomerName,
		Phone:          req.Phone,
		Email:          req.Email,
		Address:        req.Address,
		City:           req.City,
		Province:       req.Province,
		PostalCode:     req.PostalCode,
		IDCardNumber:   req.IDCardNumber,
		TaxNumber:      req.TaxNumber,
		CustomerType:   req.CustomerType,
		BirthDate:      birthDate,
		Occupation:     req.Occupation,
		IncomeRange:    req.IncomeRange,
		CreatedBy:      createdBy,
		IsActive:       true,
		Notes:          req.Notes,
	}

	return s.customerRepo.Create(ctx, customer)
}

// GetByID retrieves a customer by ID
func (s *CustomerService) GetByID(ctx context.Context, id int) (*master.Customer, error) {
	return s.customerRepo.GetByID(ctx, id)
}

// Update updates a customer
func (s *CustomerService) Update(ctx context.Context, id int, req *master.CustomerUpdateRequest) (*master.Customer, error) {
	// Get existing customer
	existing, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate customer type if provided
	if req.CustomerType != nil && !req.CustomerType.IsValid() {
		return nil, fmt.Errorf("invalid customer type")
	}

	// Apply updates
	if req.CustomerName != nil {
		existing.CustomerName = *req.CustomerName
	}
	if req.Phone != nil {
		existing.Phone = *req.Phone
	}
	if req.Email != nil {
		existing.Email = req.Email
	}
	if req.Address != nil {
		existing.Address = req.Address
	}
	if req.City != nil {
		existing.City = req.City
	}
	if req.Province != nil {
		existing.Province = req.Province
	}
	if req.PostalCode != nil {
		existing.PostalCode = req.PostalCode
	}
	if req.CustomerType != nil {
		existing.CustomerType = *req.CustomerType
	}
	if req.Occupation != nil {
		existing.Occupation = req.Occupation
	}
	if req.IncomeRange != nil {
		existing.IncomeRange = req.IncomeRange
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	if req.Notes != nil {
		existing.Notes = req.Notes
	}

	return s.customerRepo.Update(ctx, id, existing)
}

// Delete soft deletes a customer
func (s *CustomerService) Delete(ctx context.Context, id int) error {
	return s.customerRepo.Delete(ctx, id)
}

// List retrieves customers with filtering and pagination
func (s *CustomerService) List(ctx context.Context, params *master.CustomerFilterParams) ([]master.CustomerListItem, *common.PaginationMeta, error) {
	customers, total, err := s.customerRepo.List(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	meta := &common.PaginationMeta{
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: params.GetTotalPages(total),
		HasMore:    params.GetHasMore(total),
	}

	return customers, meta, nil
}