package main

import (
	"context"
	"easy-queue-go/src/internal/config"
	"easy-queue-go/src/internal/log"
	"easy-queue-go/src/internal/routes"
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
	singletons.Initialize(ctx)

	// Setup router
	router := routes.SetupRouter(tracingConfig.ServiceName)

	// Start server
	log.Info(ctx, "Starting server on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(ctx, "Failed to start server", zap.Error(err))
	}
}
