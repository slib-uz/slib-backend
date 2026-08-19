package repository

import "slib.uz/src/core/domain/entity"

type InstructionRepository interface {
	CreateOrUpdate(instruction *entity.InstructionEntity) error
	Delete(id uint) error
	GetAll() ([]*entity.InstructionEntity, error)
	GetByKey(key string) (*entity.InstructionEntity, error)
}
