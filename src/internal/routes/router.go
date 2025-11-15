package routes

import (
	"easy-queue-go/src/internal/handlers"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// SetupRouter configura e retorna o router com todas as rotas definidas
func SetupRouter(serviceName string) *gin.Engine {
	router := gin.Default()

	// Add OpenTelemetry middleware for automatic tracing
	router.Use(otelgin.Middleware(serviceName))

	// Grupo de rotas de health
	router.GET("/health", handlers.HealthCheck)

	return router
}
