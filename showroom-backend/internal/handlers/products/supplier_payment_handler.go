package products

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/middleware"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
	productService "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services/products"
)

// SupplierPaymentHandler handles supplier payment HTTP requests
type SupplierPaymentHandler struct {
	supplierPaymentService *productService.SupplierPaymentService
}

// NewSupplierPaymentHandler creates a new supplier payment handler
func NewSupplierPaymentHandler(supplierPaymentService *productService.SupplierPaymentService) *SupplierPaymentHandler {
	return &SupplierPaymentHandler{
		supplierPaymentService: supplierPaymentService,
	}
}

// CreateSupplierPayment handles supplier payment creation
func (h *SupplierPaymentHandler) CreateSupplierPayment(c *gin.Context) {
	var req products.SupplierPaymentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
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

	payment, err := h.supplierPaymentService.CreateSupplierPayment(c.Request.Context(), &req, processedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Supplier payment creation failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Supplier payment created successfully", payment,
	))
}

// GetSupplierPayment handles getting a specific supplier payment
func (h *SupplierPaymentHandler) GetSupplierPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid payment ID", "Payment ID must be a valid number",
		))
		return
	}

	payment, err := h.supplierPaymentService.GetSupplierPayment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Supplier payment not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Supplier payment retrieved successfully", payment,
	))
}

// UpdateSupplierPayment handles updating a supplier payment
func (h *SupplierPaymentHandler) UpdateSupplierPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid payment ID", "Payment ID must be a valid number",
		))
		return
	}

	var req products.SupplierPaymentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	updatedPayment, err := h.supplierPaymentService.UpdateSupplierPayment(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to update supplier payment", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Supplier payment updated successfully", updatedPayment,
	))
}

// DeleteSupplierPayment handles deleting a supplier payment
func (h *SupplierPaymentHandler) DeleteSupplierPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid payment ID", "Payment ID must be a valid number",
		))
		return
	}

	err = h.supplierPaymentService.DeleteSupplierPayment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to delete supplier payment", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Supplier payment deleted successfully", nil,
	))
}

// ListSupplierPayments handles listing supplier payments with pagination
func (h *SupplierPaymentHandler) ListSupplierPayments(c *gin.Context) {
	// Parse query parameters
	var params products.SupplierPaymentFilterParams
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

	payments, err := h.supplierPaymentService.ListSupplierPayments(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to list supplier payments", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Supplier payments retrieved successfully", payments,
	))
}

// GetSupplierPayments handles getting payments for a specific supplier
func (h *SupplierPaymentHandler) GetSupplierPayments(c *gin.Context) {
	supplierIDStr := c.Param("supplierId")
	supplierID, err := strconv.Atoi(supplierIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid supplier ID", "Supplier ID must be a valid number",
		))
		return
	}

	// Parse query parameters
	var params products.SupplierPaymentFilterParams
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

	payments, err := h.supplierPaymentService.GetSupplierPayments(c.Request.Context(), supplierID, &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get supplier payments", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Supplier payments retrieved successfully", payments,
	))
}

// GetPOPayments handles getting payments for a specific purchase order
func (h *SupplierPaymentHandler) GetPOPayments(c *gin.Context) {
	poIDStr := c.Param("poId")
	poID, err := strconv.Atoi(poIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid PO ID", "PO ID must be a valid number",
		))
		return
	}

	// Parse query parameters
	var params products.SupplierPaymentFilterParams
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

	payments, err := h.supplierPaymentService.GetPOPayments(c.Request.Context(), poID, &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get PO payments", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"PO payments retrieved successfully", payments,
	))
}

// GetOverduePayments handles getting overdue payments
func (h *SupplierPaymentHandler) GetOverduePayments(c *gin.Context) {
	// Parse query parameters
	var params products.SupplierPaymentFilterParams
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

	payments, err := h.supplierPaymentService.GetOverduePayments(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get overdue payments", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Overdue payments retrieved successfully", payments,
	))
}

// ProcessPayment handles processing a payment for an invoice
func (h *SupplierPaymentHandler) ProcessPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid payment ID", "Payment ID must be a valid number",
		))
		return
	}

	var req productService.PaymentProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	err = h.supplierPaymentService.ProcessPayment(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to process payment", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Payment processed successfully", nil,
	))
}

// UpdatePaymentStatus handles updating the payment status
func (h *SupplierPaymentHandler) UpdatePaymentStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid payment ID", "Payment ID must be a valid number",
		))
		return
	}

	var req UpdatePaymentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	err = h.supplierPaymentService.UpdatePaymentStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to update payment status", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Payment status updated successfully", nil,
	))
}

// GetPaymentSummary handles getting payment summary
func (h *SupplierPaymentHandler) GetPaymentSummary(c *gin.Context) {
	var supplierID *int
	
	// Check if supplier ID is provided in query parameters
	supplierIDStr := c.Query("supplier_id")
	if supplierIDStr != "" {
		id, err := strconv.Atoi(supplierIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(
				"Invalid supplier ID", "Supplier ID must be a valid number",
			))
			return
		}
		supplierID = &id
	}

	summary, err := h.supplierPaymentService.GetPaymentSummary(c.Request.Context(), supplierID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get payment summary", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Payment summary retrieved successfully", summary,
	))
}

// CreatePaymentFromPO handles creating a payment from a purchase order
func (h *SupplierPaymentHandler) CreatePaymentFromPO(c *gin.Context) {
	poIDStr := c.Param("poId")
	poID, err := strconv.Atoi(poIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid PO ID", "PO ID must be a valid number",
		))
		return
	}

	var req productService.POPaymentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
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

	payment, err := h.supplierPaymentService.CreatePaymentFromPO(c.Request.Context(), poID, &req, processedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to create payment from PO", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Payment created from PO successfully", payment,
	))
}

// CalculatePaymentTerms handles calculating payment terms
func (h *SupplierPaymentHandler) CalculatePaymentTerms(c *gin.Context) {
	var req PaymentTermsCalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	terms := h.supplierPaymentService.CalculatePaymentTerms(req.InvoiceDate, req.TermsDays)

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Payment terms calculated successfully", terms,
	))
}

// UpdateOverduePayments handles updating overdue status for all payments
func (h *SupplierPaymentHandler) UpdateOverduePayments(c *gin.Context) {
	err := h.supplierPaymentService.UpdateOverduePayments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to update overdue payments", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Overdue payments updated successfully", nil,
	))
}

// UpdatePaymentStatusRequest represents a request to update payment status
type UpdatePaymentStatusRequest struct {
	Status products.PaymentStatus `json:"status" binding:"required"`
}

// PaymentTermsCalculationRequest represents a request to calculate payment terms
type PaymentTermsCalculationRequest struct {
	InvoiceDate time.Time `json:"invoice_date" binding:"required"`
	TermsDays   int       `json:"terms_days" binding:"required,min=0"`
}