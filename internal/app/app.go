package app

import (
	"context"
	"net/http"

	"github.com/aleksandrpnshkn/gophermart/internal/config"
	"github.com/aleksandrpnshkn/gophermart/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func Run(
	ctx context.Context,
	config *config.Config,
	logger *zap.Logger,
) error {
	router := chi.NewRouter()

	router.Use(middleware.SetHeader("Content-Type", "application/json"))

	router.NotFound(handlers.NotFound())

	router.Get("/api/ping", handlers.Ping())

	logger.Info("Running app...")

	err := http.ListenAndServe(config.RunAddress, router)

	return err
}
