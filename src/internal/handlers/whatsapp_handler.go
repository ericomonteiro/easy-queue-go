package handlers

import (
	"easy-queue-go/src/internal/log"
	"easy-queue-go/src/internal/models"
	"easy-queue-go/src/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// WhatsAppHandler manages HTTP requests related to WhatsApp
type WhatsAppHandler struct {
	whatsappService services.WhatsAppService
}

// NewWhatsAppHandler creates a new instance of WhatsAppHandler
func NewWhatsAppHandler(whatsappService services.WhatsAppService) *WhatsAppHandler {
	return &WhatsAppHandler{
		whatsappService: whatsappService,
	}
}

// SendMessage godoc
// @Summary Send a WhatsApp message (Debug)
// @Description Sends a WhatsApp message to a phone number. This is a debug endpoint for testing.
// @Tags whatsapp-debug
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body models.SendWhatsAppMessageRequest true "Message data"
// @Success 200 {object} models.WhatsAppMessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /debug/whatsapp/send [post]
func (h *WhatsAppHandler) SendMessage(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.SendWhatsAppMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn(ctx, "Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	log.Info(ctx, "Debug: Sending WhatsApp message",
		zap.String("to", req.To),
		zap.String("type", string(req.Type)),
	)

	response, err := h.whatsappService.SendMessage(ctx, &req)
	if err != nil {
		log.Error(ctx, "Failed to send WhatsApp message", zap.Error(err))
		
		// If response exists, return it with the error
		if response != nil {
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "send_failed",
			Message: err.Error(),
		})
		return
	}

	log.Info(ctx, "WhatsApp message sent successfully",
		zap.String("message_id", response.MessageID),
		zap.String("to", response.To),
	)

	c.JSON(http.StatusOK, response)
}

// SendTextMessage godoc
// @Summary Send a simple text message (Debug)
// @Description Sends a simple text message via WhatsApp. This is a debug endpoint for quick testing.
// @Tags whatsapp-debug
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{to=string,message=string} true "Phone number and message"
// @Success 200 {object} models.WhatsAppMessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /debug/whatsapp/send-text [post]
func (h *WhatsAppHandler) SendTextMessage(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		To      string `json:"to" binding:"required"`
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn(ctx, "Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	log.Info(ctx, "Debug: Sending simple text message",
		zap.String("to", req.To),
	)

	response, err := h.whatsappService.SendTextMessage(ctx, req.To, req.Message)
	if err != nil {
		log.Error(ctx, "Failed to send text message", zap.Error(err))
		
		if response != nil {
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "send_failed",
			Message: err.Error(),
		})
		return
	}

	log.Info(ctx, "Text message sent successfully",
		zap.String("message_id", response.MessageID),
		zap.String("to", response.To),
	)

	c.JSON(http.StatusOK, response)
}

// SendTemplateMessage godoc
// @Summary Send a template message (Debug)
// @Description Sends a template message via WhatsApp. Templates must be pre-approved by Meta.
// @Tags whatsapp-debug
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{to=string,template=models.WhatsAppTemplateRequest} true "Phone number and template"
// @Success 200 {object} models.WhatsAppMessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /debug/whatsapp/send-template [post]
func (h *WhatsAppHandler) SendTemplateMessage(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		To       string                           `json:"to" binding:"required"`
		Template *models.WhatsAppTemplateRequest  `json:"template" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn(ctx, "Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	log.Info(ctx, "Debug: Sending template message",
		zap.String("to", req.To),
		zap.String("template_name", req.Template.Name),
		zap.String("language", req.Template.Language),
	)

	response, err := h.whatsappService.SendTemplateMessage(ctx, req.To, req.Template)
	if err != nil {
		log.Error(ctx, "Failed to send template message", zap.Error(err))
		
		if response != nil {
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "send_failed",
			Message: err.Error(),
		})
		return
	}

	log.Info(ctx, "Template message sent successfully",
		zap.String("message_id", response.MessageID),
		zap.String("to", response.To),
		zap.String("template_name", req.Template.Name),
	)

	c.JSON(http.StatusOK, response)
}

// VerifyWebhook godoc
// @Summary Verify WhatsApp webhook
// @Description Verifies the webhook subscription from WhatsApp (called by Meta)
// @Tags whatsapp
// @Accept json
// @Produce plain
// @Param hub.mode query string true "Subscription mode"
// @Param hub.verify_token query string true "Verify token"
// @Param hub.challenge query string true "Challenge string"
// @Success 200 {string} string "Challenge string"
// @Failure 403 {object} ErrorResponse
// @Router /whatsapp/webhook [get]
func (h *WhatsAppHandler) VerifyWebhook(c *gin.Context) {
	ctx := c.Request.Context()

	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	log.Info(ctx, "WhatsApp webhook verification request",
		zap.String("mode", mode),
	)

	challengeResponse, err := h.whatsappService.VerifyWebhook(mode, token, challenge)
	if err != nil {
		log.Warn(ctx, "Webhook verification failed", zap.Error(err))
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "verification_failed",
			Message: "Invalid verification token",
		})
		return
	}

	log.Info(ctx, "Webhook verification successful")
	c.String(http.StatusOK, challengeResponse)
}

// ReceiveWebhook godoc
// @Summary Receive WhatsApp webhook events
// @Description Receives webhook events from WhatsApp (called by Meta when messages are received)
// @Tags whatsapp
// @Accept json
// @Produce json
// @Param payload body models.WhatsAppWebhookPayload true "Webhook payload"
// @Success 200 {object} object{status=string}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /whatsapp/webhook [post]
func (h *WhatsAppHandler) ReceiveWebhook(c *gin.Context) {
	ctx := c.Request.Context()

	var payload models.WhatsAppWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Warn(ctx, "Invalid webhook payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_payload",
			Message: err.Error(),
		})
		return
	}

	log.Info(ctx, "Received WhatsApp webhook",
		zap.String("object", payload.Object),
		zap.Int("entries", len(payload.Entry)),
	)

	if err := h.whatsappService.ProcessWebhook(ctx, &payload); err != nil {
		log.Error(ctx, "Failed to process webhook", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "processing_failed",
			Message: err.Error(),
		})
		return
	}

	log.Info(ctx, "Webhook processed successfully")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetStatus godoc
// @Summary Get WhatsApp integration status (Debug)
// @Description Returns the status of the WhatsApp integration
// @Tags whatsapp-debug
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,configured=bool}
// @Failure 401 {object} ErrorResponse
// @Router /debug/whatsapp/status [get]
func (h *WhatsAppHandler) GetStatus(c *gin.Context) {
	ctx := c.Request.Context()

	log.Info(ctx, "Debug: Getting WhatsApp status")

	// Simple status check
	c.JSON(http.StatusOK, gin.H{
		"status":     "active",
		"configured": h.whatsappService != nil,
		"message":    "WhatsApp integration is configured and ready",
	})
}
