package services

import (
	"context"
	"net/http"

	"github.com/aleksandrpnshkn/gophermart/internal/responses"
	"github.com/go-playground/validator/v10"
)

type Responser struct {
	uni *AppUni
}

func (r *Responser) WriteSuccess(ctx context.Context, res http.ResponseWriter) {
	res.WriteHeader(http.StatusOK)

	rawResponseData, err := responses.EncodeOkResponse()
	if err == nil {
		res.Write(rawResponseData)
	}
}

func (r *Responser) WriteNoContent(ctx context.Context, res http.ResponseWriter) {
	res.WriteHeader(http.StatusNoContent)
}

func (r *Responser) WriteConflict(ctx context.Context, res http.ResponseWriter) {
	r.writeError(ctx, res, http.StatusConflict, "conflict")
}

func (r *Responser) WriteUnauthorizedError(ctx context.Context, res http.ResponseWriter) {
	r.writeError(ctx, res, http.StatusUnauthorized, "unauthorized")
}

func (r *Responser) WriteNotFoundError(ctx context.Context, res http.ResponseWriter) {
	r.writeError(ctx, res, http.StatusNotFound, "not found")
}

func (r *Responser) WriteBadRequestError(ctx context.Context, res http.ResponseWriter) {
	r.writeError(ctx, res, http.StatusBadRequest, "bad request")
}

func (r *Responser) WriteEmptyValidationError(ctx context.Context, res http.ResponseWriter) {
	res.WriteHeader(http.StatusUnprocessableEntity)
	errors := make(map[string]string)

	rawResponseData, err := responses.EncodeValidationErrorResponse("invalid data", errors)
	if err == nil {
		res.Write(rawResponseData)
	}
}

func (r *Responser) WriteValidationError(ctx context.Context, res http.ResponseWriter, err error) {
	trans := r.uni.ResolveUserTrans(ctx)

	validationErrors := err.(validator.ValidationErrors)
	errors := make(map[string]string, len(validationErrors))
	for _, e := range validationErrors {
		errors[e.Field()] = e.Translate(trans)
	}

	res.WriteHeader(http.StatusUnprocessableEntity)

	rawResponseData, err := responses.EncodeValidationErrorResponse("invalid data", errors)
	if err == nil {
		res.Write(rawResponseData)
	}
}

func (r *Responser) WriteInternalServerError(ctx context.Context, res http.ResponseWriter) {
	r.writeError(ctx, res, http.StatusInternalServerError, "internal server error")
}

func (r *Responser) writeError(ctx context.Context, res http.ResponseWriter, status int, message string) {
	res.WriteHeader(status)

	rawResponseData, err := responses.EncodeErrorResponse(message)
	if err == nil {
		res.Write(rawResponseData)
	}
}

func NewResponser(uni *AppUni) *Responser {
	return &Responser{
		uni: uni,
	}
}
