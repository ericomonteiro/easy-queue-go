package models

import "time"

// WhatsAppMessageType represents the type of WhatsApp message
type WhatsAppMessageType string

const (
	WhatsAppMessageTypeText     WhatsAppMessageType = "text"
	WhatsAppMessageTypeTemplate WhatsAppMessageType = "template"
	WhatsAppMessageTypeImage    WhatsAppMessageType = "image"
	WhatsAppMessageTypeDocument WhatsAppMessageType = "document"
)

// SendWhatsAppMessageRequest represents a request to send a WhatsApp message
type SendWhatsAppMessageRequest struct {
	To               string                       `json:"to" binding:"required"`
	Type             WhatsAppMessageType          `json:"type" binding:"required,oneof=text template image document"`
	Message          string                       `json:"message" binding:"required_if=Type text"`
	Template         *WhatsAppTemplateRequest     `json:"template" binding:"required_if=Type template"`
}

// WhatsAppTextMessage represents the text content of a WhatsApp message
type WhatsAppTextMessage struct {
	Body string `json:"body"`
}

// WhatsAppTemplateRequest represents a template message request
type WhatsAppTemplateRequest struct {
	Name       string                          `json:"name" binding:"required"`
	Language   string                          `json:"language" binding:"required"`
	Components []WhatsAppTemplateComponent     `json:"components,omitempty"`
}

// WhatsAppTemplateComponent represents a component in a template
type WhatsAppTemplateComponent struct {
	Type       string                           `json:"type" binding:"required,oneof=header body button"`
	Parameters []WhatsAppTemplateParameter      `json:"parameters,omitempty"`
}

// WhatsAppTemplateParameter represents a parameter in a template component
type WhatsAppTemplateParameter struct {
	Type string `json:"type" binding:"required,oneof=text currency date_time image document video"`
	Text string `json:"text,omitempty"`
}

// WhatsAppTemplate represents the template structure for API payload
type WhatsAppTemplate struct {
	Name       string                       `json:"name"`
	Language   WhatsAppTemplateLanguage     `json:"language"`
	Components []WhatsAppTemplateComponent  `json:"components,omitempty"`
}

// WhatsAppTemplateLanguage represents the language code for a template
type WhatsAppTemplateLanguage struct {
	Code string `json:"code"`
}

// WhatsAppMessagePayload represents the payload sent to WhatsApp API
type WhatsAppMessagePayload struct {
	MessagingProduct string                `json:"messaging_product"`
	To               string                `json:"to"`
	Type             string                `json:"type"`
	Text             *WhatsAppTextMessage  `json:"text,omitempty"`
	Template         *WhatsAppTemplate     `json:"template,omitempty"`
}

// WhatsAppAPIResponse represents the response from WhatsApp API
type WhatsAppAPIResponse struct {
	MessagingProduct string `json:"messaging_product"`
	Contacts         []struct {
		Input string `json:"input"`
		WaID  string `json:"wa_id"`
	} `json:"contacts"`
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
}

// WhatsAppMessageResponse represents the response after sending a message
type WhatsAppMessageResponse struct {
	Success   bool      `json:"success"`
	MessageID string    `json:"message_id,omitempty"`
	To        string    `json:"to"`
	Type      string    `json:"type"`
	SentAt    time.Time `json:"sent_at"`
	Error     string    `json:"error,omitempty"`
}

// WhatsAppWebhookPayload represents the webhook payload from WhatsApp
type WhatsAppWebhookPayload struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberID      string `json:"phone_number_id"`
				} `json:"metadata"`
				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					} `json:"profile"`
					WaID string `json:"wa_id"`
				} `json:"contacts"`
				Messages []struct {
					From      string `json:"from"`
					ID        string `json:"id"`
					Timestamp string `json:"timestamp"`
					Type      string `json:"type"`
					Text      struct {
						Body string `json:"body"`
					} `json:"text"`
				} `json:"messages"`
			} `json:"value"`
			Field string `json:"field"`
		} `json:"changes"`
	} `json:"entry"`
}

// WhatsAppWebhookVerification represents the webhook verification request
type WhatsAppWebhookVerification struct {
	Mode      string `form:"hub.mode"`
	Token     string `form:"hub.verify_token"`
	Challenge string `form:"hub.challenge"`
}
