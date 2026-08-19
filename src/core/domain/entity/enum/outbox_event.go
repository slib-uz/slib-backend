package enum

type OutboxEvent string

const (
	OutboxEventUseServiceBalance     OutboxEvent = "USE_SERVICE_BALANCE"
	OutboxEventReverseServiceBalance OutboxEvent = "REVERSE_SERVICE_BALANCE"
)

func (this OutboxEvent) String() string {
	return string(this)
}
