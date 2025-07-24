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

// WorkDetailHandler handles repair work detail endpoints
type WorkDetailHandler struct {
	service *services.RepairService
}

// NewWorkDetailHandler creates a new work detail handler
func NewWorkDetailHandler(service *services.RepairService) *WorkDetailHandler {
	return &WorkDetailHandler{
		service: service,
	}
}

// CreateWorkDetail creates a new repair work detail
// @Summary Create repair work detail
// @Description Create a new repair work detail task
// @Tags Repair Work Details
// @Accept json
// @Produce json
// @Param request body repair.CreateRepairWorkDetailRequest true "Work detail data"
// @Success 201 {object} common.APIResponse{data=repair.RepairWorkDetail}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-details [post]
func (h *WorkDetailHandler) CreateWorkDetail(c *gin.Context) {
	var req repair.CreateRepairWorkDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	workDetail, err := h.service.CreateWorkDetail(c.Request.Context(), &req)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create work detail", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Work detail created successfully", workDetail)
}

// GetWorkDetails lists repair work details
// @Summary List repair work details
// @Description Get a paginated list of repair work details
// @Tags Repair Work Details
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param work_order_id query int false "Filter by work order ID"
// @Param damage_id query int false "Filter by damage ID"
// @Param status query string false "Filter by task status"
// @Param mechanic_id query int false "Filter by mechanic ID"
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-details [get]
func (h *WorkDetailHandler) GetWorkDetails(c *gin.Context) {
	params := &repair.RepairWorkDetailFilterParams{}
	
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
	if workOrderID := c.Query("work_order_id"); workOrderID != "" {
		if woid, err := strconv.Atoi(workOrderID); err == nil {
			params.WorkOrderID = woid
		}
	}
	if damageID := c.Query("damage_id"); damageID != "" {
		if did, err := strconv.Atoi(damageID); err == nil {
			params.DamageID = did
		}
	}
	if mechanicID := c.Query("mechanic_id"); mechanicID != "" {
		if mid, err := strconv.Atoi(mechanicID); err == nil {
			params.MechanicID = mid
		}
	}
	params.Status = c.Query("status")

	workDetails, err := h.service.ListWorkDetails(c.Request.Context(), params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get work details", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Work details retrieved successfully", workDetails)
}

// GetWorkDetail retrieves a repair work detail by ID
// @Summary Get repair work detail
// @Description Get a repair work detail by ID
// @Tags Repair Work Details
// @Accept json
// @Produce json
// @Param id path int true "Work Detail ID"
// @Success 200 {object} common.APIResponse{data=repair.RepairWorkDetail}
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-details/{id} [get]
func (h *WorkDetailHandler) GetWorkDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid work detail ID", err.Error())
		return
	}

	workDetail, err := h.service.GetWorkDetail(c.Request.Context(), id)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Work detail not found", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Work detail retrieved successfully", workDetail)
}

// UpdateWorkDetail updates a repair work detail
// @Summary Update repair work detail
// @Description Update a repair work detail
// @Tags Repair Work Details
// @Accept json
// @Produce json
// @Param id path int true "Work Detail ID"
// @Param request body repair.UpdateRepairWorkDetailRequest true "Work detail data"
// @Success 200 {object} common.APIResponse{data=repair.RepairWorkDetail}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-details/{id} [put]
func (h *WorkDetailHandler) UpdateWorkDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid work detail ID", err.Error())
		return
	}

	var req repair.UpdateRepairWorkDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	workDetail, err := h.service.UpdateWorkDetail(c.Request.Context(), id, &req)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update work detail", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Work detail updated successfully", workDetail)
}

// UpdateProgress updates work detail progress
// @Summary Update work detail progress
// @Description Update progress for a repair work detail task
// @Tags Repair Work Details
// @Accept json
// @Produce json
// @Param id path int true "Work Detail ID"
// @Param request body repair.WorkDetailProgressRequest true "Progress data"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-details/{id}/progress [post]
func (h *WorkDetailHandler) UpdateProgress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid work detail ID", err.Error())
		return
	}

	var req repair.WorkDetailProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	err = h.service.UpdateWorkDetailProgress(c.Request.Context(), id, &req)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update progress", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Progress updated successfully", nil)
}

// AssignMechanic assigns a mechanic to a work detail
// @Summary Assign mechanic to work detail
// @Description Assign a mechanic to a repair work detail task
// @Tags Repair Work Details
// @Accept json
// @Produce json
// @Param id path int true "Work Detail ID"
// @Param request body repair.WorkDetailAssignmentRequest true "Assignment data"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-details/{id}/assign [post]
func (h *WorkDetailHandler) AssignMechanic(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid work detail ID", err.Error())
		return
	}

	var req repair.WorkDetailAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	err = h.service.AssignWorkDetailMechanic(c.Request.Context(), id, &req)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to assign mechanic", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Mechanic assigned successfully", nil)
}

// PerformQualityCheck performs quality check on work detail
// @Summary Perform quality check
// @Description Perform quality check on completed work detail
// @Tags Repair Work Details
// @Accept json
// @Produce json
// @Param id path int true "Work Detail ID"
// @Param request body repair.WorkDetailQualityCheckRequest true "Quality check data"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-details/{id}/quality-check [post]
func (h *WorkDetailHandler) PerformQualityCheck(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid work detail ID", err.Error())
		return
	}

	var req repair.WorkDetailQualityCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	userID := utils.GetUserIDFromContext(c)
	err = h.service.PerformQualityCheck(c.Request.Context(), id, &req, userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to perform quality check", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Quality check performed successfully", nil)
}

// GetWorkDetailsByWorkOrder gets work details for a specific work order
// @Summary Get work details by work order
// @Description Get repair work details for a specific work order
// @Tags Repair Work Details
// @Accept json
// @Produce json
// @Param work_order_id path int true "Work Order ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-orders/{work_order_id}/details [get]
func (h *WorkDetailHandler) GetWorkDetailsByWorkOrder(c *gin.Context) {
	workOrderID, err := strconv.Atoi(c.Param("work_order_id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid work order ID", err.Error())
		return
	}

	params := &repair.RepairWorkDetailFilterParams{}
	
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

	workDetails, err := h.service.GetWorkDetailsByWorkOrder(c.Request.Context(), workOrderID, params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get work details by work order", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Work details retrieved successfully", workDetails)
}

// GetWorkDetailSummary gets work detail summary for a work order
// @Summary Get work detail summary
// @Description Get work detail summary for a specific work order
// @Tags Repair Work Details
// @Accept json
// @Produce json
// @Param work_order_id path int true "Work Order ID"
// @Success 200 {object} common.APIResponse{data=repair.WorkDetailSummary}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-orders/{work_order_id}/detail-summary [get]
func (h *WorkDetailHandler) GetWorkDetailSummary(c *gin.Context) {
	workOrderID, err := strconv.Atoi(c.Param("work_order_id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid work order ID", err.Error())
		return
	}

	summary, err := h.service.GetWorkDetailSummary(c.Request.Context(), workOrderID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get work detail summary", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Work detail summary retrieved successfully", summary)
}