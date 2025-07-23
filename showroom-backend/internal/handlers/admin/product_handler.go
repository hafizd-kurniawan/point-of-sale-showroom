package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/middleware"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/products"
	productService "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services/products"
)

// ProductHandler handles product HTTP requests
type ProductHandler struct {
	productService *productService.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService *productService.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct handles product creation
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req products.ProductSparePartCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	createdBy := middleware.GetCurrentUserID(c)
	if createdBy == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Invalid user", "Creator user ID not found",
		))
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Product creation failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Product created successfully", product,
	))
}

// GetProducts handles product listing
func (h *ProductHandler) GetProducts(c *gin.Context) {
	var params products.ProductSparePartFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid query parameters", err.Error(),
		))
		return
	}

	response, err := h.productService.ListProducts(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve products", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProduct handles retrieving a specific product
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Product ID must be a valid integer",
		))
		return
	}

	product, err := h.productService.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Product not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Product retrieved successfully", product,
	))
}

// UpdateProduct handles product updates
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Product ID must be a valid integer",
		))
		return
	}

	var req products.ProductSparePartUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	product, err := h.productService.UpdateProduct(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Product update failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Product updated successfully", product,
	))
}

// DeleteProduct handles product deletion
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Product ID must be a valid integer",
		))
		return
	}

	err = h.productService.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Product deletion failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Product deleted successfully", nil,
	))
}

// GetLowStockProducts handles retrieving products with low stock
func (h *ProductHandler) GetLowStockProducts(c *gin.Context) {
	var params products.ProductSparePartFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid query parameters", err.Error(),
		))
		return
	}

	response, err := h.productService.GetLowStockProducts(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve low stock products", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, response)
}

// SearchProducts handles product search
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Missing search query", "Query parameter 'q' is required",
		))
		return
	}

	var params products.ProductSparePartFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid query parameters", err.Error(),
		))
		return
	}

	response, err := h.productService.SearchProducts(c.Request.Context(), query, &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Search failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProductByCode handles retrieving product by code
func (h *ProductHandler) GetProductByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid code", "Product code is required",
		))
		return
	}

	product, err := h.productService.GetProductByCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Product not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Product retrieved successfully", product,
	))
}

// GetProductByBarcode handles retrieving product by barcode
func (h *ProductHandler) GetProductByBarcode(c *gin.Context) {
	barcode := c.Param("barcode")
	if barcode == "" {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid barcode", "Product barcode is required",
		))
		return
	}

	product, err := h.productService.GetProductByBarcode(c.Request.Context(), barcode)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Product not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Product retrieved successfully", product,
	))
}

// AdjustStock handles stock adjustment
func (h *ProductHandler) AdjustStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Product ID must be a valid integer",
		))
		return
	}

	var req struct {
		NewQuantity int    `json:"new_quantity" binding:"required,min=0"`
		Reason      string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	adjustedBy := middleware.GetCurrentUserID(c)
	if adjustedBy == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Invalid user", "User ID not found",
		))
		return
	}

	err = h.productService.AdjustStock(c.Request.Context(), id, req.NewQuantity, req.Reason, adjustedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Stock adjustment failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Stock adjusted successfully", nil,
	))
}

// GetStockMovementHistory handles retrieving stock movement history
func (h *ProductHandler) GetStockMovementHistory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Product ID must be a valid integer",
		))
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	movements, err := h.productService.GetStockMovementHistory(c.Request.Context(), id, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to get stock movement history", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Stock movement history retrieved successfully", movements,
	))
}