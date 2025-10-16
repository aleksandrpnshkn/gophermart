package services

import (
	"context"
)

type RunnableQueue interface {
	RunWorker(ctx context.Context)
}

func RunQueue(ctx context.Context, queue RunnableQueue, workersNum int) {
	for i := 0; i < workersNum; i++ {
		go func() {
			queue.RunWorker(ctx)
		}()
	}
}
