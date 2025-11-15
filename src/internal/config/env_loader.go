package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadEnvFile loads the appropriate .env file based on the ENVIRONMENT variable
// If ENVIRONMENT is not set, defaults to "local"
// Looks for .env.{environment} file, falls back to .env if not found
func LoadEnvFile() error {
	// Get environment, default to "local"
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "local"
	}

	// Get project root (assuming we're in src/internal/config)
	projectRoot, err := getProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	// Try to load .env.{environment} first
	envFile := filepath.Join(projectRoot, fmt.Sprintf(".env.%s", environment))
	if err := godotenv.Load(envFile); err == nil {
		fmt.Printf("Loaded environment file: .env.%s\n", environment)
		return nil
	}

	// Fall back to .env
	envFile = filepath.Join(projectRoot, ".env")
	if err := godotenv.Load(envFile); err == nil {
		fmt.Printf("Loaded environment file: .env\n")
		return nil
	}

	// If no .env file exists, just continue (environment variables might be set externally)
	fmt.Printf("No .env file found, using system environment variables\n")
	return nil
}

// getProjectRoot finds the project root by looking for go.mod
func getProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up the directory tree until we find go.mod
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}
