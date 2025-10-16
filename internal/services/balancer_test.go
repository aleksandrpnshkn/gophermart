package services

import (
	"context"
	"testing"

	"github.com/aleksandrpnshkn/gophermart/internal/mocks"
	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"github.com/aleksandrpnshkn/gophermart/internal/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestBalancer(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := zap.NewExample()

	user := models.User{
		ID:    1,
		Login: "admin",
		Hash:  types.PasswordHash("hash"),
	}

	t.Run("invalid amount", func(t *testing.T) {
		balanceStorage := mocks.NewMockBalanceStorage(ctrl)
		ordersService := mocks.NewMockOrdersService(ctrl)

		balancer := NewBalancer(ordersService, balanceStorage, logger)

		err := balancer.Withdraw(context.Background(), "123", 123.123, user)

		assert.ErrorIs(t, err, ErrBalanceBadPrecision)
	})

	t.Run("negative amount", func(t *testing.T) {
		balanceStorage := mocks.NewMockBalanceStorage(ctrl)
		ordersService := mocks.NewMockOrdersService(ctrl)

		balancer := NewBalancer(ordersService, balanceStorage, logger)

		err := balancer.Withdraw(context.Background(), "123", -123, user)

		assert.ErrorIs(t, err, ErrBalanceNegativeAmount)
	})
}
