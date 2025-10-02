package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/aleksandrpnshkn/gophermart/internal/requests"
	"github.com/aleksandrpnshkn/gophermart/internal/responses"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func Withdraw(
	responser *services.Responser,
	validate *validator.Validate,
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

		rawRequestData, err := io.ReadAll(req.Body)
		if err != nil {
			responser.WriteBadRequestError(ctx, res)
			return
		}
		defer req.Body.Close()

		var requestData requests.Withdraw
		err = json.Unmarshal(rawRequestData, &requestData)
		if err != nil {
			responser.WriteBadRequestError(ctx, res)
			return
		}

		err = validate.StructCtx(ctx, requestData)
		if err != nil {
			responser.WriteValidationError(ctx, res, err)
			return
		}

		err = balancer.Withdraw(ctx, requestData.OrderNumber, requestData.Amount, user)
		if err != nil {
			if errors.Is(err, services.ErrBalanceNotEnoughFunds) {
				res.WriteHeader(http.StatusPaymentRequired)
				return
			}

			logger.Error("failed to withdraw", zap.Error(err))
			responser.WriteInternalServerError(ctx, res)
			return
		}

		res.WriteHeader(http.StatusOK)

		rawResponseData, _ := responses.EncodeOkResponse()
		res.Write(rawResponseData)
	}
}
