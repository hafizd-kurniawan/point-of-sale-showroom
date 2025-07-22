package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/auth"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/utils"
)

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(jwtManager *utils.JWTManager, sessionRepo interfaces.UserSessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
				"Authentication required", "Missing authorization header",
			))
			c.Abort()
			return
		}

		// Check Bearer prefix
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
				"Authentication failed", "Invalid authorization header format",
			))
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Validate JWT token
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			if err == utils.ErrExpiredToken {
				c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
					"Token expired", "Please refresh your token",
				))
			} else {
				c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
					"Authentication failed", "Invalid token",
				))
			}
			c.Abort()
			return
		}

		// Verify session is still active
		session, err := sessionRepo.GetByID(c.Request.Context(), claims.SessionID)
		if err != nil || !session.IsActive {
			c.JSON(http.StatusUnauthorized, common.NewErrorResponse(
				"Session invalid", "Session has been terminated",
			))
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("session_id", claims.SessionID)
		c.Set("claims", claims)

		c.Next()
	}
}

// RequireRole creates a role-based authorization middleware
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, common.NewErrorResponse(
				"Access denied", "User role not found",
			))
			c.Abort()
			return
		}

		roleStr := userRole.(string)
		for _, role := range roles {
			if roleStr == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, common.NewErrorResponse(
			"Access denied", "Insufficient permissions",
		))
		c.Abort()
	}
}

// GetCurrentUser retrieves the current user from context
func GetCurrentUser(c *gin.Context) *auth.TokenClaims {
	if claims, exists := c.Get("claims"); exists {
		return claims.(*auth.TokenClaims)
	}
	return nil
}

// GetCurrentUserID retrieves the current user ID from context
func GetCurrentUserID(c *gin.Context) int {
	if userID, exists := c.Get("user_id"); exists {
		return userID.(int)
	}
	return 0
}

// GetCurrentSessionID retrieves the current session ID from context
func GetCurrentSessionID(c *gin.Context) int {
	if sessionID, exists := c.Get("session_id"); exists {
		return sessionID.(int)
	}
	return 0
}