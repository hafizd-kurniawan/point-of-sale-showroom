package master

import (
	"context"
	"fmt"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
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

// Create creates a new supplier
func (s *SupplierService) Create(ctx context.Context, req *master.SupplierCreateRequest, createdBy int) (*master.Supplier, error) {
	// Validate supplier type
	if !req.SupplierType.IsValid() {
		return nil, fmt.Errorf("invalid supplier type")
	}

	// Generate supplier code
	code, err := s.supplierRepo.GetNextSupplierCode(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate supplier code: %w", err)
	}

	// Create supplier
	supplier := &master.Supplier{
		SupplierCode:    code,
		SupplierName:    req.SupplierName,
		ContactPerson:   req.ContactPerson,
		Phone:           req.Phone,
		Email:           req.Email,
		Address:         req.Address,
		City:            req.City,
		Province:        req.Province,
		PostalCode:      req.PostalCode,
		TaxNumber:       req.TaxNumber,
		SupplierType:    req.SupplierType,
		CreditLimit:     req.CreditLimit,
		CreditTermDays:  req.CreditTermDays,
		CreatedBy:       createdBy,
		IsActive:        true,
		Notes:           req.Notes,
	}

	return s.supplierRepo.Create(ctx, supplier)
}

// GetByID retrieves a supplier by ID
func (s *SupplierService) GetByID(ctx context.Context, id int) (*master.Supplier, error) {
	return s.supplierRepo.GetByID(ctx, id)
}

// Update updates a supplier
func (s *SupplierService) Update(ctx context.Context, id int, req *master.SupplierUpdateRequest) (*master.Supplier, error) {
	// Get existing supplier
	existing, err := s.supplierRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate supplier type if provided
	if req.SupplierType != nil && !req.SupplierType.IsValid() {
		return nil, fmt.Errorf("invalid supplier type")
	}

	// Apply updates
	if req.SupplierName != nil {
		existing.SupplierName = *req.SupplierName
	}
	if req.ContactPerson != nil {
		existing.ContactPerson = *req.ContactPerson
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
	if req.SupplierType != nil {
		existing.SupplierType = *req.SupplierType
	}
	if req.CreditLimit != nil {
		existing.CreditLimit = req.CreditLimit
	}
	if req.CreditTermDays != nil {
		existing.CreditTermDays = req.CreditTermDays
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	if req.Notes != nil {
		existing.Notes = req.Notes
	}

	return s.supplierRepo.Update(ctx, id, existing)
}

// Delete soft deletes a supplier
func (s *SupplierService) Delete(ctx context.Context, id int) error {
	return s.supplierRepo.Delete(ctx, id)
}

// List retrieves suppliers with filtering and pagination
func (s *SupplierService) List(ctx context.Context, params *master.SupplierFilterParams) ([]master.SupplierListItem, *common.PaginationMeta, error) {
	suppliers, total, err := s.supplierRepo.List(ctx, params)
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

	return suppliers, meta, nil
}