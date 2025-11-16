package services

import (
	"context"
	"easy-queue-go/src/internal/log"
	"easy-queue-go/src/internal/models"
	"easy-queue-go/src/internal/repositories"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var businessTracer = otel.Tracer("business-service")

// BusinessService defines the interface for business operations
type BusinessService interface {
	CreateBusiness(ctx context.Context, ownerID uuid.UUID, req *models.CreateBusinessRequest) (*models.BusinessResponse, error)
	GetBusinessByID(ctx context.Context, id uuid.UUID) (*models.BusinessResponse, error)
	GetBusinessesByOwner(ctx context.Context, ownerID uuid.UUID) ([]*models.BusinessResponse, error)
	ListAllBusinesses(ctx context.Context) ([]*models.BusinessResponse, error)
	UpdateBusiness(ctx context.Context, id uuid.UUID, ownerID uuid.UUID, req *models.UpdateBusinessRequest) (*models.BusinessResponse, error)
	DeleteBusiness(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) error
}

// businessService implements BusinessService
type businessService struct {
	businessRepo repositories.BusinessRepository
	userRepo     repositories.UserRepository
}

// NewBusinessService creates a new instance of BusinessService
func NewBusinessService(businessRepo repositories.BusinessRepository, userRepo repositories.UserRepository) BusinessService {
	return &businessService{
		businessRepo: businessRepo,
		userRepo:     userRepo,
	}
}

// CreateBusiness creates a new business with validation
func (s *businessService) CreateBusiness(ctx context.Context, ownerID uuid.UUID, req *models.CreateBusinessRequest) (*models.BusinessResponse, error) {
	ctx, span := businessTracer.Start(ctx, "BusinessService.CreateBusiness",
		trace.WithAttributes(
			attribute.String("owner_id", ownerID.String()),
			attribute.String("business_name", req.Name),
		),
	)
	defer span.End()

	log.Info(ctx, "Creating new business",
		zap.String("owner_id", ownerID.String()),
		zap.String("name", req.Name),
	)

	// Verify that the owner exists and has the Business Owner role
	owner, err := s.userRepo.FindByID(ctx, ownerID)
	if err != nil {
		log.Error(ctx, "Failed to find owner user",
			zap.Error(err),
			zap.String("owner_id", ownerID.String()),
		)
		span.RecordError(err)
		return nil, fmt.Errorf("owner user not found")
	}

	if !owner.HasRole(models.RoleBusinessOwner) {
		log.Warn(ctx, "User does not have Business Owner role",
			zap.String("owner_id", ownerID.String()),
			zap.Any("roles", owner.Roles),
		)
		return nil, fmt.Errorf("user must have Business Owner role to create a business")
	}

	// Create the business
	business := req.ToBusiness(ownerID)

	// Save to database
	if err := s.businessRepo.Create(ctx, business); err != nil {
		log.Error(ctx, "Failed to create business in database",
			zap.Error(err),
			zap.String("name", req.Name),
		)
		span.RecordError(err)
		return nil, fmt.Errorf("failed to create business: %w", err)
	}

	log.Info(ctx, "Business created successfully",
		zap.String("business_id", business.ID.String()),
		zap.String("name", business.Name),
	)

	span.SetAttributes(attribute.String("business_id", business.ID.String()))

	return business.ToResponse(), nil
}

// GetBusinessByID retrieves a business by ID
func (s *businessService) GetBusinessByID(ctx context.Context, id uuid.UUID) (*models.BusinessResponse, error) {
	ctx, span := businessTracer.Start(ctx, "BusinessService.GetBusinessByID",
		trace.WithAttributes(
			attribute.String("business_id", id.String()),
		),
	)
	defer span.End()

	log.Info(ctx, "Getting business by ID", zap.String("business_id", id.String()))

	business, err := s.businessRepo.FindByID(ctx, id)
	if err != nil {
		log.Error(ctx, "Failed to get business by ID",
			zap.Error(err),
			zap.String("business_id", id.String()),
		)
		span.RecordError(err)
		return nil, err
	}

	return business.ToResponse(), nil
}

// GetBusinessesByOwner retrieves all businesses owned by a specific user
func (s *businessService) GetBusinessesByOwner(ctx context.Context, ownerID uuid.UUID) ([]*models.BusinessResponse, error) {
	ctx, span := businessTracer.Start(ctx, "BusinessService.GetBusinessesByOwner",
		trace.WithAttributes(
			attribute.String("owner_id", ownerID.String()),
		),
	)
	defer span.End()

	log.Info(ctx, "Getting businesses by owner", zap.String("owner_id", ownerID.String()))

	businesses, err := s.businessRepo.FindByOwnerID(ctx, ownerID)
	if err != nil {
		log.Error(ctx, "Failed to get businesses by owner",
			zap.Error(err),
			zap.String("owner_id", ownerID.String()),
		)
		span.RecordError(err)
		return nil, err
	}

	// Convert to response format
	responses := make([]*models.BusinessResponse, len(businesses))
	for i, business := range businesses {
		responses[i] = business.ToResponse()
	}

	log.Info(ctx, "Successfully retrieved businesses by owner",
		zap.Int("count", len(responses)),
		zap.String("owner_id", ownerID.String()),
	)
	span.SetAttributes(attribute.Int("business_count", len(responses)))

	return responses, nil
}

// ListAllBusinesses returns all businesses in the system
func (s *businessService) ListAllBusinesses(ctx context.Context) ([]*models.BusinessResponse, error) {
	ctx, span := businessTracer.Start(ctx, "BusinessService.ListAllBusinesses")
	defer span.End()

	log.Info(ctx, "Listing all businesses")

	businesses, err := s.businessRepo.FindAll(ctx)
	if err != nil {
		log.Error(ctx, "Failed to list all businesses", zap.Error(err))
		span.RecordError(err)
		return nil, err
	}

	// Convert to response format
	responses := make([]*models.BusinessResponse, len(businesses))
	for i, business := range businesses {
		responses[i] = business.ToResponse()
	}

	log.Info(ctx, "Successfully listed all businesses", zap.Int("count", len(responses)))
	span.SetAttributes(attribute.Int("business_count", len(responses)))

	return responses, nil
}

// UpdateBusiness updates an existing business
func (s *businessService) UpdateBusiness(ctx context.Context, id uuid.UUID, ownerID uuid.UUID, req *models.UpdateBusinessRequest) (*models.BusinessResponse, error) {
	ctx, span := businessTracer.Start(ctx, "BusinessService.UpdateBusiness",
		trace.WithAttributes(
			attribute.String("business_id", id.String()),
			attribute.String("owner_id", ownerID.String()),
		),
	)
	defer span.End()

	log.Info(ctx, "Updating business",
		zap.String("business_id", id.String()),
		zap.String("owner_id", ownerID.String()),
	)

	// Get existing business
	business, err := s.businessRepo.FindByID(ctx, id)
	if err != nil {
		log.Error(ctx, "Failed to find business",
			zap.Error(err),
			zap.String("business_id", id.String()),
		)
		span.RecordError(err)
		return nil, err
	}

	// Verify ownership
	if business.OwnerID != ownerID {
		log.Warn(ctx, "User is not the owner of this business",
			zap.String("business_id", id.String()),
			zap.String("owner_id", ownerID.String()),
			zap.String("actual_owner_id", business.OwnerID.String()),
		)
		return nil, fmt.Errorf("you are not authorized to update this business")
	}

	business.ApplyUpdateRequest(req)

	// Save to database
	if err := s.businessRepo.Update(ctx, business); err != nil {
		log.Error(ctx, "Failed to update business in database",
			zap.Error(err),
			zap.String("business_id", id.String()),
		)
		span.RecordError(err)
		return nil, fmt.Errorf("failed to update business: %w", err)
	}

	log.Info(ctx, "Business updated successfully",
		zap.String("business_id", business.ID.String()),
	)

	return business.ToResponse(), nil
}

// DeleteBusiness deletes a business
func (s *businessService) DeleteBusiness(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) error {
	ctx, span := businessTracer.Start(ctx, "BusinessService.DeleteBusiness",
		trace.WithAttributes(
			attribute.String("business_id", id.String()),
			attribute.String("owner_id", ownerID.String()),
		),
	)
	defer span.End()

	log.Info(ctx, "Deleting business",
		zap.String("business_id", id.String()),
		zap.String("owner_id", ownerID.String()),
	)

	// Get existing business
	business, err := s.businessRepo.FindByID(ctx, id)
	if err != nil {
		log.Error(ctx, "Failed to find business",
			zap.Error(err),
			zap.String("business_id", id.String()),
		)
		span.RecordError(err)
		return err
	}

	// Verify ownership
	if business.OwnerID != ownerID {
		log.Warn(ctx, "User is not the owner of this business",
			zap.String("business_id", id.String()),
			zap.String("owner_id", ownerID.String()),
			zap.String("actual_owner_id", business.OwnerID.String()),
		)
		return fmt.Errorf("you are not authorized to delete this business")
	}

	// Delete from database
	if err := s.businessRepo.Delete(ctx, id); err != nil {
		log.Error(ctx, "Failed to delete business from database",
			zap.Error(err),
			zap.String("business_id", id.String()),
		)
		span.RecordError(err)
		return fmt.Errorf("failed to delete business: %w", err)
	}

	log.Info(ctx, "Business deleted successfully",
		zap.String("business_id", id.String()),
	)

	return nil
}
