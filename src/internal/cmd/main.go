package main

import (
	"context"
	"easy-queue-go/src/internal/log"
	"easy-queue-go/src/internal/routes"
	"easy-queue-go/src/internal/singletons"

	"go.uber.org/zap"
)

func main() {
	// Initialize zap logger
	ctx := log.Initialize(context.Background())

	// Initialize singletons (DB, configs, etc)
	singletons.Initialize(ctx)

	// Setup router
	router := routes.SetupRouter()

	// Start server
	log.Info(ctx, "Starting server on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(ctx, "Failed to start server", zap.Error(err))
	}
}
