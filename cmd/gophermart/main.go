package main

import (
	"context"
	"log"
	"os"

	"github.com/aleksandrpnshkn/gophermart/internal/app"
	"github.com/aleksandrpnshkn/gophermart/internal/config"
	"github.com/aleksandrpnshkn/gophermart/internal/logs"
	"github.com/aleksandrpnshkn/gophermart/internal/storage"
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

	usersStorage, err := storage.NewUsersStorage(ctx, config.DatabaseURI, logger)
	if err != nil {
		logger.Fatal("failed to init users storage", zap.Error(err))
	}
	defer usersStorage.Close()

	err = app.Run(ctx, config, logger, usersStorage)
	if err != nil {
		logger.Fatal("failed to run app", zap.Error(err))
	}
}
