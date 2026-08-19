package models

import (
	"time"

	"gorm.io/gorm"
	enum2 "slib.uz/src/core/domain/entity/enum"
)

type TaskResultModel struct {
	gorm.Model

	TaskID      string     `gorm:"index;not null"`
	Task        enum2.Task `gorm:"not null;size:32"`
	Payload     string
	Status      enum2.TaskStatus `gorm:"size:32;default:PENDING"`
	Error       *string
	RetryCount  int
	CompletedAt *time.Time
}

func NewTaskResultModel(taskID string, task enum2.Task, payload string, status enum2.TaskStatus, error *string, retryCount int, completedAt *time.Time) TaskResultModel {
	return TaskResultModel{TaskID: taskID, Task: task, Payload: payload, Status: status, Error: error, RetryCount: retryCount, CompletedAt: completedAt}
}
