package vehicle_purchase

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/vehicle_purchase"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/utils"
)

// TransactionHandler handles vehicle purchase transaction endpoints
type TransactionHandler struct {
	service *services.VehiclePurchaseService
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(service *services.VehiclePurchaseService) *TransactionHandler {
	return &TransactionHandler{
		service: service,
	}
}

// CreateTransaction creates a new vehicle purchase transaction
// @Summary Create vehicle purchase transaction
// @Description Create a new vehicle purchase transaction
// @Tags Vehicle Purchase Transactions
// @Accept json
// @Produce json
// @Param request body vehicle_purchase.CreateVehiclePurchaseTransactionRequest true "Transaction data"
// @Success 201 {object} common.APIResponse{data=vehicle_purchase.VehiclePurchaseTransaction}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req vehicle_purchase.CreateVehiclePurchaseTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	userID := utils.GetUserIDFromContext(c)
	transaction, err := h.service.CreateTransaction(c.Request.Context(), &req, userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Transaction created successfully", transaction)
}

// GetTransactions lists vehicle purchase transactions
// @Summary List vehicle purchase transactions
// @Description Get a paginated list of vehicle purchase transactions
// @Tags Vehicle Purchase Transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status"
// @Param customer_id query int false "Filter by customer ID"
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/transactions [get]
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	params := &vehicle_purchase.VehiclePurchaseTransactionFilterParams{}
	
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
	params.Status = c.Query("status")
	if customerID := c.Query("customer_id"); customerID != "" {
		if cid, err := strconv.Atoi(customerID); err == nil {
			params.CustomerID = cid
		}
	}

	transactions, err := h.service.ListTransactions(c.Request.Context(), params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get transactions", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Transactions retrieved successfully", transactions)
}

// GetTransaction retrieves a vehicle purchase transaction by ID
// @Summary Get vehicle purchase transaction
// @Description Get a vehicle purchase transaction by ID
// @Tags Vehicle Purchase Transactions
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} common.APIResponse{data=vehicle_purchase.VehiclePurchaseTransaction}
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/transactions/{id} [get]
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID", err.Error())
		return
	}

	transaction, err := h.service.GetTransaction(c.Request.Context(), id)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Transaction not found", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Transaction retrieved successfully", transaction)
}

// GetTransactionByNumber retrieves a vehicle purchase transaction by number
// @Summary Get vehicle purchase transaction by number
// @Description Get a vehicle purchase transaction by transaction number
// @Tags Vehicle Purchase Transactions
// @Accept json
// @Produce json
// @Param number path string true "Transaction Number"
// @Success 200 {object} common.APIResponse{data=vehicle_purchase.VehiclePurchaseTransaction}
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/transactions/number/{number} [get]
func (h *TransactionHandler) GetTransactionByNumber(c *gin.Context) {
	number := c.Param("number")
	if number == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Transaction number is required", "")
		return
	}

	transaction, err := h.service.GetTransactionByNumber(c.Request.Context(), number)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Transaction not found", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Transaction retrieved successfully", transaction)
}

// GetTransactionByVIN retrieves a vehicle purchase transaction by VIN
// @Summary Get vehicle purchase transaction by VIN
// @Description Get a vehicle purchase transaction by VIN number
// @Tags Vehicle Purchase Transactions
// @Accept json
// @Produce json
// @Param vin path string true "VIN Number"
// @Success 200 {object} common.APIResponse{data=vehicle_purchase.VehiclePurchaseTransaction}
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/transactions/vin/{vin} [get]
func (h *TransactionHandler) GetTransactionByVIN(c *gin.Context) {
	vin := c.Param("vin")
	if vin == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "VIN number is required", "")
		return
	}

	transaction, err := h.service.GetTransactionByVIN(c.Request.Context(), vin)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Transaction not found", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Transaction retrieved successfully", transaction)
}

// UpdateTransaction updates a vehicle purchase transaction
// @Summary Update vehicle purchase transaction
// @Description Update a vehicle purchase transaction
// @Tags Vehicle Purchase Transactions
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Param request body vehicle_purchase.UpdateVehiclePurchaseTransactionRequest true "Transaction data"
// @Success 200 {object} common.APIResponse{data=vehicle_purchase.VehiclePurchaseTransaction}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/transactions/{id} [put]
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID", err.Error())
		return
	}

	var req vehicle_purchase.UpdateVehiclePurchaseTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	transaction, err := h.service.UpdateTransaction(c.Request.Context(), id, &req)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update transaction", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Transaction updated successfully", transaction)
}

// CompleteInspection completes vehicle inspection
// @Summary Complete vehicle inspection
// @Description Complete inspection for a vehicle purchase transaction
// @Tags Vehicle Purchase Transactions
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Param request body vehicle_purchase.TransactionInspectionRequest true "Inspection data"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/transactions/{id}/inspect [post]
func (h *TransactionHandler) CompleteInspection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID", err.Error())
		return
	}

	var req vehicle_purchase.TransactionInspectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	userID := utils.GetUserIDFromContext(c)
	err = h.service.CompleteInspection(c.Request.Context(), id, &req, userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to complete inspection", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Inspection completed successfully", nil)
}

// ProcessApproval processes transaction approval
// @Summary Process transaction approval
// @Description Process approval for a vehicle purchase transaction
// @Tags Vehicle Purchase Transactions
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Param request body vehicle_purchase.TransactionStatusApprovalRequest true "Approval data"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/transactions/{id}/approve [post]
func (h *TransactionHandler) ProcessApproval(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID", err.Error())
		return
	}

	var req vehicle_purchase.TransactionStatusApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	userID := utils.GetUserIDFromContext(c)
	err = h.service.ProcessApproval(c.Request.Context(), id, &req, userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to process approval", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Approval processed successfully", nil)
}

// GetPendingInspections gets transactions pending inspection
// @Summary Get pending inspections
// @Description Get vehicle purchase transactions pending inspection
// @Tags Vehicle Purchase Transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/transactions/pending-inspection [get]
func (h *TransactionHandler) GetPendingInspections(c *gin.Context) {
	params := &vehicle_purchase.VehiclePurchaseTransactionFilterParams{}
	
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

	transactions, err := h.service.GetPendingInspections(c.Request.Context(), params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get pending inspections", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Pending inspections retrieved successfully", transactions)
}

// GetPendingApprovals gets transactions pending approval
// @Summary Get pending approvals
// @Description Get vehicle purchase transactions pending approval
// @Tags Vehicle Purchase Transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/transactions/pending-approval [get]
func (h *TransactionHandler) GetPendingApprovals(c *gin.Context) {
	params := &vehicle_purchase.VehiclePurchaseTransactionFilterParams{}
	
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

	transactions, err := h.service.GetPendingApprovals(c.Request.Context(), params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get pending approvals", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Pending approvals retrieved successfully", transactions)
}

// GetDashboardStats gets dashboard statistics
// @Summary Get dashboard statistics
// @Description Get vehicle purchase dashboard statistics
// @Tags Vehicle Purchase Transactions
// @Accept json
// @Produce json
// @Success 200 {object} common.APIResponse{data=vehicle_purchase.TransactionDashboardStats}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/dashboard [get]
func (h *TransactionHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.service.GetDashboardStats(c.Request.Context())
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get dashboard stats", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Dashboard stats retrieved successfully", stats)
}