package tasks

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/tasks/publisher"
)

type SendTelegramAlertTask struct {
	taskPublisher publisher.TaskPublisher
}

// @inject
func NewSendTelegramAlertTask(taskPublisher publisher.TaskPublisher) *SendTelegramAlertTask {
	return &SendTelegramAlertTask{taskPublisher: taskPublisher}
}

func (this *SendTelegramAlertTask) Run(message string) error {
	alertEntity := entity2.NewTelegramAlertEntity(message)
	taskData := entity2.NewTaskEntity[any](enum.TaskSendTelegramAlert, alertEntity)
	return this.taskPublisher.Publish(taskData, 0)
}
