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
	"easy-queue-go/src/internal/singletons"
	"easy-queue-go/src/internal/tracing"

	"go.uber.org/zap"
)

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

	// Initialize singletons (DB, configs, etc)
	instances := singletons.Initialize(ctx)

	// Get database pool
	dbClient, ok := instances.DB.(*database.Client)
	if !ok {
		log.Fatal(ctx, "Failed to cast database client")
	}
	pool := dbClient.Pool()

	// Get configs
	configs := instances.Configs

	// Initialize dependencies
	userRepo := repositories.NewUserRepository(pool)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Initialize auth service and handler
	authService := services.NewAuthService(userRepo, services.AuthServiceConfig{
		JWTSecret:       configs.JWT.Secret,
		AccessTokenTTL:  configs.JWT.AccessTokenTTL,
		RefreshTokenTTL: configs.JWT.RefreshTokenTTL,
	})
	authHandler := handlers.NewAuthHandler(authService)

	// Setup router
	router := routes.SetupRouter(tracingConfig.ServiceName, userHandler, authHandler, authService)

	// Start server
	log.Info(ctx, "Starting server on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(ctx, "Failed to start server", zap.Error(err))
	}
}
