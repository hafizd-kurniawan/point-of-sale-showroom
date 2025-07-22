package auth

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/auth"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/middleware"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/utils"
)

// Handler handles authentication HTTP requests
type Handler struct {
	authService *services.AuthService
}

// NewHandler creates a new authentication handler
func NewHandler(authService *services.AuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}

// Login handles user login
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Login request"
// @Success 200 {object} common.APIResponse{data=auth.LoginResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	// Get client IP and user agent
	ipAddress := utils.GetIPAddress(c.Request)
	userAgent := utils.GetUserAgent(c.Request)

	// Authenticate user
	response, err := h.authService.Login(c.Request.Context(), &req, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Login failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Login successful", response,
	))
}

// Logout handles user logout
// @Summary User logout
// @Description Logout user and invalidate session
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} common.APIResponse{data=auth.LogoutResponse}
// @Failure 401 {object} common.ErrorResponse
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	sessionID := middleware.GetCurrentSessionID(c)
	if sessionID == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Invalid session", "Session ID not found",
		))
		return
	}

	err := h.authService.Logout(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewErrorResponse(
			"Logout failed", err.Error(),
		))
		return
	}

	response := &auth.LogoutResponse{
		Message: "Logout successful",
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Logout successful", response,
	))
}

// Me handles getting current user information
// @Summary Get current user
// @Description Get current authenticated user information
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} common.APIResponse{data=user.User}
// @Failure 401 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /auth/me [get]
func (h *Handler) Me(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Invalid user", "User ID not found",
		))
		return
	}

	user, err := h.authService.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"User not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"User information retrieved successfully", user,
	))
}

// Profile handles getting current user profile with sessions
// @Summary Get user profile
// @Description Get current user profile with recent sessions
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} common.APIResponse{data=auth.ProfileResponse}
// @Failure 401 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /auth/profile [get]
func (h *Handler) Profile(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Invalid user", "User ID not found",
		))
		return
	}

	profile, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewErrorResponse(
			"Profile not found", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Profile retrieved successfully", profile,
	))
}

// ChangePassword handles password change
// @Summary Change password
// @Description Change current user's password
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body auth.ChangePasswordRequest true "Change password request"
// @Success 200 {object} common.APIResponse{data=auth.ChangePasswordResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Router /auth/change-password [post]
func (h *Handler) ChangePassword(c *gin.Context) {
	var req auth.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(
			"Validation failed", "Invalid request data", err.Error(),
		))
		return
	}

	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Invalid user", "User ID not found",
		))
		return
	}

	err := h.authService.ChangePassword(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(
			"Password change failed", err.Error(),
		))
		return
	}

	response := &auth.ChangePasswordResponse{
		Message: "Password changed successfully",
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Password changed successfully", response,
	))
}

// RefreshToken handles token refresh
// @Summary Refresh token
// @Description Refresh JWT token
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} common.APIResponse{data=auth.RefreshTokenResponse}
// @Failure 401 {object} common.ErrorResponse
// @Router /auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	// Get current token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Token required", "Authorization header is required",
		))
		return
	}

	// Extract token
	tokenParts := c.Request.Header.Get("Authorization")
	if len(tokenParts) < 7 || tokenParts[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Invalid token format", "Token must be in Bearer format",
		))
		return
	}
	oldToken := tokenParts[7:]

	// Refresh token
	response, err := h.authService.RefreshToken(c.Request.Context(), oldToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
			"Token refresh failed", err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse(
		"Token refreshed successfully", response,
	))
}

// parseIntParam parses integer parameter from URL
func parseIntParam(c *gin.Context, param string) (int, error) {
	value := c.Param(param)
	return strconv.Atoi(value)
}