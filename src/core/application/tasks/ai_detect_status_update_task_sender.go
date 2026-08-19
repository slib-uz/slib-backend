package tasks

import (
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/tasks/publisher"
)

type AiDetectStatusUpdateTask struct {
	publisher.TaskPublisher
}

// @inject
func NewAiDetectStatusUpdateTask(taskPublisher publisher.TaskPublisher) *AiDetectStatusUpdateTask {
	return &AiDetectStatusUpdateTask{TaskPublisher: taskPublisher}
}

func (this *AiDetectStatusUpdateTask) Run(externalID uint) error {
	task := entity.NewTaskEntity[any](enum.TaskAiDetectStatusUpdate, externalID)

	processIn := 2 * time.Minute
	uniqueIn := 2 * time.Minute

	return this.PublishWithOptions(task, 720, &processIn, &uniqueIn)
}
