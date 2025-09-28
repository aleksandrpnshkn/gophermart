package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aleksandrpnshkn/gophermart/internal/app"
	"github.com/aleksandrpnshkn/gophermart/internal/config"
	"github.com/aleksandrpnshkn/gophermart/internal/logs"
	"github.com/aleksandrpnshkn/gophermart/internal/storage"
	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"
)

func main() {
	config := config.New()
	ctx := context.Background()

	logger, err := logs.NewLogger(config.LogLevel)
	if err != nil {
		log.Printf("failed to create app logger: %v", err)
		os.Exit(1)
	}
	defer logger.Sync()

	err = storage.RunMigrations(config.DatabaseURI)
	if err == nil {
		logger.Info("successfully migrated")
	} else if errors.Is(err, migrate.ErrNoChange) {
		logger.Info("nothing to migrate")
	} else {
		logger.Fatal("failed to run SQL migrations", zap.Error(err))
	}

	usersStorage, err := storage.NewUsersStorage(ctx, config.DatabaseURI, logger)
	if err != nil {
		logger.Fatal("failed to init users storage", zap.Error(err))
	}
	defer usersStorage.Close()

	ordersStorage, err := storage.NewOrdersStorage(ctx, config.DatabaseURI, logger)
	if err != nil {
		logger.Fatal("failed to init orders storage", zap.Error(err))
	}
	defer ordersStorage.Close()

	err = app.Run(ctx, config, logger, usersStorage, ordersStorage)
	if err != nil {
		logger.Fatal("failed to run app", zap.Error(err))
	}
}
