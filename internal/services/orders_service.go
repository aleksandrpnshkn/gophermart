package services

import (
	"context"
	"errors"
	"time"

	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"github.com/aleksandrpnshkn/gophermart/internal/storage/orders"
	"github.com/aleksandrpnshkn/gophermart/internal/types"
)

type OrdersService struct {
	ordersStorage orders.Storage
}

var (
	ErrOrderAlreadyCreated              = errors.New("order already created")
	ErrOrderAlreadyCreatedByAnotherUser = errors.New("order already created by another user")
)

func (o *OrdersService) Add(ctx context.Context, orderNumber string, user models.User) error {
	order := models.Order{
		OrderNumber: orderNumber,
		UserID:      user.ID,
		Accrual:     0,
		Status:      types.OrderStatusNew,
		UploadedAt:  time.Now(),
	}

	err := o.ordersStorage.Create(ctx, order)
	if err != nil {
		if errors.Is(err, orders.ErrOrderAlreadyCreated) {
			return ErrOrderAlreadyCreated
		}
		if errors.Is(err, orders.ErrOrderAlreadyCreatedByAnotherUser) {
			return ErrOrderAlreadyCreatedByAnotherUser
		}
		return err
	}

	return nil
}

func NewOrdersService(ordersStorage orders.Storage) *OrdersService {
	return &OrdersService{
		ordersStorage: ordersStorage,
	}
}
