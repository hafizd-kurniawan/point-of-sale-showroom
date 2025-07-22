package master

import (
	"context"
	"fmt"

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

// Create creates a new vehicle brand
func (s *VehicleBrandService) Create(ctx context.Context, req *master.VehicleBrandCreateRequest, createdBy int) (*master.VehicleBrand, error) {
	// Check if brand name already exists
	exists, err := s.brandRepo.ExistsByName(ctx, req.BrandName)
	if err != nil {
		return nil, fmt.Errorf("failed to check brand name: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("vehicle brand with this name already exists")
	}

	// Generate brand code
	code, err := s.brandRepo.GetNextBrandCode(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate brand code: %w", err)
	}

	// Create brand
	brand := &master.VehicleBrand{
		BrandCode:     code,
		BrandName:     req.BrandName,
		CountryOrigin: req.CountryOrigin,
		LogoImage:     req.LogoImage,
		CreatedBy:     createdBy,
		IsActive:      true,
	}

	return s.brandRepo.Create(ctx, brand)
}

// GetByID retrieves a vehicle brand by ID
func (s *VehicleBrandService) GetByID(ctx context.Context, id int) (*master.VehicleBrand, error) {
	return s.brandRepo.GetByID(ctx, id)
}

// Update updates a vehicle brand
func (s *VehicleBrandService) Update(ctx context.Context, id int, req *master.VehicleBrandUpdateRequest) (*master.VehicleBrand, error) {
	// Get existing brand
	existing, err := s.brandRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if brand name already exists (excluding current brand)
	if req.BrandName != nil {
		exists, err := s.brandRepo.ExistsByNameExcludingID(ctx, *req.BrandName, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check brand name: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("vehicle brand with this name already exists")
		}
		existing.BrandName = *req.BrandName
	}

	// Apply updates
	if req.CountryOrigin != nil {
		existing.CountryOrigin = req.CountryOrigin
	}
	if req.LogoImage != nil {
		existing.LogoImage = req.LogoImage
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	return s.brandRepo.Update(ctx, id, existing)
}

// Delete soft deletes a vehicle brand
func (s *VehicleBrandService) Delete(ctx context.Context, id int) error {
	return s.brandRepo.Delete(ctx, id)
}

// ListActive retrieves all active vehicle brands
func (s *VehicleBrandService) ListActive(ctx context.Context) ([]master.VehicleBrand, error) {
	return s.brandRepo.ListActive(ctx)
}

// List retrieves vehicle brands with optional filtering
func (s *VehicleBrandService) List(ctx context.Context, isActive *bool) ([]master.VehicleBrand, error) {
	return s.brandRepo.List(ctx, isActive)
}