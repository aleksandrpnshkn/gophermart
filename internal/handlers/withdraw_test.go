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

func TestWithdraw(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uni := services.NewAppUni()
	responser := services.NewResponser(uni)
	validate := services.NewValidate(uni)
	logger := zap.NewExample()

	user := models.User{
		ID:    1,
		Login: "admin",
		Hash:  types.PasswordHash("hash"),
	}

	t.Run("widthdraw successfully", func(t *testing.T) {
		auther := mocks.NewMockAuther(ctrl)
		auther.EXPECT().
			FromUserContext(gomock.Any()).
			Return(user, nil)

		newOrder := models.Order{
			OrderNumber: "2377225624",
			UserID:      user.ID,
		}

		ordersService := mocks.NewMockIOrdersService(ctrl)
		ordersService.EXPECT().
			Add(gomock.Any(), "2377225624", user).
			Return(newOrder, nil)

		balance := models.Balance{
			Current:   decimal.NewFromInt(1000),
			Withdrawn: decimal.NewFromInt(0),
		}

		expectedWithdrawal := models.BalanceChange{
			OrderNumber: "2377225624",
			UserID:      user.ID,
			Amount:      decimal.NewFromInt(-751),
		}

		balanceStorage := mocks.NewMockBalanceStorage(ctrl)
		balanceStorage.EXPECT().
			GetBalance(gomock.Any(), user).
			Return(balance, nil)
		balanceStorage.EXPECT().
			Withdraw(gomock.Any(), expectedWithdrawal).
			Return(nil)

		balancer := services.NewBalancer(ordersService, balanceStorage, logger)

		handler := Withdraw(responser, validate, auther, balancer, logger)

		apitest.New().
			HandlerFunc(handler).
			Post("/api/user/balance/withdraw").
			ContentType("application/json").
			Body(`{
                "order": "2377225624",
                "sum": 751
            }`).
			Expect(t).
			Status(http.StatusOK).
			End()
	})
}
