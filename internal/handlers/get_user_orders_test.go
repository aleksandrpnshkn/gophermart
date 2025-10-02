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

func TestGetUserOrders(t *testing.T) {
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

	t.Run("new order added", func(t *testing.T) {
		auther := mocks.NewMockAuther(ctrl)
		auther.EXPECT().
			FromUserContext(gomock.Any()).
			Return(user, nil)

		loc, _ := time.LoadLocation("Europe/Moscow")
		uploadedAt := time.Date(2020, 12, 10, 15, 15, 45, 0, loc)

		orders := []models.Order{
			{
				OrderNumber: "3",
				Status:      types.OrderStatusNew,
				UploadedAt:  uploadedAt,
			},
			{
				OrderNumber: "2",
				Status:      types.OrderStatusProcessed,
				UploadedAt:  uploadedAt,
				Accrual:     decimal.NewFromFloat32(3.50),
			},
			{
				OrderNumber: "1",
				Status:      types.OrderStatusProcessed,
				UploadedAt:  uploadedAt,
				Accrual:     decimal.NewFromFloat32(3),
			},
		}

		accrualService := mocks.NewMockIAccrualService(ctrl)
		ordersStorage := mocks.NewMockOrdersStorage(ctrl)
		ordersStorage.EXPECT().
			GetUserOrders(gomock.Any(), gomock.Any()).
			Return(orders, nil)
		ordersService := services.NewOrdersService(ordersStorage, accrualService, logger)

		handler := GetUserOrders(responser, auther, logger, ordersService)

		apitest.New().
			HandlerFunc(handler).
			Get("/api/user/orders").
			Expect(t).
			Status(http.StatusOK).
			Body(`[
                {
                    "number": "3",
                    "status": "NEW",
                    "uploaded_at": "2020-12-10T15:15:45+03:00"
                },
                {
                    "number": "2",
                    "status": "PROCESSED",
                    "accrual": 3.50,
                    "uploaded_at": "2020-12-10T15:15:45+03:00"
                },
                {
                    "number": "1",
                    "status": "PROCESSED",
                    "accrual": 3,
                    "uploaded_at": "2020-12-10T15:15:45+03:00"
                }
            ]`).
			End()
	})

	t.Run("user has no orders", func(t *testing.T) {
		auther := mocks.NewMockAuther(ctrl)
		auther.EXPECT().
			FromUserContext(gomock.Any()).
			Return(user, nil)

		accrualService := mocks.NewMockIAccrualService(ctrl)
		ordersStorage := mocks.NewMockOrdersStorage(ctrl)
		ordersStorage.EXPECT().
			GetUserOrders(gomock.Any(), gomock.Any()).
			Return([]models.Order{}, nil)
		ordersService := services.NewOrdersService(ordersStorage, accrualService, logger)

		handler := GetUserOrders(responser, auther, logger, ordersService)

		apitest.New().
			HandlerFunc(handler).
			Get("/api/user/orders").
			Expect(t).
			Status(http.StatusNoContent).
			Body("").
			End()
	})

}
