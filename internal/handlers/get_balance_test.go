package handlers

import (
	"net/http"
	"testing"

	"github.com/aleksandrpnshkn/gophermart/internal/mocks"
	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"github.com/aleksandrpnshkn/gophermart/internal/types"
	"github.com/shopspring/decimal"
	"github.com/steinfletcher/apitest"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestGetBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uni := services.NewAppUni()
	responser := services.NewResponser(uni)
	logger := zap.NewExample()

	user := models.User{
		ID:    1,
		Login: "admin",
		Hash:  types.PasswordHash("hash"),
	}

	t.Run("get user balance", func(t *testing.T) {
		auther := mocks.NewMockAuther(ctrl)
		auther.EXPECT().
			FromUserContext(gomock.Any()).
			Return(user, nil)

		balance := models.Balance{
			Current:   decimal.NewFromFloat(500.5),
			Withdrawn: decimal.NewFromInt(42),
		}

		balancer := mocks.NewMockIBalancer(ctrl)
		balancer.EXPECT().
			GetBalance(gomock.Any(), user).
			Return(balance, nil)

		handler := GetBalance(responser, auther, balancer, logger)

		apitest.New().
			HandlerFunc(handler).
			Post("/api/user/balance").
			Expect(t).
			Status(http.StatusOK).
			Body(`{
                "current": 500.5,
                "withdrawn": 42
            }`).
			End()
	})
}
