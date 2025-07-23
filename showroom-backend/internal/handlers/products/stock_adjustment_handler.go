package products

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/middleware"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
	productService "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services/products"
)

// StockAdjustmentHandler handles stock adjustment HTTP requests
type StockAdjustmentHandler struct {
	stockAdjustmentService *productService.StockAdjustmentService
}

// NewStockAdjustmentHandler creates a new stock adjustment handler
func NewStockAdjustmentHandler(stockAdjustmentService *productService.StockAdjustmentService) *StockAdjustmentHandler {
	return &StockAdjustmentHandler{
		stockAdjustmentService: stockAdjustmentService,
	}
}

// CreateStockAdjustment handles stock adjustment creation
func (h *StockAdjustmentHandler) CreateStockAdjustment(c *gin.Context) {
	var req products.StockAdjustmentCreateRequest
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

	adjustment, err := h.stockAdjustmentService.CreateStockAdjustment(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Stock adjustment creation failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse("Stock adjustment created successfully", adjustment))
}

// GetStockAdjustments handles stock adjustments list with pagination and filtering
func (h *StockAdjustmentHandler) GetStockAdjustments(c *gin.Context) {
	var params products.StockAdjustmentFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Query parameter validation failed", err.Error(),
		))
		return
	}

	result, err := h.stockAdjustmentService.GetStockAdjustments(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get stock adjustments", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Stock adjustments retrieved successfully", result))
}

// GetStockAdjustmentByID handles stock adjustment retrieval by ID
func (h *StockAdjustmentHandler) GetStockAdjustmentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid adjustment ID", "Adjustment ID must be a valid integer",
		))
		return
	}

	adjustment, err := h.stockAdjustmentService.GetStockAdjustmentByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Stock adjustment not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Stock adjustment retrieved successfully", adjustment))
}

// UpdateStockAdjustment handles stock adjustment updates
func (h *StockAdjustmentHandler) UpdateStockAdjustment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid adjustment ID", "Adjustment ID must be a valid integer",
		))
		return
	}

	var req products.StockAdjustmentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	adjustment, err := h.stockAdjustmentService.UpdateStockAdjustment(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Stock adjustment update failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Stock adjustment updated successfully", adjustment))
}

// ApproveStockAdjustment handles stock adjustment approval
func (h *StockAdjustmentHandler) ApproveStockAdjustment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid adjustment ID", "Adjustment ID must be a valid integer",
		))
		return
	}

	approvedBy := middleware.GetCurrentUserID(c)
	if approvedBy == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Invalid user", "Approver user ID not found",
		))
		return
	}

	err = h.stockAdjustmentService.ApproveStockAdjustment(c.Request.Context(), id, approvedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Stock adjustment approval failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Stock adjustment approved successfully", nil))
}

// DeleteStockAdjustment handles stock adjustment deletion
func (h *StockAdjustmentHandler) DeleteStockAdjustment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid adjustment ID", "Adjustment ID must be a valid integer",
		))
		return
	}

	err = h.stockAdjustmentService.DeleteStockAdjustment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to delete stock adjustment", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Stock adjustment deleted successfully", nil))
}

// GetStockAdjustmentsByProduct handles stock adjustments by product
func (h *StockAdjustmentHandler) GetStockAdjustmentsByProduct(c *gin.Context) {
	productIDStr := c.Param("productId")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid product ID", "Product ID must be a valid integer",
		))
		return
	}

	var params products.StockAdjustmentFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Query parameter validation failed", err.Error(),
		))
		return
	}

	result, err := h.stockAdjustmentService.GetStockAdjustmentsByProduct(c.Request.Context(), productID, &params)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to get product adjustments", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Product adjustments retrieved successfully", result))
}

// GetPendingApprovalAdjustments handles pending approval adjustments
func (h *StockAdjustmentHandler) GetPendingApprovalAdjustments(c *gin.Context) {
	var params products.StockAdjustmentFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Query parameter validation failed", err.Error(),
		))
		return
	}

	result, err := h.stockAdjustmentService.GetPendingApprovalAdjustments(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get pending adjustments", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Pending adjustments retrieved successfully", result))
}

// GetVarianceReport handles variance report
func (h *StockAdjustmentHandler) GetVarianceReport(c *gin.Context) {
	var params products.StockAdjustmentFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Query parameter validation failed", err.Error(),
		))
		return
	}

	result, err := h.stockAdjustmentService.GetVarianceReport(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get variance report", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Variance report retrieved successfully", result))
}