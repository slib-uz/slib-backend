package mapper

import (
	"context"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/async"
)

// toAsynqHandler converts a domain AsyncTaskHandler to an asynq.Handler
type toAsynqHandler struct {
	handler async.AsyncTaskHandler
}

func (h *toAsynqHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	task := entity.NewAsyncTask(enum.Task(t.Type()), t.Payload())
	taskID, ok := asynq.GetTaskID(ctx)
	if !ok {
		taskID = uuid.NewString()
	}
	task.TaskID = taskID
	return h.handler.ProcessTask(ctx, task)
}

func ToAsynqHandler(handler async.AsyncTaskHandler) asynq.Handler {
	return &toAsynqHandler{handler: handler}
}

// toDomainHandler converts an asynq.Handler to a domain AsyncTaskHandler
type toDomainHandler struct {
	handler asynq.Handler
}

func (h *toDomainHandler) ProcessTask(ctx context.Context, task *entity.AsyncTask) error {
	t := asynq.NewTask(string(task.TaskType), task.Payload)
	return h.handler.ProcessTask(ctx, t)
}

func ToDomainHandler(handler asynq.Handler) async.AsyncTaskHandler {
	return &toDomainHandler{handler: handler}
}
