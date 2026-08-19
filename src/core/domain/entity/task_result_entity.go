package entity

import (
	"time"

	enum2 "slib.uz/src/core/domain/entity/enum"
)

type TaskResultEntity struct {
	TaskID     string           `json:"task_id"`
	Task       enum2.Task       `json:"task"`
	Payload    string           `json:"payload"`
	Status     enum2.TaskStatus `json:"status"`
	Error      *string          `json:"error"`
	RetryCount int              `json:"retires"`

	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

func NewTaskResultEntity(taskID string, task enum2.Task, payload string, status enum2.TaskStatus, error *string, retryCount int, createdAt *time.Time, updatedAt *time.Time, completedAt *time.Time) *TaskResultEntity {
	return &TaskResultEntity{TaskID: taskID, Task: task, Payload: payload, Status: status, Error: error, RetryCount: retryCount, CreatedAt: createdAt, UpdatedAt: updatedAt, CompletedAt: completedAt}
}
