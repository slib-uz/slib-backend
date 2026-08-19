package producer

import (
	"fmt"

	ck "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"slib.uz/src/infrastructure/config"
)

// @inject
func NewConfluentProducer(env *config.Config) *ck.Producer {
	cf := &ck.ConfigMap{
		"bootstrap.servers":   env.KafkaBrokers,
		"enable.idempotence":  true,
		"acks":                "all",
		"retries":             5,
		"linger.ms":           10,
		"batch.num.messages":  10000,
		"go.delivery.reports": true,
	}

	_ = cf.SetKey("security.protocol", "SASL_PLAINTEXT")
	_ = cf.SetKey("sasl.mechanisms", "PLAIN")
	_ = cf.SetKey("sasl.username", env.KafkaUsername)
	_ = cf.SetKey("sasl.password", env.KafkaPassword)

	p, err := ck.NewProducer(cf)
	if err != nil {
		fmt.Println("Failed to create producer:", err)
	}

	return p
}
