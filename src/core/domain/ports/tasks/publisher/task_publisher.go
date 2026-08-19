package publisher

import (
	"time"

	"slib.uz/src/core/domain/entity"
)

type TaskPublisher interface {
	Publish(task *entity.TaskEntity[any], maxRetryCount int) error
	PublishWithOptions(task *entity.TaskEntity[any], maxRetryCount int, processIn, unique *time.Duration) error
}
