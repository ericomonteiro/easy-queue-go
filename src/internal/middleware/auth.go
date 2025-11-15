package middleware

import (
	"easy-queue-go/src/internal/log"
	"easy-queue-go/src/internal/models"
	"easy-queue-go/src/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// UserClaimsKey is the context key for storing JWT claims
	UserClaimsKey ContextKey = "user_claims"
)

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Warn(ctx, "Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			c.Abort()
			return
		}

		// Check Bearer token format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Warn(ctx, "Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := authService.ValidateToken(ctx, tokenString, models.TokenTypeAccess)
		if err != nil {
			log.Warn(ctx, "Invalid or expired token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		// Store claims in context for use in handlers
		c.Set(string(UserClaimsKey), claims)

		log.Info(ctx, "User authenticated",
			zap.String("user_id", claims.UserID.String()),
			zap.String("email", claims.Email),
			zap.String("role", string(claims.Role)),
		)

		c.Next()
	}
}

// RequireRole creates a middleware that checks if the user has one of the required roles
func RequireRole(roles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get claims from context
		claimsValue, exists := c.Get(string(UserClaimsKey))
		if !exists {
			log.Error(ctx, "User claims not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		claims, ok := claimsValue.(*models.JWTClaims)
		if !ok {
			log.Error(ctx, "Invalid user claims type in context")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			c.Abort()
			return
		}

		// Check if user has one of the required roles
		hasRole := false
		for _, role := range roles {
			if claims.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			log.Warn(ctx, "User does not have required role",
				zap.String("user_id", claims.UserID.String()),
				zap.String("user_role", string(claims.Role)),
			)
			c.JSON(http.StatusForbidden, gin.H{
				"error": "insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserClaims is a helper function to extract JWT claims from gin context
func GetUserClaims(c *gin.Context) (*models.JWTClaims, bool) {
	claimsValue, exists := c.Get(string(UserClaimsKey))
	if !exists {
		return nil, false
	}

	claims, ok := claimsValue.(*models.JWTClaims)
	return claims, ok
}
