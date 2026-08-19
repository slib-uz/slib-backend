package messaging

type Producer interface {
	Publish(topic, key string, value []byte, headers map[string]string) error
	PublishSync(topic, key string, value []byte) error
}
