package tasks

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/tasks/publisher"
)

type OutboxEventTaskSender struct {
	publisher.TaskPublisher
}

// @inject
func NewOutboxEventTaskSender(taskPublisher publisher.TaskPublisher) *OutboxEventTaskSender {
	return &OutboxEventTaskSender{TaskPublisher: taskPublisher}
}

func (this *OutboxEventTaskSender) Run(event *entity.OutboxEventEntity) error {
	data := entity.NewTaskEntity[any](enum.TaskOutboxEvent, event)
	return this.Publish(data, 3)
}
