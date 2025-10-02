package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/aleksandrpnshkn/gophermart/internal/responses"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"go.uber.org/zap"
)

func AddOrder(
	responser *services.Responser,
	auther services.Auther,
	logger *zap.Logger,
	ordersService services.IOrdersService,
	ordersQueue services.OrdersQueue,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		user, err := auther.FromUserContext(ctx)
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
