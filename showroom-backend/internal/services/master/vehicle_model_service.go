package master

import (
	"context"
	"fmt"
	"strings"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// VehicleModelService handles vehicle model business logic
type VehicleModelService struct {
	modelRepo    interfaces.VehicleModelRepository
	brandRepo    interfaces.VehicleBrandRepository
	categoryRepo interfaces.VehicleCategoryRepository
}

// NewVehicleModelService creates a new vehicle model service
func NewVehicleModelService(
	modelRepo interfaces.VehicleModelRepository,
	brandRepo interfaces.VehicleBrandRepository,
	categoryRepo interfaces.VehicleCategoryRepository,
) *VehicleModelService {
	return &VehicleModelService{
		modelRepo:    modelRepo,
		brandRepo:    brandRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateVehicleModel creates a new vehicle model
func (s *VehicleModelService) CreateVehicleModel(ctx context.Context, req *master.VehicleModelCreateRequest, createdBy int) (*master.VehicleModel, error) {
	// Validate brand exists
	_, err := s.brandRepo.GetByID(ctx, req.BrandID)
	if err != nil {
		return nil, fmt.Errorf("invalid brand ID: %w", err)
	}

	// Validate category exists
	_, err = s.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	// Generate model code
	code, err := s.modelRepo.GenerateCode(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate model code: %w", err)
	}

	// Create model entity
	model := &master.VehicleModel{
		ModelCode:      code,
		ModelName:      strings.TrimSpace(req.ModelName),
		BrandID:        req.BrandID,
		CategoryID:     req.CategoryID,
		ModelYear:      req.ModelYear,
		EngineCapacity: req.EngineCapacity,
		FuelType:       strings.TrimSpace(req.FuelType),
		Transmission:   strings.TrimSpace(req.Transmission),
		SeatCapacity:   req.SeatCapacity,
		Color:          strings.TrimSpace(req.Color),
		Price:          req.Price,
		Description:    req.Description,
		ImageURL:       req.ImageURL,
		CreatedBy:      createdBy,
	}

	return s.modelRepo.Create(ctx, model)
}

// GetVehicleModel retrieves a vehicle model by ID
func (s *VehicleModelService) GetVehicleModel(ctx context.Context, id int) (*master.VehicleModel, error) {
	return s.modelRepo.GetByID(ctx, id)
}

// GetVehicleModelByCode retrieves a vehicle model by code
func (s *VehicleModelService) GetVehicleModelByCode(ctx context.Context, code string) (*master.VehicleModel, error) {
	return s.modelRepo.GetByCode(ctx, code)
}

// UpdateVehicleModel updates a vehicle model
func (s *VehicleModelService) UpdateVehicleModel(ctx context.Context, id int, req *master.VehicleModelUpdateRequest) (*master.VehicleModel, error) {
	// Get existing model
	existing, err := s.modelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate brand exists if changed
	if req.BrandID != nil && *req.BrandID != existing.BrandID {
		_, err := s.brandRepo.GetByID(ctx, *req.BrandID)
		if err != nil {
			return nil, fmt.Errorf("invalid brand ID: %w", err)
		}
	}

	// Validate category exists if changed
	if req.CategoryID != nil && *req.CategoryID != existing.CategoryID {
		_, err := s.categoryRepo.GetByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category ID: %w", err)
		}
	}

	// Update fields
	updatedModel := &master.VehicleModel{
		ModelCode:      existing.ModelCode,
		ModelName:      existing.ModelName,
		BrandID:        existing.BrandID,
		CategoryID:     existing.CategoryID,
		ModelYear:      existing.ModelYear,
		EngineCapacity: existing.EngineCapacity,
		FuelType:       existing.FuelType,
		Transmission:   existing.Transmission,
		SeatCapacity:   existing.SeatCapacity,
		Color:          existing.Color,
		Price:          existing.Price,
		Description:    existing.Description,
		ImageURL:       existing.ImageURL,
		IsActive:       existing.IsActive,
		CreatedAt:      existing.CreatedAt,
		CreatedBy:      existing.CreatedBy,
	}

	if req.ModelName != nil {
		updatedModel.ModelName = strings.TrimSpace(*req.ModelName)
	}
	if req.BrandID != nil {
		updatedModel.BrandID = *req.BrandID
	}
	if req.CategoryID != nil {
		updatedModel.CategoryID = *req.CategoryID
	}
	if req.ModelYear != nil {
		updatedModel.ModelYear = *req.ModelYear
	}
	if req.EngineCapacity != nil {
		updatedModel.EngineCapacity = req.EngineCapacity
	}
	if req.FuelType != nil {
		updatedModel.FuelType = strings.TrimSpace(*req.FuelType)
	}
	if req.Transmission != nil {
		updatedModel.Transmission = strings.TrimSpace(*req.Transmission)
	}
	if req.SeatCapacity != nil {
		updatedModel.SeatCapacity = *req.SeatCapacity
	}
	if req.Color != nil {
		updatedModel.Color = strings.TrimSpace(*req.Color)
	}
	if req.Price != nil {
		updatedModel.Price = *req.Price
	}
	if req.Description != nil {
		updatedModel.Description = req.Description
	}
	if req.ImageURL != nil {
		updatedModel.ImageURL = req.ImageURL
	}
	if req.IsActive != nil {
		updatedModel.IsActive = *req.IsActive
	}

	return s.modelRepo.Update(ctx, id, updatedModel)
}

// DeleteVehicleModel soft deletes a vehicle model
func (s *VehicleModelService) DeleteVehicleModel(ctx context.Context, id int) error {
	return s.modelRepo.Delete(ctx, id)
}

// ListVehicleModels retrieves vehicle models with filtering and pagination
func (s *VehicleModelService) ListVehicleModels(ctx context.Context, params *master.VehicleModelFilterParams) (*common.PaginatedResponse, error) {
	// Validate pagination parameters
	params.Validate()

	return s.modelRepo.List(ctx, params)
}