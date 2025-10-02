package handlers

import (
	"net/http"
	"testing"

	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"github.com/steinfletcher/apitest"
)

func TestNotFound(t *testing.T) {
	responser := services.NewResponser(services.NewAppUni())

	apitest.New().
		HandlerFunc(NotFound(responser)).
		Get("/404").
		Expect(t).
		Body(`{"error": {"message":"not found"}}`).
		Status(http.StatusNotFound).
		End()
}
