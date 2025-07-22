package master

import (
	"context"
	"fmt"
	"strings"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// SupplierService handles supplier business logic
type SupplierService struct {
	supplierRepo interfaces.SupplierRepository
}

// NewSupplierService creates a new supplier service
func NewSupplierService(supplierRepo interfaces.SupplierRepository) *SupplierService {
	return &SupplierService{
		supplierRepo: supplierRepo,
	}
}

// CreateSupplier creates a new supplier
func (s *SupplierService) CreateSupplier(ctx context.Context, req *master.SupplierCreateRequest, createdBy int) (*master.Supplier, error) {
	// Validate supplier type
	if !req.SupplierType.IsValid() {
		return nil, fmt.Errorf("invalid supplier type: %s", req.SupplierType)
	}

	// Check for duplicate phone
	phoneExists, err := s.supplierRepo.IsPhoneExists(ctx, req.Phone, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to check phone existence: %w", err)
	}
	if phoneExists {
		return nil, fmt.Errorf("phone number already exists")
	}

	// Check for duplicate email if provided
	if req.Email != nil && *req.Email != "" {
		emailExists, err := s.supplierRepo.IsEmailExists(ctx, *req.Email, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to check email existence: %w", err)
		}
		if emailExists {
			return nil, fmt.Errorf("email already exists")
		}
	}

	// Generate supplier code
	code, err := s.supplierRepo.GenerateCode(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate supplier code: %w", err)
	}

	// Create supplier entity
	supplier := &master.Supplier{
		SupplierCode:  code,
		SupplierName:  strings.TrimSpace(req.SupplierName),
		SupplierType:  req.SupplierType,
		Phone:         strings.TrimSpace(req.Phone),
		Email:         req.Email,
		Address:       strings.TrimSpace(req.Address),
		City:          strings.TrimSpace(req.City),
		PostalCode:    req.PostalCode,
		TaxNumber:     req.TaxNumber,
		ContactPerson: strings.TrimSpace(req.ContactPerson),
		BankAccount:   req.BankAccount,
		PaymentTerms:  req.PaymentTerms,
		Notes:         req.Notes,
		CreatedBy:     createdBy,
	}

	return s.supplierRepo.Create(ctx, supplier)
}

// GetSupplier retrieves a supplier by ID
func (s *SupplierService) GetSupplier(ctx context.Context, id int) (*master.Supplier, error) {
	return s.supplierRepo.GetByID(ctx, id)
}

// GetSupplierByCode retrieves a supplier by code
func (s *SupplierService) GetSupplierByCode(ctx context.Context, code string) (*master.Supplier, error) {
	return s.supplierRepo.GetByCode(ctx, code)
}

// UpdateSupplier updates a supplier
func (s *SupplierService) UpdateSupplier(ctx context.Context, id int, req *master.SupplierUpdateRequest) (*master.Supplier, error) {
	// Get existing supplier
	existing, err := s.supplierRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate supplier type if provided
	if req.SupplierType != nil && !req.SupplierType.IsValid() {
		return nil, fmt.Errorf("invalid supplier type: %s", *req.SupplierType)
	}

	// Check for duplicate phone if changed
	if req.Phone != nil && *req.Phone != existing.Phone {
		phoneExists, err := s.supplierRepo.IsPhoneExists(ctx, *req.Phone, id)
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
			emailExists, err := s.supplierRepo.IsEmailExists(ctx, *req.Email, id)
			if err != nil {
				return nil, fmt.Errorf("failed to check email existence: %w", err)
			}
			if emailExists {
				return nil, fmt.Errorf("email already exists")
			}
		}
	}

	// Update fields
	updatedSupplier := &master.Supplier{
		SupplierCode:  existing.SupplierCode,
		SupplierName:  existing.SupplierName,
		SupplierType:  existing.SupplierType,
		Phone:         existing.Phone,
		Email:         existing.Email,
		Address:       existing.Address,
		City:          existing.City,
		PostalCode:    existing.PostalCode,
		TaxNumber:     existing.TaxNumber,
		ContactPerson: existing.ContactPerson,
		BankAccount:   existing.BankAccount,
		PaymentTerms:  existing.PaymentTerms,
		Notes:         existing.Notes,
		IsActive:      existing.IsActive,
		CreatedAt:     existing.CreatedAt,
		CreatedBy:     existing.CreatedBy,
	}

	if req.SupplierName != nil {
		updatedSupplier.SupplierName = strings.TrimSpace(*req.SupplierName)
	}
	if req.SupplierType != nil {
		updatedSupplier.SupplierType = *req.SupplierType
	}
	if req.Phone != nil {
		updatedSupplier.Phone = strings.TrimSpace(*req.Phone)
	}
	if req.Email != nil {
		updatedSupplier.Email = req.Email
	}
	if req.Address != nil {
		updatedSupplier.Address = strings.TrimSpace(*req.Address)
	}
	if req.City != nil {
		updatedSupplier.City = strings.TrimSpace(*req.City)
	}
	if req.PostalCode != nil {
		updatedSupplier.PostalCode = req.PostalCode
	}
	if req.TaxNumber != nil {
		updatedSupplier.TaxNumber = req.TaxNumber
	}
	if req.ContactPerson != nil {
		updatedSupplier.ContactPerson = strings.TrimSpace(*req.ContactPerson)
	}
	if req.BankAccount != nil {
		updatedSupplier.BankAccount = req.BankAccount
	}
	if req.PaymentTerms != nil {
		updatedSupplier.PaymentTerms = req.PaymentTerms
	}
	if req.Notes != nil {
		updatedSupplier.Notes = req.Notes
	}
	if req.IsActive != nil {
		updatedSupplier.IsActive = *req.IsActive
	}

	return s.supplierRepo.Update(ctx, id, updatedSupplier)
}

// DeleteSupplier soft deletes a supplier
func (s *SupplierService) DeleteSupplier(ctx context.Context, id int) error {
	return s.supplierRepo.Delete(ctx, id)
}

// ListSuppliers retrieves suppliers with filtering and pagination
func (s *SupplierService) ListSuppliers(ctx context.Context, params *master.SupplierFilterParams) (*common.PaginatedResponse, error) {
	// Validate pagination parameters
	params.Validate()

	return s.supplierRepo.List(ctx, params)
}