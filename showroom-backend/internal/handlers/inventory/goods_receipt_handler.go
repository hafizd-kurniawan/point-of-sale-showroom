package inventory

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	inventoryModels "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	inventoryServices "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services/inventory"
)

// GoodsReceiptHandler handles HTTP requests for goods receipts
type GoodsReceiptHandler struct {
	goodsReceiptService *inventoryServices.GoodsReceiptService
}

// NewGoodsReceiptHandler creates a new goods receipt handler
func NewGoodsReceiptHandler(goodsReceiptService *inventoryServices.GoodsReceiptService) *GoodsReceiptHandler {
	return &GoodsReceiptHandler{
		goodsReceiptService: goodsReceiptService,
	}
}

// CreateGoodsReceipt handles POST /goods-receipts
func (h *GoodsReceiptHandler) CreateGoodsReceipt(c *gin.Context) {
	var req inventoryModels.GoodsReceiptCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Get user ID from context (assuming it's set by auth middleware)
	userID := c.GetInt("user_id")
	if userID == 0 {
		userID = 1 // Default for demo
	}

	receipt, err := h.goodsReceiptService.CreateGoodsReceipt(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Status:  "error",
			Message: "Failed to create goods receipt",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse{
		Status:  "success",
		Message: "Goods receipt created successfully",
		Data:    receipt,
	})
}

// GetGoodsReceipt handles GET /goods-receipts/:id
func (h *GoodsReceiptHandler) GetGoodsReceipt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid goods receipt ID",
			Error:   err.Error(),
		})
		return
	}

	receipt, err := h.goodsReceiptService.GetGoodsReceipt(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.ErrorResponse{
			Status:  "error",
			Message: "Goods receipt not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse{
		Status: "success",
		Data:   receipt,
	})
}

// ListGoodsReceipts handles GET /goods-receipts
func (h *GoodsReceiptHandler) ListGoodsReceipts(c *gin.Context) {
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
	params := &inventoryModels.GoodsReceiptFilterParams{}
	params.Page = page
	params.Limit = limit

	// Parse status filter
	if status := c.Query("status"); status != "" {
		receiptStatus := inventoryModels.ReceiptStatus(status)
		params.ReceiptStatus = &receiptStatus
	}

	// Parse search query
	if search := c.Query("search"); search != "" {
		params.Search = search
	}

	receipts, total, err := h.goodsReceiptService.ListGoodsReceipts(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve goods receipts",
			Error:   err.Error(),
		})
		return
	}

	// Calculate pagination info
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
		"Goods receipts retrieved successfully", receipts, meta,
	))
}

// UpdateGoodsReceipt handles PUT /goods-receipts/:id
func (h *GoodsReceiptHandler) UpdateGoodsReceipt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid goods receipt ID",
			Error:   err.Error(),
		})
		return
	}

	var req inventoryModels.GoodsReceiptUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	receipt, err := h.goodsReceiptService.UpdateGoodsReceipt(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Status:  "error",
			Message: "Failed to update goods receipt",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse{
		Status:  "success",
		Message: "Goods receipt updated successfully",
		Data:    receipt,
	})
}

// DeleteGoodsReceipt handles DELETE /goods-receipts/:id
func (h *GoodsReceiptHandler) DeleteGoodsReceipt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid goods receipt ID",
			Error:   err.Error(),
		})
		return
	}

	err = h.goodsReceiptService.DeleteGoodsReceipt(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete goods receipt",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse{
		Status:  "success",
		Message: "Goods receipt deleted successfully",
	})
}