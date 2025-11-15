package config

import (
	"fmt"
	"time"
)

// JWTConfig holds the JWT configuration
type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// LoadJWTConfig loads the JWT configuration from environment variables
func LoadJWTConfig() (*JWTConfig, error) {
	secret := getEnv("JWT_SECRET", "")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	// Parse access token TTL (default: 15 minutes)
	accessTTLStr := getEnv("JWT_ACCESS_TOKEN_TTL", "15m")
	accessTTL, err := time.ParseDuration(accessTTLStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_TOKEN_TTL: %w", err)
	}

	// Parse refresh token TTL (default: 7 days)
	refreshTTLStr := getEnv("JWT_REFRESH_TOKEN_TTL", "168h")
	refreshTTL, err := time.ParseDuration(refreshTTLStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_TOKEN_TTL: %w", err)
	}

	return &JWTConfig{
		Secret:          secret,
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
	}, nil
}
