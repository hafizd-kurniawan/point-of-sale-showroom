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

// DamageHandler handles vehicle damage endpoints
type DamageHandler struct {
	service *services.RepairService
}

// NewDamageHandler creates a new damage handler
func NewDamageHandler(service *services.RepairService) *DamageHandler {
	return &DamageHandler{
		service: service,
	}
}

// CreateDamage creates a new vehicle damage record
// @Summary Create vehicle damage
// @Description Create a new vehicle damage record
// @Tags Vehicle Damages
// @Accept json
// @Produce json
// @Param request body repair.CreateVehicleDamageRequest true "Damage data"
// @Success 201 {object} common.APIResponse{data=repair.VehicleDamage}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/damages [post]
func (h *DamageHandler) CreateDamage(c *gin.Context) {
	var req repair.CreateVehicleDamageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	userID := utils.GetUserIDFromContext(c)
	damage, err := h.service.CreateDamage(c.Request.Context(), &req, userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create damage", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Damage created successfully", damage)
}

// GetDamages lists vehicle damages
// @Summary List vehicle damages
// @Description Get a paginated list of vehicle damages
// @Tags Vehicle Damages
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param transaction_id query int false "Filter by transaction ID"
// @Param category query string false "Filter by damage category"
// @Param severity query string false "Filter by damage severity"
// @Param status query string false "Filter by damage status"
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/damages [get]
func (h *DamageHandler) GetDamages(c *gin.Context) {
	params := &repair.VehicleDamageFilterParams{}
	
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
	params.Category = c.Query("category")
	params.Severity = c.Query("severity")
	params.Status = c.Query("status")

	damages, err := h.service.ListDamages(c.Request.Context(), params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get damages", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Damages retrieved successfully", damages)
}

// GetDamage retrieves a vehicle damage by ID
// @Summary Get vehicle damage
// @Description Get a vehicle damage by ID
// @Tags Vehicle Damages
// @Accept json
// @Produce json
// @Param id path int true "Damage ID"
// @Success 200 {object} common.APIResponse{data=repair.VehicleDamage}
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/damages/{id} [get]
func (h *DamageHandler) GetDamage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid damage ID", err.Error())
		return
	}

	damage, err := h.service.GetDamage(c.Request.Context(), id)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Damage not found", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Damage retrieved successfully", damage)
}

// UpdateDamage updates a vehicle damage record
// @Summary Update vehicle damage
// @Description Update a vehicle damage record
// @Tags Vehicle Damages
// @Accept json
// @Produce json
// @Param id path int true "Damage ID"
// @Param request body repair.UpdateVehicleDamageRequest true "Damage data"
// @Success 200 {object} common.APIResponse{data=repair.VehicleDamage}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/damages/{id} [put]
func (h *DamageHandler) UpdateDamage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid damage ID", err.Error())
		return
	}

	var req repair.UpdateVehicleDamageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	damage, err := h.service.UpdateDamage(c.Request.Context(), id, &req)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update damage", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Damage updated successfully", damage)
}

// AssessDamage assesses damage with specific criteria
// @Summary Assess damage
// @Description Assess damage with specific criteria and cost estimation
// @Tags Vehicle Damages
// @Accept json
// @Produce json
// @Param id path int true "Damage ID"
// @Param request body repair.DamageAssessmentRequest true "Assessment data"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/damages/{id}/assess [post]
func (h *DamageHandler) AssessDamage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid damage ID", err.Error())
		return
	}

	var req repair.DamageAssessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	err = h.service.AssessDamage(c.Request.Context(), id, &req)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to assess damage", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Damage assessed successfully", nil)
}

// GetHighPriorityDamages gets high priority damages
// @Summary Get high priority damages
// @Description Get vehicle damages with high priority or critical severity
// @Tags Vehicle Damages
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/damages/high-priority [get]
func (h *DamageHandler) GetHighPriorityDamages(c *gin.Context) {
	params := &repair.VehicleDamageFilterParams{}
	
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

	damages, err := h.service.GetHighPriorityDamages(c.Request.Context(), params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get high priority damages", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "High priority damages retrieved successfully", damages)
}

// GetDamageSummary gets damage summary for a transaction
// @Summary Get damage summary
// @Description Get damage summary for a specific transaction
// @Tags Vehicle Damages
// @Accept json
// @Produce json
// @Param transaction_id path int true "Transaction ID"
// @Success 200 {object} common.APIResponse{data=repair.DamageSummary}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/repairs/transactions/{transaction_id}/damage-summary [get]
func (h *DamageHandler) GetDamageSummary(c *gin.Context) {
	transactionID, err := strconv.Atoi(c.Param("transaction_id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID", err.Error())
		return
	}

	summary, err := h.service.GetDamageSummary(c.Request.Context(), transactionID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get damage summary", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Damage summary retrieved successfully", summary)
}