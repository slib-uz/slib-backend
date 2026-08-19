package tasks

import (
	"fmt"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/tasks/publisher"
)

type SendArticleKafkaTaskSender struct {
	publisher publisher.TaskPublisher
}

// @inject
func NewSendArticleKafkaTaskSender(publisher publisher.TaskPublisher) *SendArticleKafkaTaskSender {
	return &SendArticleKafkaTaskSender{publisher: publisher}
}

func (this *SendArticleKafkaTaskSender) Run(applicationID uint) error {
	fmt.Println("SendArticleKafkaTaskSender.Run")
	task := entity.NewTaskEntity[any](enum.TaskSendArticleKafka, applicationID)
	return this.publisher.Publish(task, 5)
}
