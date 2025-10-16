package handlers

import (
	"net/http"

	"github.com/aleksandrpnshkn/gophermart/internal/services"
)

func NotFound(responser *services.Responser) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		responser.WriteNotFoundError(req.Context(), res)
	}
}
