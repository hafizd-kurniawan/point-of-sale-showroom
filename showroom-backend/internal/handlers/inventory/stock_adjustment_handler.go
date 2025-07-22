package inventory

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	inventoryModels "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	inventoryServices "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services/inventory"
)

// StockAdjustmentHandler handles HTTP requests for stock adjustments
type StockAdjustmentHandler struct {
	stockAdjustmentService *inventoryServices.StockAdjustmentService
}

// NewStockAdjustmentHandler creates a new stock adjustment handler
func NewStockAdjustmentHandler(stockAdjustmentService *inventoryServices.StockAdjustmentService) *StockAdjustmentHandler {
	return &StockAdjustmentHandler{
		stockAdjustmentService: stockAdjustmentService,
	}
}

// CreateStockAdjustment handles POST /stock-adjustments
func (h *StockAdjustmentHandler) CreateStockAdjustment(c *gin.Context) {
	var req inventoryModels.StockAdjustmentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Get user ID from context
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, common.ErrorResponse{
			Status:  "error",
			Message: "Unauthorized",
			Error:   "Invalid authentication",
		})
		return
	}

	adjustment, err := h.stockAdjustmentService.CreateStockAdjustment(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Status:  "error",
			Message: "Failed to create stock adjustment",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse{
		Status:  "success",
		Message: "Stock adjustment created successfully",
		Data:    adjustment,
	})
}

// GetStockAdjustment handles GET /stock-adjustments/:id
func (h *StockAdjustmentHandler) GetStockAdjustment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid adjustment ID",
		})
		return
	}

	adjustment, err := h.stockAdjustmentService.GetStockAdjustment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.ErrorResponse{
			Status:  "error",
			Message: "Stock adjustment not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse{
		Status:  "success",
		Message: "Stock adjustment retrieved successfully",
		Data:    adjustment,
	})
}

// GetStockAdjustments handles GET /stock-adjustments
func (h *StockAdjustmentHandler) GetStockAdjustments(c *gin.Context) {
	var params inventoryModels.StockAdjustmentFilterParams

	// Parse query parameters
	if productIDStr := c.Query("product_id"); productIDStr != "" {
		productID, err := strconv.Atoi(productIDStr)
		if err == nil {
			params.ProductID = &productID
		}
	}

	if adjustmentType := c.Query("adjustment_type"); adjustmentType != "" {
		adjType := inventoryModels.AdjustmentType(adjustmentType)
		params.AdjustmentType = &adjType
	}

	if approvedByStr := c.Query("approved_by"); approvedByStr != "" {
		approvedBy, err := strconv.Atoi(approvedByStr)
		if err == nil {
			params.ApprovedBy = &approvedBy
		}
	}

	adjustments, total, err := h.stockAdjustmentService.ListStockAdjustments(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve stock adjustments",
			Error:   err.Error(),
		})
		return
	}
totalPages := (total + limit - 1) / limit
hasMore := page < totalPages

meta := common.PaginationMeta{
Page:       page,
Limit:      limit,
Total:      total,
TotalPages: totalPages,
HasMore:    hasMore,
}

c.JSON(http.StatusOK, common.NewPaginationResponse(
"Stock adjustments retrieved successfully", adjustments, meta,
))
}

// UpdateStockAdjustment handles PUT /stock-adjustments/:id
func (h *StockAdjustmentHandler) UpdateStockAdjustment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid adjustment ID",
		})
		return
	}

	var req inventoryModels.StockAdjustmentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	adjustment, err := h.stockAdjustmentService.UpdateStockAdjustment(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Status:  "error",
			Message: "Failed to update stock adjustment",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse{
		Status:  "success",
		Message: "Stock adjustment updated successfully",
		Data:    adjustment,
	})
}

// DeleteStockAdjustment handles DELETE /stock-adjustments/:id
func (h *StockAdjustmentHandler) DeleteStockAdjustment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid adjustment ID",
		})
		return
	}

	err = h.stockAdjustmentService.DeleteStockAdjustment(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete stock adjustment",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse{
		Status:  "success",
		Message: "Stock adjustment deleted successfully",
	})
}

// ApproveStockAdjustment handles POST /stock-adjustments/:id/approve
func (h *StockAdjustmentHandler) ApproveStockAdjustment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid adjustment ID",
		})
		return
	}

	var req inventoryModels.StockAdjustmentApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	adjustment, err := h.stockAdjustmentService.ApproveStockAdjustment(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Status:  "error",
			Message: "Failed to approve stock adjustment",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse{
		Status:  "success",
		Message: "Stock adjustment approved successfully",
		Data:    adjustment,
	})
}

// GetStockAdjustmentsByProduct handles GET /products/:product_id/adjustments
func (h *StockAdjustmentHandler) GetStockAdjustmentsByProduct(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid product ID",
		})
		return
	}

	filter := inventoryModels.StockAdjustmentFilter{
		ProductID: &productID,
	}

	// Parse pagination
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	adjustments, total, err := h.stockAdjustmentService.GetStockAdjustments(&filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve stock adjustments",
			Error:   err.Error(),
		})
		return
	}

totalPages := (total + limit - 1) / limit
hasMore := page < totalPages

meta := common.PaginationMeta{
Page:       page,
Limit:      limit,
Total:      total,
TotalPages: totalPages,
HasMore:    hasMore,
}

c.JSON(http.StatusOK, common.NewPaginationResponse(
"Stock adjustments retrieved successfully", adjustments, meta,
))
}