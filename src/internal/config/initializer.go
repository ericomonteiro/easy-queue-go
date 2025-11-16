package config

import (
	"log"
	"os"
)

type Configs struct {
	DB       *DBConfig
	JWT      *JWTConfig
	WhatsApp *WhatsAppConfig
}

// InitializeConfigs initializes the configs
func InitializeConfigs() *Configs {
	dbConfigs := LoadDBConfigs()
	
	jwtConfig, err := LoadJWTConfig()
	if err != nil {
		log.Fatalf("Failed to load JWT config: %v", err)
	}

	whatsappConfig, err := LoadWhatsAppConfig()
	if err != nil {
		log.Printf("Warning: Failed to load WhatsApp config: %v (WhatsApp features will be disabled)", err)
	}

	return &Configs{
		DB:       dbConfigs,
		JWT:      jwtConfig,
		WhatsApp: whatsappConfig,
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
