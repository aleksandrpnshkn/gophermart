package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aleksandrpnshkn/gophermart/internal/responses"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"go.uber.org/zap"
)

func GetUserOrders(
	responser *services.Responser,
	userReceiver UserReceiver,
	logger *zap.Logger,
	ordersService OrdersService,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user, err := userReceiver.FromContext(ctx)
		if err != nil {
			logger.Error("failed to get user", zap.Error(err))
			responser.WriteInternalServerError(ctx, res)
			return
		}

		orders, err := ordersService.GetUserOrders(ctx, user)
		if err != nil {
			logger.Error("failed to get user orders",
				zap.Int64("user_id", user.ID),
				zap.Error(err),
			)
			responser.WriteInternalServerError(ctx, res)
			return
		}

		if len(orders) == 0 {
			responser.WriteNoContent(ctx, res)
			return
		}

		responseData := []responses.Order{}

		for _, order := range orders {
			orderData := responses.Order{
				OrderNumber: order.OrderNumber,
				Status:      string(order.Status),
				UploadedAt:  order.UploadedAt.Format(time.RFC3339),
			}

			if ordersService.HasProcessedStatus(order) {
				orderData.Accrual, _ = order.Accrual.Float64()
			}

			responseData = append(responseData, orderData)
		}

		rawResponseData, err := json.Marshal(responseData)
		if err != nil {
			logger.Error("failed to marshal user orders",
				zap.Int64("user_id", user.ID),
				zap.Error(err),
			)
			responser.WriteInternalServerError(ctx, res)
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Write(rawResponseData)
	}
}
