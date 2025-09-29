package app

import (
	"context"
	"net/http"

	"github.com/aleksandrpnshkn/gophermart/internal/config"
	"github.com/aleksandrpnshkn/gophermart/internal/handlers"
	"github.com/aleksandrpnshkn/gophermart/internal/middlewares"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"github.com/aleksandrpnshkn/gophermart/internal/storage/orders"
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
	ordersStorage orders.Storage,
) error {
	router := chi.NewRouter()

	uni := services.NewAppUni()
	responser := services.NewResponser(uni)
	validate := services.NewValidate(uni)
	auther := services.NewAuther(usersStorage, config.JwtSecretKey)

	accrualService := services.NewAccrualService(logger, config.AccrualSystemAddress)
	ordersService := services.NewOrdersService(ordersStorage, accrualService, logger)

	ordersProceessor := services.NewOrdersProcessor(ordersService)
	ordersQueue := services.NewOrdersQueue(
		ctx,
		ordersProceessor,
		logger,
	)

	router.Use(middlewares.NewLogMiddleware(logger))
	router.Use(middleware.SetHeader("Content-Type", "application/json"))

	router.NotFound(handlers.NotFound(responser))

	router.Get("/api/ping", handlers.Ping())

	router.Post("/api/user/login", handlers.Login(ctx, responser, validate, auther, logger))
	router.Post("/api/user/register", handlers.Register(ctx, responser, validate, auther, logger))

	router.Group(func(router chi.Router) {
		router.Use(middlewares.NewAuthMiddleware(logger, auther))

		router.Post("/api/user/orders", handlers.AddOrder(ctx, responser, auther, logger, ordersService, ordersQueue))
	})

	logger.Info("running app...")

	err := http.ListenAndServe(config.RunAddress, router)

	return err
}
