package services

import (
	"context"
	"testing"
	"time"

	"github.com/aleksandrpnshkn/gophermart/internal/mocks"
	"github.com/aleksandrpnshkn/gophermart/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestQueue(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := zap.NewExample()

	o1 := models.Order{
		OrderNumber: "111",
	}
	o2 := models.Order{
		OrderNumber: "222",
	}

	jobsDelay := time.Duration(0)

	t.Run("jobs completed", func(t *testing.T) {
		testProcessor := mocks.NewMockOrderJobProcessor(ctrl)
		testProcessor.EXPECT().GetName().AnyTimes().Return("jobs completed")
		testProcessor.EXPECT().Process(gomock.Any(), o1).Return(o1, nil)
		testProcessor.EXPECT().Process(gomock.Any(), o2).Return(o2, nil)

		queue := NewOrdersQueue(testProcessor, logger, jobsDelay)
		defer queue.Stop()

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := queue.Add(ctx, o1)
		assert.NoError(t, err, "order 1 should be added")
		err = queue.Add(ctx, o2)
		assert.NoError(t, err, "order 2 should be added")

		queue.RunWorker(ctx)
	})

	t.Run("job retried", func(t *testing.T) {
		testProcessor := mocks.NewMockOrderJobProcessor(ctrl)
		testProcessor.EXPECT().GetName().AnyTimes().Return("job retried")
		testProcessor.EXPECT().Process(gomock.Any(), o1).Return(o1, ErrJobRetry)
		testProcessor.EXPECT().Process(gomock.Any(), o1).Return(o1, nil)

		queue := NewOrdersQueue(testProcessor, logger, jobsDelay)
		defer queue.Stop()

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := queue.Add(ctx, o1)
		assert.NoError(t, err, "order 1 should be added")

		queue.RunWorker(ctx)
	})

	t.Run("worker retried", func(t *testing.T) {
		testProcessor := mocks.NewMockOrderJobProcessor(ctrl)
		testProcessor.EXPECT().GetName().AnyTimes().Return("worker retried")
		testProcessor.EXPECT().Process(gomock.Any(), o1).Return(o1, &ErrWorkerRetry{
			RetryAfter: 0,
		})
		testProcessor.EXPECT().Process(gomock.Any(), o1).Return(o1, nil)

		queue := NewOrdersQueue(testProcessor, logger, jobsDelay)
		defer queue.Stop()

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := queue.Add(ctx, o1)
		assert.NoError(t, err, "order 1 should be added")

		queue.RunWorker(ctx)
	})

	t.Run("worker paused", func(t *testing.T) {
		testProcessor := mocks.NewMockOrderJobProcessor(ctrl)
		testProcessor.EXPECT().GetName().AnyTimes().Return("worker paused")
		testProcessor.EXPECT().Process(gomock.Any(), o1).Return(o1, &ErrWorkerRetry{
			RetryAfter: 1,
		})

		queue := NewOrdersQueue(testProcessor, logger, jobsDelay)
		defer queue.Stop()

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := queue.Add(ctx, o1)
		assert.NoError(t, err, "order 1 should be added")

		err = queue.Add(ctx, o2)
		assert.NoError(t, err, "order 2 should be added")

		queue.RunWorker(ctx)
	})
}
