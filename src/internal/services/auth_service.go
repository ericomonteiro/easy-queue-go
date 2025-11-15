package services

import (
	"context"
	"easy-queue-go/src/internal/log"
	"easy-queue-go/src/internal/models"
	"easy-queue-go/src/internal/repositories"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.RefreshTokenResponse, error)
	ValidateToken(ctx context.Context, tokenString string, tokenType models.TokenType) (*models.JWTClaims, error)
}

// authService implements AuthService
type authService struct {
	userRepo         repositories.UserRepository
	jwtSecret        string
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
}

// AuthServiceConfig holds configuration for the auth service
type AuthServiceConfig struct {
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRepo repositories.UserRepository, config AuthServiceConfig) AuthService {
	return &authService{
		userRepo:        userRepo,
		jwtSecret:       config.JWTSecret,
		accessTokenTTL:  config.AccessTokenTTL,
		refreshTokenTTL: config.RefreshTokenTTL,
	}
}

// Login authenticates a user and returns JWT tokens
func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	ctx, span := tracer.Start(ctx, "AuthService.Login",
		trace.WithAttributes(
			attribute.String("email", req.Email),
		),
	)
	defer span.End()

	log.Info(ctx, "User login attempt", zap.String("email", req.Email))

	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		log.Warn(ctx, "Login failed: user not found", zap.String("email", req.Email))
		span.RecordError(err)
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		log.Warn(ctx, "Login failed: user is inactive",
			zap.String("email", req.Email),
			zap.String("user_id", user.ID.String()),
		)
		return nil, fmt.Errorf("user account is inactive")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		log.Warn(ctx, "Login failed: invalid password", zap.String("email", req.Email))
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate access token
	accessToken, err := s.generateToken(user, models.TokenTypeAccess, s.accessTokenTTL)
	if err != nil {
		log.Error(ctx, "Failed to generate access token", zap.Error(err))
		span.RecordError(err)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.generateToken(user, models.TokenTypeRefresh, s.refreshTokenTTL)
	if err != nil {
		log.Error(ctx, "Failed to generate refresh token", zap.Error(err))
		span.RecordError(err)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	log.Info(ctx, "User logged in successfully",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email),
	)

	span.SetAttributes(attribute.String("user_id", user.ID.String()))

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.accessTokenTTL.Seconds()),
		User:         user.ToResponse(),
	}, nil
}

// RefreshToken generates new access and refresh tokens using a valid refresh token
func (s *authService) RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.RefreshTokenResponse, error) {
	ctx, span := tracer.Start(ctx, "AuthService.RefreshToken")
	defer span.End()

	log.Info(ctx, "Token refresh attempt")

	// Validate refresh token
	claims, err := s.ValidateToken(ctx, req.RefreshToken, models.TokenTypeRefresh)
	if err != nil {
		log.Warn(ctx, "Token refresh failed: invalid refresh token", zap.Error(err))
		span.RecordError(err)
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Get user from database to ensure they still exist and are active
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		log.Warn(ctx, "Token refresh failed: user not found",
			zap.String("user_id", claims.UserID.String()),
		)
		span.RecordError(err)
		return nil, fmt.Errorf("user not found")
	}

	if !user.IsActive {
		log.Warn(ctx, "Token refresh failed: user is inactive",
			zap.String("user_id", user.ID.String()),
		)
		return nil, fmt.Errorf("user account is inactive")
	}

	// Generate new access token
	accessToken, err := s.generateToken(user, models.TokenTypeAccess, s.accessTokenTTL)
	if err != nil {
		log.Error(ctx, "Failed to generate new access token", zap.Error(err))
		span.RecordError(err)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate new refresh token
	refreshToken, err := s.generateToken(user, models.TokenTypeRefresh, s.refreshTokenTTL)
	if err != nil {
		log.Error(ctx, "Failed to generate new refresh token", zap.Error(err))
		span.RecordError(err)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	log.Info(ctx, "Tokens refreshed successfully",
		zap.String("user_id", user.ID.String()),
	)

	span.SetAttributes(attribute.String("user_id", user.ID.String()))

	return &models.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.accessTokenTTL.Seconds()),
	}, nil
}

// ValidateToken validates a JWT token and returns its claims
func (s *authService) ValidateToken(ctx context.Context, tokenString string, expectedType models.TokenType) (*models.JWTClaims, error) {
	_, span := tracer.Start(ctx, "AuthService.ValidateToken",
		trace.WithAttributes(
			attribute.String("token_type", string(expectedType)),
		),
	)
	defer span.End()

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*models.JWTClaims)
	if !ok || !token.Valid {
		err := fmt.Errorf("invalid token claims")
		span.RecordError(err)
		return nil, err
	}

	// Verify token type
	if claims.Type != expectedType {
		err := fmt.Errorf("invalid token type: expected %s, got %s", expectedType, claims.Type)
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(attribute.String("user_id", claims.UserID.String()))

	return claims, nil
}

// generateToken creates a new JWT token for a user
func (s *authService) generateToken(user *models.User, tokenType models.TokenType, ttl time.Duration) (string, error) {
	now := time.Now()
	expiresAt := now.Add(ttl)

	claims := &models.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Roles:  user.Roles,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "easy-queue-go",
			Subject:   user.ID.String(),
			ID:        uuid.New().String(), // Unique token ID (jti)
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}
