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

// ProductCategoryHandler handles product category HTTP requests
type ProductCategoryHandler struct {
	categoryService *masterService.ProductCategoryService
}

// NewProductCategoryHandler creates a new product category handler
func NewProductCategoryHandler(categoryService *masterService.ProductCategoryService) *ProductCategoryHandler {
	return &ProductCategoryHandler{
		categoryService: categoryService,
	}
}

// CreateProductCategory handles product category creation
func (h *ProductCategoryHandler) CreateProductCategory(c *gin.Context) {
	var req master.ProductCategoryCreateRequest
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

	category, err := h.categoryService.CreateProductCategory(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Product category creation failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Product category created successfully", category,
	))
}

// GetProductCategories handles product category list with filtering and pagination
func (h *ProductCategoryHandler) GetProductCategories(c *gin.Context) {
	var params master.ProductCategoryFilterParams

	// Bind query parameters
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Failed to parse query parameters", err.Error(),
		))
		return
	}

	// Handle parent_id parameter
	if parentIDStr := c.Query("parent_id"); parentIDStr != "" {
		if parentID, err := strconv.Atoi(parentIDStr); err == nil {
			params.ParentID = &parentID
		}
	}

	// Handle level parameter
	if levelStr := c.Query("level"); levelStr != "" {
		if level, err := strconv.Atoi(levelStr); err == nil {
			params.Level = &level
		}
	}

	// Handle is_active parameter
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			params.IsActive = &isActive
		}
	}

	result, err := h.categoryService.ListProductCategories(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve product categories", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Product categories retrieved successfully", result,
	))
}

// GetProductCategory handles getting a single product category by ID
func (h *ProductCategoryHandler) GetProductCategory(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid product category ID", "Product category ID must be a valid integer",
		))
		return
	}

	category, err := h.categoryService.GetProductCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Product category not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Product category retrieved successfully", category,
	))
}

// UpdateProductCategory handles product category update
func (h *ProductCategoryHandler) UpdateProductCategory(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid product category ID", "Product category ID must be a valid integer",
		))
		return
	}

	var req master.ProductCategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	updatedCategory, err := h.categoryService.UpdateProductCategory(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Product category update failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Product category updated successfully", updatedCategory,
	))
}

// DeleteProductCategory handles product category deletion (soft delete)
func (h *ProductCategoryHandler) DeleteProductCategory(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid product category ID", "Product category ID must be a valid integer",
		))
		return
	}

	err = h.categoryService.DeleteProductCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Product category deletion failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Product category deleted successfully", nil,
	))
}

// GetProductCategoryTree handles getting product categories in hierarchical structure
func (h *ProductCategoryHandler) GetProductCategoryTree(c *gin.Context) {
	tree, err := h.categoryService.GetProductCategoryTree(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve product category tree", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Product category tree retrieved successfully", tree,
	))
}

// GetProductCategoryChildren handles getting child categories of a parent
func (h *ProductCategoryHandler) GetProductCategoryChildren(c *gin.Context) {
	parentID, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid parent category ID", "Parent category ID must be a valid integer",
		))
		return
	}

	children, err := h.categoryService.GetProductCategoryChildren(c.Request.Context(), parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve child categories", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Child categories retrieved successfully", children,
	))
}