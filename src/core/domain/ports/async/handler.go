package async

import (
	"context"

	"slib.uz/src/core/domain/entity"
)

type AsyncTaskHandler interface {
	ProcessTask(ctx context.Context, task *entity.AsyncTask) error
}

type AsyncTaskHandlerFunc func(ctx context.Context, task *entity.AsyncTask) error

func (f AsyncTaskHandlerFunc) ProcessTask(ctx context.Context, task *entity.AsyncTask) error {
	return f(ctx, task)
}
