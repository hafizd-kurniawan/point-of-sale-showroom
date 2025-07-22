package inventory

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/middleware"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/inventory"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services/inventory"
)

// ProductHandler handles product HTTP requests
type ProductHandler struct {
	productService *inventory.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService *inventory.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct handles product creation
// @Summary Create product
// @Description Create a new spare part product
// @Tags inventory
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body inventory.ProductSparePartCreateRequest true "Create product request"
// @Success 201 {object} common.APIResponse{data=inventory.ProductSparePart}
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Router /inventory/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req inventory.ProductSparePartCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	createdBy := middleware.GetCurrentUserID(c)
	if createdBy == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Unauthorized", "Invalid authentication", "USER_NOT_AUTHENTICATED",
		))
		return
	}

	product, err := h.productService.Create(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to create product", err.Error(), "CREATION_FAILED",
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Product created successfully", product,
	))
}

// GetProduct handles getting a product by ID
// @Summary Get product
// @Description Get a product by ID
// @Tags inventory
// @Security BearerAuth
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} common.APIResponse{data=inventory.ProductSparePart}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /inventory/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Product ID must be a valid integer", "INVALID_ID",
		))
		return
	}

	product, err := h.productService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Product not found", err.Error(), "PRODUCT_NOT_FOUND",
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Product retrieved successfully", product,
	))
}

// GetProductByCode handles getting a product by code
// @Summary Get product by code
// @Description Get a product by product code
// @Tags inventory
// @Security BearerAuth
// @Produce json
// @Param code path string true "Product Code"
// @Success 200 {object} common.APIResponse{data=inventory.ProductSparePart}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /inventory/products/code/{code} [get]
func (h *ProductHandler) GetProductByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid code", "Product code is required", "INVALID_CODE",
		))
		return
	}

	product, err := h.productService.GetByCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Product not found", err.Error(), "PRODUCT_NOT_FOUND",
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Product retrieved successfully", product,
	))
}

// GetProductByBarcode handles getting a product by barcode
// @Summary Get product by barcode
// @Description Get a product by barcode
// @Tags inventory
// @Security BearerAuth
// @Produce json
// @Param barcode path string true "Product Barcode"
// @Success 200 {object} common.APIResponse{data=inventory.ProductSparePart}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /inventory/products/barcode/{barcode} [get]
func (h *ProductHandler) GetProductByBarcode(c *gin.Context) {
	barcode := c.Param("barcode")
	if barcode == "" {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid barcode", "Product barcode is required", "INVALID_BARCODE",
		))
		return
	}

	product, err := h.productService.GetByBarcode(c.Request.Context(), barcode)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Product not found", err.Error(), "PRODUCT_NOT_FOUND",
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Product retrieved successfully", product,
	))
}

// UpdateProduct handles product updates
// @Summary Update product
// @Description Update a product
// @Tags inventory
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body inventory.ProductSparePartUpdateRequest true "Update product request"
// @Success 200 {object} common.APIResponse{data=inventory.ProductSparePart}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /inventory/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Product ID must be a valid integer", "INVALID_ID",
		))
		return
	}

	var req inventory.ProductSparePartUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	product, err := h.productService.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to update product", err.Error(), "UPDATE_FAILED",
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Product updated successfully", product,
	))
}

// DeleteProduct handles product deletion
// @Summary Delete product
// @Description Delete a product
// @Tags inventory
// @Security BearerAuth
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /inventory/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Product ID must be a valid integer", "INVALID_ID",
		))
		return
	}

	err = h.productService.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to delete product", err.Error(), "DELETION_FAILED",
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Product deleted successfully", nil,
	))
}

// GetProducts handles listing products with filtering and pagination
// @Summary List products
// @Description Get products with filtering and pagination
// @Tags inventory
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search query"
// @Param brand_id query int false "Brand ID filter"
// @Param category_id query int false "Category ID filter"
// @Param location_rack query string false "Location filter"
// @Param low_stock query bool false "Low stock filter"
// @Param is_active query bool false "Active status filter"
// @Success 200 {object} common.PaginatedResponse{data=[]inventory.ProductSparePartListItem}
// @Failure 400 {object} common.ErrorResponse
// @Router /inventory/products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	var params inventory.ProductSparePartFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Invalid filter parameters", err.Error(),
		))
		return
	}

	products, total, err := h.productService.List(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve products", err.Error(), "RETRIEVAL_FAILED",
		))
		return
	}

	response := common.NewPaginatedResponse(
		"Products retrieved successfully",
		products,
		total,
		params.Page,
		params.Limit,
	)

	c.JSON(http.StatusOK, response)
}

// GetLowStockProducts handles getting products with low stock
// @Summary Get low stock products
// @Description Get products that are below minimum stock level
// @Tags inventory
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.PaginatedResponse{data=[]inventory.ProductSparePartListItem}
// @Failure 400 {object} common.ErrorResponse
// @Router /inventory/products/low-stock [get]
func (h *ProductHandler) GetLowStockProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, total, err := h.productService.GetLowStockProducts(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve low stock products", err.Error(), "RETRIEVAL_FAILED",
		))
		return
	}

	response := common.NewPaginatedResponse(
		"Low stock products retrieved successfully",
		products,
		total,
		page,
		limit,
	)

	c.JSON(http.StatusOK, response)
}

// SearchProducts handles product search
// @Summary Search products
// @Description Search products by name, code, or description
// @Tags inventory
// @Security BearerAuth
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.PaginatedResponse{data=[]inventory.ProductSparePartListItem}
// @Failure 400 {object} common.ErrorResponse
// @Router /inventory/products/search [get]
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Missing search query", "Search query parameter 'q' is required", "MISSING_QUERY",
		))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, total, err := h.productService.Search(c.Request.Context(), query, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to search products", err.Error(), "SEARCH_FAILED",
		))
		return
	}

	response := common.NewPaginatedResponse(
		"Product search completed successfully",
		products,
		total,
		page,
		limit,
	)

	c.JSON(http.StatusOK, response)
}

// AdjustStock handles stock adjustments
// @Summary Adjust product stock
// @Description Adjust stock quantity for a product
// @Tags inventory
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body struct{NewQuantity int `json:"new_quantity" binding:"required,min=0"`;Reason string `json:"reason" binding:"required"`;Notes *string `json:"notes,omitempty"`} true "Stock adjustment request"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Router /inventory/products/{id}/adjust-stock [post]
func (h *ProductHandler) AdjustStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid ID", "Product ID must be a valid integer", "INVALID_ID",
		))
		return
	}

	var req struct {
		NewQuantity int     `json:"new_quantity" binding:"required,min=0"`
		Reason      string  `json:"reason" binding:"required"`
		Notes       *string `json:"notes,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	processedBy := middleware.GetCurrentUserID(c)
	if processedBy == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Unauthorized", "Invalid authentication", "USER_NOT_AUTHENTICATED",
		))
		return
	}

	err = h.productService.AdjustStock(c.Request.Context(), id, req.NewQuantity, processedBy, req.Reason, req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to adjust stock", err.Error(), "ADJUSTMENT_FAILED",
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Stock adjusted successfully", nil,
	))
}

// GetInventoryValue handles getting total inventory value
// @Summary Get inventory value
// @Description Get total inventory value
// @Tags inventory
// @Security BearerAuth
// @Produce json
// @Success 200 {object} common.APIResponse{data=map[string]float64}
// @Failure 500 {object} common.ErrorResponse
// @Router /inventory/products/inventory-value [get]
func (h *ProductHandler) GetInventoryValue(c *gin.Context) {
	value, err := h.productService.GetInventoryValue(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to calculate inventory value", err.Error(), "CALCULATION_FAILED",
		))
		return
	}

	data := map[string]float64{
		"total_inventory_value": value,
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Inventory value calculated successfully", data,
	))
}