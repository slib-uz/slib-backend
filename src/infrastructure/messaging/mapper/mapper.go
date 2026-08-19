package mapper

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"slib.uz/src/core/domain/entity"
)

func KafkaToMessage(m *kafka.Message) entity.MessageEntity {
	topic := ""
	headers := make(map[string][]byte)

	if m.TopicPartition.Topic != nil {
		topic = *m.TopicPartition.Topic
	}

	if m.Headers != nil {
		for _, header := range m.Headers {
			headers[header.Key] = header.Value
		}
	}

	return entity.NewMessageEntity(
		m.Key,
		m.Value,
		m.Timestamp,
		headers,
		topic,
	)
}
