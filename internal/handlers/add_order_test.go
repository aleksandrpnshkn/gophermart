package handlers

import (
	"net/http"
	"testing"

	"github.com/aleksandrpnshkn/gophermart/internal/mocks"
	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"github.com/aleksandrpnshkn/gophermart/internal/services"
	"github.com/aleksandrpnshkn/gophermart/internal/storage/orders"
	"github.com/aleksandrpnshkn/gophermart/internal/types"
	"github.com/steinfletcher/apitest"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestAddOrder(t *testing.T) {
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

		accrualService := mocks.NewMockIAccrualService(ctrl)

		ordersStorage := mocks.NewMockOrdersStorage(ctrl)
		ordersStorage.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return(models.Order{}, nil)
		ordersService := services.NewOrdersService(ordersStorage, accrualService, logger)

		ordersQueue := mocks.NewMockOrdersQueue(ctrl)
		ordersQueue.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil)

		handler := AddOrder(responser, auther, logger, ordersService, ordersQueue)

		apitest.New().
			HandlerFunc(handler).
			Post("/api/user/orders").
			ContentType("text/plain").
			Body("125").
			Expect(t).
			Status(http.StatusAccepted).
			End()
	})

	t.Run("invalid order number", func(t *testing.T) {
		auther := mocks.NewMockAuther(ctrl)
		auther.EXPECT().
			FromUserContext(gomock.Any()).
			Return(user, nil)

		ordersService := mocks.NewMockIOrdersService(ctrl)
		ordersQueue := mocks.NewMockOrdersQueue(ctrl)

		handler := AddOrder(responser, auther, logger, ordersService, ordersQueue)

		apitest.New().
			HandlerFunc(handler).
			Post("/api/user/orders").
			ContentType("text/plain").
			Body("bad order number").
			Expect(t).
			Status(http.StatusUnprocessableEntity).
			End()
	})

	t.Run("order already created", func(t *testing.T) {
		auther := mocks.NewMockAuther(ctrl)
		auther.EXPECT().
			FromUserContext(gomock.Any()).
			Return(user, nil)

		accrualService := mocks.NewMockIAccrualService(ctrl)

		ordersStorage := mocks.NewMockOrdersStorage(ctrl)
		ordersStorage.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return(models.Order{}, orders.ErrOrderAlreadyCreated)
		ordersService := services.NewOrdersService(ordersStorage, accrualService, logger)

		ordersQueue := mocks.NewMockOrdersQueue(ctrl)

		handler := AddOrder(responser, auther, logger, ordersService, ordersQueue)

		apitest.New().
			HandlerFunc(handler).
			Post("/api/user/orders").
			ContentType("text/plain").
			Body("125").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("order already created by another user", func(t *testing.T) {
		auther := mocks.NewMockAuther(ctrl)
		auther.EXPECT().
			FromUserContext(gomock.Any()).
			Return(user, nil)

		accrualService := mocks.NewMockIAccrualService(ctrl)

		ordersStorage := mocks.NewMockOrdersStorage(ctrl)
		ordersStorage.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return(models.Order{}, orders.ErrOrderAlreadyCreatedByAnotherUser)
		ordersService := services.NewOrdersService(ordersStorage, accrualService, logger)

		ordersQueue := mocks.NewMockOrdersQueue(ctrl)

		handler := AddOrder(responser, auther, logger, ordersService, ordersQueue)

		apitest.New().
			HandlerFunc(handler).
			Post("/api/user/orders").
			ContentType("text/plain").
			Body("125").
			Expect(t).
			Status(http.StatusConflict).
			End()
	})
}
