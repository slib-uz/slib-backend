package entity

import "time"

type MessageEntity struct {
	Key       []byte
	Value     []byte
	Timestamp time.Time
	Headers   map[string][]byte
	Topic     string // channel or topic name
}

func NewMessageEntity(key []byte, value []byte, timestamp time.Time, headers map[string][]byte, topic string) MessageEntity {
	return MessageEntity{Key: key, Value: value, Timestamp: timestamp, Headers: headers, Topic: topic}
}
