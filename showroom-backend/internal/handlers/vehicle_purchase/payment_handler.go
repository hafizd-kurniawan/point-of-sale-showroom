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

// PaymentHandler handles vehicle purchase payment endpoints
type PaymentHandler struct {
	service *services.VehiclePurchaseService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(service *services.VehiclePurchaseService) *PaymentHandler {
	return &PaymentHandler{
		service: service,
	}
}

// CreatePayment creates a new vehicle purchase payment
// @Summary Create vehicle purchase payment
// @Description Create a new vehicle purchase payment
// @Tags Vehicle Purchase Payments
// @Accept json
// @Produce json
// @Param request body vehicle_purchase.CreateVehiclePurchasePaymentRequest true "Payment data"
// @Success 201 {object} common.APIResponse{data=vehicle_purchase.VehiclePurchasePayment}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/payments [post]
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req vehicle_purchase.CreateVehiclePurchasePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	userID := utils.GetUserIDFromContext(c)
	payment, err := h.service.CreatePayment(c.Request.Context(), &req, userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create payment", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Payment created successfully", payment)
}

// GetPayments lists vehicle purchase payments
// @Summary List vehicle purchase payments
// @Description Get a paginated list of vehicle purchase payments
// @Tags Vehicle Purchase Payments
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status"
// @Param transaction_id query int false "Filter by transaction ID"
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/payments [get]
func (h *PaymentHandler) GetPayments(c *gin.Context) {
	params := &vehicle_purchase.VehiclePurchasePaymentFilterParams{}
	
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
	if transactionID := c.Query("transaction_id"); transactionID != "" {
		if tid, err := strconv.Atoi(transactionID); err == nil {
			params.TransactionID = tid
		}
	}

	payments, err := h.service.ListPayments(c.Request.Context(), params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get payments", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Payments retrieved successfully", payments)
}

// GetPayment retrieves a vehicle purchase payment by ID
// @Summary Get vehicle purchase payment
// @Description Get a vehicle purchase payment by ID
// @Tags Vehicle Purchase Payments
// @Accept json
// @Produce json
// @Param id path int true "Payment ID"
// @Success 200 {object} common.APIResponse{data=vehicle_purchase.VehiclePurchasePayment}
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/payments/{id} [get]
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid payment ID", err.Error())
		return
	}

	payment, err := h.service.GetPayment(c.Request.Context(), id)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Payment not found", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Payment retrieved successfully", payment)
}

// GetPaymentByNumber retrieves a vehicle purchase payment by number
// @Summary Get vehicle purchase payment by number
// @Description Get a vehicle purchase payment by payment number
// @Tags Vehicle Purchase Payments
// @Accept json
// @Produce json
// @Param number path string true "Payment Number"
// @Success 200 {object} common.APIResponse{data=vehicle_purchase.VehiclePurchasePayment}
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/payments/number/{number} [get]
func (h *PaymentHandler) GetPaymentByNumber(c *gin.Context) {
	number := c.Param("number")
	if number == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Payment number is required", "")
		return
	}

	payment, err := h.service.GetPaymentByNumber(c.Request.Context(), number)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Payment not found", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Payment retrieved successfully", payment)
}

// ProcessPayment processes a payment
// @Summary Process payment
// @Description Process a vehicle purchase payment
// @Tags Vehicle Purchase Payments
// @Accept json
// @Produce json
// @Param id path int true "Payment ID"
// @Param request body vehicle_purchase.PaymentProcessRequest true "Payment process data"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/payments/{id}/process [post]
func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid payment ID", err.Error())
		return
	}

	var req vehicle_purchase.PaymentProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	userID := utils.GetUserIDFromContext(c)
	err = h.service.ProcessPayment(c.Request.Context(), id, &req, userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to process payment", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Payment processed successfully", nil)
}

// ProcessPaymentApproval processes payment approval
// @Summary Process payment approval
// @Description Process approval for a vehicle purchase payment
// @Tags Vehicle Purchase Payments
// @Accept json
// @Produce json
// @Param id path int true "Payment ID"
// @Param request body vehicle_purchase.PaymentApprovalRequest true "Payment approval data"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/payments/{id}/approve [post]
func (h *PaymentHandler) ProcessPaymentApproval(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid payment ID", err.Error())
		return
	}

	var req vehicle_purchase.PaymentApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	userID := utils.GetUserIDFromContext(c)
	err = h.service.ProcessPaymentApproval(c.Request.Context(), id, &req, userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to process payment approval", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Payment approval processed successfully", nil)
}

// GetPaymentsByTransaction gets payments for a specific transaction
// @Summary Get payments by transaction
// @Description Get vehicle purchase payments for a specific transaction
// @Tags Vehicle Purchase Payments
// @Accept json
// @Produce json
// @Param transaction_id path int true "Transaction ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/transactions/{transaction_id}/payments [get]
func (h *PaymentHandler) GetPaymentsByTransaction(c *gin.Context) {
	transactionID, err := strconv.Atoi(c.Param("transaction_id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID", err.Error())
		return
	}

	params := &vehicle_purchase.VehiclePurchasePaymentFilterParams{}
	
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

	payments, err := h.service.GetPaymentsByTransaction(c.Request.Context(), transactionID, params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get payments by transaction", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Payments retrieved successfully", payments)
}

// GetPendingApprovals gets payments pending approval
// @Summary Get pending payment approvals
// @Description Get vehicle purchase payments pending approval
// @Tags Vehicle Purchase Payments
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/payments/pending-approval [get]
func (h *PaymentHandler) GetPendingApprovals(c *gin.Context) {
	params := &vehicle_purchase.VehiclePurchasePaymentFilterParams{}
	
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

	payments, err := h.service.GetPendingPaymentApprovals(c.Request.Context(), params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get pending approvals", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Pending approvals retrieved successfully", payments)
}

// GetOverduePayments gets overdue payments
// @Summary Get overdue payments
// @Description Get overdue vehicle purchase payments
// @Tags Vehicle Purchase Payments
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/payments/overdue [get]
func (h *PaymentHandler) GetOverduePayments(c *gin.Context) {
	params := &vehicle_purchase.VehiclePurchasePaymentFilterParams{}
	
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

	payments, err := h.service.GetOverduePayments(c.Request.Context(), params)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get overdue payments", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Overdue payments retrieved successfully", payments)
}

// GetPaymentSummary gets payment summary for a transaction
// @Summary Get payment summary
// @Description Get payment summary for a specific transaction
// @Tags Vehicle Purchase Payments
// @Accept json
// @Produce json
// @Param transaction_id path int true "Transaction ID"
// @Success 200 {object} common.APIResponse{data=vehicle_purchase.PaymentSummary}
// @Failure 500 {object} common.ErrorResponse
// @Router /api/v1/vehicle-purchases/transactions/{transaction_id}/payment-summary [get]
func (h *PaymentHandler) GetPaymentSummary(c *gin.Context) {
	transactionID, err := strconv.Atoi(c.Param("transaction_id"))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID", err.Error())
		return
	}

	summary, err := h.service.GetPaymentSummary(c.Request.Context(), transactionID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get payment summary", err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Payment summary retrieved successfully", summary)
}