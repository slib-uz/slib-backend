package repository

import (
	"slib.uz/src/core/domain/entity"
)

type TaskResultRepository interface {
	Create(task *entity.TaskResultEntity) error
	Update(task *entity.TaskResultEntity) error
	UpdateByTaskID(taskID string, updates map[string]any) error
	GetByTaskID(taskID string) (*entity.TaskResultEntity, error)
	UpdateOrCreate(task *entity.TaskResultEntity) error
}
