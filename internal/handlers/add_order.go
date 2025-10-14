package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"github.com/aleksandrpnshkn/gophermart/internal/responses"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"go.uber.org/zap"
)

type UserReceiver interface {
	FromContext(ctx context.Context) (models.User, error)
}

type OrdersQueue interface {
	Add(ctx context.Context, order models.Order) error

	Stop()
}

type OrdersService interface {
	Add(ctx context.Context, orderNumber string, user models.User) (models.Order, error)

	UpdateAccrual(ctx context.Context, order models.Order) (models.Order, error)

	GetUserOrders(ctx context.Context, user models.User) ([]models.Order, error)

	HasProcessedStatus(order models.Order) bool
}

func AddOrder(
	responser *services.Responser,
	userReciever UserReceiver,
	logger *zap.Logger,
	ordersService OrdersService,
	ordersQueue OrdersQueue,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		user, err := userReciever.FromContext(ctx)
		if err != nil {
			logger.Error("failed to get user", zap.Error(err))
			responser.WriteInternalServerError(ctx, res)
			return
		}

		rawOrderNumber, err := io.ReadAll(req.Body)
		if err != nil {
			responser.WriteBadRequestError(ctx, res)
			return
		}
		defer req.Body.Close()

		orderNumber := string(rawOrderNumber)
		if !services.IsValidLuhnNumber(orderNumber) {
			responser.WriteEmptyValidationError(ctx, res)
			return
		}

		order, err := ordersService.Add(ctx, orderNumber, user)
		if err != nil {
			if errors.Is(err, services.ErrOrderAlreadyCreated) {
				responser.WriteSuccess(ctx, res)
				return
			}
			if errors.Is(err, services.ErrOrderAlreadyCreatedByAnotherUser) {
				responser.WriteConflict(ctx, res)
				return
			}

			logger.Error("failed to create order",
				zap.String("order_number", orderNumber),
				zap.Error(err),
			)
			responser.WriteInternalServerError(ctx, res)
			return
		}

		err = ordersQueue.Add(ctx, order)
		if err != nil {
			logger.Error("failed to add order to queue",
				zap.String("order_number", orderNumber),
				zap.Error(err),
			)
			responser.WriteInternalServerError(ctx, res)
			return
		}

		rawResponseData, _ := responses.EncodeOkResponse()

		res.WriteHeader(http.StatusAccepted)
		res.Write(rawResponseData)
	}
}
