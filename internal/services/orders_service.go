package services

import (
	"context"
	"errors"
	"time"

	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"github.com/aleksandrpnshkn/gophermart/internal/storage/orders"
	"github.com/aleksandrpnshkn/gophermart/internal/types"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type IOrdersService interface {
	Add(ctx context.Context, orderNumber string, user models.User) (models.Order, error)

	UpdateAccrual(ctx context.Context, order models.Order) (models.Order, error)

	GetUserOrders(ctx context.Context, user models.User) ([]models.Order, error)

	HasProcessedStatus(order models.Order) bool
}

type OrdersService struct {
	ordersStorage  orders.Storage
	accrualService IAccrualService
	logger         *zap.Logger
}

var (
	ErrOrderAlreadyCreated              = errors.New("order already created")
	ErrOrderAlreadyCreatedByAnotherUser = errors.New("order already created by another user")
)

func (o *OrdersService) Add(ctx context.Context, orderNumber string, user models.User) (models.Order, error) {
	order := models.Order{
		OrderNumber: orderNumber,
		UserID:      user.ID,
		Accrual:     decimal.NewFromInt(0),
		Status:      types.OrderStatusNew,
		UploadedAt:  time.Now(),
	}

	order, err := o.ordersStorage.Create(ctx, order)
	if err != nil {
		if errors.Is(err, orders.ErrOrderAlreadyCreated) {
			return order, ErrOrderAlreadyCreated
		}
		if errors.Is(err, orders.ErrOrderAlreadyCreatedByAnotherUser) {
			return models.Order{}, ErrOrderAlreadyCreatedByAnotherUser
		}
		return models.Order{}, err
	}

	return order, nil
}

func (o *OrdersService) UpdateAccrual(
	ctx context.Context,
	order models.Order,
) (models.Order, error) {
	if order.Status == types.OrderStatusNew {
		order.Status = types.OrderStatusProcessing
		err := o.ordersStorage.UpdateStatus(ctx, order)
		if err != nil {
			o.logger.Error("failed to set processing status",
				zap.String("order_number", order.OrderNumber),
				zap.Error(err),
			)
			return order, err
		}
	}

	if order.Status != types.OrderStatusProcessing {
		o.logger.Error("tried to update unexpected order status",
			zap.String("order_number", order.OrderNumber),
			zap.String("order_status", string(order.Status)),
		)
		return order, errors.New("tried to update unexpected order status")
	}

	if order.Accrual.IsZero() {
		accrual, err := o.accrualService.GetAccrual(ctx, order.OrderNumber)
		if err != nil && !errors.Is(err, ErrAccrualInvalidStatus) {
			o.logger.Error("failed to get accrual",
				zap.String("order_number", order.OrderNumber),
				zap.Error(err),
			)
			return order, err
		}

		if errors.Is(err, ErrAccrualInvalidStatus) {
			order.Status = types.OrderStatusInvalid
		} else {
			order.Accrual = accrual
			order.Status = types.OrderStatusProcessed
		}
	}

	if order.Accrual.IsZero() {
		err := o.ordersStorage.UpdateStatus(ctx, order)
		if err != nil {
			o.logger.Error("failed to set processed status",
				zap.String("order_number", order.OrderNumber),
				zap.String("order_status", string(order.Status)),
				zap.Error(err),
			)
			return order, err
		}
	} else {
		err := o.ordersStorage.UpdateAccrual(ctx, order)
		if err != nil {
			o.logger.Error("failed to update accrual",
				zap.String("order_number", order.OrderNumber),
				zap.String("order_status", string(order.Status)),
				zap.String("accrual", order.Accrual.String()),
				zap.Error(err),
			)
			return order, err
		}
	}

	return order, nil
}

func (o *OrdersService) GetUserOrders(
	ctx context.Context,
	user models.User,
) ([]models.Order, error) {
	orders, err := o.ordersStorage.GetUserOrders(ctx, user)
	if err != nil {
		return []models.Order{}, err
	}

	return orders, nil
}

func (o *OrdersService) HasProcessedStatus(order models.Order) bool {
	return order.Status == types.OrderStatusProcessed ||
		order.Status == types.OrderStatusInvalid
}

func NewOrdersService(
	ordersStorage orders.Storage,
	accrualService IAccrualService,
	logger *zap.Logger,
) *OrdersService {
	return &OrdersService{
		ordersStorage:  ordersStorage,
		accrualService: accrualService,
		logger:         logger,
	}
}
