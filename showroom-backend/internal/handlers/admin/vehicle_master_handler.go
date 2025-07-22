package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/middleware"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
	masterService "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services/master"
)

// VehicleMasterHandler handles vehicle master data HTTP requests
type VehicleMasterHandler struct {
	brandService    *masterService.VehicleBrandService
	categoryService *masterService.VehicleCategoryService
	modelService    *masterService.VehicleModelService
}

// NewVehicleMasterHandler creates a new vehicle master handler
func NewVehicleMasterHandler(
	brandService *masterService.VehicleBrandService,
	categoryService *masterService.VehicleCategoryService,
	modelService *masterService.VehicleModelService,
) *VehicleMasterHandler {
	return &VehicleMasterHandler{
		brandService:    brandService,
		categoryService: categoryService,
		modelService:    modelService,
	}
}

// ===== VEHICLE BRANDS =====

// CreateVehicleBrand handles vehicle brand creation
func (h *VehicleMasterHandler) CreateVehicleBrand(c *gin.Context) {
	var req master.VehicleBrandCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	createdBy := middleware.GetCurrentUserID(c)
	if createdBy == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Invalid user", "Creator user ID not found",
		))
		return
	}

	brand, err := h.brandService.CreateVehicleBrand(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Vehicle brand creation failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Vehicle brand created successfully", brand,
	))
}

// GetVehicleBrands handles vehicle brand list with filtering and pagination
func (h *VehicleMasterHandler) GetVehicleBrands(c *gin.Context) {
	var params master.VehicleBrandFilterParams

	// Bind query parameters
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Failed to parse query parameters", err.Error(),
		))
		return
	}

	// Handle is_active parameter
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			params.IsActive = &isActive
		}
	}

	result, err := h.brandService.ListVehicleBrands(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve vehicle brands", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Vehicle brands retrieved successfully", result,
	))
}

// GetVehicleBrand handles getting a single vehicle brand by ID
func (h *VehicleMasterHandler) GetVehicleBrand(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid vehicle brand ID", "Vehicle brand ID must be a valid integer",
		))
		return
	}

	brand, err := h.brandService.GetVehicleBrand(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Vehicle brand not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Vehicle brand retrieved successfully", brand,
	))
}

// UpdateVehicleBrand handles vehicle brand update
func (h *VehicleMasterHandler) UpdateVehicleBrand(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid vehicle brand ID", "Vehicle brand ID must be a valid integer",
		))
		return
	}

	var req master.VehicleBrandUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	updatedBrand, err := h.brandService.UpdateVehicleBrand(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Vehicle brand update failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Vehicle brand updated successfully", updatedBrand,
	))
}

// DeleteVehicleBrand handles vehicle brand deletion (soft delete)
func (h *VehicleMasterHandler) DeleteVehicleBrand(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid vehicle brand ID", "Vehicle brand ID must be a valid integer",
		))
		return
	}

	err = h.brandService.DeleteVehicleBrand(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Vehicle brand deletion failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Vehicle brand deleted successfully", nil,
	))
}

// ===== VEHICLE CATEGORIES =====

// CreateVehicleCategory handles vehicle category creation
func (h *VehicleMasterHandler) CreateVehicleCategory(c *gin.Context) {
	var req master.VehicleCategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	createdBy := middleware.GetCurrentUserID(c)
	if createdBy == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Invalid user", "Creator user ID not found",
		))
		return
	}

	category, err := h.categoryService.CreateVehicleCategory(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Vehicle category creation failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Vehicle category created successfully", category,
	))
}

// GetVehicleCategories handles vehicle category list with filtering and pagination
func (h *VehicleMasterHandler) GetVehicleCategories(c *gin.Context) {
	var params master.VehicleCategoryFilterParams

	// Bind query parameters
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Failed to parse query parameters", err.Error(),
		))
		return
	}

	// Handle is_active parameter
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			params.IsActive = &isActive
		}
	}

	result, err := h.categoryService.ListVehicleCategories(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve vehicle categories", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Vehicle categories retrieved successfully", result,
	))
}

// GetVehicleCategory handles getting a single vehicle category by ID
func (h *VehicleMasterHandler) GetVehicleCategory(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid vehicle category ID", "Vehicle category ID must be a valid integer",
		))
		return
	}

	category, err := h.categoryService.GetVehicleCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Vehicle category not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Vehicle category retrieved successfully", category,
	))
}

// UpdateVehicleCategory handles vehicle category update
func (h *VehicleMasterHandler) UpdateVehicleCategory(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid vehicle category ID", "Vehicle category ID must be a valid integer",
		))
		return
	}

	var req master.VehicleCategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	updatedCategory, err := h.categoryService.UpdateVehicleCategory(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Vehicle category update failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Vehicle category updated successfully", updatedCategory,
	))
}

// DeleteVehicleCategory handles vehicle category deletion (soft delete)
func (h *VehicleMasterHandler) DeleteVehicleCategory(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid vehicle category ID", "Vehicle category ID must be a valid integer",
		))
		return
	}

	err = h.categoryService.DeleteVehicleCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Vehicle category deletion failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Vehicle category deleted successfully", nil,
	))
}

// ===== VEHICLE MODELS =====

// CreateVehicleModel handles vehicle model creation
func (h *VehicleMasterHandler) CreateVehicleModel(c *gin.Context) {
	var req master.VehicleModelCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	createdBy := middleware.GetCurrentUserID(c)
	if createdBy == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Invalid user", "Creator user ID not found",
		))
		return
	}

	model, err := h.modelService.CreateVehicleModel(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Vehicle model creation failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Vehicle model created successfully", model,
	))
}

// GetVehicleModels handles vehicle model list with filtering and pagination
func (h *VehicleMasterHandler) GetVehicleModels(c *gin.Context) {
	var params master.VehicleModelFilterParams

	// Bind query parameters
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Failed to parse query parameters", err.Error(),
		))
		return
	}

	// Handle brand_id parameter
	if brandIDStr := c.Query("brand_id"); brandIDStr != "" {
		if brandID, err := strconv.Atoi(brandIDStr); err == nil {
			params.BrandID = &brandID
		}
	}

	// Handle category_id parameter
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.Atoi(categoryIDStr); err == nil {
			params.CategoryID = &categoryID
		}
	}

	// Handle model_year parameter
	if modelYearStr := c.Query("model_year"); modelYearStr != "" {
		if modelYear, err := strconv.Atoi(modelYearStr); err == nil {
			params.ModelYear = &modelYear
		}
	}

	// Handle min_price parameter
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			params.MinPrice = &minPrice
		}
	}

	// Handle max_price parameter
	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			params.MaxPrice = &maxPrice
		}
	}

	// Handle is_active parameter
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			params.IsActive = &isActive
		}
	}

	result, err := h.modelService.ListVehicleModels(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve vehicle models", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Vehicle models retrieved successfully", result,
	))
}

// GetVehicleModel handles getting a single vehicle model by ID
func (h *VehicleMasterHandler) GetVehicleModel(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid vehicle model ID", "Vehicle model ID must be a valid integer",
		))
		return
	}

	model, err := h.modelService.GetVehicleModel(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Vehicle model not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Vehicle model retrieved successfully", model,
	))
}

// UpdateVehicleModel handles vehicle model update
func (h *VehicleMasterHandler) UpdateVehicleModel(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid vehicle model ID", "Vehicle model ID must be a valid integer",
		))
		return
	}

	var req master.VehicleModelUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	updatedModel, err := h.modelService.UpdateVehicleModel(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Vehicle model update failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Vehicle model updated successfully", updatedModel,
	))
}

// DeleteVehicleModel handles vehicle model deletion (soft delete)
func (h *VehicleMasterHandler) DeleteVehicleModel(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid vehicle model ID", "Vehicle model ID must be a valid integer",
		))
		return
	}

	err = h.modelService.DeleteVehicleModel(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Vehicle model deletion failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Vehicle model deleted successfully", nil,
	))
}