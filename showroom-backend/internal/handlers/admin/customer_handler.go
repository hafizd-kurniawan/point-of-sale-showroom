package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/middleware"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/master"
	masterService "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services/master"
)

// CustomerHandler handles customer HTTP requests
type CustomerHandler struct {
	customerService *masterService.CustomerService
}

// NewCustomerHandler creates a new customer handler
func NewCustomerHandler(customerService *masterService.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}

// CreateCustomer handles customer creation
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
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
			"Invalid user", "Creator user ID not found",
		))
		return
	}

	customer, err := h.customerService.CreateCustomer(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Customer creation failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"Customer created successfully", customer,
	))
}

// GetCustomers handles customer list with filtering and pagination
func (h *CustomerHandler) GetCustomers(c *gin.Context) {
	var params master.CustomerFilterParams

	// Bind query parameters
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Failed to parse query parameters", err.Error(),
		))
		return
	}

	// Handle customer_type parameter
	if customerTypeStr := c.Query("customer_type"); customerTypeStr != "" {
		customerType := master.CustomerType(customerTypeStr)
		if !customerType.IsValid() {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(
				"Invalid customer type", "Customer type must be one of: individual, corporate",
			))
			return
		}
		params.CustomerType = &customerType
	}

	// Handle is_active parameter
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			params.IsActive = &isActive
		}
	}

	result, err := h.customerService.ListCustomers(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve customers", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Customers retrieved successfully", result,
	))
}

// GetCustomer handles getting a single customer by ID
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid customer ID", "Customer ID must be a valid integer",
		))
		return
	}

	customer, err := h.customerService.GetCustomer(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Customer not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Customer retrieved successfully", customer,
	))
}

// UpdateCustomer handles customer update
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid customer ID", "Customer ID must be a valid integer",
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

	updatedCustomer, err := h.customerService.UpdateCustomer(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Customer update failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Customer updated successfully", updatedCustomer,
	))
}

// DeleteCustomer handles customer deletion (soft delete)
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid customer ID", "Customer ID must be a valid integer",
		))
		return
	}

	err = h.customerService.DeleteCustomer(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Customer deletion failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Customer deleted successfully", nil,
	))
}