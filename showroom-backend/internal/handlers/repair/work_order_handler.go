package repair

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/repair"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/utils"
)

// WorkOrderHandler handles repair work order endpoints
type WorkOrderHandler struct {
	service *services.RepairService
}

// NewWorkOrderHandler creates a new work order handler
func NewWorkOrderHandler(service *services.RepairService) *WorkOrderHandler {
	return &WorkOrderHandler{
		service: service,
	}
}

// CreateWorkOrder creates a new repair work order
// @Summary Create repair work order
// @Description Create a new repair work order
// @Tags Repair Work Orders
// @Accept json
// @Produce json
// @Param request body repair.CreateRepairWorkOrderRequest true "Work order data"
// @Success 201 {object} common.APIResponse{data=repair.RepairWorkOrder}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-orders [post]
func (h *WorkOrderHandler) CreateWorkOrder(c *gin.Context) {
	var req repair.CreateRepairWorkOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	userID := utils.GetUserIDFromContext(c)
	workOrder, err := h.service.CreateWorkOrder(c.Request.Context(), &req, userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create work order", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Work order created successfully", workOrder)
}

// GetWorkOrders lists repair work orders
// @Summary List repair work orders
// @Description Get a paginated list of repair work orders
// @Tags Repair Work Orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param transaction_id query int false "Filter by transaction ID"
// @Param status query string false "Filter by status"
// @Param mechanic_id query int false "Filter by mechanic ID"
// @Param priority query string false "Filter by priority level"
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-orders [get]
func (h *WorkOrderHandler) GetWorkOrders(c *gin.Context) {
	params := &repair.RepairWorkOrderFilterParams{}
	
	// Parse pagination
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			params.Page = p
		}
	}
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			params.Limit = l
		}
	}
	
	// Parse filters
	if transactionID := c.Query("transaction_id"); transactionID != "" {
		if tid, err := strconv.Atoi(transactionID); err == nil {
			params.TransactionID = tid
		}
	}
	if mechanicID := c.Query("mechanic_id"); mechanicID != "" {
		if mid, err := strconv.Atoi(mechanicID); err == nil {
			params.MechanicID = mid
		}
	}
	params.Status = c.Query("status")
	params.Priority = c.Query("priority")

	workOrders, err := h.service.ListWorkOrders(c.Request.Context(), params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get work orders", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Work orders retrieved successfully", workOrders)
}

// GetWorkOrder retrieves a repair work order by ID
// @Summary Get repair work order
// @Description Get a repair work order by ID
// @Tags Repair Work Orders
// @Accept json
// @Produce json
// @Param id path int true "Work Order ID"
// @Success 200 {object} common.APIResponse{data=repair.RepairWorkOrder}
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-orders/{id} [get]
func (h *WorkOrderHandler) GetWorkOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid work order ID", err.Error())
		return
	}

	workOrder, err := h.service.GetWorkOrder(c.Request.Context(), id)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Work order not found", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Work order retrieved successfully", workOrder)
}

// GetWorkOrderByNumber retrieves a repair work order by number
// @Summary Get repair work order by number
// @Description Get a repair work order by work order number
// @Tags Repair Work Orders
// @Accept json
// @Produce json
// @Param number path string true "Work Order Number"
// @Success 200 {object} common.APIResponse{data=repair.RepairWorkOrder}
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-orders/number/{number} [get]
func (h *WorkOrderHandler) GetWorkOrderByNumber(c *gin.Context) {
	number := c.Param("number")
	if number == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Work order number is required", "")
		return
	}

	workOrder, err := h.service.GetWorkOrderByNumber(c.Request.Context(), number)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Work order not found", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Work order retrieved successfully", workOrder)
}

// UpdateWorkOrder updates a repair work order
// @Summary Update repair work order
// @Description Update a repair work order
// @Tags Repair Work Orders
// @Accept json
// @Produce json
// @Param id path int true "Work Order ID"
// @Param request body repair.UpdateRepairWorkOrderRequest true "Work order data"
// @Success 200 {object} common.APIResponse{data=repair.RepairWorkOrder}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-orders/{id} [put]
func (h *WorkOrderHandler) UpdateWorkOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid work order ID", err.Error())
		return
	}

	var req repair.UpdateRepairWorkOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	workOrder, err := h.service.UpdateWorkOrder(c.Request.Context(), id, &req)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update work order", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Work order updated successfully", workOrder)
}

// AssignMechanic assigns a mechanic to a work order
// @Summary Assign mechanic to work order
// @Description Assign a mechanic to a repair work order
// @Tags Repair Work Orders
// @Accept json
// @Produce json
// @Param id path int true "Work Order ID"
// @Param request body repair.WorkOrderAssignmentRequest true "Assignment data"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-orders/{id}/assign [post]
func (h *WorkOrderHandler) AssignMechanic(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid work order ID", err.Error())
		return
	}

	var req repair.WorkOrderAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	err = h.service.AssignMechanic(c.Request.Context(), id, &req)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to assign mechanic", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Mechanic assigned successfully", nil)
}

// ProcessApproval processes work order approval
// @Summary Process work order approval
// @Description Process approval for a repair work order
// @Tags Repair Work Orders
// @Accept json
// @Produce json
// @Param id path int true "Work Order ID"
// @Param request body repair.WorkOrderApprovalRequest true "Approval data"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-orders/{id}/approve [post]
func (h *WorkOrderHandler) ProcessApproval(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid work order ID", err.Error())
		return
	}

	var req repair.WorkOrderApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	userID := utils.GetUserIDFromContext(c)
	err = h.service.ProcessWorkOrderApproval(c.Request.Context(), id, &req, userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to process approval", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Approval processed successfully", nil)
}

// GetPendingApprovals gets work orders pending approval
// @Summary Get pending work order approvals
// @Description Get repair work orders pending approval
// @Tags Repair Work Orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-orders/pending-approval [get]
func (h *WorkOrderHandler) GetPendingApprovals(c *gin.Context) {
	params := &repair.RepairWorkOrderFilterParams{}
	
	// Parse pagination
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			params.Page = p
		}
	}
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			params.Limit = l
		}
	}

	workOrders, err := h.service.GetPendingWorkOrderApprovals(c.Request.Context(), params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get pending approvals", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Pending approvals retrieved successfully", workOrders)
}

// GetWorkOrdersByMechanic gets work orders assigned to a specific mechanic
// @Summary Get work orders by mechanic
// @Description Get repair work orders assigned to a specific mechanic
// @Tags Repair Work Orders
// @Accept json
// @Produce json
// @Param mechanic_id path int true "Mechanic ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-orders/mechanic/{mechanic_id} [get]
func (h *WorkOrderHandler) GetWorkOrdersByMechanic(c *gin.Context) {
	mechanicID, err := strconv.Atoi(c.Param("mechanic_id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid mechanic ID", err.Error())
		return
	}

	params := &repair.RepairWorkOrderFilterParams{}
	
	// Parse pagination
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			params.Page = p
		}
	}
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			params.Limit = l
		}
	}

	workOrders, err := h.service.GetWorkOrdersByMechanic(c.Request.Context(), mechanicID, params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get work orders by mechanic", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Work orders retrieved successfully", workOrders)
}

// GetWorkOrderSummary gets work order summary for a transaction
// @Summary Get work order summary
// @Description Get work order summary for a specific transaction
// @Tags Repair Work Orders
// @Accept json
// @Produce json
// @Param transaction_id path int true "Transaction ID"
// @Success 200 {object} common.APIResponse{data=repair.WorkOrderSummary}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/transactions/{transaction_id}/work-order-summary [get]
func (h *WorkOrderHandler) GetWorkOrderSummary(c *gin.Context) {
	transactionID, err := strconv.Atoi(c.Param("transaction_id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID", err.Error())
		return
	}

	summary, err := h.service.GetWorkOrderSummary(c.Request.Context(), transactionID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get work order summary", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Work order summary retrieved successfully", summary)
}