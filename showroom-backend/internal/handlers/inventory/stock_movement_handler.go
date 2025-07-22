package inventory

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	inventoryModels "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	inventoryServices "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services/inventory"
)

// StockMovementHandler handles HTTP requests for stock movements
type StockMovementHandler struct {
	stockMovementService *inventoryServices.StockMovementService
}

// NewStockMovementHandler creates a new stock movement handler
func NewStockMovementHandler(stockMovementService *inventoryServices.StockMovementService) *StockMovementHandler {
	return &StockMovementHandler{
		stockMovementService: stockMovementService,
	}
}

// GetStockMovement handles GET /stock-movements/:id
func (h *StockMovementHandler) GetStockMovement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid stock movement ID",
			Error:   err.Error(),
		})
		return
	}

	movement, err := h.stockMovementService.GetStockMovement(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.ErrorResponse{
			Status:  "error",
			Message: "Stock movement not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse{
		Status: "success",
		Data:   movement,
	})
}

// ListStockMovements handles GET /stock-movements
func (h *StockMovementHandler) ListStockMovements(c *gin.Context) {
	// Parse pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Parse filter parameters
	params := &inventoryModels.StockMovementFilterParams{}
	params.Page = page
	params.Limit = limit

	// Parse product ID filter
	if productIDStr := c.Query("product_id"); productIDStr != "" {
		if productID, err := strconv.Atoi(productIDStr); err == nil {
			params.ProductID = &productID
		}
	}

	// Parse movement type filter
	if movementType := c.Query("movement_type"); movementType != "" {
		mt := inventoryModels.MovementType(movementType)
		params.MovementType = &mt
	}

	// Parse search query
	if search := c.Query("search"); search != "" {
		params.Search = search
	}

	movements, total, err := h.stockMovementService.ListStockMovements(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve stock movements",
			Error:   err.Error(),
		})
		return
	}

	// Calculate pagination info
	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, common.SuccessResponse{
		Status: "success",
		Data:   movements,
		Meta: map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

// GetProductStockHistory handles GET /stock-movements/product/:id/history
func (h *StockMovementHandler) GetProductStockHistory(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid product ID",
			Error:   err.Error(),
		})
		return
	}

	// Parse pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	movements, total, err := h.stockMovementService.GetProductStockHistory(c.Request.Context(), productID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve product stock history",
			Error:   err.Error(),
		})
		return
	}

	// Calculate pagination info
	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, common.SuccessResponse{
		Status: "success",
		Data:   movements,
		Meta: map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

// CreateManualAdjustment handles POST /stock-movements/adjustment
func (h *StockMovementHandler) CreateManualAdjustment(c *gin.Context) {
	var req struct {
		ProductID          int    `json:"product_id" binding:"required"`
		AdjustmentQuantity int    `json:"adjustment_quantity" binding:"required"`
		Reason             string `json:"reason" binding:"required"`
		Notes              string `json:"notes,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Get user ID from context
	userID := c.GetInt("user_id")
	if userID == 0 {
		userID = 1 // Default for demo
	}

	movement, err := h.stockMovementService.CreateManualAdjustment(
		c.Request.Context(),
		req.ProductID,
		req.AdjustmentQuantity,
		req.Reason,
		userID,
	)
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
		Data:    movement,
	})
}

// CreateTransfer handles POST /stock-movements/transfer
func (h *StockMovementHandler) CreateTransfer(c *gin.Context) {
	var req struct {
		ProductID    int    `json:"product_id" binding:"required"`
		Quantity     int    `json:"quantity" binding:"required,min=1"`
		FromLocation string `json:"from_location" binding:"required"`
		ToLocation   string `json:"to_location" binding:"required"`
		Reason       string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Get user ID from context
	userID := c.GetInt("user_id")
	if userID == 0 {
		userID = 1 // Default for demo
	}

	movement, err := h.stockMovementService.CreateTransfer(
		c.Request.Context(),
		req.ProductID,
		req.Quantity,
		req.FromLocation,
		req.ToLocation,
		req.Reason,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Status:  "error",
			Message: "Failed to create stock transfer",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse{
		Status:  "success",
		Message: "Stock transfer created successfully",
		Data:    movement,
	})
}