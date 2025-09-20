package handlers

import (
	"net/http"

	"github.com/aleksandrpnshkn/gophermart/internal/handlers/responses"
)

func Ping() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		rawResponseData, _ := responses.EncodeOkResponse()

		res.WriteHeader(http.StatusOK)
		res.Write(rawResponseData)
	}
}
