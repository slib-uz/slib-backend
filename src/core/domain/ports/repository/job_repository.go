package repository

import (
	"slib.uz/src/core/domain/entity"
)

type JobRepository interface {
	GetByUserID(userID uint) ([]*entity.JobEntity, error)
	GetByUserIDSingle(userID uint) (*entity.JobEntity, error)
	Create(job *entity.JobEntity) (*entity.JobEntity, error)
	Update(userID uint, job *entity.JobEntity) (*entity.JobEntity, error)
	UpdateOrCreate(userID uint, entity []*entity.JobEntity) error
	DeleteLastByUserID(userID uint) error
}
