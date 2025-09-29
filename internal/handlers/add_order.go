package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"regexp"

	"github.com/aleksandrpnshkn/gophermart/internal/responses"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"go.uber.org/zap"
)

const orderNumberRegexPattern = "^[0-9]+$"

func AddOrder(
	ctx context.Context,
	responser *services.Responser,
	auther services.Auther,
	logger *zap.Logger,
	ordersService services.IOrdersService,
	ordersQueue services.OrdersQueue,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		user, err := auther.FromUserContext(req.Context())
		if err != nil {
			logger.Error("failed to get user", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		rawOrderNumber, err := io.ReadAll(req.Body)
		if err != nil {
			responser.WriteBadRequestError(ctx, res)
			return
		}
		defer req.Body.Close()

		isValidOrderNumber, err := regexp.Match(orderNumberRegexPattern, rawOrderNumber)
		if err != nil {
			logger.Error("failed to validate order number", zap.Error(err))
			responser.WriteInternalServerError(ctx, res)
			return
		}
		if !isValidOrderNumber {
			responser.WriteEmptyValidationError(ctx, res)
			return
		}
		orderNumber := string(rawOrderNumber)

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
