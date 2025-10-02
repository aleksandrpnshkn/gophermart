package middlewares

import (
	"net/http"
	"testing"

	"github.com/aleksandrpnshkn/gophermart/internal/mocks"
	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"github.com/aleksandrpnshkn/gophermart/internal/types"
	"github.com/steinfletcher/apitest"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uni := services.NewAppUni()
	responser := services.NewResponser(uni)

	t.Run("client sent valid token", func(t *testing.T) {
		testToken := types.RawToken("testToken")
		testUser := models.User{
			ID: 123,
		}

		auther := mocks.NewMockAuther(ctrl)
		auther.EXPECT().ParseToken(gomock.Any(), testToken).Return(testUser, nil)
		handler := NewAuthMiddleware(responser, zap.NewExample(), auther)(testOkHandler())

		apitest.New().
			Handler(handler).
			Post("/").
			Cookie(AuthCookieName, string(testToken)).
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("client sent invalid token", func(t *testing.T) {
		auther := services.NewAuther(mocks.NewMockUsersStorage(ctrl), "secretkey")
		handler := NewAuthMiddleware(responser, zap.NewExample(), auther)(testOkHandler())

		apitest.New().
			Handler(handler).
			Post("/").
			Cookie(AuthCookieName, "wrong token").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("client not sent token", func(t *testing.T) {
		auther := services.NewAuther(mocks.NewMockUsersStorage(ctrl), "secretkey")
		handler := NewAuthMiddleware(responser, zap.NewExample(), auther)(testOkHandler())

		apitest.New().
			Handler(handler).
			Post("/").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})
}

func testOkHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
	})
}
