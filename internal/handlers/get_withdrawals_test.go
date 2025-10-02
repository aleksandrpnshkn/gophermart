package handlers

import (
	"net/http"
	"testing"
	"time"

	"github.com/aleksandrpnshkn/gophermart/internal/mocks"
	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"github.com/aleksandrpnshkn/gophermart/internal/types"
	"github.com/shopspring/decimal"
	"github.com/steinfletcher/apitest"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestGetWithdrawals(t *testing.T) {
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

	t.Run("get user withdrawals", func(t *testing.T) {
		auther := mocks.NewMockAuther(ctrl)
		auther.EXPECT().
			FromUserContext(gomock.Any()).
			Return(user, nil)

		loc, _ := time.LoadLocation("Europe/Moscow")
		withdrawals := []models.BalanceChange{
			{
				OrderNumber: "2377225624",
				UserID:      user.ID,
				Amount:      decimal.NewFromInt(-500),
				ProcessedAt: time.Date(2020, 12, 9, 16, 9, 57, 0, loc),
			},
		}

		balancer := mocks.NewMockIBalancer(ctrl)
		balancer.EXPECT().
			GetWithdrawals(gomock.Any(), user).
			Return(withdrawals, nil)

		handler := GetWithdrawals(responser, auther, balancer, logger)

		apitest.New().
			HandlerFunc(handler).
			Post("/api/user/withdrawals").
			Expect(t).
			Status(http.StatusOK).
			Body(`[
                {
                    "order": "2377225624",
                    "sum": 500,
                    "processed_at": "2020-12-09T16:09:57+03:00"
                }
            ]`).
			End()
	})
}
