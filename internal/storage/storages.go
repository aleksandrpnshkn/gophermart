package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/aleksandrpnshkn/gophermart/internal/storage/balance"
	"github.com/aleksandrpnshkn/gophermart/internal/storage/orders"
	"github.com/aleksandrpnshkn/gophermart/internal/storage/users"
	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"
)

type Storages struct {
	Orders  orders.Storage
	Users   users.Storage
	Balance balance.Storage
}

func (s *Storages) Close() error {
	err := s.Orders.Close()
	if err != nil {
		return err
	}

	err = s.Users.Close()
	if err != nil {
		return err
	}

	err = s.Balance.Close()
	if err != nil {
		return err
	}

	return nil
}

func NewStorages(
	ctx context.Context,
	databaseDSN string,
	logger *zap.Logger,
) (*Storages, error) {
	err := RunMigrations(databaseDSN)
	if err == nil {
		logger.Info("successfully migrated")
	} else if errors.Is(err, migrate.ErrNoChange) {
		logger.Info("nothing to migrate")
	} else {
		return nil, fmt.Errorf("failed to run SQL migrations: %w", err)
	}

	ordersStorage, err := orders.NewSQLStorage(ctx, databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to init orders SQL storage: %w", err)
	}

	usersStorage, err := users.NewSQLStorage(ctx, databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to init users SQL storage: %w", err)
	}

	balanceStorage, err := balance.NewSQLStorage(ctx, databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to init balance SQL storage: %w", err)
	}

	return &Storages{
		Orders:  ordersStorage,
		Users:   usersStorage,
		Balance: balanceStorage,
	}, nil
}
