package producer

import (
	"fmt"

	ck "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// @inject
func NewProducerAdmin(p *ck.Producer) *ck.AdminClient {
	admin, err := ck.NewAdminClientFromProducer(p)
	if err != nil {
		p.Close()
		fmt.Println("[NewProducerAdmin] Failed to create admin:", err)
	}
	return admin
}
