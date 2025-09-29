package services

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type IAccrualService interface {
	GetAccrual(ctx context.Context, orderNumber string) (decimal.Decimal, error)
}

type AccrualService struct {
	client  *http.Client
	baseURL string
	logger  *zap.Logger
}

var (
	ErrAccrualInvalidStatus      = errors.New("order is invalid")
	ErrAccrualUnknownStatus      = errors.New("order has unknown status")
	ErrAccrualNotProcessedStatus = errors.New("order is not processed yet")
	ErrAccrualOrderNotFound      = errors.New("order not found")
	ErrAccrualFailedToGet        = errors.New("failed to get accrual")
	ErrAccrualUnexpectedError    = errors.New("failed to get accrual with unexpected error")
)

type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

const (
	// заказ зарегистрирован, но вознаграждение не рассчитано
	statusRegistered = "REGISTERED"

	// заказ не принят к расчёту, и вознаграждение не будет начислено
	statusInvalid = "INVALID"

	// расчёт начисления в процессе
	statusProcessing = "PROCESSING"

	// расчёт начисления окончен
	statusProcessed = "PROCESSED"
)

func (a *AccrualService) GetAccrual(
	ctx context.Context,
	orderNumber string,
) (decimal.Decimal, error) {
	zero := decimal.NewFromInt(0)

	a.logger.Debug("sending accrual request...",
		zap.String("order_number", orderNumber),
	)

	res, err := a.client.Get(a.baseURL + "/api/orders/" + orderNumber)
	if err != nil {
		return zero, err
	}
	defer res.Body.Close()

	a.logger.Debug("accrual responsed",
		zap.String("order_number", orderNumber),
		zap.Int("status_code", res.StatusCode),
	)

	if res.StatusCode == http.StatusNoContent {
		return zero, ErrAccrualOrderNotFound
	}

	if res.StatusCode == http.StatusTooManyRequests ||
		res.StatusCode >= http.StatusInternalServerError {
		return zero, ErrAccrualFailedToGet
	}

	if res.StatusCode != http.StatusOK {
		a.logger.Error("unexpected status code from accrual",
			zap.String("order_number", orderNumber),
			zap.Int("status_code", res.StatusCode),
		)
		return zero, ErrAccrualUnexpectedError
	}

	rawResponse, err := io.ReadAll(res.Body)
	if err != nil {
		a.logger.Error("failed to read response body from accrual",
			zap.String("order_number", orderNumber),
			zap.Error(err),
		)
		return zero, ErrAccrualUnexpectedError
	}

	var accrualOrder AccrualResponse
	err = json.Unmarshal(rawResponse, &accrualOrder)
	if err != nil {
		a.logger.Error("failed to unmarshal response body from accrual",
			zap.String("order_number", orderNumber),
			zap.Error(err),
		)
		return zero, ErrAccrualUnexpectedError
	}

	switch accrualOrder.Status {
	case statusInvalid:
		return zero, ErrAccrualInvalidStatus
	case statusProcessed:
		return decimal.NewFromFloat(accrualOrder.Accrual), nil
	case statusRegistered, statusProcessing:
		return zero, ErrAccrualNotProcessedStatus
	default:
		a.logger.Error("unknown order status from accrual",
			zap.String("order_number", orderNumber),
			zap.String("order_status", accrualOrder.Status),
		)
		return zero, ErrAccrualUnknownStatus
	}
}

func NewAccrualService(
	client *http.Client,
	logger *zap.Logger,
	baseURL string,
) *AccrualService {
	return &AccrualService{
		client:  client,
		baseURL: baseURL,
		logger:  logger,
	}
}
