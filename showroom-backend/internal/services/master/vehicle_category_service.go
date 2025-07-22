package master

import (
	"context"
	"fmt"
	"strings"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// VehicleCategoryService handles vehicle category business logic
type VehicleCategoryService struct {
	categoryRepo interfaces.VehicleCategoryRepository
}

// NewVehicleCategoryService creates a new vehicle category service
func NewVehicleCategoryService(categoryRepo interfaces.VehicleCategoryRepository) *VehicleCategoryService {
	return &VehicleCategoryService{
		categoryRepo: categoryRepo,
	}
}

// CreateVehicleCategory creates a new vehicle category
func (s *VehicleCategoryService) CreateVehicleCategory(ctx context.Context, req *master.VehicleCategoryCreateRequest, createdBy int) (*master.VehicleCategory, error) {
	// Check for duplicate category name
	nameExists, err := s.categoryRepo.IsNameExists(ctx, req.CategoryName, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to check category name existence: %w", err)
	}
	if nameExists {
		return nil, fmt.Errorf("category name already exists")
	}

	// Generate category code
	code, err := s.categoryRepo.GenerateCode(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate category code: %w", err)
	}

	// Create category entity
	category := &master.VehicleCategory{
		CategoryCode: code,
		CategoryName: strings.TrimSpace(req.CategoryName),
		Description:  req.Description,
		CreatedBy:    createdBy,
	}

	return s.categoryRepo.Create(ctx, category)
}

// GetVehicleCategory retrieves a vehicle category by ID
func (s *VehicleCategoryService) GetVehicleCategory(ctx context.Context, id int) (*master.VehicleCategory, error) {
	return s.categoryRepo.GetByID(ctx, id)
}

// GetVehicleCategoryByCode retrieves a vehicle category by code
func (s *VehicleCategoryService) GetVehicleCategoryByCode(ctx context.Context, code string) (*master.VehicleCategory, error) {
	return s.categoryRepo.GetByCode(ctx, code)
}

// UpdateVehicleCategory updates a vehicle category
func (s *VehicleCategoryService) UpdateVehicleCategory(ctx context.Context, id int, req *master.VehicleCategoryUpdateRequest) (*master.VehicleCategory, error) {
	// Get existing category
	existing, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check for duplicate category name if changed
	if req.CategoryName != nil && *req.CategoryName != existing.CategoryName {
		nameExists, err := s.categoryRepo.IsNameExists(ctx, *req.CategoryName, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check category name existence: %w", err)
		}
		if nameExists {
			return nil, fmt.Errorf("category name already exists")
		}
	}

	// Update fields
	updatedCategory := &master.VehicleCategory{
		CategoryCode: existing.CategoryCode,
		CategoryName: existing.CategoryName,
		Description:  existing.Description,
		IsActive:     existing.IsActive,
		CreatedAt:    existing.CreatedAt,
		CreatedBy:    existing.CreatedBy,
	}

	if req.CategoryName != nil {
		updatedCategory.CategoryName = strings.TrimSpace(*req.CategoryName)
	}
	if req.Description != nil {
		updatedCategory.Description = req.Description
	}
	if req.IsActive != nil {
		updatedCategory.IsActive = *req.IsActive
	}

	return s.categoryRepo.Update(ctx, id, updatedCategory)
}

// DeleteVehicleCategory soft deletes a vehicle category
func (s *VehicleCategoryService) DeleteVehicleCategory(ctx context.Context, id int) error {
	return s.categoryRepo.Delete(ctx, id)
}

// ListVehicleCategories retrieves vehicle categories with filtering and pagination
func (s *VehicleCategoryService) ListVehicleCategories(ctx context.Context, params *master.VehicleCategoryFilterParams) (*common.PaginatedResponse, error) {
	// Validate pagination parameters
	params.Validate()

	return s.categoryRepo.List(ctx, params)
}