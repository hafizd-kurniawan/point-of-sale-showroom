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

// GoodsReceiptHandler handles goods receipt HTTP requests
type GoodsReceiptHandler struct {
	goodsReceiptService *productService.GoodsReceiptService
}

// NewGoodsReceiptHandler creates a new goods receipt handler
func NewGoodsReceiptHandler(goodsReceiptService *productService.GoodsReceiptService) *GoodsReceiptHandler {
	return &GoodsReceiptHandler{
		goodsReceiptService: goodsReceiptService,
	}
}

// CreateGoodsReceipt handles goods receipt creation
func (h *GoodsReceiptHandler) CreateGoodsReceipt(c *gin.Context) {
	var req products.GoodsReceiptCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	receivedBy := middleware.GetCurrentUserID(c)
	if receivedBy == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Invalid user", "Receiver user ID not found",
		))
		return
	}

	receipt, err := h.goodsReceiptService.CreateGoodsReceipt(c.Request.Context(), &req, receivedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Goods receipt creation failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse("Goods receipt created successfully", receipt))
}

// GetGoodsReceipts handles goods receipts list with pagination and filtering
func (h *GoodsReceiptHandler) GetGoodsReceipts(c *gin.Context) {
	var params products.GoodsReceiptFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Query parameter validation failed", err.Error(),
		))
		return
	}

	result, err := h.goodsReceiptService.GetGoodsReceipts(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get goods receipts", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Goods receipts retrieved successfully", result))
}

// GetGoodsReceiptByID handles goods receipt retrieval by ID
func (h *GoodsReceiptHandler) GetGoodsReceiptByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid receipt ID", "Receipt ID must be a valid integer",
		))
		return
	}

	receipt, details, err := h.goodsReceiptService.GetGoodsReceiptByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Goods receipt not found", err.Error(),
		))
		return
	}

	response := map[string]interface{}{
		"receipt": receipt,
		"details": details,
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Goods receipt retrieved successfully", response))
}

// UpdateGoodsReceipt handles goods receipt updates
func (h *GoodsReceiptHandler) UpdateGoodsReceipt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid receipt ID", "Receipt ID must be a valid integer",
		))
		return
	}

	var req products.GoodsReceiptUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	receipt, err := h.goodsReceiptService.UpdateGoodsReceipt(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Goods receipt update failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Goods receipt updated successfully", receipt))
}

// AddGoodsReceiptDetails handles adding details to a goods receipt
func (h *GoodsReceiptHandler) AddGoodsReceiptDetails(c *gin.Context) {
	receiptIDStr := c.Param("id")
	receiptID, err := strconv.Atoi(receiptIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid receipt ID", "Receipt ID must be a valid integer",
		))
		return
	}

	var details []products.GoodsReceiptDetailCreateRequest
	if err := c.ShouldBindJSON(&details); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	if len(details) == 0 {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Empty request", "At least one receipt detail is required",
		))
		return
	}

	err = h.goodsReceiptService.AddGoodsReceiptDetails(c.Request.Context(), receiptID, details)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to add receipt details", err.Error(),
		))
		return
	}

	response := map[string]interface{}{
		"receipt_id":    receiptID,
		"details_count": len(details),
		"message":       "Receipt details added successfully",
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse("Receipt details added successfully", response))
}

// DeleteGoodsReceipt handles goods receipt deletion
func (h *GoodsReceiptHandler) DeleteGoodsReceipt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid receipt ID", "Receipt ID must be a valid integer",
		))
		return
	}

	err = h.goodsReceiptService.DeleteGoodsReceipt(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to delete goods receipt", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Goods receipt deleted successfully", nil))
}

// GetGoodsReceiptsByPO handles goods receipts by purchase order
func (h *GoodsReceiptHandler) GetGoodsReceiptsByPO(c *gin.Context) {
	poIDStr := c.Param("poId")
	poID, err := strconv.Atoi(poIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid PO ID", "PO ID must be a valid integer",
		))
		return
	}

	var params products.GoodsReceiptFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Query parameter validation failed", err.Error(),
		))
		return
	}

	result, err := h.goodsReceiptService.GetGoodsReceiptsByPO(c.Request.Context(), poID, &params)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to get PO receipts", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("PO receipts retrieved successfully", result))
}

// UpdateReceiptDetailQuantities handles receipt detail quantity updates
func (h *GoodsReceiptHandler) UpdateReceiptDetailQuantities(c *gin.Context) {
	detailIDStr := c.Param("detailId")
	detailID, err := strconv.Atoi(detailIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid detail ID", "Detail ID must be a valid integer",
		))
		return
	}

	var req struct {
		QuantityAccepted int `json:"quantity_accepted" binding:"required,min=0"`
		QuantityRejected int `json:"quantity_rejected" binding:"min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	err = h.goodsReceiptService.UpdateReceiptDetailQuantities(c.Request.Context(), detailID, req.QuantityAccepted, req.QuantityRejected)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to update quantities", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Quantities updated successfully", nil))
}