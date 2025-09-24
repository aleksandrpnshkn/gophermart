package handlers

import (
	"context"
	"net/http"
	"testing"

	"github.com/aleksandrpnshkn/gophermart/internal/middlewares"
	"github.com/aleksandrpnshkn/gophermart/internal/mocks"
	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"github.com/aleksandrpnshkn/gophermart/internal/types"
	"github.com/steinfletcher/apitest"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uni := services.NewAppUni()
	validate := services.NewValidate(uni)
	responser := services.NewResponser(uni)
	logger := zap.NewExample()

	t.Run("empty data", func(t *testing.T) {
		auther := mocks.NewMockAuther(ctrl)

		handler := Register(context.Background(), responser, validate, auther, logger)

		apitest.New().
			HandlerFunc(handler).
			Post("/api/user/register").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("invalid data", func(t *testing.T) {
		auther := mocks.NewMockAuther(ctrl)

		handler := Register(context.Background(), responser, validate, auther, logger)

		apitest.New().
			HandlerFunc(handler).
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
	})

	t.Run("user registered", func(t *testing.T) {
		user := models.User{
			ID:    1,
			Login: "admin",
			Hash:  types.PasswordHash("blablahash"),
		}
		rawToken := types.RawToken("token")

		auther := mocks.NewMockAuther(ctrl)
		auther.EXPECT().
			RegisterUser(gomock.Any(), "admin", "secret").
			Return(user, rawToken, nil)

		handler := Register(context.Background(), responser, validate, auther, logger)

		apitest.New().
			HandlerFunc(handler).
			Post("/api/user/register").
			Body(`{
                "login": "admin",
                "password": "secret"
            }`).
			Expect(t).
			Status(http.StatusOK).
			CookiePresent(middlewares.AuthCookieName).
			Cookie(middlewares.AuthCookieName, string(rawToken)).
			End()
	})

	t.Run("user already exists", func(t *testing.T) {
		auther := mocks.NewMockAuther(ctrl)
		auther.EXPECT().
			RegisterUser(gomock.Any(), "admin", "secret").
			Return(models.User{}, types.RawToken(""), services.ErrLoginAlreadyExists)

		handler := Register(context.Background(), responser, validate, auther, logger)

		apitest.New().
			HandlerFunc(handler).
			Post("/api/user/register").
			Body(`{
                "login": "admin",
                "password": "secret"
            }`).
			Expect(t).
			Status(http.StatusConflict).
			CookieNotPresent(middlewares.AuthCookieName).
			End()
	})
}
