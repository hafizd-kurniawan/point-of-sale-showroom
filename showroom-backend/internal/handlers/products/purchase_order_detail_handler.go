package products

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
)

// PurchaseOrderDetailHandler handles purchase order detail HTTP requests
type PurchaseOrderDetailHandler struct {
	poDetailRepo interfaces.PurchaseOrderDetailRepository
}

// NewPurchaseOrderDetailHandler creates a new purchase order detail handler
func NewPurchaseOrderDetailHandler(poDetailRepo interfaces.PurchaseOrderDetailRepository) *PurchaseOrderDetailHandler {
	return &PurchaseOrderDetailHandler{
		poDetailRepo: poDetailRepo,
	}
}

// CreatePODetail handles adding a line item to a purchase order
func (h *PurchaseOrderDetailHandler) CreatePODetail(c *gin.Context) {
	poIDStr := c.Param("id")
	poID, err := strconv.Atoi(poIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid PO ID", "PO ID must be a valid number",
		))
		return
	}

	var req products.PurchaseOrderDetailCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	// Create detail model
	detail := &products.PurchaseOrderDetail{
		POID:            poID,
		ProductID:       req.ProductID,
		ItemDescription: req.ItemDescription,
		QuantityOrdered: req.QuantityOrdered,
		UnitCost:        req.UnitCost,
		ExpectedDate:    req.ExpectedDate,
		ItemNotes:       req.ItemNotes,
	}

	createdDetail, err := h.poDetailRepo.Create(c.Request.Context(), detail)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to create PO detail", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"PO detail created successfully", createdDetail,
	))
}

// GetPODetails handles getting details for a purchase order
func (h *PurchaseOrderDetailHandler) GetPODetails(c *gin.Context) {
	poIDStr := c.Param("id")
	poID, err := strconv.Atoi(poIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid PO ID", "PO ID must be a valid number",
		))
		return
	}

	// Parse query parameters
	var params products.PurchaseOrderDetailFilterParams
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

	details, err := h.poDetailRepo.GetByPOID(c.Request.Context(), poID, &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get PO details", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"PO details retrieved successfully", details,
	))
}

// GetPODetail handles getting a specific purchase order detail
func (h *PurchaseOrderDetailHandler) GetPODetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid detail ID", "Detail ID must be a valid number",
		))
		return
	}

	detail, err := h.poDetailRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"PO detail not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"PO detail retrieved successfully", detail,
	))
}

// UpdatePODetail handles updating a purchase order detail
func (h *PurchaseOrderDetailHandler) UpdatePODetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid detail ID", "Detail ID must be a valid number",
		))
		return
	}

	var req products.PurchaseOrderDetailUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	// Get existing detail
	existing, err := h.poDetailRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"PO detail not found", err.Error(),
		))
		return
	}

	// Update fields if provided
	if req.ProductID != nil {
		existing.ProductID = *req.ProductID
	}
	if req.ItemDescription != nil {
		existing.ItemDescription = req.ItemDescription
	}
	if req.QuantityOrdered != nil {
		existing.QuantityOrdered = *req.QuantityOrdered
	}
	if req.UnitCost != nil {
		existing.UnitCost = *req.UnitCost
	}
	if req.ExpectedDate != nil {
		existing.ExpectedDate = req.ExpectedDate
	}
	if req.ItemNotes != nil {
		existing.ItemNotes = req.ItemNotes
	}

	updatedDetail, err := h.poDetailRepo.Update(c.Request.Context(), id, existing)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to update PO detail", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"PO detail updated successfully", updatedDetail,
	))
}

// DeletePODetail handles deleting a purchase order detail
func (h *PurchaseOrderDetailHandler) DeletePODetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid detail ID", "Detail ID must be a valid number",
		))
		return
	}

	err = h.poDetailRepo.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to delete PO detail", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"PO detail deleted successfully", nil,
	))
}

// GetProductPODetails handles getting PO details for a specific product
func (h *PurchaseOrderDetailHandler) GetProductPODetails(c *gin.Context) {
	productIDStr := c.Param("productId")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid product ID", "Product ID must be a valid number",
		))
		return
	}

	// Parse query parameters
	var params products.PurchaseOrderDetailFilterParams
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

	details, err := h.poDetailRepo.GetByProductID(c.Request.Context(), productID, &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get product PO details", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Product PO details retrieved successfully", details,
	))
}

// GetPendingReceiptItems handles getting items pending receipt for a PO
func (h *PurchaseOrderDetailHandler) GetPendingReceiptItems(c *gin.Context) {
	poIDStr := c.Param("id")
	poID, err := strconv.Atoi(poIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid PO ID", "PO ID must be a valid number",
		))
		return
	}

	items, err := h.poDetailRepo.GetPendingReceiptItems(c.Request.Context(), poID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to get pending receipt items", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Pending receipt items retrieved successfully", items,
	))
}

// BulkCreatePODetails handles creating multiple PO details at once
func (h *PurchaseOrderDetailHandler) BulkCreatePODetails(c *gin.Context) {
	poIDStr := c.Param("id")
	poID, err := strconv.Atoi(poIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid PO ID", "PO ID must be a valid number",
		))
		return
	}

	var requests []products.PurchaseOrderDetailCreateRequest
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	// Convert requests to detail models
	var details []products.PurchaseOrderDetail
	for _, req := range requests {
		detail := products.PurchaseOrderDetail{
			POID:            poID,
			ProductID:       req.ProductID,
			ItemDescription: req.ItemDescription,
			QuantityOrdered: req.QuantityOrdered,
			UnitCost:        req.UnitCost,
			ExpectedDate:    req.ExpectedDate,
			ItemNotes:       req.ItemNotes,
		}
		details = append(details, detail)
	}

	err = h.poDetailRepo.BulkCreate(c.Request.Context(), details)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to create PO details", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"PO details created successfully", nil,
	))
}