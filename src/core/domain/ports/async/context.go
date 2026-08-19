package async

import "context"

type AsyncContext interface {
	GetRetryCount(ctx context.Context) int
	GetMaxRetry(ctx context.Context) int
	GetTaskID(ctx context.Context) string
}
