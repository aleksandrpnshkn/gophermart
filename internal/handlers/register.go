package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/aleksandrpnshkn/gophermart/internal/handlers/requests"
	"github.com/aleksandrpnshkn/gophermart/internal/handlers/responses"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"github.com/go-playground/validator/v10"
)

func Register(
	ctx context.Context,
	responser *services.Responser,
	validate *validator.Validate,
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

		rawResponseData, _ := responses.EncodeOkResponse()

		res.WriteHeader(http.StatusOK)
		res.Write(rawResponseData)
	}
}
