package middleware

import (
	"easy-queue-go/src/internal/log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LoggerMiddleware adds a logger with request-id to the context
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate a unique request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set the request ID in the response header
		c.Header("X-Request-ID", requestID)

		// Get or create a logger with the request ID field
		ctx := c.Request.Context()
		ctx = log.Initialize(ctx)
		
		// Add request ID field to the logger
		ctx = log.WithField(ctx, zap.String("request_id", requestID))

		// Update the request context
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
