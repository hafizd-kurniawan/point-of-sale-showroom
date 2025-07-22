package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/config"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/dto/common"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/handlers/admin"
	authHandlers "github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/handlers/auth"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/repositories/interfaces"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/routes"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/services"
	"github.com/hafizd-kurniawan/point-of-sale-showroom/showroom-backend/internal/utils"
)

func TestHealthCheck(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response common.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response.Status)
}

// Note: Full integration tests with database would require:
// - Test database setup
// - Database migrations
// - Seed test data
// - Proper cleanup after tests

func setupTestRouter() http.Handler {
	// Use in-memory or test database configuration
	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "Test App",
			Version: "1.0.0",
		},
		JWT: config.JWTConfig{
			SecretKey:      "test-secret-key",
			ExpirationHour: 24,
		},
		Server: config.ServerConfig{
			Env: "test",
		},
	}

	// For integration tests, you might want to use a test database
	// For now, we'll test endpoints that don't require database
	
	// Mock repositories for testing (nil for basic endpoint tests)
	var userRepo interfaces.UserRepository
	var sessionRepo interfaces.UserSessionRepository

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(cfg.JWT.SecretKey, cfg.JWT.GetExpiration())

	// Initialize services
	authService := services.NewAuthService(userRepo, sessionRepo, jwtManager)
	userService := services.NewUserService(userRepo, sessionRepo)

	// Initialize handlers
	authHandler := authHandlers.NewHandler(authService)
	adminHandler := admin.NewHandler(userService)

	// Initialize router
	router := routes.NewRouter(authHandler, adminHandler, jwtManager, sessionRepo, cfg)

	return router.SetupRoutes()
}

// Note: For full integration tests with database, you would:
// 1. Set up a test database
// 2. Run migrations
// 3. Seed test data
// 4. Run tests
// 5. Clean up database