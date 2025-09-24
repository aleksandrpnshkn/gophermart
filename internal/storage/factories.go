package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/aleksandrpnshkn/gophermart/internal/storage/users"
	"github.com/golang-migrate/migrate/v4"
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

	err = runMigrations(databaseDSN)
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("nothing to migrate")
			return storage, nil
		} else {
			return nil, fmt.Errorf("failed to run SQL migrations: %w", err)
		}
	}

	logger.Info("database successfully migrated")
	return storage, nil
}
