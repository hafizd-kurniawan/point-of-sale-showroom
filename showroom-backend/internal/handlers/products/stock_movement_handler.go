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

// StockMovementHandler handles stock movement HTTP requests
type StockMovementHandler struct {
	stockService *productService.StockService
}

// NewStockMovementHandler creates a new stock movement handler
func NewStockMovementHandler(stockService *productService.StockService) *StockMovementHandler {
	return &StockMovementHandler{
		stockService: stockService,
	}
}

// CreateStockMovement handles stock movement creation
func (h *StockMovementHandler) CreateStockMovement(c *gin.Context) {
	var req products.StockMovementCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	processedBy := middleware.GetCurrentUserID(c)
	if processedBy == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Invalid user", "Processor user ID not found",
		))
		return
	}

	movement, err := h.stockService.CreateStockMovement(c.Request.Context(), &req, processedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Stock movement creation failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Stock movement created successfully", movement,
	))
}

// GetStockMovement handles getting a specific stock movement
func (h *StockMovementHandler) GetStockMovement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid movement ID", "Movement ID must be a valid number",
		))
		return
	}

	movement, err := h.stockService.GetStockMovement(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Stock movement not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Stock movement retrieved successfully", movement,
	))
}

// ListStockMovements handles listing stock movements with pagination
func (h *StockMovementHandler) ListStockMovements(c *gin.Context) {
	// Parse query parameters
	var params products.StockMovementFilterParams
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

	movements, err := h.stockService.ListStockMovements(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to list stock movements", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Stock movements retrieved successfully", movements,
	))
}

// GetProductStockMovements handles getting stock movements for a specific product
func (h *StockMovementHandler) GetProductStockMovements(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid product ID", "Product ID must be a valid number",
		))
		return
	}

	// Parse query parameters
	var params products.StockMovementFilterParams
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

	movements, err := h.stockService.GetProductStockMovements(c.Request.Context(), productID, &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get product stock movements", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Product stock movements retrieved successfully", movements,
	))
}

// GetProductStockHistory handles getting recent stock movement history for a product
func (h *StockMovementHandler) GetProductStockHistory(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid product ID", "Product ID must be a valid number",
		))
		return
	}

	// Parse limit parameter
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	movements, err := h.stockService.GetMovementHistory(c.Request.Context(), productID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get movement history", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Movement history retrieved successfully", movements,
	))
}

// GetCurrentStock handles getting current stock quantity for a product
func (h *StockMovementHandler) GetCurrentStock(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid product ID", "Product ID must be a valid number",
		))
		return
	}

	stock, err := h.stockService.GetCurrentStock(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get current stock", err.Error(),
		))
		return
	}

	response := map[string]interface{}{
		"product_id":     productID,
		"current_stock":  stock,
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Current stock retrieved successfully", response,
	))
}

// TransferStock handles stock transfer between locations
func (h *StockMovementHandler) TransferStock(c *gin.Context) {
	var req TransferStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	processedBy := middleware.GetCurrentUserID(c)
	if processedBy == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Invalid user", "Processor user ID not found",
		))
		return
	}

	err := h.stockService.TransferStock(
		c.Request.Context(),
		req.ProductID,
		req.Quantity,
		req.FromLocation,
		req.ToLocation,
		processedBy,
		req.Notes,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Stock transfer failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Stock transfer completed successfully", nil,
	))
}

// TransferStockRequest represents a request to transfer stock between locations
type TransferStockRequest struct {
	ProductID    int     `json:"product_id" binding:"required,min=1"`
	Quantity     int     `json:"quantity" binding:"required,min=1"`
	FromLocation string  `json:"from_location" binding:"required"`
	ToLocation   string  `json:"to_location" binding:"required"`
	Notes        *string `json:"notes,omitempty"`
}