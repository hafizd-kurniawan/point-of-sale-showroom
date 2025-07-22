package master

import (
	"context"
	"fmt"
	"strings"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
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

// CreateCustomer creates a new customer
func (s *CustomerService) CreateCustomer(ctx context.Context, req *master.CustomerCreateRequest, createdBy int) (*master.Customer, error) {
	// Validate customer type
	if !req.CustomerType.IsValid() {
		return nil, fmt.Errorf("invalid customer type: %s", req.CustomerType)
	}

	// Check for duplicate phone
	phoneExists, err := s.customerRepo.IsPhoneExists(ctx, req.Phone, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to check phone existence: %w", err)
	}
	if phoneExists {
		return nil, fmt.Errorf("phone number already exists")
	}

	// Check for duplicate email if provided
	if req.Email != nil && *req.Email != "" {
		emailExists, err := s.customerRepo.IsEmailExists(ctx, *req.Email, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to check email existence: %w", err)
		}
		if emailExists {
			return nil, fmt.Errorf("email already exists")
		}
	}

	// Generate customer code
	code, err := s.customerRepo.GenerateCode(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate customer code: %w", err)
	}

	// Create customer entity
	customer := &master.Customer{
		CustomerCode:  code,
		CustomerName:  strings.TrimSpace(req.CustomerName),
		CustomerType:  req.CustomerType,
		Phone:         strings.TrimSpace(req.Phone),
		Email:         req.Email,
		Address:       strings.TrimSpace(req.Address),
		City:          strings.TrimSpace(req.City),
		PostalCode:    req.PostalCode,
		TaxNumber:     req.TaxNumber,
		ContactPerson: req.ContactPerson,
		Notes:         req.Notes,
		CreatedBy:     createdBy,
	}

	return s.customerRepo.Create(ctx, customer)
}

// GetCustomer retrieves a customer by ID
func (s *CustomerService) GetCustomer(ctx context.Context, id int) (*master.Customer, error) {
	return s.customerRepo.GetByID(ctx, id)
}

// GetCustomerByCode retrieves a customer by code
func (s *CustomerService) GetCustomerByCode(ctx context.Context, code string) (*master.Customer, error) {
	return s.customerRepo.GetByCode(ctx, code)
}

// UpdateCustomer updates a customer
func (s *CustomerService) UpdateCustomer(ctx context.Context, id int, req *master.CustomerUpdateRequest) (*master.Customer, error) {
	// Get existing customer
	existing, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate customer type if provided
	if req.CustomerType != nil && !req.CustomerType.IsValid() {
		return nil, fmt.Errorf("invalid customer type: %s", *req.CustomerType)
	}

	// Check for duplicate phone if changed
	if req.Phone != nil && *req.Phone != existing.Phone {
		phoneExists, err := s.customerRepo.IsPhoneExists(ctx, *req.Phone, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check phone existence: %w", err)
		}
		if phoneExists {
			return nil, fmt.Errorf("phone number already exists")
		}
	}

	// Check for duplicate email if changed
	if req.Email != nil && *req.Email != "" {
		currentEmail := ""
		if existing.Email != nil {
			currentEmail = *existing.Email
		}
		if *req.Email != currentEmail {
			emailExists, err := s.customerRepo.IsEmailExists(ctx, *req.Email, id)
			if err != nil {
				return nil, fmt.Errorf("failed to check email existence: %w", err)
			}
			if emailExists {
				return nil, fmt.Errorf("email already exists")
			}
		}
	}

	// Update fields
	updatedCustomer := &master.Customer{
		CustomerCode: existing.CustomerCode,
		CustomerName: existing.CustomerName,
		CustomerType: existing.CustomerType,
		Phone:        existing.Phone,
		Email:        existing.Email,
		Address:      existing.Address,
		City:         existing.City,
		PostalCode:   existing.PostalCode,
		TaxNumber:    existing.TaxNumber,
		ContactPerson: existing.ContactPerson,
		Notes:        existing.Notes,
		IsActive:     existing.IsActive,
		CreatedAt:    existing.CreatedAt,
		CreatedBy:    existing.CreatedBy,
	}

	if req.CustomerName != nil {
		updatedCustomer.CustomerName = strings.TrimSpace(*req.CustomerName)
	}
	if req.CustomerType != nil {
		updatedCustomer.CustomerType = *req.CustomerType
	}
	if req.Phone != nil {
		updatedCustomer.Phone = strings.TrimSpace(*req.Phone)
	}
	if req.Email != nil {
		updatedCustomer.Email = req.Email
	}
	if req.Address != nil {
		updatedCustomer.Address = strings.TrimSpace(*req.Address)
	}
	if req.City != nil {
		updatedCustomer.City = strings.TrimSpace(*req.City)
	}
	if req.PostalCode != nil {
		updatedCustomer.PostalCode = req.PostalCode
	}
	if req.TaxNumber != nil {
		updatedCustomer.TaxNumber = req.TaxNumber
	}
	if req.ContactPerson != nil {
		updatedCustomer.ContactPerson = req.ContactPerson
	}
	if req.Notes != nil {
		updatedCustomer.Notes = req.Notes
	}
	if req.IsActive != nil {
		updatedCustomer.IsActive = *req.IsActive
	}

	return s.customerRepo.Update(ctx, id, updatedCustomer)
}

// DeleteCustomer soft deletes a customer
func (s *CustomerService) DeleteCustomer(ctx context.Context, id int) error {
	return s.customerRepo.Delete(ctx, id)
}

// ListCustomers retrieves customers with filtering and pagination
func (s *CustomerService) ListCustomers(ctx context.Context, params *master.CustomerFilterParams) (*common.PaginatedResponse, error) {
	// Validate pagination parameters
	params.Validate()

	return s.customerRepo.List(ctx, params)
}