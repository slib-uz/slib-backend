package tasks

import (
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/tasks/publisher"
)

type AntiPlagStatusUpdateTask struct {
	publisher.TaskPublisher
}

// @inject
func NewAntiPlagStatusUpdateTask(taskPublisher publisher.TaskPublisher) *AntiPlagStatusUpdateTask {
	return &AntiPlagStatusUpdateTask{TaskPublisher: taskPublisher}
}

func (this *AntiPlagStatusUpdateTask) Run(externalID uint) error {
	task := entity.NewTaskEntity[any](enum.TaskAntiPlagStatusUpdate, externalID)

	processIn := 2 * time.Minute // Time after which the task should be processed
	uniqueIn := 2 * time.Minute  // Time after which the task should be considered unique

	return this.PublishWithOptions(task, 720, &processIn, &uniqueIn)
}
