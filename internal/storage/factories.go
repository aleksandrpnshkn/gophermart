package storage

import (
	"context"
	"fmt"

	"github.com/aleksandrpnshkn/gophermart/internal/storage/orders"
	"github.com/aleksandrpnshkn/gophermart/internal/storage/users"
	"go.uber.org/zap"
)

func NewUsersStorage(
	ctx context.Context,
	databaseDSN string,
	logger *zap.Logger,
) (users.Storage, error) {
	var storage users.Storage

	storage, err := users.NewSQLStorage(ctx, databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to init users SQL storage: %w", err)
	}

	return storage, nil
}

func NewOrdersStorage(
	ctx context.Context,
	databaseDSN string,
	logger *zap.Logger,
) (orders.Storage, error) {
	var storage orders.Storage

	storage, err := orders.NewSQLStorage(ctx, databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to init users SQL storage: %w", err)
	}

	return storage, nil
}
