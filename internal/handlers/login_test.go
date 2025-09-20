package handlers

import (
	"context"
	"net/http"
	"testing"

	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"github.com/steinfletcher/apitest"
)

func TestLogin(t *testing.T) {
	uni := services.NewAppUni()
	validate := services.NewValidate(uni)
	responser := services.NewResponser(uni)

	apitest.New("empty data").
		HandlerFunc(Login(context.Background(), responser, validate)).
		Post("/api/user/login").
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	apitest.New("invalid data").
		HandlerFunc(Login(context.Background(), responser, validate)).
		Post("/api/user/login").
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
