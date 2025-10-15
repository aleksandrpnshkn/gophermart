package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

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

type ErrAccrualFailedToGetWithRetry struct {
	RetryAfter int
}

func (e *ErrAccrualFailedToGetWithRetry) Error() string {
	return fmt.Sprintf("failed to get accrual, should retry after %d", e.RetryAfter)
}

type AccrualService struct {
	client  *http.Client
	baseURL string
	logger  *zap.Logger
}

func (a *AccrualService) GetAccrual(
	ctx context.Context,
	orderNumber string,
) (decimal.Decimal, error) {
	zero := decimal.NewFromInt(0)

	a.logger.Debug("sending accrual request...",
		zap.String("order_number", orderNumber),
	)

	accrualURL, err := url.JoinPath(a.baseURL, "/api/orders/", orderNumber)
	if err != nil {
		return zero, err
	}

	res, err := a.client.Get(accrualURL)
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
		rawRetryAfter := res.Header.Get("Retry-After")

		a.logger.Error("failed to get accrual, should retry later",
			zap.String("order_number", orderNumber),
			zap.Int("status_code", res.StatusCode),
			zap.String("retry_after", rawRetryAfter),
		)

		if rawRetryAfter != "" {
			retryAfter, err := strconv.Atoi(rawRetryAfter)
			if err != nil && retryAfter > 0 {
				return zero, &ErrAccrualFailedToGetWithRetry{
					RetryAfter: retryAfter,
				}
			}
		}

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
