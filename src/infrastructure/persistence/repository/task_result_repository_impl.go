package repository

import (
	"errors"

	"gorm.io/gorm"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	infraError "slib.uz/src/infrastructure/errors"
	"slib.uz/src/infrastructure/persistence/mapper"
	"slib.uz/src/infrastructure/persistence/models"
)

type TaskResultRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewTaskResultRepository(r *BaseRepository) repository.TaskResultRepository {
	return &TaskResultRepositoryImpl{BaseRepository: r}
}

func (this *TaskResultRepositoryImpl) Create(task *entity.TaskResultEntity) error {
	var _model = mapper.TaskResultEntityToModel(task)
	return this.db().Create(&_model).Error
}

func (this *TaskResultRepositoryImpl) Update(task *entity.TaskResultEntity) error {
	var _model = mapper.TaskResultEntityToModel(task)
	return this.db().Model(&_model).Where("task_id = ?", task.TaskID).Updates(_model).Error
}

func (this *TaskResultRepositoryImpl) UpdateByTaskID(taskID string, updates map[string]any) error {
	return this.db().Model(&models.TaskResultModel{}).Where("task_id = ?", taskID).Updates(updates).Error
}

func (this *TaskResultRepositoryImpl) GetByTaskID(taskID string) (*entity.TaskResultEntity, error) {
	var taskResult models.TaskResultModel
	if err := this.db().Where("task_id = ?", taskID).First(&taskResult).Error; err != nil {
		return nil, infraError.Wrap(err)
	}
	return mapper.TaskResultModelToEntity(&taskResult), nil
}

func (this *TaskResultRepositoryImpl) UpdateOrCreate(task *entity.TaskResultEntity) error {

	var _model = mapper.TaskResultEntityToModel(task)

	if err := this.db().Where("task_id = ?", task.TaskID).First(&_model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return this.Create(task)
		}
		return err
	}
	return this.Update(task)
}
