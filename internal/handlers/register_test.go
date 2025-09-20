package handlers

import (
	"context"
	"net/http"
	"testing"

	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"github.com/steinfletcher/apitest"
)

func TestRegister(t *testing.T) {
	uni := services.NewAppUni()
	validate := services.NewValidate(uni)
	responser := services.NewResponser(uni)

	apitest.New("empty data").
		HandlerFunc(Register(context.Background(), responser, validate)).
		Post("/api/user/register").
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	apitest.New("invalid data").
		HandlerFunc(Register(context.Background(), responser, validate)).
		Post("/api/user/register").
		Body(`{
            "login": "admin"
        }`).
		Expect(t).
		Status(http.StatusBadRequest).
		Body(`{
		    "error": {
		        "message": "invalid data",
		        "invalid_fields": [
		            {
		                "field": "password",
		                "message": "password is a required field"
		            }
		        ]
		    }
		}`).
		End()
}
