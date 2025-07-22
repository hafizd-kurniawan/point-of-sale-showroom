package inventory

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/middleware"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services/inventory"
)

// PurchaseOrderHandler handles purchase order HTTP requests
type PurchaseOrderHandler struct {
	poService *inventory.PurchaseOrderService
}

// NewPurchaseOrderHandler creates a new purchase order handler
func NewPurchaseOrderHandler(poService *inventory.PurchaseOrderService) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{
		poService: poService,
	}
}

// CreatePurchaseOrder handles purchase order creation
// @Summary Create purchase order
// @Description Create a new purchase order
// @Tags inventory
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body inventory.PurchaseOrderPartCreateRequest true "Create purchase order request"
// @Success 201 {object} common.APIResponse{data=inventory.PurchaseOrderPart}
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Router /inventory/purchase-orders [post]
func (h *PurchaseOrderHandler) CreatePurchaseOrder(c *gin.Context) {
	var req inventory.PurchaseOrderPartCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	createdBy := middleware.GetCurrentUserID(c)
	if createdBy == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Unauthorized", "Invalid authentication", "USER_NOT_AUTHENTICATED",
		))
		return
	}

	po, err := h.poService.Create(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to create purchase order", err.Error(), "CREATION_FAILED",
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Purchase order created successfully", po,
	))
}

// GetPurchaseOrder handles getting a purchase order by ID
// @Summary Get purchase order
// @Description Get a purchase order by ID
// @Tags inventory
// @Security BearerAuth
// @Produce json
// @Param id path int true "Purchase Order ID"
// @Success 200 {object} common.APIResponse{data=inventory.PurchaseOrderPart}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /inventory/purchase-orders/{id} [get]
func (h *PurchaseOrderHandler) GetPurchaseOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Purchase order ID must be a valid integer", "INVALID_ID",
		))
		return
	}

	po, err := h.poService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Purchase order not found", err.Error(), "PO_NOT_FOUND",
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Purchase order retrieved successfully", po,
	))
}

// GetPurchaseOrderByNumber handles getting a purchase order by number
// @Summary Get purchase order by number
// @Description Get a purchase order by PO number
// @Tags inventory
// @Security BearerAuth
// @Produce json
// @Param number path string true "Purchase Order Number"
// @Success 200 {object} common.APIResponse{data=inventory.PurchaseOrderPart}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /inventory/purchase-orders/number/{number} [get]
func (h *PurchaseOrderHandler) GetPurchaseOrderByNumber(c *gin.Context) {
	number := c.Param("number")
	if number == "" {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid number", "Purchase order number is required", "INVALID_NUMBER",
		))
		return
	}

	po, err := h.poService.GetByNumber(c.Request.Context(), number)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Purchase order not found", err.Error(), "PO_NOT_FOUND",
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Purchase order retrieved successfully", po,
	))
}

// UpdatePurchaseOrder handles purchase order updates
// @Summary Update purchase order
// @Description Update a purchase order
// @Tags inventory
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Purchase Order ID"
// @Param request body inventory.PurchaseOrderPartUpdateRequest true "Update purchase order request"
// @Success 200 {object} common.APIResponse{data=inventory.PurchaseOrderPart}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /inventory/purchase-orders/{id} [put]
func (h *PurchaseOrderHandler) UpdatePurchaseOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Purchase order ID must be a valid integer", "INVALID_ID",
		))
		return
	}

	var req inventory.PurchaseOrderPartUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	po, err := h.poService.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to update purchase order", err.Error(), "UPDATE_FAILED",
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Purchase order updated successfully", po,
	))
}

// DeletePurchaseOrder handles purchase order deletion
// @Summary Delete purchase order
// @Description Delete a purchase order
// @Tags inventory
// @Security BearerAuth
// @Produce json
// @Param id path int true "Purchase Order ID"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /inventory/purchase-orders/{id} [delete]
func (h *PurchaseOrderHandler) DeletePurchaseOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Purchase order ID must be a valid integer", "INVALID_ID",
		))
		return
	}

	err = h.poService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to delete purchase order", err.Error(), "DELETION_FAILED",
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Purchase order deleted successfully", nil,
	))
}

// GetPurchaseOrders handles listing purchase orders with filtering and pagination
// @Summary List purchase orders
// @Description Get purchase orders with filtering and pagination
// @Tags inventory
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search query"
// @Param supplier_id query int false "Supplier ID filter"
// @Param status query string false "Status filter"
// @Param po_type query string false "PO type filter"
// @Success 200 {object} common.PaginatedResponse{data=[]inventory.PurchaseOrderPartListItem}
// @Failure 400 {object} common.ErrorResponse
// @Router /inventory/purchase-orders [get]
func (h *PurchaseOrderHandler) GetPurchaseOrders(c *gin.Context) {
	var params inventory.PurchaseOrderPartFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Invalid filter parameters", err.Error(),
		))
		return
	}

	orders, total, err := h.poService.List(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve purchase orders", err.Error(), "RETRIEVAL_FAILED",
		))
		return
	}

	response := common.NewPaginatedResponse(
		"Purchase orders retrieved successfully",
		orders,
		total,
		params.Page,
		params.Limit,
	)

	c.JSON(http.StatusOK, response)
}

// SearchPurchaseOrders handles purchase order search
// @Summary Search purchase orders
// @Description Search purchase orders by number, supplier name, or notes
// @Tags inventory
// @Security BearerAuth
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.PaginatedResponse{data=[]inventory.PurchaseOrderPartListItem}
// @Failure 400 {object} common.ErrorResponse
// @Router /inventory/purchase-orders/search [get]
func (h *PurchaseOrderHandler) SearchPurchaseOrders(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Missing search query", "Search query parameter 'q' is required", "MISSING_QUERY",
		))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, total, err := h.poService.Search(c.Request.Context(), query, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to search purchase orders", err.Error(), "SEARCH_FAILED",
		))
		return
	}

	response := common.NewPaginatedResponse(
		"Purchase order search completed successfully",
		orders,
		total,
		page,
		limit,
	)

	c.JSON(http.StatusOK, response)
}

// SendToSupplier handles sending a purchase order to supplier
// @Summary Send purchase order to supplier
// @Description Send a purchase order to supplier
// @Tags inventory
// @Security BearerAuth
// @Produce json
// @Param id path int true "Purchase Order ID"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Router /inventory/purchase-orders/{id}/send [post]
func (h *PurchaseOrderHandler) SendToSupplier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Purchase order ID must be a valid integer", "INVALID_ID",
		))
		return
	}

	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Unauthorized", "Invalid authentication", "USER_NOT_AUTHENTICATED",
		))
		return
	}

	err = h.poService.SendToSupplier(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to send purchase order", err.Error(), "SEND_FAILED",
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Purchase order sent to supplier successfully", nil,
	))
}

// ApprovePurchaseOrder handles approving a purchase order
// @Summary Approve purchase order
// @Description Approve a purchase order
// @Tags inventory
// @Security BearerAuth
// @Produce json
// @Param id path int true "Purchase Order ID"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Router /inventory/purchase-orders/{id}/approve [post]
func (h *PurchaseOrderHandler) ApprovePurchaseOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Purchase order ID must be a valid integer", "INVALID_ID",
		))
		return
	}

	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Unauthorized", "Invalid authentication", "USER_NOT_AUTHENTICATED",
		))
		return
	}

	err = h.poService.Approve(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to approve purchase order", err.Error(), "APPROVAL_FAILED",
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Purchase order approved successfully", nil,
	))
}

// CancelPurchaseOrder handles cancelling a purchase order
// @Summary Cancel purchase order
// @Description Cancel a purchase order
// @Tags inventory
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Purchase Order ID"
// @Param request body struct{Reason string `json:"reason" binding:"required"`} true "Cancellation reason"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Router /inventory/purchase-orders/{id}/cancel [post]
func (h *PurchaseOrderHandler) CancelPurchaseOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Purchase order ID must be a valid integer", "INVALID_ID",
		))
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Cancellation reason is required", err.Error(),
		))
		return
	}

	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Unauthorized", "Invalid authentication", "USER_NOT_AUTHENTICATED",
		))
		return
	}

	err = h.poService.Cancel(c.Request.Context(), id, userID, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to cancel purchase order", err.Error(), "CANCELLATION_FAILED",
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Purchase order cancelled successfully", nil,
	))
}

// GetWorkflowActions handles getting available workflow actions for a purchase order
// @Summary Get workflow actions
// @Description Get available workflow actions for a purchase order
// @Tags inventory
// @Security BearerAuth
// @Produce json
// @Param id path int true "Purchase Order ID"
// @Success 200 {object} common.APIResponse{data=[]string}
// @Failure 400 {object} common.ErrorResponse
// @Router /inventory/purchase-orders/{id}/actions [get]
func (h *PurchaseOrderHandler) GetWorkflowActions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Purchase order ID must be a valid integer", "INVALID_ID",
		))
		return
	}

	actions, err := h.poService.GetWorkflowActions(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to get workflow actions", err.Error(), "ACTIONS_FAILED",
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Workflow actions retrieved successfully", actions,
	))
}

// GetPendingApproval handles getting purchase orders pending approval
// @Summary Get pending approval purchase orders
// @Description Get purchase orders that are pending approval
// @Tags inventory
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.PaginatedResponse{data=[]inventory.PurchaseOrderPartListItem}
// @Failure 400 {object} common.ErrorResponse
// @Router /inventory/purchase-orders/pending-approval [get]
func (h *PurchaseOrderHandler) GetPendingApproval(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, total, err := h.poService.GetPendingApproval(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve pending approval purchase orders", err.Error(), "RETRIEVAL_FAILED",
		))
		return
	}

	response := common.NewPaginatedResponse(
		"Pending approval purchase orders retrieved successfully",
		orders,
		total,
		page,
		limit,
	)

	c.JSON(http.StatusOK, response)
}