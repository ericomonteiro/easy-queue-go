package singletons

import (
	"context"
	"easy-queue-go/src/internal/config"
	"easy-queue-go/src/internal/infra"
	"easy-queue-go/src/internal/infra/database"
	"easy-queue-go/src/internal/log"

	"go.uber.org/zap"
)

type Instances struct {
	DB infra.IDataBase
}

func Initialize(ctx context.Context) *Instances {
	instances := new(Instances)

	configs := loadAllConfigs(ctx)

	loadAllInfra(ctx, instances, configs)

	return instances
}

func loadAllConfigs(ctx context.Context) *config.Configs {
	log.Info(ctx, "Loading configs")
	return config.InitializeConfigs()
}

func loadAllInfra(ctx context.Context, instances *Instances, configs *config.Configs) {
	db, err := database.NewClient(ctx, configs.DB)
	if err != nil {
		log.Panic(ctx, "Failed to connect to database", zap.Error(err))
	}
	instances.DB = db
}
