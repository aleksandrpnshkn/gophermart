package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/aleksandrpnshkn/gophermart/internal/middlewares"
	"github.com/aleksandrpnshkn/gophermart/internal/requests"
	"github.com/aleksandrpnshkn/gophermart/internal/responses"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func Register(
	ctx context.Context,
	responser *services.Responser,
	validate *validator.Validate,
	auther services.Auther,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		rawRequestData, err := io.ReadAll(req.Body)
		if err != nil {
			responser.WriteBadRequestError(ctx, res)
			return
		}
		defer req.Body.Close()

		var requestData requests.Register
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

		_, token, err := auther.RegisterUser(ctx, requestData.Login, requestData.Password)
		if err != nil {
			if errors.Is(err, services.ErrLoginAlreadyExists) {
				res.WriteHeader(http.StatusConflict)
				return
			}

			logger.Error("failed to register user", zap.Error(err))
			responser.WriteInternalServerError(ctx, res)
			return
		}

		middlewares.SetAuthCookie(res, token)

		res.WriteHeader(http.StatusOK)

		rawResponseData, _ := responses.EncodeOkResponse()
		res.Write(rawResponseData)
	}
}
