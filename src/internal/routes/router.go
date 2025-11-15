package routes

import (
	"easy-queue-go/src/internal/handlers"
	"easy-queue-go/src/internal/middleware"
	"easy-queue-go/src/internal/models"
	"easy-queue-go/src/internal/services"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// SetupRouter configures and returns the router with all routes defined
func SetupRouter(
	serviceName string,
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	authService services.AuthService,
) *gin.Engine {
	router := gin.Default()

	// Add logger middleware with request-id
	router.Use(middleware.LoggerMiddleware())

	// Add OpenTelemetry middleware for automatic tracing
	router.Use(otelgin.Middleware(serviceName))

	// Public routes
	router.GET("/health", handlers.HealthCheck)

	// Auth routes (public)
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.RefreshToken)
	}

	// User registration (public)
	router.POST("/users", userHandler.CreateUser)

	// Protected routes - require authentication
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		// User routes (authenticated users can access their own data)
		usersGroup := protected.Group("/users")
		{
			usersGroup.GET("/:id", userHandler.GetUserByID)
			usersGroup.GET("/by-email", userHandler.GetUserByEmail)
		}

		// Admin-only routes
		adminGroup := protected.Group("/admin")
		adminGroup.Use(middleware.RequireRole(models.RoleAdmin))
		{
			adminGroup.GET("/users", userHandler.ListAllUsers)
		}
	}

	return router
}
