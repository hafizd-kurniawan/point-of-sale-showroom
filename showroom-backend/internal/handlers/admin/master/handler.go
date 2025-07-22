package master

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	masterDto "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/master"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/middleware"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
	masterServices "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services/master"
)

// Handler handles master data HTTP requests
type Handler struct {
	customerService     *masterServices.CustomerService
	supplierService     *masterServices.SupplierService
	vehicleBrandService *masterServices.VehicleBrandService
}

// NewHandler creates a new master data handler
func NewHandler(
	customerService *masterServices.CustomerService,
	supplierService *masterServices.SupplierService,
	vehicleBrandService *masterServices.VehicleBrandService,
) *Handler {
	return &Handler{
		customerService:     customerService,
		supplierService:     supplierService,
		vehicleBrandService: vehicleBrandService,
	}
}

// Customer endpoints

// CreateCustomer handles customer creation
func (h *Handler) CreateCustomer(c *gin.Context) {
	var req master.CustomerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	createdBy := middleware.GetCurrentUserID(c)
	if createdBy == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Unauthorized", "User not found in context",
		))
		return
	}

	customer, err := h.customerService.Create(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to create customer", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Customer created successfully", customer,
	))
}

// GetCustomers handles customer list retrieval with filtering
func (h *Handler) GetCustomers(c *gin.Context) {
	var params master.CustomerFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid query parameters", err.Error(),
		))
		return
	}

	customers, meta, err := h.customerService.List(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve customers", err.Error(),
		))
		return
	}

	response := masterDto.CustomerListResponse{
		Customers:      customers,
		PaginationMeta: *meta,
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Customers retrieved successfully", response,
	))
}

// GetCustomer handles customer retrieval by ID
func (h *Handler) GetCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid customer ID", "Customer ID must be a valid number",
		))
		return
	}

	customer, err := h.customerService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "customer not found" {
			c.JSON(http.StatusNotFound, common.NewErrorResponse(
				"Customer not found", err.Error(),
			))
			return
		}
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve customer", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Customer retrieved successfully", customer,
	))
}

// UpdateCustomer handles customer updates
func (h *Handler) UpdateCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid customer ID", "Customer ID must be a valid number",
		))
		return
	}

	var req master.CustomerUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	customer, err := h.customerService.Update(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "customer not found" {
			c.JSON(http.StatusNotFound, common.NewErrorResponse(
				"Customer not found", err.Error(),
			))
			return
		}
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to update customer", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Customer updated successfully", customer,
	))
}

// DeleteCustomer handles customer deletion
func (h *Handler) DeleteCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid customer ID", "Customer ID must be a valid number",
		))
		return
	}

	err = h.customerService.Delete(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "customer not found" {
			c.JSON(http.StatusNotFound, common.NewErrorResponse(
				"Customer not found", err.Error(),
			))
			return
		}
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to delete customer", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Customer deleted successfully", nil,
	))
}

// Supplier endpoints

// CreateSupplier handles supplier creation
func (h *Handler) CreateSupplier(c *gin.Context) {
	var req master.SupplierCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	createdBy := middleware.GetCurrentUserID(c)
	if createdBy == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Unauthorized", "User not found in context",
		))
		return
	}

	supplier, err := h.supplierService.Create(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to create supplier", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Supplier created successfully", supplier,
	))
}

// GetSuppliers handles supplier list retrieval with filtering
func (h *Handler) GetSuppliers(c *gin.Context) {
	var params master.SupplierFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid query parameters", err.Error(),
		))
		return
	}

	suppliers, meta, err := h.supplierService.List(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve suppliers", err.Error(),
		))
		return
	}

	response := masterDto.SupplierListResponse{
		Suppliers:      suppliers,
		PaginationMeta: *meta,
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Suppliers retrieved successfully", response,
	))
}

// GetSupplier handles supplier retrieval by ID
func (h *Handler) GetSupplier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid supplier ID", "Supplier ID must be a valid number",
		))
		return
	}

	supplier, err := h.supplierService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "supplier not found" {
			c.JSON(http.StatusNotFound, common.NewErrorResponse(
				"Supplier not found", err.Error(),
			))
			return
		}
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve supplier", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Supplier retrieved successfully", supplier,
	))
}

// UpdateSupplier handles supplier updates
func (h *Handler) UpdateSupplier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid supplier ID", "Supplier ID must be a valid number",
		))
		return
	}

	var req master.SupplierUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	supplier, err := h.supplierService.Update(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "supplier not found" {
			c.JSON(http.StatusNotFound, common.NewErrorResponse(
				"Supplier not found", err.Error(),
			))
			return
		}
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to update supplier", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Supplier updated successfully", supplier,
	))
}

// DeleteSupplier handles supplier deletion
func (h *Handler) DeleteSupplier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid supplier ID", "Supplier ID must be a valid number",
		))
		return
	}

	err = h.supplierService.Delete(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "supplier not found" {
			c.JSON(http.StatusNotFound, common.NewErrorResponse(
				"Supplier not found", err.Error(),
			))
			return
		}
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to delete supplier", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Supplier deleted successfully", nil,
	))
}

// Vehicle Brand endpoints

// CreateVehicleBrand handles vehicle brand creation
func (h *Handler) CreateVehicleBrand(c *gin.Context) {
	var req master.VehicleBrandCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	createdBy := middleware.GetCurrentUserID(c)
	if createdBy == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Unauthorized", "User not found in context",
		))
		return
	}

	brand, err := h.vehicleBrandService.Create(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Failed to create vehicle brand", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Vehicle brand created successfully", brand,
	))
}

// GetVehicleBrands handles vehicle brand list retrieval
func (h *Handler) GetVehicleBrands(c *gin.Context) {
	brands, err := h.vehicleBrandService.ListActive(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve vehicle brands", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Vehicle brands retrieved successfully", brands,
	))
}