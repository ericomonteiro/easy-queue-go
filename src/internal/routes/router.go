package routes

import (
	"easy-queue-go/src/internal/handlers"
	"easy-queue-go/src/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// SetupRouter configura e retorna o router com todas as rotas definidas
func SetupRouter(serviceName string, userHandler *handlers.UserHandler) *gin.Engine {
	router := gin.Default()

	// Add logger middleware with request-id
	router.Use(middleware.LoggerMiddleware())

	// Add OpenTelemetry middleware for automatic tracing
	router.Use(otelgin.Middleware(serviceName))

	// Grupo de rotas de health
	router.GET("/health", handlers.HealthCheck)

	// Grupo de rotas de usuários
	usersGroup := router.Group("/users")
	{
		usersGroup.POST("", userHandler.CreateUser)
		usersGroup.GET("/:id", userHandler.GetUserByID)
		usersGroup.GET("/by-email", userHandler.GetUserByEmail)
	}

	return router
}
