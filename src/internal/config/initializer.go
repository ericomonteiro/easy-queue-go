package config

import "os"

type Configs struct {
	DB *DBConfig
}

// InitializeConfigs initializes the configs
func InitializeConfigs() *Configs {
	dbConfigs := LoadDBConfigs()

	return &Configs{
		DB: dbConfigs,
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
