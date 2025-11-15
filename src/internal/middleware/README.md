# Middleware

## Logger Middleware

The `LoggerMiddleware` adds a logger with a `request_id` field to the request context.

### Features

- **Request ID Generation**: Automatically generates a unique UUID for each request
- **X-Request-ID Header Support**: Uses the `X-Request-ID` header if provided, otherwise generates a new one
- **Response Header**: Sets the `X-Request-ID` header in the response for request tracing
- **Context Logger**: Adds a logger with the request ID to the context for use throughout the request lifecycle

### Usage

The middleware is automatically applied to all routes in the router setup:

```go
router.Use(middleware.LoggerMiddleware())
```

### Logging in Handlers

All log calls will automatically include the `request_id` field:

```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    ctx := c.Request.Context()
    
    // This log will automatically include the request_id
    log.Info(ctx, "User created successfully",
        zap.String("user_id", user.ID.String()),
        zap.String("email", user.Email),
    )
}
```

### Example Log Output

```json
{
  "level": "info",
  "ts": 1700000000.123456,
  "caller": "handlers/user_handler.go:87",
  "msg": "User created successfully",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "user@example.com"
}
```
