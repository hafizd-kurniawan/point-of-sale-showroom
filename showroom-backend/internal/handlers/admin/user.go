package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/middleware"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/user"
	commonModels "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/models/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services"
)

// Handler handles admin HTTP requests
type Handler struct {
	userService *services.UserService
}

// NewHandler creates a new admin handler
func NewHandler(userService *services.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

// CreateUser handles user creation
// @Summary Create user
// @Description Create a new user account
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body user.UserCreateRequest true "Create user request"
// @Success 201 {object} common.APIResponse{data=user.User}
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Router /admin/users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var req user.UserCreateRequest
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

	newUser, err := h.userService.CreateUser(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"User creation failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, common.NewSuccessResponse(
		"User created successfully", newUser,
	))
}

// GetUsers handles user list with filtering and pagination
// @Summary List users
// @Description Get paginated list of users with filtering
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param role query string false "Filter by role"
// @Param is_active query bool false "Filter by active status"
// @Param search query string false "Search by username, email, full name, or phone"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Router /admin/users [get]
func (h *Handler) GetUsers(c *gin.Context) {
	var params user.UserFilterParams

	// Bind query parameters
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Invalid query parameters", "Failed to parse query parameters", err.Error(),
		))
		return
	}

	// Handle role parameter
	if roleStr := c.Query("role"); roleStr != "" {
		role := commonModels.UserRole(roleStr)
		if !role.IsValid() {
			c.JSON(http.StatusBadRequest, common.NewErrorResponse(
				"Invalid role", "Role must be one of: admin, sales, cashier, mechanic, manager",
			))
			return
		}
		params.Role = &role
	}

	// Handle is_active parameter
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			params.IsActive = &isActive
		}
	}

	result, err := h.userService.ListUsers(c.Request.Context(), &params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve users", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Users retrieved successfully", result,
	))
}

// GetUser handles getting a single user by ID
// @Summary Get user by ID
// @Description Get user details by ID
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} common.APIResponse{data=user.User}
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /admin/users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid user ID", "User ID must be a valid integer",
		))
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"User not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"User retrieved successfully", user,
	))
}

// UpdateUser handles user update
// @Summary Update user
// @Description Update user information
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body user.UserUpdateRequest true "Update user request"
// @Success 200 {object} common.APIResponse{data=user.User}
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /admin/users/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid user ID", "User ID must be a valid integer",
		))
		return
	}

	var req user.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	updatedUser, err := h.userService.UpdateUser(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"User update failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"User updated successfully", updatedUser,
	))
}

// DeleteUser handles user deletion (soft delete)
// @Summary Delete user
// @Description Soft delete user account
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} common.APIResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /admin/users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid user ID", "User ID must be a valid integer",
		))
		return
	}

	err = h.userService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"User deletion failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"User deleted successfully", nil,
	))
}

// GetUsersByRole handles getting users by role
// @Summary Get users by role
// @Description Get users filtered by specific role
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param role path string true "User role"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Router /admin/users/role/{role} [get]
func (h *Handler) GetUsersByRole(c *gin.Context) {
	roleStr := c.Param("role")
	role := commonModels.UserRole(roleStr)
	if !role.IsValid() {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid role", "Role must be one of: admin, sales, cashier, mechanic, manager",
		))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.userService.GetUsersByRole(c.Request.Context(), roleStr, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve users", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Users retrieved successfully", result,
	))
}

// GetUserSessions handles getting user sessions
// @Summary Get user sessions
// @Description Get user's active and recent sessions
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} common.APIResponse{data=common.PaginatedResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Router /admin/users/{id}/sessions [get]
func (h *Handler) GetUserSessions(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid user ID", "User ID must be a valid integer",
		))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.userService.GetUserSessions(c.Request.Context(), id, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to retrieve user sessions", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"User sessions retrieved successfully", result,
	))
}

// RevokeUserSessions handles revoking all user sessions
// @Summary Revoke user sessions
// @Description Revoke all active sessions for a user
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} common.APIResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Router /admin/users/{id}/sessions [delete]
func (h *Handler) RevokeUserSessions(c *gin.Context) {
	id, err := parseIntParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Invalid user ID", "User ID must be a valid integer",
		))
		return
	}

	err = h.userService.RevokeUserSessions(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Failed to revoke user sessions", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"User sessions revoked successfully", nil,
	))
}

// parseIntParam parses integer parameter from URL
func parseIntParam(c *gin.Context, param string) (int, error) {
	value := c.Param(param)
	return strconv.Atoi(value)
}