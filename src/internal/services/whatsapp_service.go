package services

import (
	"bytes"
	"context"
	"easy-queue-go/src/internal/config"
	"easy-queue-go/src/internal/log"
	"easy-queue-go/src/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

var whatsappTracer = otel.Tracer("whatsapp-service")

// WhatsAppService defines the interface for WhatsApp operations
type WhatsAppService interface {
	SendTextMessage(ctx context.Context, to, message string) (*models.WhatsAppMessageResponse, error)
	SendTemplateMessage(ctx context.Context, to string, template *models.WhatsAppTemplateRequest) (*models.WhatsAppMessageResponse, error)
	SendMessage(ctx context.Context, req *models.SendWhatsAppMessageRequest) (*models.WhatsAppMessageResponse, error)
	VerifyWebhook(mode, token, challenge string) (string, error)
	ProcessWebhook(ctx context.Context, payload *models.WhatsAppWebhookPayload) error
}

type whatsappService struct {
	config     *config.WhatsAppConfig
	httpClient *http.Client
}

// NewWhatsAppService creates a new WhatsApp service instance
func NewWhatsAppService(config *config.WhatsAppConfig) WhatsAppService {
	return &whatsappService{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendTextMessage sends a simple text message via WhatsApp
func (s *whatsappService) SendTextMessage(ctx context.Context, to, message string) (*models.WhatsAppMessageResponse, error) {
	ctx, span := whatsappTracer.Start(ctx, "whatsapp.SendTextMessage")
	defer span.End()

	span.SetAttributes(
		attribute.String("whatsapp.to", to),
		attribute.String("whatsapp.type", "text"),
	)

	log.Info(ctx, "Sending WhatsApp text message",
		zap.String("to", to),
		zap.Int("message_length", len(message)),
	)

	// Build the API URL
	url := fmt.Sprintf("%s/%s/%s/messages",
		s.config.APIURL,
		s.config.APIVersion,
		s.config.PhoneNumberID,
	)

	// Build the payload
	payload := models.WhatsAppMessagePayload{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "text",
		Text: &models.WhatsAppTextMessage{
			Body: message,
		},
	}

	// Marshal payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Error(ctx, "Failed to marshal WhatsApp payload", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error(ctx, "Failed to create HTTP request", zap.Error(err))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.AccessToken))

	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Error(ctx, "Failed to send WhatsApp message", zap.Error(err))
		return nil, fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(ctx, "Failed to read response body", zap.Error(err))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Error(ctx, "WhatsApp API returned error",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(body)),
		)
		return &models.WhatsAppMessageResponse{
			Success: false,
			To:      to,
			Type:    "text",
			SentAt:  time.Now(),
			Error:   fmt.Sprintf("API error: %s", string(body)),
		}, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var apiResponse models.WhatsAppAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Error(ctx, "Failed to parse WhatsApp API response", zap.Error(err))
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Build response
	messageID := ""
	if len(apiResponse.Messages) > 0 {
		messageID = apiResponse.Messages[0].ID
	}

	response := &models.WhatsAppMessageResponse{
		Success:   true,
		MessageID: messageID,
		To:        to,
		Type:      "text",
		SentAt:    time.Now(),
	}

	log.Info(ctx, "WhatsApp message sent successfully",
		zap.String("message_id", messageID),
		zap.String("to", to),
	)

	span.SetAttributes(
		attribute.String("whatsapp.message_id", messageID),
		attribute.Bool("whatsapp.success", true),
	)

	return response, nil
}

// SendTemplateMessage sends a template message via WhatsApp
func (s *whatsappService) SendTemplateMessage(ctx context.Context, to string, template *models.WhatsAppTemplateRequest) (*models.WhatsAppMessageResponse, error) {
	ctx, span := whatsappTracer.Start(ctx, "whatsapp.SendTemplateMessage")
	defer span.End()

	span.SetAttributes(
		attribute.String("whatsapp.to", to),
		attribute.String("whatsapp.type", "template"),
		attribute.String("whatsapp.template_name", template.Name),
		attribute.String("whatsapp.template_language", template.Language),
	)

	log.Info(ctx, "Sending WhatsApp template message",
		zap.String("to", to),
		zap.String("template_name", template.Name),
		zap.String("language", template.Language),
	)

	// Build the API URL
	url := fmt.Sprintf("%s/%s/%s/messages",
		s.config.APIURL,
		s.config.APIVersion,
		s.config.PhoneNumberID,
	)

	// Build the template payload
	payload := models.WhatsAppMessagePayload{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "template",
		Template: &models.WhatsAppTemplate{
			Name: template.Name,
			Language: models.WhatsAppTemplateLanguage{
				Code: template.Language,
			},
			Components: template.Components,
		},
	}

	// Marshal payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Error(ctx, "Failed to marshal WhatsApp template payload", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	log.Debug(ctx, "Template payload", zap.String("json", string(jsonData)))

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error(ctx, "Failed to create HTTP request", zap.Error(err))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.AccessToken))

	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Error(ctx, "Failed to send WhatsApp template message", zap.Error(err))
		return nil, fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(ctx, "Failed to read response body", zap.Error(err))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Error(ctx, "WhatsApp API returned error",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(body)),
		)
		return &models.WhatsAppMessageResponse{
			Success: false,
			To:      to,
			Type:    "template",
			SentAt:  time.Now(),
			Error:   fmt.Sprintf("API error: %s", string(body)),
		}, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var apiResponse models.WhatsAppAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Error(ctx, "Failed to parse WhatsApp API response", zap.Error(err))
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Build response
	messageID := ""
	if len(apiResponse.Messages) > 0 {
		messageID = apiResponse.Messages[0].ID
	}

	response := &models.WhatsAppMessageResponse{
		Success:   true,
		MessageID: messageID,
		To:        to,
		Type:      "template",
		SentAt:    time.Now(),
	}

	log.Info(ctx, "WhatsApp template message sent successfully",
		zap.String("message_id", messageID),
		zap.String("to", to),
		zap.String("template_name", template.Name),
	)

	span.SetAttributes(
		attribute.String("whatsapp.message_id", messageID),
		attribute.Bool("whatsapp.success", true),
	)

	return response, nil
}

// SendMessage sends a WhatsApp message based on the request type
func (s *whatsappService) SendMessage(ctx context.Context, req *models.SendWhatsAppMessageRequest) (*models.WhatsAppMessageResponse, error) {
	ctx, span := whatsappTracer.Start(ctx, "whatsapp.SendMessage")
	defer span.End()

	span.SetAttributes(
		attribute.String("whatsapp.to", req.To),
		attribute.String("whatsapp.type", string(req.Type)),
	)

	log.Info(ctx, "Sending WhatsApp message",
		zap.String("to", req.To),
		zap.String("type", string(req.Type)),
	)

	// Route to appropriate handler based on message type
	switch req.Type {
	case models.WhatsAppMessageTypeText:
		return s.SendTextMessage(ctx, req.To, req.Message)
	case models.WhatsAppMessageTypeTemplate:
		if req.Template == nil {
			return nil, fmt.Errorf("template is required for template messages")
		}
		return s.SendTemplateMessage(ctx, req.To, req.Template)
	case models.WhatsAppMessageTypeImage, models.WhatsAppMessageTypeDocument:
		return nil, fmt.Errorf("message type %s not yet implemented", req.Type)
	default:
		return nil, fmt.Errorf("unsupported message type: %s", req.Type)
	}
}

// VerifyWebhook verifies the webhook subscription from WhatsApp
func (s *whatsappService) VerifyWebhook(mode, token, challenge string) (string, error) {
	if mode != "subscribe" {
		return "", fmt.Errorf("invalid mode: %s", mode)
	}

	if token != s.config.WebhookToken {
		return "", fmt.Errorf("invalid verify token")
	}

	return challenge, nil
}

// ProcessWebhook processes incoming webhook events from WhatsApp
func (s *whatsappService) ProcessWebhook(ctx context.Context, payload *models.WhatsAppWebhookPayload) error {
	ctx, span := whatsappTracer.Start(ctx, "whatsapp.ProcessWebhook")
	defer span.End()

	log.Info(ctx, "Processing WhatsApp webhook",
		zap.String("object", payload.Object),
		zap.Int("entries", len(payload.Entry)),
	)

	// Process each entry
	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			if change.Field != "messages" {
				continue
			}

			// Process each message
			for _, message := range change.Value.Messages {
				log.Info(ctx, "Received WhatsApp message",
					zap.String("from", message.From),
					zap.String("message_id", message.ID),
					zap.String("type", message.Type),
					zap.String("text", message.Text.Body),
				)

				span.SetAttributes(
					attribute.String("whatsapp.from", message.From),
					attribute.String("whatsapp.message_id", message.ID),
					attribute.String("whatsapp.type", message.Type),
				)

				// TODO: Process the message (e.g., send to queue, trigger business logic)
				// For now, just log it
			}
		}
	}

	return nil
}
