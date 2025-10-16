package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aleksandrpnshkn/gophermart/internal/responses"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"go.uber.org/zap"
)

func GetWithdrawals(
	responser *services.Responser,
	userReceiver UserReceiver,
	withdrawer Withdrawer,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		user, err := userReceiver.FromContext(ctx)
		if err != nil {
			logger.Error("failed to get user", zap.Error(err))
			responser.WriteInternalServerError(ctx, res)
			return
		}

		withdrawals, err := withdrawer.GetWithdrawals(ctx, user)
		if err != nil {
			logger.Error("failed to get user withdrawals", zap.Error(err))
			responser.WriteInternalServerError(ctx, res)
			return
		}

		responseData := []responses.Withdraw{}
		for _, balanceChange := range withdrawals {
			responseData = append(responseData, responses.Withdraw{
				OrderNumber: balanceChange.OrderNumber,
				Sum:         balanceChange.Amount.Abs().InexactFloat64(),
				ProcessedAt: balanceChange.ProcessedAt.Format(time.RFC3339),
			})
		}

		rawResponseData, err := json.Marshal(responseData)
		if err != nil {
			logger.Error("failed to marshal user balance",
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
