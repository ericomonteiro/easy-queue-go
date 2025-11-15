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

// UserService define a interface para operações de negócio de usuário
type UserService interface {
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserResponse, error)
}

// userService implementa UserService
type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService cria uma nova instância de UserService
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser cria um novo usuário com validações de negócio
func (s *userService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error) {
	ctx, span := tracer.Start(ctx, "UserService.CreateUser",
		trace.WithAttributes(
			attribute.String("email", req.Email),
			attribute.String("role", string(req.Role)),
		),
	)
	defer span.End()

	log.Info(ctx, "Creating new user",
		zap.String("email", req.Email),
		zap.String("role", string(req.Role)),
	)

	// Validar se o email já existe
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		log.Warn(ctx, "User with email already exists",
			zap.String("email", req.Email),
		)
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash da senha
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		log.Error(ctx, "Failed to hash password", zap.Error(err))
		span.RecordError(err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Criar o usuário
	now := time.Now()
	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Phone:        req.Phone,
		Role:         req.Role,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Salvar no banco
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

// GetUserByID busca um usuário pelo ID
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

// GetUserByEmail busca um usuário pelo email
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

// hashPassword gera um hash bcrypt da senha
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword verifica se a senha corresponde ao hash
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
