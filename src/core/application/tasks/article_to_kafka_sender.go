package tasks

import (
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/tasks/publisher"
)

type SendArticleToKafkaTaskSender struct {
	publisher publisher.TaskPublisher
}

// @inject
func NewSendArticleToKafkaTask(publisher publisher.TaskPublisher) *SendArticleToKafkaTaskSender {
	return &SendArticleToKafkaTaskSender{publisher: publisher}
}

func (this *SendArticleToKafkaTaskSender) Run(msg *entity2.MessageEntity) error {
	task := entity2.NewTaskEntity[any](enum.TaskSendArticleToKafka, msg)
	return this.publisher.Publish(task, 5)
}
