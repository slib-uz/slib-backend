package entity

import (
	"slib.uz/src/core/domain/entity/enum"
)

type TaskEntity[T any] struct {
	//TaskID  string    `json:"task_id"`
	Task    enum.Task `json:"task"`
	Payload T         `json:"payload"`
}

func NewTaskEntity[T any](task enum.Task, payload T) *TaskEntity[T] {
	return &TaskEntity[T]{Task: task, Payload: payload}
}
