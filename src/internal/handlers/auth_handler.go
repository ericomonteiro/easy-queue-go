package handlers

import (
	"easy-queue-go/src/internal/log"
	"easy-queue-go/src/internal/models"
	"easy-queue-go/src/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles user login requests
// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn(ctx, "Invalid login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	response, err := h.authService.Login(ctx, &req)
	if err != nil {
		log.Error(ctx, "Login failed", zap.Error(err))
		
		// Return 401 for authentication failures
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Info(ctx, "User logged in successfully",
		zap.String("user_id", response.User.ID.String()),
		zap.String("email", response.User.Email),
	)

	c.JSON(http.StatusOK, response)
}

// RefreshToken handles token refresh requests
// @Summary Refresh access token
// @Description Generate new access and refresh tokens using a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} models.RefreshTokenResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn(ctx, "Invalid refresh token request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	response, err := h.authService.RefreshToken(ctx, &req)
	if err != nil {
		log.Error(ctx, "Token refresh failed", zap.Error(err))
		
		// Return 401 for invalid refresh tokens
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Info(ctx, "Tokens refreshed successfully")

	c.JSON(http.StatusOK, response)
}
