package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/aleksandrpnshkn/gophermart/internal/middlewares"
	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"github.com/aleksandrpnshkn/gophermart/internal/requests"
	"github.com/aleksandrpnshkn/gophermart/internal/responses"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"github.com/aleksandrpnshkn/gophermart/internal/types"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type UserLoginer interface {
	LoginUser(ctx context.Context, login string, password string) (models.User, types.RawToken, error)
}

func Login(
	responser *services.Responser,
	validate *validator.Validate,
	userLoginer UserLoginer,
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

		_, rawToken, err := userLoginer.LoginUser(ctx, requestData.Login, requestData.Password)
		if err != nil {
			if errors.Is(err, services.ErrBadCredentials) {
				responser.WriteUnauthorizedError(ctx, res)
				return
			}

			logger.Error("failed to login user", zap.Error(err))
			responser.WriteInternalServerError(ctx, res)
			return
		}

		middlewares.SetAuthCookie(res, rawToken)

		rawResponseData, _ := responses.EncodeOkResponse()

		res.WriteHeader(http.StatusOK)
		res.Write(rawResponseData)
	}
}
