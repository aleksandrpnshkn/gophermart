package balance

import (
	"context"

	"github.com/aleksandrpnshkn/gophermart/internal/models"
)

type Storage interface {
	Ping(ctx context.Context) error

	Withdraw(ctx context.Context, withdraw models.BalanceChange) error

	GetBalance(ctx context.Context, user models.User) (models.Balance, error)

	GetWithdrawals(ctx context.Context, user models.User) ([]models.BalanceChange, error)

	Close() error
}
