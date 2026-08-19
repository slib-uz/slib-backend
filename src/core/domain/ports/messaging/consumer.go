package messaging

import (
	"slib.uz/src/core/domain/entity"
)

type Consumer interface {
	SubscribeTopics([]string) error
	Listen(onMessage func(message entity.MessageEntity))
	Close()
}
