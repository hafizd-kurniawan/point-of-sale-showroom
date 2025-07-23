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

// QualityInspectionHandler handles quality inspection endpoints
type QualityInspectionHandler struct {
	service *services.RepairService
}

// NewQualityInspectionHandler creates a new quality inspection handler
func NewQualityInspectionHandler(service *services.RepairService) *QualityInspectionHandler {
	return &QualityInspectionHandler{
		service: service,
	}
}

// CreateInspection creates a new quality inspection
// @Summary Create quality inspection
// @Description Create a new quality inspection
// @Tags Quality Inspections
// @Accept json
// @Produce json
// @Param request body repair.CreateQualityInspectionRequest true "Inspection data"
// @Success 201 {object} common.APIResponse{data=repair.QualityInspection}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/inspections [post]
func (h *QualityInspectionHandler) CreateInspection(c *gin.Context) {
	var req repair.CreateQualityInspectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	inspection, err := h.service.CreateInspection(c.Request.Context(), &req)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create inspection", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Inspection created successfully", inspection)
}

// GetInspections lists quality inspections
// @Summary List quality inspections
// @Description Get a paginated list of quality inspections
// @Tags Quality Inspections
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param work_order_id query int false "Filter by work order ID"
// @Param inspector_id query int false "Filter by inspector ID"
// @Param status query string false "Filter by inspection status"
// @Param type query string false "Filter by inspection type"
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/inspections [get]
func (h *QualityInspectionHandler) GetInspections(c *gin.Context) {
	params := &repair.QualityInspectionFilterParams{}
	
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
	if inspectorID := c.Query("inspector_id"); inspectorID != "" {
		if iid, err := strconv.Atoi(inspectorID); err == nil {
			params.InspectorID = iid
		}
	}
	params.Status = c.Query("status")
	params.InspectionType = c.Query("type")

	inspections, err := h.service.ListInspections(c.Request.Context(), params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get inspections", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Inspections retrieved successfully", inspections)
}

// GetInspection retrieves a quality inspection by ID
// @Summary Get quality inspection
// @Description Get a quality inspection by ID
// @Tags Quality Inspections
// @Accept json
// @Produce json
// @Param id path int true "Inspection ID"
// @Success 200 {object} common.APIResponse{data=repair.QualityInspection}
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/inspections/{id} [get]
func (h *QualityInspectionHandler) GetInspection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid inspection ID", err.Error())
		return
	}

	inspection, err := h.service.GetInspection(c.Request.Context(), id)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Inspection not found", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Inspection retrieved successfully", inspection)
}

// SignOffInspection signs off an inspection
// @Summary Sign off inspection
// @Description Sign off and complete a quality inspection
// @Tags Quality Inspections
// @Accept json
// @Produce json
// @Param id path int true "Inspection ID"
// @Param request body repair.InspectionSignOffRequest true "Sign off data"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/inspections/{id}/sign-off [post]
func (h *QualityInspectionHandler) SignOffInspection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid inspection ID", err.Error())
		return
	}

	var req repair.InspectionSignOffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	userID := utils.GetUserIDFromContext(c)
	err = h.service.SignOffInspection(c.Request.Context(), id, &req, userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to sign off inspection", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Inspection signed off successfully", nil)
}

// ScheduleRework schedules rework for failed inspection
// @Summary Schedule rework
// @Description Schedule rework for a failed quality inspection
// @Tags Quality Inspections
// @Accept json
// @Produce json
// @Param id path int true "Inspection ID"
// @Param request body repair.InspectionReworkRequest true "Rework data"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/inspections/{id}/schedule-rework [post]
func (h *QualityInspectionHandler) ScheduleRework(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid inspection ID", err.Error())
		return
	}

	var req repair.InspectionReworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	err = h.service.ScheduleRework(c.Request.Context(), id, &req)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to schedule rework", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Rework scheduled successfully", nil)
}

// ScheduleInspection schedules a new inspection
// @Summary Schedule inspection
// @Description Schedule a new quality inspection
// @Tags Quality Inspections
// @Accept json
// @Produce json
// @Param request body repair.InspectionScheduleRequest true "Schedule data"
// @Success 201 {object} common.APIResponse{data=repair.QualityInspection}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/inspections/schedule [post]
func (h *QualityInspectionHandler) ScheduleInspection(c *gin.Context) {
	var req repair.InspectionScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	inspection, err := h.service.ScheduleInspection(c.Request.Context(), &req)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to schedule inspection", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Inspection scheduled successfully", inspection)
}

// GetInspectionsByWorkOrder gets inspections for a specific work order
// @Summary Get inspections by work order
// @Description Get quality inspections for a specific work order
// @Tags Quality Inspections
// @Accept json
// @Produce json
// @Param work_order_id path int true "Work Order ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-orders/{work_order_id}/inspections [get]
func (h *QualityInspectionHandler) GetInspectionsByWorkOrder(c *gin.Context) {
	workOrderID, err := strconv.Atoi(c.Param("work_order_id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid work order ID", err.Error())
		return
	}

	params := &repair.QualityInspectionFilterParams{}
	
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

	inspections, err := h.service.GetInspectionsByWorkOrder(c.Request.Context(), workOrderID, params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get inspections by work order", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Inspections retrieved successfully", inspections)
}

// GetFailedInspections gets failed inspections
// @Summary Get failed inspections
// @Description Get quality inspections that failed
// @Tags Quality Inspections
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/inspections/failed [get]
func (h *QualityInspectionHandler) GetFailedInspections(c *gin.Context) {
	params := &repair.QualityInspectionFilterParams{}
	
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

	inspections, err := h.service.GetFailedInspections(c.Request.Context(), params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get failed inspections", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Failed inspections retrieved successfully", inspections)
}

// GetReworkRequired gets inspections requiring rework
// @Summary Get rework required
// @Description Get quality inspections that require rework
// @Tags Quality Inspections
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/inspections/rework-required [get]
func (h *QualityInspectionHandler) GetReworkRequired(c *gin.Context) {
	params := &repair.QualityInspectionFilterParams{}
	
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

	inspections, err := h.service.GetReworkRequired(c.Request.Context(), params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get rework required", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Rework required inspections retrieved successfully", inspections)
}

// GetQualityMetrics gets quality metrics for a work order
// @Summary Get quality metrics
// @Description Get quality metrics for a specific work order
// @Tags Quality Inspections
// @Accept json
// @Produce json
// @Param work_order_id path int true "Work Order ID"
// @Success 200 {object} common.APIResponse{data=repair.QualityMetrics}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-orders/{work_order_id}/quality-metrics [get]
func (h *QualityInspectionHandler) GetQualityMetrics(c *gin.Context) {
	workOrderID, err := strconv.Atoi(c.Param("work_order_id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid work order ID", err.Error())
		return
	}

	metrics, err := h.service.GetQualityMetrics(c.Request.Context(), workOrderID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get quality metrics", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Quality metrics retrieved successfully", metrics)
}

// GetInspectionDashboard gets inspection dashboard data
// @Summary Get inspection dashboard
// @Description Get quality inspection dashboard data
// @Tags Quality Inspections
// @Accept json
// @Produce json
// @Success 200 {object} common.APIResponse{data=repair.InspectionDashboard}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/inspections/dashboard [get]
func (h *QualityInspectionHandler) GetInspectionDashboard(c *gin.Context) {
	dashboard, err := h.service.GetInspectionDashboard(c.Request.Context())
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get inspection dashboard", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Inspection dashboard retrieved successfully", dashboard)
}