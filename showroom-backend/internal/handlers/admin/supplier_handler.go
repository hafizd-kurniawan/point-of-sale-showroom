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

// SupplierHandler handles supplier HTTP requests
type SupplierHandler struct {
	supplierService *masterService.SupplierService
}

// NewSupplierHandler creates a new supplier handler
func NewSupplierHandler(supplierService *masterService.SupplierService) *SupplierHandler {
	return &SupplierHandler{
		supplierService: supplierService,
	}
}

// CreateSupplier handles supplier creation
func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var req master.SupplierCreateRequest
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

	supplier, err := h.supplierService.CreateSupplier(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Supplier creation failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Supplier created successfully", supplier,
	))
}

// GetSuppliers handles supplier list with filtering and pagination
func (h *SupplierHandler) GetSuppliers(c *gin.Context) {
	var params master.SupplierFilterParams

	// Bind query parameters
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Failed to parse query parameters", err.Error(),
		))
		return
	}

	// Handle supplier_type parameter
	if supplierTypeStr := c.Query("supplier_type"); supplierTypeStr != "" {
		supplierType := master.SupplierType(supplierTypeStr)
		if !supplierType.IsValid() {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(
				"Invalid supplier type", "Supplier type must be one of: parts, vehicle, both",
			))
			return
		}
		params.SupplierType = &supplierType
	}

	// Handle is_active parameter
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			params.IsActive = &isActive
		}
	}

	result, err := h.supplierService.ListSuppliers(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve suppliers", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Suppliers retrieved successfully", result,
	))
}

// GetSupplier handles getting a single supplier by ID
func (h *SupplierHandler) GetSupplier(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid supplier ID", "Supplier ID must be a valid integer",
		))
		return
	}

	supplier, err := h.supplierService.GetSupplier(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Supplier not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Supplier retrieved successfully", supplier,
	))
}

// UpdateSupplier handles supplier update
func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid supplier ID", "Supplier ID must be a valid integer",
		))
		return
	}

	var req master.SupplierUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	updatedSupplier, err := h.supplierService.UpdateSupplier(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Supplier update failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Supplier updated successfully", updatedSupplier,
	))
}

// DeleteSupplier handles supplier deletion (soft delete)
func (h *SupplierHandler) DeleteSupplier(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid supplier ID", "Supplier ID must be a valid integer",
		))
		return
	}

	err = h.supplierService.DeleteSupplier(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Supplier deletion failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Supplier deleted successfully", nil,
	))
}