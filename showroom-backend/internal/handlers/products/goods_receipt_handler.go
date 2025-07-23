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

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Goods receipt created successfully", receipt,
	))
}

// GetGoodsReceipt handles getting a specific goods receipt
func (h *GoodsReceiptHandler) GetGoodsReceipt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid receipt ID", "Receipt ID must be a valid number",
		))
		return
	}

	receipt, err := h.goodsReceiptService.GetGoodsReceipt(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Goods receipt not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Goods receipt retrieved successfully", receipt,
	))
}

// UpdateGoodsReceipt handles updating a goods receipt
func (h *GoodsReceiptHandler) UpdateGoodsReceipt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid receipt ID", "Receipt ID must be a valid number",
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

	updatedReceipt, err := h.goodsReceiptService.UpdateGoodsReceipt(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to update goods receipt", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Goods receipt updated successfully", updatedReceipt,
	))
}

// DeleteGoodsReceipt handles deleting a goods receipt
func (h *GoodsReceiptHandler) DeleteGoodsReceipt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid receipt ID", "Receipt ID must be a valid number",
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

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Goods receipt deleted successfully", nil,
	))
}

// ListGoodsReceipts handles listing goods receipts with pagination
func (h *GoodsReceiptHandler) ListGoodsReceipts(c *gin.Context) {
	// Parse query parameters
	var params products.GoodsReceiptFilterParams
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

	receipts, err := h.goodsReceiptService.ListGoodsReceipts(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to list goods receipts", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Goods receipts retrieved successfully", receipts,
	))
}

// GetGoodsReceiptsByPO handles getting goods receipts for a specific PO
func (h *GoodsReceiptHandler) GetGoodsReceiptsByPO(c *gin.Context) {
	poIDStr := c.Param("poId")
	poID, err := strconv.Atoi(poIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid PO ID", "PO ID must be a valid number",
		))
		return
	}

	// Parse query parameters
	var params products.GoodsReceiptFilterParams
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

	receipts, err := h.goodsReceiptService.GetGoodsReceiptsByPO(c.Request.Context(), poID, &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get goods receipts by PO", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"PO goods receipts retrieved successfully", receipts,
	))
}

// ProcessGoodsReceipt handles processing a goods receipt and updating stock
func (h *GoodsReceiptHandler) ProcessGoodsReceipt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid receipt ID", "Receipt ID must be a valid number",
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

	err = h.goodsReceiptService.ProcessGoodsReceipt(c.Request.Context(), id, processedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to process goods receipt", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Goods receipt processed successfully", nil,
	))
}

// AddReceiptDetail handles adding a detail line to a goods receipt
func (h *GoodsReceiptHandler) AddReceiptDetail(c *gin.Context) {
	receiptIDStr := c.Param("id")
	receiptID, err := strconv.Atoi(receiptIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid receipt ID", "Receipt ID must be a valid number",
		))
		return
	}

	var req products.GoodsReceiptDetailCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	detail, err := h.goodsReceiptService.AddReceiptDetail(c.Request.Context(), receiptID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to add receipt detail", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Receipt detail added successfully", detail,
	))
}

// GetReceiptDetails handles getting all details for a goods receipt
func (h *GoodsReceiptHandler) GetReceiptDetails(c *gin.Context) {
	receiptIDStr := c.Param("id")
	receiptID, err := strconv.Atoi(receiptIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid receipt ID", "Receipt ID must be a valid number",
		))
		return
	}

	details, err := h.goodsReceiptService.GetReceiptDetails(c.Request.Context(), receiptID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get receipt details", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Receipt details retrieved successfully", details,
	))
}

// GetPendingReceiptItems handles getting items pending receipt for a PO
func (h *GoodsReceiptHandler) GetPendingReceiptItems(c *gin.Context) {
	poIDStr := c.Param("poId")
	poID, err := strconv.Atoi(poIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid PO ID", "PO ID must be a valid number",
		))
		return
	}

	items, err := h.goodsReceiptService.GetPendingReceiptItems(c.Request.Context(), poID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get pending receipt items", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Pending receipt items retrieved successfully", items,
	))
}

// BulkReceiveItems handles creating receipt details for multiple items at once
func (h *GoodsReceiptHandler) BulkReceiveItems(c *gin.Context) {
	receiptIDStr := c.Param("id")
	receiptID, err := strconv.Atoi(receiptIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid receipt ID", "Receipt ID must be a valid number",
		))
		return
	}

	var requests []products.GoodsReceiptDetailCreateRequest
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	err = h.goodsReceiptService.BulkReceiveItems(c.Request.Context(), receiptID, requests)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to bulk receive items", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Items received successfully", nil,
	))
}