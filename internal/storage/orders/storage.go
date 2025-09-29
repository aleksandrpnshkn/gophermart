package orders

import (
	"context"
	"errors"

	"github.com/aleksandrpnshkn/gophermart/internal/models"
)

type Storage interface {
	Ping(ctx context.Context) error

	GetByNumber(ctx context.Context, orderNumber string) (models.Order, error)

	GetUserOrders(ctx context.Context, user models.User) ([]models.Order, error)

	Create(ctx context.Context, order models.Order) error

	Update(ctx context.Context, order models.Order) error

	Close() error
}

var (
	ErrOrderNotFound                    = errors.New("order not found")
	ErrOrderAlreadyCreated              = errors.New("order already created")
	ErrOrderAlreadyCreatedByAnotherUser = errors.New("order already created by another user")
)
