package services

import (
	"context"
	"easy-queue-go/src/internal/log"
	"easy-queue-go/src/internal/models"
	"easy-queue-go/src/internal/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var tracer = otel.Tracer("user-service")

// UserService defines the interface for user business operations
type UserService interface {
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserResponse, error)
	ListAllUsers(ctx context.Context) ([]*models.UserResponse, error)
}

// userService implements UserService
type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user with business validations
func (s *userService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error) {
	ctx, span := tracer.Start(ctx, "UserService.CreateUser",
		trace.WithAttributes(
			attribute.String("email", req.Email),
		),
	)
	defer span.End()

	log.Info(ctx, "Creating new user",
		zap.String("email", req.Email),
		zap.Int("roles_count", len(req.Roles)),
	)

	// Validate if email already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		log.Warn(ctx, "User with email already exists",
			zap.String("email", req.Email),
		)
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash the password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		log.Error(ctx, "Failed to hash password", zap.Error(err))
		span.RecordError(err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create the user
	now := time.Now()
	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Phone:        req.Phone,
		Roles:        req.Roles,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Save to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		log.Error(ctx, "Failed to create user in database",
			zap.Error(err),
			zap.String("email", req.Email),
		)
		span.RecordError(err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	log.Info(ctx, "User created successfully",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email),
	)

	span.SetAttributes(attribute.String("user_id", user.ID.String()))

	return user.ToResponse(), nil
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserResponse, error) {
	ctx, span := tracer.Start(ctx, "UserService.GetUserByID",
		trace.WithAttributes(
			attribute.String("user_id", id.String()),
		),
	)
	defer span.End()

	log.Info(ctx, "Getting user by ID", zap.String("user_id", id.String()))

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		log.Error(ctx, "Failed to get user by ID",
			zap.Error(err),
			zap.String("user_id", id.String()),
		)
		span.RecordError(err)
		return nil, err
	}

	return user.ToResponse(), nil
}

// GetUserByEmail retrieves a user by email
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.UserResponse, error) {
	ctx, span := tracer.Start(ctx, "UserService.GetUserByEmail",
		trace.WithAttributes(
			attribute.String("email", email),
		),
	)
	defer span.End()

	log.Info(ctx, "Getting user by email", zap.String("email", email))

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		log.Error(ctx, "Failed to get user by email",
			zap.Error(err),
			zap.String("email", email),
		)
		span.RecordError(err)
		return nil, err
	}

	return user.ToResponse(), nil
}

// ListAllUsers returns all users in the system
func (s *userService) ListAllUsers(ctx context.Context) ([]*models.UserResponse, error) {
	ctx, span := tracer.Start(ctx, "UserService.ListAllUsers")
	defer span.End()

	log.Info(ctx, "Listing all users")

	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		log.Error(ctx, "Failed to list all users", zap.Error(err))
		span.RecordError(err)
		return nil, err
	}

	// Convert to response format
	responses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	log.Info(ctx, "Successfully listed all users", zap.Int("count", len(responses)))
	span.SetAttributes(attribute.Int("user_count", len(responses)))

	return responses, nil
}

// hashPassword generates a bcrypt hash of the password
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword verifies if the password matches the hash
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
