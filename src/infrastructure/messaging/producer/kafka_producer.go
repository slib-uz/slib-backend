package producer

import (
	"context"
	"log"
	"time"

	ck "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"slib.uz/src/core/domain/ports/messaging"
)

type KafkaProducer struct {
	p      *ck.Producer
	admin  *ck.AdminClient
	ctx    context.Context
	cancel context.CancelFunc
}

// @inject
func NewKafkaProducer(p *ck.Producer, admin *ck.AdminClient) messaging.Producer {

	ctx, cancel := context.WithCancel(context.Background())

	producer := &KafkaProducer{p: p, admin: admin, ctx: ctx, cancel: cancel}
	go producer.deliveryHandler()

	return producer
}

func (this *KafkaProducer) PublishSync(topic, key string, value []byte) error {
	deliveryChan := make(chan ck.Event)

	err := this.p.Produce(&ck.Message{
		TopicPartition: ck.TopicPartition{Topic: &topic, Partition: ck.PartitionAny},
		Key:            []byte(key),
		Value:          value,
		// Headers:        hdrs,
	}, deliveryChan)

	if err != nil {
		log.Printf("❌ Kafka produce error: %v", err)
		return err
	}

	e := <-deliveryChan
	m := e.(*ck.Message)

	if m.TopicPartition.Error != nil {
		log.Printf("❌ Failed to deliver message to topic %s: %v", topic, m.TopicPartition.Error)
		close(deliveryChan)
		return m.TopicPartition.Error
	}

	close(deliveryChan)
	return nil
}

func (this *KafkaProducer) Publish(topic, key string, value []byte, headers map[string]string) error {

	var hdrs []ck.Header
	for k, v := range headers {
		hdrs = append(hdrs, ck.Header{Key: k, Value: []byte(v)})
	}

	return this.p.Produce(&ck.Message{
		TopicPartition: ck.TopicPartition{Topic: &topic, Partition: ck.PartitionAny},
		Key:            []byte(key),
		Value:          value,
		Headers:        hdrs,
	}, nil)
}

func (this *KafkaProducer) deliveryHandler() {
	for {
		select {
		case <-this.ctx.Done():
			return
		case ev, ok := <-this.p.Events():
			if !ok {
				return
			}
			switch e := ev.(type) {
			case *ck.Message:
				if e.TopicPartition.Error != nil {
					log.Printf("❌ Failed to deliver DLQ: %v", e.TopicPartition.Error)
				} else {
					log.Printf("✅ DLQ delivered to partition and offset %v   %v", e.TopicPartition.Partition, e.TopicPartition.Offset)
				}
			case ck.Error:
				log.Printf("⚠️ Kafka producer error: %v", e)
			}
		}
	}
}

func (this *KafkaProducer) EnsureTopic(topic string, partitions int, replicationFactor int) error {
	spec := ck.TopicSpecification{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}

	ctx, cancel := context.WithTimeout(this.ctx, 10*time.Second)
	defer cancel()

	results, err := this.admin.CreateTopics(ctx, []ck.TopicSpecification{spec})
	if err != nil {
		return err
	}

	if len(results) == 1 {
		if results[0].Error.Code() == ck.ErrTopicAlreadyExists {
			return nil
		}
		if results[0].Error.Code() != ck.ErrNoError {
			return results[0].Error
		}
	}

	return nil
}

func (this *KafkaProducer) Close() {
	this.cancel()
	this.p.Flush(10_000) // 10s kutadi
	this.admin.Close()
	this.p.Close()
}
