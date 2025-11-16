package main

import (
	"context"
	"easy-queue-go/src/internal/config"
	"easy-queue-go/src/internal/handlers"
	"easy-queue-go/src/internal/infra/database"
	"easy-queue-go/src/internal/log"
	"easy-queue-go/src/internal/repositories"
	"easy-queue-go/src/internal/routes"
	"easy-queue-go/src/internal/services"
	"easy-queue-go/src/internal/tracing"

	"go.uber.org/zap"
)

// @title Easy Queue API
// @version 1.0
// @description API for Easy Queue - A queue management system
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@easyqueue.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables from .env file
	if err := config.LoadEnvFile(); err != nil {
		panic(err)
	}

	// Initialize zap logger
	ctx := log.Initialize(context.Background())

	// Initialize tracing
	tracingConfig := tracing.LoadConfigFromEnv()
	shutdownTracing, err := tracing.Initialize(ctx, tracingConfig)
	if err != nil {
		log.Fatal(ctx, "Failed to initialize tracing", zap.Error(err))
	}
	defer func() {
		if err := shutdownTracing(ctx); err != nil {
			log.Error(ctx, "Failed to shutdown tracing", zap.Error(err))
		}
	}()

	if tracingConfig.Enabled {
		log.Info(ctx, "Tracing initialized",
			zap.String("service", tracingConfig.ServiceName),
			zap.String("endpoint", tracingConfig.OTLPEndpoint),
		)
	}

	// Load configs
	log.Info(ctx, "Loading configs")
	configs := config.InitializeConfigs()

	// Initialize database
	dbClient, err := database.NewClient(ctx, configs.DB)
	if err != nil {
		log.Panic(ctx, "Failed to connect to database", zap.Error(err))
	}

	// Get database pool
	client, ok := dbClient.(*database.Client)
	if !ok {
		log.Fatal(ctx, "Failed to cast database client")
	}
	pool := client.Pool()

	// Initialize dependencies
	userRepo := repositories.NewUserRepository(pool)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Initialize business dependencies
	businessRepo := repositories.NewBusinessRepository(pool)
	businessService := services.NewBusinessService(businessRepo, userRepo)
	businessHandler := handlers.NewBusinessHandler(businessService)

	// Initialize auth service
	authService := services.NewAuthService(userRepo, services.AuthServiceConfig{
		JWTSecret:       configs.JWT.Secret,
		AccessTokenTTL:  configs.JWT.AccessTokenTTL,
		RefreshTokenTTL: configs.JWT.RefreshTokenTTL,
	})

	// Initialize auth handler
	authHandler := handlers.NewAuthHandler(authService)

	// Initialize WhatsApp service and handler
	var whatsappHandler *handlers.WhatsAppHandler
	if configs.WhatsApp != nil {
		whatsappService := services.NewWhatsAppService(configs.WhatsApp)
		whatsappHandler = handlers.NewWhatsAppHandler(whatsappService)
		log.Info(ctx, "WhatsApp integration initialized",
			zap.String("phone_number_id", configs.WhatsApp.PhoneNumberID),
		)
	} else {
		log.Warn(ctx, "WhatsApp integration not configured - WhatsApp features will be disabled")
		// Create a dummy handler to avoid nil pointer issues
		whatsappHandler = handlers.NewWhatsAppHandler(nil)
	}

	// Setup router
	router := routes.SetupRouter(tracingConfig.ServiceName, userHandler, authHandler, businessHandler, whatsappHandler, authService)

	// Start server
	log.Info(ctx, "Starting server on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(ctx, "Failed to start server", zap.Error(err))
	}
}
