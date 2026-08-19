package mapper

import (
	"github.com/hibiken/asynq"
	"slib.uz/src/core/domain/ports/async"
)

func ToAsynqMiddleware(middleware async.AsyncTaskMiddleware) func(asynq.Handler) asynq.Handler {
	return func(next asynq.Handler) asynq.Handler {
		return ToAsynqHandler(middleware(ToDomainHandler(next)))
	}
}
