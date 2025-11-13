package main

import (
	"context"
	"easy-queue-go/src/internal/log"
	"easy-queue-go/src/internal/singletons"
)

func main() {
	// Initialize zap logger
	ctx := log.Initialize(context.Background())

	singletons.Initialize(ctx)
}
