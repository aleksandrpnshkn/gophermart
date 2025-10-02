package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aleksandrpnshkn/gophermart/internal/app"
	"github.com/aleksandrpnshkn/gophermart/internal/config"
	"github.com/aleksandrpnshkn/gophermart/internal/logs"
	"go.uber.org/zap"
)

func main() {
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config := config.New()

	logger, err := logs.NewLogger(config.LogLevel)
	if err != nil {
		log.Printf("failed to create app logger: %v", err)
		os.Exit(1)
	}
	defer logger.Sync()

	err = app.Run(rootCtx, config, logger)
	if err != nil {
		logger.Fatal("failed to run app", zap.Error(err))
	}
}
