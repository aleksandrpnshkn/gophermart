package handlers

import (
	"net/http"

	"github.com/aleksandrpnshkn/gophermart/internal/handlers/responses"
)

func NotFound() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		rawResponseData, _ := responses.EncodeNotFoundResponse()

		res.WriteHeader(http.StatusNotFound)
		res.Write(rawResponseData)
	}
}
