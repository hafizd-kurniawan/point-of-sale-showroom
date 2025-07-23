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

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Stock adjustment created successfully", adjustment,
	))
}

// GetStockAdjustment handles getting a specific stock adjustment
func (h *StockAdjustmentHandler) GetStockAdjustment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid adjustment ID", "Adjustment ID must be a valid number",
		))
		return
	}

	adjustment, err := h.stockAdjustmentService.GetStockAdjustment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Stock adjustment not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Stock adjustment retrieved successfully", adjustment,
	))
}

// UpdateStockAdjustment handles updating a stock adjustment
func (h *StockAdjustmentHandler) UpdateStockAdjustment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid adjustment ID", "Adjustment ID must be a valid number",
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

	updatedAdjustment, err := h.stockAdjustmentService.UpdateStockAdjustment(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to update stock adjustment", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Stock adjustment updated successfully", updatedAdjustment,
	))
}

// DeleteStockAdjustment handles deleting a stock adjustment
func (h *StockAdjustmentHandler) DeleteStockAdjustment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid adjustment ID", "Adjustment ID must be a valid number",
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

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Stock adjustment deleted successfully", nil,
	))
}

// ListStockAdjustments handles listing stock adjustments with pagination
func (h *StockAdjustmentHandler) ListStockAdjustments(c *gin.Context) {
	// Parse query parameters
	var params products.StockAdjustmentFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid query parameters", err.Error(),
		))
		return
	}

	// Set default pagination if not provided
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	adjustments, err := h.stockAdjustmentService.ListStockAdjustments(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to list stock adjustments", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Stock adjustments retrieved successfully", adjustments,
	))
}

// GetProductStockAdjustments handles getting stock adjustments for a specific product
func (h *StockAdjustmentHandler) GetProductStockAdjustments(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid product ID", "Product ID must be a valid number",
		))
		return
	}

	// Parse query parameters
	var params products.StockAdjustmentFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid query parameters", err.Error(),
		))
		return
	}

	// Set default pagination if not provided
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	adjustments, err := h.stockAdjustmentService.GetProductStockAdjustments(c.Request.Context(), productID, &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get product stock adjustments", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Product stock adjustments retrieved successfully", adjustments,
	))
}

// GetPendingAdjustments handles getting stock adjustments pending approval
func (h *StockAdjustmentHandler) GetPendingAdjustments(c *gin.Context) {
	// Parse query parameters
	var params products.StockAdjustmentFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid query parameters", err.Error(),
		))
		return
	}

	// Set default pagination if not provided
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	adjustments, err := h.stockAdjustmentService.GetPendingAdjustments(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get pending adjustments", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Pending adjustments retrieved successfully", adjustments,
	))
}

// ApproveStockAdjustment handles approving a stock adjustment
func (h *StockAdjustmentHandler) ApproveStockAdjustment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid adjustment ID", "Adjustment ID must be a valid number",
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

	// Validate adjustment workflow first
	err = h.stockAdjustmentService.ValidateAdjustmentWorkflow(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Adjustment validation failed", err.Error(),
		))
		return
	}

	err = h.stockAdjustmentService.ApproveStockAdjustment(c.Request.Context(), id, approvedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to approve stock adjustment", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Stock adjustment approved successfully", nil,
	))
}

// GetVarianceReport handles getting variance report for stock adjustments
func (h *StockAdjustmentHandler) GetVarianceReport(c *gin.Context) {
	// Parse query parameters
	var params products.StockAdjustmentFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid query parameters", err.Error(),
		))
		return
	}

	// Set default pagination if not provided
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	report, err := h.stockAdjustmentService.GetVarianceReport(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get variance report", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Variance report retrieved successfully", report,
	))
}

// CreatePhysicalCountAdjustments handles creating multiple adjustments from physical count
func (h *StockAdjustmentHandler) CreatePhysicalCountAdjustments(c *gin.Context) {
	var req PhysicalCountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	countedBy := middleware.GetCurrentUserID(c)
	if countedBy == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Invalid user", "Counter user ID not found",
		))
		return
	}

	adjustments, err := h.stockAdjustmentService.CreatePhysicalCountAdjustments(
		c.Request.Context(),
		req.Items,
		countedBy,
		req.Reason,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to create physical count adjustments", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Physical count adjustments created successfully", adjustments,
	))
}

// BulkApproveAdjustments handles approving multiple stock adjustments at once
func (h *StockAdjustmentHandler) BulkApproveAdjustments(c *gin.Context) {
	var req BulkApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
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

	err := h.stockAdjustmentService.BulkApproveAdjustments(c.Request.Context(), req.AdjustmentIDs, approvedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to bulk approve adjustments", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Adjustments approved successfully", nil,
	))
}

// PhysicalCountRequest represents a request to create adjustments from physical count
type PhysicalCountRequest struct {
	Items  []productService.PhysicalCountItem `json:"items" binding:"required,dive"`
	Reason string                             `json:"reason" binding:"required"`
}

// BulkApprovalRequest represents a request to approve multiple adjustments
type BulkApprovalRequest struct {
	AdjustmentIDs []int `json:"adjustment_ids" binding:"required,dive,min=1"`
}