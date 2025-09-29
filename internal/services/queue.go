package services

import (
	"context"
	"errors"
	"time"

	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"go.uber.org/zap"
)

type OrdersQueue interface {
	Add(ctx context.Context, order models.Order) error
	Stop()
}

type MemoryOrdersQueue struct {
	processor OrderJobProcessor
	logger    *zap.Logger

	jobsQueue  chan models.Order
	jobTimeout time.Duration
	jobsDelay  time.Duration
}

var (
	ErrJobRetry = errors.New("retriable error occurred during job processing")
)

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

func (q *MemoryOrdersQueue) runWorker(ctx context.Context) {
	q.logger.Info("running queue worker...",
		zap.String("job_name", q.processor.GetName()),
	)

	for {
		select {
		case <-ctx.Done():
			return
		case order := <-q.jobsQueue:
			jobCtx, cancel := context.WithTimeout(ctx, q.jobTimeout)
			defer cancel()

			order, err := q.processor.Process(jobCtx, order)
			if err != nil {
				if errors.Is(err, ErrJobRetry) {
					q.logger.Info("retrying order job",
						zap.String("order_number", order.OrderNumber),
						zap.String("job_name", q.processor.GetName()),
					)

					go func() {
						retryCtx, cancel := context.WithTimeout(ctx, q.jobTimeout)
						defer cancel()

						q.retry(retryCtx, order)
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

func (q *MemoryOrdersQueue) retry(ctx context.Context, order models.Order) {
	select {
	case <-ctx.Done():
		q.logger.Error("worker stopped before order job retry",
			zap.String("order_number", order.OrderNumber),
			zap.String("job_name", q.processor.GetName()),
		)
		return
	case <-time.After(q.jobsDelay):
		q.Add(ctx, order)
		return
	}
}

func NewOrdersQueue(
	ctx context.Context,
	ordersProcessor OrderJobProcessor,
	logger *zap.Logger,
) *MemoryOrdersQueue {
	workersNum := 3
	jobTimeout := 10 * time.Second
	jobsQueue := make(chan models.Order, 100)

	queue := MemoryOrdersQueue{
		processor: ordersProcessor,
		logger:    logger,

		jobsQueue:  jobsQueue,
		jobTimeout: jobTimeout,
	}

	for i := 0; i < workersNum; i++ {
		go func() {
			queue.runWorker(ctx)
		}()
	}

	return &queue
}
