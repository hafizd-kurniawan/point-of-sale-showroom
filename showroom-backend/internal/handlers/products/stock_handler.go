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

// StockHandler handles stock management HTTP requests
type StockHandler struct {
	stockService *productService.StockService
}

// NewStockHandler creates a new stock handler
func NewStockHandler(stockService *productService.StockService) *StockHandler {
	return &StockHandler{
		stockService: stockService,
	}
}

// GetStockMovements handles stock movements list with pagination and filtering
func (h *StockHandler) GetStockMovements(c *gin.Context) {
	var params products.StockMovementFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Query parameter validation failed", err.Error(),
		))
		return
	}

	result, err := h.stockService.GetStockMovements(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get stock movements", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Stock movements retrieved successfully", result))
}

// GetStockMovementByID handles stock movement retrieval by ID
func (h *StockHandler) GetStockMovementByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid movement ID", "Movement ID must be a valid integer",
		))
		return
	}

	movement, err := h.stockService.GetStockMovementByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Stock movement not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Stock movement retrieved successfully", movement))
}

// GetStockMovementsByProduct handles stock movements by product
func (h *StockHandler) GetStockMovementsByProduct(c *gin.Context) {
	productIDStr := c.Param("productId")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid product ID", "Product ID must be a valid integer",
		))
		return
	}

	var params products.StockMovementFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Query parameter validation failed", err.Error(),
		))
		return
	}

	result, err := h.stockService.GetStockMovementsByProduct(c.Request.Context(), productID, &params)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to get product stock movements", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Product stock movements retrieved successfully", result))
}

// GetProductMovementHistory handles product movement history
func (h *StockHandler) GetProductMovementHistory(c *gin.Context) {
	productIDStr := c.Param("productId")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid product ID", "Product ID must be a valid integer",
		))
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	movements, err := h.stockService.GetProductMovementHistory(c.Request.Context(), productID, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to get movement history", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Movement history retrieved successfully", movements))
}

// GetCurrentStock handles current stock retrieval for a product
func (h *StockHandler) GetCurrentStock(c *gin.Context) {
	productIDStr := c.Param("productId")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid product ID", "Product ID must be a valid integer",
		))
		return
	}

	stock, err := h.stockService.GetCurrentStock(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to get current stock", err.Error(),
		))
		return
	}

	response := map[string]interface{}{
		"product_id":     productID,
		"current_stock":  stock,
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Current stock retrieved successfully", response))
}

// CreateManualStockMovement handles manual stock movement creation
func (h *StockHandler) CreateManualStockMovement(c *gin.Context) {
	var req products.StockMovementCreateRequest
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

	movement, err := h.stockService.CreateManualStockMovement(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Stock movement creation failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse("Stock movement created successfully", movement))
}

// GetStockMovementsByReference handles stock movements by reference
func (h *StockHandler) GetStockMovementsByReference(c *gin.Context) {
	referenceType := products.ReferenceType(c.Query("reference_type"))
	referenceIDStr := c.Query("reference_id")

	if referenceType == "" || referenceIDStr == "" {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Missing parameters", "reference_type and reference_id are required",
		))
		return
	}

	referenceID, err := strconv.Atoi(referenceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid reference ID", "Reference ID must be a valid integer",
		))
		return
	}

	movements, err := h.stockService.GetStockMovementsByReference(c.Request.Context(), referenceType, referenceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to get stock movements", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Stock movements retrieved successfully", movements))
}

// GetLowStockProducts handles low stock products retrieval
func (h *StockHandler) GetLowStockProducts(c *gin.Context) {
	var params products.ProductSparePartFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Query parameter validation failed", err.Error(),
		))
		return
	}

	result, err := h.stockService.GetLowStockProducts(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get low stock products", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Low stock products retrieved successfully", result))
}

// GetStockSummary handles stock summary retrieval
func (h *StockHandler) GetStockSummary(c *gin.Context) {
	var productID *int
	productIDStr := c.Query("product_id")
	if productIDStr != "" {
		id, err := strconv.Atoi(productIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(
				"Invalid product ID", "Product ID must be a valid integer",
			))
			return
		}
		productID = &id
	}

	summary, err := h.stockService.GetStockSummary(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to get stock summary", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Stock summary retrieved successfully", summary))
}

// BulkCreateMovements handles bulk stock movement creation
func (h *StockHandler) BulkCreateMovements(c *gin.Context) {
	var movements []products.StockMovement
	if err := c.ShouldBindJSON(&movements); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	if len(movements) == 0 {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Empty request", "At least one movement is required",
		))
		return
	}

	err := h.stockService.BulkCreateMovements(c.Request.Context(), movements)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Bulk movement creation failed", err.Error(),
		))
		return
	}

	response := map[string]interface{}{
		"created_count": len(movements),
		"message":      "All movements created successfully",
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse("Bulk movements created successfully", response))
}