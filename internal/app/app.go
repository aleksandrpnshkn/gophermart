package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/aleksandrpnshkn/gophermart/internal/config"
	"github.com/aleksandrpnshkn/gophermart/internal/handlers"
	"github.com/aleksandrpnshkn/gophermart/internal/middlewares"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"github.com/aleksandrpnshkn/gophermart/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

const (
	orderQueueWorkersCount = 3
	orderQueueJobsDelay    = 10 * time.Second

	shutdownTimeout = 20 * time.Second
)

func Run(
	rootCtx context.Context,
	config *config.Config,
	logger *zap.Logger,
) error {
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	storages, err := storage.NewStorages(appCtx, config.DatabaseURI, logger)
	if err != nil {
		return err
	}
	defer storages.Close()

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
		ordersProceessor,
		logger,
		orderQueueJobsDelay,
	)
	services.RunQueue(appCtx, ordersQueue, orderQueueWorkersCount)

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

	server := http.Server{
		Addr:    config.RunAddress,
		Handler: router,
		BaseContext: func(l net.Listener) context.Context {
			return appCtx
		},
	}

	go func() {
		logger.Info("server listening...",
			zap.String("addr", config.RunAddress),
		)

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Error("failed to run server", zap.Error(err))
			appCancel()
		}
	}()

	select {
	case <-appCtx.Done():
		return errors.New("app context prematurely cancelled")
	case <-rootCtx.Done():
		logger.Info("received shutdown signal, shutting down...")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()
	err = server.Shutdown(shutdownCtx)
	if err != nil {
		logger.Error("failed to shutdown server", zap.Error(err))
		return err
	}

	logger.Info("canceling app context manually...")
	appCancel()

	logger.Info("closing storages manually...")
	err = storages.Close()
	if err != nil {
		logger.Error("failed to close storages", zap.Error(err))
		return err
	}

	logger.Info("stopping orders queue...")
	ordersQueue.Stop()

	return nil
}
