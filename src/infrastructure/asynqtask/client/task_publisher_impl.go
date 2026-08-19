package client

import (
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/domain/ports/tasks/publisher"
	"slib.uz/src/core/utils"
)

type TaskPublisherImpl struct {
	client     *asynq.Client
	repository repository.TaskResultRepository
}

// @inject
func NewTaskPublisher(client *asynq.Client, repository repository.TaskResultRepository) publisher.TaskPublisher {
	return &TaskPublisherImpl{client: client, repository: repository}
}

func (this *TaskPublisherImpl) Publish(task *entity.TaskEntity[any], maxRetryCount int) error {

	taskID := this.createPending(task)

	asynqTask := asynq.NewTask(string(task.Task), utils.JsonMarshal(task))

	_, err := this.client.Enqueue(
		asynqTask, asynq.MaxRetry(maxRetryCount),
		asynq.TaskID(taskID),
	)

	return err
}

func (this *TaskPublisherImpl) PublishWithOptions(task *entity.TaskEntity[any], maxRetryCount int, processIn, unique *time.Duration) error {
	taskID := this.createPending(task)

	asynqTask := asynq.NewTask(string(task.Task), utils.JsonMarshal(task))

	var options = []asynq.Option{
		asynq.MaxRetry(maxRetryCount),
		asynq.TaskID(taskID),
	}

	if processIn != nil {
		options = append(options, asynq.ProcessIn(*processIn))
	}
	if unique != nil {
		options = append(options, asynq.Unique(*unique))
	}

	_, err := this.client.Enqueue(
		asynqTask,
		options...,
	)

	return err
}
func (this *TaskPublisherImpl) createPending(task *entity.TaskEntity[any]) string {
	taskID := uuid.NewString()

	taskResult := entity.NewTaskResultEntity(taskID, task.Task, string(utils.JsonMarshal(task.Payload)), enum.TaskStatusPending, nil, 0, nil, nil, nil)
	_ = this.repository.Create(taskResult)

	return taskID
}
