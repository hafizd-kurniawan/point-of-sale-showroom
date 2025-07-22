package master

import (
	"context"
	"fmt"
	"strings"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// VehicleBrandService handles vehicle brand business logic
type VehicleBrandService struct {
	brandRepo interfaces.VehicleBrandRepository
}

// NewVehicleBrandService creates a new vehicle brand service
func NewVehicleBrandService(brandRepo interfaces.VehicleBrandRepository) *VehicleBrandService {
	return &VehicleBrandService{
		brandRepo: brandRepo,
	}
}

// CreateVehicleBrand creates a new vehicle brand
func (s *VehicleBrandService) CreateVehicleBrand(ctx context.Context, req *master.VehicleBrandCreateRequest, createdBy int) (*master.VehicleBrand, error) {
	// Check for duplicate brand name
	nameExists, err := s.brandRepo.IsNameExists(ctx, req.BrandName, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to check brand name existence: %w", err)
	}
	if nameExists {
		return nil, fmt.Errorf("brand name already exists")
	}

	// Generate brand code
	code, err := s.brandRepo.GenerateCode(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate brand code: %w", err)
	}

	// Create brand entity
	brand := &master.VehicleBrand{
		BrandCode:     code,
		BrandName:     strings.TrimSpace(req.BrandName),
		CountryOrigin: strings.TrimSpace(req.CountryOrigin),
		Description:   req.Description,
		LogoURL:       req.LogoURL,
		CreatedBy:     createdBy,
	}

	return s.brandRepo.Create(ctx, brand)
}

// GetVehicleBrand retrieves a vehicle brand by ID
func (s *VehicleBrandService) GetVehicleBrand(ctx context.Context, id int) (*master.VehicleBrand, error) {
	return s.brandRepo.GetByID(ctx, id)
}

// GetVehicleBrandByCode retrieves a vehicle brand by code
func (s *VehicleBrandService) GetVehicleBrandByCode(ctx context.Context, code string) (*master.VehicleBrand, error) {
	return s.brandRepo.GetByCode(ctx, code)
}

// UpdateVehicleBrand updates a vehicle brand
func (s *VehicleBrandService) UpdateVehicleBrand(ctx context.Context, id int, req *master.VehicleBrandUpdateRequest) (*master.VehicleBrand, error) {
	// Get existing brand
	existing, err := s.brandRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check for duplicate brand name if changed
	if req.BrandName != nil && *req.BrandName != existing.BrandName {
		nameExists, err := s.brandRepo.IsNameExists(ctx, *req.BrandName, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check brand name existence: %w", err)
		}
		if nameExists {
			return nil, fmt.Errorf("brand name already exists")
		}
	}

	// Update fields
	updatedBrand := &master.VehicleBrand{
		BrandCode:     existing.BrandCode,
		BrandName:     existing.BrandName,
		CountryOrigin: existing.CountryOrigin,
		Description:   existing.Description,
		LogoURL:       existing.LogoURL,
		IsActive:      existing.IsActive,
		CreatedAt:     existing.CreatedAt,
		CreatedBy:     existing.CreatedBy,
	}

	if req.BrandName != nil {
		updatedBrand.BrandName = strings.TrimSpace(*req.BrandName)
	}
	if req.CountryOrigin != nil {
		updatedBrand.CountryOrigin = strings.TrimSpace(*req.CountryOrigin)
	}
	if req.Description != nil {
		updatedBrand.Description = req.Description
	}
	if req.LogoURL != nil {
		updatedBrand.LogoURL = req.LogoURL
	}
	if req.IsActive != nil {
		updatedBrand.IsActive = *req.IsActive
	}

	return s.brandRepo.Update(ctx, id, updatedBrand)
}

// DeleteVehicleBrand soft deletes a vehicle brand
func (s *VehicleBrandService) DeleteVehicleBrand(ctx context.Context, id int) error {
	return s.brandRepo.Delete(ctx, id)
}

// ListVehicleBrands retrieves vehicle brands with filtering and pagination
func (s *VehicleBrandService) ListVehicleBrands(ctx context.Context, params *master.VehicleBrandFilterParams) (*common.PaginatedResponse, error) {
	// Validate pagination parameters
	params.Validate()

	return s.brandRepo.List(ctx, params)
}