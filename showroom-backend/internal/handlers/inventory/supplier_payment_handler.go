package inventory

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/middleware"
	inventoryModels "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	inventoryServices "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services/inventory"
)

// SupplierPaymentHandler handles HTTP requests for supplier payments
type SupplierPaymentHandler struct {
	supplierPaymentService *inventoryServices.SupplierPaymentService
}

// NewSupplierPaymentHandler creates a new supplier payment handler
func NewSupplierPaymentHandler(supplierPaymentService *inventoryServices.SupplierPaymentService) *SupplierPaymentHandler {
	return &SupplierPaymentHandler{
		supplierPaymentService: supplierPaymentService,
	}
}

// CreateSupplierPayment handles POST /supplier-payments
func (h *SupplierPaymentHandler) CreateSupplierPayment(c *gin.Context) {
	var req inventoryModels.SupplierPaymentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Get user ID from context
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, common.ErrorResponse{
			Status:  "error",
			Message: "Unauthorized",
			Error:   "Invalid authentication",
		})
		return
	}

	payment, err := h.supplierPaymentService.CreateSupplierPayment(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Status:  "error",
			Message: "Failed to create supplier payment",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, common.SuccessResponse{
		Status:  "success",
		Message: "Supplier payment created successfully",
		Data:    payment,
	})
}

// GetSupplierPayment handles GET /supplier-payments/:id
func (h *SupplierPaymentHandler) GetSupplierPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid payment ID",
		})
		return
	}

	payment, err := h.supplierPaymentService.GetSupplierPayment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.ErrorResponse{
			Status:  "error",
			Message: "Supplier payment not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse{
		Status:  "success",
		Message: "Supplier payment retrieved successfully",
		Data:    payment,
	})
}

// GetSupplierPayments handles GET /supplier-payments
func (h *SupplierPaymentHandler) GetSupplierPayments(c *gin.Context) {
	var filter inventoryModels.SupplierPaymentFilterParams

	// Parse query parameters
	if supplierIDStr := c.Query("supplier_id"); supplierIDStr != "" {
		supplierID, err := strconv.Atoi(supplierIDStr)
		if err == nil {
			filter.SupplierID = &supplierID
		}
	}

	if poIDStr := c.Query("po_id"); poIDStr != "" {
		poID, err := strconv.Atoi(poIDStr)
		if err == nil {
			filter.PoID = &poID
		}
	}

	if status := c.Query("status"); status != "" {
		paymentStatus := inventoryModels.PaymentStatus(status)
		filter.PaymentStatus = &paymentStatus
	}

	if method := c.Query("method"); method != "" {
		paymentMethod := inventoryModels.PaymentMethod(method)
		filter.PaymentMethod = &paymentMethod
	}

	// Parse pagination
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	payments, total, err := h.supplierPaymentService.ListSupplierPayments(c.Request.Context(), &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve supplier payments",
			Error:   err.Error(),
		})
		return
	}

totalPages := (total + limit - 1) / limit
hasMore := page < totalPages

meta := common.PaginationMeta{
Page:       page,
Limit:      limit,
Total:      total,
TotalPages: totalPages,
HasMore:    hasMore,
}

c.JSON(http.StatusOK, common.NewPaginationResponse(
"Supplier payments retrieved successfully", payments, meta,
))
}

// UpdateSupplierPayment handles PUT /supplier-payments/:id
func (h *SupplierPaymentHandler) UpdateSupplierPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid payment ID",
		})
		return
	}

	var req inventoryModels.SupplierPaymentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	payment, err := h.supplierPaymentService.UpdateSupplierPayment(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Status:  "error",
			Message: "Failed to update supplier payment",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse{
		Status:  "success",
		Message: "Supplier payment updated successfully",
		Data:    payment,
	})
}

// DeleteSupplierPayment handles DELETE /supplier-payments/:id
func (h *SupplierPaymentHandler) DeleteSupplierPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid payment ID",
		})
		return
	}

	err = h.supplierPaymentService.DeleteSupplierPayment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete supplier payment",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse{
		Status:  "success",
		Message: "Supplier payment deleted successfully",
	})
}

// GetPaymentsBySupplier handles GET /suppliers/:supplier_id/payments
func (h *SupplierPaymentHandler) GetPaymentsBySupplier(c *gin.Context) {
	supplierIDStr := c.Param("supplier_id")
	supplierID, err := strconv.Atoi(supplierIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{
			Status:  "error",
			Message: "Invalid supplier ID",
		})
		return
	}

	filter := inventoryModels.SupplierPaymentFilterParams{
		SupplierID: &supplierID,
	}

	// Parse pagination
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	payments, total, err := h.supplierPaymentService.ListSupplierPayments(c.Request.Context(), &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve supplier payments",
			Error:   err.Error(),
		})
		return
	}

totalPages := (total + limit - 1) / limit
hasMore := page < totalPages

meta := common.PaginationMeta{
Page:       page,
Limit:      limit,
Total:      total,
TotalPages: totalPages,
HasMore:    hasMore,
}

c.JSON(http.StatusOK, common.NewPaginationResponse(
"Supplier payments retrieved successfully", payments, meta,
))
}