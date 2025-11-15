package handlers

import (
	"easy-queue-go/src/internal/log"
	"easy-queue-go/src/internal/models"
	"easy-queue-go/src/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// BusinessHandler manages HTTP requests related to businesses
type BusinessHandler struct {
	businessService services.BusinessService
}

// NewBusinessHandler creates a new instance of BusinessHandler
func NewBusinessHandler(businessService services.BusinessService) *BusinessHandler {
	return &BusinessHandler{
		businessService: businessService,
	}
}

// CreateBusiness godoc
// @Summary Creates a new business
// @Description Creates a new business for the authenticated business owner
// @Tags businesses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param business body models.CreateBusinessRequest true "Business data"
// @Success 201 {object} models.BusinessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /businesses [post]
func (h *BusinessHandler) CreateBusiness(c *gin.Context) {
	ctx := c.Request.Context()

	// Get JWT claims from context
	jwtClaims, ok := GetClaimsFromContext(c)
	if !ok {
		return
	}

	var req models.CreateBusinessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn(ctx, "Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	business, err := h.businessService.CreateBusiness(ctx, jwtClaims.UserID, &req)
	if err != nil {
		log.Error(ctx, "Failed to create business", zap.Error(err))

		if err.Error() == "user must have Business Owner role to create a business" {
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "forbidden",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create business",
		})
		return
	}

	log.Info(ctx, "Business created successfully via HTTP",
		zap.String("business_id", business.ID.String()),
		zap.String("name", business.Name),
	)

	c.JSON(http.StatusCreated, business)
}

// GetBusinessByID godoc
// @Summary Retrieves a business by ID
// @Description Returns the data of a specific business by ID
// @Tags businesses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Business ID (UUID)"
// @Success 200 {object} models.BusinessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /businesses/{id} [get]
func (h *BusinessHandler) GetBusinessByID(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Warn(ctx, "Invalid business ID format", zap.String("id", idParam))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid UUID format",
		})
		return
	}

	business, err := h.businessService.GetBusinessByID(ctx, id)
	if err != nil {
		log.Error(ctx, "Failed to get business", zap.Error(err), zap.String("id", id.String()))

		if err.Error() == "business not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "business_not_found",
				Message: "Business not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get business",
		})
		return
	}

	c.JSON(http.StatusOK, business)
}

// GetMyBusinesses godoc
// @Summary Retrieves all businesses owned by the authenticated user
// @Description Returns a list of all businesses owned by the current user
// @Tags businesses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.BusinessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /businesses/my [get]
func (h *BusinessHandler) GetMyBusinesses(c *gin.Context) {
	ctx := c.Request.Context()

	// Get JWT claims from context
	jwtClaims, ok := GetClaimsFromContext(c)
	if !ok {
		return
	}

	businesses, err := h.businessService.GetBusinessesByOwner(ctx, jwtClaims.UserID)
	if err != nil {
		log.Error(ctx, "Failed to get businesses", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get businesses",
		})
		return
	}

	log.Info(ctx, "Successfully retrieved user's businesses",
		zap.String("user_id", jwtClaims.UserID.String()),
		zap.Int("count", len(businesses)),
	)

	c.JSON(http.StatusOK, businesses)
}

// ListAllBusinesses godoc
// @Summary Lists all businesses (Admin only)
// @Description Returns a list of all businesses in the system
// @Tags businesses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.BusinessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/businesses [get]
func (h *BusinessHandler) ListAllBusinesses(c *gin.Context) {
	ctx := c.Request.Context()

	log.Info(ctx, "Listing all businesses")

	businesses, err := h.businessService.ListAllBusinesses(ctx)
	if err != nil {
		log.Error(ctx, "Failed to list all businesses", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to list businesses",
		})
		return
	}

	log.Info(ctx, "Successfully listed all businesses", zap.Int("count", len(businesses)))
	c.JSON(http.StatusOK, businesses)
}

// UpdateBusiness godoc
// @Summary Updates a business
// @Description Updates an existing business (owner only)
// @Tags businesses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Business ID (UUID)"
// @Param business body models.UpdateBusinessRequest true "Business data to update"
// @Success 200 {object} models.BusinessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /businesses/{id} [put]
func (h *BusinessHandler) UpdateBusiness(c *gin.Context) {
	ctx := c.Request.Context()

	// Get JWT claims from context
	jwtClaims, ok := GetClaimsFromContext(c)
	if !ok {
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Warn(ctx, "Invalid business ID format", zap.String("id", idParam))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid UUID format",
		})
		return
	}

	var req models.UpdateBusinessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn(ctx, "Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	business, err := h.businessService.UpdateBusiness(ctx, id, jwtClaims.UserID, &req)
	if err != nil {
		log.Error(ctx, "Failed to update business", zap.Error(err))

		if err.Error() == "business not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "business_not_found",
				Message: "Business not found",
			})
			return
		}

		if err.Error() == "you are not authorized to update this business" {
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "forbidden",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to update business",
		})
		return
	}

	log.Info(ctx, "Business updated successfully via HTTP",
		zap.String("business_id", business.ID.String()),
	)

	c.JSON(http.StatusOK, business)
}

// DeleteBusiness godoc
// @Summary Deletes a business
// @Description Deletes an existing business (owner only)
// @Tags businesses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Business ID (UUID)"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /businesses/{id} [delete]
func (h *BusinessHandler) DeleteBusiness(c *gin.Context) {
	ctx := c.Request.Context()

	// Get JWT claims from context
	jwtClaims, ok := GetClaimsFromContext(c)
	if !ok {
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Warn(ctx, "Invalid business ID format", zap.String("id", idParam))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid UUID format",
		})
		return
	}

	err = h.businessService.DeleteBusiness(ctx, id, jwtClaims.UserID)
	if err != nil {
		log.Error(ctx, "Failed to delete business", zap.Error(err))

		if err.Error() == "business not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "business_not_found",
				Message: "Business not found",
			})
			return
		}

		if err.Error() == "you are not authorized to delete this business" {
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "forbidden",
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to delete business",
		})
		return
	}

	log.Info(ctx, "Business deleted successfully via HTTP",
		zap.String("business_id", id.String()),
	)

	c.Status(http.StatusNoContent)
}
