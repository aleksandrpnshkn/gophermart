package app

import (
	"context"
	"net/http"

	"github.com/aleksandrpnshkn/gophermart/internal/config"
	"github.com/aleksandrpnshkn/gophermart/internal/handlers"
	"github.com/aleksandrpnshkn/gophermart/internal/middlewares"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"github.com/aleksandrpnshkn/gophermart/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func Run(
	ctx context.Context,
	config *config.Config,
	logger *zap.Logger,
	storages *storage.Storages,
) error {
	router := chi.NewRouter()

	uni := services.NewAppUni()
	responser := services.NewResponser(uni)
	validate := services.NewValidate(uni)
	auther := services.NewAuther(storages.Users, config.JwtSecretKey)

	accrualClient := http.Client{}
	accrualService := services.NewAccrualService(&accrualClient, logger, config.AccrualSystemAddress)
	ordersService := services.NewOrdersService(storages.Orders, accrualService, logger)

	ordersProceessor := services.NewOrdersProcessor(ordersService)
	ordersQueue := services.NewOrdersQueue(
		ctx,
		ordersProceessor,
		logger,
	)

	balancer := services.NewBalancer(ordersService, storages.Balance, logger)

	router.Use(middlewares.NewLogMiddleware(logger))
	router.Use(middleware.SetHeader("Content-Type", "application/json"))

	router.NotFound(handlers.NotFound(responser))

	router.Get("/api/ping", handlers.Ping())

	router.Post("/api/user/login", handlers.Login(responser, validate, auther, logger))
	router.Post("/api/user/register", handlers.Register(responser, validate, auther, logger))

	router.Group(func(router chi.Router) {
		router.Use(middlewares.NewAuthMiddleware(responser, logger, auther))

		router.Post("/api/user/orders", handlers.AddOrder(responser, auther, logger, ordersService, ordersQueue))
		router.Get("/api/user/orders", handlers.GetUserOrders(responser, auther, logger, ordersService))

		router.Get("/api/user/balance", handlers.GetBalance(responser, auther, balancer, logger))
		router.Post("/api/user/balance/withdraw", handlers.Withdraw(responser, validate, auther, balancer, logger))
		router.Get("/api/user/withdrawals", handlers.GetWithdrawals(responser, auther, balancer, logger))
	})

	logger.Info("running app...")

	err := http.ListenAndServe(config.RunAddress, router)

	return err
}
