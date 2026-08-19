package session

type Atomic interface {
	Transaction(fn func(tx Tx) error) error
}
