package tasks

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/tasks/publisher"
)

type SendNotificationTask struct {
	taskPublisher publisher.TaskPublisher
}

// @inject
func NewSendNotificationTask(taskPublisher publisher.TaskPublisher) *SendNotificationTask {
	return &SendNotificationTask{taskPublisher: taskPublisher}
}

func (this *SendNotificationTask) Run(notification *entity2.NotificationEntity) error {
	taskData := entity2.NewTaskEntity[any](enum.TaskSendNotification, notification)
	return this.taskPublisher.Publish(taskData, 0)
}
