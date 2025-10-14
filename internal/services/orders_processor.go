package services

import (
	"context"
	"errors"

	"github.com/aleksandrpnshkn/gophermart/internal/models"
)

type GetAccrualProcessor struct {
	ordersService *OrdersService
}

func (p *GetAccrualProcessor) GetName() string {
	return "get_accrual"
}

func (p *GetAccrualProcessor) Process(
	ctx context.Context,
	order models.Order,
) (models.Order, error) {
	order, err := p.ordersService.UpdateAccrual(ctx, order)
	if err != nil {
		if errors.Is(err, ErrAccrualNotProcessedStatus) ||
			errors.Is(err, ErrAccrualFailedToGet) {
			return order, ErrJobRetry
		}

		return order, err
	}

	return order, nil
}

func NewOrdersProcessor(ordersService *OrdersService) *GetAccrualProcessor {
	return &GetAccrualProcessor{
		ordersService: ordersService,
	}
}
