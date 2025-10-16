package services

import (
	"context"
	"errors"

	"github.com/aleksandrpnshkn/gophermart/internal/models"
	balancePackage "github.com/aleksandrpnshkn/gophermart/internal/storage/balance"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var (
	ErrBalanceNotEnoughFunds = errors.New("not enough funds on user balance")
	ErrBalanceBadPrecision   = errors.New("amount contains more than two digits after the dot")
	ErrBalanceNegativeAmount = errors.New("cannot withdraw negative or zero amount")
)

type OrderAdder interface {
	Add(ctx context.Context, orderNumber string, user models.User) (models.Order, error)
}

type BalanceService struct {
	ordersAdder    OrderAdder
	balanceStorage balancePackage.Storage
	logger         *zap.Logger
}

func (b *BalanceService) Withdraw(
	ctx context.Context,
	orderNumber string,
	sumRaw float64,
	user models.User,
) error {
	if sumRaw <= 0 {
		return ErrBalanceNegativeAmount
	}

	sum := decimal.NewFromFloat(sumRaw)
	if !sum.Equal(sum.Truncate(2)) {
		return ErrBalanceBadPrecision
	}

	balance, err := b.GetBalance(ctx, user)
	if err != nil {
		return err
	}
	if balance.Current.LessThan(sum) {
		return ErrBalanceNotEnoughFunds
	}

	order, err := b.ordersAdder.Add(ctx, orderNumber, user)
	if err != nil && !errors.Is(err, ErrOrderAlreadyCreated) {
		b.logger.Error("failed to add order for withdrawal", zap.Error(err))
		return err
	}

	withdraw := models.BalanceChange{
		OrderNumber: order.OrderNumber,
		UserID:      order.UserID,
		Amount:      sum.Neg(),
	}
	err = b.balanceStorage.Withdraw(ctx, withdraw)
	if err != nil {
		if errors.Is(err, balancePackage.ErrNotEnoughFunds) {
			return ErrBalanceNotEnoughFunds
		}

		b.logger.Error("failed to add order for withdrawal", zap.Error(err))
		return err
	}

	return nil
}

func (b *BalanceService) GetBalance(
	ctx context.Context,
	user models.User,
) (models.Balance, error) {
	return b.balanceStorage.GetBalance(ctx, user)
}

func (b *BalanceService) GetWithdrawals(
	ctx context.Context,
	user models.User,
) ([]models.BalanceChange, error) {
	return b.balanceStorage.GetWithdrawals(ctx, user)
}

func NewBalancer(
	ordersService OrderAdder,
	balanceStorage balancePackage.Storage,
	logger *zap.Logger,
) *BalanceService {
	return &BalanceService{
		ordersAdder:    ordersService,
		balanceStorage: balanceStorage,
		logger:         logger,
	}
}
