package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"go.uber.org/zap"
)

type OrderJobProcessor interface {
	GetName() string

	Process(ctx context.Context, order models.Order) (models.Order, error)
}

var (
	ErrJobRetry = errors.New("retriable error occurred during job processing")
)

type ErrWorkerRetry struct {
	RetryAfter int
}

func (e *ErrWorkerRetry) Error() string {
	return fmt.Sprintf("worker should retry jobs after %d", e.RetryAfter)
}

type MemoryOrdersQueue struct {
	processor OrderJobProcessor
	logger    *zap.Logger

	jobsQueue  chan models.Order
	jobTimeout time.Duration
	jobsDelay  time.Duration
}

func (q *MemoryOrdersQueue) Add(ctx context.Context, order models.Order) error {
	select {
	case <-ctx.Done():
		return errors.New("failed to add to order queue")
	case q.jobsQueue <- order:
		return nil
	}
}

// завершать после остановки хэндлеров
func (q *MemoryOrdersQueue) Stop() {
	close(q.jobsQueue)
}

func (q *MemoryOrdersQueue) RunWorker(ctx context.Context) {
	q.logger.Info("running queue worker...",
		zap.String("job_name", q.processor.GetName()),
	)

	for {
		// лайфхак для простоты тестов - избавиться от рандомности в select при закрытом контексте
		if ctx.Err() != nil {
			return
		}

		select {
		case <-ctx.Done():
			return
		case order := <-q.jobsQueue:
			jobCtx, cancel := context.WithTimeout(ctx, q.jobTimeout)
			defer cancel()

			order, err := q.processor.Process(jobCtx, order)
			if err != nil {
				var e *ErrWorkerRetry
				if errors.As(err, &e) {
					q.logger.Info("pausing worker...",
						zap.String("order_number", order.OrderNumber),
						zap.String("job_name", q.processor.GetName()),
						zap.Int("retry_delay", e.RetryAfter),
						zap.Error(e),
					)

					workerDelay := time.Second * time.Duration(e.RetryAfter)

					err := q.retry(ctx, order, workerDelay)
					if err != nil {
						q.logger.Error("failed to retry job after worker was paused",
							zap.String("order_number", order.OrderNumber),
							zap.String("job_name", q.processor.GetName()),
							zap.Error(err),
						)
					}
				} else if errors.Is(err, ErrJobRetry) {
					q.logger.Info("retrying order job",
						zap.String("order_number", order.OrderNumber),
						zap.String("job_name", q.processor.GetName()),
					)

					go func() {
						retryCtx, cancel := context.WithTimeout(ctx, q.jobTimeout)
						defer cancel()

						err := q.retry(retryCtx, order, q.jobsDelay)
						if err != nil {
							q.logger.Error("failed to retry job",
								zap.String("order_number", order.OrderNumber),
								zap.String("job_name", q.processor.GetName()),
								zap.Error(err),
							)
						}
					}()
				} else {
					q.logger.Error("order job failed",
						zap.String("order_number", order.OrderNumber),
						zap.String("job_name", q.processor.GetName()),
					)
				}
			} else {
				q.logger.Debug("order job finished",
					zap.String("order_number", order.OrderNumber),
					zap.String("job_name", q.processor.GetName()),
				)
			}
		}
	}
}

func (q *MemoryOrdersQueue) retry(
	ctx context.Context,
	order models.Order,
	delay time.Duration,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(delay):
		return q.Add(ctx, order)
	}
}

func NewOrdersQueue(
	ordersProcessor OrderJobProcessor,
	logger *zap.Logger,
	jobsDelay time.Duration,
) *MemoryOrdersQueue {
	jobTimeout := 10 * time.Second
	jobsQueue := make(chan models.Order, 100)

	queue := MemoryOrdersQueue{
		processor: ordersProcessor,
		logger:    logger,

		jobsQueue:  jobsQueue,
		jobTimeout: jobTimeout,
		jobsDelay:  jobsDelay,
	}

	return &queue
}
