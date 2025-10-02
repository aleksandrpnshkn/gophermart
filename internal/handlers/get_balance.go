package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aleksandrpnshkn/gophermart/internal/responses"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"go.uber.org/zap"
)

func GetBalance(
	responser *services.Responser,
	auther services.Auther,
	balancer services.IBalancer,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		user, err := auther.FromUserContext(ctx)
		if err != nil {
			logger.Error("failed to get user", zap.Error(err))
			responser.WriteInternalServerError(ctx, res)
			return
		}

		balance, err := balancer.GetBalance(ctx, user)
		if err != nil {
			logger.Error("failed to get balance", zap.Error(err))
			responser.WriteInternalServerError(ctx, res)
			return
		}

		responseData := responses.Balance{
			Current:   balance.Current.InexactFloat64(),
			Withdrawn: balance.Withdrawn.InexactFloat64(),
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
