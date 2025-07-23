package products

import (
	"net/http"
	"strconv"

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

	c.JSON(http.StatusCreated, common.NewSuccessResponse("Supplier payment created successfully", payment))
}

// GetSupplierPayments handles supplier payments list with pagination and filtering
func (h *SupplierPaymentHandler) GetSupplierPayments(c *gin.Context) {
	var params products.SupplierPaymentFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Query parameter validation failed", err.Error(),
		))
		return
	}

	result, err := h.supplierPaymentService.GetSupplierPayments(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get supplier payments", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Supplier payments retrieved successfully", result))
}

// GetSupplierPaymentByID handles supplier payment retrieval by ID
func (h *SupplierPaymentHandler) GetSupplierPaymentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid payment ID", "Payment ID must be a valid integer",
		))
		return
	}

	payment, err := h.supplierPaymentService.GetSupplierPaymentByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Supplier payment not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Supplier payment retrieved successfully", payment))
}

// UpdateSupplierPayment handles supplier payment updates
func (h *SupplierPaymentHandler) UpdateSupplierPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid payment ID", "Payment ID must be a valid integer",
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

	payment, err := h.supplierPaymentService.UpdateSupplierPayment(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Supplier payment update failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Supplier payment updated successfully", payment))
}

// AddPayment handles adding payment to existing record
func (h *SupplierPaymentHandler) AddPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid payment ID", "Payment ID must be a valid integer",
		))
		return
	}

	var req struct {
		Amount           float64                `json:"amount" binding:"required,min=0.01"`
		PaymentMethod    products.PaymentMethod `json:"payment_method" binding:"required"`
		PaymentReference *string                `json:"payment_reference,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	err = h.supplierPaymentService.AddPayment(c.Request.Context(), id, req.Amount, req.PaymentMethod, req.PaymentReference)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to add payment", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Payment added successfully", nil))
}

// UpdatePaymentStatus handles payment status updates
func (h *SupplierPaymentHandler) UpdatePaymentStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid payment ID", "Payment ID must be a valid integer",
		))
		return
	}

	var req struct {
		Status products.PaymentStatus `json:"status" binding:"required"`
	}

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

	c.JSON(http.StatusOK, common.NewSuccessResponse("Payment status updated successfully", nil))
}

// DeleteSupplierPayment handles supplier payment deletion
func (h *SupplierPaymentHandler) DeleteSupplierPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid payment ID", "Payment ID must be a valid integer",
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

	c.JSON(http.StatusOK, common.NewSuccessResponse("Supplier payment deleted successfully", nil))
}

// GetSupplierPaymentsBySupplier handles payments by supplier
func (h *SupplierPaymentHandler) GetSupplierPaymentsBySupplier(c *gin.Context) {
	supplierIDStr := c.Param("supplierId")
	supplierID, err := strconv.Atoi(supplierIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid supplier ID", "Supplier ID must be a valid integer",
		))
		return
	}

	var params products.SupplierPaymentFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Query parameter validation failed", err.Error(),
		))
		return
	}

	result, err := h.supplierPaymentService.GetSupplierPaymentsBySupplier(c.Request.Context(), supplierID, &params)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to get supplier payments", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Supplier payments retrieved successfully", result))
}

// GetOverduePayments handles overdue payments
func (h *SupplierPaymentHandler) GetOverduePayments(c *gin.Context) {
	var params products.SupplierPaymentFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Query parameter validation failed", err.Error(),
		))
		return
	}

	result, err := h.supplierPaymentService.GetOverduePayments(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get overdue payments", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Overdue payments retrieved successfully", result))
}

// GetPaymentSummary handles payment summary
func (h *SupplierPaymentHandler) GetPaymentSummary(c *gin.Context) {
	var supplierID *int
	supplierIDStr := c.Query("supplier_id")
	if supplierIDStr != "" {
		id, err := strconv.Atoi(supplierIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(
				"Invalid supplier ID", "Supplier ID must be a valid integer",
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

	c.JSON(http.StatusOK, common.NewSuccessResponse("Payment summary retrieved successfully", summary))
}

// ProcessOverduePayments handles updating overdue status
func (h *SupplierPaymentHandler) ProcessOverduePayments(c *gin.Context) {
	err := h.supplierPaymentService.ProcessOverduePayments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to process overdue payments", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse("Overdue payments processed successfully", nil))
}