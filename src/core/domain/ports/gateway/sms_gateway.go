package gateway

type SmsGateway interface {
	Send(phone string, message string) error
}
