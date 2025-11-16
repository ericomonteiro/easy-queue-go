# WhatsApp Integration

This document describes how to integrate and use the WhatsApp Business API in Easy Queue.

## Overview

The WhatsApp integration allows you to:
- Send text messages to customers
- Send template messages (pre-approved by Meta)
- Receive messages from customers via webhooks
- Process incoming messages automatically

## Configuration

### 1. Environment Variables

Add the following variables to your `.env` file:

```bash
# WhatsApp Configuration
WHATSAPP_ACCESS_TOKEN=your-whatsapp-access-token
WHATSAPP_PHONE_NUMBER_ID=your-phone-number-id
WHATSAPP_BUSINESS_ID=your-business-account-id
WHATSAPP_WEBHOOK_TOKEN=your-custom-webhook-verify-token
WHATSAPP_API_VERSION=v18.0
WHATSAPP_API_URL=https://graph.facebook.com
```

### 2. Getting Your Credentials

1. Go to [Meta for Developers](https://developers.facebook.com)
2. Create or select your app
3. Add the WhatsApp product
4. Navigate to **WhatsApp → Getting Started**
5. Copy the following values:
   - **Access Token** → `WHATSAPP_ACCESS_TOKEN`
   - **Phone Number ID** → `WHATSAPP_PHONE_NUMBER_ID`
   - **Business Account ID** → `WHATSAPP_BUSINESS_ID`
6. Create a custom token for webhook verification → `WHATSAPP_WEBHOOK_TOKEN`

## API Endpoints

### Debug Endpoints (Authenticated)

These endpoints are for testing and debugging. They require authentication.

#### 1. Check Status

```http
GET /debug/whatsapp/status
Authorization: Bearer <your-jwt-token>
```

**Response:**
```json
{
  "status": "active",
  "configured": true,
  "message": "WhatsApp integration is configured and ready"
}
```

#### 2. Send Text Message (Simple)

```http
POST /debug/whatsapp/send-text
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "to": "5511999999999",
  "message": "Hello from Easy Queue!"
}
```

**Response:**
```json
{
  "success": true,
  "message_id": "wamid.HBgNNTUxMTk5OTk5OTk5ORUCABIYIDNBNzE4QjQxRjQwRDhGNzg5OTcyMjI4OTk5OTk5OTk5AA==",
  "to": "5511999999999",
  "type": "text",
  "sent_at": "2024-01-15T10:30:00Z"
}
```

#### 3. Send Message (Advanced)

```http
POST /debug/whatsapp/send
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "to": "5511999999999",
  "type": "text",
  "message": "Hello from Easy Queue!"
}
```

**Supported Types:**
- `text` - Simple text message ✅
- `template` - Template message ✅
- `image` - Image message (coming soon)
- `document` - Document message (coming soon)

#### 4. Send Template Message

```http
POST /debug/whatsapp/send-template
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "to": "5511999999999",
  "template": {
    "name": "hello_world",
    "language": "en_US"
  }
}
```

**With Parameters:**
```json
{
  "to": "5511999999999",
  "template": {
    "name": "order_confirmation",
    "language": "pt_BR",
    "components": [
      {
        "type": "body",
        "parameters": [
          {
            "type": "text",
            "text": "João Silva"
          },
          {
            "type": "text",
            "text": "12345"
          }
        ]
      }
    ]
  }
}
```

**Response:**
```json
{
  "success": true,
  "message_id": "wamid.HBgNNTUxMTk5OTk5OTk5ORUCABIYIDNBNzE4QjQxRjQwRDhGNzg5OTcyMjI4OTk5OTk5OTk5AA==",
  "to": "5511999999999",
  "type": "template",
  "sent_at": "2024-01-15T10:30:00Z"
}
```

### Webhook Endpoints (Public)

These endpoints are called by Meta and must be publicly accessible.

#### 1. Verify Webhook

```http
GET /whatsapp/webhook?hub.mode=subscribe&hub.verify_token=your-token&hub.challenge=challenge-string
```

This endpoint is called by Meta to verify your webhook URL.

#### 2. Receive Messages

```http
POST /whatsapp/webhook
Content-Type: application/json

{
  "object": "whatsapp_business_account",
  "entry": [...]
}
```

This endpoint receives incoming messages from WhatsApp.

## Testing

### Using cURL

```bash
# 1. Login to get JWT token
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# 2. Send a test message
curl -X POST http://localhost:8080/debug/whatsapp/send-text \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "to": "5511999999999",
    "message": "Hello from Easy Queue!"
  }'
```

### Phone Number Format

Always use the international format without `+` or spaces:
- ✅ Correct: `5511999999999` (Brazil)
- ✅ Correct: `14155551234` (USA)
- ❌ Wrong: `+55 11 99999-9999`
- ❌ Wrong: `(11) 99999-9999`

## Creating Message Templates

Templates must be created and approved by Meta before you can use them.

### 1. Create a Template

1. Go to [Meta Business Manager](https://business.facebook.com)
2. Navigate to **WhatsApp Manager → Message Templates**
3. Click **Create Template**
4. Fill in the template details:
   - **Name**: `order_confirmation` (lowercase, underscores only)
   - **Category**: Choose appropriate category (MARKETING, UTILITY, AUTHENTICATION)
   - **Language**: Select language (e.g., Portuguese - Brazil)
   - **Content**: Add your message with variables

### 2. Template Example

**Template Name:** `order_confirmation`  
**Language:** `pt_BR`  
**Content:**
```
Olá {{1}}, seu pedido #{{2}} foi confirmado! 
Obrigado por comprar com a gente.
```

### 3. Using Variables

Variables in templates are represented as `{{1}}`, `{{2}}`, etc. When sending, you provide the values:

```json
{
  "to": "5511999999999",
  "template": {
    "name": "order_confirmation",
    "language": "pt_BR",
    "components": [
      {
        "type": "body",
        "parameters": [
          {"type": "text", "text": "João Silva"},
          {"type": "text", "text": "12345"}
        ]
      }
    ]
  }
}
```

### 4. Template Approval

- Templates are usually approved within 24-48 hours
- You'll receive an email notification when approved
- Once approved, you can use them immediately

### 5. Common Template Categories

- **MARKETING**: Promotional messages, offers, announcements
- **UTILITY**: Account updates, order updates, alerts
- **AUTHENTICATION**: OTP codes, verification messages

## Webhook Setup

To receive messages, you need to configure the webhook in Meta for Developers:

1. Go to **WhatsApp → Configuration**
2. Click **Edit** in the Webhook section
3. Enter your webhook URL: `https://your-domain.com/whatsapp/webhook`
4. Enter your verify token (same as `WHATSAPP_WEBHOOK_TOKEN`)
5. Click **Verify and Save**
6. Subscribe to the `messages` field

### ngrok for Local Testing

For local development, use ngrok to expose your local server:

```bash
# Install ngrok
brew install ngrok

# Start your server
go run src/internal/cmd/main.go

# In another terminal, expose port 8080
ngrok http 8080

# Use the ngrok URL in Meta webhook configuration
# Example: https://abc123.ngrok.io/whatsapp/webhook
```

## Architecture

```
┌─────────────┐
│   Meta API  │
│  (WhatsApp) │
└──────┬──────┘
       │
       │ Webhook
       ▼
┌─────────────────┐
│ WhatsApp Handler│
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│WhatsApp Service │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   HTTP Client   │
│  (Send Messages)│
└─────────────────┘
```

## Error Handling

The service handles common errors:

- **Invalid phone number**: Returns 400 Bad Request
- **API rate limit**: Returns 500 with error details
- **Invalid access token**: Returns 500 with authentication error
- **Network timeout**: Returns 500 after 30 seconds

## Message Costs

WhatsApp charges per conversation (24-hour window):

- **Free tier**: 1,000 conversations/month
- **After free tier**: ~$0.005 - $0.09 per conversation (varies by country)

A conversation starts when:
1. You send a template message to a user
2. A user sends you a message (and you reply within 24h)

## Limitations

- Template messages require approval from Meta (24-48h)
- You can only send proactive messages using approved templates
- Messages to users must be within 24h of their last message (or use templates)
- Rate limits apply (check Meta documentation)

## Next Steps

- [ ] Implement template message support
- [ ] Add image and document support
- [ ] Create queue integration for async processing
- [ ] Add message status tracking
- [ ] Implement conversation management
- [ ] Add analytics and reporting

## Troubleshooting

### "WhatsApp integration not configured"

Make sure all required environment variables are set in your `.env` file.

### "API returned status 401"

Your access token is invalid or expired. Generate a new token in Meta for Developers.

### "API returned status 400"

Check the phone number format. It must be in international format without `+` or spaces.

### Webhook not receiving messages

1. Verify the webhook URL is publicly accessible
2. Check that the verify token matches
3. Ensure you subscribed to the `messages` field
4. Check server logs for errors

## Resources

- [WhatsApp Business API Documentation](https://developers.facebook.com/docs/whatsapp)
- [Meta for Developers](https://developers.facebook.com)
- [WhatsApp Business Platform](https://business.whatsapp.com)
