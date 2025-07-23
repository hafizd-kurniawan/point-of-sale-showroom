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

// PartsUsageHandler handles repair parts usage endpoints
type PartsUsageHandler struct {
	service *services.RepairService
}

// NewPartsUsageHandler creates a new parts usage handler
func NewPartsUsageHandler(service *services.RepairService) *PartsUsageHandler {
	return &PartsUsageHandler{
		service: service,
	}
}

// CreatePartsUsage creates a new repair parts usage record
// @Summary Create repair parts usage
// @Description Create a new repair parts usage record
// @Tags Repair Parts Usage
// @Accept json
// @Produce json
// @Param request body repair.CreateRepairPartsUsageRequest true "Parts usage data"
// @Success 201 {object} common.APIResponse{data=repair.RepairPartsUsage}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/parts-usage [post]
func (h *PartsUsageHandler) CreatePartsUsage(c *gin.Context) {
	var req repair.CreateRepairPartsUsageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	userID := utils.GetUserIDFromContext(c)
	partsUsage, err := h.service.CreatePartsUsage(c.Request.Context(), &req, userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create parts usage", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Parts usage created successfully", partsUsage)
}

// GetPartsUsages lists repair parts usage records
// @Summary List repair parts usage
// @Description Get a paginated list of repair parts usage records
// @Tags Repair Parts Usage
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param work_detail_id query int false "Filter by work detail ID"
// @Param product_id query int false "Filter by product ID"
// @Param usage_type query string false "Filter by usage type"
// @Param status query string false "Filter by usage status"
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/parts-usage [get]
func (h *PartsUsageHandler) GetPartsUsages(c *gin.Context) {
	params := &repair.RepairPartsUsageFilterParams{}
	
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
	if workDetailID := c.Query("work_detail_id"); workDetailID != "" {
		if wdid, err := strconv.Atoi(workDetailID); err == nil {
			params.WorkDetailID = wdid
		}
	}
	if productID := c.Query("product_id"); productID != "" {
		if pid, err := strconv.Atoi(productID); err == nil {
			params.ProductID = pid
		}
	}
	params.UsageType = c.Query("usage_type")
	params.Status = c.Query("status")

	partsUsages, err := h.service.ListPartsUsage(c.Request.Context(), params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get parts usage", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Parts usage retrieved successfully", partsUsages)
}

// GetPartsUsage retrieves a repair parts usage by ID
// @Summary Get repair parts usage
// @Description Get a repair parts usage by ID
// @Tags Repair Parts Usage
// @Accept json
// @Produce json
// @Param id path int true "Parts Usage ID"
// @Success 200 {object} common.APIResponse{data=repair.RepairPartsUsage}
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/parts-usage/{id} [get]
func (h *PartsUsageHandler) GetPartsUsage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid parts usage ID", err.Error())
		return
	}

	partsUsage, err := h.service.GetPartsUsage(c.Request.Context(), id)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Parts usage not found", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Parts usage retrieved successfully", partsUsage)
}

// ProcessApproval processes parts usage approval
// @Summary Process parts usage approval
// @Description Process approval for a repair parts usage
// @Tags Repair Parts Usage
// @Accept json
// @Produce json
// @Param id path int true "Parts Usage ID"
// @Param request body repair.PartsUsageApprovalRequest true "Approval data"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/parts-usage/{id}/approve [post]
func (h *PartsUsageHandler) ProcessApproval(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid parts usage ID", err.Error())
		return
	}

	var req repair.PartsUsageApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	userID := utils.GetUserIDFromContext(c)
	err = h.service.ProcessPartsUsageApproval(c.Request.Context(), id, &req, userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to process approval", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Approval processed successfully", nil)
}

// IssuePartsForRepair issues parts for repair
// @Summary Issue parts for repair
// @Description Issue parts for a repair work detail
// @Tags Repair Parts Usage
// @Accept json
// @Produce json
// @Param work_detail_id path int true "Work Detail ID"
// @Param request body repair.PartsUsageIssueRequest true "Parts issue data"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-details/{work_detail_id}/issue-parts [post]
func (h *PartsUsageHandler) IssuePartsForRepair(c *gin.Context) {
	workDetailID, err := strconv.Atoi(c.Param("work_detail_id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid work detail ID", err.Error())
		return
	}

	var req repair.PartsUsageIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	userID := utils.GetUserIDFromContext(c)
	err = h.service.IssuePartsForRepair(c.Request.Context(), workDetailID, &req, userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to issue parts", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Parts issued successfully", nil)
}

// GetPartsUsageByWorkDetail gets parts usage for a specific work detail
// @Summary Get parts usage by work detail
// @Description Get repair parts usage for a specific work detail
// @Tags Repair Parts Usage
// @Accept json
// @Produce json
// @Param work_detail_id path int true "Work Detail ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-details/{work_detail_id}/parts-usage [get]
func (h *PartsUsageHandler) GetPartsUsageByWorkDetail(c *gin.Context) {
	workDetailID, err := strconv.Atoi(c.Param("work_detail_id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid work detail ID", err.Error())
		return
	}

	params := &repair.RepairPartsUsageFilterParams{}
	
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

	partsUsages, err := h.service.GetPartsUsageByWorkDetail(c.Request.Context(), workDetailID, params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get parts usage by work detail", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Parts usage retrieved successfully", partsUsages)
}

// GetPartsUsageSummary gets parts usage summary for a work order
// @Summary Get parts usage summary
// @Description Get parts usage summary for a specific work order
// @Tags Repair Parts Usage
// @Accept json
// @Produce json
// @Param work_order_id path int true "Work Order ID"
// @Success 200 {object} common.APIResponse{data=repair.PartsUsageSummary}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-orders/{work_order_id}/parts-summary [get]
func (h *PartsUsageHandler) GetPartsUsageSummary(c *gin.Context) {
	workOrderID, err := strconv.Atoi(c.Param("work_order_id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid work order ID", err.Error())
		return
	}

	summary, err := h.service.GetPartsUsageSummary(c.Request.Context(), workOrderID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get parts usage summary", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Parts usage summary retrieved successfully", summary)
}

// GetInventoryImpact gets inventory impact for a work order
// @Summary Get inventory impact
// @Description Get inventory impact for a specific work order
// @Tags Repair Parts Usage
// @Accept json
// @Produce json
// @Param work_order_id path int true "Work Order ID"
// @Success 200 {object} common.APIResponse{data=[]repair.PartsInventoryImpact}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/work-orders/{work_order_id}/inventory-impact [get]
func (h *PartsUsageHandler) GetInventoryImpact(c *gin.Context) {
	workOrderID, err := strconv.Atoi(c.Param("work_order_id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid work order ID", err.Error())
		return
	}

	impact, err := h.service.GetInventoryImpact(c.Request.Context(), workOrderID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get inventory impact", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Inventory impact retrieved successfully", impact)
}