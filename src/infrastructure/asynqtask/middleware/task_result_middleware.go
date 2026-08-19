package middleware

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/async"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/logger"
)

type TaskResultMiddleware struct {
	repository repository.TaskResultRepository
	asyncCtx   async.AsyncContext
	fileLog    *logger.AsyncLogger
	consoleLog *logger.ConsoleLogger
}

// @inject
func NewTaskResultMiddleware(
	repository repository.TaskResultRepository,
	asyncCtx async.AsyncContext,
	fileLog *logger.AsyncLogger,
	consoleLog *logger.ConsoleLogger,
) *TaskResultMiddleware {
	return &TaskResultMiddleware{
		repository: repository,
		asyncCtx:   asyncCtx,
		fileLog:    fileLog,
		consoleLog: consoleLog,
	}
}

func (m *TaskResultMiddleware) Wrap(next async.AsyncTaskHandler) async.AsyncTaskHandler {
	return async.AsyncTaskHandlerFunc(func(ctx context.Context, task *entity.AsyncTask) error {
		start := time.Now()

		taskID := task.TaskID
		retried := m.asyncCtx.GetRetryCount(ctx)
		maxRetry := m.asyncCtx.GetMaxRetry(ctx)

		m.log(taskID, string(task.TaskType), enum.TaskStatusInProgress, 0, nil)
		m.updateStatus(taskID, enum.TaskStatusInProgress, 0, nil)

		err := next.ProcessTask(ctx, task)
		latency := time.Since(start)

		if err != nil {
			errMessage := err.Error()
			if retried < maxRetry {
				m.log(taskID, string(task.TaskType), enum.TaskStatusRetrying, latency, &errMessage)
				m.updateStatus(taskID, enum.TaskStatusRetrying, retried+1, &errMessage)
			} else {
				m.log(taskID, string(task.TaskType), enum.TaskStatusFailed, latency, &errMessage)
				m.updateStatus(taskID, enum.TaskStatusFailed, retried+1, &errMessage)
			}
			return err
		}

		m.log(taskID, string(task.TaskType), enum.TaskStatusCompleted, latency, nil)
		m.updateStatus(taskID, enum.TaskStatusCompleted, 0, nil)
		return nil
	})
}

func (m *TaskResultMiddleware) updateStatus(taskID string, status enum.TaskStatus, retryCount int, errMsg *string) {
	var completedAt *time.Time
	if status == enum.TaskStatusCompleted || status == enum.TaskStatusFailed {
		now := time.Now()
		completedAt = &now
	}

	updates := map[string]any{
		"status":       status,
		"retry_count":  retryCount,
		"completed_at": completedAt,
	}
	if errMsg != nil {
		updates["error"] = *errMsg
	}

	if err := m.repository.UpdateByTaskID(taskID, updates); err != nil {
		m.fileLog.Error("[TaskResultMiddleware] DB update failed",
			zap.String("taskID", taskID),
			zap.String("status", string(status)),
			zap.Error(err),
		)
	}
}

func (m *TaskResultMiddleware) log(taskID string, taskType string, status enum.TaskStatus, latency time.Duration, errMsg *string) {
	fields := []zap.Field{
		zap.String("taskID", taskID),
		zap.String("task", taskType),
		zap.String("status", string(status)),
		zap.Duration("latency", latency),
	}
	if errMsg != nil {
		fields = append(fields, zap.String("error", *errMsg))
	}

	m.fileLog.Info("[Task]", fields...)

	msg := fmt.Sprintf("%s | %s | %s | %s",
		colorTaskStatus(status),
		taskType,
		taskID,
		colorLatency(latency),
	)
	if errMsg != nil {
		msg += fmt.Sprintf(" | err=%s", *errMsg)
	}
	m.consoleLog.Info(msg)
}

func colorTaskStatus(status enum.TaskStatus) string {
	switch status {
	case enum.TaskStatusCompleted:
		return "\033[32m" + string(status) + "\033[0m"
	case enum.TaskStatusFailed:
		return "\033[31m" + string(status) + "\033[0m"
	case enum.TaskStatusRetrying:
		return "\033[33m" + string(status) + "\033[0m"
	case enum.TaskStatusInProgress:
		return "\033[34m" + string(status) + "\033[0m"
	default:
		return string(status)
	}
}

func colorLatency(d time.Duration) string {
	ms := d.Milliseconds()
	switch {
	case ms < 100:
		return fmt.Sprintf("\033[32m%dms\033[0m", ms)
	case ms < 500:
		return fmt.Sprintf("\033[33m%dms\033[0m", ms)
	default:
		return fmt.Sprintf("\033[31m%dms\033[0m", ms)
	}
}
