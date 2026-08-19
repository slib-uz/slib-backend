package asynqtask

import (
	"context"

	"github.com/hibiken/asynq"
	"slib.uz/src/core/domain/ports/async"
)

type AsyncContextImpl struct{}

// @inject
func NewAsyncContextImpl() async.AsyncContext {
	return &AsyncContextImpl{}
}

func (c *AsyncContextImpl) GetTaskID(ctx context.Context) string {
	id, ok := asynq.GetTaskID(ctx)
	if !ok {
		return ""
	}
	return id
}

func (c *AsyncContextImpl) GetRetryCount(ctx context.Context) int {
	count, ok := asynq.GetRetryCount(ctx)
	if !ok {
		return 0
	}
	return count
}

func (c *AsyncContextImpl) GetMaxRetry(ctx context.Context) int {
	count, ok := asynq.GetMaxRetry(ctx)
	if !ok {
		return 0
	}
	return count
}
