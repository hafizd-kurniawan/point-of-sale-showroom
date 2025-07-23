package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/middleware"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
	productService "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services/products"
)

// PurchaseOrderHandler handles purchase order HTTP requests
type PurchaseOrderHandler struct {
	poService *productService.PurchaseOrderService
}

// NewPurchaseOrderHandler creates a new purchase order handler
func NewPurchaseOrderHandler(poService *productService.PurchaseOrderService) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{
		poService: poService,
	}
}

// CreatePurchaseOrder handles purchase order creation
func (h *PurchaseOrderHandler) CreatePurchaseOrder(c *gin.Context) {
	var req products.PurchaseOrderPartsCreateRequest
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

	po, err := h.poService.CreatePurchaseOrder(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Purchase order creation failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Purchase order created successfully", po,
	))
}

// GetPurchaseOrders handles purchase order listing
func (h *PurchaseOrderHandler) GetPurchaseOrders(c *gin.Context) {
	var params products.PurchaseOrderPartsFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid query parameters", err.Error(),
		))
		return
	}

	response, err := h.poService.ListPurchaseOrders(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve purchase orders", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetPurchaseOrder handles retrieving a specific purchase order
func (h *PurchaseOrderHandler) GetPurchaseOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Purchase order ID must be a valid integer",
		))
		return
	}

	po, err := h.poService.GetPurchaseOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Purchase order not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Purchase order retrieved successfully", po,
	))
}

// UpdatePurchaseOrder handles purchase order updates
func (h *PurchaseOrderHandler) UpdatePurchaseOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Purchase order ID must be a valid integer",
		))
		return
	}

	var req products.PurchaseOrderPartsUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	po, err := h.poService.UpdatePurchaseOrder(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Purchase order update failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Purchase order updated successfully", po,
	))
}

// AddLineItem handles adding line items to purchase order
func (h *PurchaseOrderHandler) AddLineItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Purchase order ID must be a valid integer",
		))
		return
	}

	var req products.PurchaseOrderDetailCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	detail, err := h.poService.AddLineItem(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Line item addition failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Line item added successfully", detail,
	))
}

// GetPurchaseOrderDetails handles retrieving purchase order line items
func (h *PurchaseOrderHandler) GetPurchaseOrderDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Purchase order ID must be a valid integer",
		))
		return
	}

	var params products.PurchaseOrderDetailFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid query parameters", err.Error(),
		))
		return
	}

	response, err := h.poService.GetPurchaseOrderDetails(c.Request.Context(), id, &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve purchase order details", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, response)
}

// ApprovePurchaseOrder handles purchase order approval
func (h *PurchaseOrderHandler) ApprovePurchaseOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Purchase order ID must be a valid integer",
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

	err = h.poService.ApprovePurchaseOrder(c.Request.Context(), id, approvedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Purchase order approval failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Purchase order approved successfully", nil,
	))
}

// SendPurchaseOrder handles sending purchase order to supplier
func (h *PurchaseOrderHandler) SendPurchaseOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Purchase order ID must be a valid integer",
		))
		return
	}

	err = h.poService.SendPurchaseOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to send purchase order", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Purchase order sent to supplier successfully", nil,
	))
}

// CancelPurchaseOrder handles purchase order cancellation
func (h *PurchaseOrderHandler) CancelPurchaseOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Purchase order ID must be a valid integer",
		))
		return
	}

	err = h.poService.CancelPurchaseOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Purchase order cancellation failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Purchase order cancelled successfully", nil,
	))
}

// GetPendingApproval handles retrieving purchase orders pending approval
func (h *PurchaseOrderHandler) GetPendingApproval(c *gin.Context) {
	var params products.PurchaseOrderPartsFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid query parameters", err.Error(),
		))
		return
	}

	response, err := h.poService.GetPendingApprovalPOs(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve pending approval purchase orders", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateGoodsReceipt handles goods receipt creation
func (h *PurchaseOrderHandler) CreateGoodsReceipt(c *gin.Context) {
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

	receipt, err := h.poService.CreateGoodsReceipt(c.Request.Context(), &req, receivedBy)
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

// ProcessReceiptItem handles processing individual items in goods receipt
func (h *PurchaseOrderHandler) ProcessReceiptItem(c *gin.Context) {
	receiptIDStr := c.Param("receiptId")
	receiptID, err := strconv.Atoi(receiptIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid receipt ID", "Receipt ID must be a valid integer",
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

	err = h.poService.ProcessGoodsReceiptItem(c.Request.Context(), receiptID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Receipt item processing failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Receipt item processed successfully", nil,
	))
}