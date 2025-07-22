package master

import (
	"context"
	"fmt"
	"strings"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// ProductCategoryService handles product category business logic
type ProductCategoryService struct {
	categoryRepo interfaces.ProductCategoryRepository
}

// NewProductCategoryService creates a new product category service
func NewProductCategoryService(categoryRepo interfaces.ProductCategoryRepository) *ProductCategoryService {
	return &ProductCategoryService{
		categoryRepo: categoryRepo,
	}
}

// CreateProductCategory creates a new product category
func (s *ProductCategoryService) CreateProductCategory(ctx context.Context, req *master.ProductCategoryCreateRequest, createdBy int) (*master.ProductCategory, error) {
	// Check for duplicate category name within the same parent
	nameExists, err := s.categoryRepo.IsNameExists(ctx, req.CategoryName, req.ParentID, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to check category name existence: %w", err)
	}
	if nameExists {
		return nil, fmt.Errorf("category name already exists in the same parent category")
	}

	// Validate parent category exists if provided
	var level int = 1
	var path string
	if req.ParentID != nil {
		parent, err := s.categoryRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent category ID: %w", err)
		}
		level = parent.Level + 1
		path = parent.Path
	}

	// Generate category code
	code, err := s.categoryRepo.GenerateCode(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate category code: %w", err)
	}

	// Build path
	if path == "" {
		path = code
	} else {
		path = path + "/" + code
	}

	// Create category entity
	category := &master.ProductCategory{
		CategoryCode: code,
		CategoryName: strings.TrimSpace(req.CategoryName),
		Description:  req.Description,
		ParentID:     req.ParentID,
		Level:        level,
		Path:         path,
		CreatedBy:    createdBy,
	}

	return s.categoryRepo.Create(ctx, category)
}

// GetProductCategory retrieves a product category by ID
func (s *ProductCategoryService) GetProductCategory(ctx context.Context, id int) (*master.ProductCategory, error) {
	return s.categoryRepo.GetByID(ctx, id)
}

// GetProductCategoryByCode retrieves a product category by code
func (s *ProductCategoryService) GetProductCategoryByCode(ctx context.Context, code string) (*master.ProductCategory, error) {
	return s.categoryRepo.GetByCode(ctx, code)
}

// UpdateProductCategory updates a product category
func (s *ProductCategoryService) UpdateProductCategory(ctx context.Context, id int, req *master.ProductCategoryUpdateRequest) (*master.ProductCategory, error) {
	// Get existing category
	existing, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check for duplicate category name if changed
	if req.CategoryName != nil && *req.CategoryName != existing.CategoryName {
		nameExists, err := s.categoryRepo.IsNameExists(ctx, *req.CategoryName, existing.ParentID, id)
		if err != nil {
			return nil, fmt.Errorf("failed to check category name existence: %w", err)
		}
		if nameExists {
			return nil, fmt.Errorf("category name already exists in the same parent category")
		}
	}

	// Handle parent change
	var level int = existing.Level
	var path string = existing.Path
	if req.ParentID != nil {
		// Prevent circular reference
		if *req.ParentID == id {
			return nil, fmt.Errorf("category cannot be its own parent")
		}

		if *req.ParentID != 0 {
			// Validate parent exists
			parent, err := s.categoryRepo.GetByID(ctx, *req.ParentID)
			if err != nil {
				return nil, fmt.Errorf("invalid parent category ID: %w", err)
			}

			// Check if the new parent is not a descendant of current category
			if strings.Contains(parent.Path, existing.CategoryCode) {
				return nil, fmt.Errorf("cannot move category to its own descendant")
			}

			level = parent.Level + 1
			path = parent.Path + "/" + existing.CategoryCode
		} else {
			// Moving to root level
			level = 1
			path = existing.CategoryCode
		}
	}

	// Update fields
	updatedCategory := &master.ProductCategory{
		CategoryCode: existing.CategoryCode,
		CategoryName: existing.CategoryName,
		Description:  existing.Description,
		ParentID:     existing.ParentID,
		Level:        level,
		Path:         path,
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
	if req.ParentID != nil {
		if *req.ParentID == 0 {
			updatedCategory.ParentID = nil
		} else {
			updatedCategory.ParentID = req.ParentID
		}
	}
	if req.IsActive != nil {
		updatedCategory.IsActive = *req.IsActive
	}

	return s.categoryRepo.Update(ctx, id, updatedCategory)
}

// DeleteProductCategory soft deletes a product category
func (s *ProductCategoryService) DeleteProductCategory(ctx context.Context, id int) error {
	// Check if category has children
	children, err := s.categoryRepo.GetChildren(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check for child categories: %w", err)
	}
	if len(children) > 0 {
		return fmt.Errorf("cannot delete category with child categories")
	}

	return s.categoryRepo.Delete(ctx, id)
}

// ListProductCategories retrieves product categories with filtering and pagination
func (s *ProductCategoryService) ListProductCategories(ctx context.Context, params *master.ProductCategoryFilterParams) (*common.PaginatedResponse, error) {
	// Validate pagination parameters
	params.Validate()

	return s.categoryRepo.List(ctx, params)
}

// GetProductCategoryTree retrieves product categories in hierarchical structure
func (s *ProductCategoryService) GetProductCategoryTree(ctx context.Context) ([]master.ProductCategoryTree, error) {
	return s.categoryRepo.GetTree(ctx)
}

// GetProductCategoryChildren retrieves child categories of a parent
func (s *ProductCategoryService) GetProductCategoryChildren(ctx context.Context, parentID int) ([]master.ProductCategory, error) {
	return s.categoryRepo.GetChildren(ctx, parentID)
}