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
	ordersService *services.OrdersService,
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

		err = ordersService.Add(ctx, orderNumber, user)
		if err != nil {
			if errors.Is(err, services.ErrOrderAlreadyCreated) {
				responser.WriteSuccess(ctx, res)
				return
			}
			if errors.Is(err, services.ErrOrderAlreadyCreatedByAnotherUser) {
				responser.WriteConflict(ctx, res)
				return
			}

			logger.Error("failed to create order", zap.Error(err))
			responser.WriteInternalServerError(ctx, res)
			return
		}

		// добавить в воркер

		// Статусы заказов должны быть персистентными.
		// Неполная асинхронщина, надо проверить заказы в БД.
		// Чтобы проверка была успешной надо либо сразу же складывать заказы в БД, либо держать кэш принятых заказов и батчами их грузить в БД.
		// Можно сразу сделать батчинг на входе в БД, чтобы не дудосить БД.
		// Батчинг должен быть универсальным и для сохранения в БД, и для обновления.
		// В батчинге не должно быть одного и того же заказа в разных состояниях, чтобы не приходилось разбираться с гонками.
		// Реализация очереди может быть одна, но самих очередей две:
		// updateOrderQueue := queues.New() // сюда по идее могут поступать заказы на любом этапе
		// checkOrderQueue := queues.New()

		// добавить заказ в очередь на обработку

		rawResponseData, _ := responses.EncodeOkResponse()

		res.WriteHeader(http.StatusAccepted)
		res.Write(rawResponseData)
	}
}
