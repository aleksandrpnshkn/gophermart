package app

import (
	"context"
	"net/http"

	"github.com/aleksandrpnshkn/gophermart/internal/config"
	"github.com/aleksandrpnshkn/gophermart/internal/handlers"
	"github.com/aleksandrpnshkn/gophermart/internal/middlewares"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"github.com/aleksandrpnshkn/gophermart/internal/storage/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func Run(
	ctx context.Context,
	config *config.Config,
	logger *zap.Logger,
	usersStorage users.Storage,
) error {
	router := chi.NewRouter()

	uni := services.NewAppUni()
	responser := services.NewResponser(uni)
	validate := services.NewValidate(uni)
	auther := services.NewAuther(usersStorage, config.JwtSecretKey)

	router.Use(middlewares.NewLogMiddleware(logger))
	router.Use(middleware.SetHeader("Content-Type", "application/json"))

	router.NotFound(handlers.NotFound(responser))

	router.Get("/api/ping", handlers.Ping())

	router.Post("/api/user/login", handlers.Login(ctx, responser, validate, auther, logger))
	router.Post("/api/user/register", handlers.Register(ctx, responser, validate, auther, logger))

	logger.Info("Running app...")

	err := http.ListenAndServe(config.RunAddress, router)

	return err
}
