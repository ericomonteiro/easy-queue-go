package config

import (
	"fmt"
	"os"
)

// WhatsAppConfig holds the WhatsApp API configuration
type WhatsAppConfig struct {
	AccessToken   string
	PhoneNumberID string
	BusinessID    string
	WebhookToken  string
	APIVersion    string
	APIURL        string
	AppID         string
	AppSecret     string
}

// LoadWhatsAppConfig loads WhatsApp configuration from environment variables
func LoadWhatsAppConfig() (*WhatsAppConfig, error) {
	accessToken := os.Getenv("WHATSAPP_ACCESS_TOKEN")
	if accessToken == "" {
		return nil, fmt.Errorf("WHATSAPP_ACCESS_TOKEN is required")
	}

	phoneNumberID := os.Getenv("WHATSAPP_PHONE_NUMBER_ID")
	if phoneNumberID == "" {
		return nil, fmt.Errorf("WHATSAPP_PHONE_NUMBER_ID is required")
	}

	businessID := getEnv("WHATSAPP_BUSINESS_ID", "")
	webhookToken := getEnv("WHATSAPP_WEBHOOK_TOKEN", "easy-queue-webhook-token")
	apiVersion := getEnv("WHATSAPP_API_VERSION", "v18.0")
	apiURL := getEnv("WHATSAPP_API_URL", "https://graph.facebook.com")
	appID := getEnv("WHATSAPP_APP_ID", "")
	appSecret := getEnv("WHATSAPP_APP_SECRET", "")

	return &WhatsAppConfig{
		AccessToken:   accessToken,
		PhoneNumberID: phoneNumberID,
		BusinessID:    businessID,
		WebhookToken:  webhookToken,
		APIVersion:    apiVersion,
		APIURL:        apiURL,
		AppID:         appID,
		AppSecret:     appSecret,
	}, nil
}
