package gateway

type LiteracyGateway interface {
	SpellCheck(file []byte) ([]byte, error)
}
