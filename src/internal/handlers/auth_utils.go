package handlers

import (
	"easy-queue-go/src/internal/log"
	"easy-queue-go/src/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetClaimsFromContext extracts JWT claims from the Gin context
// Returns the claims and true if successful, or responds with error and returns nil, false
func GetClaimsFromContext(c *gin.Context) (*models.JWTClaims, bool) {
	ctx := c.Request.Context()

	claims, exists := c.Get("user_claims")
	if !exists {
		log.Warn(ctx, "Claims not found in context")
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "Authentication required",
		})
		return nil, false
	}

	jwtClaims, ok := claims.(*models.JWTClaims)
	if !ok {
		log.Error(ctx, "Failed to cast claims to JWTClaims")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Invalid authentication data",
		})
		return nil, false
	}

	return jwtClaims, true
}

// RequireBusinessOwnerRole checks if the user has Business Owner role
// Returns true if the user has the role, or responds with error and returns false
func RequireBusinessOwnerRole(c *gin.Context, claims *models.JWTClaims) bool {
	ctx := c.Request.Context()

	if !claims.HasRole(models.RoleBusinessOwner) {
		log.Warn(ctx, "User does not have Business Owner role")
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "forbidden",
			Message: "Business Owner role required",
		})
		return false
	}

	return true
}
